import sys
import yaml
import os
from pathlib import Path

from unittest import expectedFailure
import yaml
import pytest

current_file = Path(__file__).parent.parent.resolve()
test_path = str(os.path.join(str(current_file), "tests", "config_test"))
path = os.path.join(str(current_file), "src")
sys.path.insert(0, path)

from config_parser import validate_config, parse_config
from config_parser import Config
from errors import ValidationError
from constants import CONFIG_FILE_PATH
import constants as const


def test_validate_config():
    """Checks the config file with all the required datas"""
    with open(os.path.join(test_path, "config_1_P.yaml"), "r") as file:
        is_valid, errors = validate_config(yaml.safe_load(file))
        assert is_valid == True
        assert errors == {}


def test_validate_config_without_searches():
    """Validates if config has search_description field in it."""
    with open(os.path.join(test_path, "config_2_P.yaml"), "r") as file:
        is_valid, errors = validate_config(yaml.safe_load(file))
        assert is_valid == True
        assert errors == {}


def test_validate_config_missing_parameter():
    """Validates if config has missing parameters in it"""
    with open(os.path.join(test_path, "config_4_F.yaml"), "r") as file:
        is_valid, errors = validate_config(yaml.safe_load(file))
        assert is_valid == False
        assert errors == {
            "index_clean_up_age_days": ["required field"],
            "index_roll_over_size_gb": ["required field"],
        }


def test_validate_config_invalid_data_type():
    """Checks if the config has a valid data type"""
    with open(os.path.join(test_path, "config_5_F.yaml"), "r") as file:
        is_valid, errors = validate_config(yaml.safe_load(file))
        assert is_valid == False
        assert errors == {
            "states": [{3: [{"ingestion_rate_gb_per_hr": ["must be of number type"]}]}]
        }


def test_validate_config_missing_nested_key():
    """Validates config against the list of dictionary in search_description with schema"""
    with open(os.path.join(test_path, "config_6_F.yaml"), "r") as file:
        is_valid, errors = validate_config(yaml.safe_load(file))
        assert is_valid == False
        assert errors == {
            "search_description": [
                {"simple": [{"heap_load_percent": ["required field"]}]}
            ]
        }


def test_parse_config():
    """Checks if it's a valid config and place the file in the simulator/src/main and return"""
    fp = open(os.path.join(test_path, "config_1_P.yaml"), "r")
    all_configs = yaml.safe_load(fp.read())

    simulation_frequency_minutes = all_configs.pop(const.SIMULATION_FREQUENCY_MINUTES)
    randomness_percentage = all_configs.pop(const.DATA_INGESTION_RANDOMNESS_PERCENTAGE)
    states = all_configs.pop(const.STATES)
    search_description = all_configs.pop(const.SEARCH_DESCRIPTION)
    stats = all_configs

    expected_config = Config(
        stats,
        states,
        search_description,
        simulation_frequency_minutes,
        randomness_percentage,
    )
    config = parse_config(os.path.join(test_path, "config_1_P.yaml"))
    assert (
        expected_config.simulation_frequency_minutes
        == config.simulation_frequency_minutes
    )
    assert expected_config.randomness_percentage == config.randomness_percentage
    assert expected_config.states == config.states
    assert expected_config.stats == config.stats
    assert expected_config.search_description["simple"].search_stat.__getattribute__("cpu_load_percent") == config.search_description["simple"].search_stat.__getattribute__("cpu_load_percent")
    assert expected_config.search_description["simple"].search_stat.__getattribute__("memory_load_percent") == config.search_description["simple"].search_stat.__getattribute__("memory_load_percent")
    assert expected_config.search_description["simple"].search_stat.__getattribute__("heap_load_percent") == config.search_description["simple"].search_stat.__getattribute__("heap_load_percent")
    assert expected_config.search_description["medium"].search_stat.__getattribute__("cpu_load_percent") == config.search_description["medium"].search_stat.__getattribute__("cpu_load_percent")
    assert expected_config.search_description["medium"].search_stat.__getattribute__("memory_load_percent") == config.search_description["medium"].search_stat.__getattribute__("memory_load_percent")
    assert expected_config.search_description["medium"].search_stat.__getattribute__("heap_load_percent") == config.search_description["medium"].search_stat.__getattribute__("heap_load_percent")
    assert expected_config.search_description["complex"].search_stat.__getattribute__("cpu_load_percent") == config.search_description["complex"].search_stat.__getattribute__("cpu_load_percent")
    assert expected_config.search_description["complex"].search_stat.__getattribute__("memory_load_percent") == config.search_description["complex"].search_stat.__getattribute__("memory_load_percent")
    assert expected_config.search_description["complex"].search_stat.__getattribute__("heap_load_percent") == config.search_description["complex"].search_stat.__getattribute__("heap_load_percent")


def test_parse_config_error_reading_config():
    """Checks the config is complete or not"""
    with pytest.raises(ValidationError) as e:
        parse_config(os.path.join(test_path, "config_5_F.yaml"))
        assert "error reading config file - " == e


def test_parse_config_validate_error():
    """If required field is not there in config and dont place in src path"""
    with pytest.raises(ValidationError) as e:
        parse_config(os.path.join(test_path, "config_4_F.yaml"))
        assert "Error validating config file - " == e
