import json
import sys
import os
from pathlib import Path

current_file = Path(__file__).parent.parent.resolve()
path = os.path.join(str(current_file), "src")
sys.path.insert(0, path)

import app


def test_violated_count():
    """Validates violated_count function with proper stat_name, duration, threshold with status code"""
    stat_name = "cpu"
    duration = 90
    threshold = 100.0
    expected_response = True
    response = app.violated_count(stat_name, duration, threshold)
    key = "ViolatedCount"
    if key in response.get_data(as_text=True):
        actual_response = True
    assert expected_response == actual_response
    assert response.status_code == 200


def test_violated_count_with_improper_stat_name():
    """Validates violated_count function with proper stat_name, duration, threshold and check errors"""
    stat_name = "ram"
    duration = 90
    threshold = 100.0
    response = app.violated_count(stat_name, duration, threshold)
    assert response.status_code == 404
    assert response.get_data(as_text=True) == f"stat not found - {stat_name}"


def test_violated_count_with_improper_duration():
    """Validates violated_count function with proper stat_name, duration, threshold and check errors"""
    stat_name = "cpu"
    duration = 10000
    threshold = 100.0
    response = app.violated_count(stat_name, duration, threshold)
    assert response.status_code == 400
    assert response.get_data(as_text=True) == '"Not enough Data points"'


def test_average():
    """Validates average function with proper stat_name, duration and the keys of endpoint"""
    stat_name = "cpu"
    duration = 90
    response = app.average(stat_name, duration)
    required_output_keys = {
        "avg",
        "max",
        "min",
    }
    
    #convert response into dictionary
    response_dictionary = json.loads(response.get_data(as_text=True))
    
    #assert and compare required_output_keys with response dictionary 
    assert required_output_keys == response_dictionary.keys() 
    assert response.status_code == 200


def test_average_with_improper_stat_name():
    """Validates average function with improper stat_name"""
    stat_name = "ram"
    duration = 90
    response = app.average(stat_name, duration)
    assert response.status_code == 404
    assert response.get_data(as_text=True) == f"stat not found - {stat_name}"
    

def test_average_with_improper_duration():
    """Validates average function with improper duration"""
    stat_name = "cpu"
    duration = 0
    response = app.average(stat_name, duration)
    assert response.status_code == 400
    assert response.get_data(as_text=True) == '"Not enough Data points"'

def test_current():
    """Validates current function and checks for the requested API endpoints returned"""
    stat_name = "cpu"
    response = app.current(stat_name)
    expected_response = True
    key = "current"
    if key in response.get_data(as_text=True):
        actual_response = True
    assert response.status_code == 200
    assert expected_response == actual_response


def test_current_with_improper_stat_name():
    """Validates current function with improper stat_name"""
    stat_name = "ram"
    response = app.current(stat_name)
    assert response.status_code == 404
    assert response.get_data(as_text=True) == f"stat not found - {stat_name}"


def test_current_all():
    """Validates current_all function and checks for the requested API endpoints returned"""
    response = app.current_all()
    required_output_keys = {
        "NumNodes",
        "ClusterStatus",
        "NumActiveShards",
        "NumActivePrimaryShards",
        "NumInitializingShards",
        "NumUnassignedShards",
        "NumRelocatingShards",
        "NumMasterNodes",
        "NumActiveDataNodes",
    }
    #convert response into dictionary
    response_dictionary = json.loads(response.get_data(as_text=True))
    assert response.status_code == 200
    #assert and perform set operation to compare response dictionary 
    assert required_output_keys == response_dictionary.keys() 
