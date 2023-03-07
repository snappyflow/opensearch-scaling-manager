### Open-search Scaling Manager

------

#### Overview

------

OpenSearch scaling manager is used to elastically scale a cluster to ensure optimum cluster performance and expenses involved. Scaling Manager can automatically scale up or scale down an OpenSearch node based on the effect of load on metric in cluster. Scaling Manager can be used to automate manual scale up, scale down and reduce the manual effort to achieve the same. Scale up, Scale down can happen whenever  it meets the criteria which is mentioned by the user. In addition to this there is event based scaling where as scale up, scale down happens at specific time.

**List of features:**

- Automatic Scaling 

  Parameters supported are,

  1. CPU usage
  2. Mem usage
  3. Heap usage
  4. Disk usage
  5. Shards

- Event based Scaling 

#### Brief explanation, Architecture of Scaling Manager

------

- Lets consider the cluster has 3 nodes. OpenSearch scaling manager is now installed in each of the nodes present in cluster. Node metrics(CPU, Mem, Heap, Disk utilization) is monitored in each nodes present in the cluster and Cluster Metrics(Number of nodes, Shards) of the cluster is also monitored.

- Rules are specified by the user in config.yaml file like what should be maximum usage of CPU, Mem, Heap, Disk, Maximum nodes allowed, Maximum nodes allowed etc. and those rules are verified across the resource utilized in cluster.

- Now scaling manager will check the resource utilization and if the utilization is more than the rules which user specifies in config.yaml it scales up a new node to the cluster in order to accommodate the high  resource utilization.

  ​	For Example if the average cpu usage is more than 80% across the decision_period mentioned in the cluster, you have to scale_up a node in order to bring the cpu usage less, this applies for other metrics(Mem, Heap, Disk) as well. 
  ​	When it comes to scale_down if the average cpu usage is less than 30% across you have to scale_down a node and similar to other metrics as well.

  1. If cpu_util > 80, scale_up a node
     If mem_util > 90, scale_up a node

  2. If cpu_util < 80, scale_down a node
     If mem_util < 90, scale_down a node

     <img src="https://github.com/maplelabs/opensearch-scaling-manager/blob/release_v0.1_dev/images/ScalingManager_Architecture.png?raw=true" alt="Scaling_Manager_Architecture">

     

#### Working Principle of Scaling Manager

------

Scaling manager has following modules

- Fetch Metrics
- Recommendation
- Trigger
- Provision
- State

**Fetch Metrics:** 

- Scaling Manager code is deployed and all the nodes available in the cluster will be running the fetch metrics code.
- Only the current master node of the cluster will collect the cluster level data and each node will collect the node level data.
- Node metrics (Usage of CPU, Mem, Heap, Disk etc.) are collected for each node and those are aggregated for cluster level.
- In addition to the aggregated data, Cluster metrics (Number of nodes, Cluster Status, Shards) are collected and both are indexed into Elasticsearch.
- Old data is purged periodically from the index where the duration can be specified by the user.
- Collected metrics is fetched from recommendation engine periodically.

**Recommendation:** 

- Collected metrics are now checked against the rules which are specified by the user in config.yaml file.

- Rules have the following details, 

  1. Metric - Metrics can be CPU, Mem, Heap, Disk, Shard 

  2. Limit - Limit indicates the threshold value of metric, If this threshold is achieved for a given metric for the decision period then the rule will be activated.

  3. Stat - Stat indicates the statistics on which the evaluation of the rule will happen. For CPU and Mem the values can be:

     ​         Avg: The average CPU or MEM value will be calculated for a given decision period.

     ​         Count: The number of occurences where CPU or MEM value crossed the threshold limit.

  4. Decision Period - Decision Period indicates the time in minutes for which a rule is evaluated.

  5. Occurrences - Occurrences indicate the number of time a rule reached the threshold limit for a give decision period.

