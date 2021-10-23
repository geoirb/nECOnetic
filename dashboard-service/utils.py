import plotly.graph_objects as go


def get_error_plot():
    """
    Возвращает фон с надписью об ошибке

    Заглушка для использования в случаях, когда из базы не пришли нужные данные
    :return: фон с надписью об ошибке
    """
    fig = go.Figure().add_annotation(x=2, y=2, text="No Data to Display",
                                     font=dict(family="sans serif", size=25, color="rgba(0, 0, 0, 1)"),
                                     showarrow=False, yshift=10)
    fig.update_layout(paper_bgcolor="rgba(0, 0, 0, 0)",
                      plot_bgcolor="rgba(0, 0, 0, 0)",
                      xaxis={'visible': False},
                      yaxis={'visible': False}
                      )

    return fig


def get_empty_plot():
    """
    Возвращает пустой график

    Для использования в случаях, когда не выбраны станции для отображения
    графиков распределения компонентов по времени
    :return: пустой график
    """
    fig = go.Figure()
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
