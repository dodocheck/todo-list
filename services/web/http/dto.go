package web

import (
	"encoding/json"
	"time"
)

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func NewErrorDTO(errorMsg string) ErrorDTO {
	return ErrorDTO{
		Message: errorMsg,
		Time:    time.Now()}
}

func (e *ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

type TaskDTO struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
