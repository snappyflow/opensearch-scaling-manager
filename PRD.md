#               OpenSearch Scaling Manager 



## Overview

OpenSearch scale manager is used to elastically scale a cluster to ensure optimum cluster performance and expenses involved.



## Approach

Cluster Scaling can happen based on the following rules

- **Scale-Up**: 

  - If the load is gradually increasing scale up by 1 node.
  - If the load increase is significant and persistent then scale up by more than 1 node.

- **Rule based cluster sizing**

  - Ensure X Nodes are present in the Cluster at a given time.
    - 5 Nodes, 9 A.M, Monday
    - 3 Nodes,  9 A.M, Saturday
  - Consider Cluster sizing if the number of Nodes present does not match the number specified in the rule.

- **Scale down**

  - If the load is gradually falling and ingest is low, Scale down by 1 Node.

    

**Detecting gradual increase in load**

- Monotonic Increase in CPU Utilization and CPU Utilization > Threshold.

- Monotonic increase in Memory Utilization and Memory Utilization > Threshold.

- Monotonic increase in JVM Utilization and JVM Utilization > Threshold

  

**Detecting sudden increase in load** 

- Calculate the ingest rate based on the data increase in current indexes. This rate is expressed on GB/day.

- If ingest rate is steady and increasing.

- Calculate the number of Nodes corresponding to Ingest Rate.

- If number of Nodes required > current Cluster size, increase the cluster by **N** nodes, 

  where N = number of Nodes required - Current number of Nodes in the Cluster.

  

**Rule based Cluster sizing**

- Rule can be executed only if Current Ingest Rate < Ingest Rate that new Cluster can handle.

- Else the rule is ignored with a Warning log.

  

**Detecting the gradual decrease in load**

- Check if Cluster Utilization is consistently low and monotonically reducing.

- Check Ingest Rate and if warrants that the Cluster size has to be reduced.

- If above checks are true, then trigger Scale-down by 1 Node.

  

There are 2 key modules for lifecycle manager:

1. **Analyzer**: Performs analysis and recommends action to be take.

2. **Executor**: Executes the recommended action.

   

**Avoiding Conflict of Rules**

- Each scheme generates recommendations based on its analysis. There could be situations where multiple recommendations are generated in an analysis cycle. Analyzer has to reduce these multiple recommendations into a single recommendation.
- One approach to reduce multiple recommendation to a single recommendation, is to choose the largest cluster size recommended.
- Any new recommendation can be acted upon only when the previous recommendation has been executed and the cluster is stable.

**Sensor Failure**

- The solution completely depends on input signals
- If any input signal is not available, Scaling manager generates error logs and does not make any recommendation.

**Executor Failure**

- If there is any error in performing scale-up/scale-down, Scaling manager generates error logs and does not generate further recommendations until error state is cleared.

  

**Simulator**

- Simulator is a critical module to test,validate and refine solution.
- Simulator will model Elasticsearch through the variables that matter to scaling manager
  - CPU Utilization
  - RAM Utilization
  - JVM Utilization
  - Ingest Rate
  - Number of Nodes
- Ingest rate is represented using a cyclical function with some randomness.
- Recommendation affects the number of nodes in the cluster and in some cases the executor may fail.



## Architecture

The solution includes the following components

- Fetch Metrics

- Recommendation

- Trigger

- Provision

- State

  

**Fetch Metrics**

- Code is deployed on each node available in the cluster.

- Only the current master node will have the privilege to execute the code.

- Monitored metrics are indexed into Elasticsearch.

- Old data is purged periodically from the index.

  

**Recommendation**

- Performs checks based on the configuration file.
- Based on the checks, provides a Scale-up-by-1 or Scale-down-by-1 recommendation.
- The data is maintained in a command queue.



**Trigger**

- Checks the state of the cluster
- If the cluster is in normal state
  - Give command to Provision module
  - Updates state = provisioning
  - Log "provision triggered - Up/Down Number of  Nodes"
- Else
  - Clear the command queue, commands are ignored since the cluster health criteria is not satisfied.



**Provision**

- Receives the command from the trigger module and updates state = Provision
- Take action based on provisioning command. The execution engine has to be cloud independent
- If provisioning completed successfully, update state = provision completed.
- If provisioning failed, update state = provisioning failed.



