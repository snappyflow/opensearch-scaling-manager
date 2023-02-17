from shard import Shard
import datetime

class Index:
    """
    Essentially a collection of shards along with other
    parameters and methods that transform and store the
    state of shards associated to this index.
    """

    def __init__(
            self,
            primary_shards_count: int,
            replica_shards_count: int,
            index_id: int,
            time: datetime.datetime
    ):
        """
        Initialize the index object
        :param primary_shard_count: number of primary shards per index
        :param replica_shard_count: number of replica shards per index
        :param index_id: unique node identifier for a index
        """
        self.primary_shards_count = primary_shards_count
        self.replica_shards_count = replica_shards_count
        self.index_id = index_id
        self.rolled_over = False
        self.index_size = 0
        self.created_at = time
        self.time_elapsed_last_roll_over = time
        self.shards = self.initialize_shards(primary_shards_count, replica_shards_count)
        

    @property
    def max_shard_size(self):
        """
        Returns the maximum size a shard can have, based on
        shards per and index and roll over size
        :return: shard size in bytes
        """
        # Todo: establish relation between shard size,
        #  shards per and index and roll over size
        return

    @property
    def host(self):
        return self.host

    @host.setter
    def host(self, value):
        self._host = value
    
    def initialize_shards(self, primary_count, replica_count):
        """
        The function creates shard objects and associates it witht the index
        Number of shard objects for a given index is the sum of primary and 
        replica count
        :return shards: A list of shard objects for a given index
        """

        shards = []

        for primary in range(primary_count):

            primary_shard = Shard("Primary", self.index_id)
            shards.append(primary_shard)

            for replica in range(replica_count):
                replica_shard = Shard("Replica",  self.index_id)
                shards.append(replica_shard)
            
        return shards

    def get_index_primary_size(self):
        """
        Calculates the size occupied by 
        the primary shards on the index 
        """
        size = 0

        for shard in self.shards:
            if shard.type == 'Primary':
                size+= shard.shard_size
        
        return size
    
    def get_index_size(self):
        """
        Evaluates and returns the size occupied
        by the shards of the index
        """
        if self.rolled_over:
            return self.shards[0].shard_size
    
        size = 0
        for shard in self.shards:
                size+= shard.shard_size
        
        return size
        