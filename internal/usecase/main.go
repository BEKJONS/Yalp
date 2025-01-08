package usecase

import (
	"yalp_ulab/config"
	"yalp_ulab/internal/usecase/repo"
	"yalp_ulab/pkg/logger"
	"yalp_ulab/pkg/postgres"
)

// UseCase -.
type UseCase struct {
	UserRepo    UserRepoI
	SessionRepo SessionRepoI
}

// New -.
func New(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		UserRepo:    repo.NewUserRepo(pg, config, logger),
		SessionRepo: repo.NewSessionRepo(pg, config, logger),
	}
}