**State**

- Check the current state and status of cluster and update state
- If state == provision completed
  - Check if all system metrics are in a normal state
    - Cluster state is green
    - No relocating shards
  - Update State = Normal



**Installation of the solution**

- Package to be installed on all nodes of Elasticsearch cluster.
- Should automatically create an Index if not present.
- On uninstalling, index should be deleted.



**Configuration inputs**

- Cluster details

  - IP/DNS
  - Elasticsearch Credentials
  - Cloud Type - AWS, GCP, Azure
  - Cloud Credentials/ IAM role
  - Base node type 
  - CPUs per node
  - RAM per node in GB
  - Disk per node in GB
  - Maximum nodes allowed

- Task details

  - Scale-up-by-1 [ Triggered based on OR of the conditions]
    - Data nodes required
    - CPU Utilization Rule
      - Enabled
      - Limit - 80%
      - Stat - Average/Count
      - Decision Period - 60 minutes.
    - Memory Utilization Rule
      - Enabled
      - Limit - 90%
      - Stat - Average/Count
      - Decision Period - 60 minutes
    - Shard Rule
      - Enabled
      - Limit - 900
      - Decision Period - 60 minutes
  - Scale-down-by-1 [Triggered based on AND of the conditions]
    - Data nodes required.
    - CPU Utilization Rule
      - Enabled
      - Limit - 60%
      - Stat - Average/Count
      - Decision Period - 2 days
    - Memory Utilization Rule
      - Enabled
      - Limit - 60%
      - Stat - Average/Count
      - Decision Period - 2 days
    - Shard Rule
      - Enabled
      - Limit - 900
      - Decision Period - 2 days
  - Implementation of the tasks
    - Implements following functions
      - Average(time period): Returns struct Cluster Stats which contains CPU, Memory, Shards - Maximum, Minimum , Average.
      - Count(time period,threshold_list): Returns struct Cluster Stats which contains number of times violated for CPU, Memory, Shards.
      - A module looks at all the tasks and invokes build stats, then it evaluates each task which compares with the cluster average.

- Fetch Metrics 

  - Metrics required

    - Node level
      - CPU utilization
      - RAM Utilization
      - Heap Utilization
      - Disk Utilization
      - Shards count
    - Cluster level
      - Active Master nodes count
      - Active Data nodes count
      - Overall cluster status
      - Initializing shards count
      - Unassigned shards count
      - Relocating shards count

  - Periodicity: 5 minutes

  - Maintains data for 72 hours in the Elasticsearch index.

    

- Elasticsearch simulator

  - This module has following functions

    - Model behavior based on configuration file

    - Provides Metrics based on the API call

    - Alter behavior based on provisioning commands

      

  - Configuration file

    - Data node count

    - Shards

    - Index count

    - Shards per Index 

    - Aggregate data rate function

      - Saw tooth(max,min,max_time,min_time)
      - Add random(min %, max %)

    - Index roll-over size

    - Index clean-up age

    - Disk function

      - Derived from shards
      - Disk clean-up limits

    - Search events

      - Simple search

        - Probability
        - CPU Load
        - Memory Load

      - Medium search

        - Probability 
        - CPU Load
        - Memory Load

      - Complex search

        - Probability

        - CPU Load

        - Memory Load

          

  - Behavior

    - Initial focus on Index and Shards

      - Define Initial indexes

      - Number of Shards = N * Indexes *shards per Index *2

      - Calculate Index size every cycle based on

        - Ingest function
        - Random value

      - If Index size = Threshold, then Roll over

      - If disk size limit reached, clean up old indexes

      - Equally distribute index and shard across all nodes.

        

    - Disk size

      - Totals disk size = sum of all indexes

      - Equal distribution of disk size per node

        

    - CPU Utilization

      - CPU Load in cores = function(Ingest Rate) + function(CPU load from search)

      - CPU Utilization of cluster = used cores/ total CPU cores

      - CPU Utilization of a node is same as the cluster since it is equally distributed

        

    - Memory Utilization

      - Memory Load = function(Ingest Rate) + function(Memory load from search)
      - Memory Utilization of a cluster = used memory / total memory
      - Memory Utilization of a node is same as the cluster since it is equally distributed

