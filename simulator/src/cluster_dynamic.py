class ClusterDynamic:
    """
    Represents the Dynamic state of the cluster
    """

    def __init__(self):
        self.ClusterStatus = "Green"
        self.NumActiveDataNodes = 0
        self.NumUnassignedShards = 0
        self.NumInitializingShards = 0
        self.NumNodes = 0
        self.NumActiveShards = 0
        self.NumRelocatingShards = 0
        self.NumMasterNodes = 0
        self.NumActivePrimaryShards = 0
