package providers

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
)

type CouchbaseServiceProvider struct {
}

func (receiver *CouchbaseServiceProvider) Register(app foundation.Application) {

	app.Singleton("couchbase", func(app foundation.Application) (any, error) {
		config := facades.Config()
		connectionString := fmt.Sprintf("%v", config.Env("COUCHBASE_CONNECTION_STRING", "localhost"))
		bucketName := fmt.Sprintf("%v", config.Env("COUCHBASE_BUCKET_NAME", "paygate"))
		username := fmt.Sprintf("%v", config.Env("COUCHBASE_USERNAME", "admax"))
		password := fmt.Sprintf("%v", config.Env("COUCHBASE_PASSWORD", "fcD!1234"))

		cluster, err := gocb.Connect(connectionString, gocb.ClusterOptions{
			Authenticator: gocb.PasswordAuthenticator{
				Username: username,
				Password: password,
			},
		})
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		bucket := cluster.Bucket(bucketName)

		err = bucket.WaitUntilReady(5*time.Second, nil)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		return bucket, nil
	})

}

func (receiver *CouchbaseServiceProvider) Boot(app foundation.Application) {

}
