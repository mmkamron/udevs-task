package main

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/mmkamron/miniTwitter/cmd/api/docs"
	"github.com/mmkamron/miniTwitter/internal/data"
	"github.com/mmkamron/miniTwitter/internal/pkg/config"
	database "github.com/mmkamron/miniTwitter/internal/pkg/db"
)

type application struct {
	config *config.Config
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	cfg := config.Load("./config/local.yaml")
	db, err := database.Load(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = database.MigrateDB(cfg.DBurl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: data.NewModels(db),
	}
	if err := app.serve(); err != nil {
		log.Fatal(err.Error())
	}
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Println("starting the server")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	log.Println("stopped server")
	return nil
}
