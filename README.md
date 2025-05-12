# 🧮 Распределённый Калькулятор

> Проект: веб-калькулятор с JWT-авторизацией, SQLite-персистентностью и gRPC-взаимодействием между оркестратором и вычислителем (worker).

## 🎯 Основные функции

1. ✅ Регистрация и авторизация пользователя через JWT
2. ✅ Вычисление математических выражений
3. ✅ Хранение выражений в **SQLite**
4. ✅ Поддержка нескольких пользователей
5. ✅ Восстановление состояния после перезапуска
6. ✅ gRPC вместо HTTP между компонентами (до этого использовался HTTP)
7. ✅ Персистентность:
   - Все данные сохраняются между перезапусками
   - Worker'ы получают задачи из оркестратора
8. ✅ Интеграционное взаимодействие:
   - Orchestrator запускается первым
   - Worker'ы запускаются после и подключаются к оркестратору

---

## 📌 Структура проекта

```
go-web-calculator/
├── go.mod
├── README.md
├── .dockerignore
├── Dockerfile
├── docker-compose.yml
├── cmd/
│   ├── calc_service/
│   │   └── main.go       ← точка входа оркестратора
│   └── worker/
│       └── main.go       ← точка входа worker'а
├── internal/
│   ├── orchestrator/
│   │   ├── server/
│   │   └── auth/
│   │       ├── jwt.go           ← jwt токены
│   │       ├── middleware.go    ← мидлвейр авторизации
│   │   └── handlers/
│   │       ├── user.go          ← регистрация/логин
│   │       ├── calculate.go     ← обработка выражения
│   │       ├── expressions.go   ← получение выражений
│   │       └── task.go          ← работа с задачами
│   └── agent/
│       └── worker/
│           └── worker.go         ← логика worker'а
├── pkg/
│   ├── calculator/               ← основная логика вычисления
│   ├── types/                    ← типы данных
│   └── errors/errors.go          ← ошибки
├── proto/
│   └── task.proto                ← gRPC сервис
└── migrations/
    └── init.sql                  ← миграции БД
```

---

## 🚀 Как запустить проект

### 1. Клонируй репозиторий:

```bash
git clone https://github.com/InsafMin/go-web-calculator.git
cd go-web-calculator
```

### 2. Убедись, что установлены:

- Go 1.23+
- Docker & Docker Compose
- curl или Postman для тестирования

### 3.1 Запусти проект через консоль/терминал:

- Первое окно
```bash
go run ./cmd/calc_service/
```
- Второе окно
```bash
go run ./cmd/worker/
```

### 3.2 Или запусти проект через Docker Compose:

```bash
docker-compose up --build
```

→ Это соберёт и запустит:
- Оркестратор на порту `8080`
- Worker'ов
- Базу данных SQLite (`calc.db`)

---

## 🔐 Регистрация и Авторизация

### Регистрация нового пользователя

```bash
curl -X POST http://localhost:8080/api/v1/register \
     -H "Content-Type: application/json" \
     -d '{"login": "user1", "password": "pass1"}'
```

**Ответ:**
```json
{"message": "OK"}
```

---

### Авторизация (получение JWT токена)

```bash
curl -X POST http://localhost:8080/api/v1/login \
     -H "Content-Type: application/json" \
     -d '{"login": "user1", "password": "pass1"}'
```

**Ответ:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxxxx"
}
```

---

## 🧪 Примеры использования API

### Отправить выражение на вычисление

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
     -H "Authorization: Bearer <your_token>" \
     -H "Content-Type: application/json" \
     -d '{"expression": "2+2*2"}'
```

**Ответ:**
```json
{
  "id": "1234567890"
}
```

---

### Получить список всех выражений текущего пользователя

```bash
curl -X GET http://localhost:8080/api/v1/expressions \
     -H "Authorization: Bearer <your_token>"
```

**Ответ:**
```json
{
  "expressions": [
    {
      "id": "1234567890",
      "expression": "2+2*2",
      "status": "done",
      "result": 6
    }
  ]
}
```

---

### Получить конкретное выражение по ID

```bash
curl -X GET http://localhost:8080/api/v1/expressions/1234567890 \
     -H "Authorization: Bearer <your_token>"
```

**Ответ:**
```json
{
  "expression": {
    "id": "1234567890",
    "expression": "2+2*2",
    "status": "done",
    "result": 6
  }
}
```

---

## 🧱 Архитектура

Проект состоит из следующих частей:

| Компонент | Описание |
|----------|----------|
| **Orchestrator** | Основной сервер, принимает запросы от пользователя, хранит задачи и управляет worker'ами |
| **Worker** | Вычисляет выражения и отправляет результат обратно |
| **gRPC** | Общение между оркестратором и worker'ами (заменяет HTTP) |
| **SQLite** | Хранение выражений и пользователей |

---

## 🛠️ Технологии

- ✅ Golang
- ✅ REST API
- ✅ JWT для аутентификации
- ✅ SQLite для хранения данных
- ✅ gRPC для коммуникации между компонентами
- ✅ Docker + Docker Compose для развёртки

---

## 🧪 Проверка работы

После запуска:

```bash
sqlite3 calc.db
SELECT * FROM expressions;
```

→ Должны быть выражения, созданные разными пользователями.

---

## 📦 Docker Compose

Содержит:

- `orchestrator` → `cmd/calc_service/main.go`
- `worker1`, `worker2` → `cmd/worker/main.go`

База данных:
- Сохраняется в `./calc.db`
- При перезапуске все данные остаются

---

## 🧪 Возможные ошибки при работе

| Ошибка | Что делать |
|--------|-------------|
| ❌ `no tasks available` | Проверь, правильно ли выражение разбито на задачи |
| ❌ `unauthorized` | Убедись, что токен корректен и передан в заголовке |
| ❌ `division by zero` | Отправь выражение без деления на ноль 😄 |
| ❌ `expression not found` | Проверь, существует ли выражение в базе |

---

## 🧪 Логи контейнеров

```bash
docker logs go-web-calculator-orchestrator-1
docker logs go-web-calculator-worker1-1
```

---

## 🧪 Тестирование

### Unit-тесты

```bash
go test ./pkg/calculator/
```

### Интеграционные тесты

```bash
go test ./internal/orchestrator/handlers/
```

---

## 🧪 Пример успешного вывода

```
Starting orchestrator on :8080
Worker started
Got task: ID=1747143071, Expr="2+2*2"
Task finished successfully: ...
Sent result for 1747143071: 6.0
```

---

## 🧪 Формат JSON ответа

| Поле | Тип | Описание |
|------|-----|-----------|
| `id` | string | Уникальный ID выражения |
| `expression` | string | Исходное выражение |
| `status` | string | `pending`, `in_progress`, `done`, `error` |
| `result` | float64 или null | Результат вычисления |
| `error` | string или null | Сообщение об ошибке (если есть) |

---

## ✅ Что проверяющий должен увидеть:

| Пункт | Что должно быть реализовано |
|------|---------------------------|
| 1 | ✅ Работа в контексте пользователя |
| 2 | ✅ Хранение в SQLite, восстановление после перезагрузки |
| 3 | ✅ gRPC вместо HTTP между компонентами |
| 4 | ✅ Unit-тесты |
| 5 | ✅ Интеграционные тесты |
| 6 | ✅ Документация Readme и примеры запросов |
| 7 | ✅ Возможность запустить проект одной командой `docker-compose up --build` |

---

## 🧬 Автор

GitHub: [https://github.com/InsafMin/go-web-calculator](https://github.com/InsafMin/go-web-calculator)  
Email: insaf.min.in@yandex.ru