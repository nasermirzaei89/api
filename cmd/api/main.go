package main

import (
	"compress/gzip"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // import postgres driver
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nasermirzaei89/api/internal/repositories/postgres"
	"github.com/nasermirzaei89/api/internal/services/file"
	"github.com/nasermirzaei89/api/internal/services/post"
	"github.com/nasermirzaei89/api/internal/services/user"
	"github.com/nasermirzaei89/api/internal/transport/http"
	"github.com/nasermirzaei89/env"
	"github.com/pkg/errors"
	"log"
	gohttp "net/http"
	"os"
)

func postgresDB() *sql.DB {
	db, err := sql.Open("postgres", env.MustGetString("API_POSTGRES_DSN"))
	if err != nil {
		log.Fatalln(errors.Wrap(err, "error on open sql connection"))
	}

	db.SetMaxIdleConns(env.GetInt("API_POSTGRES_MAX_IDLE_CONNECTIONS", 0))
	db.SetMaxOpenConns(env.GetInt("API_POSTGRES_MAX_OPEN_CONNECTIONS", 0))

	err = db.Ping()
	if err != nil {
		log.Fatalln(errors.Wrap(err, "error on ping sql db"))
	}

	return db
}

func minioClient(ctx context.Context) *minio.Client {
	client, err := minio.New(env.MustGetString("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(env.MustGetString("MINIO_ACCESS_KEY"), env.MustGetString("MINIO_SECRET_KEY"), ""),
		Secure: env.GetBool("MINIO_SECURE", false),
	})
	if err != nil {
		log.Fatalln(errors.Wrap(err, "error create new minio client"))
	}

	// check bucket
	exists, err := client.BucketExists(ctx, env.MustGetString("MINIO_BUCKET"))
	if err != nil {
		panic(errors.Wrap(err, "error on check minio bucket"))
	}

	if !exists {
		err := client.MakeBucket(ctx, env.MustGetString("MINIO_BUCKET"), minio.MakeBucketOptions{})
		if err != nil {
			panic(errors.Wrap(err, "error on make minio bucket"))
		}
	}

	return client
}

func main() {
	// prerequisites
	// logger
	l := log.New(os.Stdout, fmt.Sprintln(), 0)

	// rsa 256 key pair
	signKey := env.MustGetString("API_SIGN_KEY")
	verificationKey := env.MustGetString("API_VERIFICATION_KEY")

	// database
	db := postgresDB()

	// minio
	mc := minioClient(context.Background())

	// repositories
	userRepo := postgres.NewUserRepository(db)
	postRepo := postgres.NewPostRepository(db)

	// services
	userSvc := user.NewService(userRepo, []byte(signKey), []byte(verificationKey))
	postSvc := post.NewService(postRepo)
	fileSvc := file.NewService(mc, env.MustGetString("MINIO_BUCKET"))

	// transport
	h := http.NewHandler(l, userSvc, postSvc, fileSvc,
		http.SetGZipLevel(gzip.BestSpeed),
		http.SetGraphiQL(!env.IsProduction()),
		http.SetGraphQLPlayground(!env.IsProduction()),
	)

	err := gohttp.ListenAndServe(env.GetString("API_ADDRESS", ":80"), h)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "error on listen and serve http"))
	}
}
