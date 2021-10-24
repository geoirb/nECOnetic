# nECOnetic

## Разворачивание системы

* Запуск системы
```
make up
```
* Миграция данных
```
make 
```

## Data-service

maintainers:

- Морочев Георгий morochev.g@gmail.coma

Функционал:
* Хранение и управление списком станций
* Обработка, управление и хранение данных со станций экомониторинга и профилемера
* Хранение и доступ к расчетных данных

Документация:
 * [REST API](data-service/docs/API.md)
 * [TODO](data-service/docs/TODO.md)


## Model-service

Функционал:
* Обработка данных мониторинга
* Возвращение предсказанных значений в формате .xlsx
* Запись предсказанных значений в БД

Документация:
 * [REST API](model-service/docs/API_model.md)

## !!!WARNING!!!

Для корректной обработки файлов с данными, необходимо чтобы были:
- правильные и единообразные названия параметров измерения концентраци
- дата и время в формате: dd/mm/yyyy hh:mm(:ss - не обязательно)


