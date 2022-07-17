package services

import (
	"github.com/SuvorovSergey/gitlab-automerge/internal/config"
	"github.com/SuvorovSergey/gitlab-automerge/internal/services/gitlab"
)

type Services struct {
	Gitlab *gitlab.Gitlab
}

var services Services

func NewServices(cfg *config.Config) *Services {
	services.Gitlab = gitlab.NewGitlab(&cfg.Gitlab)
	return &services
}
