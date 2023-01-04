from datetime import datetime
from node import Node
from index import Index
import random
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
        status: str = "green",
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
        self.cluster_disk_size_used = 0
        self.heap_usage_percent = heap_usage_percent
        self.total_nodes_count = total_nodes_count
        self.active_data_nodes = active_data_nodes
        self.master_eligible_nodes_count = master_eligible_nodes_count
        self.index_count = index_count
        self.index_roll_over_size_gb = index_roll_over_size_gb
        self.index_clean_up_age_in_minutes = index_clean_up_age_days
        self.total_shard_count = total_shard_count
        self.primary_shards_per_index = primary_shards_per_index
        self.replica_shards_per_index = replica_shards_per_index
        self.rolled_over_index_id = -1
        self.total_shards_per_index = primary_shards_per_index * (
            1 + replica_shards_per_index
        )
        self.initializing_shards = initializing_shards_count
        self.relocating_shards = relocating_shards_count
        self.unassigned_shards = unassigned_shards_count
        self.active_shards = active_shards_count
        self.date_time = datetime.now()
        self._ingestion_rate = 0
        self._simple_query_rate = 0
        self._medium_query_rate = 0
        self._complex_query_rate = 0
        self.active_primary_shards = active_primary_shards
        self.nodes = self.initialize_nodes(
                                            total_nodes_count, 
                                            index_count, 
                                            primary_shards_per_index,
                                            replica_shards_per_index
                                            )
        self.indices = self.initialize_indices(index_count,primary_shards_per_index,replica_shards_per_index)
        self.allocate_shards_to_node()


    # TODO: Define methods for controlling cluster behaviour,
    #  node addition, removal etc
    def add_nodes(self, nodes=1):
        self.total_nodes_count += nodes
        self.status = constants.CLUSTER_STATE_YELLOW
        # Todo - simulate effect on shards

    def remove_nodes(self, nodes=1):
        self.total_nodes_count -= nodes
        self.status = constants.CLUSTER_STATE_YELLOW
        # Todo - simulate effect on shards
    
    def initialize_nodes(
        self, 
        total_nodes_count,
        index_count,
        primary_shards_count,
        replica_shards_count 
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
                node = Node(index_count, primary_shards_count, replica_shards_count,i)
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

    def create_index(self, primary_count,replica_count):
        """
        Creates an index object with specified number
        of primary and replica shards
        """
        index = Index(primary_count, replica_count, len(self.indices))
        self.indices.append(index)


    def allocate_shards_to_node(self):
        """
        Allocates shards arbitrarily to nodes,
        This creates shards to node mapping
        """
        node_id = self.get_node_id()
        
        for node in self.nodes:
            node.shards_on_node.clear()

        for index in self.indices:
            for shard in index.shards: 
                id = random.choice(node_id)
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

        for index in self.indices:
           total_size+= index.get_index_size()

        return total_size
