import pandas as pd
import numpy as np
import os
import matplotlib.pyplot as plt
import fbprophet
import matplotlib.dates as mdates

from sklearn.impute import SimpleImputer
from sklearn.preprocessing import MinMaxScaler
from sklearn.preprocessing import LabelEncoder
from sklearn.metrics import mean_squared_error
from keras.models import Sequential
from keras.layers import Dense
from keras.layers import LSTM
from keras.layers import RepeatVector
from keras.layers import TimeDistributed

class ModelManager:
    def __init__(self, input_timestemps=504, output_timestemps=72, target_params=['CO', 'NO2', 'NO', 'PM10', 'PM2.5']):
        self.input_timestemps=input_timestemps
        self.output_timestemps=output_timestemps

        self.target_params = target_params
        self.target_features = target_params.__len__()

        self.X_samples = None
        self.y_samples = None
        self.input_features = None

        self.model = None
        self.learning_history = None

        self.scaler = None
        self.scaler_y = None              
    
    def predict_series(self, input_series):
        if self.model and self.scaler and self.scaler_y:
            raw_prediction = self.model.predict(input_series)
            return self.scaler_y.inverse_trasform(raw_prediction)                                  
        else:
            print('Initialize model first with using methods init_model or load_model')

    def _init_samples(self, data):
        variables_names = data.columns.to_list()
        number_of_columns = data.shape[1]

        col_values, col_names = [], []
        y_values, y_col_names = [], []
        for i in range(self.input_timestemps, 0, -1):
            col_values.append(data.shift(i))
            col_names += [('{col_name} (t-{step})'.format(col_name=col_name, step=i)) for col_name in variables_names]
        for i in range(0, self.output_timestemps, 1):
            y_values.append(data[self.target_features].shift(-i))
            y_col_names += [('{col_name} (t+{step})'.format(col_name=col_name, step=i)) for col_name in variables_names if col_name in target_features]

        X_samples = pd.concat(col_values, axis=1)
        X_samples.columns = col_names
        #X_samples.dropna(inplace=True)

        y_samples = pd.concat(y_values, axis=1)
        y_samples.columns = y_col_names
        #y_samples.dropna(inplace=True)

        self.X_samples, self.y_samples = X_samples, y_samples

    def _prepare_data_for_learning(self, df_station, df_pr, input_timestemps=504, output_timestemps=72):
        #  Make timestap index 
        df_station['datetime'] = pd.to_datetime(df_station['Дата и время'], format='%d/%m/%Y %H:%M:%S')
        df_station.set_index('datetime', inplace=True)

        df_pr.reindex(df_station.index)

        # Merge station data and profilemers
        X = df_station.join(df_pr)

        # Drop features if all values are missing
        df_t = X.isnull().all()
        drop_features = df_t[df_t == True].index.to_list()

        # Drop features if there are only unique value
        for feature in X.columns:
            if len(X[feature].unique()) < 2:
                drop_features.append(feature)

        # Drop depricated time data
        depricated_datetime = ['Дата и время', 'data time']
        for feature in depricated_datetime:
            if feature in X.columns:
                drop_features.append(feature)

        print('Features to drop:', drop_features)
        X = X.drop(columns=drop_features)
        X.dropna(inplace=True)

        # Drop features if there are only unique value
        for feature in X.columns:
            if len(X[feature].unique()) < 2:
                drop_features.append(feature)

        
        TARGET_PARAMS = ['CO', 'NO2', 'NO', 'PM10', 'PM2.5']
        input_features = X.columns.shape[0]
        target_features = TARGET_PARAMS.__len__()

        # Make samples for supervised learning from initial dataframe 
        # 1 measurment as input, 1 hour as ouptut
        _init_samples(X, TARGET_PARAMS, n_input=input_timestemps, n_output=output_timestemps)

        nan_rows = X_samples.isnull().any(axis=1) | y_samples.isnull().any(axis=1)
        X_samples, y_samples = X_samples.drop(X_samples[nan_rows].index), y_samples.drop(y_samples[nan_rows].index)

        # normalize features

        self.scaler = MinMaxScaler(feature_range=(0, 1))
        X_samples = self.scaler.fit_transform(X_samples)

        columns_names_y = y_samples.columns.to_list()
        self.scaler_y = MinMaxScaler(feature_range=(0, 1))
        y_samples = self.scaler_y.fit_transform(y_samples)

        X_samples = X_samples.reshape(X_samples.shape[0], input_timestemps, input_features)
        y_samples = y_samples.reshape(y_samples.shape[0], output_timestemps, target_features)

        return X_samples, y_samples
    

    def fit_model(self):
        if self.X_samples and self.y_samples:
            X_train, X_test, y_train, y_test = train_test_split(X_samples, y_samples, test_size=0.25) 
            # define model
            model = Sequential()
            model.add(LSTM(200, activation='relu', input_shape=(self.input_timestemps, self.input_features)))
            model.add(RepeatVector(self.output_timestemps))
            model.add(LSTM(200, activation='relu', return_sequences=True))
            model.add(TimeDistributed(Dense(self.target_features)))
            model.compile(optimizer='adam', loss='mse')
            # fit model
            history = model.fit(X_train, y_train, epochs=100, batch_size=128, validation_data=(X_test, y_test), verbose=1, shuffle=False)
            self.model = model
            self.learning_history = history
        else:
            print('Initialize samples first')