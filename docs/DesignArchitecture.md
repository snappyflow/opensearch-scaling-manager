# Design Architecture

- [Design Architecture](#design-architecture)
  - [List of features:](#list-of-features)
  - [Working Principle of Scaling Manager](#working-principle-of-scaling-manager)
  - [Scaling Manager Flow Diagram](#scaling-manager-flow-diagram)
  - [Scaling Manager Architecture](#scaling-manager-architecture)
  - [Crypto](#crypto)
  - [Scale Up and Scale Down](#scale-up-and-scale-down)
  - [Scaling Manager Configuration](#scaling-manager-configuration)

## **List of features:**

  - Automatic Scaling 

    Parameters supported are,

    1. CPU usage
    2. Mem usage
    3. Heap usage
    4. Disk usage
    5. Shards

  - Event based Scaling 

- Lets consider the cluster has 3 nodes. OpenSearch scaling manager is now installed in each of the nodes present in cluster. Node metrics(CPU, Mem, Heap, Disk utilization) is monitored in each nodes present in the cluster and Cluster Metrics(Number of nodes, Shards) of the cluster is also monitored.

- Rules are specified by the user in config.yaml file like what should be maximum usage of CPU, Mem, Heap, Disk, Maximum nodes allowed, Maximum nodes allowed etc. and those rules are verified across the resource utilized in cluster.

- Now scaling manager will check the resource utilization and if the utilization is more than the rules which user specifies in config.yaml it scales up a new node to the cluster in order to accommodate the high  resource utilization.

  For Example if the average cpu usage is more than 80% across the decision_period mentioned in the cluster, you have to scale_up a node in order to bring the cpu usage less, this applies for other metrics(Mem, Heap, Disk) as well. 
  When it comes to scale_down if the average cpu usage is less than 30% across you have to scale_down a node and similar to other metrics as well.

  1. If cpu_util > 80, scale_up a node
     If mem_util > 90, scale_up a node
  2. If cpu_util < 80, scale_down a node
     If mem_util < 90, scale_down a node

<img src="https://github.com/maplelabs/opensearch-scaling-manager/blob/master/images/ScalingManagerArchitecture.png" alt="ScalingManagerArchitecture">

## Working Principle of Scaling Manager

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

     Avg: The average CPU or MEM value will be calculated for a given decision period.

     Count: The number of occurences where CPU or MEM value crossed the threshold limit.

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

  

## Scaling Manager Flow Diagram 

<img src="https://github.com/maplelabs/opensearch-scaling-manager/blob/master/images/DetailedFlowScalingManager.png" alt="DetailedFlowScalingManager">



## Scaling Manager Architecture

1. Scaling Manager is deployed in all the nodes in cluster. Lets say cluster has 3 nodes. Now resource utilization went high and there is a need of new node in cluster.
2. When a new node is added to the cluster ansible scripts will run in new node and it will install Scaling Manger, OpenSearch, All the necessary details which is needed and the new node details will be added to the available nodes list in order to monitor it

<img src="https://github.com/maplelabs/opensearch-scaling-manager/blob/master/images/BasicFlowScalingManager.png" alt="BasicFlowScalingManager">



## Crypto

- Crypto is used to convert your credentials like username, password of os_credentials, cloud credentials in a encrypted way to maintain confidentiality of your data. 
- Crypto generates a secret key which will be used for encryption and decryption of credentials which is a random string of length 16. 
- Credentials(1. os_credentials - os_admin_username, os_admin_password, 2.cloud_credentials - secret_key, access_key) will be taken from config.yaml to encrypt and encode the credentials and store those encoded data in config file.
- The credentials is now encrypted and stored in config.
- The secret key which is stored is decrypted and decoded back.
- The secret stored will undergo process for scramble and unscramble by converting into matrix and interchanging it.
- Once the credentials are encrypted,updated in config file, the file is updated over all the nodes present in the cluster. By this way if master node goes down other node which can become as master will have the encrypted data.  



## Scale Up and Scale Down

**Scale Up** 

- New node that is added to the cluster will be configured with all the requirements such as OpenSearch, Security groups, sudo aspects, ssh aspects etc... in order to communicate with the other nodes in the cluster.
- When the process of scale_up is completed by provision and when scale_up is recommended again in specified decision time, it will discard the scale_up since there was already successful provision done. So it discards provision until next polling interval.
- When provision(scale_up) is recommended and the cluster has reached maximum number of nodes(specified in config.yaml), scaling manager will not scale up until max_nodes_allowed is increased manually by user in config.yaml and it will log the message to notify the user to increase the size.

**Scale Down** 

- Identifies node which is other than master node to remove from the cluster and stores the node IP. Configuring (Reallocating the shards to other nodes) to remove the node from cluster. Remove the node and terminate the instance.
- When task == scale_down && Cluster_Status != green, recommendation(task) can not be provisioned as open search cluster is unhealthy for a scale_down.
- When provision(scale_down) is recommended and the cluster has reached minimum number of nodes(specified in config.yaml), scaling manager will not scale down until min_nodes_allowed is decreased manually by user in config.yaml and it will log the message to notify the user to decrease the size.

<img src="https://github.com/maplelabs/opensearch-scaling-manager/blob/master/images/ScaleUpScaleDown.png" alt="ScaleUpScaleDown">

## Scaling Manager Configuration

Please check [config file](https://github.com/maplelabs/opensearch-scaling-manager/blob/master/docs/Config.md) to know more about scaling manager configuration
