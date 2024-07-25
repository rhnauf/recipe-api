package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/rhnauf/recipe-api/internal/api"
	"log"
	"os"
	"os/signal"
	"time"
)

type App struct {
	port string
}

func (a *App) initConfiguration() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error reading env files =>", err)
	}
	a.port = os.Getenv("APP_PORT")
}

func (a *App) runWebServer() {
	//_, dbDispose := db.NewDatabase()
	//defer dbDispose()

	handler := api.NewAPI()
	srv := handler.Server(a.port)

	go func() { _ = srv.ListenAndServe() }()

	log.Println("STARTED API AT PORT", a.port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	srv.Shutdown(ctx)

	log.Println("SHUT DOWN GRACEFULLY")
}

func main() {
	a := App{}

	a.initConfiguration()

	a.runWebServer()
}
