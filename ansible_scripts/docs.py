import requests
import json
import sys
import time
from requests.auth import HTTPBasicAuth
url = "http://"+sys.argv[4]+":9200/_nodes/"+sys.argv[3]+"/stats/indices"
health_url = "http://"+sys.argv[4]+":9200/_cluster/health"
docs_count = 1
while docs_count != 0:
  unassigned_shards = requests.get(health_url, auth = HTTPBasicAuth(sys.argv[1],sys.argv[2]))
  json_resp = unassigned_shards.json()
  if json_resp['unassigned_shards'] > 0:
    print("Still moving shards across the cluster.")
    print("\nWARNING!! But there are unassigned shards in the cluster: ", json_resp['unassigned_shards'])
    print("\n Please take a look and remove the exclusion of node from _cluster/settings and restart scaling manager if you do not want to remove this node\n")
  response = requests.get(url,auth = HTTPBasicAuth(sys.argv[1],sys.argv[2]))
  json_obj = response.json()
  node_id= list(json_obj['nodes'].keys())
  docs_count= json_obj['nodes'][node_id[0]]['indices']['docs']['count']
  print(docs_count)
  time.sleep(5)
