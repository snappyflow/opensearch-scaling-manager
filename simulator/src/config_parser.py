import os
import sys
from pathlib import Path

import yaml
from cerberus import Validator

import constants
from cluster import Cluster
from data_ingestion import DataIngestion, State
from search import SearchState
from search import Search
from search import SearchDescription
from errors import ValidationError


class Config:
    """
    Class to hold all the configurations provided in config file together
    """

    def __init__(
        self,
        stats: dict,
        states: list[dict],
        search_description: dict[dict],
        simulation_frequency_minutes: int,
        randomness_percentage: int,
    ):
        """
        Initialise the Config object
        :param states: Different states over a period of time.
        :param search_description: search description list specified in config file
        :param randomness_percentage:
        :param stats: cluster stats specified in config file
        :param simulation_frequency_minutes: interval between two simulated points
        """
        self.cluster = Cluster(**stats)
        self.simulation_frequency_minutes = simulation_frequency_minutes
        # state_object = State(90,"time",90,{},90)
        all_states = [
            State(position=state["position"],
                  time_hh_mm_ss=state["time_hh_mm_ss"],
                  ingestion_rate_gb_per_hr=state["ingestion_rate_gb_per_hr"])
            for state in states
        ]
        self.data_function = DataIngestion(all_states, randomness_percentage)
        self.search_description = [
            SearchDescription(**specs, search_type=search_type)
            for search_type, specs in search_description.items()
        ]
        self.searches = Search([
            SearchState(position=state["position"],
                        time_hh_mm_ss=state["time_hh_mm_ss"],
                        search_type=search_type,
                        count=search_count)
            for state in states
            for search_type, search_count in state["searches"].items()
        ])


def get_source_code_dir():
    """
    get the parent directory of simulator code
    :return: parent directory of simulator code
    """
    return Path(__file__).parent.resolve()


def validate_config(all_configs: dict):
    """
    Validate dictionary of configs (read from config file) against the defined schema
    :param all_configs: dictionary containing all items from config yaml
    :return: tuple containing validation state of configuration (True/False) and
             dictionary of errors
             eg. (True, {})
    """
    # Fetching the dir path to add to the schema file name
    source_code_dir: Path = get_source_code_dir()
    schema_path = os.path.join(source_code_dir, constants.SCHEMA_FILE_NAME)
    schema = eval(open(schema_path, "r").read())

    # validating config file against the schema
    validator = Validator(schema)
    return validator.validate(all_configs, schema), validator.errors


def parse_config(config_file_path: str):
    """
    Read and parse the config file into objects,
    that can work with simulator
    :param config_file_path: path of the yaml file
    :return: object of Config class
    """
    # Fetching the config file from the specified path
    fp = open(config_file_path, "r")

    # Error handling mechanism for incompletely filled config file
    try:
        # Loading the config file content to dictionary to validate
        all_configs = yaml.safe_load(fp.read())
    except Exception as e:
        fp.close()
        raise ValidationError("error reading config file - " + str(e))

    fp.close()

    # Perform Validation of the config file
    is_valid, errors = validate_config(all_configs)

    if not is_valid:
        raise ValidationError("Error validating config file - " + str(errors))

    # Extract the configurations from the file to form Config object
    simulation_frequency_minutes = all_configs.pop(
        constants.SIMULATION_FREQUENCY_MINUTES
    )
    randomness_percentage = all_configs.pop(
        constants.DATA_INGESTION_RANDOMNESS_PERCENTAGE
    )
    states = all_configs.pop(constants.STATES)
    search_description = all_configs.pop(constants.SEARCH_DESCRIPTION)
    stats = all_configs
    config = Config(
        stats,
        states,
        search_description,
        simulation_frequency_minutes,
        randomness_percentage,
    )
    return config
