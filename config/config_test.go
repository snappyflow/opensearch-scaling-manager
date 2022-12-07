package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestClusterName(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: 1cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestClusterIpAddress(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.257, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestClusterOsCredentials(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.254, cluster_name: cluster-1, os_credentials: , cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestClusterCloudCredentials(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.254, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: , base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestClusterCloudType(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AW, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestClusterBaseNodeType(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: , number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestClusterNumCpusPerNode(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 0, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestClusterRAMPerNodeInGB(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 10, ram_per_node_in_gb: 0, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestClusterDiskPerNodeInGB(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 10, ram_per_node_in_gb: 10, disk_per_node_in_gb: 0, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestClusterNumMaxNodesAllowed(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 10, ram_per_node_in_gb: 10, disk_per_node_in_gb: 10, number_max_nodes_allowed: 0}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestTask(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: []}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestTaskName(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestTaskOperator(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: O, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestRule(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: []}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestRuleMetric(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cp, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestRuleStat(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: C, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestRuleDecisionPeriod(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, occurences: 10, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 0}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestRuleOccurences(t *testing.T) {
	yamlString := `{cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: cpu, limit: 1, stat: COUNT, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 9}, {metric: shard, stat: TERM, limit: 900, decision_period: 10}]}]}`
	var config = new(ConfigStruct)
	err := yaml.Unmarshal([]byte(yamlString), &config)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}
	//t.Logf("%#v",config)
	err = validation(*config)
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
}

func TestConfig(t *testing.T) {
	config, err := GetConfig("../config.yaml")
	if err != nil {
		t.Fail()
		t.Logf("expected validation got %v", err)
	}
	t.Log(config)
}
