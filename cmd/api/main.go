package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // import postgres driver
	"github.com/nasermirzaei89/api/internal/repositories/postgres"
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

func main() {
	l := log.New(os.Stdout, fmt.Sprintln(), 0)

	signKey := env.MustGetString("API_SIGN_KEY")
	verificationKey := env.MustGetString("API_VERIFICATION_KEY")

	db := postgresDB()

	userRepo := postgres.NewUserRepository(db)

	userSvc := user.NewService(userRepo, []byte(signKey), []byte(verificationKey))

	h := http.NewHandler(l, userSvc)

	err := gohttp.ListenAndServe(env.GetString("API_ADDRESS", ":80"), h)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "error on listen and serve http"))
	}
}
