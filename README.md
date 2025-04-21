![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=Prometheus&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

# Avito PVZ Task · Backend на Go 1.23

## Что реализовано
| Функция                                       | Статус  |
|-----------------------------------------------|:-------:|
| Регистрация/логин dummy‑токеном               |    ✔    |
| CRUD‑операции ПВЗ / приёмок / товаров         |    ✔    |
| Валидация ролей (``moderator``, ``employee``) |    ✔    |
| Фильтр/пагинация списка ПВЗ                   |    ✔    |
| Удаление товара LIFO и закрытие приёмки       |    ✔    |
| gRPC‑метод GetPVZList (порт ``3000``)         |    ✔    |
| Логирование (Zap)                             |    ✔    |
| PostgreSQL без ORM (sqlx+squirrel)            |    ✔    |
| Unit‑coverage ≥ 75 %                          | ✔(≈82%) |
| Интеграционный тест «full flow»               |    ✔    |
| Линтеры                                       |    ✔    |
| Docker / Docker‑Compose                       |    ✔    |

---

## Стек
- Язык:	Go 1.23
- HTTP‑роутер:	chi
- Логирование:	zap
- БД:	PostgreSQL 15 (``sqlx``, ``squirrel``)
- Миграции:	``psql`` + init‑container
- gRPC:	``google.golang.org/grpc``, ``protoc``
- Метрики:	Prometheus client
- Тесты/моки:	``testing``, ``testify``, ``sqlmock``, ``gomock``
- Линт:	golangci‑lint
- Контейнеры:	Docker & Docker‑Compose        
---

## Структура
```
.
├── cmd/
│   ├── service/           # HTTP + gRPC + /metrics
│   └── integration_test/  # e2e flow
│   
├── internal/
│   ├── api/               # HTTP‑хендлеры
│   ├── auth/              # Токены
│   ├── db/                # Репозиторий
│   ├── grpc/              # gRPC‑server
│   ├── logging/           # zap‑wrapper
│   ├── metrics/           # Prom‑middleware
│   └── model/             # Модельки
│   
├── migrations/            # *.sql DDL
├── pkg/proto/…            # сгенерированный gRPC
├── prometheus/            # prometheus
├── proto/                 # *.proto
├── scripts/               # init‑SQL для PG
├── docker‑compose.yml     # docker-compose
├── Dockerfile             # Dockerfile
└── swagger.yaml           # OpenAPI
```

---

## Запуск

### Docker-Compose
```bash
git clone https://github.com/51mans0n/avito-pvz-task.git
cd avito-pvz-task
docker compose up -d --build
# REST      → http://localhost:8080
# gRPC      → localhost:3000
# Prometheus→ http://localhost:9000   (UI)
```

### Локально
- установить PostgreSQL, Go ≥ 1.23, protoc ≥ 25
- создать БД master:master@localhost:5432/master
- выполнить миграции psql -f migrations/*.sql
```bash
go run ./cmd/service        # HTTP + gRPC + metrics
```

---

## REST эндпоинты
| Метод |                                                  URL                                                  |                Роль                 |    Описание     |
|-------|:-----------------------------------------------------------------------------------------------------:|:-----------------------------------:|:---------------:|
| POST  |                                      /dummyLogin ?role=moderator                                      |                  -                  |   Тест‑токен    |
| POST  |                                                 /pvz                                                  |              moderator              |   Создать ПВЗ   |
| GET   |                                    /pvz ?page=&limit=&startDate=&…                                    |         employee/moderator          |     Список      |
| POST  |                                              /receptions                                              |              employee               | Открыть приёмку |
| POST  |                                               /products                                               |              employee               | Добавить товар  |
| POST  |                                     /pvz/{id}/delete_last_product                                     |              employee               |  LIFO‑удаление  |
| POST  |                                    /pvz/{id}/close_last_reception                                     |              employee               | Закрыть приёмку |
---

## gRPC

Service ``pvz.v1.PVZService``

Method ``GetPVZList``

Порт ``3000``
```
Проверка:
grpcurl -plaintext localhost:3000 pvz.v1.PVZService/GetPVZList
```

---

## Метрики Prometheus
|             Метрика             |    Тип    |         Лейблы         |
|:-------------------------------:|:---------:|:----------------------:|
|      http_requests_total	       |  Counter  | ``method,path,status`` |
| http_request_duration_seconds	  | Histogram |    ``method,path``     |
|       pvz_created_total	        |  Counter  |           -            |
|    receptions_created_total	    |  Counter  |           -            |
|     products_created_total	     |  Counter  |           -            |

---

## Тесты и линт
```bash
# линтеры
golangci-lint run

# все тесты + покрытие
go test ./... -cover
# HTML‑отчёт
go tool cover -html=cover.out
```

---

## To Do
- Полноценная регистрация/логин с bcrypt + JWT(скрыть в config/env)
- Настроить кодогенерацию DTO endpoint'ов по openapi схеме

---

## Автор
- Максим Скороход · Алматы
- GitHub https://github.com/51mans0n 
- Telegram: Simanson