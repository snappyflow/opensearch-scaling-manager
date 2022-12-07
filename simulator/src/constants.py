
SCHEMA_FILE_NAME = 'schema.py' # Todo : Check if needed

CLUSTER_STATE_GREEN = 'green'
CLUSTER_STATE_YELLOW = 'yellow'
CLUSTER_STATE_RED = 'red'

CLUSTER_STATES = [CLUSTER_STATE_GREEN, CLUSTER_STATE_YELLOW, CLUSTER_STATE_RED]

HIGH_INGESTION_RATE_GB_PER_HOUR = 60

STAT_REQUEST = {
    'cpu': 'cpu_usage_percent',
    'mem': 'memory_usage_percent',
    'status': 'status'
} # Todo : Add remaining stats  

APP_PORT = 5000

SCHEMA_PATH = '\schema.py'
CONFIG_PATH = '\config.yaml'

DATA_INGESTION = "data_ingestion"
SEARCHES = "searches"
DATA_GENERATION_INTERVAL_MINUTES = "data_generation_interval_minutes"