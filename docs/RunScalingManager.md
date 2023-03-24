## Scaling Manager Pre-Requisites

- Cluster with OpenSearch installed.
- OpenSearch version - 1.2.4 and above. 
- Go version - 1.19
- Ansible Version - 2.9.
- Cluster credentials (Username, Password) to access the OpenSearch.
- Cloud credential  (Username, Password). 
- In AWS we can create a instance by templates which is provided by Domain that is used.
- Launch Template - AWS launch template to spin a new node which has the necessary tags.
- Template ID format (lt-xxxxxxxxxxxxxxxxx.)
- Security certificate to have regex in it to accept the new node.
- PEM file.
- SSH aspect - If cloud type is AWS then Security group is configured in such a way that newly spin up node should be reached via ssh.
- Sudo permission - All the nodes, jump host should have sudo permission by which task could be performed with sudo access between nodes and run ansible playbook with sudo access on jump host. Sudo password can be empty which is preferable.



## SSH details

- Download any remote computing toolbox like MobaXterm. 
- Click Session -> SSH -> Remote host.
- Enter Remote host details, mention the username.
- Click Advanced SSH Settings, choose the PEM file that you have for credentials and click OK.
- Login using the Cluster and Jump host details.
- Once you login, download the latest build and update the following files 
  1. GNUmakefile 
  2. install_scaling_manager
  3. scaling_manager.tar.gz
- Now you are ready with your latest build to see the actual working of Scaling Manager.



## Scripts to run Scaling Manager

**Inventory file** -  Defines the hosts and groups of hosts upon which commands, modules, and tasks in a playbook operate.

**Populate inventory.yaml**

User can should mention the master node IP address in populate inventory.yaml command. When you mention the master node IP it will collect the details about all the nodes present in the cluster. 

For example lets say there is a cluster with 5 nodes. Following are the IP address of nodes,

- Node1 IP - 10.10.10.1 (Master node)
- Node2 IP - 10.10.10.2 (Non-Master node)
- Node3 IP - 10.10.10.3 (Non-Master node)
- Node4 IP - 10.10.10.4 (Non-Master node)
- Node5 IP - 10.10.10.5 (Non-Master node)

Initially inventory.yaml will be empty. Now, when you run the populate inventory command with master node IP all the IP address of the cluster will be populated inside inventory.yaml.

```master_node_ip
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "populate_inventory_yaml" -e master_node_ip=0.0.0.0 -e os_user=USERNAME -e os_pass=PASSWORD
```

master_node_ip = IP address of master node,
os_user = Appropriate username,
os_pass = Appropriate password

**Build, Pack**

In case to use password based authentication, Use the following command

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "build_and_pack" -kK
```

-kK is used for password authentication

In case to use key based authentication, Use the following command

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "update_pem" --key-file USERPEMFILEPATH.pem -e pem_path="user-dev-aws-ssh.pem"
```

USERPEMFILEPATH = Please provide appropriate pem file path here.

Key based and password based commands does the same work but the way you do is different. You can use either of the commands to execute.

1. In password based authentication you should mention the credentials through which the command executes.
2. In key based authentication you can specify the pem file path which has all the credentials so that you need not mention the credentials in command.

**Installation**

Password based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "install" -kK
```

- Update should be performed when provision is not happening, then stop, install, start the service. These steps can be performed in the same command as well.  

- Stop, Install, Start in same command 

  ```
  sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "stop,install,start" -kK
  ```

Key based authentication command

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "install" --key-file USERPEMFILEPATH.pem -e src_bin_path="."
```

USERPEMFILEPATH = Please provide appropriate pem file path here.

**Update Config**

Password based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "update_config" -kK
```

Key based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "update_config" --key-file USERPEMFILEPATH.pem -e config_path="config.yaml"
```

USERPEMFILEPATH = Please provide appropriate pem file path here.

**Update pem**

Password based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "update_pem" -kK
```

Key based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "update_pem" --key-file USERPEMFILEPATH.pem -e pem_path="user-dev-aws-ssh.pem"
```

USERPEMFILEPATH = Please provide appropriate pem file path here.

**Start**

Password based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "start" -kK
```

Key based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "start" --key-file USERPEMFILEPATH.pem -e src_bin_path="."
```

USERPEMFILEPATH = Please provide appropriate pem file path here.

**Stop**

Password based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "stop" -kK
```

Key based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "stop" --key-file USERPEMFILEPATH.pem -e src_bin_path="."
```

USERPEMFILEPATH = Please provide appropriate pem file path here.

- Stop command works quick when there is no provisioning happening/provisioning is completed.
- When provisioning is in process and the stop command is executed it waits till provisioning is completed. To know the status user can do Ctrl+C and check the status of the cluster. 

**Status**

Password based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "status" -kK
```

Key based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "status" --key-file USERPEMFILEPATH.pem -e src_bin_path="."
```

USERPEMFILEPATH = Please provide appropriate pem file path here.

**Uninstall**

Password based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "uninstall" -kK
```

Key based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "uninstall" --key-file USERPEMFILEPATH.pem -e src_bin_path="."
```

USERPEMFILEPATH = Please provide appropriate pem file path here.
