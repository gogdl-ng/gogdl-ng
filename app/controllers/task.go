package controllers

import (
	"net/http"
	"strconv"

	"github.com/LegendaryB/gogdl-ng/app/environment"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/LegendaryB/gogdl-ng/app/models/task"
	"github.com/gorilla/mux"
	"github.com/qkgo/yin"
)

type TaskPostBody struct {
	DriveId string `json:"driveId"`
}

func GetAllTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _ := yin.Event(w, r)
		tasks, err := task.Repository.GetAll()

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

		id, err := strconv.ParseInt(param, 10, 64)

		if err != nil || id <= 0 {
			res.SendStatus(http.StatusBadRequest)
			return
		}

		t, err := task.Repository.Get(id)

		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}

		if t.Id <= 0 {
			res.SendStatus(http.StatusNotFound)
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

		res.SendJSON(yin.H{
			"result": "",
		})

		folder, err := gdrive.Folder(body.DriveId)

		if err != nil {
			res.SendStatus(http.StatusNotFound)
			return
		}

		path, err := environment.CreateTaskDirectory(folder.Name)

		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}

		insert := task.Task{
			DriveId:   folder.Id,
			DriveName: folder.Name,
			LocalPath: path,
		}

		t, err := task.Repository.Create(insert)

		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}

		res.SendJSON(t)
	}
}

func DeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _ := yin.Event(w, r)
		param := mux.Vars(r)["id"]

		id, err := strconv.ParseInt(param, 10, 64)

		if err != nil || id <= 0 {
			res.SendStatus(http.StatusBadRequest)
			return
		}

		if err != nil {
			res.SendStatus(http.StatusNotFound)
			return
		}

		err = task.Repository.Delete(id)

		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}

		res.SendStatus(http.StatusOK)
	}
}
