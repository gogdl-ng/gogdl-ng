package app

import (
	"log"
	"net/http"

	"github.com/LegendaryB/gogdl-ng/app/controllers"
	"github.com/LegendaryB/gogdl-ng/app/middlewares"
	"github.com/LegendaryB/gogdl-ng/app/persistence"
	"github.com/gorilla/mux"
)

type App struct {
	Router    *mux.Router
	DbContext *persistence.DbContext
}

func (app *App) Run() {
	dbContext, err := persistence.NewDbContext()

	if err != nil {
		log.Fatal("Failed to create database context!")
	}

	app.Router = createRouter()
	app.DbContext = dbContext

	app.Router.Use(middlewares.JSONMiddleware)
	controllers.AddRoutes(app.Router, app.DbContext.Tasks)

	log.Fatal(http.ListenAndServe(":3200", app.Router))
}

func createRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	return router.PathPrefix("/api/v1").Subrouter()
}
