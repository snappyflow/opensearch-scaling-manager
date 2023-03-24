# Scaling Manager Trouble Shooting Guide

- [Scaling Manager Trouble Shooting Guide](#scaling-manager-trouble-shooting-guide)
  - [Scenario 1](#scenario-1)
  - [Scenario 2](#scenario-2)
  - [Scenario 3](#scenario-3)
  - [Scenario 4](#scenario-4)
  - [Scenario 5](#scenario-5)


## Scenario 1

Installation of Scaling Manager not completed successfully on a new node which is added into the existing cluster.

**Explanation**

Let's assume there is 5node in a cluster and scaling manager spins up a new node and joins the new node to the cluster. Scaling manager and OpenSearch will be installed in new node, during installation when the master node goes down and the newly elected node becomes as master. Now the installation is not yet completed and the new node has become as master which causes this issue of not completing the installation.

**Solution to resolve**

**Step 1** - Run populate inventory.yaml command 

```master_node_ip
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "populate_inventory_yaml" -e master_node_ip=0.0.0.0 -e os_user=USERNAME -e os_pass=PASSWORD
```

master_node_ip = IP address of master node,
os_user = Appropriate username,
os_pass = Appropriate password

**Step 2** - Commenting IP of already present nodes in cluster except new node.

- Open inventory.yaml, comment already present nodes in cluster except new node which got added.
- This commenting of IP is to avoid installation of scaling manager, OpenSearch and other dependencies again on the nodes which has already installed in it.

**Step 2** - Installation

**Uninstall and Installation**

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

**Install**

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



## Scenario 2

If scaling manager fail to start on a node with following issue:

" osConnection.go:76: Unable to ping OpenSearch! Error: 401 Unauthorizedsystemd[1]: scaling_manager.service: Main process exited, code=exited, status=1/FAILURE "

**Explanation**

This may happen when .secret.txt file is not present or deleted. In this case we should run install_scaling_manager.service with update_config tag to update the config file with plain text. As soon as there is a change detected in config file, master node will send back the encrypted config file and secret file on all the nodes and then you can start your application.

**Solution to resolve**

**Step 1** - Run populate inventory.yaml command 

```master_node_ip
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "populate_inventory_yaml" -e master_node_ip=0.0.0.0 -e os_user=USERNAME -e os_pass=PASSWORD
```

master_node_ip = IP address of master node,
os_user = Appropriate username,
os_pass = Appropriate password

**Step 2-** Update Config file and retry

- Password based authentication command 

  ```
  sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "update_config" -kK
  ```

- Key based authentication command 

  ```
  sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "update_config" --key-file USERPEMFILEPATH.pem -e config_path="config.yaml"
  ```

USERPEMFILEPATH = Please provide appropriate pem file path here.



## Scenario 3

Config update failed due to node unreachable

**Explanation**

This may happen due to the case where the entry for that host is deleted or not added in case of addition of node.

**Solution to resolve**

Step 1 - Add the nodes in "etc hosts" manually 

Step 2 - Update Config file

- Password based authentication command 

  ```
  sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "update_config" -kK
  ```

- Key based authentication command 

  ```
  sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "update_config" --key-file USERPEMFILEPATH.pem -e config_path="config.yaml"
  ```

USERPEMFILEPATH = Please provide appropriate pem file path here.



## Scenario 4

Failed to connect to the host via ssh. Node is unreachable while updating the config and secret or spinning up a new node from master node. 

**Explanation**

This may happen due to the ssh is not allowed because the underlying security group does not support ssh.

**Solution to resolve**

Add the correct source in the security group which allow ssh.



## Scenario 5

Node will not be added due to the security configuration.

**Explanation and Solution to resolve**

Please go back and check if the wild card is enabled on configuration of CN on all the node.
