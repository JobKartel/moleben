package assembly

import (
	"moleben/config"
	"moleben/controller"
	"moleben/repository"
	"moleben/router"
	"moleben/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct{ Router *router.Router }

func Build(cfg config.Config, pool *pgxpool.Pool) *App {
	repo := repository.NewPostgresRepo(pool)

	client := repository.NewClient(cfg.BaseURL, cfg.APIKey, cfg.Model, cfg.Provider, cfg.AppReferer, cfg.AppTitle)

	svc := service.NewChatService(repo, client, cfg.SystemPrompt)

	ctrl := controller.NewChatController(svc)

	r := router.New(ctrl)

	return &App{Router: r}
}
