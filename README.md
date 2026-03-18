# TestForOzon

gRPC-сервис на Go для получения данных по валютам на указанную дату (статические данные + динамика курса) с хранением в PostgreSQL.

## Микросервис отвечающий за http, который связывается с помощью gRPC 

https://github.com/L1irik259/TestForOzonHTTPService


## Пример работы

![Demo](https://github.com/L1irik259/TestForOzon/blob/main/pkg/assets/gifForOzonTest.gif)

## Структура проекта

```
TestForOzon/
├── cmd/app/               # Точка входа приложения
│   ├── main.go           # Инициализация сервиса
│   ├── main_test.go      # E2E тесты
│   └── .env              # Переменные окружения
├── internal/
│   ├── domain/           # Бизнес-логика (Item, ItemStaticData, ItemDynamicData)
│   ├── adapter/          # Слой доступа к БД
│   ├── service/          # Бизнес-сервисы (ItemService)
│   ├── mapper/           # Мапперы (Proto ↔ Domain)
│   └── transport/
│       ├── service/      # gRPC сервис (OzonService)
│       └── proto/        # Proto файлы и сгенерированный код
├── pkg/transport/        # Вспомогательные компоненты транспорта
├── go.mod                # Зависимости проекта
├── go.sum                # Хешь зависимостей
├── Taskfile.yml          # Задачи для запуска
├── application.yml       # Docker Compose конфиг
└── Readme.md             # Этот файл
```

## Что делает сервис

- Реализует gRPC метод `GetItem`.
- Принимает дату и возвращает список валют на эту дату.
- Конвертирует domain-модели в protobuf через mapper.
- Работает с PostgreSQL через adapter-слой.

## Важный момент по формату даты

В `GetItem` используется парсинг:

- `time.Parse("02/01/2006", req.Date)`

То есть ожидаемый формат: **`DD/MM/YYYY`** (например, `18/03/2026`).

> В `ozon.proto` комментарий сейчас указан как `DD-MM-YYYY`, это не совпадает с текущей реализацией.

## Быстрый запуск

### 1) Поднять PostgreSQL

```powershell
docker compose -f .\application.yml up -d
```

### 2) Проверить env

Файл: `cmd/app/.env`

```env
DATABASE_URL=postgres://postgres:test123@localhost:5432/postgres?sslmode=disable
```

### 3) Запустить сервис

```powershell
go run .\cmd\app
```

## Запуск через Taskfile / Makefile

Если в проекте есть Taskfile:

```powershell
task run
```

Если в проекте есть Makefile:

```powershell
make run
```
