# Simulator 

Simulator is a module which mimics the real cluster and get the metrics like CPU, RAM, HEAP, SHARD etc.. by which the scaling manager can test with instead of stats from real cluster.



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

[opensearch-scaling-manager/config.yaml at release_v0.1_dev · maplelabs/opensearch-scaling-manager (github.com)](https://github.com/maplelabs/opensearch-scaling-manager/blob/release_v0.1_dev/simulator/src/config.yaml)



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