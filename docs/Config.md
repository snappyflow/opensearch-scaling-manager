### Scaling Manager Configuration

------

The user can specify some key features of an OpenSearch Cluster for simulator through the config file provided. The functionalities supported are :

1. User Configuration - user_config.
2. Cluster Details - cluster_details (Details of the cluster).
3. Task Details - task_details (Scale up / Scale down details).

**user_config:**

**monitor_with_logs:** Field that contains bool value which specifies whether to monitor with logs or not.

**monitor_with_simulator:** Field that contains bool value which specifies whether to monitor with simulator or not.

**purge_old_docs_after_hours:** Duration which indicates to delete the documents once it exceed the specified hours.

**recommendation_polling_interval_in_secs:**  recommendation_polling_interval_in_secs indicates the time in seconds for which polling will be repeated.

**fetchmetrics_polling_interval_in_secs:** fetchmetrics_polling_interval_in_secs indicates the time in seconds for which the metrics will be fetched from the cluster and repeated in the interval.

**is_accelerated:** Field that contains bool value which accelerates the time.



**cluster_details:**

**cluster_name:** Name of the cluster. 

**cloud_type:** Name of the cloud infrastructure.

**max_nodes_allowed:** Maximum number of nodes allowed for the cluster.

**min_nodes_allowed:** Minimum number of nodes allowed for the cluster.

**launch_template_id:** ID by which launch template can be identified and deployed.

**launch_template_version:** Version of the launch template used.

**os_user:** Used in ansible for copy files with user.

**os_group:** Used in ansible for copy files with group.

**os_version:** OpenSearch version which needs to be used.

**os_home:** Default OpenSearch user info.

**domain_name:** Configure hostnames for OpenSearch nodes which is required to configure SSL.

**os_credentials:** 

​	**os_admin_username:** Username for the OpenSearch for connecting. This can be set to empty if the security is disable in OpenSearch.

​	**os_admin_password:** Password for the OpenSearch for connecting. This can be set to empty if the security is disable in OpenSearch.

 **cloud_credentials:**

​	**pem_file_path:** Path where the pem file is located. 

​	**secret_key:** Secret key for cluster.

​	**access_key:** Access key for cluster.

​	**region:** Region at which AWS is used.

​	**role_arn:** AWS IAM role of user which has permissions to spin a node.

**jvm_factor:** Specify the percent of RAM to be allocated to HEAP.



**task_details:** 

Tasks supports two types of scaling 

1. Metric based scaling 
2. Event based scaling 

(Metric based scaling)

- **task_name:** Task name indicates the name of the task to recommend by the recommendation engine.
  **operator:** Operator indicates the logical operation needs to be performed while executing the rules.
  **rules:** Rules indicates list of rules to evaluate the criteria for the recommendation engine.

  - **metric:** Metric indicates the name of the metric. These can be CpuUtil, MemUtil, ShardUtil, DiskUtil

    **limit:** Limit indicates the threshold value for a metric.

    **stat:** Stat indicates the statistics on which the evaluation of the rule will happen. These can be AVG, COUNT.

    **decision_period:** Decision Period indicates the time in minutes for which a rule is evaluated.

    **occurrences_percent:** Percent at which metrics crossed the limit for the specified decision_period. 

(Event based scaling)

- **task_name:** Task name indicates the name of the task to recommend by the recommendation engine.

  **operator:** EVENT

  **rules:**

  **scheduling_time:** Specifies the cron job time at which the task happens

  

## Sample config.yaml

[config.yaml](https://github.com/maplelabs/opensearch-scaling-manager/blob/master/config.yaml)
