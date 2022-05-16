from flask import jsonify
import numpy as np

def is_valid(request_data):
    if 'characteristics' not in request_data:        
        return False    
    data = request_data['characteristics']
    for zole in data:                
        for characteristic in ['size', 'consistance', 'stability']:
            if characteristic in zole:
                if type(zole[characteristic]) != int and  type(zole[characteristic]) != float:                                   
                    return False
                if characteristic=='consistance':
                    if zole[characteristic]<0 or zole[characteristic]>100:
                        return False
                else:
                    if zole[characteristic]<0:
                        return False
            else:
                return False
    return True

def convert_to_array(data):
    zoles = list()
    for i in data:        
        zole =[float(i['size']), float(i['consistance']), float(i['stability'])]        
        zoles.append(zole)
        
    return np.array(zoles)

def convert_to_json(result):  
    conditions=[]
    response_data = {"status":200, "message":"OK"}          
    for i in result:
        zole_cond = {
            "temperature":i[0],
            "time":i[1],
            "c_acid":i[2],
            "c_ti":i[3],
            "acid":"hno3" if i[4]==1 else "h2so4",
            "treatment":i[5]
        }
        conditions.append(zole_cond)
    response_data["conditions"]=conditions
    return jsonify(response_data)