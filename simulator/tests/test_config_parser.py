import sys
import yaml
import os
from pathlib import Path

current_file = Path(__file__).parent.parent.resolve()
test_path = str(os.path.join(str(current_file),"tests","config_test"))
path = os.path.join(str(current_file),"src")
sys.path.insert(0, path)

from config_parser import validate_config


def test_validate_config():
    with open(os.path.join(test_path,"config_1_P.yaml"),"r") as file:
        is_valid,errors = validate_config(yaml.safe_load(file))
        assert is_valid == True
        assert errors == {}


def test_validate_config_without_searches():
    with open(os.path.join(test_path,"config_2_P.yaml"),"r") as file:
        is_valid,errors = validate_config(yaml.safe_load(file))
        assert is_valid == True
        assert errors == {}


def test_validate_config_without_data_ingestion():
    with open(os.path.join(test_path,"config_3_P.yaml"),"r") as file:
        is_valid,errors = validate_config(yaml.safe_load(file))
        assert is_valid == True
        assert errors == {}


def test_validate_config_missing_parameter():
    with open(os.path.join(test_path,"config_4_F.yaml"),"r") as file:
        is_valid,errors = validate_config(yaml.safe_load(file))
        assert is_valid == False
        assert errors == {
            "index_clean_up_age_days": ["required field"],
            "index_roll_over_size_gb": ["required field"],
        }


def test_validate_config_invalid_data_type():
    with open(os.path.join(test_path,"config_5_F.yaml"),"r") as file:
        is_valid,errors = validate_config(yaml.safe_load(file))
        assert is_valid == False
        assert errors == {
            "data_ingestion": [
                {
                    "states": [
                        {3: [{"ingestion_rate_gb_per_hr": ["must be of number type"]}]}
                    ]
                }
            ]
        }


def test_validate_config_missing_nested_key():
    with open(os.path.join(test_path,"config_6_F.yaml"),"r") as file:
        is_valid,errors = validate_config(yaml.safe_load(file))
        assert is_valid == False
        assert errors == {"searches": [{0: [{"probability": ["required field"]}]}]}
