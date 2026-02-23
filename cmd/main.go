package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/EduGoGroup/edugo-api-admin-new/docs"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/config"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/container"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/infrastructure/http/handler"
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
// @description Clean administration API for EduGo
// @host localhost:8081
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

	// 3. Initialize logger (using standard log for now; in production use edugo-shared/logger)
	// For a clean start, we use a simple logger adapter
	appLogger := newSimpleLogger()

	// 4. Create dependency container
	c := container.NewContainer(gormDB, appLogger, cfg)
	defer func() { _ = c.Close() }()

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
	healthHandler := handler.NewHealthHandler(gormDB, Version)
	r.GET("/health", healthHandler.Health)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ==================== PUBLIC ROUTES ====================
	v1Public := r.Group("/api/v1")
	{
		v1Public.GET("/health", healthHandler.Health)

		authGroup := v1Public.Group("/auth")
		{
			authGroup.POST("/login", c.AuthHandler.Login)
			authGroup.POST("/refresh", c.AuthHandler.Refresh)
			authGroup.POST("/verify", c.VerifyHandler.VerifyToken)
		}
	}

	// ==================== PROTECTED ROUTES (JWT required) ====================
	v1 := r.Group("/api/v1")
	v1.Use(ginmiddleware.JWTAuthMiddleware(c.JWTManager))
	{
		// Auth (protected)
		v1.POST("/auth/logout", c.AuthHandler.Logout)
		v1.POST("/auth/switch-context", c.AuthHandler.SwitchContext)
		v1.GET("/auth/contexts", c.AuthHandler.GetAvailableContexts)

		// Schools
		schools := v1.Group("/schools")
		{
			schools.POST("", ginmiddleware.RequirePermission(enum.PermissionSchoolsCreate), c.SchoolHandler.CreateSchool)
			schools.GET("", ginmiddleware.RequirePermission(enum.PermissionSchoolsRead), c.SchoolHandler.ListSchools)
			schools.GET("/code/:code", ginmiddleware.RequirePermission(enum.PermissionSchoolsRead), c.SchoolHandler.GetSchoolByCode)

			// Academic Units nested under school
			schools.POST("/:id/units", ginmiddleware.RequirePermission(enum.PermissionUnitsCreate), c.AcademicUnitHandler.CreateUnit)
			schools.GET("/:id/units", ginmiddleware.RequirePermission(enum.PermissionUnitsRead), c.AcademicUnitHandler.ListUnitsBySchool)
			schools.GET("/:id/units/tree", ginmiddleware.RequirePermission(enum.PermissionUnitsRead), c.AcademicUnitHandler.GetUnitTree)
			schools.GET("/:id/units/by-type", ginmiddleware.RequirePermission(enum.PermissionUnitsRead), c.AcademicUnitHandler.ListUnitsByType)

			// School CRUD
			schools.GET("/:id", ginmiddleware.RequirePermission(enum.PermissionSchoolsRead), c.SchoolHandler.GetSchool)
			schools.PUT("/:id", ginmiddleware.RequirePermission(enum.PermissionSchoolsUpdate), c.SchoolHandler.UpdateSchool)
			schools.DELETE("/:id", ginmiddleware.RequirePermission(enum.PermissionSchoolsDelete), c.SchoolHandler.DeleteSchool)
		}

		// Academic Units (standalone)
		units := v1.Group("/units")
		{
			units.GET("/:id", ginmiddleware.RequirePermission(enum.PermissionUnitsRead), c.AcademicUnitHandler.GetUnit)
			units.PUT("/:id", ginmiddleware.RequirePermission(enum.PermissionUnitsUpdate), c.AcademicUnitHandler.UpdateUnit)
			units.DELETE("/:id", ginmiddleware.RequirePermission(enum.PermissionUnitsDelete), c.AcademicUnitHandler.DeleteUnit)
			units.POST("/:id/restore", ginmiddleware.RequirePermission(enum.PermissionUnitsUpdate), c.AcademicUnitHandler.RestoreUnit)
			units.GET("/:id/hierarchy-path", ginmiddleware.RequirePermission(enum.PermissionUnitsRead), c.AcademicUnitHandler.GetHierarchyPath)
		}

		// Memberships
		memberships := v1.Group("/memberships")
		{
			memberships.POST("", c.MembershipHandler.CreateMembership)
			memberships.GET("", c.MembershipHandler.ListMembershipsByUnit)
			memberships.GET("/by-role", c.MembershipHandler.ListMembershipsByRole)
			memberships.GET("/:id", c.MembershipHandler.GetMembership)
			memberships.PUT("/:id", c.MembershipHandler.UpdateMembership)
			memberships.DELETE("/:id", c.MembershipHandler.DeleteMembership)
			memberships.POST("/:id/expire", c.MembershipHandler.ExpireMembership)
		}

		// Users CRUD
		users := v1.Group("/users")
		{
			users.POST("", ginmiddleware.RequirePermission(enum.PermissionUsersUpdate), c.UserHandler.CreateUser)
			users.GET("", ginmiddleware.RequirePermission(enum.PermissionUsersRead), c.UserHandler.ListUsers)
			users.GET("/:user_id", ginmiddleware.RequirePermission(enum.PermissionUsersRead), c.UserHandler.GetUser)
			users.PATCH("/:user_id", ginmiddleware.RequirePermission(enum.PermissionUsersUpdate), c.UserHandler.UpdateUser)
			users.DELETE("/:user_id", ginmiddleware.RequirePermission(enum.PermissionUsersUpdate), c.UserHandler.DeleteUser)

			// User sub-resources
			users.GET("/:user_id/memberships", c.MembershipHandler.ListMembershipsByUser)
			users.GET("/:user_id/roles", ginmiddleware.RequirePermission(enum.PermissionUsersRead), c.RoleHandler.GetUserRoles)
			users.POST("/:user_id/roles", ginmiddleware.RequirePermission(enum.PermissionUsersUpdate), c.RoleHandler.GrantRole)
			users.DELETE("/:user_id/roles/:role_id", ginmiddleware.RequirePermission(enum.PermissionUsersUpdate), c.RoleHandler.RevokeRole)
		}

		// Stats
		stats := v1.Group("/stats")
		{
			stats.GET("/global", ginmiddleware.RequirePermission(enum.PermissionPermissionsMgmtRead), c.StatsHandler.GetGlobalStats)
		}

		// Materials (admin moderation)
		materials := v1.Group("/materials")
		{
			materials.DELETE("/:id", ginmiddleware.RequirePermission(enum.PermissionPermissionsMgmtUpdate), c.MaterialHandler.DeleteMaterial)
		}

		// Menu
		v1.GET("/menu", c.MenuHandler.GetUserMenu)
		menu := v1.Group("/menu")
		{
			menu.GET("/full", ginmiddleware.RequirePermission(enum.PermissionPermissionsMgmtRead), c.MenuHandler.GetFullMenu)
		}

		// Resources
		resources := v1.Group("/resources")
		{
			resources.GET("", ginmiddleware.RequirePermission(enum.PermissionPermissionsMgmtRead), c.ResourceHandler.ListResources)
			resources.GET("/:id", ginmiddleware.RequirePermission(enum.PermissionPermissionsMgmtRead), c.ResourceHandler.GetResource)
			resources.POST("", ginmiddleware.RequirePermission(enum.PermissionPermissionsMgmtUpdate), c.ResourceHandler.CreateResource)
			resources.PUT("/:id", ginmiddleware.RequirePermission(enum.PermissionPermissionsMgmtUpdate), c.ResourceHandler.UpdateResource)
		}

		// Permissions
		permissionsGroup := v1.Group("/permissions")
		{
			permissionsGroup.GET("", ginmiddleware.RequirePermission(enum.PermissionPermissionsMgmtRead), c.PermissionHandler.ListPermissions)
			permissionsGroup.GET("/:id", ginmiddleware.RequirePermission(enum.PermissionPermissionsMgmtRead), c.PermissionHandler.GetPermission)
		}

		// Roles
		roles := v1.Group("/roles")
		{
			roles.GET("", ginmiddleware.RequirePermission(enum.PermissionUsersRead), c.RoleHandler.ListRoles)
			roles.GET("/:id", ginmiddleware.RequirePermission(enum.PermissionUsersRead), c.RoleHandler.GetRole)
			roles.GET("/:id/permissions", ginmiddleware.RequirePermission(enum.PermissionUsersRead), c.RoleHandler.GetRolePermissions)
		}

		// Subjects
		subjects := v1.Group("/subjects")
		{
			subjects.POST("", c.SubjectHandler.CreateSubject)
			subjects.GET("", c.SubjectHandler.ListSubjects)
			subjects.GET("/:id", c.SubjectHandler.GetSubject)
			subjects.PATCH("/:id", c.SubjectHandler.UpdateSubject)
			subjects.DELETE("/:id", c.SubjectHandler.DeleteSubject)
		}

		// Guardian Relations
		guardianRelations := v1.Group("/guardian-relations")
		{
			guardianRelations.POST("", c.GuardianHandler.CreateGuardianRelation)
			guardianRelations.GET("/:id", c.GuardianHandler.GetGuardianRelation)
			guardianRelations.PUT("/:id", c.GuardianHandler.UpdateGuardianRelation)
			guardianRelations.DELETE("/:id", c.GuardianHandler.DeleteGuardianRelation)
		}
		guardians := v1.Group("/guardians")
		{
			guardians.GET("/:guardian_id/relations", c.GuardianHandler.GetGuardianRelations)
		}
		students := v1.Group("/students")
		{
			students.GET("/:student_id/guardians", c.GuardianHandler.GetStudentGuardians)
		}

		// Screen Config
		screenConfig := v1.Group("/screen-config")
		{
			templates := screenConfig.Group("/templates")
			{
				templates.POST("", c.ScreenConfigHandler.CreateTemplate)
				templates.GET("", c.ScreenConfigHandler.ListTemplates)
				templates.GET("/:id", c.ScreenConfigHandler.GetTemplate)
				templates.PUT("/:id", c.ScreenConfigHandler.UpdateTemplate)
				templates.DELETE("/:id", c.ScreenConfigHandler.DeleteTemplate)
			}
			instances := screenConfig.Group("/instances")
			{
				instances.POST("", c.ScreenConfigHandler.CreateInstance)
				instances.GET("", c.ScreenConfigHandler.ListInstances)
				instances.GET("/:id", c.ScreenConfigHandler.GetInstance)
				instances.GET("/key/:key", c.ScreenConfigHandler.GetInstanceByKey)
				instances.PUT("/:id", c.ScreenConfigHandler.UpdateInstance)
				instances.DELETE("/:id", c.ScreenConfigHandler.DeleteInstance)
			}
			resolve := screenConfig.Group("/resolve")
			{
				resolve.GET("/key/:key", c.ScreenConfigHandler.ResolveScreenByKey)
			}
			resourceScreens := screenConfig.Group("/resource-screens")
			{
				resourceScreens.POST("", c.ScreenConfigHandler.LinkScreenToResource)
				resourceScreens.GET("/:resourceId", c.ScreenConfigHandler.GetScreensForResource)
				resourceScreens.DELETE("/:id", c.ScreenConfigHandler.UnlinkScreen)
			}
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
