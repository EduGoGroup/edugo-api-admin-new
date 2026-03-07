package container

import (
	"github.com/EduGoGroup/edugo-api-admin-new/internal/application/service"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/client"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/config"
	"github.com/EduGoGroup/edugo-api-admin-new/internal/infrastructure/http/handler"
	pgRepo "github.com/EduGoGroup/edugo-api-admin-new/internal/infrastructure/persistence/postgres/repository"
	"github.com/EduGoGroup/edugo-shared/audit"
	auditpostgres "github.com/EduGoGroup/edugo-shared/audit/postgres"
	"github.com/EduGoGroup/edugo-shared/logger"
	sharedrepopg "github.com/EduGoGroup/edugo-shared/repository"
	"gorm.io/gorm"
)

// Container is the dependency injection container
type Container struct {
	DB     *gorm.DB
	Logger logger.Logger

	// Clients
	AuthClient *client.AuthClient
	IAMClient  *client.IAMClient

	// Audit
	AuditLogger audit.AuditLogger

	// Handlers
	SchoolHandler       *handler.SchoolHandler
	AcademicUnitHandler *handler.AcademicUnitHandler
	MembershipHandler   *handler.MembershipHandler
	SubjectHandler      *handler.SubjectHandler
	GuardianHandler     *handler.GuardianHandler
	UserHandler         *handler.UserHandler
	StatsHandler        *handler.StatsHandler
	MaterialHandler     *handler.MaterialHandler
	ConceptTypeHandler  *handler.ConceptTypeHandler
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
	conceptTypeRepo := pgRepo.NewPostgresConceptTypeRepository(db)
	conceptDefRepo := pgRepo.NewPostgresConceptDefinitionRepository(db)
	schoolConceptRepo := pgRepo.NewPostgresSchoolConceptRepository(db)

	// Audit logger
	auditLogger := auditpostgres.NewPostgresAuditLogger(db, "admin-api")
	c.AuditLogger = auditLogger

	// Services
	schoolService := service.NewSchoolService(schoolRepo, conceptTypeRepo, conceptDefRepo, schoolConceptRepo, log, cfg.Defaults.School, auditLogger)
	unitService := service.NewAcademicUnitService(unitRepo, schoolRepo, log, auditLogger)
	membershipService := service.NewMembershipService(membershipRepo, log, auditLogger)
	subjectService := service.NewSubjectService(subjectRepo, log, auditLogger)
	guardianService := service.NewGuardianService(guardianRepo, log, auditLogger)
	userService := service.NewUserService(userRepo, log, auditLogger)
	statsService := service.NewStatsService(statsRepo, log)
	materialService := service.NewMaterialService(materialRepo, log)
	conceptTypeService := service.NewConceptTypeService(conceptTypeRepo, conceptDefRepo, schoolConceptRepo, log, auditLogger)

	// Handlers
	c.SchoolHandler = handler.NewSchoolHandler(schoolService, log)
	c.AcademicUnitHandler = handler.NewAcademicUnitHandler(unitService, log)
	c.MembershipHandler = handler.NewMembershipHandler(membershipService, log)
	c.SubjectHandler = handler.NewSubjectHandler(subjectService, log)
	c.GuardianHandler = handler.NewGuardianHandler(guardianService, log)
	c.UserHandler = handler.NewUserHandler(userService, log)
	c.StatsHandler = handler.NewStatsHandler(statsService, log)
	c.MaterialHandler = handler.NewMaterialHandler(materialService, log)
	c.ConceptTypeHandler = handler.NewConceptTypeHandler(conceptTypeService, log)
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
