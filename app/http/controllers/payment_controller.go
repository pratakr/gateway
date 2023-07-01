package controllers

import (
	"fmt"
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
	config := facades.Config()
	connectionString := fmt.Sprintf("%v", config.Env("COUCHBASE_CONNECTION_STRING", "localhost"))
	bucketName := fmt.Sprintf("%v", config.Env("COUCHBASE_BUCKET_NAME", "paygate"))
	username := fmt.Sprintf("%v", config.Env("COUCHBASE_USERNAME", "admax"))
	password := fmt.Sprintf("%v", config.Env("COUCHBASE_PASSWORD", "fcD!1234"))

	type RequestInput struct {
		Product    string  `form:"product" json:"product"`
		Amount     float64 `form:"amount" json:"amount"`
		Channel    string  `form:"channel" json:"channel"`
		CardName   string  `form:"card_name" json:"card_name"`
		CardNumber string  `form:"card_number" json:"card_number"`
	}
	var requestInput RequestInput
	ctx.Request().Bind(&requestInput)

	cluster, err := gocb.Connect(connectionString, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
		ctx.Response().Json(400, http.Json{"error": "1000", "message": "error"})
	}

	bucket := cluster.Bucket(bucketName)
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatal(err)
		ctx.Response().Json(401, http.Json{"error": "1000", "message": "error"})
	}

	// col := app.Couchbase().bucket.Collection("payments")

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
