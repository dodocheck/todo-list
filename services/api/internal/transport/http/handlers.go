package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/app"
	"github.com/dodocheck/go-pet-project-1/shared/contracts"
)

type HttpHandlers struct {
	service     *app.Service
	closeServer func() error
}

func NewHttpHandlers(service *app.Service) *HttpHandlers {
	return &HttpHandlers{
		service:     service,
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

	taskImportData := contracts.TaskImportData{
		Title: taskDTO.Title,
		Text:  taskDTO.Text}

	ctx := r.Context()
	createdTask, err := h.service.AddTask(ctx, taskImportData)
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
	ctx := r.Context()
	tasks, err := h.service.ListAllTasks(ctx)
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

	ctx := r.Context()
	if err := h.service.RemoveTask(ctx, idDTO.Id); err != nil {
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

	ctx := r.Context()
	updatedTask, err := h.service.MarkTaskFinished(ctx, idDTO.Id)
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
