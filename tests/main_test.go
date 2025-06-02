package main

import (
	"encoding/json"
	"final_project/pkg/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func GetProjectDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		return "", fmt.Errorf("не удалось получить информацию о вызывающем файле")
	}

	dir := filepath.Dir(filename)
	projectDir := filepath.Dir(filepath.Dir(dir))
	return projectDir, nil

}

func TestHomeHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	homeHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Hello, this is the homepage!"

	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}

func TestTimeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/time", nil)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	timeHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Current time:"

	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/index?name=TestUser", nil)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	rr := httptest.NewRecorder()
	indexHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Hello, TestUser!"

	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	req, err = http.NewRequest("GET", "/index", nil)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr = httptest.NewRecorder()
	indexHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected = "Hello, Guest!"

	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}

func TestApiDataHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/data", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	apiDataHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")

	if contentType != "application/json" {
		t.Errorf("Handler returned unexpected Content-Type: got %v want %v", contentType, "application/json")
	}

	var data map[string]string
	err = json.NewDecoder(rr.Body).Decode(&data)

	if err != nil {
		t.Fatalf("Could not decode response body: %v", err)
	}

	expectedMessage := "This is JSON data from the API"

	if data["message"] != expectedMessage {
		t.Errorf("Handler returned unexpected message: got %v want %v", data["message"], expectedMessage)
	}
}

func TestGetCurrentTime(t *testing.T) {

	currentTime := utils.GetCurrentTime()
	if currentTime == "" {
		t.Errorf("GetCurrentTime returned an empty string")
	}
	fmt.Println("Current Time:", currentTime)

}

func TestLogRequest(t *testing.T) {

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	logChan := make(chan logData, 1)
	oldLogChan := logChan

	logMutex.Lock()
	logChan = make(chan logData, 1)
	logMutex.Unlock()

	logRequest(req)

	logEntry := <-logChan

	logMutex.Lock()
	logChan = oldLogChan
	logMutex.Unlock()

	if logEntry.url != "/test" {
		t.Errorf("logRequest logged incorrect URL: got %v, want %v", logEntry.url, "/test")
	}

	if logEntry.method != "GET" {
		t.Errorf("logRequest logged incorrect method: got %v, want %v", logEntry.method, "GET")
	}

	if logEntry.ip == "" {
		t.Errorf("logRequest logged empty IP address")
	}

	close(logChan)

}

func TestLogProcessor(t *testing.T) {

	tempFile, err := os.CreateTemp("", "testlog")
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}

	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	logChan := make(chan logData, 1)
	logChan <- logData{
		url:    "/test",
		method: "GET",
		ip:     "127.0.0.1",
	}
	close(logChan)
	logMutex.Lock()
	oldOutput := log.Writer()
	log.SetOutput(tempFile)

	defer log.SetOutput(oldOutput)

	go logProcessor()
	logMutex.Unlock()

	close(logChan)

	content, err := io.ReadFile(tempFile.Name())

	if err != nil {
		t.Fatalf("Could not read temp file: %v", err)
	}
}
