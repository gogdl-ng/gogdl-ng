package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/LegendaryB/gogdl-ng/controller"
	"github.com/LegendaryB/gogdl-ng/model/task"
	"github.com/gorilla/mux"
)

func main() {
	db, _ := sql.Open("sqlite3", "./gogdl.db")
	tasks := task.NewStore(db)

	router := mux.NewRouter()

	controller.AddRoutes(router, tasks)
	log.Fatal(http.ListenAndServe(":3200", router))
}
