package utilities

import (
        "context"
        "encoding/json"
        os "scaling_manager/opensearch"
        "scaling_manager/logger"
        "hash/fnv"

)

var log = new(logger.LOG)

func init() {
        log.Init("logger")
        log.Info.Println("Utilities module initiated")
}

// Input: opensearch client and context
// Description: Checks if the current node is the master node.
// Output: A boolean value, true if current node is master, false if it is not.
func CheckIfMaster(ctx context.Context, nodeId string) bool {
        var clusterStateInterface map[string]interface{} //To store the cluster state info and parse for master node ID
        var nodeStatsInterface map[string]interface{}    //To store current node stats and parse for current node ID

        //Create cluster state request and fetch cluster state
        clusterState, err := os.GetClusterState(ctx)
        if err != nil {
                panic(err)
        }

        //Decoding the response and dumping in the cluster state interface
        decodeErr := json.NewDecoder(clusterState.Body).Decode(&clusterStateInterface)
        if decodeErr != nil {
                log.Panic.Println("decode Error: ", decodeErr)
                panic(err)
        }

        //Parsing interface to get the id of the master node
        masterNode := clusterStateInterface["master_node"].(string)

        if nodeId != "" {
                return masterNode == nodeId
        }

        nodes := []string{"_local"}

        //Creating node stats request and fetching the node stats for the current node
        nodeStatReq, err := os.GetNodeStats(nodes, nil, ctx)
        if err != nil {
                log.Panic.Println("Node stat fetch error: ", err)
                panic(err)
        }

        //Decoding the response and dumping the node stats in the interface
        nodeDecodeErr := json.NewDecoder(nodeStatReq.Body).Decode(&nodeStatsInterface)
        if nodeDecodeErr != nil {
                log.Panic.Println("decode Error: ", nodeDecodeErr)
                panic(err)
        }

        //Parsing for the node id of the current node
        var currentNode string
        for node := range nodeStatsInterface["nodes"].(map[string]interface{}) {
                currentNode = node
                break
        }
        return masterNode == currentNode
}

func GetClusterId() string {
        var clusterStatsInterface map[string]interface{}
        resp, err := os.GetClusterStats(context.Background())
        if err != nil {
                log.Error.Println("cluster Stats fetch ERROR:", err)
        }

        decodeErr := json.NewDecoder(resp.Body).Decode(&clusterStatsInterface)
        if decodeErr != nil {
                log.Error.Println("decode Error: ", decodeErr)
        }

        return clusterStatsInterface["cluster_uuid"].(string)
}

func GetNodes() map[string]interface{} {
        var nodeStatsInterface map[string]interface{}
        var nodeMap map[string]interface{}

        nodes := []string{"_all"}
        metrics := []string{}

        nodeStatResp, err := os.GetNodeStats(nodes, metrics, context.Background())
        if err != nil {
                log.Error.Println("Node stat fetch error: ", err)
        }

        decodeErr := json.NewDecoder(nodeStatResp.Body).Decode(&nodeStatsInterface)
        if decodeErr != nil {
                log.Error.Println("decode Error: ", decodeErr)
        }

        for node, nodeInfo := range nodeStatsInterface["nodes"].(map[string]interface{}) {
                var nodeMap map[string]interface{}
                nodeInfoMap := nodeInfo.(map[string]interface{})
                nodeMap[node] = map[string]string{"name": nodeInfoMap["name"].(string), "hostIp": nodeInfoMap["ip"].(string) }
        }

        return nodeMap

}

// Input: string
//
// Description: Returns a hashed value of the string passed as input
//
// Output: uint32 (Hashed value of string)

func Hash(s string) uint32 {
        h := fnv.New32a()
        h.Write([]byte(s))
        return h.Sum32()
}

func ParseNodeId(mapp map[string]interface{}) string {
        for node := range mapp {
                return node
        }
        return ""
}
