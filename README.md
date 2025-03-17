# Golang Archive API

## Описание
Этот проект — REST API на Golang для загрузки, сжатия и архивации фото и видео.  
Использует **FFmpeg** для сжатия видео и **Go image processing** для обработки изображений.

## Возможности
- Загрузка файлов (видео и фото)
- Сжатие изображений (JPEG, PNG)
- Сжатие видео с помощью FFmpeg (MP4, AVI, MKV и др.)
- Архивация сжатых файлов в ZIP и отправка пользователю

## Технологии
- **Golang** — основной язык
- **FFmpeg** — обработка видео
- **Slog** — логирование
- **Fiber/Chi** — маршрутизация (в зависимости от использования)
- **Multipart Form** — обработка загруженных файлов

## Структура проекта
```plaintext
GOLANGARCHIVE/
│── .vscode/                 # Конфигурация для VSCode
│── cmd/                     # Точка входа в приложение
│   ├── main.go
│── config/                  # Конфигурация проекта
│   ├── config.go
│   ├── config.yml
│── internal/                # Основная бизнес-логика
│   ├── handlers/            # Обработчики HTTP-запросов
│   │   ├── middleware/      # Middleware
│   │   │   ├── logger.go
│   │   ├── upload/
│   │   │   ├── upload.go
│   ├── routes/              # Маршруты API
│   │   ├── routes.go
│   ├── service/             # Основная бизнес-логика
│   │   ├── archive/         # Работа с архивами
│   │   │   ├── archive.go
│   │   ├── comp/            # Сжатие фото/видео
│   │   │   ├── comp.go
│   │   ├── upload/          # Логика загрузки
│   │   │   ├── upload.go
│   ├── structures/          # Определение структур данных
│   │   ├── server/          # Настройки сервера
│   │   │   ├── server.go
│── pkg/                     # Вспомогательные библиотеки
│   ├── api/               
│   │   ├── response.go
│   ├── generatename/        # Генерация уникальных имен
│   │   ├── generatename.go
│   ├── logger/              # Логирование
│   │   ├── handlers/
│   │   │   ├── slogpretty.go
│   │   ├── sl.go
│── uploads/                 
│── .gitignore
│── go.mod
│── go.sum

```
## Установка и запуск
**1** Установить FFmpeg (если его нет)
**https://ffmpeg.org/download.html**
**2** Установить зависимости
go mod tidy
**3** Проверить конфиг и указать в переменных среды CONFIG_PATH
config.yaml
**4** Запустить сервер
go run cmd/main.go
## Использование API
POST /upload
Content-Type: multipart/form-data
Body: file=[Ваш файл]
