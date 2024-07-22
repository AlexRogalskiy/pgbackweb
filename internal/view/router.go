package view

import (
	"github.com/eduardolat/pgbackweb/internal/service"
	"github.com/eduardolat/pgbackweb/internal/view/api"
	"github.com/eduardolat/pgbackweb/internal/view/static"
	"github.com/eduardolat/pgbackweb/internal/view/web"
	"github.com/labstack/echo/v4"
)

func MountRouter(app *echo.Echo, servs *service.Service) {
	app.StaticFS("", static.StaticFs)
	api.MountRouter(app.Group("/api"), servs)
	web.MountRouter(app.Group(""), servs)
}