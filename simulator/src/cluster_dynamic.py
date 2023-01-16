class ClusterDynamic:
    """
    Represents the Dynamic state of the cluster
    """
    def __init__(self):
        self.cluster_status = "Green"
        self.unassigned_shards = 0
        self.total_nodes = 0
        self.active_shards_count = 0
        self.relocating_shards = 0
        self.active_primary_shards = 0
    