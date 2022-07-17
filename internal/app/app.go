package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SuvorovSergey/gitlab-automerge/internal/config"
	"github.com/SuvorovSergey/gitlab-automerge/internal/entity"
	"github.com/SuvorovSergey/gitlab-automerge/internal/services"
	"github.com/SuvorovSergey/gitlab-automerge/pkg/scheduler"
)

type App struct {
	cfg      *config.Config
	services *services.Services
}

var app App
var projects []entity.Project

func Run(cfg *config.Config) {
	configure(cfg)
	//init
	projects = app.services.Gitlab.ProjectsWithConfig()
	
	//schedule updating projects
	go scheduler.DoEvery(cfg.Scheduler.Update, func() {		
		projects = app.services.Gitlab.ProjectsWithConfig()		
	})

	//shedule merge requests
	go scheduler.DoEvery(cfg.Scheduler.Merge, func() {
		app.services.Gitlab.AcceptAllMergeRequests(projects)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	log.Printf("Caught signal %s. Shutting down...", sig)
}

func configure(cfg *config.Config) {
	app.cfg = cfg
	app.services = services.NewServices(cfg)
}
