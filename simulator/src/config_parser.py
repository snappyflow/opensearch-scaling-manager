import os
import sys
from pathlib import Path

import yaml
from cerberus import Validator

import constants
from cluster import Cluster
from data_ingestion import DataIngestion, State
from search import Search
from errors import ValidationError


class Config:
    """
    Class to hold all the configurations provided in config file together
    """
    def __init__(
        self,
        stats: dict,
        data_ingestion: dict[dict],
        searches: list[dict],
        simulation_frequency_minutes: int,
    ):
        """
        Initialise the Config object
        :param stats: cluster stats specified in config file
        :param data_ingestion: data ingestion mapping specified in config file
        :param searches: searches list specified in config file
        :param simulation_frequency_minutes: interval between two simulated points
        """
        self.cluster = Cluster(**stats)
        self.simulation_frequency_minutes = simulation_frequency_minutes
        all_states = [State(**state) for state in data_ingestion.get(constants.DATA_INGESTION_STATES)]
        randomness_percentage = data_ingestion.get(constants.DATA_INGESTION_RANDOMNESS_PERCENTAGE)
        self.data_function = DataIngestion(all_states, randomness_percentage)
        self.searches = [Search(**specs) for specs in searches]


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
        raise ValidationError('error reading config file - ' + str(e))

    fp.close()

    source_code_dir: Path = get_source_code_dir()

    # Perform Validation of the config file
    is_valid, errors = validate_config(all_configs)

    if not is_valid:
        raise ValidationError('Error validating config file - ' + str(errors))

    # Extract the configurations from the file to form Config object
    simulation_frequency_minutes = all_configs.pop(
        constants.SIMULATION_FREQUENCY_MINUTES
    )
    data_ingestion = all_configs.pop(constants.DATA_INGESTION)
    searches = all_configs.pop(constants.SEARCHES)
    stats = all_configs

    return Config(stats, data_ingestion, searches, simulation_frequency_minutes)
