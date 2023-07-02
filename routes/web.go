package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) {
		ctx.Response().Json(http.StatusOK, http.Json{
			"Hello": "Goravel",
		})
	})

	healthcheckController := controllers.NewHealthcheckController()
	userController := controllers.NewUserController()
	paymentController := controllers.NewPaymentController()
	facades.Route().Get("/users/{id}", userController.Show)
	facades.Route().Post("/users", userController.Create)

	facades.Route().Get("/healthcheck", healthcheckController.Index)

	facades.Route().Prefix("/api").Group(func(route route.Route) {
		route.Post("/payments", paymentController.Create)
	})
}
