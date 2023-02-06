package provision

import (
	"context"
	"io/ioutil"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"gopkg.in/yaml.v2"
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

	// Create varsFile.yml
	yamlData, err := yaml.Marshal(&clusterCfg)
	if err != nil {
		log.Error.Println("Error while Marshaling. %v", err)
		return err
	}

	varsFile := "varsFile.yaml"
	err = ioutil.WriteFile(varsFile, yamlData, 0644)
	if err != nil {
		log.Panic.Println("Error while Marshaling. %v", err)
		return err
	}

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		User: username,
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory:     hosts,
		ExtraVarsFile: []string{"@" + varsFile},
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

	// Create varsFile.yml
	yamlData, err := yaml.Marshal(&clusterCfg)
	if err != nil {
		log.Error.Println("Error while Marshaling. %v", err)
		return err
	}

	varsFile := "varsFile.yaml"
	err = ioutil.WriteFile(varsFile, yamlData, 0644)
	if err != nil {
		log.Panic.Println("Error while Marshaling. %v", err)
		return err
	}

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		User: username,
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory:     hosts,
		ExtraVarsFile: []string{"@" + varsFile},
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
