import dash_bootstrap_components as dbc
import dash_core_components as dcc
import dash_html_components as html

import data_loader

from constants import COMPONENTS, DEFAULT_COMPONENT, DEFAULT_STATION, DISABLED_COMPONENTS, PERIODS

stations_df = data_loader.get_stations()

if not stations_df.empty:
    stations = list(stations_df['name'].unique())
else:
    stations = []


main_layout = html.Div(
    [
        dbc.Row([dbc.Col(html.Div([html.H1('MoscowEcoMap')], className='header-1'),
                         align='end', xs=12, sm=12, md=12, lg=12, xl=12)]),
        dbc.Row(dbc.Col(html.Div([html.H5('Экологическая карта города Москвы')], className='header-2'),
                        align='end', xs=12, sm=12, md=12, lg=12, xl=12)),
        dbc.Row(
            [
                dbc.Col(html.Div([dcc.Graph(
                    id='map',
                )], className='map'), width={'size': 'auto'}, xs=12, sm=12, md=12, lg=6, xl=6),
                dbc.Col([
                    html.Div([
                        html.Div([
                            html.Div([
                                html.Div([
                                    dcc.Dropdown(
                                        id='stations',
                                        options=[{'label': s, 'value': s} for s in stations],
                                        value=DEFAULT_STATION,
                                        multi=True,
                                        clearable=False,
                                        className='stations-dropdown'
                                    )], className='stations-dropdown-container'),
                                html.Div([
                                    dcc.Dropdown(
                                        id='component',
                                        options=[{'label': c, 'value': c, 'disabled': c in DISABLED_COMPONENTS} for c in COMPONENTS],
                                        value=COMPONENTS[0] if COMPONENTS else DEFAULT_COMPONENT,
                                        multi=False,
                                        clearable=False,
                                        className='component-dropdown'
                                    )]),
                            ], className='dropdowns'),
                            dbc.FormGroup([
                                dcc.RadioItems(
                                    id='period',
                                    options=[{'label': p, 'value': p} for p in PERIODS],
                                    value='day',
                                    labelClassName='mr-2'
                                )], className='period-radioitems'
                            )], className='selectors'),
                        html.Div(
                            [
                                html.Div(
                                    [
                                        html.H5('Данные с детекторов', className='components-past-header'),
                                        html.Div([dcc.Graph(id='components-by-period')],
                                                 className='components-past-plot'),
                                    ], className='components-past'
                                ),

                                html.Div(
                                    [
                                        html.H5('Прогноз на сутки', className='components-future-header'),
                                        html.Div([
                                            dcc.Graph(id='components-by-future-period')
                                        ], className='components-future-plot')
                                    ], className='components-future'
                                )
                            ],
                            className='components-plots'),
                    ], className='right-column')], width={'size': 12}, xs=12, sm=12, md=12, lg=6, xl=6
                ),
            ], align='center'
        ),
    ], className='wrapper')
