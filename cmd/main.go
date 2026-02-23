package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/EduGoGroup/edugo-api-admin-new/docs"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/client"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/config"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/container"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/infrastructure/http/middleware"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/EduGoGroup/edugo-shared/logger"
	ginmiddleware "github.com/EduGoGroup/edugo-shared/middleware/gin"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

// @title EduGo API Admin New
// @version 1.0
// @description Administration API for EduGo - schools, academic units, memberships, users, subjects, guardians, stats, materials
// @host localhost:8060
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token. Example: "Bearer eyJhbGci..."
func main() {
	log.Printf("EduGo API Admin New starting... (Version: %s, Build: %s)", Version, BuildTime)

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// 2. Connect to PostgreSQL via GORM
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=auth,iam,academic,content,assessment,ui_config,public",
		cfg.Database.Postgres.Host, cfg.Database.Postgres.Port, cfg.Database.Postgres.User,
		cfg.Database.Postgres.Password, cfg.Database.Postgres.Database, cfg.Database.Postgres.SSLMode)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL via GORM: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Error getting underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.Postgres.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.Postgres.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("Error pinging PostgreSQL: %v", err)
	}
	log.Println("PostgreSQL connected successfully via GORM")

	// 3. Initialize logger
	appLogger := newSimpleLogger()

	// 4. Create dependency container
	cont := container.NewContainer(gormDB, appLogger, cfg)
	defer func() { _ = cont.Close() }()

	// 5. Configure Swagger host dynamically from config
	if cfg.Server.SwaggerHost != "" {
		docs.SwaggerInfo.Host = cfg.Server.SwaggerHost
	} else {
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", cfg.Server.Port)
	}

	// 6. Configure Gin
	r := gin.Default()

	// CORS middleware
	r.Use(middleware.CORSMiddleware(&cfg.CORS))

	// Error handler middleware
	r.Use(middleware.ErrorHandler(appLogger))

	// Health check (root for infra probes + /api/v1 for Swagger)
	r.GET("/health", cont.HealthHandler.Health)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ==================== PUBLIC ROUTES ====================
	v1Public := r.Group("/api/v1")
	{
		v1Public.GET("/health", cont.HealthHandler.Health)
	}

	// ==================== PROTECTED ROUTES (JWT required) ====================
	v1 := r.Group("/api/v1")
	v1.Use(middleware.RemoteAuthMiddleware(cont.AuthClient))
	{
		// Schools
		schools := v1.Group("/schools")
		{
			schools.POST("", ginmiddleware.RequirePermission(enum.PermissionSchoolsCreate), cont.SchoolHandler.CreateSchool)
			schools.GET("", ginmiddleware.RequirePermission(enum.PermissionSchoolsRead), cont.SchoolHandler.ListSchools)
			schools.GET("/code/:code", ginmiddleware.RequirePermission(enum.PermissionSchoolsRead), cont.SchoolHandler.GetSchoolByCode)

			// Academic Units nested under school
			schools.POST("/:id/units", ginmiddleware.RequirePermission(enum.PermissionUnitsCreate), cont.AcademicUnitHandler.CreateUnit)
			schools.GET("/:id/units", ginmiddleware.RequirePermission(enum.PermissionUnitsRead), cont.AcademicUnitHandler.ListUnitsBySchool)
			schools.GET("/:id/units/tree", ginmiddleware.RequirePermission(enum.PermissionUnitsRead), cont.AcademicUnitHandler.GetUnitTree)
			schools.GET("/:id/units/by-type", ginmiddleware.RequirePermission(enum.PermissionUnitsRead), cont.AcademicUnitHandler.ListUnitsByType)

			// School CRUD
			schools.GET("/:id", ginmiddleware.RequirePermission(enum.PermissionSchoolsRead), cont.SchoolHandler.GetSchool)
			schools.PUT("/:id", ginmiddleware.RequirePermission(enum.PermissionSchoolsUpdate), cont.SchoolHandler.UpdateSchool)
			schools.DELETE("/:id", ginmiddleware.RequirePermission(enum.PermissionSchoolsDelete), cont.SchoolHandler.DeleteSchool)
		}

		// Academic Units (standalone)
		units := v1.Group("/units")
		{
			units.GET("/:id", ginmiddleware.RequirePermission(enum.PermissionUnitsRead), cont.AcademicUnitHandler.GetUnit)
			units.PUT("/:id", ginmiddleware.RequirePermission(enum.PermissionUnitsUpdate), cont.AcademicUnitHandler.UpdateUnit)
			units.DELETE("/:id", ginmiddleware.RequirePermission(enum.PermissionUnitsDelete), cont.AcademicUnitHandler.DeleteUnit)
			units.POST("/:id/restore", ginmiddleware.RequirePermission(enum.PermissionUnitsUpdate), cont.AcademicUnitHandler.RestoreUnit)
			units.GET("/:id/hierarchy-path", ginmiddleware.RequirePermission(enum.PermissionUnitsRead), cont.AcademicUnitHandler.GetHierarchyPath)
		}

		// Memberships
		memberships := v1.Group("/memberships")
		{
			memberships.POST("", cont.MembershipHandler.CreateMembership)
			memberships.GET("", cont.MembershipHandler.ListMembershipsByUnit)
			memberships.GET("/by-role", cont.MembershipHandler.ListMembershipsByRole)
			memberships.GET("/:id", cont.MembershipHandler.GetMembership)
			memberships.PUT("/:id", cont.MembershipHandler.UpdateMembership)
			memberships.DELETE("/:id", cont.MembershipHandler.DeleteMembership)
			memberships.POST("/:id/expire", cont.MembershipHandler.ExpireMembership)
		}

		// Users CRUD
		users := v1.Group("/users")
		{
			users.POST("", ginmiddleware.RequirePermission(enum.PermissionUsersUpdate), cont.UserHandler.CreateUser)
			users.GET("", ginmiddleware.RequirePermission(enum.PermissionUsersRead), cont.UserHandler.ListUsers)
			users.GET("/:user_id", ginmiddleware.RequirePermission(enum.PermissionUsersRead), cont.UserHandler.GetUser)
			users.PATCH("/:user_id", ginmiddleware.RequirePermission(enum.PermissionUsersUpdate), cont.UserHandler.UpdateUser)
			users.DELETE("/:user_id", ginmiddleware.RequirePermission(enum.PermissionUsersUpdate), cont.UserHandler.DeleteUser)

			// User sub-resources
			users.GET("/:user_id/memberships", cont.MembershipHandler.ListMembershipsByUser)

			// IAM Proxy routes (delegate to iam-platform)
			users.GET("/:user_id/roles", ginmiddleware.RequirePermission(enum.PermissionUsersRead), iamProxyGetUserRoles(cont))
			users.POST("/:user_id/roles", ginmiddleware.RequirePermission(enum.PermissionUsersUpdate), iamProxyGrantRole(cont))
			users.DELETE("/:user_id/roles/:role_id", ginmiddleware.RequirePermission(enum.PermissionUsersUpdate), iamProxyRevokeRole(cont))
		}

		// Stats
		stats := v1.Group("/stats")
		{
			stats.GET("/global", ginmiddleware.RequirePermission(enum.PermissionPermissionsMgmtRead), cont.StatsHandler.GetGlobalStats)
		}

		// Materials (admin moderation)
		materials := v1.Group("/materials")
		{
			materials.DELETE("/:id", ginmiddleware.RequirePermission(enum.PermissionPermissionsMgmtUpdate), cont.MaterialHandler.DeleteMaterial)
		}

		// Subjects
		subjects := v1.Group("/subjects")
		{
			subjects.POST("", cont.SubjectHandler.CreateSubject)
			subjects.GET("", cont.SubjectHandler.ListSubjects)
			subjects.GET("/:id", cont.SubjectHandler.GetSubject)
			subjects.PATCH("/:id", cont.SubjectHandler.UpdateSubject)
			subjects.DELETE("/:id", cont.SubjectHandler.DeleteSubject)
		}

		// Guardian Relations
		guardianRelations := v1.Group("/guardian-relations")
		{
			guardianRelations.POST("", cont.GuardianHandler.CreateGuardianRelation)
			guardianRelations.GET("/:id", cont.GuardianHandler.GetGuardianRelation)
			guardianRelations.PUT("/:id", cont.GuardianHandler.UpdateGuardianRelation)
			guardianRelations.DELETE("/:id", cont.GuardianHandler.DeleteGuardianRelation)
		}
		guardians := v1.Group("/guardians")
		{
			guardians.GET("/:guardian_id/relations", cont.GuardianHandler.GetGuardianRelations)
		}
		students := v1.Group("/students")
		{
			students.GET("/:student_id/guardians", cont.GuardianHandler.GetStudentGuardians)
		}
	}

	// 7. Start HTTP server with graceful shutdown
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("Server listening on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

// ==================== IAM Proxy Handlers ====================

func iamProxyGetUserRoles(cont *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		roles, err := cont.IAMClient.GetUserRoles(c.Request.Context(), token, userID)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, roles)
	}
}

