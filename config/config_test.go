package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMonitorWithLogs(t *testing.T) {
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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

func TestMonitorWithSimulator(t *testing.T) {
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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

func TestPollingIntervalSecs(t *testing.T) {
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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

func TestClusterName(t *testing.T) {
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
	yamlString := `{user_config: {monitor_with_logs: true, monitor_with_simulator: false, purge_old_docs_after_hours: 50, polling_interval_in_secs: 10, is_accelerated: false}, cluster_details: {ip_address: 10.81.1.225, cluster_name: cluster-1, os_credentials: {os_admin_username: elastic, os_admin_password: changeme}, os_user: ubuntu ,os_version: 2.3.0, os_home: /usr/share/opensearch, domain_name: snappyflow.com, cloud_type: AWS, cloud_credentials: {secret_key: secret_key, access_key: access_key}, base_node_type: t2x.large, number_cpus_per_node: 5, ram_per_node_in_gb: 10, disk_per_node_in_gb: 100, number_max_nodes_allowed: 2}, task_details: [{task_name: scale_up_by_1, operator: OR, rules: [{metric: CpuUtil, limit: 2, stat: AVG, decision_period: 9}, {metric: CpuUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}, {metric: RamUtil, limit: 1, stat: COUNT, occurrences: 10, decision_period: 9}]}]}`
	config := new(ConfigStruct)
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
