package provision

import (
	"context"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
)

// Input:
//
//	username (string): Username string to be used to ssh into the host inventory
//	hosts (string): The file name of hosts file to pass to ansible playbook
//
// Description:
//
//	Calls the ansible script responsible for adding a new node into the Opensearch cluster and configuring it.
//
// Return:
//
//	(error): Returns error if any
func CallScaleUp(username string, hosts string) error {

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		User: username,
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: hosts,
	}

	ansiblePlaybookPrivilegeEscalationOptions := &options.AnsiblePrivilegeEscalationOptions{
		Become: true,
	}

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:                  []string{"provision/ansible_scripts/scaleUpPlaybook.yml"},
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

	err := playbook.Run(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

// Input:
//
//	username (string): Username string to be used to ssh into the host inventory
//	hosts (string): The file name of hosts file to pass to ansible playbook
//
// Description:
//
//	Calls the ansible script responsible for removing a node from the Opensearch cluster.
//
// Return:
//
//	(error): Returns error if any
func CallScaleDown(username string, hosts string) error {

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		User: username,
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: hosts,
	}

	ansiblePlaybookPrivilegeEscalationOptions := &options.AnsiblePrivilegeEscalationOptions{
		Become: true,
	}

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:                  []string{"provision/ansible_scripts/scaleDownPlaybook.yml"},
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

	err := playbook.Run(context.TODO())
	if err != nil {
		return err
	}
	return nil
}
