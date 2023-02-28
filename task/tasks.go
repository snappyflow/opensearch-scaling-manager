package task

// This struct contains the task to be perforrmed by the recommendation and set of rules wrt the action.
type Task struct {
	// TaskName indicates the name of the task to recommend by the recommendation engine.
	TaskName string `yaml:"task_name" validate:"required,isValidTaskName"`
	// Rules indicates list of rules to evaluate the criteria for the recomm+endation engine.
	Rules []Rule `yaml:"rules" validate:"gt=0,dive"`
	// Operator indicates the logical operation needs to be performed while executing the rules
	Operator string `yaml:"operator" validate:"required,oneof=AND OR EVENT"`
}

// This struct contains the rule.
type Rule struct {
	// Metic indicates the name of the metric. These can be:
	//      Cpu
	//      Mem
	//      Shard
	Metric string `yaml:"metric" validate:"required,oneof=CpuUtil RamUtil HeapUtil DiskUtil NumShards"`
	// Limit indicates the threshold value for a metric.
	// If this threshold is achieved for a given metric for the decision periond then the rule will be activated.
	Limit float32 `yaml:"limit" validate:"required"`
	// Stat indicates the statistics on which the evaluation of the rule will happen.
	// For Cpu and Mem the values can be:
	//              Avg: The average CPU or MEM value will be calculated for a given decision period.
	//              Count: The number of occurences where CPU or MEM value crossed the threshold limit.
	//              Term:
	// For rule: Shard, the stat will not be applicable as the shard will be calculated across the cluster and is not a statistical value.
	Stat string `yaml:"stat" validate:"required,oneof=AVG COUNT TERM"`
	// DecisionPeriod indicates the time in minutes for which a rule is evalated.
	DecisionPeriod int `yaml:"decision_period" validate:"required,min=1"`
	// Occurrences indicate the number of time a rule reached the threshold limit for a give decision period.
	// It will be applicable only when the Stat is set to Count.
	Occurrences int `yaml:"occurrences" validate:"required_if=Stat COUNT"`
	// Scheduling time indicates cron time expression to schedule scaling operations
	// Example:
	// SchedulingTime = "30 5 * * 1-5"
	// In the above example the cron job will run at 5:30 AM from Mon-Fri of every month
	SchedulingTime string `yaml:"scheduling_time" validate:"required_if=Operator EVENT"`
	// NumNodesRequired specifies the integer value of number of nodes to be present in cluster for event based scaling operations
	NumNodesRequired int `yaml:"number_of_node" validate:"required_if=Operator EVENT"`
}

// This struct contains the task details which is set of actions.
type TaskDetails struct {
	// Tasks indicates list of task.
	// A task indicates what operation needs to be recommended by recommendation engine.
	// As of now tasks can be of two types:
	//
	//      scale_up_by_1
	//      scale_down_by_1
	Tasks []Task `yaml:"action" validate:"gt=0,dive"`
}
