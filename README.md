# nECOnetic

## Разворачивание системы

<!-- TODO: начало -->

* Настройка базы данных
```
docker exec -it deployment_mongo_1 bash

mongo

rs.initiate(
{
    _id : 'neconetic',
    members: [
    { _id : 0, host : "0.0.0.0:27017" },
    ]
}
)

db.isMaster()
```

## Data-service

maintainers:

- Морочев Георгий morochev.g@gmail.coma

description:
 
 Данный сервис занимается хранинением измерений

