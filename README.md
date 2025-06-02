# Final Project

This is a final project for demonstrating Go web development skills.

## Features

*   Serves a homepage with a welcome message.

*   Provides the current server time.

*   Greets users by name from a URL parameter.

*   Returns JSON data from an API endpoint.

*   Logs requests using goroutines.

## Instructions

1.  Clone the repository.

2.  Run `go mod init final_project` to initialize the Go module.

3.  Run `go run ./cmd/main.go` to start the server.

4.  Access the application in your browser:

    *   Homepage: `http://localhost:8080/`

    *   Time: `http://localhost:8080/time`

    *   Index: `http://localhost:8080/index?name=YourName`

    *   API: `http://localhost:8080/api/data`

## Testing

To run the tests, use the following command:

bash

go test ./…

## Logging

Request logs are written to the `request.log` file.

Дополнительные замечания и улучшения:

Обработка ошибок:  В коде добавлена более детальная обработка ошибок и логирование.

Конфигурация:  Использование переменных окружения для порта, что делает приложение более гибким.  Можно добавить и другие параметры конфигурации через переменные окружения.

Логирование: Добавлен mutex для защиты от race conditions при параллельной записи в файл логов.

Обработка статики:  Добавлен обработчик для статических файлов, таких как CSS и JS.

.env file: Добавлена поддержка .env file, для загрузки переменных окружения.

GetProjectDir: Добавлена функция для получения абсолютного пути до корневой директории проекта.

Этот код предоставляет базовую структуру и функциональность для вашего проекта.  Вы можете расширить его, добавляя новые функции, улучшая обработку ошибок и реализуя более сложные шаблоны.
