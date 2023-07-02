package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type HealthcheckController struct {
	//Dependent services
}

func NewHealthcheckController() *HealthcheckController {
	return &HealthcheckController{
		//Inject services
	}
}

func (r *HealthcheckController) Index(ctx http.Context) {
	ctx.Response().Json(http.StatusOK, http.Json{
		"status": "ok",
	})
}
