# Slowly

Сервис, эмулирующий запуск долгоживущего процесса с возможностью управления 
таймаутом запроса.

## Использование

Сборка:

```
go build -o ./target/slowly ./cmd/slowly/main.go
```

Запуск:

```
Usage of ./slowly:
  -addr string
        server address (default ":8080")
  -timeout duration
        server request timeout (default 5s)
```

Примеры запросов можно посмотреть в [requests.http](requests.http).

## Тестирование

Запуск юнит-тестов:

```
go test -tags=unit ./...
```

Запуск интеграционных-тестов:

```
go test -tags=integration ./...
```

При запуске интеграционных-тестов можно задать адрес тестового сервиса с помощью 
переменной окружения `SLOWLY_TEST_ADDR` (по умолчанию `localhost:9090`).