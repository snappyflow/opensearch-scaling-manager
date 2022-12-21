import os
import copy
import random
from datetime import datetime, timedelta

import constants
from config_parser import get_source_code_dir
from cluster import Cluster
from data_ingestion import DataIngestion
from search import Search


class Simulator:
    """
    Runs simulation on a passed cluster object
    Takes care of:
        - triggering of events
        - altering the states of nodes and cluster based on events
    """
    def __init__(
            self,
            cluster: Cluster,
            data_ingestion: DataIngestion,
            searches: list[Search],
            frequency_minutes: int,
            elapsed_time_minutes: int = 0,
    ):
        """
        Initialize the Simulator object
        :param cluster: cluster object to run the simulations on
        :param data_ingestion: data ingestion object to simulate on cluster
        :param searches: search query objects to simulate on cluster
        :param frequency_minutes: the difference between two resultant simulated points
        :param elapsed_time_minutes: if provided, the cluster
            represents the state after the elapsed minutes
        """
        self.cluster = cluster
        self.data_ingestion = data_ingestion
        self.searches = searches
        self.elapsed_time_minutes = elapsed_time_minutes
        self.frequency_minutes = frequency_minutes

    def aggregate_data(
            self,
            duration_minutes,
            start_time_hh_mm_ss: str = '00_00_00'
    ):
        # first collect all data aggregation events
        x, y = self.data_ingestion.aggregate_data(start_time_hh_mm_ss, duration_minutes, self.frequency_minutes)
        return x, y

    def cpu_used_for_ingestion(self, ingestion):
        return min(ingestion / self.cluster.total_nodes_count * random.randrange(1, 15) / 100 * 100, 100)

    def memory_used_for_ingestion(self, ingestion):
        return min(ingestion / self.cluster.total_nodes_count * random.randrange(5, 12) / 100 * 100, 100)

    def cluster_state_for_ingestion(self, ingestion):
        if ingestion < constants.HIGH_INGESTION_RATE_GB_PER_HOUR:
            return random.choice([constants.CLUSTER_STATE_GREEN] * 20 + [constants.CLUSTER_STATE_YELLOW])
        if self.cluster.status == constants.CLUSTER_STATE_RED:
            return random.choice([constants.CLUSTER_STATE_YELLOW] + [constants.CLUSTER_STATE_RED]*5)
        return random.choice(
            [constants.CLUSTER_STATE_GREEN] * 20 + [constants.CLUSTER_STATE_YELLOW] * 10 + [constants.CLUSTER_STATE_RED])

    def run(self, duration_minutes):
        resultant_cluster_objects = []
        data_x, data_y = self.aggregate_data(duration_minutes)
        now = datetime.now()
        date_obj = now - timedelta(
                hours=now.hour,
                minutes=now.minute,
                seconds=now.second,
                microseconds=now.microsecond
            )
        for instantaneous_data_ingestion_rate in data_y:
            self.cluster._ingestion_rate = instantaneous_data_ingestion_rate
            self.cluster.cpu_usage_percent = self.cpu_used_for_ingestion(instantaneous_data_ingestion_rate)
            self.cluster.memory_usage_percent = self.memory_used_for_ingestion(instantaneous_data_ingestion_rate)
            self.cluster.status = self.cluster_state_for_ingestion(instantaneous_data_ingestion_rate)
            # Todo: simulate effect on remaining cluster parameters 
            date_time = date_obj + timedelta(minutes=self.elapsed_time_minutes)
            self.cluster.date_time = date_time
            resultant_cluster_objects.append(copy.deepcopy(self.cluster))
            self.elapsed_time_minutes += self.frequency_minutes
        return resultant_cluster_objects

    @staticmethod
    def create_provisioning_lock():
        lock_file_path = os.path.join(get_source_code_dir(), constants.PROVISION_LOCK_FILE_NAME)
        expiry_time = datetime.now() + timedelta(seconds=random.randint(15, 65))
        with open(lock_file_path, 'w') as file_handler:
            file_handler.write(expiry_time.isoformat())
        return expiry_time.isoformat()

    @staticmethod
    def is_provision_in_progress():
        lock_file_path = os.path.join(get_source_code_dir(), constants.PROVISION_LOCK_FILE_NAME)
        if os.path.exists(lock_file_path):
            with open(lock_file_path, 'r') as file_handler:
                if datetime.now() > datetime.fromisoformat(file_handler.read()):
                    file_handler.close()
                    Simulator.remove_provisioning_lock()
                    return False
                else:
                    return True
        return False

    @staticmethod
    def remove_provisioning_lock():
        lock_file_path = os.path.join(get_source_code_dir(), constants.PROVISION_LOCK_FILE_NAME)
        if os.path.exists(lock_file_path):
            os.remove(lock_file_path)
