# Open Search Cluster Simulator
Open Search Simulator is a python sub project that attempts to mimic to behavior of an AWS on which OpenSearch is deployed.
It exposes a set of APIs that let's user get and set cluster paramaters like cpu usage statistics, number of nodes of cluster, etc. 



## Configurations
The user can specify some key features of an OpenSearch Cluster through the config file provided. The functionalities supported are:
1.	Cluster state based on ingestion rate states
2.	Searches being performed on the cluster



## Simulator Behavior 
As simulator starts, it generates and stores the data points corresponding to the entire day and stores them in a internal database. Based on the user inputs (through APIs), the data points are fetched or re-generated.

Currently, the equation of cluster stats is purely experiment driven.

<img src="https://github.com/maplelabs/opensearch-scaling-manager/blob/master/images/Simulator.png?raw=true" alt="Simulator_Architecture">



## APIs
Simulator provide the following APIs to interact with it:


| Path                                                 | Description                                                                               | Method | Path Parameters                                                              | Request Body         | Response                                     |
|------------------------------------------------------|-------------------------------------------------------------------------------------------|--------|------------------------------------------------------------------------------|----------------------|----------------------------------------------|
| `/stats/avg?metric=<stat_name>&duration=<duration>`                  | Returns the average value of a stat for the last specified duration.                      | GET    | __stat_name__: string <br/> __duration__: integer                            | None                 | `{"avg": float, "min": float, "max": float}` |
| `/stats/violated?metric=<stat_name>&duration=<duration>&threshold=<threshold>` | Returns the number of time, a stat crossed the threshold duration the specified duration. | GET    | __stat_name__: string <br/> __duration__: integer <br/> __threshold__: float | None                 | `{"ViolatedCount": int}`                     |
| `/stats/current/metric=<stat_name>`                         | Returns the most recent value of a stat.                                                  | GET    | __stat_name__: string                                                        | None                 | `{"current": float}`                         |
| `/provision/addnode`                                 | Ask the simulator to perform a node addition.                                             | POST   | None                                                                         | `{"nodes": integer}` | `{'expiry': ISO Date time}`                  |
| `/provision/remnode`                                 | Ask the simulator to perform a node removal.                                              | POST   | None                                                                         | `{"nodes": integer}` | `{'expiry': ISO Date time}`                  |


## Future Scope
There are currently a lot of missing features that are expected from the simulator.
Ability to specify the type of machine is one such. This would require profiling of 
instance types based on their specifications.


## Getting Started
1. Dowload the project repository.
2. cd into simulator folder i.e cd simulator
3. Create a python virtual environment, i.e python -m venv venv
4. Activate the virtual environemnt, i.e source venv/bin/activate
5. install dependencies using pip, i.e pip install -r requirements.txt
6. Start the simulator using appy.py i.e python src/app.py
7. Request metrics from simultor, eg. curl http://127.0.0.1:5000/stats/avg/60