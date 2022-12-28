# Open Search Cluster Simulator
Open Search Simulator is an attempt to mimic to behavior of an AWS on which OpenSearch is deployed. 



## Configurations
The user can specify some key features of an OpenSearch Cluster through the config file provided. The functionalities supported are:
1.	Cluster state based on ingestion rate states
2.	Searches being performed on the cluster



## Simulator Behavior 
As simulator starts, it generates and stores the data points corresponding to the entire day and stores them in a internal database. Based on the user inputs (through APIs), the data points are fetched or re-generated.



## APIs
Simulator provide the following APIs to interact with it


| Path                                                 | Description                                                                               | Method | Path Parameters                                                              | Request Body         | Response                                     |
|------------------------------------------------------|-------------------------------------------------------------------------------------------|--------|------------------------------------------------------------------------------|----------------------|----------------------------------------------|
| `/stats/avg/{stat_name}/{duration}`                  | Returns the average value of a stat for the last specified duration.                      | GET    | __stat_name__: string <br/> __duration__: integer                            | None                 | `{"avg": float, "min": float, "max": float}` |
| `/stats/violated/{stat_name}/{duration}/{threshold}` | Returns the number of time, a stat crossed the threshold duration the specified duration. | GET    | __stat_name__: string <br/> __duration__: integer <br/> __threshold__: float | None                 | `{"ViolatedCount": int}`                     |
| `/stats/current/{stat_name}`                         | Returns the most recent value of a stat.                                                  | GET    | __stat_name__: string                                                        | None                 | `{"current": float}`                         |
| `/provision/addnode`                                 | Ask the simulator to perform a node addition.                                             | POST   | None                                                                         | `{"nodes": integer}` | `{'expiry': ISO Date time}`                  |
| `/provision/remnode`                                 | Ask the simulator to perform a node removal.                                              | POST   | None                                                                         | `{"nodes": integer}` | `{'expiry': ISO Date time}`                  |

