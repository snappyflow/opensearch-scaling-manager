// This package consists of all the data structure required for defining a task.
// Tasks are set of Actions.
// The actions can have list of rules.
// The recommendation engine will parse these rules and recommend the action if rules meets the criteria.
// Multiple rules can be added inside an action and like wise multiple actions can be added inside a task.
package task

import (
	"encoding/json"
	"fmt"
	"log"
	"scaling_manager/cluster"
)

// This struct contains the task to be perforrmed by the recommendation and set of rules wrt the action.
type Task struct {
	// TaskName indicates the name of the task to recommend by the recommendation engine.
	TaskName string `yaml:"task_name" validate:"required,isValidTaskName"`
	// Rules indicates list of rules to evaluate the criteria for the recomm+endation engine.
	Rules []Rule `yaml:"rules" validate:"gt=0,dive"`
	// Operator indicates the logical operation needs to be performed while executing the rules
	Operator string `yaml:"operator" validate:"required,oneof=AND OR"`
}

// This struct contains the rule.
type Rule struct {
	// Metic indicates the name of the metric. These can be:
	// 	Cpu
	//	Mem
	//	Shard
	Metric string `yaml:"metric" validate:"required,oneof=cpu mem heap disk shard"`
	// Limit indicates the threshold value for a metric.
	// If this threshold is achieved for a given metric for the decision periond then the rule will be activated.
	Limit float32 `yaml:"limit" validate:"required"`
	// Stat indicates the statistics on which the evaluation of the rule will happen.
	// For Cpu and Mem the values can be:
	//		Avg: The average CPU or MEM value will be calculated for a given decision period.
	//		Count: The number of occurences where CPU or MEM value crossed the threshold limit.
	//		Term:
	// For rule: Shard, the stat will not be applicable as the shard will be calculated across the cluster and is not a statistical value.
	Stat string `yaml:"stat" validate:"required,oneof=AVG COUNT TERM"`
	// DecisionPeriod indicates the time in minutes for which a rule is evalated.
	DecisionPeriod int `yaml:"decision_period" validate:"required,min=1"`
	// Occurences indicate the number of time a rule reached the threshold limit for a give decision period.
	// It will be applicable only when the Stat is set to Count.
	Occurences int `yaml:"occurences" validate:"required_if=Stat COUNT"`
}

// This struct contains the task details which is set of actions.
type TaskDetails struct {
	// Tasks indicates list of task.
	// A task indicates what operation needs to be recommended by recommendation engine.
	// As of now tasks can be of two types:
	//
	//	scale_up_by_1
	//	scale_down_by_1
	Tasks []Task `yaml:"action" validate:"gt=0,dive"`
}

// Inputs:
// Caller: Object of TaskDetails
// Description:
//
//		EvaluateTask will go through all the tasks one by one. and
//		It check if the task are meeting the criteria based on rules and operator.
//		If the task is meeting the criteria then it will push the task to recommendation queue.
//
// Return:

func (t TaskDetails) EvaluateTask() {
	for _, v := range t.Tasks {
		isRecommendeTask := v.GetNextTask()
		if isRecommendeTask {
			v.PushToRecommendationQueue()
		} else {
			fmt.Printf("The %s task is not recommended as rules are not satisfied", v.TaskName)
			fmt.Println()
		}
	}
}

// Inputs:
// Caller: Object of Task
// Description:
//
//		GetNextTask will get the Task and iterate through all the rules inside the task.
//		Based on the operator it will check if it should iterate through all the rules or not.
//		It will call GetNextRule while iterating through the rules.
//		Based on the result GetNextTask will check if a task can be recommended or not.
//
// Return:
//
//		Return if a task can be recommended or not(bool)

