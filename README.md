# HezzlTest - Контроль товаров в рекламных кампаниях.

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/doc/)
[![Postgres](https://img.shields.io/badge/Postgres-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/docs/)
[![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/docs)
[![Gin](https://img.shields.io/badge/Gin-Golang-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://gin-gonic.com/docs/)
[![ClickHouse](https://img.shields.io/badge/ClickHouse-FFCC01?style=for-the-badge&logo=clickhouse&logoColor=black)](https://clickhouse.com/docs/en/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://docs.docker.com/)
[![NATS](https://img.shields.io/badge/NATS-199bfc?style=for-the-badge&logo=nats&logoColor=white)](https://docs.nats.io/)

Сервис для управления торговыми рекламными кампаниями.

## Особенности
- Worker для отправки событий в Clickhouse через Nats jetStream
- Интеграционные тесты для postgres и redis
- Redis для GET запросов
- Конфигурация приложения через viper (возможность легко поменять кфг под прод)
- CRUD для всех сущностей PostgreSQL
- Чистая архитектура с разделением слоёв
- Максимальная степень изоляции для всех транзакций PostgreSQL
- RESTful API
- Доп. оптимизация по btree индексу для ускорения расчёта priority при добавлении новых goods


## Структура проекта (Clean Architecture)
├── cmd # Точка входа  
├── config # Конфигурация приложения  
├── integration-test # Интеграционные тесты  
├── internal  
│ ├── app # Инициализация и закрытие приложения  
│ ├── controller # Логика обработчиков   
│ ├── entity # Бизнес-сущности  
│ ├── infrastructure # слой инфры (клиент nats)  
│ ├── logger # настройка логгера  
│ ├── repository # Интерфейсы хранилища  
│ │ ├── clickhouse # Методы для работы с Сlickhouse  
│ │ ├── events	# Методы для отправки эвентов в nats  
│ │ ├── postgres	# Методы для работы с postgreSQL  
│ │ ├── redis	# Методы для работы с redis  
│ ├── usecase  # Интерфейсы и реализация бизнес-логики  
│ ├── worker # Логика отправки событий в Сlickhouse  
├── migrations # Миграции Postgres  
| ├── clickhouse # Миграции Сlickhouse  
├── scripts # скрипты для запуска (entrypoint для Dockerfile)  
└── .env # Переменные окружения для docker-compose  
└── docker-compose.yml # Настройка контейнерной среды  
└── Dockerfile # Файл для создания образа нашего сервиса  
└── go.mod  # Модули и зависимости  
└── go.sum  # Модули и зависимости  
└── README.md # Описание проекта

## Требования
![Docker](https://img.shields.io/badge/Docker-Required-blue?logo=docker&style=flat)

## Установка и запуск
```bash
# Клонировать репозиторий
git clone https://github.com/paxaf/HezzlTest.git
cd hezzlTest

# Запустить все сервисы
docker-compose up -d --build #old
docker compose up -d --build #new
```
## API Endpoints

| Метод   | Путь           | Описание                 |
|---------|----------------|--------------------------|
| POST    | `/projects`      | Создать кампанию         |
| GET     | `/projets`      | Получить все кампании      |
| GET     | `/project/:id`      | Получить кампанию по id  |
| DELETE  | `/project/:id`  | Удалить кампанию           |
| POST  | `/goods`  | Создать товар           |
| GET  | `/goods`  | Получить все товары           |
| GET  | ``/:project_id/goods``  | Получить все товары по id кампании          |
| GET  | `/goods/search/:any`  | Поиск товаров с соответсвием по имени           |
| GET  | `/goods/:id`  | Получить товар по id           |
| PATCH  | `/goods`  | Обновить информацию о товаре           |
| DELETE  | `/goods/:id`  | Удалить товар           |

## Тело запросов

1. `POST /goods` - Создать товар
```JSON
{
  "name": "Товар 1", // обязательное поле
  "description": "Описание товара",
  "project_id": 123 // обязательное и должно ссылаться на существующий projects(id)
}
```

2. `PATCH /goods` - Обновить товар
```JSON
{
	// все поля обязательные
  "id": 42, 
  "name": "Обновленный товар", // не пустое
  "description": "Новое описание товара", // можно пустым
  "priority": 10 // gt=0
}
```
3. `POST /projects` - Создать проект
```JSON
{
  "name": "Новая кампания", // не пустое
}
```
4. `PATCH /projects` - Создать проект
```JSON
{
	"id"
  "name": "Новая кампания", // не пустое
}
```
## Запуск тестов
Перед запуском интеграционных тестов убедитесь что у вас запущен Docker на машине.
```bash
go test ./... -race -v
```
Если по каким то причинам `gcc` у вас в нет или флаг `-race` не работает корректно на вашей локальной машине. Тогда для теста с флагом `-race` можно открыть `Dockerfile` изменить 
```Dockerfile 
CGO_ENABLED=1 # Устанавливаем значение на 1
...
RUN go test ./... -race -cover -v -short # убираем `#` перед этой строчкой
```
Запускаем команду 
```bash
docker-compose build | tee build.log
```

И смотрим как проходят тесты с флагом `-race`.
Интеграционные тесты так не запустятся, только если нактатить в билдер docker и поднять. Поэтому там флаг -short.
Перед билдом бинарника для использования в контейнере ``обязательно`` возвращаем всё в исходное состояние, иначе проект не запустится.

