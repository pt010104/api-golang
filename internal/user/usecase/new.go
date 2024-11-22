package usecase

import (
	"github.com/pt010104/api-golang/internal/email"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/log"
)

type implUsecase struct {
	l         log.Logger
	repo      user.Repo
	emailUC   email.UseCase
	redisRepo user.Redis
}

func New(l log.Logger, repo user.Repo, emailUC email.UseCase, redisRepo user.Redis) implUsecase {
	return implUsecase{
		l:         l,
		repo:      repo,
		emailUC:   emailUC,
		redisRepo: redisRepo,
	}
}