func iamProxyGrantRole(cont *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		var grantReq client.GrantRoleRequest
		if err := c.ShouldBindJSON(&grantReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		result, err := cont.IAMClient.GrantRole(c.Request.Context(), token, userID, grantReq)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func iamProxyRevokeRole(cont *container.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		roleID := c.Param("role_id")
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if err := cont.IAMClient.RevokeRole(c.Request.Context(), token, userID, roleID); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// simpleLogger adapts standard log to the logger.Logger interface
type simpleLogger struct{}

func newSimpleLogger() *simpleLogger { return &simpleLogger{} }

func (l *simpleLogger) Debug(msg string, keysAndValues ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, keysAndValues)
}
func (l *simpleLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Printf("[INFO] %s %v", msg, keysAndValues)
}
func (l *simpleLogger) Warn(msg string, keysAndValues ...interface{}) {
	log.Printf("[WARN] %s %v", msg, keysAndValues)
}
func (l *simpleLogger) Error(msg string, keysAndValues ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, keysAndValues)
}
func (l *simpleLogger) Fatal(msg string, keysAndValues ...interface{}) {
	log.Fatalf("[FATAL] %s %v", msg, keysAndValues)
}
func (l *simpleLogger) With(_ ...interface{}) logger.Logger {
	return l
}
func (l *simpleLogger) Sync() error { return nil }
