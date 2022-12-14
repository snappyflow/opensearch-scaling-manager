import sys
import yaml
from pathlib import Path

current_file: Path = Path(__file__).parent.parent.resolve()
path = str(current_file) + "\src"
print(path)
sys.path.insert(0, path)

from config_parser import validate_config


def test_validate_config():
    is_valid, errors = validate_config(
        yaml.safe_load(open("./config_test/config_1_P.yaml", "r"))
    )
    assert is_valid == True
    assert errors == {}


def test_validate_config_without_searches():
    is_valid, errors = validate_config(
        yaml.safe_load(open("./config_test/config_2_P.yaml", "r"))
    )
    assert is_valid == True
    assert errors == {}


def test_validate_config_without_data_ingestion():
    is_valid, errors = validate_config(
        yaml.safe_load(open("./config_test/config_3_P.yaml", "r"))
    )
    assert is_valid == True
    assert errors == {}


def test_validate_config_missing_parameter():
    is_valid, errors = validate_config(
        yaml.safe_load(open("./config_test/config_4_F.yaml", "r"))
    )
    assert is_valid == False
    assert errors == {
        "index_clean_up_age_days": ["required field"],
        "index_roll_over_size_gb": ["required field"],
    }


def test_validate_config_invalid_data_type():
    is_valid, errors = validate_config(
        yaml.safe_load(open("./config_test/config_5_F.yaml", "r"))
    )
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
    is_valid, errors = validate_config(
        yaml.safe_load(open("./config_test/config_6_F.yaml", "r"))
    )
    assert is_valid == False
    assert errors == {"searches": [{0: [{"probability": ["required field"]}]}]}
