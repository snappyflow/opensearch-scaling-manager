# Open-search Scaling Manager

Open Search Simulator is an attempt to mimic to behavior of an AWS on which OpenSearch is deployed.



### Working Principle of Scaling Manager

------

Scaling manager has following modules

- Fetch Metrics
- Recommendation
- Trigger
- Provision
- State

![Scaling_Manager_Architecture](https://lucid.app/publicSegments/view/12de2241-e528-4fb2-a891-194ebd2d9c95/image.png)



**Fetch Metrics:** 

- Scaling Manager  code is deployed on each node available in the cluster where OpenSearch is installed.
- Only the current master node of the cluster will have the privilege to execute the code.
- Fetch Metrics collects the monitored metrics of the cluster (Usage of cpu, ram, heap, shard etc.)  and those are indexed into Elasticsearch
- Old data is purged periodically from the index.
- Collected metrics is next passed to the recommendation module.



**Recommendation:** 

- Collected metrics are now checked against the rules which are specified by the user in config.yaml file.
- If the metrics are satisfied against the rules,  provides a Scale-up-by-1 or Scale-down-by-1 recommendation.
- The data(Scale-up-by-1 or Scale-down-by-1) is maintained in a command queue.
- Data is next passed to the trigger module.



**Trigger:** 

- Checks the state of the cluster

- If the cluster is in normal state

  - Give command of Scale-up-by-1 or Scale-down-by-1 to Provision module
  - Updates state = provisioning
  - Log "provision triggered - Up/Down Number of Nodes"

- Else

  - Clear the command queue, commands are ignored since the cluster health criteria is not satisfied.

  

**Provision:**

- Receives the command from the trigger module and updates state = Provision
- Take action based on provisioning command(Scale-up-by-1 or Scale-down-by-1) i.e spin n number node in a cluster/delete n number of node in a cluster. The execution engine has to be cloud independent
- If provisioning is completed successfully, update state = provision_completed.
- If provisioning failed, update state = provisioning_failed.



**State:** 

- Check the current state and status of cluster and update state
- If state == provision completed
  - Check if all system metrics are in a normal state
    - Cluster state is green
    - No relocating shards
  - Update State = Normal



### Scaling Manager Flow Diagram 

------

![Scaling_Manager_Flow_diagram](https://lucid.app/publicSegments/view/b8e022c2-8adf-4737-82d8-f3869d61a86a/image.png)



### Scaling Manager Configuration

------

The user can specify some key features of an OpenSearch Cluster for simulator through the config file provided. The functionalities supported are:

**user_config:**

​	**monitor_with_logs:** Field that contains bool value which specifies whether to monitor with logs or not

​	**monitor_with_simulator:** Field that contains bool value which specifies whether to monitor with simulator or not

​	**purge_old_docs_after_hours:** Duration which indicates to delete the documents once it exceed the specified hours

​	**polling_interval_in_secs:**  polling_interval_in_secs indicates the time in seconds for which polling will be repeated

​	**is_accelerated:** Field that contains bool value which accelerates the time

**cluster_details:**

​	**ip_address:** IP address of the cluster 

​	**launch_template_id:** 

​	**launch_template_version:** 

​	**cluster_name:** Name of the cluster 

​	**os_credentials:** 

​		**os_admin_username:** Username for the OpenSearch for connecting. This can be set to empty if the security is disable in OpenSearch

​		**os_admin_password:** Password for the OpenSearch for connecting. This can be set to empty if the security is disable in OpenSearch

​	**os_user:** SSH login username

​	**os_version:** OpenSearch version

​	**os_home: **Default OpenSearch user info

​	**domain_name:** Configure hostnames for OpenSearch nodes which is required to configure SSL

​	**cloud_type:** Name of the cloud infrastructure

​	 **cloud_credentials:**

​		**secret_key:** Secret key for cluster

​		**access_key:** Access key for cluster

​	 **base_node_type:** t2x.large

​	 **number_cpus_per_node:** Total number of CPU present per node

​	 **ram_per_node_in_gb:** Size of RAM used per node (GB)

​	**disk_per_node_in_gb:** Size of DISK used per node (GB)

​	**max_nodes_allowed:** Maximum number of nodes allowed for the cluster

​	**min_nodes_allowed:** Minimum number of nodes allowed for the cluster

**task_details:** Field that contains details on what task should be performed i.e scale_up_by_1 or scale_down_by_1

- **task_name:** Task name indicates the name of the task to recommend by the recommendation engine.
  **operator:** Operator indicates the logical operation needs to be performed while executing the rules
  **rules:** Rules indicates list of rules to evaluate the criteria for the recommendation engine.

  - **metric:** Metric indicates the name of the metric. These can be CpuUtil, MemUtil, ShardUtil, DiskUtil
    **limit: **Limit indicates the threshold value for a metric.
    **stat:** Stat indicates the statistics on which the evaluation of the rule will happen. These can be AVG, COUNT
    **decision_period:** Decision Period indicates the time in minutes for which a rule is evaluated.

  

### Sample config.yaml

------

[opensearch-scaling-manager/config.yaml at release_v0.1_dev · maplelabs/opensearch-scaling-manager (github.com)](https://github.com/maplelabs/opensearch-scaling-manager/blob/release_v0.1_dev/config.yaml)



### Scaling Manager Pre-Requisites

------

- Cluster with OpenSearch installed 
- OpenSearch version 
- Cluster credentials (Username, Password) to access the OpenSearch
- Cloud credential  (Username, Password) 
- Launch Template - AWS launch template to spin a new node which has the necessary tags
- Security certificate to have regex in it to accept the new node 
- PEM file 



### Ansible Scripts For Scaling Manager

------

Build, Pack

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yml --tags "build_and_pack" -kK
```

Installation

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yml --tags "install" -kK
```

Update Config

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yml --tags "update_config" -kK
```

Start 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yml --tags "start" -kK
```

Stop

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yml --tags "stop" -kK
```

Uninstall

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yml --tags "uninstall" -kK
```



### Build, Packaging and installation

------

To install the scaling manager please download the source code using following command:

```
git clone https://github.com/maplelabs/opensearch-scaling-manager.git -b release_v0.1_dev
```



Run the following commands to build and install the scaling manager

```
cd opensearch-scaling-manager/
# Build the scaling_manager module.
sudo make build
# Package the scaling_manager module and create a tarball.
sudo make pack
# Install the scaling_manager module and create systemd service.
sudo make install
```



To start scaling manager run the following command:

```
sudo systemctl start scaling_manager
```



To stop the scaling manager run the following command:

```
sudo systemctl stop scaling_manager
```





### Simulator 

------

Simulator is a module which mimics the real cluster and get the metrics like CPU, RAM, HEAP, SHARD etc.. by which the scaling manager can test with instead of stats from real cluster.

Find more about Simulator here [opensearch-scaling-manager/readme_simulator.md at release_v0.1_dev · Manojkumar-Chandru-ML/opensearch-scaling-manager (github.com)](https://github.com/Manojkumar-Chandru-ML/opensearch-scaling-manager/blob/release_v0.1_dev/docs/readme_simulator.md)
