class Shard:
    """
    Shard acts as the smallest unit of data aggregation.
    They can either be primary or replica.
    """

    def __init__(
            self,
            type_of_shard: str,
            index_id: int
    ):
        # Todo: Decide initializing, relocating, unassigned and
        #  active shards are to be properties of individual shard
        #  or managed in cluster
        """
        initializes the shard object
        :param index_id: unique index identifier to which shard will be mapped
        :param type_of_shard: "primary" or "replica" shard
        :param state: Can have the following values only
            INITIALIZING: The shard is recovering from a peer shard or gateway.
            RELOCATING: The shard is relocating.
            STARTED: The shard has started.
            UNASSIGNED: The shard is not assigned to any node.
        """
        self.type = type_of_shard
        self.shard_size = 0
        self.node_id = -1
        self.index_id = index_id
        self.state = "unassigned"


    # @property
    # def total(self):
    #     """
    #
    #     :return:
    #     """
    #     # Todo: Derive total number of shards from initializing,
    #     #  relocating, unassigned and active shards
    #     return self.initializing + self.unassigned + self.relocating

    @property
    def host(self):
        return self.host

    @host.setter
    def host(self, value):
        self._host = value