func (t Task) GetNextTask() bool {
	var isRecommendedTask bool = true
	var isRecommendedRule bool

	// We should have a mechanism to check if we have enough data points for evaluating the rules.
	// If we do not have enough data point for evaluating rule then we should not recommend the task.
	// In case of AND condition if we do not have enough data point for even one rule then the for
	// loop should be broken.
	// This can be considered while implementation.
	if t.Operator == "AND" {
		// Here we will add go routine.
		// So that all the rules getMetrics will be fetched in concurrent way
		for _, v := range t.Rules {
			isRecommendedRule = v.GetNextRule()
			isRecommendedTask = isRecommendedRule && isRecommendedTask
		}
	} else if t.Operator == "OR" {
		// Here we will add go routine.
		// So that all the rules getMetrics will be fetched in concurrent way
		// There is a possibility that each rule is taking time and the matching rule may be at last in or condition
		for _, v := range t.Rules {
			isRecommendedRule = v.GetNextRule()
			if isRecommendedRule {
				break
			}
		}
		isRecommendedTask = isRecommendedRule
	}
	return isRecommendedTask
}

// Input:
// Caller: Object of Rule
// Description:
//
//		GetNextRule will fetch the metrics based on the rules MetricName and Stats using GetMetrics
//		Then it will evaluate if the rule is meeting the criteria or not using EvaluateRule
//
// Return:
//		Return if a rule is meeting the criteria or not(bool)

func (r Rule) GetNextRule() bool {
	cluster := r.GetMetrics()
	isRecommended := r.EvaluateRule(cluster)
	return isRecommended
}

// Input:
// Caller: Object of Rule
// Description:
//
//		GetMetrics will be getting the metrics for a metricName based on its stats
//		If the stat is Avg then it will call GetClusterAvg which will provide MetricViolatedCountCluster struct.
//		If the stat is Count or Term then it will call GetClusterCount which will provide MetricViolatedCountCluster struct.
//		At last it marshal the structure such that uniform data can be used across multiple methods.
//
// Return:
//		Return marshal form of either MetricStatsCluster or MetricViolatedCountCluster struct([]byte)

func (r Rule) GetMetrics() []byte {
	var clusterStats cluster.MetricStats
	var clusterCount cluster.MetricViolatedCount
	var clusterMetric []byte
	var jsonErr error

	if r.Stat == "AVG" {
		clusterStats = cluster.GetClusterAvg(r.Metric, r.DecisionPeriod)
		clusterMetric, jsonErr = json.MarshalIndent(clusterStats, "", "\t")
		if jsonErr != nil {
			log.Fatal("Error converting struct to json: ", jsonErr)
		}
	} else if r.Stat == "COUNT" || r.Stat == "TERM" {
		clusterCount = cluster.GetClusterCount(r.Metric, r.DecisionPeriod, r.Limit)
		clusterMetric, jsonErr = json.MarshalIndent(clusterCount, "", "\t")
		if jsonErr != nil {
			log.Fatal("Error converting struct to json: ", jsonErr)
		}
	}

	return clusterMetric
}

// Input: clusterMetric []byte: Marshal struct containing clusterMetric details based on stats.
// Caller: Object of Rule
// Description:
//
//		EvaluateRule will be compare the collected metric and mentioned rule
//		It will then decide if rules are meeting the criteria or not and return the result.
// Return:
//		Return whether a rule is meeting the criteria or not(bool)

func (r Rule) EvaluateRule(clusterMetric []byte) bool {
	if r.Stat == "AVG" {
		var clusterStats cluster.MetricStats
		err := json.Unmarshal(clusterMetric, &clusterStats)
		if err != nil {
			log.Fatal("Error converting struct to json: ", err)
		}
		if clusterStats.Avg > r.Limit {
			return true
		} else {
			return false
		}
	} else if r.Stat == "COUNT" || r.Stat == "TERM" {
		var clusterStats cluster.MetricViolatedCount
		err := json.Unmarshal(clusterMetric, &clusterStats)
		if err != nil {
			log.Fatal("Error converting struct to json: ", err)
		}
		if r.Stat == "COUNT" {
			if clusterStats.ViolatedCount > r.Occurences {
				return true
			} else {
				return false
			}
		} else if r.Stat == "TERM" {
			if clusterStats.ViolatedCount > int(r.Limit) {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

// Input:
// Caller: Object of Task
// Description:
//
//		PushToRecommendationQueue will be pushing the task which matches the criteria to recommendation queue.
//
// Return:

func (task Task) PushToRecommendationQueue() {
	fmt.Printf("The %s task is recommended and will be pushed to the queue", task.TaskName)
	fmt.Println()
}
