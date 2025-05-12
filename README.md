# 🧮 Распределённый Калькулятор

## Описание
Это распределённое приложение-калькулятор с поддержкой:

- ✅ Регистрации и авторизации через JWT  
- ✅ Хранения выражений в базе данных SQLite  
- ✅ Поддержки нескольких пользователей  
- ✅ Асинхронной обработки выражений  
- ✅ Персистентности (сохранение состояния после перезапуска)  
- ✅ REST API  
- ✅ Worker'ами для вычислений  

> В будущем можно добавить gRPC вместо HTTP между компонентами для повышения производительности.

---

## 📌 Функционал

| Компонент | Возможности |
|----------|-------------|
| Orchestrator | Принимает запросы от пользователя, управляет задачами |
| Worker | Вычисляет задачи асинхронно |
| База данных | Хранит данные пользователей и выражений в `SQLite` |

---

## 🚀 Как использовать

### Установка

1. Склонируй репозиторий:
```bash
git clone https://github.com/InsafMin/go-web-calculator.git
cd go-web-calculator
```

2. Убедись, что установлен Go 1.22+

3. Запусти orchestrator:
```bash
go run ./cmd/
```

4. В другом терминале запусти worker:
```bash
export ORCHESTRATOR_URL=http://localhost:8080
go run ./internal/agent/worker/
```

5. Можно запустить несколько worker'ов:
```bash
go run ./internal/agent/worker/
```

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

## 📡 API Эндпоинты

| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/api/v1/register` | Регистрация пользователя |
| POST | `/api/v1/login` | Получение JWT токена |
| POST | `/api/v1/calculate` | Отправка выражения на вычисление |
| GET  | `/api/v1/expressions` | Получить все выражения текущего пользователя |
| GET  | `/api/v1/expressions/{id}` | Получить конкретное выражение |
| GET / POST | `/internal/task` | Обмен задачами между orchestrator'ом и worker'ом |

---

## 🧪 Примеры запросов

### 1. Регистрация

```bash
curl -X POST http://localhost:8080/api/v1/register \
     -H "Content-Type: application/json" \
     -d '{"login": "user", "password": "1234"}'
```

**Ожидаемый ответ:**

```json
{"message": "OK"}
```

---

### 2. Авторизация

```bash
curl -X POST http://localhost:8080/api/v1/login \
     -H "Content-Type: application/json" \
     -d '{"login": "user", "password": "1234"}'
```

**Ответ:**

```json
{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxxxx"}
```

Сохраните токен — он нужен для последующих запросов.

---

### 3. Вычисление выражения

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxxxx" \
     -H "Content-Type: application/json" \
     -d '{"expression": "2+2*2"}'
```

**Ответ:**

```json
{"id": "1746998389082085000"}
```

---

### 4. Получить список выражений пользователя

```bash
curl -X GET http://localhost:8080/api/v1/expressions \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxxxx"
```

**Ответ:**

```json
{
  "expressions": [
    {
      "id": "1746998389082085000",
      "expression": "2+2*2",
      "status": "done",
      "result": 6
    }
  ]
}
```

---

### 5. Получить одно выражение

```bash
curl -X GET http://localhost:8080/api/v1/expressions/1746998389082085000 \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xxxxx"
```

**Ответ:**

```json
{
  "expression": {
    "id": "1746998389082085000",
    "expression": "2+2*2",
    "status": "done",
    "result": 6
  }
}
```

---

## 🛠️ Конфигурация

Можно настроить через переменные окружения:

| Переменная | По умолчанию | Назначение |
|------------|---------------|------------|
| `COMPUTING_POWER` | `4` | Число worker'ов |
| `ORCHESTRATOR_URL` | `http://localhost:8080` | URL orchestrator'а |
| `TIME_ADDITION_MS` | `100` | Время сложения (мс) |
| `TIME_SUBTRACTION_MS` | `100` | Время вычитания (мс) |
| `TIME_MULTIPLICATIONS_MS` | `200` | Время умножения (мс) |
| `TIME_DIVISIONS_MS` | `200` | Время деления (мс) |

---

## 🧾 Структура проекта

```
github.com/InsafMin/go-web-calculator/
├── cmd/
│   └── main.go ← точка входа сервиса
├── internal/
│   ├── auth/         ← JWT авторизация
│   ├── db/           ← работа с SQLite
│   ├── orchestrator/
│   │   ├── server/   ← запуск сервера
│   │   └── handlers/ ← обработчики API
│   └── agent/
│       └── worker/   ← worker для вычисления задач
├── pkg/
│   ├── calculator/   ← логика калькулятора
│   ├── errors/       ← ошибки
│   └── types/        ← общие структуры
├── migrations/
│   └── init.sql      ← миграции БД
├── go.mod
└── README.md
```

---

## 🗃️ Хранилище данных

Используется SQLite с двумя таблицами:

### Таблица `users`

| Поле | Тип | Описание |
|------|------|-----------|
| id | INTEGER | Уникальный ID |
| login | TEXT | Логин пользователя |
| password_hash | TEXT | Хэш пароля |

### Таблица `expressions`

| Поле | Тип | Описание |
|------|------|-----------|
| id | TEXT | Уникальный ID выражения |
| user_id | INTEGER | Ссылка на пользователя |
| expression | TEXT | Исходное выражение |
| status | TEXT | pending / in_progress / done |
| result | REAL | Результат вычисления (`NULL`, если ещё не готов) |

---

## 🧪 Тестирование

Проект покрыт:
- ✅ Unit-тестами (логика калькулятора, парсер RPN)
- ✅ Интеграционными тестами (API, взаимодействие с БД)

Пример запуска всех тестов:
```bash
go test ./...
```

---

## 🧩 TODO

- [ ] Перевести взаимодействие между orchestrator и worker на gRPC
- [ ] Добавить Makefile для удобного запуска и тестирования
- [ ] Добавить Docker Compose
- [ ] Покрыть тестами маршруты и middleware

---

## 🧬 Требования

- Go 1.22+
- SQLite
- curl или Postman для тестирования

---

## 📞 Автор

GitHub: [https://github.com/InsafMin/go-web-calculator](https://github.com/InsafMin/go-web-calculator)  
Email: insaf.min.in@yandex.ru

---

## ✅ Лицензия

MIT License – свободное использование и модификация разрешены.