
SCHEMA_FILE_NAME = 'schema.py' # Todo : Check if needed

CLUSTER_STATE_GREEN = 'green'
CLUSTER_STATE_YELLOW = 'yellow'
CLUSTER_STATE_RED = 'red'

CLUSTER_STATES = [CLUSTER_STATE_GREEN, CLUSTER_STATE_YELLOW, CLUSTER_STATE_RED]

HIGH_INGESTION_RATE_GB_PER_HOUR = 60

CLUSTER_STATE = 'status'
CPU_USAGE_PERCENT = 'cpu_usage_percent'
MEMORY_USAGE_PERCENT = 'memory_usage_percent'
STAT_REQUEST = {
    'cpu': CPU_USAGE_PERCENT,
    'mem': MEMORY_USAGE_PERCENT,
    'status': CLUSTER_STATE

} # Todo : Add remaining stats  

APP_PORT = 5000

CONFIG_PATH = 'config.yaml'

DATA_INGESTION = "data_ingestion"
SEARCHES = "searches"
DATA_GENERATION_INTERVAL_MINUTES = "data_generation_interval_minutes"

PROVISION_LOCK_FILE_NAME = 'provision.lock'