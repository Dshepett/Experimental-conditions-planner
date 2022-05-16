from flask import Flask, make_response, request
from planner import Planner

from utils import convert_to_array, convert_to_json, is_valid

app = Flask(__name__)
app.config['JSON_SORT_KEYS'] = False

planner : Planner

@app.route('/calculate', methods=['POST'])
def calculate():
    json_data=request.get_json()    
    if not is_valid(json_data):
        return make_response({"status":404,"message":"incorrect data"},404)
    data = json_data['characteristics']
    zoles = convert_to_array(data)    
    result = planner.calculate(zoles)
    return make_response(convert_to_json(result),200)


if __name__ == "__main__":  
    import os
    from dotenv import load_dotenv

    # загрузка переменных из .env
    load_dotenv()

    port = int(os.getenv('PORT'))
    host = os.getenv('HOST')
    path_to_data = os.getenv('PATH_TO_DATA')
    path_to_models = os.getenv('PATH_TO_MODELS')
    planner = Planner()
    planner.build(path_to_data=path_to_data, path_to_models=path_to_models)          
    app.run(host = host,port=port)    