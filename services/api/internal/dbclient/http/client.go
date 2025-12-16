package dbhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dodocheck/go-pet-project-1/services/api/pb"
)

type DBClient struct {
	dbUrl      string
	httpClient *http.Client
}

func NewDBClient(dbUrl string) *DBClient {
	return &DBClient{
		dbUrl: dbUrl,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		}}
}

func (c *DBClient) AddTask(ctx context.Context, task pb.TaskImportData) (pb.TaskExportData, error) {
	b, err := json.Marshal(task)
	if err != nil {
		return pb.TaskExportData{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.dbUrl+"/tasks", bytes.NewReader(b))
	if err != nil {
		return pb.TaskExportData{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return pb.TaskExportData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return pb.TaskExportData{}, errors.New("db-service returned unexpected status " + resp.Status)
	}

	var createdTask pb.TaskExportData
	if err := json.NewDecoder(resp.Body).Decode(&createdTask); err != nil {
		return pb.TaskExportData{}, err
	}

	return createdTask, nil
}

func (c *DBClient) RemoveTask(ctx context.Context, id int) error {
	idDTO := struct {
		Id int `json:"id"`
	}{
		Id: id}

	b, err := json.Marshal(&idDTO)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.dbUrl+"/tasks", bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errors.New("db-service returned unexpected status " + resp.Status)
	}

	return nil
}

func (c *DBClient) ListAllTasks(ctx context.Context) ([]pb.TaskExportData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.dbUrl+"/tasks", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("db-service returned unexpected status " + resp.Status)
	}

	var tasks []pb.TaskExportData
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *DBClient) MarkTaskFinished(ctx context.Context, id int) (pb.TaskExportData, error) {
	idDTO := struct {
		Id int `json:"id"`
	}{Id: id}

	b, err := json.Marshal(&idDTO)
	if err != nil {
		return pb.TaskExportData{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.dbUrl+"/tasks", bytes.NewReader(b))
	if err != nil {
		return pb.TaskExportData{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return pb.TaskExportData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return pb.TaskExportData{}, errors.New("db-service returned unexpected status " + resp.Status)
	}

	var updatedTask pb.TaskExportData
	if err := json.NewDecoder(resp.Body).Decode(&updatedTask); err != nil {
		return pb.TaskExportData{}, err
	}

	return updatedTask, nil
}
