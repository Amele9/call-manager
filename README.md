# call-manager

### Описание

Сервис для обработки заявок. Может сохранять и удалять заявки. Может обновлять статусы заявок.

### Как запустить

Сначала необходимо запустить Docker клиент. После этого нужно создать config/configuration.yml и настроить параметры
приложения (можно скопировать из примера). А теперь запускаем наши сервисы: `docker compose up`.

#### Параметры

```yaml
port: 6689
connectionString: "postgres://username:password@postgres:5432/calls?sslmode=disable"
```

### API

#### Добавить заявку

POST http://localhost:6689/calls/

Тело запроса:

```json
{
  "client_name": "client_name1",
  "phone_number": "+77001234567",
  "description": "description1",
  "status": "open"
}
```

#### Получить все заявки

GET http://localhost:6689/calls/

#### Получить заявку по ID

GET http://localhost:6689/calls/1

#### Обновить статус заявки

PATCH http://localhost:6689/calls/1/status

#### Удалить заявку

DELETE http://localhost:6689/calls/1
