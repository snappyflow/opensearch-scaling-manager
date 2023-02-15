# This is the main module of OpenSearch Cluster simulator
# The intent of this module is to simulate the behaviour of a cluster
# under varies loading conditions

import os
import sys

import constants
from config_parser import parse_config, get_source_code_dir
from open_search_simulator import Simulator
from cluster import Cluster
from data_ingestion import State, DataIngestion
from plotter import plot_data_points


if __name__ == '__main__':
    configs = parse_config(os.path.join(get_source_code_dir(), constants.CONFIG_FILE_PATH))
    all_states = [State(**state) for state in configs.data_ingestion.get(constants.DATA_INGESTION_STATES)]
    randomness_percentage = configs.data_ingestion.get(constants.DATA_INGESTION_RANDOMNESS_PERCENTAGE)

    data_function = DataIngestion(all_states, randomness_percentage)

    cluster = Cluster(**configs.stats)

    sim = Simulator(cluster, data_function, configs.searches, configs.data_generation_interval_minutes)
    cluster_objects = sim.run(24*60)
    plot_data_points(cluster_objects)
