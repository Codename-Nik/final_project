package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"final_project/pkg/handlers"
)

var logger *log.Logger

func init() {
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Ошибка открытия файла логов:", err)
	}
	logger = log.New(file, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.Printf("Метод: %s, Путь: %s, IP: %s, User-Agent: %s",
			r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		next(w, r)
		logger.Printf("Запрос обработан за %v", time.Since(start))
	}
}

func main() {
	http.HandleFunc("/", loggingMiddleware(handlers.HomeHandler))
	http.HandleFunc("/time", loggingMiddleware(handlers.TimeHandler))
	http.HandleFunc("/index", loggingMiddleware(handlers.IndexHandler))
	http.HandleFunc("/api", loggingMiddleware(handlers.APIPageHandler))
	http.HandleFunc("/api/data", loggingMiddleware(handlers.APIDataHandler))

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	startMsg := `
╔══════════════════════════════════════════════════════╗
║          СЕРВЕР УСПЕШНО ЗАПУЩЕН                     ║
╠══════════════════════════════════════════════════════╣
║  Порт: :8080                                        ║
║  Адрес: http://localhost:8080                       ║
╠══════════════════════════════════════════════════════╣
║  Маршруты:                                          ║
║  • /           - Главная страница                   ║
║  • /time       - Текущее время                      ║
║  • /index      - Приветствие                        ║
║  • /api        - Страница API                       ║
║  • /api/data   - JSON API                           ║
╚══════════════════════════════════════════════════════╝
`

	fmt.Println(startMsg)

	logger.Println("Сервер запущен на порту 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
		log.Fatal("Ошибка запуска сервера:", err)
	}

}
