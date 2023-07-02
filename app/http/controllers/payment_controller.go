package controllers

import (
	"goravel/app/models"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type PaymentController struct {
	//Dependent services
}

func NewPaymentController() *PaymentController {
	return &PaymentController{
		//Inject services
	}
}

func (r *PaymentController) Index(ctx http.Context) {
}

func (r *PaymentController) Create(ctx http.Context) {

	type RequestInput struct {
		Product    string  `form:"product" json:"product"`
		Amount     float64 `form:"amount" json:"amount"`
		Channel    string  `form:"channel" json:"channel"`
		CardName   string  `form:"card_name" json:"card_name"`
		CardNumber string  `form:"card_number" json:"card_number"`
	}
	var requestInput RequestInput
	ctx.Request().Bind(&requestInput)

	buc, err := facades.App().Make("couchbase")
	if err != nil {
		log.Fatal(err)
		ctx.Response().Json(402, http.Json{"error": "900", "message": "error"})
	}
	var bucket *gocb.Bucket = buc.(*gocb.Bucket)

	col := bucket.Collection("payments")

	id := uuid.New().String()
	log.Println("debug:" + time.Now().String())
	_, err = col.Upsert(id,
		models.Payment{
			UserId:    1,
			Product:   requestInput.Product,
			Amount:    requestInput.Amount,
			Fee:       0,
			Total:     0,
			CreatedAt: carbon.Now().String(),
		}, nil)
	if err != nil {
		log.Println(err)
		ctx.Response().Json(402, http.Json{"error": "1000", "message": "error"})
	}

	getResult, err := col.Get(id, nil)
	if err != nil {
		log.Fatal(err)
		ctx.Response().Json(403, http.Json{"error": "1000", "message": "error"})
	}

	var inDB models.Payment
	err = getResult.Content(&inDB)
	if err != nil {
		log.Fatal(err)
		ctx.Response().Json(404, http.Json{"error": "1000", "message": "error"})
	}

	ctx.Response().Success().Json(inDB)
}