- For a particular task there is two operators(OR,AND). When the operator is OR, task is recommended if any of the mentioned rules is satisfied when the operator is AND, task is recommended only if all the mentioned rules is satisfied. 

- If the metrics are satisfied against the rules, i.e If usage is more than the rules specified then recommendation of Scale-up-by-1 comes as a task or If usage is less then recommendation of Scale-down-by-1 comes as a task.

- The recommendation data(Scale-up-by-1 or Scale-down-by-1) is maintained in a command queue.

- Recommended data is next passed to the trigger module.

**Trigger:** 

- Gets the Task from the recommendation Queue.

- Checks the state of the cluster.

- If the cluster is in normal state

  - Give command of Scale-up-by-1 or Scale-down-by-1 Provision module.

  - Before giving command of Scale-down-by-1, check is made if clusterstatus != green recommendation can not be provisioned untill clusterstate becomes green.

  - When the command is given then the states are updated normal to provisioning

    i.e ("states = normal" to "states = provisioning").

  - Log "provision triggered - Up/Down Number of Nodes".

- Else

  - Clear the command queue and goes to the next recommendation in queue, commands are ignored since the cluster health is in normal state.

**Provision:**

- Receives the command from the trigger module.
- For recommendation to be provisioned state should be "state = normal" when it is normal provisioning starts and it updates "state = provisioning" and it indicates whether scaleup / scaledown process is happening.
- Take action based on provisioning command(Scale-up-by-1 or Scale-down-by-1) i.e spin up a  new node in a cluster/delete a node in a cluster. 
- Scale up will invoke commands to create a VM based on cloud type. Then it will configure the OpenSearch on newly created nodes and add the newly spinned up node to list of nodes available. Check is made if node is added to cluster, if it is added install and start scaling manager on new node. 
- Scale down will terminate number of node, before scale down it identifies which node should be terminated(It should not be a master  node, 1st node other than master node will be terminated).
- The execution engine has to be cloud independent.
- If provisioning is completed successfully, update "state = provision_completed".
- Again the state is set back to "state = normal" for next provision to happen.
- If provisioning failed, update state = provisioning_failed.
- All the step by step process of scale_up/ scale_down is been logged into OpenSearch where you can check what is the status of provision, At what time did the provision take place, Is the provision successful or failed, reason for failure etc.

**State:** 

- Check the current state and status of cluster and update state.

- If state == provision completed.

  - Check if all system metrics are in a normal state
    - Cluster state is green.
    - No relocating shards.
  - Update "state = Normal".

  

#### Scaling Manager Flow Diagram 

------

