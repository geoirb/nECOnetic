# REST API

## Формат ответа

```json
{
    "is_ok": bool,
    "payload": {},
}
```

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

payload:
```json
{
    "id":string,
    "name": string,
    "lat": float64,
    "lon": float64
}
```

пример запроса:
```bash
curl --request POST \
  --url http://127.0.0.1:8000/api/v1/data-service/station \
  --header 'Content-Type: application/json' \
  --data '	{
		"name": "Академика Анохина",
		"lat":  55.658163,
		"lon":  37.471434
	}'
```

## Получение списка станций

URI: /api/v1/data-service/station
Method: GET

payload:
```json
{
    "station": [
        {
        "id": "616feb34e3bf1d6ceefeb706",
        "name": "Академика Анохина",
        "lat": 55.658163,
        "lon": 37.471434
        },
        {
        "id": "617001935df807ab45002de9",
        "name": "Глебовская",
        "lat": 55.811801,
        "lon": 37.71249
        },
        {
        "id": "617001935df807ab45002dea",
        "name": "Коптевский",
        "lat": 55.833222,
        "lon": 37.525158
        },
        {}
    ]
}
```

пример запроса
```bash
curl --request GET \
  --url http://127.0.0.1:8000/api/v1/data-service/station
```

## Добавить измерения в систему

URI: /api/v1/data-service/station/data
Method: POST

body:
format: Multipart
|  ключ   | тип значения |                  описание                   |
|:-------:|:------------:|:-------------------------------------------:|
| station |    string    | имя станции на которой данные были получены |
|  type   |    string    |   тип данных(`eco`,`temperature`,`wind` )   |
|  data   |     file     |               файл с данными                |


Описание типов данных, которые принимает система
| тиа данных  |                       описание                       | расширение файла |
|:-----------:|:----------------------------------------------------:|:----------------:|
|     eco     |  данные экомониторинга(концентрации веществ и т.д.)  |      .xlsx       |
| temperature | данные с профилемера по тепературе на разных высотах |       .txt       |
|    wind     |    данные измерения ветра(скорость и направление)    |      .xlsx       |

Требование к файлу:
```
Для корректной обработке файлов с данными, необходимо чтобы были:
- правильные и единообразные названы параметров измерения концентраци
- дата и время в формате: dd/mm/yyyy hh:mm(:ss - не обязательно)
```

пример запроса
```bash
curl --request POST \
  --url http://127.0.01:8000/api/v1/data-service/station/data \
  --header 'Content-Type: multipart/form-data; boundary=---011000010111000001101001' \
  --form 'station=Академика Анохина' \
  --form type=eco \
  --form 'data=@/home/geoirb/project/nECOnetic/dataset/Академика Анохина 2020.xlsx'
```

## Добавить расчеты в систему

URI: /api/v1/data-service/predict
Method: POST


Тело запроса:
```json
{
  "data": [
    {
      "station_id": "616feb34e3bf1d6ceefeb706",
      "timestamp": 1609443600,
      "measurement": {
        "-T-": -1.1,
        "CO": 0.4,
        "NO": 0.001,
        "NO2": 0.0235,
        "PM2.5": 0.024,
        "_V_": 228,
        "| V |": 0.5,
        "Влажность": 91,
        "Давление": 746.6,
        "Осадки": 0
      }
    },
    {}
  ]
}
```

## Получить список измерений


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

payload:
```json
{
    "id":string,
    "name": string,
    "lat": float64,
    "lon": float64
}
```

пример запроса:
```bash
curl --request POST \
  --url http://127.0.0.1:8000/api/v1/data-service/station \
  --header 'Content-Type: application/json' \
  --data '	{
		"name": "Академика Анохина",
		"lat":  55.658163,
		"lon":  37.471434
	}'
```

## Получение списка данных экомониторинга

URI: /api/v1/data-service/station/eco-data?station={`string`}&timestamp_from={`int64`}&timestamp_till={`int64`}&measurement={measurement_1}&measurement={measurement_2}
Method: GET

Параметры запроса:
>station - станция, по которой нужно получить данные экомониторинга, если параметр не указан то фильтрация по станциям не проводится
>
>timestamp_from - начало интервала, за который нужно получить данные экомониторинга
>
>timestamp_till - конец интервала, за который нужно получить данные экомониторинга
>
>measurement -  список измерений, которые должны присутствовать в выборке

payload:
```json
{
  "data": [
    {
      "station_id": "616feb34e3bf1d6ceefeb706",
      "timestamp": 1609443600,
      "measurement": {
        "-T-": -1.1,
        "CO": 0.4,
        "NO": 0.001,
        "NO2": 0.0235,
        "PM2.5": 0.024,
        "_V_": 228,
        "| V |": 0.5,
        "Влажность": 91,
        "Давление": 746.6,
        "Осадки": 0
      }
    },
    {
      "station_id": "616feb34e3bf1d6ceefeb706",
      "timestamp": 1609444800,
      "predicted_measurement": {
        "-T-": -1.3,
        "CO": 0.4,
        "NO": 0.001,
        "NO2": 0.0207,
        "PM2.5": 0.023,
        "_V_": 0,
        "| V |": 0.9,
        "Влажность": 91,
        "Давление": 746.6,
        "Осадки": 0
      }
    }
  ]
}
```
> measurement - измеренные данные экомониторинга, набор полей меняется в зависимости от того какие данные пришли в файле
>
> predicted_measurement - расчетные данные экомониторинга, выдаются только в том случа, если нет данных со станций

пример запроса
```bash
curl --request GET \
  --url 'http://127.0.01:8000/api/v1/data-service/station/eco-data?station="Академика Анохина"&timestamp_from=1609443600&timestamp_till=1609443600'
```
<!-- 
## Запустить расчет параметров


URI: /api/v1/data-service/predict?station={`string`}&timestamp_from={`int64`}&timestamp_till={`int64`}
Method: GET

Параметры запроса:
>station - станция, для которой необходимо произвести расчеты, если данный параметр не указан, то расчет идет по всем стациям
>
>timestamp_from - начало интервала данных, по которым будет произведен расчет
>
>timestamp_till - конец интервала данных, по которым будет произведен расчет -->
