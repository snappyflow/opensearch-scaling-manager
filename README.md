# Open-search Scaling Manager

Open Search Simulator is an attempt to mimic to behavior of an AWS on which OpenSearch is deployed.



### Simulator Configurations

------

The user can specify some key features of an OpenSearch Cluster for simulator through the config file provided. The functionalities supported are:



#### 1.Cluster Stats

------

**cluster_name:** Name of the cluster that is to be used

**cluster_hostname:** Host name of the cluster

**cluster_ip_address:** IP address of the cluster

**node_machine_type_identifier:** Defines the type of the instance or node deployed in a cluster

**total_nodes_count:** Total number of nodes present in the cluster

**active_data_nodes:** Number of active data nodes in total number of nodes present in cluster

**min_nodes_in_cluster:** Minimum number of nodes that the cluster must have to perform the necessary tasks

**master_eligible_nodes_count:** Nodes that are eligible to become master whenever the present master node goes down

**heap_memory_factor:**

**index_count:** Number of index that cluster must have

**primary_shards_per_index:** Number of primary shards that is present in index

**replica_shards_per_index:** Number of replica shards that is present in index(replica of data that represents each primary shard)

**index_roll_over_size_gb:** Specific size at which index will roll over to new index when it exceeds

**index_roll_over_hours:** Specific time in hour at which index will roll over to new index when it exceeds

**total_disk_size_gb:** Total number of size in GB that the disk should have

**simulation_frequency_minutes:** Time interval that the simulator will run the data simulation



#### 2.Data Ingestion

------

Specify data ingestion with respect to time of the day to represent pattern for entire day(24hrs)

**states:** States is an array where user can provide multiple data points through out a day

**day:** Day is an array which contains multiple hour of for the day and also can contain multiple days

**position:** For a day there can be any number of position where it contains time_hh_mm_ss, ingestion_rate_gb_per_hr, searches

**time_hh_mm_ss:** Time interval of the position

**ingestion_rate_gb_per_hr:** Amount of data that has been ingested for the particular interval of time that is defined in time_hh_mm_ss

**searches:** Contains the types of searches that needs to be made, if the config has certain searches it takes the corresponding values. Three types of searches are simple, medium, complex

**index:**

​	**count:** Number of index to add at the specified time interval



#### 3.Randomness Percentage

------

**randomness_percentage:**  Percentage at which the stats value needs to be differing while simulating.



#### 4.Search Description

------

**search_description:** Specify searches along with their type, probability and load inflected on the cluster. Three level of search_description are simple, medium, complex.

**simple:**

​	**cpu_load_percent:** Percentage at which cpu must be used if search_description is simple

​	**memory_load_percent: **Percentage at which memory must be used if search_description is simple

​	**heap_load_percent: **Percentage at which heap must be used  if search_description is simple

**medium:**

​	**cpu_load_percent:** Percentage at which cpu must be used if search_description is medium

​	**memory_load_percent: **Percentage at which memory must be used if search_description is medium

​	**heap_load_percent: **Percentage at which heap must be used if search_description is medium

**complex:**

​	**cpu_load_percent:** Percentage at which cpu must be used if search_description is complex

​	**memory_load_percent: **Percentage at which memory must be used if search_description is complex

​	**heap_load_percent: **Percentage at which heap must be used if search_description is complex



### Sample cofig.yaml

------

https://maplelabsblr-my.sharepoint.com/:u:/g/personal/manojkumar_chandru_maplelabs_com/Efh45zUUBO1NsYKrZX6ycRIB14WPcLR7NQ-RRMXKRqqAqw?e=E63khe



### Simulator Behavior

------

As simulator starts, it generates and stores the data points corresponding to the entire day and stores them in a internal database. Based on the user inputs (through APIs), the data points are fetched or re-generated.



### Installation and Executing Simulator

------

To install the simulator please download the source code using following command:

```
git clone https://github.com/maplelabs/opensearch-scaling-manager.git -b release_v0.1_dev
```



Execute the following commands to run and install the simulator

```python
cd opensearch-scaling-manager/simulator
# Path to simulator module.

python -m venv venv
# Creating virtual environment.

.\venv\Scripts\activate
# Activatinng virtual environment.

pip install -r .\requirements.txt
# Install every requirements for simulator.

cd src
# Path to execute simulator.

python app.py
# Run entire simulator module.
```



### APIs

------

Simulator provide the following APIs to interact with it

| Path               | Query Parameters                                             | Description                                                  | Method | Request Body       | Response                                   |
| :----------------- | ------------------------------------------------------------ | ------------------------------------------------------------ | ------ | ------------------ | ------------------------------------------ |
| /stats/avg         | {key,value} = {metric:string},{duration:int}                 | Returns the average value of a stat for the last specified duration. | GET    | None               | {"avg": float, "min": float, "max": float} |
| /stats/violated    | {key,value} = {metric:string},{duration:int},{threshold:float} | Returns the number of time, a stat crossed the threshold duration the specified duration. | GET    | None               | {"ViolatedCount": int}                     |
| /stats/current     | {key,value} = {metric:string},{duration:int}                 | Returns the most recent value of a stat.                     | GET    | None               | {"current": float}                         |
| /provision/addnode | None                                                         | Ask the simulator to perform a node addition.                | POST   | {"nodes": integer} | {"nodes": int}                             |
| /provision/remnode | None                                                         | Ask the simulator to perform a node removal.                 | POST   | {"nodes": integer} | {"nodes": int}                             |



## Scaling Manager Configuration

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

https://maplelabsblr-my.sharepoint.com/:u:/g/personal/manojkumar_chandru_maplelabs_com/EaY5h9sz2qJAruB2SzaR15UBABaJB3TVOsTE3V4E_G_Rnw?e=n5dgTk



### Working Principle of Scaling Manager

------

Scaling manager has following modules

- Fetch Metrics
- Recommendation
- Trigger
- Provision
- State

**Fetch Metrics:** Collects the metrics of the cluster where the code is deployed and those metrics are indexed into Elasticsearch 

**Recommendation:** Based on the config file of the Scaling Manager it evaluates the rules against the merics that are fetched in Fetch Metrics and recommends the task (scale_up_by_1 or scale_down_by_1 )

**Trigger:** Checks the cluster state if the state is normal gives command to provision module, update state=provisioning and logs "provision triggered - Up/Down Number of Nodes". If cluster state is not normal it clears the command queue, commands are ignored since the cluster health criteria is not satisfied.

**Provision:** Command of the state is received by provisioning and takes action based on provisioning command.

- If provisioning completed successfully, update state = provision completed.
- If provisioning failed, update state = provisioning failed.

**State:** Checks the current state and status of cluster and update state.If state == provision completed

- Check if all system metrics are in a normal state
  - Cluster state is green
  - No relocating shards
- Update State = Normal



### Scaling Manager Flow Diagram 

------

https://lucid.app/publicSegments/view/b8e022c2-8adf-4737-82d8-f3869d61a86a/image.png



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