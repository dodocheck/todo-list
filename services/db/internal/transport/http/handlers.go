package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dodocheck/go-pet-project-1/pb"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/app"
)

type HttpHandlers struct {
	controller  app.Controller
	closeServer func() error
}

func NewHttpHandlers(controller app.Controller) *HttpHandlers {
	return &HttpHandlers{
		controller:  controller,
		closeServer: nil}
}

func (h *HttpHandlers) SetCloseServerFunc(f func() error) {
	h.closeServer = f
}

/*
pattern: /tasks
method: POST
info: JSON in HTTP request body

success:
  - status code: 201 Created
  - response body: JSON represented created data

failure:
  - status code: 400, 404, 500
  - response body: JSON with error + time
*/
func (h *HttpHandlers) handleAddTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO

	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errorDTO := NewErrorDTO(err.Error())
		http.Error(w, errorDTO.ToString(), http.StatusBadRequest)
		return
	}

	taskImportData := pb.TaskImportData{
		Title: taskDTO.Title,
		Text:  taskDTO.Text}

	createdTask, err := h.controller.AddTask(taskImportData)
	if err != nil {
		errorDTO := NewErrorDTO(err.Error())
		http.Error(w, errorDTO.ToString(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	b, err := json.MarshalIndent(createdTask, "", "    ")
	if err != nil {
		errorDTO := NewErrorDTO(err.Error())
		http.Error(w, errorDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		log.Println("Failed to send http answer:", err)
		return
	}
}

/*
pattern: /tasks
method: GET
info: -

success:
  - status code: 200 Ok
  - response body: JSON represented found data

failure:
  - status code: 500
  - response body: JSON with error + time
*/
func (h *HttpHandlers) handleListAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.controller.ListAllTasks()
	if err != nil {
		errorDTO := NewErrorDTO(err.Error())
		http.Error(w, errorDTO.ToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		errorDTO := NewErrorDTO(err.Error())
		http.Error(w, errorDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		log.Println("Failed to send http answer:", err)
		return
	}
}

/*
pattern: /tasks
method: DELETE
info: JSON in HTTP request body

success:
  - status code: 204 No Content
  - response body: -

failure:
  - status code: 400, 404, 429, 500
  - response body: JSON with error + time
*/
func (h *HttpHandlers) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	var idDTO struct {
		Id int
	}
	if err := json.NewDecoder(r.Body).Decode(&idDTO); err != nil {
		errorDTO := NewErrorDTO(err.Error())
		http.Error(w, errorDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := h.controller.DeleteTask(idDTO.Id); err != nil {
		errorDTO := NewErrorDTO(err.Error())
		http.Error(w, errorDTO.ToString(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
pattern: /tasks
method: PATCH
info: JSON in HTTP request body

success:
  - status code: 200 Ok
  - response body: JSON represented updated data

failure:
  - status code: 400, 404, 429, 500
  - response body: JSON with error + time
*/
func (h *HttpHandlers) handleFinishTask(w http.ResponseWriter, r *http.Request) {
	var idDTO struct {
		Id int
	}
	if err := json.NewDecoder(r.Body).Decode(&idDTO); err != nil {
		errorDTO := NewErrorDTO(err.Error())
		http.Error(w, errorDTO.ToString(), http.StatusBadRequest)
		return
	}

	updatedTask, err := h.controller.MarkTaskFinished(idDTO.Id)
	if err != nil {
		errorDTO := NewErrorDTO(err.Error())
		http.Error(w, errorDTO.ToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(updatedTask, "", "    ")
	if err != nil {
		errorDTO := NewErrorDTO(err.Error())
		http.Error(w, errorDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		log.Println("Failed to send http answer:", err)
		return
	}
}
