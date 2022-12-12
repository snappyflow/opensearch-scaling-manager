# Open Search Cluster Simulator
Open Search Simulator is an attempt to mimic to behavior of an AWS on which OpenSearch is deployed. 



## Configurations
The user can specify some key features of an OpenSearch Cluster through the config file provided. The functionalities supported are:
1.	Cluster state based on ingestion rate states
2.	Searches being performed on the cluster



## Simulator Behavior 
As simulator starts, it generates and stores the data points corresponding to the entire day and stores them in a internal database. Based on the user inputs (through APIs), the data points are fetched or re-generated.

When provisioning is progress, i.e. when a new node is being added or an old node is being removed, the cluster state becomes yellow for sometime and then becomes green again 


## APIs
Simulator provide the following APIs to interact with it


|           |                                     |
|-----------|-------------------------------------|
| __Path__  | `/stats/avg/{stat_name}/{duration}` |


**Description :** Returns the average value of a stat for the last specified 	duration.

**Type :** GET 

**Path Parameters :**

stat_name: string

duration: integer

**Response :** 
`{"avg": float, "min": float, "max": float}`


|           |                                                 |
|-----------|-------------------------------------------------|
| __Path__  | `/stats/avg/{stat_name}/{duration}/{threshold}` |

**Description :** Returns the number of time, a stat crossed the threshold duration the specified duration.

**Type :** GET 

**Path Parameters :**

stat_name: string

duration: integer

threshold: float

**Response :** 
`{"ViolatedCount": int}`


|          |                              |
|----------|------------------------------|
| __Path__ | `/stats/current/{stat_name}` |
**Description :** Returns the most recent value of a stat

**Type :** GET 

**Path Parameters :**

stat_name: string

**Response :** 
`{"current": float}`


|          |                      |
|----------|----------------------|
| __Path__ | `/provision/addnode` |

**Description :** Ask to simulator to represent a node addition

**Type :** POST

**Request Body :**
`{"nodes": integer}`

**Response :** 
`{  
'expiry': ISO Date time
}`
