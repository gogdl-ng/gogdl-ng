package controller

import (
	"net/http"
	"strconv"

	"github.com/LegendaryB/gogdl-ng/model/task"
	"github.com/LegendaryB/gogdl-ng/models"
	"github.com/gorilla/mux"
	"github.com/qkgo/yin"
)

type TaskPostBody struct {
	FolderId   string `json:"folderId"`
	FolderName string `json:"folderName"`
}

func GetAllTasks(store task.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _ := yin.Event(w, r)
		tasks, err := store.GetAll()

		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}

		if len(tasks) <= 0 {
			res.SendStatus(http.StatusNoContent)
			return
		}

		res.SendJSON(tasks)
	}
}

func GetTask(store task.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _ := yin.Event(w, r)
		param := mux.Vars(r)["id"]

		id, err := strconv.ParseInt(param, 10, 64)

		if err != nil || id <= 0 {
			res.SendStatus(http.StatusBadRequest)
			return
		}

		t, err := store.Get(id)

		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}

		if t.ID <= 0 {
			res.SendStatus(http.StatusNotFound)
			return
		}

		res.SendJSON(t)
	}
}

func CreateTask(store task.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, req := yin.Event(w, r)
		body := TaskPostBody{}
		req.BindBody(&body)

		insert := models.Task{
			FolderId:   body.FolderId,
			FolderName: body.FolderName,
		}

		t, err := store.Create(insert)

		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}

		res.SendJSON(t)
	}
}

func DeleteTask(store task.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _ := yin.Event(w, r)
		param := mux.Vars(r)["id"]

		id, err := strconv.ParseInt(param, 10, 64)

		if err != nil || id <= 0 {
			res.SendStatus(http.StatusBadRequest)
			return
		}

		t, err := store.Get(id)

		if err != nil {
			res.SendStatus(http.StatusNotFound)
			return
		}

		if t.Status == models.Processing.String() {
			res.SendStatus(http.StatusMethodNotAllowed)
			return
		}

		err = store.Delete(id)

		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}

		res.SendStatus(http.StatusOK)
	}
}