![Scaling_Manager_Flow_diagram](https://github.com/maplelabs/opensearch-scaling-manager/blob/release_v0.1_dev/images/Detailed_flow_ScalingManager.png?raw=true)



#### Scaling Manager Architecture

------

1. Scaling Manager is deployed in all the nodes in cluster. Lets say cluster has 3 nodes. Now resource utilization went high and there is a need of new node in cluster.
2. When a new node is added to the cluster ansible scripts will run in new node and it will install Scaling Manger, OpenSearch, All the necessary details which is needed and the new node details will be added to the available nodes list in order to monitor it

![Scaling_Manager_flow](https://github.com/maplelabs/opensearch-scaling-manager/blob/release_v0.1_dev/images/Basic_flow_ScalingManager.png?raw=true)



#### Crypto

------

- Crypto is used to convert your credentials like username, password of os_credentials, cloud credentials in a encrypted way to maintain confidentiality of your data. 
- Crypto generates a secret key which will be used for encryption and decryption of credentials which is a random string of length 16. 
- Credentials(1. os_credentials - os_admin_username, os_admin_password, 2.cloud_credentials - secret_key, access_key) will be taken from config.yaml to encrypt and encode the credentials and store those encoded data in config file.
- The credentials is now encrypted and stored in config.
- The secret key which is stored is decrypted and decoded back.
- The secret stored will undergo process for scramble and unscramble by converting into matrix and interchanging it.
- Once the credentials are encrypted,updated in config file, the file is updated over all the nodes present in the cluster. By this way if master node goes down other node which can become as master will have the encrypted data.  



#### Scale Up and Scale Down

------

##### **Scale Up** 

- New node that is added to the cluster will be configured with all the requirements such as OpenSearch, Security groups, sudo aspects, ssh aspects etc... in order to communicate with the other nodes in the cluster.
- When the process of scale_up is completed by provision and when scale_up is recommended again in specified decision time, it will discard the scale_up since there was already successful provision done. So it discards provision until next polling interval.
- When provision(scale_up) is recommended and the cluster has reached maximum number of nodes(specified in config.yaml), scaling manager will not scale up until max_nodes_allowed is increased manually by user in config.yaml and it will log the message to notify the user to increase the size.

##### **Scale Down** 

- Identifies node which is other than master node to remove from the cluster and stores the node IP. Configuring (Reallocating the shards to other nodes) to remove the node from cluster. Remove the node and terminate the instance.
- When task == scale_down && Cluster_Status != green, recommendation(task) can not be provisioned as open search cluster is unhealthy for a scale_down.
- When provision(scale_down) is recommended and the cluster has reached minimum number of nodes(specified in config.yaml), scaling manager will not scale down until min_nodes_allowed is decreased manually by user in config.yaml and it will log the message to notify the user to decrease the size.

<img src="https://github.com/maplelabs/opensearch-scaling-manager/blob/release_v0.1_dev/images/Scale_up-Scale_down.png?raw=true" alt="Scale_up,Scale_down" style="zoom:150%;" />



#### Scaling Manager Configuration

------

The user can specify some key features of an OpenSearch Cluster for simulator through the config file provided. The functionalities supported are:

**user_config:**

​	**monitor_with_logs:** Field that contains bool value which specifies whether to monitor with logs or not.

​	**monitor_with_simulator:** Field that contains bool value which specifies whether to monitor with simulator or not.

​	**purge_old_docs_after_hours:** Duration which indicates to delete the documents once it exceed the specified hours.

​	**polling_interval_in_secs:**  polling_interval_in_secs indicates the time in seconds for which polling will be repeated.

​	**is_accelerated:** Field that contains bool value which accelerates the time.

**cluster_details:**

​	**ip_address:** IP address of the cluster.

​	**launch_template_id:** ID by which launch template can be identified and deployed.

​	**launch_template_version:** Version of the launch template used.

​	**cluster_name:** Name of the cluster. 

​	**os_credentials:** 

​		**os_admin_username:** Username for the OpenSearch for connecting. This can be set to empty if the security is disable in OpenSearch.

​		**os_admin_password:** Password for the OpenSearch for connecting. This can be set to empty if the security is disable in OpenSearch.

​	**os_user:** Used in ansible for copy files with user.

​	**os_group:** Used in ansible for copy files with group.

​	**os_version:** OpenSearch version which needs to be used.

​	**os_home:** Default OpenSearch user info.

​	**domain_name:** Configure hostnames for OpenSearch nodes which is required to configure SSL.

​	**cloud_type:** Name of the cloud infrastructure.

​	 **cloud_credentials:**

​		**secret_key:** Secret key for cluster.

​		**access_key:** Access key for cluster.

​		**pem_file_path:** Path where the pem file is located. 

​		**region:** Region at which AWS is used.

​	 **base_node_type:** t2x.large.

​	 **number_cpus_per_node:** Total number of CPU present per node.

​	 **ram_per_node_in_gb:** Size of RAM used per node (GB).

​	**disk_per_node_in_gb:** Size of DISK used per node (GB).

​	**max_nodes_allowed:** Maximum number of nodes allowed for the cluster.

​	**min_nodes_allowed:** Minimum number of nodes allowed for the cluster.

**task_details:** Field that contains details on what task should be performed i.e scale_up_by_1 or scale_down_by_1.

- **task_name:** Task name indicates the name of the task to recommend by the recommendation engine.
  **operator:** Operator indicates the logical operation needs to be performed while executing the rules.
  **rules:** Rules indicates list of rules to evaluate the criteria for the recommendation engine.

  - **metric:** Metric indicates the name of the metric. These can be CpuUtil, MemUtil, ShardUtil, DiskUtil

    **limit:** Limit indicates the threshold value for a metric.

    **stat:** Stat indicates the statistics on which the evaluation of the rule will happen. These can be AVG, COUNT.

    **decision_period:** Decision Period indicates the time in minutes for which a rule is evaluated.

  

#### Sample config.yaml

------

[opensearch-scaling-manager/config.yaml at release_v0.1_dev · maplelabs/opensearch-scaling-manager (github.com)](https://github.com/maplelabs/opensearch-scaling-manager/blob/release_v0.1_dev/config.yaml)



#### Scaling Manager Pre-Requisites

------

- Cluster with OpenSearch installed.
- OpenSearch version - 1.2.4 and above. 
- Go version - 1.19
- Ansible Version - 2.9.
- Cluster credentials (Username, Password) to access the OpenSearch.
- Cloud credential  (Username, Password). 
- In AWS we can create a instance by templates which is provided by Domain that is used.
- Launch Template - AWS launch template to spin a new node which has the necessary tags.
- Template ID format (lt-xxxxxxxxxxxxxxxxx.)
- Security certificate to have regex in it to accept the new node.
- PEM file.
- SSH aspect - If cloud type is AWS then Security group is configured in such a way that newly spin up node should be reached via ssh.
- Sudo permission - All the nodes, jump host should have sudo permission by which task could be performed with sudo access between nodes and run ansible playbook with sudo access on jump host. Sudo password can be empty which is preferable.



#### Jump Host login details

------

- Download any remote computing toolbox like MobaXterm. 
- Click Session -> SSH -> Remote host.
- Enter Remote host details, mention the username.
- Click Advanced SSH Settings, choose the PEM file that is present in your local and click OK.
- Login using the Cluster and Jump host details.



#### Build and Installation of Scaling Manager

------

**Inventory file** -  Defines the hosts and groups of hosts upon which commands, modules, and tasks in a playbook operate.

**Populate inventory.yaml**

```master_node_ip
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "populate_inventory_yaml" -e master_node_ip=0.0.0.0 -e os_user=USERNAME -e os_pass=PASSWORD
```

master_node_ip = IP address of master node,
os_user = Appropriate username,
os_pass = Appropriate password

**Build, Pack**

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "build_and_pack" -kK
```

-kK is used for password authentication

In case to use key based authentication, Use the following command

```
udo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "update_pem" --key-file user-dev-aws-ssh.pem -e pem_path="user-dev-aws-ssh.pem"
```

**Installation**

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "install" -kK
```

- Update should be performed when provision is not happening, then stop, install, start the service. These steps can be performed in the same command as well.  

- Stop, Install, Start in same command 

  ```
  sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "stop,install,start" -kK
  ```

**Update Config**

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "update_config" -kK
```

**Start** 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "start" -kK
```

**Stop**

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "stop" -kK
```

**Status**

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "status" -kK
```

**Uninstall**

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "uninstall" -kK
```

- Stop command works quick when there is no provisioning happening/provisioning is completed.
- When provisioning is in process and the stop command is executed it waits till provisioning is completed. To know the status user can do Ctrl+C and check the status of the cluster. 



#### Simulator 

------

Find more about Simulator here [opensearch-scaling-manager/readme_simulator.md at release_v0.1_dev · maplelabs/opensearch-scaling-manager (github.com)](https://github.com/maplelabs/opensearch-scaling-manager/blob/release_v0.1_dev/docs/readme_simulator.md)
