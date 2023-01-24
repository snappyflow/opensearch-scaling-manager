package provision

import (
        "context"
        "encoding/json"
        "fmt"
        osutils "scaling_manager/opensearchUtils"
        utils "scaling_manager/utilities"

)

// A global variable which stores the document ID of the State document that will to stored and fetched frm Opensearch
var docId string

func InitializeDocId() {
        docId = fmt.Sprint(utils.Hash(utils.GetClusterId()))
}

// Input:
// Description:
//
//      GetCurrentState will update the state variable pointer such that it is insync with the updated values.
//      Reads the document from Opensearch and updates the Struct
//
// Return:
//

func (s *State) GetCurrentState() {
        // Get the document.

        searchResponse, err := osutils.SearchDoc(docId, context.Background())
        if err != nil {
                log.Panic.Println("failed to search document: ", err)
                panic(err)
        }
        var stateInterface map[string]interface{}
        log.Debug.Println("Get resp: ", searchResponse)
        if searchResponse.Status() == "404 Not Found" {
                //Setting the initial state
                s.CurrentState = "normal"
                s.UpdateState()
                return
        }
        jsonErr := json.NewDecoder(searchResponse.Body).Decode(&stateInterface)
        if jsonErr != nil {
                log.Panic.Println("Unable to decode the response into interface: ", jsonErr)
                panic(jsonErr)
        }
        // convert map to json
        jsonString, errr := json.Marshal(stateInterface["_source"].(map[string]interface{}))
        if errr != nil {
                log.Panic.Println("Unable to unmarshal interface: ", errr)
                panic(errr)
        }

        // convert json to struct
        json.Unmarshal(jsonString, s)
}

// Input:
//
// Description:
//
//      Updates the opensearch document with the values in state Struct pointer.
//
// Return:

func (s *State) UpdateState() {
        // Update the document.

        state, err := json.Marshal(s)
        if err != nil {
                log.Panic.Println("json.Marshal ERROR: ", err)
                panic(err)
        }
        content := string(state)

        updateResponse, err := osutils.UpdateDoc(docId, content, context.Background())
        if err != nil {
                log.Panic.Println("failed to update document: ", err)
                panic(err)
        }
        log.Debug.Println("Update resp: ", updateResponse)
}
