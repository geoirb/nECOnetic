# REST API

## Добавление станций мониторинга в систему

URI: /api/v1/data-service/station
Method: POST

body:
```json
{
    "name": string,
    "lat": float64,
    "lon": float64
}
```

answer: