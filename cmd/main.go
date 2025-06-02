package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"final_project/pkg/utils"

	"github.com/joho/godotenv"
)

type Data struct {
	Message string `json:"message"`
}

type logData struct {
	url    string
	method string
	ip     string
}

var (
	logMutex sync.Mutex
	logChan  chan logData
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Error loading .env file: %v\n", err)
	}

	logChan = make(chan logData, 100)

	go logProcessor()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/api/data", apiDataHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s\n", port)
	}
	log.Printf("Starting server on port %s\n", port)
	err = http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	tmpl, err := template.ParseFiles("templates/home.html")

	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Println("Error loading template:", err)
		return
	}

	data := map[string]string{"Message": "Hello, this is the homepage!"}
	err = tmpl.Execute(w, data)

	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
	}

}

func timeHandler(w http.ResponseWriter, r *http.Request) {

	logRequest(r)

	currentTime := utils.GetCurrentTime()

	fmt.Fprintf(w, "Current time: %s", currentTime)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	name := r.URL.Query().Get("name")

	if name == "" {
		name = "Guest"
	}

	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := map[string]string{"Name": name}
	err = tmpl.Execute(w, data)

	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func apiDataHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	data := Data{Message: "This is JSON data from the API"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

func logRequest(r *http.Request) {
	ip := r.RemoteAddr

	logChan <- logData{
		url:    r.URL.Path,
		method: r.Method,
		ip:     ip,
	}
}

func logProcessor() {
	logFile, err := os.OpenFile("request.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	defer logFile.Close()

	for logEntry := range logChan {
		logMutex.Lock()
		logMessage := fmt.Sprintf("URL: %s, Method: %s, IP: %s, Goroutines: %d\n", logEntry.url, logEntry.method, logEntry.ip, runtime.NumGoroutine())
		_, err := logFile.WriteString(logMessage)

		if err != nil {

			log.Printf("Failed to write to log file: %v", err)

		}
		logMutex.Unlock()
		log.Println(logMessage)

	}

}

func GetProjectDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		return "", fmt.Errorf("не удалось получить информацию о вызывающем файле")
	}

	dir := filepath.Dir(filename)
	projectDir := filepath.Dir(dir)
	return projectDir, nil
}
