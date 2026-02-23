package container

import (
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/client"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/config"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/infrastructure/http/handler"
	pgRepo "github.com/EduGoGroup/edugo-api-admin-new/internal/infrastructure/persistence/postgres/repository"
	"github.com/EduGoGroup/edugo-shared/logger"
	sharedrepopg "github.com/EduGoGroup/edugo-shared/repository/postgres"
	"gorm.io/gorm"
)

// Container is the dependency injection container
type Container struct {
	DB     *gorm.DB
	Logger logger.Logger

	// Clients
	AuthClient *client.AuthClient
	IAMClient  *client.IAMClient

	// Handlers
	SchoolHandler       *handler.SchoolHandler
	AcademicUnitHandler *handler.AcademicUnitHandler
	MembershipHandler   *handler.MembershipHandler
	SubjectHandler      *handler.SubjectHandler
	GuardianHandler     *handler.GuardianHandler
	UserHandler         *handler.UserHandler
	StatsHandler        *handler.StatsHandler
	MaterialHandler     *handler.MaterialHandler
	HealthHandler       *handler.HealthHandler
}

// NewContainer creates a new container and initializes all dependencies
func NewContainer(db *gorm.DB, log logger.Logger, cfg *config.Config) *Container {
	c := &Container{
		DB:     db,
		Logger: log,
	}

	// Auth Client (local JWT + optional remote fallback via IAM Platform)
	c.AuthClient = client.NewAuthClient(client.AuthClientConfig{
		JWTSecret:       cfg.Auth.JWT.Secret,
		JWTIssuer:       cfg.Auth.JWT.Issuer,
		BaseURL:         cfg.Auth.APIIamPlatform.BaseURL,
		Timeout:         cfg.Auth.APIIamPlatform.Timeout,
		RemoteEnabled:   cfg.Auth.APIIamPlatform.RemoteEnabled,
		FallbackEnabled: cfg.Auth.APIIamPlatform.FallbackEnabled,
		CacheTTL:        cfg.Auth.APIIamPlatform.CacheTTL,
		CacheEnabled:    cfg.Auth.APIIamPlatform.CacheEnabled,
	})

	// IAM Client (for role operations proxy)
	c.IAMClient = client.NewIAMClient(client.IAMClientConfig{
		BaseURL: cfg.Auth.APIIamPlatform.BaseURL,
		Timeout: cfg.Auth.APIIamPlatform.Timeout,
	})

	// Shared repositories (from edugo-shared/repository)
	schoolRepo := sharedrepopg.NewPostgresSchoolRepository(db)
	userRepo := sharedrepopg.NewPostgresUserRepository(db)
	membershipRepo := sharedrepopg.NewPostgresMembershipRepository(db)

	// Local repositories
	unitRepo := pgRepo.NewPostgresAcademicUnitRepository(db)
	subjectRepo := pgRepo.NewPostgresSubjectRepository(db)
	guardianRepo := pgRepo.NewPostgresGuardianRepository(db)
	statsRepo := pgRepo.NewPostgresStatsRepository(db)
	materialRepo := pgRepo.NewPostgresMaterialRepository(db)

	// Services
	schoolService := service.NewSchoolService(schoolRepo, log, cfg.Defaults.School)
	unitService := service.NewAcademicUnitService(unitRepo, schoolRepo, log)
	membershipService := service.NewMembershipService(membershipRepo, log)
	subjectService := service.NewSubjectService(subjectRepo, log)
	guardianService := service.NewGuardianService(guardianRepo, log)
	userService := service.NewUserService(userRepo, log)
	statsService := service.NewStatsService(statsRepo, log)
	materialService := service.NewMaterialService(materialRepo, log)

	// Handlers
	c.SchoolHandler = handler.NewSchoolHandler(schoolService, log)
	c.AcademicUnitHandler = handler.NewAcademicUnitHandler(unitService, log)
	c.MembershipHandler = handler.NewMembershipHandler(membershipService, log)
	c.SubjectHandler = handler.NewSubjectHandler(subjectService, log)
	c.GuardianHandler = handler.NewGuardianHandler(guardianService, log)
	c.UserHandler = handler.NewUserHandler(userService, log)
	c.StatsHandler = handler.NewStatsHandler(statsService, log)
	c.MaterialHandler = handler.NewMaterialHandler(materialService, log)
	c.HealthHandler = handler.NewHealthHandler(db, "dev")

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
