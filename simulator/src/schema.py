{
    "cluster_name": {"required": True, "type": "string"},
    "cluster_hostname": {"required": True, "type": "string"},
    "cluster_ip_address": {"required": True, "type": "string"},
    "node_machine_type_identifier": {"required": True, "type": "string"},
    "total_nodes_count": {"required": True, "type": "number"},
    "data_nodes_count": {"required": True, "type": "number"},
    "master_eligible_nodes_count": {"required": True, "type": "number"},
    "index_count": {"required": True, "type": "number"},
    "shards_per_index": {"required": True, "type": "number"},
    "index_roll_over_size_gb": {"required": True, "type": "number"},
    "index_clean_up_age_days": {"required": True, "type": "number"},
    "simulation_frequency_minutes": {"required": True, "type": "number"},
    "data_ingestion": {
        "required": False,
        "type": "dict",
        "schema": {
            "states": {
                "required": True,
                "type": "list",
                "schema": {
                    "type": "dict",
                    "schema": {
                        "position": {"required": True, "type": "number"},
                        "time_hh_mm_ss": {"required": True, "type": "string"},
                        "ingestion_rate_gb_per_hr": {
                            "required": True,
                            "type": "number",
                        },
                    },
                },
            },
            "randomness_percentage": {"required": True, "type": "number"},
        },
    },
    "searches": {
        "required": False,
        "type": "list",
        "schema": {
            "type": "dict",
            "schema": {
                "search_type": {"required": True, "type": "string"},
                "probability": {"required": True, "type": "number"},
                "cpu_load_percent": {"required": True, "type": "number"},
                "memory_load_percent": {"required": True, "type": "number"},
            },
        },
    },
}
