# Обработчик транзакций платежной системы

REST API сервис для обработки транзакций между кошельками

## 📌 Основные функции

- Перевод средств между кошельками: POST /api/send
- Просмотр истории транзакций: GET /api/transactions?count=N
- Проверка баланса кошелька:  GET /api/wallet/{address}/balance
- Автоматическое создание 10 тестовых кошельков при первом запуске

## 🚀 Быстрый старт

### Требования
- Go 1.23+
- PostgreSQL 15+
- Docker (опционально)

### Установка
```bash
git clone https://github.com/Stacypapo/golang-testtask
cd golang-testtask
go mod download
go mod tidy
```

### Конфигурация
Создайте .env файл в корневой директории и добавьте следующие значения:
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=<your db password>
DB_NAME=<name of db>
DB_SSLMODE=disable
```

### Запуск
```bash
go run cmd/main.go
```

### 🐳 Запуск в Docker
```bash
docker-compose up --build
```

## 🧪 Тестирование
В корневой директории выполните:
```bash
go test -v ./internal/handler ./internal/service ./internal/repository    
```

## 🔒 Безопасность
- Валидация всех входящих параметров
- Защита от SQL-инъекций
- Проверка достаточности баланса перед переводом

## 📚 Документация
Документация по API представлена в Swagger: http://localhost:8080/swagger/index.html 