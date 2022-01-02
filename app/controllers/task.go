package controllers

import (
	"net/http"
	"strconv"

	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/LegendaryB/gogdl-ng/app/models/task"
	"github.com/LegendaryB/gogdl-ng/app/persistence"
	"github.com/gorilla/mux"
	"github.com/qkgo/yin"
)

type TaskPostBody struct {
	DriveId string `json:"driveId"`
}

func GetTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _ := yin.Event(w, r)
		tasks, err := persistence.Tasks.All()

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

func GetTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _ := yin.Event(w, r)
		param := mux.Vars(r)["id"]

		id, err := strconv.Atoi(param)

		if err != nil || id <= 0 {
			res.SendStatus(http.StatusBadRequest)
			return
		}

		t := &task.Task{}

		if err = persistence.Tasks.Get(id, t); err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}

		res.SendJSON(t)
	}
}

func CreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, req := yin.Event(w, r)
		body := TaskPostBody{}
		req.BindBody(&body)

		folder, err := gdrive.GetFilesByFolderId(body.DriveId)

		if err != nil {
			res.SendStatus(http.StatusNotFound)
			return
		}

		task := &task.Task{
			Name:   folder.Name,
			Status: task.New,
			Files:  folder.Files,
		}

		_, err = persistence.Tasks.Add(task)

		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}

		res.SendJSON(task)
	}
}

func DeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _ := yin.Event(w, r)
		param := mux.Vars(r)["id"]

		id, err := strconv.Atoi(param)

		if err != nil || id <= 0 {
			res.SendStatus(http.StatusBadRequest)
			return
		}

		if err != nil {
			res.SendStatus(http.StatusNotFound)
			return
		}

		err = persistence.Tasks.Delete(id)

		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}

		res.SendStatus(http.StatusOK)
	}
}
