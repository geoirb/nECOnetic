import datetime
import os

import numpy as np
import pandas as pd

import pymongo

from typing import Any, Dict, List, Optional

from pymongo import MongoClient

from constants import COMPONENTS, DATETIME_COL, REQUIRED_COLUMNS, STATION_COL

DB_HOST = os.getenv("DB_HOST", "localhost")
DB_PORT = os.getenv("DB_PORT", "27017")

DB_NAME = 'neconetic'
DB_USER = 'neconetic'
DB_PASSWORD = 'neconetic_secrete'

client = MongoClient(f"mongodb://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/", serverSelectionTimeoutMS=30000)

db = client.get_database(DB_NAME)


def get_measurements_df(measurements_list: List[Dict[str, Any]], id_column: str) -> pd.DataFrame:
    """
    Возвращает датафрейм с данными измерений

    Вспомогательная функция для формирования датафрейма из данных,
    представленных в формате, в котором они поступают из БД. Достает нужные
    поля, создает из них датафрейм и, если какого-то из нужных полей не
    оказалось в базе, добавляет пустую колонку с названием этого поля.
    :param measurements_list: список записей из БД
    :param id_column: название колонки с индексом
    :return: датафрейм с данными измерений
    """
    data = []

    for item in measurements_list:
        m = item['measurement']
        m['timestamp'] = item['timestamp']
        m[id_column] = item[id_column]
        data.append(m)

    measurements_df = pd.DataFrame(data)

    for c in COMPONENTS:
        if c not in measurements_df:
            measurements_df[c] = np.nan

    return measurements_df


def get_stations(stations: Optional[List[str]] = None) -> pd.DataFrame:
    """
    Достает из БД информацию о станциях и кладет ее в датафрейм

    :param stations: список названий станций, информацию о которых достаем
    из БД
    :return: датафрейм с информацией о станциях
    """
    stations_data = db.get_collection('station')

    if stations:
        query = {
            "name": {
                "$in": stations,
            }
        }
    else:
        query = {}

    stations_list = list(stations_data.find(query))
    stations_df = pd.DataFrame(stations_list)

    return stations_df


def get_measurements_data(stations=None, from_date=None, to_date=None) -> pd.DataFrame:
    """
    Достает данные об измерениях компонентов из БД и кладет их в датафрейм

    :param stations: список станций, для которых запрашиваются данные
    :param from_date: дата, с которой запроашиваются данные
    :param to_date: дата, до которой запрашиваются данные
    :return: датафрейм с информацией о компонентах
    """
    from_date = from_date or 1
    to_date = to_date or datetime.datetime.now().timestamp()

    stations_df = get_stations(stations)

    if not stations_df.empty:
        stations_ids = list(stations_df['_id'].unique())
    else:
        stations_ids = []

    measurements = db.get_collection('eco-data')

    query = {
        "station_id": {
            "$in": stations_ids
        },
        "timestamp": {
            "$gt": from_date,
            "$lt": to_date
        }
    }

    sorting = [("station_id", pymongo.ASCENDING), ("timestamp", pymongo.DESCENDING)]
    measurements_list = list(measurements.find(query).sort(sorting))
    measurements_df = get_measurements_df(measurements_list, id_column="station_id")

    if any((stations_df.empty, measurements_df.empty)):
        return pd.DataFrame()

    df = pd.merge(measurements_df, stations_df, how='left', left_on="station_id", right_on="_id")
    df = df.rename(columns={
        "name": STATION_COL,
        "timestamp": DATETIME_COL,
    })
    df[DATETIME_COL] = pd.to_datetime(df[DATETIME_COL], unit='s')

    return df[REQUIRED_COLUMNS]


def get_newest_measurements() -> pd.DataFrame:
    """
    Достает данные о последних измерениях компонентов из БД

    :return: датафрейм с информацией о последних измерениях компонентов
    """
    stations_df = get_stations()

    measurements = db.get_collection('eco-data')

    sort_query = {
        "$sort": {"timestamp": -1}
    }

    group_query = {
        "$group": {
            "_id": "$station_id",
            "measurement": {"$first": "$measurement"},
            "timestamp": {"$first": "$timestamp"}
        }
    }

    measurements_list = list(measurements.aggregate([
        sort_query,
        group_query
    ]))
    measurements_df = get_measurements_df(measurements_list, id_column="_id")

    if any((stations_df.empty, measurements_df.empty)):
        return pd.DataFrame()

    df = pd.merge(measurements_df, stations_df, how='left', on='_id')
    df = df.rename(columns={
        "name": STATION_COL,
        "timestamp": DATETIME_COL,
    })
    df[DATETIME_COL] = pd.to_datetime(df[DATETIME_COL], unit='s')

    return df[REQUIRED_COLUMNS]
