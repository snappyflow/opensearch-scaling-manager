package provision

import (
	"context"

	"encoding/json"
	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"scaling_manager/config"
)

// Input:
//
//	username (string): Username string to be used to ssh into the host inventory
//	hosts (string): The file name of hosts file to pass to ansible playbook
//	clusterCfg (config.ClusterDetails): Opensearch cluster details for configuring
//
// Description:
//
//	Calls the ansible script responsible for adding a new node into the Opensearch cluster and configuring it.
//
// Return:
//
//	(error): Returns error if any
func CallScaleUp(username string, hosts string, clusterCfg config.ClusterDetails) error {

	jsonData, err := json.Marshal(&clusterCfg)
	if err != nil {
		log.Error.Println("Error while Marshaling. %v", err)
		return err
	}

	var jsonMap map[string]interface{}

	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		log.Error.Println("json parsing error")
		return err
	}

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		User: username,
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: hosts,
		ExtraVars: jsonMap,
	}

	ansiblePlaybookPrivilegeEscalationOptions := &options.AnsiblePrivilegeEscalationOptions{
		Become: true,
	}

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:                  []string{"ansible_scripts/scaleUpPlaybook.yml"},
		ConnectionOptions:          ansiblePlaybookConnectionOptions,
		PrivilegeEscalationOptions: ansiblePlaybookPrivilegeEscalationOptions,
		Options:                    ansiblePlaybookOptions,
		Exec: execute.NewDefaultExecute(
			execute.WithEnvVar("ANSIBLE_FORCE_COLOR", "true"),
			execute.WithTransformers(
				results.Prepend("Go-ansible with become"),
			),
		),
	}

	err = playbook.Run(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

// Input:
//
//	username (string): Username string to be used to ssh into the host inventory
//	hosts (string): The file name of hosts file to pass to ansible playbook
//	clusterCfg (config.ClusterDetails): Opensearch cluster details for configuring
//
// Description:
//
//	Calls the ansible script responsible for removing a node from the Opensearch cluster.
//
// Return:
//
//	(error): Returns error if any
func CallScaleDown(username string, hosts string, clusterCfg config.ClusterDetails) error {

	jsonData, err := json.Marshal(&clusterCfg)
	if err != nil {
		log.Error.Println("Error while Marshaling. %v", err)
		return err
	}

	var jsonMap map[string]interface{}

	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		log.Error.Println("json parsing error")
		return err
	}

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		User: username,
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: hosts,
		ExtraVars: jsonMap,
	}

	ansiblePlaybookPrivilegeEscalationOptions := &options.AnsiblePrivilegeEscalationOptions{
		Become: true,
	}

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:                  []string{"ansible_scripts/scaleDownPlaybook.yml"},
		ConnectionOptions:          ansiblePlaybookConnectionOptions,
		PrivilegeEscalationOptions: ansiblePlaybookPrivilegeEscalationOptions,
		Options:                    ansiblePlaybookOptions,
		Exec: execute.NewDefaultExecute(
			execute.WithEnvVar("ANSIBLE_FORCE_COLOR", "true"),
			execute.WithTransformers(
				results.Prepend("Go-ansible with become"),
			),
		),
	}

	err = playbook.Run(context.TODO())
	if err != nil {
		return err
	}
	return nil
}
