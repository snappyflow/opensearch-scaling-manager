# file paths
SCHEMA_FILE_NAME = "schema.py"
CONFIG_FILE_PATH = "config.yaml"
PROVISION_LOCK_FILE_NAME = "provision.lock"
SIMULATION_GRAPH_FILE_NAME = "simulation_diagram.png"

# config keys
# DATA_INGESTION = "data_ingestion"
STATES = "states"
SEARCH_DESCRIPTION = "search_description"
SIMULATION_FREQUENCY_MINUTES = "simulation_frequency_minutes"
DATA_INGESTION_RANDOMNESS_PERCENTAGE = "randomness_percentage"
# Todo - Shrinidhi/Manoj add remaining keys from config yaml

CLUSTER_STATE_GREEN = "green"
CLUSTER_STATE_YELLOW = "yellow"
CLUSTER_STATE_RED = "red"

CLUSTER_STATES = [CLUSTER_STATE_GREEN, CLUSTER_STATE_YELLOW, CLUSTER_STATE_RED]
INITIAL_DISK_SPACE_FACTOR = 0.01

HIGH_INGESTION_RATE_GB_PER_HOUR = 60

TIME_FORMAT = "%d-%m-%Y %H:%M:%S"
# mapping inputs for API endpoints
CLUSTER_STATE = 'status'
CPU_USAGE_PERCENT = 'cpu_usage_percent'
MEMORY_USAGE_PERCENT = 'memory_usage_percent'
DISK_USAGE_PERCENT = 'disk_usage_percent'
HEAP_USAGE_PERCENT = 'heap_usage_percent'
TOTAL_NODES_COUNT = 'total_nodes_count'
ROLLED_INDEX_SIZE = 'rolled_index_size'
STAT_REQUEST = {
    'CpuUtil': CPU_USAGE_PERCENT,
    'RamUtil': MEMORY_USAGE_PERCENT,
    'HeapUtil': HEAP_USAGE_PERCENT,
    'status': CLUSTER_STATE,
    'nodes': TOTAL_NODES_COUNT,
    'DiskUtil': DISK_USAGE_PERCENT,
    'rolled_index_size' : ROLLED_INDEX_SIZE
}  # Todo : Shrinidhi/Manoj Add remaining stats that will be queried from the recommendation engine

CLUSTER_STATE = "status"
TOTAL_NUM_NODES = "total_nodes_count"
NUM_ACTIVE_SHARD_COUNT = "active_shards_count"
NUM_ACTIVE_PRIMARY_SHARDS = "active_primary_shards"
NUM_INITIALIZING_SHARDS = "initializing_shards_count"
NUM_UNASSIGNED_SHARDS = "unassigned_shards_count"
NUM_RELOCATING_SHARDS = "relocating_shards_count"
NUM_MASTER_NODES = "master_eligible_nodes_count"
NUM_ACTIVE_DATA_NODES = "active_data_nodes"
STAT_REQUEST_CURRENT = {
    "NumNodes": TOTAL_NUM_NODES,
    "ClusterStatus": CLUSTER_STATE,
    "NumActiveShards": NUM_ACTIVE_SHARD_COUNT,
    "NumActivePrimaryShards": NUM_ACTIVE_PRIMARY_SHARDS,
    "NumInitializingShards": NUM_INITIALIZING_SHARDS,
    "NumUnassignedShards": NUM_UNASSIGNED_SHARDS,
    "NumRelocatingShards": NUM_RELOCATING_SHARDS,
    "NumMasterNodes": NUM_MASTER_NODES,
    "NumActiveDataNodes": NUM_ACTIVE_DATA_NODES,
}


APP_PORT = 5000
QUERY_ARG_LENGTH_ONE = 1
QUERY_ARG_LENGTH_TWO = 2
QUERY_ARG_LENGTH_THREE = 3
QUERY_ARG_LENGTH_FOUR = 4
PRIMARY_SHARDS_IN_ROLLOVER_INDEX = 1
REPLICA_SHARDS_IN_ROLLOVER_INDEX = 1
METRIC_PARAMETER = 'metric'
DURATION_PARAMETER = 'duration'
TIME_NOW_PARAMETER = 'time_now'
THRESHOLD_PARAMETER = 'threshold'
