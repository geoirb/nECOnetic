import os

import dash
import dash_bootstrap_components as dbc
import pandas as pd

import plotly.express as px
import plotly.graph_objects as go

import app_layouts
import data_loader

from typing import List, Union

from constants import COMPONENTS, DATETIME_COL, STATION_COL, TIMEDELTAS
from utils import get_empty_plot, get_error_plot

HOST = os.getenv("DASHBOARD_SERVICE_HOST", "localhost")
PORT = os.getenv("DASHBOARD_SERVICE_PORT", 8050)

app = dash.Dash(__name__, external_stylesheets=[dbc.themes.BOOTSTRAP],
                meta_tags=[{"name": "viewport",
                            "content": "width=device-width, initial-scale=1.0"}])

NOW = pd.to_datetime(1609457999, unit='s')  # timestamp последнего измерения, в проде заменить везде на datetime.now()
app.layout = app_layouts.main_layout


@app.callback(
    dash.dependencies.Output("map", "figure"),
    [dash.dependencies.Input("component", "value")])
def update_map(component: str) -> go.Figure:
    """
    Отрисовка карты

    :param component: название выбранного компонента
    :return: карта с маркерами, цвет которых зависит от концентрации выбранного
    компонента в данной точке
    """
    df = data_loader.get_newest_measurements()

    if df.empty:
        return get_error_plot()

    stations = df[STATION_COL].unique()

    # для каждой станции отображается последнее измеренное значение компонента
    df = df.drop_duplicates(subset=["lon", "lat"], keep="last")

    # предельные значения концентрации выбранного компонента для отображения на цветовой шкале
    cmax = df[component].max()
    cmin = df[component].min()
    df = df.fillna('-')

    # подготовка данных для отображения на карте
    data = [
        {
            "customdata": df.loc[df[STATION_COL] == station, COMPONENTS],
            "hovertemplate": f"<extra>{station}</extra> " +
                             "<br> ".join([f"{c} = " + "%{" + f"customdata[{i}]" + "}" for i, c in enumerate(COMPONENTS)]),
            "name": station,
            "showlegend": True,
            "type": "scattermapbox",
            "lat": df.loc[df[STATION_COL] == station, "lat"],
            "lon": df.loc[df[STATION_COL] == station, "lon"],
            "marker": go.scattermapbox.Marker(
                size=25,
                opacity=0.9,
                sizemode="area",
                sizemin=4,
                color=df.loc[df[STATION_COL] == station, component] if
                df.loc[df[STATION_COL] == station, component].values[0] != '-' else 'rgba(218, 215, 218, 1)',
                coloraxis="coloraxis",
            )
        } for station in stations]

    # параметры карты
    layout = go.Layout(
        hovermode="closest",
        hoverdistance=100,
        margin=dict(l=40, r=1, t=20, b=1),
        coloraxis=dict(colorscale="blugrn",
                       cmax=cmax,
                       cmin=cmin,
                       colorbar={"title": component,
                                 "titleside": "top",
                                 "thickness": 12,
                                 "ticksuffix": "",
                                 "x": 1.03,
                                 "tickformat": ".1e"}
                       ),
        mapbox_style="carto-positron",
        mapbox_center={"lat": 55.75, "lon": 37.65},
        mapbox_zoom=9.4,
        height=865,
        legend=dict(
            itemsizing="constant",
            orientation="h",
            traceorder="normal",
            font=dict(
                family="sans-serif",
                size=12,
                color="grey"
            ),
        )
    )

    fig = go.Figure(data=data, layout=layout)

    return fig


def create_time_series(dff: pd.DataFrame, component: str) -> go.Figure:
    """
    Отрисовка распределения концентрации выбранного компонента по времени

    :param dff: датафрейм с данными о концентрации компонентов
    :param component: название компонента, для которого строим распределение
    :return: график с распределением концентрации компонента
    """
    fig = px.scatter(dff, x=DATETIME_COL, y=component, color=STATION_COL, color_discrete_sequence=px.colors.qualitative.Dark2)

    fig.update_traces(mode="lines")

    fig.update_xaxes(showgrid=False)
    fig.update_yaxes(type="linear", showline=True, linewidth=2, linecolor="rgba(216, 219, 226, 1)",
                     gridcolor="rgba(216, 219, 226, 1)")

    fig.add_annotation(x=0, y=0.85, xanchor="left", yanchor="bottom",
                       xref="paper", yref="paper", showarrow=False, align="left",
                       text="")

    fig.update_layout(height=300, margin={"l": 20, "b": 30, "r": 5, "t": 25}, paper_bgcolor="rgba(0,0,0,0)",
                      plot_bgcolor="rgba(0, 0, 0, 0)")

    return fig


@app.callback(
    dash.dependencies.Output("components-by-period", "figure"),
    [dash.dependencies.Input("stations", "value"),
     dash.dependencies.Input("component", "value"),
     dash.dependencies.Input("period", "value")])
def update_linear_past_plot(stations: Union[List[str], str], component: str, period: str):
    """
    Вызов отрисовки распределения концентрации выбранного компонента по времени

    :param stations: станции, для которых отображаем распределение
    :param component: название выбранного компонента
    :param period: период, для которого строим распределение
    :return: график с распределением концентрации компонента
    """
    stations = stations if isinstance(stations, list) else [stations]
    period_delta = TIMEDELTAS.get(period)

    t = int((NOW - period_delta).timestamp())

    df = data_loader.get_measurements_data(stations, t, NOW.timestamp())

    if df.empty:
        return get_error_plot()
    if not stations:
        return get_empty_plot()

    df_by_stations = df[df[STATION_COL].isin(stations)]
    df_by_stations_and_period = df_by_stations[df_by_stations[DATETIME_COL] >= NOW - period_delta]

    return create_time_series(df_by_stations_and_period, component)


@app.callback(
    dash.dependencies.Output("components-by-future-period", "figure"),
    [dash.dependencies.Input("stations", "value"),
     dash.dependencies.Input("component", "value")])
def update_linear_future_plot(stations, component):
    """
    Вызов отрисовки предсказанного распределения концентрации выбранного компонента по времени

    :param stations: станции, для которых отображаем распределение
    :param component: название выбранного компонента
    :return: график с предсказанным распределением концентрации компонента
    """
    stations = stations if isinstance(stations, list) else [stations]
    period_delta1 = pd.Timedelta(hours=24)
    period_delta2 = pd.Timedelta(hours=48)

    period_delta1 = int((NOW - period_delta1).timestamp())
    period_delta2 = int((NOW - period_delta2).timestamp())

    df = data_loader.get_measurements_data(stations, period_delta2, period_delta1)

    if df.empty:
        return get_error_plot()
    if not stations:
        return get_empty_plot()

    df_by_stations = df[df[STATION_COL].isin(stations)]
    df_by_stations_and_period = df_by_stations[(df_by_stations[DATETIME_COL] >= NOW - pd.Timedelta(hours=48)) &
                                               (df_by_stations[DATETIME_COL] <= NOW - pd.Timedelta(hours=24))]
    df_by_stations_and_period[DATETIME_COL] = df_by_stations_and_period[DATETIME_COL].apply(lambda x: x + pd.Timedelta(days=2))

    return create_time_series(df_by_stations_and_period, component)


if __name__ == "__main__":
    app.run_server(debug=False, use_reloader=False, host=HOST, port=PORT)
