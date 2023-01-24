from index import Index
from shard import Shard

class Node:
    """
    Representation of a OpenSearch node that essentially
    is a collection of node resources like shards, cpu, memory
    and the corresponding methods altering these parameters
    """
    def __init__(
            self,
            index_count: int,
            primary_shard_count: int,
            replica_shard_count: int,
            node_id: int
    ):
        """
        Initialize the node object
        :param index_count: number of indexes on a cluster
        :param primary_shard_count: number of primary shards per index
        :param replica_shard_count: number of replica shards per index
        :param node_id: unique node identifier for a node
        """

        self.index_count = index_count
        self.primary_shard_count = primary_shard_count
        self.replica_shard_count = replica_shard_count
        self.node_id = node_id
        self.node_available = True
        self.shards_on_node = []
        

    @property
    def is_master(self):
        """
        States whether the node is mater or not
        :return: bool
        """
        return 'master' in self.roles

    @property
    def is_data(self):
        """
        States whether the node is data or not
        :return: bool
        """
        return 'data' in self.roles

    def calculate_total_node_size(self):
        """
        Calculates the total size of a node.
        This calculation is based on the sum of 
        size of shards present on the node
        """
        size = 0

        if self.node_available:
            for shard in self.shards_on_node:
                size+=shard.shard_size

        return size 

    def initialize_index(self, index_count, primary_count, replica_count):
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
