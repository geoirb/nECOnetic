import pandas as pd

COMPONENTS = ['CO', 'NO', 'NO2', 'PM2.5', 'PM10']

DATETIME_COL = 'Дата и время'

DEFAULT_COMPONENT = 'CO'
DEFAULT_STATION = 'Останкино'
DISABLED_COMPONENTS = []

PERIODS = ['year', 'month', 'week', 'day']

STATION_COL = 'станция'

TIMEDELTAS = {'year': pd.Timedelta(days=365),
              'month': pd.Timedelta(days=31),
              'week': pd.Timedelta(days=7),
              'day': pd.Timedelta(hours=24)}

REQUIRED_COLUMNS = COMPONENTS + [DATETIME_COL, 'станция', 'lat', 'lon']

