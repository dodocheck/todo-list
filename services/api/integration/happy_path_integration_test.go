//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

type createdTask struct {
	Id         int        `json:"Id"`
	Title      string     `json:"Title"`
	Text       string     `json:"Text"`
	Finished   bool       `json:"Finished"`
	CreatedAt  time.Time  `json:"CreatedAt"`
	FinishedAt *time.Time `json:"FinishedAt"`
}

func TestHappyPath_HTTP_to_gRPC_to_Postgres(t *testing.T) {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		// твой дефолт из deployment/.env: API_SERVICE_EXTERNAL_PORT=9089
		baseURL = "http://localhost:9089"
	}

	client := &http.Client{Timeout: 5 * time.Second}

	// Ждём, пока api-service реально поднимется и начнёт отвечать
	waitForAPI(t, client, baseURL)

	// 1) create
	createReq := map[string]any{
		"title": fmt.Sprintf("it-%d", time.Now().UnixNano()),
		"text":  "from integration test",
	}

	var created createdTask
	doJSON(t, client, "POST", baseURL+"/create", createReq, &created, http.StatusCreated)

	if created.Id <= 0 {
		t.Fatalf("expected created.Id > 0, got %v", created.Id)
	}
	if created.Finished {
		t.Fatalf("expected created.Finished=false, got true")
	}

	// 2) list -> должен увидеть созданную задачу
	var list1 []createdTask
	doJSON(t, client, "GET", baseURL+"/list", nil, &list1, http.StatusOK)

	if !containsID(list1, created.Id) {
		t.Fatalf("expected task id=%d in list, got %+v", created.Id, list1)
	}

	// 3) done
	doneReq := map[string]any{"Id": created.Id}
	var done createdTask
	doJSON(t, client, "PUT", baseURL+"/done", doneReq, &done, http.StatusOK)

	if !done.Finished {
		t.Fatalf("expected done.Finished=true, got false")
	}
	if done.FinishedAt == nil {
		t.Fatalf("expected done.FinishedAt != nil")
	}

	// 4) delete
	deleteReq := map[string]any{"Id": created.Id}
	doJSON(t, client, "DELETE", baseURL+"/delete", deleteReq, nil, http.StatusNoContent)

	// 5) list -> задачи уже нет
	var list2 []createdTask
	doJSON(t, client, "GET", baseURL+"/list", nil, &list2, http.StatusOK)

	if containsID(list2, created.Id) {
		t.Fatalf("expected task id=%d NOT in list after delete, got %+v", created.Id, list2)
	}
}

func waitForAPI(t *testing.T, client *http.Client, baseURL string) {
	t.Helper()

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		req, _ := http.NewRequest("GET", baseURL+"/list", nil)
		resp, err := client.Do(req)
		if err == nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("api-service did not become ready at %s", baseURL)
}

func doJSON(t *testing.T, client *http.Client, method, url string, reqBody any, respBody any, wantStatus int) {
	t.Helper()

	var body io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatalf("marshal request: %v", err)
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("%s %s failed: %v", method, url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != wantStatus {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("%s %s: want status=%d got=%d body=%s", method, url, wantStatus, resp.StatusCode, string(b))
	}

	if respBody == nil {
		io.Copy(io.Discard, resp.Body)
		return
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}
	if err := json.Unmarshal(b, respBody); err != nil {
		t.Fatalf("unmarshal response: %v, body=%s", err, string(b))
	}
}

func containsID(list []createdTask, id int) bool {
	for _, t := range list {
		if t.Id == id {
			return true
		}
	}
	return false
}
