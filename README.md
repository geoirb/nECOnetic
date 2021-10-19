# nECOnetic

## Разворачивание системы

* Запуск
```
make up
```

* Настройка базы данных
```
docker exec -it deployment_mongo_1 bash

mongo

rs.initiate(
{
    _id : 'neconetic',
    members: [
    { _id : 0, host : "127.0.0.1:27017" },
    ]
}
)

db.isMaster()
```

## Data-service

maintainers:

- Морочев Георгий morochev.g@gmail.coma

description:
 
Сервис для храниения данных со станций мониторинга, профилемера и расчетных данных
Документация:
 * [REST API](data-service/docs/API.md)
 * [TODO](data-service/docs/TODO.md)

## !!!WARNING!!!

Формат времени в файлах должен быть `dd/mm/yyyy hh:mm` или `dd/mm/yyyy hh:mm:ss`


