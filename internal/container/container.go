package container

import (
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	authHandler "github.com/EduGoGroup/edugo-api-admin-new/internal/auth/handler"
	authService "github.com/EduGoGroup/edugo-api-admin-new/internal/auth/service"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/config"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/infrastructure/http/handler"
	pgRepo "github.com/EduGoGroup/edugo-api-admin-new/internal/infrastructure/persistence/postgres/repository"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
	"gorm.io/gorm"
)

// Container is the dependency injection container
type Container struct {
	DB         *gorm.DB
	Logger     logger.Logger
	JWTManager *auth.JWTManager

	// Auth
	TokenService  *authService.TokenService
	AuthService   authService.AuthService
	AuthHandler   *authHandler.AuthHandler
	VerifyHandler *authHandler.VerifyHandler

	// Handlers
	SchoolHandler       *handler.SchoolHandler
	AcademicUnitHandler *handler.AcademicUnitHandler
	MembershipHandler   *handler.MembershipHandler
	RoleHandler         *handler.RoleHandler
	ResourceHandler     *handler.ResourceHandler
	MenuHandler         *handler.MenuHandler
	PermissionHandler   *handler.PermissionHandler
	SubjectHandler      *handler.SubjectHandler
	GuardianHandler     *handler.GuardianHandler
	ScreenConfigHandler *handler.ScreenConfigHandler
	UserHandler         *handler.UserHandler
	StatsHandler        *handler.StatsHandler
	MaterialHandler     *handler.MaterialHandler
}

// NewContainer creates a new container and initializes all dependencies
func NewContainer(db *gorm.DB, log logger.Logger, cfg *config.Config) *Container {
	c := &Container{
		DB:         db,
		Logger:     log,
		JWTManager: auth.NewJWTManager(cfg.Auth.JWT.Secret, cfg.Auth.JWT.Issuer),
	}

	// Repositories
	userRepo := pgRepo.NewPostgresUserRepository(db)
	schoolRepo := pgRepo.NewPostgresSchoolRepository(db)
	unitRepo := pgRepo.NewPostgresAcademicUnitRepository(db)
	membershipRepo := pgRepo.NewPostgresMembershipRepository(db)
	roleRepo := pgRepo.NewPostgresRoleRepository(db)
	permissionRepo := pgRepo.NewPostgresPermissionRepository(db)
	userRoleRepo := pgRepo.NewPostgresUserRoleRepository(db)
	resourceRepo := pgRepo.NewPostgresResourceRepository(db)
	subjectRepo := pgRepo.NewPostgresSubjectRepository(db)
	guardianRepo := pgRepo.NewPostgresGuardianRepository(db)
	screenTemplateRepo := pgRepo.NewPostgresScreenTemplateRepository(db)
	screenInstanceRepo := pgRepo.NewPostgresScreenInstanceRepository(db)
	resourceScreenRepo := pgRepo.NewPostgresResourceScreenRepository(db)
	statsRepo := pgRepo.NewPostgresStatsRepository(db)
	materialRepo := pgRepo.NewPostgresMaterialRepository(db)

	// Auth
	c.TokenService = authService.NewTokenService(c.JWTManager, cfg.Auth.JWT.AccessTokenDuration, cfg.Auth.JWT.RefreshTokenDuration)
	c.AuthService = authService.NewAuthService(userRepo, userRoleRepo, roleRepo, membershipRepo, schoolRepo, c.TokenService, log)
	c.AuthHandler = authHandler.NewAuthHandler(c.AuthService, log)
	c.VerifyHandler = authHandler.NewVerifyHandler(c.TokenService)

	// Services
	schoolService := service.NewSchoolService(schoolRepo, log, cfg.Defaults.School)
	unitService := service.NewAcademicUnitService(unitRepo, schoolRepo, log)
	membershipService := service.NewMembershipService(membershipRepo, log)
	roleService := service.NewRoleService(roleRepo, permissionRepo, userRoleRepo, log)
	resourceService := service.NewResourceService(resourceRepo, log)
	menuService := service.NewMenuService(resourceRepo, resourceScreenRepo, log)
	permissionService := service.NewPermissionService(permissionRepo, log)
	subjectService := service.NewSubjectService(subjectRepo, log)
	guardianService := service.NewGuardianService(guardianRepo, log)
	screenConfigService := service.NewScreenConfigService(screenTemplateRepo, screenInstanceRepo, resourceScreenRepo, log)
	userService := service.NewUserService(userRepo, log)
	statsService := service.NewStatsService(statsRepo, log)
	materialService := service.NewMaterialService(materialRepo, log)

	// Handlers
	c.SchoolHandler = handler.NewSchoolHandler(schoolService, log)
	c.AcademicUnitHandler = handler.NewAcademicUnitHandler(unitService, log)
	c.MembershipHandler = handler.NewMembershipHandler(membershipService, log)
	c.RoleHandler = handler.NewRoleHandler(roleService, log)
	c.ResourceHandler = handler.NewResourceHandler(resourceService, log)
	c.MenuHandler = handler.NewMenuHandler(menuService, log)
	c.PermissionHandler = handler.NewPermissionHandler(permissionService, log)
	c.SubjectHandler = handler.NewSubjectHandler(subjectService, log)
	c.GuardianHandler = handler.NewGuardianHandler(guardianService, log)
	c.ScreenConfigHandler = handler.NewScreenConfigHandler(screenConfigService, log)
	c.UserHandler = handler.NewUserHandler(userService, log)
	c.StatsHandler = handler.NewStatsHandler(statsService, log)
	c.MaterialHandler = handler.NewMaterialHandler(materialService, log)

	return c
}

// Close releases container resources
func (c *Container) Close() error {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
