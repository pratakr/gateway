package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/models"
)

type UserController struct {
	//Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		//Inject services
	}
}

func (r *UserController) Show(ctx http.Context) {
	ctx.Response().Success().Json(http.Json{
		"Hello": "Goravel",
	})
}

func (r *UserController) Create(ctx http.Context) {
	config := facades.Config()
	connectionString := fmt.Sprintf("%v", config.Env("COUCHBASE_CONNECTION_STRING", "localhost"))
	bucketName := fmt.Sprintf("%v", config.Env("COUCHBASE_BUCKET_NAME", "paygate"))
	username := fmt.Sprintf("%v", config.Env("COUCHBASE_USERNAME", "admax"))
	password := fmt.Sprintf("%v", config.Env("COUCHBASE_PASSWORD", "fcD!1234"))

	type UserRequest struct {
		Name    string `form:"name" json:"name"`
		Email   string `form:"email" json:"email"`
		Address string `form:"address" json:"address"`
	}

	var userRequest UserRequest
	ctx.Request().Bind(&userRequest)

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

	col := bucket.Collection("users")

	_, err = col.Upsert("u:"+userRequest.Name,
		models.User{
			Name:    userRequest.Name,
			Email:   userRequest.Email,
			Address: userRequest.Address,
		}, nil)
	if err != nil {
		log.Fatal(err)
		ctx.Response().Json(402, http.Json{"error": "1000", "message": "error"})
	}

	getResult, err := col.Get("u:"+userRequest.Name, nil)
	if err != nil {
		log.Fatal(err)
		ctx.Response().Json(403, http.Json{"error": "1000", "message": "error"})
	}

	var inUser models.User
	err = getResult.Content(&inUser)
	if err != nil {
		log.Fatal(err)
		ctx.Response().Json(404, http.Json{"error": "1000", "message": "error"})
	}

	ctx.Response().Success().Json(inUser)
}
