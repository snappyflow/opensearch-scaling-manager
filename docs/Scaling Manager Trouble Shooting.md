### Scaling Manager Trouble Shooting Guide

------

**Scenario 1**

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
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "uninstall" --key-file user-dev-aws-ssh.pem -e src_bin_path="."
```

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
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "install" --key-file user-dev-aws-ssh.pem -e src_bin_path="."
```



**Scenario 2**

If scaling manager fail to start on a node with following issue:

" osConnection.go:76: Unable to ping OpenSearch! Error: 401 Unauthorizedsystemd[1]: scaling_manager.service: Main process exited, code=exited, status=1/FAILURE "

**Explanation**

This may happen when .secret.txt file is not present or deleted. In this case we should run install_scaling_manager.service with update_config tag to update the config file with plain text. As soon as there is a change detected in config file, master node will send back the encrypted config file and secret file on all the nodes and then you can start your application.

**Solution to resolve**

Update Config file and retry

Password based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "update_config" -kK
```

Key based authentication command 

```
sudo ansible-playbook -i inventory.yaml install_scaling_manager.yaml --tags "update_config" --key-file user-dev-aws-ssh.pem -e config_path="config.yaml"
```
