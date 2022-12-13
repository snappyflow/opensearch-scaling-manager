
# file paths
SCHEMA_FILE_NAME = 'schema.py'
CONFIG_FILE_PATH = 'config.yaml'
PROVISION_LOCK_FILE_NAME = 'provision.lock'
SIMULATION_GRAPH_FILE_NAME = 'simulation_diagram.png'

# config keys
DATA_INGESTION = "data_ingestion"
DATA_INGESTION_STATES = 'states'
SEARCHES = "searches"
DATA_GENERATION_INTERVAL_MINUTES = "data_generation_interval_minutes"
DATA_INGESTION_RANDOMNESS_PERCENTAGE = 'randomness_percentage'
# Todo - Shrinidhi/Manoj add remaining keys from config yaml

CLUSTER_STATE_GREEN = 'green'
CLUSTER_STATE_YELLOW = 'yellow'
CLUSTER_STATE_RED = 'red'

CLUSTER_STATES = [CLUSTER_STATE_GREEN, CLUSTER_STATE_YELLOW, CLUSTER_STATE_RED]

HIGH_INGESTION_RATE_GB_PER_HOUR = 60

# mapping inputs for API endpoints
CLUSTER_STATE = 'status'
CPU_USAGE_PERCENT = 'cpu_usage_percent'
MEMORY_USAGE_PERCENT = 'memory_usage_percent'
STAT_REQUEST = {
    'cpu': CPU_USAGE_PERCENT,
    'mem': MEMORY_USAGE_PERCENT,
    'status': CLUSTER_STATE

}  # Todo : Shrinidhi/Manoj Add remaining stats

APP_PORT = 5000
