# gushort

Учебный сервис, реализующий функционал для сокращения ссылок и последующего редиректа на них.

Пользователь передаёт длинный URL — сервис возвращает короткую ссылку, по которой происходит редирект на оригинальный адрес.

## Стек

- Go
- SQLite
- HTTP (REST)
- Goose (миграции)

## Конфигурация

Сервис настраивается через YAML-файл конфигурации.

### Пример `config.yaml`

```yaml
env: "prod" # local | dev | prod

storage_path: "./storage.sqlite"

http_server:
  address: "0.0.0.0:8080"
  timeout: 4s
  idle_timeout: 30s
```

## Миграции

```bash
goose -dir ./migrations sqlite3 storage.sqlite up
```
