from datetime import datetime
from node import Node
from index import Index
from cluster_dynamic import ClusterDynamic
import random
import time
import constants


class Cluster:
    """
    Acts as an interface for simulation of all associated nodes
    Performs and simulates the output of all operations performed
    by the master node.
    """

    def __init__(
        self,
        cluster_name: str,
        cluster_hostname: str,
        cluster_ip_address: str,
        node_machine_type_identifier: str,
        total_nodes_count: int,
        active_data_nodes: int,
        master_eligible_nodes_count: int,
        index_count: int,
        index_roll_over_size_gb: int,
        index_clean_up_age_days: int,
        primary_shards_per_index: int,
        replica_shards_per_index: int,
        min_nodes_in_cluster: int,
        heap_memory_factor: float,
        cluster_dynamic: ClusterDynamic = None,
        status: str = "Green",
        cpu_usage_percent: float = 0,
        memory_usage_percent: float = 0,
        disk_usage_percent: float = 0,
        total_disk_size_gb: int = 0,
        heap_usage_percent: float = 0,
        total_shard_count: int = 0,
        initializing_shards_count: int = 0,
        relocating_shards_count: int = 0,
        unassigned_shards_count: int = 0,
        active_shards_count: int = 0,
        active_primary_shards: int = 0,
    ):
        """
        Initialize the cluster object
        :param primary_shards_per_index:
        :param replica_shards_per_index:
        :param cluster_name: name of the cluster
        :param cluster_hostname: name of the cluster host
        :param cluster_ip_address: ip address of the cluster
        :param node_machine_type_identifier: type of machine on which elastic search is deployed
        :param total_nodes_count: total number of nodes of the cluster
        :param active_data_nodes: total number of data nodes of the cluster
        :param master_eligible_nodes_count: total number of master eligible nodes of the cluster
        :param index_count: total number of indexes in the cluster
        :param index_roll_over_size_gb: size in gb after which the index will be rolled over
        :param index_clean_up_age_days: time in minutes after which the index will be cleaned up
        :param status: status of the cluster from "green", "yellow" or "red"
        :param cpu_usage_percent: average cluster cpu usage in percent
        :param memory_usage_percent: average cluster memory usage in percent
        :param disk_usage_percent: average cluster disk usage in percent
        :param heap_usage_percent: average cluster heap memory usage in percent
        :param total_shard_count: total numer of shards present on the cluster
        :param initializing_shards_count: total number of shards in initializing state
        :param relocating_shards_count: total number of shards in relocating state
        :param unassigned_shards_count: total number of shards in unassigned state
        :param active_shards_count: total number of shards in active state
        :param active_primary_shards: total number of primary shards in active state
        """
        self.node_machine_type_identifier = node_machine_type_identifier
        self.name = cluster_name
        self.host_name = cluster_hostname
        self.ip_address = cluster_ip_address
        self.status = status
        self.cpu_usage_percent = cpu_usage_percent
        self.memory_usage_percent = memory_usage_percent
        self.disk_usage_percent = disk_usage_percent
        self.total_disk_size_gb = total_disk_size_gb
        self.heap_memory_factor = heap_memory_factor
        self.cluster_disk_size_used = 0
        self.heap_usage_percent = heap_usage_percent
        self.total_nodes_count = total_nodes_count
        self.active_data_nodes = active_data_nodes
        self.master_eligible_nodes_count = master_eligible_nodes_count
        self.index_count = index_count
        self.index_roll_over_size_gb = index_roll_over_size_gb
        self.index_clean_up_age_in_minutes = index_clean_up_age_days
        self.primary_shards_per_index = primary_shards_per_index
        self.replica_shards_per_index = replica_shards_per_index
        self.total_shard_count = (primary_shards_per_index
                                 * (replica_shards_per_index + 1)
                                 ) * index_count
        self.rolled_over_index_id = -1
        self.cluster_dynamic = cluster_dynamic
        self.min_nodes_in_cluster = min_nodes_in_cluster
        self.total_shards_per_index = primary_shards_per_index * (
            1 + replica_shards_per_index
        )
        self.initializing_shards = initializing_shards_count
        self.relocating_shards = relocating_shards_count
        self.unassigned_shards = unassigned_shards_count
        self.unassigned_shards_list = []
        self.active_shards = active_shards_count
        self.date_time = datetime.now()
        self._ingestion_rate = 0
        self._simple_query_rate = 0
        self._medium_query_rate = 0
        self._complex_query_rate = 0
        self.rolled_index_size = 0
        self.active_primary_shards = active_primary_shards
        self.nodes = self.initialize_nodes(
            total_nodes_count,
            index_count,
            primary_shards_per_index,
            replica_shards_per_index,
        )
        self.indices = self.initialize_indices(
            index_count, primary_shards_per_index, replica_shards_per_index
        )
        self.allocate_shards_to_node()

    # TODO: Define methods for controlling cluster behaviour,
    #  node addition, removal etc
    def add_nodes(self, nodes, accelerate=False):
        """
        Adds node to cluster and performs shards rebalancing.
        Cluster state will be Yellow till rebalancing is complete.
        """
        # Update the total node count in cluster dynamic
        self.cluster_dynamic.ClusterStatus = constants.CLUSTER_STATE_YELLOW
        self.cluster_dynamic.NumMasterNodes = self.master_eligible_nodes_count + nodes
        self.cluster_dynamic.NumActiveDataNodes = self.active_data_nodes + nodes
        self.cluster_dynamic.NumNodes = self.total_nodes_count + nodes
        self.cluster_dynamic.NumActivePrimaryShards = (
            self.primary_shards_per_index * self.index_count
        )
        self.cluster_dynamic.NumActiveShards = (
            self.primary_shards_per_index * (self.replica_shards_per_index + 1)
        ) * self.index_count
        self.cluster_dynamic.NumRelocatingShards = ((
            self.primary_shards_per_index * (self.replica_shards_per_index + 1)
        ) * self.index_count) // self.total_nodes_count

        # Add the node
        for node in range(nodes):
            new_node = Node(0, 0, 0, len(self.nodes))
            existing_node_id = self.get_available_node_id()
            self.nodes.append(new_node)
            rebalancing_size = self.cluster_disk_size_used / (self.total_nodes_count + nodes)
            self.total_disk_size_gb+= (self.total_disk_size_gb/self.total_nodes_count)
            rebalance_time = self.time_function_for_rebalancing(rebalancing_size,accelerate)
            self.rebalance_shards(rebalance_time,existing_node_id, len(existing_node_id))
            self.total_nodes_count += 1
        self.cluster_dynamic.NumRelocatingShards = 0
        self.status = constants.CLUSTER_STATE_GREEN
        self.active_data_nodes+=nodes
        self.master_eligible_nodes_count+=nodes
        self.cluster_dynamic.ClusterStatus = constants.CLUSTER_STATE_GREEN
        return

        # Perform rebalancing
        self.status = constants.CLUSTER_STATE_YELLOW
        # Todo - simulate effect on shards

    def remove_nodes(self, nodes, accelerate):
        """
        Removes node from cluster, rebalances unassigned shards due to
        removed node. If sufficient nodes are not present to allocate
        unassigned shards, the cluster will be in yellow state. The
        cluster will be in Yellow state when there is shard movement
        to allocate unassigned shards to nodes.
        :param nodes: count of nodes to be added to cluster
        """
        if self.min_nodes_in_cluster > self.total_nodes_count:
            print("Cannot remove more nodes, minimum nodes required")
            return

        # Update the total node count in cluster dynamic
        self.cluster_dynamic.NumNodes = self.total_nodes_count - nodes
        self.cluster_dynamic.NumMasterNodes = self.master_eligible_nodes_count - nodes
        self.cluster_dynamic.NumActiveDataNodes = self.active_data_nodes - nodes
        # Choose a node from cluster and remove it
        for node in range(nodes):
            node_id = random.randint(0, len(self.nodes) - 1)

            while not self.nodes[node_id].node_available:
                node_id = random.randint(0, len(self.nodes) - 1)

            unassigned_shard_size = self.nodes[node_id].calculate_total_node_size()

            self.cluster_dynamic.NumActivePrimaryShards = (
                self.primary_shards_per_index * self.index_count
            )
            self.cluster_dynamic.NumUnassignedShards = len(
                self.nodes[node_id].shards_on_node
            )
            self.cluster_dynamic.NumRelocatingShards = 0
            self.cluster_dynamic.ClusterStatus = constants.CLUSTER_STATE_YELLOW
            self.cluster_dynamic.NumActiveShards = (
                self.primary_shards_per_index * (self.replica_shards_per_index + 1)
            ) * self.index_count

            # shards present on that node will be un-assigned
            for shard in self.nodes[node_id].shards_on_node:
                shard.state = "unassigned"
                self.unassigned_shards += 1
                self.unassigned_shards_list.append(shard)

            del self.nodes[node_id]
            self.total_disk_size_gb-= (self.total_disk_size_gb/self.total_nodes_count)
            self.total_nodes_count-= 1
            self.update_node_id()

        # If sufficient nodes are present
        if self.total_nodes_count >= self.replica_shards_per_index + 1:
            self.rebalance_unassigned_shards(unassigned_shard_size, accelerate)
            self.unassigned_shards_list.clear()
            self.unassigned_shards = 0
            self.cluster_dynamic.NumUnassignedShards = 0
            self.cluster_dynamic.NumRelocatingShards = 0
            self.status = constants.CLUSTER_STATE_GREEN
            self.active_data_nodes-=nodes
            self.master_eligible_nodes_count-=nodes
            self.cluster_dynamic.ClusterStatus = constants.CLUSTER_STATE_GREEN
            return

        # If sufficient nodes not present, set cluster state yellow
        self.status = constants.CLUSTER_STATE_YELLOW
        self.cluster_dynamic.ClusterStatus = "Yellow"
        # Todo - simulate effect on shards

    def initialize_nodes(
        self, total_nodes_count, index_count, primary_shards_count, replica_shards_count
    ):
        """
        Function takes the count of nodes in the cluster and creates a
        list of node objects. Each node object will have arbitrary count
        of indexes and each index will have the primary and replica shards
        as per the parameter.
        :return nodes: A list of node objects
        """
        nodes = []

        for i in range(total_nodes_count):
            # To-do: Add mechanism to distribute the index count randomnly across nodes
            node = Node(index_count, primary_shards_count, replica_shards_count, i)
            nodes.append(node)

        return nodes

    def get_node_id(self):
        """
        Function fetches the node id and returns a list of node id's
        in a cluster object
        """
        node_id = []

        for node in self.nodes:
            node_id.append(node.node_id)

        return node_id

    def update_node_id(self):
        for node_id in range(len(self.nodes)):
            self.nodes[node_id].node_id = node_id

    def get_available_node_id(self):
        node_id = []

        for node in self.nodes:
            if node.node_available:
                node_id.append(node.node_id)

        return node_id

    def initialize_indices(self, index_count, primary_count, replica_count):
        """
        The function will create index objects of the specified count
        Each index will have primary and replica shards of specified
        count
        :return index: A list of index objects
        """
        indices = []

        for i in range(index_count):
            index = Index(primary_count, replica_count, i)
            indices.append(index)

        return indices

    def create_index(self, primary_count, replica_count):
        """
        Creates an index object with specified number
        of primary and replica shards
        """
        index = Index(primary_count, replica_count, len(self.indices))
        self.indices.append(index)

    def clear_index_size(self):
        for index in self.indices:
            index.index_size = 0
            for shard in index.shards:
                shard.shard_size = 0

    def allocate_shards_to_node(self):
        """
        Allocates shards arbitrarily to nodes,
        This creates shards to node mapping
        """
        node_id = self.get_available_node_id()

        for node in self.nodes:
            node.shards_on_node.clear()

        for index in self.indices:
            for shard in index.shards:
                id = random.choice(node_id)

                while not self.nodes[id].node_available:
                    id = random.randint(0, len(self.nodes))

                shard.node_id = id
                shard.state = "started"
                self.nodes[id].shards_on_node.append(shard)

    def calculate_cluster_disk_size(self):
        """
        Evaluates the disk space occupied in the cluster
        Returns the total size used in GB for the cluster
        object
        """
        # To-Do: Total size must be taken from initial size of the cluster before ingestion
        total_size = 0

        for node in self.nodes:
            total_size += node.calculate_total_node_size()
        return total_size

    def get_unassigned_shard_size(self):
        size = 0

        for shard in self.unassigned_shards_list:
            if shard.type == "Replica":
                size += shard.shard_size

        return size

    def time_function_for_rebalancing(self, unassigned_shard_size,accelerate):
        """
        Simulates the time taken for rebalancing the shard
        The time evaluation is based on the size of unassigned shards.
        Time is accelerated in the following mannner
        60 minutes = 5 seconds
        1 minute = 5/60 seconds
        With initial assumption of 1Gb data size takes
        5 minutes to rebalance, it takes 1/12 seconds per GB of data
        The time to rebalance is evaluated for unassigned shards size.
        """
        if accelerate:
            return random.randint(0, 5)
        rebalancing_time = unassigned_shard_size * (1 / 12)
        return rebalancing_time

    def rebalance_shards(self, rebalance_time, existing_node_id_list, new_node_id):
        """
        The function simulates the rebalancing of shards when a new 
        node is added. It takes the rebalance time and simulates the 
        time elapsed for the rebalance of shards
        """
        total_rebalance_shard_count =((
            self.primary_shards_per_index * (self.replica_shards_per_index + 1)
        ) * self.index_count) // self.total_nodes_count
        
        for node_id in existing_node_id_list:
            for shard in range(total_rebalance_shard_count//self.total_nodes_count):
                rebalancing_shard = self.nodes[node_id].shards_on_node.pop()
                rebalancing_shard.node_id = new_node_id
                self.nodes[new_node_id].shards_on_node.append(rebalancing_shard)
                if rebalance_time != 0:
                    sleep_time = random.uniform(0, rebalance_time)
                    rebalance_time -= sleep_time
                    time.sleep(sleep_time)
                self.cluster_dynamic.NumRelocatingShards-=1


    def rebalance_unassigned_shards(self, unassigned_shard_size, accelerate):
        """
        Rebalances unassigned shards among available nodes.
        The time taken for shard rebalancing is simulated
        using time function
        :param unassigned_shard_size: size of unassigned shards 
        """
        # Add time function to simulate the time taken for rebalancing
        rebalance_time = self.time_function_for_rebalancing(unassigned_shard_size, accelerate)

        # Assign the shards to the available nodes on the cluster
        for shard in self.unassigned_shards_list:
            # Choose node to place the shard
            node_id = random.randint(0, len(self.nodes) - 1)

            # If the chosen node is not available then pick different node
            while not self.nodes[node_id].node_available:
                node_id = random.randint(0, len(self.nodes) - 1)

            # Update the shard state and its node id
            shard.state = "started"
            shard.node_id = node_id

            # Add the shard to the node
            self.nodes[node_id].shards_on_node.append(shard)
            self.cluster_dynamic.NumUnassignedShards -= 1
            self.cluster_dynamic.NumRelocatingShards += 1
            if rebalance_time != 0:
                sleep_time = random.uniform(0, rebalance_time)
                rebalance_time -= sleep_time
                time.sleep(sleep_time)
