import os
import copy
import random
from datetime import datetime, timedelta

import constants
from config_parser import get_source_code_dir
from cluster import Cluster
from data_ingestion import DataIngestion
from search import SearchDescription, Search
from index_addittion import IndexAddition
import time


def timeit(func):
    def inner(*args, **kwargs):
        time_start = time.time()
        ret = func(*args, **kwargs)
        time_end = time.time()
        total_time = time_end - time_start
        print("time taken for the function :", func.__name__, " is: ", total_time)
        return ret

    return inner


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
        search_description: SearchDescription,
        searches: Search,
        frequency_minutes: int,
        index_addition: IndexAddition,
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
        self.search_description = search_description
        self.elapsed_time_minutes = elapsed_time_minutes
        self.frequency_minutes = frequency_minutes
        self.searches = searches
        self.index_addition = index_addition
        self.simulated_data_rates = []
        self.simulated_search_rates = {}
        self.total_simulation_minutes = 0
        self.index_added_list = []
        self.current_index_count = 0

    def aggregate_data(self, duration_minutes, start_time_hh_mm_ss: str = "00_00_00"):
        # first collect all data aggregation events
        x, y = self.data_ingestion.aggregate_data(
            start_time_hh_mm_ss, duration_minutes, self.frequency_minutes
        )
        return x, y

    def aggregate_data_searches(
        self, duration_minutes, start_time_hh_mm_ss: str = "00_00_00"
    ):
        # first collect all data aggregation events
        x, y = self.searches.aggregate_data(
            start_time_hh_mm_ss, duration_minutes, self.frequency_minutes
        )
        return x, y

    def aggregate_index_addition(self, start_time_hh_mm_ss="00_00_00"):
        y2 = self.index_addition.aggregate_index_addition(
            self.cluster.index_count, start_time_hh_mm_ss, self.frequency_minutes
        )
        return y2

    def add_index_to_cluster(self, ind, time):
        if self.cluster.index_count >= self.index_added_list[ind]:
            return

        index_count_to_add = self.index_added_list[ind] - self.cluster.index_count

        for index in range(index_count_to_add):
            self.cluster.create_index(
                self.cluster.primary_shards_per_index,
                self.cluster.replica_shards_per_index,
                time,
            )
            index_list = []
            index_list.append(self.cluster.indices[-1])
            self.cluster.allocate_shards_to_node(index_list)
            self.cluster.index_count += 1

    def compute_cpu(self, data_rate):
        if data_rate in range(0, 21):
            cpu_rate = random.uniform(5, 20)
            return round(cpu_rate, 2)

        if data_rate in range(20, 50):
            cpu_rate = random.uniform(20, 40)
            return round(cpu_rate, 2)

        if data_rate in range(50, 80):
            cpu_rate = random.uniform(40, 60)
            return round(cpu_rate, 2)

        if data_rate in range(80, 200):
            cpu_rate = random.uniform(60, 80)
            return round(cpu_rate, 2)

        if data_rate > 200:
            cpu_rate = random.uniform(80, 90)
            return round(cpu_rate, 2)

        if data_rate <= 0:
            return round(random.uniform(5, 20), 2)

        return round(random.uniform(5, 10), 2)

    def cpu_used_for_ingestion(self, ingestion, search_count, index):
        cpu_util = self.compute_cpu(int(ingestion))
        cpu_util = cpu_util * (7 / self.cluster.total_nodes_count)

        for search_type, count_array in search_count.items():
            cpu_load_percent = (
                self.search_description[search_type].search_stat.cpu_load_percent / 100
            )
            search_factor = count_array[index] * cpu_load_percent
            search_factor = search_factor * (7 / self.cluster.total_nodes_count)
            cpu_util += search_factor

        return min(cpu_util, 100)

    def memory_used_for_ingestion(self, ingestion, search_count, index):
        memory_util = (
            ingestion/self.cluster.total_nodes_count* random.randrange(5, 12)/100 * 100
        )
        for search_type, count_array in search_count.items():
            memory_load_percent = (
                self.search_description[search_type].search_stat.memory_load_percent
                / 100
            )
            search_factor = count_array[index] * memory_load_percent
            search_factor = search_factor * (7 / self.cluster.total_nodes_count)
            memory_util += search_factor
        return min(memory_util, 98)

    def heap_used_for_ingestion(self, ingestion, search_count, index, memory_util):
        # heap_util = ingestion / self.cluster.total_nodes_count * random.randrange(5, 8) / 200 * 100
        heap_util = memory_util * (2 / 3)
        for search_type, count_array in search_count.items():
            heap_load_percent = (
                self.search_description[search_type].search_stat.heap_load_percent / 100
            )
            search_factor = count_array[index] * heap_load_percent
            search_factor = search_factor * (7 / self.cluster.total_nodes_count)
            heap_util += search_factor
        return min(heap_util, 100)

    def cluster_state_for_ingestion(self, ingestion):
        if ingestion < constants.HIGH_INGESTION_RATE_GB_PER_HOUR:
            return random.choice(
                [constants.CLUSTER_STATE_GREEN] * 20 + [constants.CLUSTER_STATE_YELLOW]
            )
        if self.cluster.status == constants.CLUSTER_STATE_RED:
            return random.choice(
                [constants.CLUSTER_STATE_YELLOW] + [constants.CLUSTER_STATE_RED] * 5
            )
        return random.choice(
            [constants.CLUSTER_STATE_GREEN] * 20
            + [constants.CLUSTER_STATE_YELLOW] * 10
            + [constants.CLUSTER_STATE_RED]
        )

    def disk_utilization_for_ingestion(self):
        return self.cluster.calculate_cluster_disk_size(self.current_index_count)

    def disk_util_for_index_roll_over(self, time):
        # for index in range(len(self.cluster.indices)):
        for index in range(self.current_index_count):
            index_size = self.cluster.indices[index].get_index_primary_size()
            roll_over_age = False
            # print('Time: ',self.cluster.indices[index].created_at + timedelta(hours = self.cluster.index_roll_over_hours) )
            if (
                self.cluster.indices[index].time_elapsed_last_roll_over
                + timedelta(hours=self.cluster.index_roll_over_hours)
                <= time
            ):
                roll_over_age = True
            if (
                index_size >= self.cluster.index_roll_over_size_gb or roll_over_age
            ) and not self.cluster.indices[index].rolled_over:
                if self.cluster.rolled_over_index_id != -1:
                    # Roll over index already exists
                    # Add the size of index with roll over index size
                    self.cluster.indices[
                        self.cluster.rolled_over_index_id
                    ].index_size += index_size
                    self.cluster.indices[self.cluster.rolled_over_index_id].shards[
                        0
                    ].shard_size += index_size
                    self.cluster.rolled_index_size = (
                        self.cluster.indices[self.cluster.rolled_over_index_id]
                        .shards[0]
                        .shard_size
                    )
                    # discard the shards of roll over index
                    for shard in range(len(self.cluster.indices[index].shards)):
                        self.cluster.indices[index].shards[shard].shard_size = 0
                    self.cluster.indices[index].time_elapsed_last_roll_over = time

                # If it is first roll over, discard replicas and retain primaries
                else:
                    node_id = self.cluster.get_available_node_id()
                    id = random.choice(node_id)
                    for shard in range(len(self.cluster.indices[index].shards)):
                        #  if self.cluster.indices[index].shards[shard].type == "Replica":
                        del self.cluster.indices[index].shards[0]
                        shard -= 1

                    # Merge the primaries
                    shard = self.cluster.indices[index].initialize_shards(1, 0)
                    shard[0].node_id = id
                    shard[0].index_id = self.cluster.indices[index].index_id
                    shard[0].shard_size = index_size
                    self.cluster.nodes[id].shards_on_node.append(shard[0])

                    # Add the primary size to roll over index size
                    self.cluster.indices[index].index_size += index_size
                    self.cluster.rolled_index_size += index_size

                    # mark the index is rolled over
                    self.cluster.indices[index].rolled_over = True

                    # set the index roll over id
                    self.cluster.rolled_over_index_id = self.cluster.indices[
                        index
                    ].index_id
                    self.cluster.indices[index].shards.append(shard[0])

                    # create a new index with similar configuration of rolled over index
                    self.cluster.create_index(
                        self.cluster.primary_shards_per_index,
                        self.cluster.replica_shards_per_index,
                        time,
                    )

                    # allocate the shards
                    self.cluster.allocate_shards_to_node()

        return self.cluster.calculate_cluster_disk_size(self.current_index_count)

    def distribute_load(self, ingestion):
        """
        The function will select an Index and distribute
        data in an arbitrary fashion.
        """
        # Repeat the process till data distribution is complete

        # if not resimulation:
        ingestion = (ingestion / 60) * self.frequency_minutes
        while int(ingestion) > 0:
            # Select an index
            # index_id = random.randint(0, len(self.cluster.indices) - 1)
            index_id = random.randint(0, self.current_index_count - 1)

            # If rolled over index is chosen, chose a different index
            while index_id == self.cluster.rolled_over_index_id:
                index_id = random.randint(0, len(self.cluster.indices) - 1)

            # Choose a part of the data and push it to index
            data_pushed_to_index_gb = random.uniform(0.1 * ingestion, ingestion)

            # subtract from the total
            ingestion -= data_pushed_to_index_gb

            # Get primary shard count and evaluate the data size to be pushed
            primary_shards_count = self.cluster.primary_shards_per_index
            data_per_shard_gb = data_pushed_to_index_gb / primary_shards_count

            # Update the size of shards
            for shard in range(len(self.cluster.indices[index_id].shards)):
                self.cluster.indices[index_id].shards[
                    shard
                ].shard_size += data_per_shard_gb

    @timeit
    def run(self, duration_minutes, start_time="00_00_00", resimulate=False, time=None):
        resultant_cluster_objects = []
        if time == None:
            now = datetime.now()
        else:
            now = time

        if start_time == "00_00_00":
            date_obj = now - timedelta(
                hours=now.hour,
                minutes=now.minute,
                seconds=now.second,
                microseconds=now.microsecond,
            )
        else:
            date_obj = now

        if resimulate:
            start_time_simulate = self.total_simulation_minutes - duration_minutes
            start_time_minute = int(
                (start_time_simulate - (start_time_simulate % self.frequency_minutes))
                / self.frequency_minutes
            )
            data_y = self.simulated_data_rates[start_time_minute - 1 :]
            data_y1 = {}
            data_y1["simple"] = self.simulated_search_rates["simple"][
                start_time_minute - 1 :
            ]
            data_y1["medium"] = self.simulated_search_rates["medium"][
                start_time_minute - 1 :
            ]
            data_y1["complex"] = self.simulated_search_rates["complex"][
                start_time_minute - 1 :
            ]
            data_y2 = self.index_added_list[start_time_minute - 1 :]
            self.cluster.date_time = date_obj

        else:
            data_x, data_y = self.aggregate_data(duration_minutes, start_time)
            data_x1, data_y1 = self.aggregate_data_searches(
                duration_minutes, start_time
            )
            data_y2 = self.aggregate_index_addition()
            self.simulated_data_rates = data_y.copy()
            self.simulated_search_rates = data_y1.copy()
            self.index_added_list = data_y2.copy()

            self.total_simulation_minutes = duration_minutes

        for index, instantaneous_data_ingestion_rate in enumerate(data_y):
            self.current_index_count = data_y2[index]
            self.cluster.instantaneous_index_count = data_y2[index]
            self.cluster._ingestion_rate = instantaneous_data_ingestion_rate
            self.cluster._simple_query_rate = data_y1["simple"][index]
            self.cluster._medium_query_rate = data_y1["medium"][index]
            self.cluster._complex_query_rate = data_y1["complex"][index]
            self.cluster.cpu_usage_percent = self.cpu_used_for_ingestion(
                instantaneous_data_ingestion_rate, data_y1, index
            )
            self.cluster.memory_usage_percent = self.memory_used_for_ingestion(
                instantaneous_data_ingestion_rate, data_y1, index
            )
            self.cluster.heap_usage_percent = self.heap_used_for_ingestion(
                instantaneous_data_ingestion_rate,
                data_y1,
                index,
                self.cluster.memory_usage_percent,
            )
            self.cluster.status = self.cluster_state_for_ingestion(
                instantaneous_data_ingestion_rate
            )
            self.add_index_to_cluster(index, self.cluster.date_time)
            self.distribute_load(instantaneous_data_ingestion_rate)
            # Todo: simulate effect on remaining cluster parameters
            self.cluster.cluster_disk_size_used = self.disk_utilization_for_ingestion()
            self.cluster.cluster_disk_size_used = self.disk_util_for_index_roll_over(
                self.cluster.date_time
            )
            # self.cluster.cluster_disk_size_used+= (constants.INITIAL_DISK_SPACE_FACTOR * self.cluster.total_disk_size_gb)
            self.cluster.disk_usage_percent = min(
                (self.cluster.cluster_disk_size_used / self.cluster.total_disk_size_gb)
                * 100,
                100,
            )
            date_time = date_obj + timedelta(minutes=self.elapsed_time_minutes)
            self.cluster.date_time = date_time
            self.cluster.active_primary_shards = (
                self.cluster.primary_shards_per_index * self.cluster.index_count
            )
            self.cluster.active_shards = (
                self.cluster.primary_shards_per_index
                * (self.cluster.replica_shards_per_index + 1)
            ) * self.cluster.index_count
            resultant_cluster_objects.append(copy.deepcopy(self.cluster))
            self.elapsed_time_minutes += self.frequency_minutes

        print("======== Size of nodes ===========")
        for node in range(len(self.cluster.nodes)):
            print(
                "Size of Node " + str(node) + " : ",
                self.cluster.nodes[node].calculate_total_node_size(),
            )

        print("========= Number of Shards in nodes ========= ")
        for node in range(len(self.cluster.nodes)):
            print(
                "node " + str(node) + ": ", len(self.cluster.nodes[node].shards_on_node)
            )

        # print("======= Size of Indexes ========")
        # for index in range(len(self.cluster.indices)):
        #     print(
        #         "Size of index " + str(index) + " : ",
        #         self.cluster.indices[index].get_index_primary_size(),
        #     )

        print("========= Index Roll over size ========")
        print(self.cluster.rolled_index_size)

        print("========= Size of Cluster ========")
        print(self.cluster.cluster_disk_size_used)

        # print("======= Index Addition List =======")
        # print(self.index_added_list)

        # print("======= Index obj ========")
        # for index in range(len(self.cluster.indices)):
        #     print(self.cluster.indices[index].__dict__)

        print("")
        return resultant_cluster_objects
