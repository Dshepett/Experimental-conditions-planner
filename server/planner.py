from json import load
from joblib import dump, load
import numpy as np
from sklearn.pipeline import Pipeline, make_pipeline
from sklearn.preprocessing import StandardScaler
from sklearn.ensemble import RandomForestClassifier, RandomForestRegressor
import pandas as pd
from os.path import exists


class Planner:
    path_to_data: str # путь до csv файла с тестовой выборкой

    path_to_models: str # путь до объектов для вычисления условий синтеза

    # объекты для вычисления условий синтеза золи
    temperature_calculator: Pipeline

    time_calculator: Pipeline

    c_acid_calculator: Pipeline

    c_ti_calculator: Pipeline

    acid_caclulator: Pipeline

    tretmant_calculator: Pipeline

    treatment_calculator: Pipeline

    scaler: StandardScaler

    # подгружает объекты для вычисления условий в программу
    def build(self, path_to_data, path_to_models):
        self.path_to_data = path_to_data
        self.path_to_models = path_to_models
        if not self.load():
            print('Fail')
            self.generate()

    # загрузка объектов для вычисения условий из локального хранилища
    def load(self):
        print('Trying to load models from file...')
        if not exists(f'{self.path_to_models}/acid_calculator.joblib'):
            return False
        if not exists(f'{self.path_to_models}/tempeareture_calculator.joblib'):
            return False
        if not exists(f'{self.path_to_models}/time_calculator.joblib'):
            return False
        if not exists(f'{self.path_to_models}/c_ti_calculator.joblib'):
            return False
        if not exists(f'{self.path_to_models}/treatment_calculator.joblib'):
            return False
        if not exists(f'{self.path_to_models}/с_acid_calculator.joblib'):
            return False
        acid_calculator = load(f'{self.path_to_models}/acid_calculator.joblib')
        self.acid_caclulator = acid_calculator
        temperature_calculator = load(
            f'{self.path_to_models}/tempeareture_calculator.joblib')
        self.temperature_calculator = temperature_calculator
        time_calculator = load(f'{self.path_to_models}/time_calculator.joblib')
        self.time_calculator = time_calculator
        c_ti_calculator = load(f'{self.path_to_models}/c_ti_calculator.joblib')
        self.c_ti_calculator = c_ti_calculator
        treatment_calculator = load(
            f'{self.path_to_models}/treatment_calculator.joblib')
        self.treatment_calculator = treatment_calculator
        c_acid_calсulator = load(
            f'{self.path_to_models}/с_acid_calculator.joblib')
        self.c_acid_calculator = c_acid_calсulator
        print('Success!')
        return True

    # вычисление и создание объектов для вычисления услловий для синтеза на основе файла с тестовой выборкой
    def generate(self):
        data = pd.read_csv(self.path_to_data, sep=";")
        res = data.iloc[:, 6:9].to_numpy(dtype=float)  # характеристики золи
        temperature = data['Т, °С'].to_numpy(dtype=float)  # температура синтеза
        time = data['t, мин'].to_numpy(dtype=float)  # время синтеза
        # концентрация кислоты
        c_acid = data['С(HNO3 || H2SO4), моль/л'].to_numpy(dtype=float)
        # концентрация титана
        c_ti = data['с(Ti4+), моль/л'].to_numpy(dtype=float)
        acid = data['HNO3'].to_numpy(dtype=int)  # вид кислоты
        treatment = data['Ультразвуковая обработка, мин'].to_numpy(dtype=float)

        print('generating acid calculator')
        acid_calculator = Pipeline([
            ('scaler', StandardScaler()),
            ('classifier', RandomForestClassifier(n_estimators=2000, n_jobs=-1))
        ])
        acid_calculator.fit(res, acid)
        self.acid_caclulator = acid_calculator
        dump(acid_calculator, f'{self.path_to_models}/acid_calculator.joblib')

        print('generating temperature calculator')
        temperature_calulator = make_pipeline(
            StandardScaler(), RandomForestRegressor(n_estimators=3000, n_jobs=-1))
        temperature_calulator.fit(res, temperature)
        self.temperature_calculator = temperature_calulator
        dump(temperature_calulator,
             f'{self.path_to_models}/tempeareture_calculator.joblib')

        print('generating time calcuator')
        time_calulator = make_pipeline(
            StandardScaler(), RandomForestRegressor(n_estimators=10000))
        time_calulator.fit(res, time)
        self.time_calculator = time_calulator
        dump(time_calulator, f'{self.path_to_models}/time_calculator.joblib')

        print('generating c_ti caclculator')
        c_ti_calulator = make_pipeline(
            StandardScaler(), RandomForestRegressor(n_estimators=10000, n_jobs=-1))
        c_ti_calulator.fit(res, c_ti)
        self.c_ti_calculator = c_ti_calulator
        dump(c_ti_calulator, f'{self.path_to_models}/c_ti_calculator.joblib')

        print('generating treatment caclculator')
        treat = list()
        for i in treatment:
            if i == 0:
                treat.append(0)
            elif i == 0.5:
                treat.append(1)
            elif i == 1:
                treat.append(2)
            elif i == 1.5:
                treat.append(3)
            else:
                treat.append(4)
        treat = np.array(treat)
        treatment_calculator = make_pipeline(
            StandardScaler(), RandomForestClassifier(n_estimators=3000, n_jobs=-1))
        treatment_calculator.fit(res, treat)
        dump(treatment_calculator,
             f'{self.path_to_models}/treatment_calculator.joblib')

        print('generating c_acid calculator')
        res_extra = data.iloc[:, 5:9].to_numpy(dtype=float)
        c_acid_calсulator = make_pipeline(
            StandardScaler(), RandomForestRegressor(n_estimators=12000, n_jobs=-1))
        c_acid_calсulator.fit(res_extra, c_acid)
        self.c_acid_calculator = c_acid_calсulator
        dump(c_acid_calсulator,
             f'{self.path_to_models}/с_acid_calculator.joblib')

    # вычисляет условия для синтеза золей
    def calculate(self, zoles):
        res = list()
        for zole in zoles:
            temperature = self.temperature_calculator.predict([zole])[0]
            time = self.time_calculator.predict([zole])[0]
            c_ti = self.c_ti_calculator.predict([zole])[0]
            acid = self.acid_caclulator.predict([zole])[0]
            c_acid = self.c_acid_calculator.predict(np.array([[acid, zole[0], zole[1], zole[2]]]))[0]
            treatment = self.treatment_calculator.predict([zole])[0]
            curr = [temperature, time, c_acid, c_ti, acid, treatment]            
            res.append(curr)
        return np.array(res)
