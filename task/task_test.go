package task

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestTaskNotRecommendedOr(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg/cpu/9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"Avg": 1,
				"Min": 0,
				"Max": 1,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg/mem/9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"Avg": 1,
				"Min": 0,
				"Max": 1,
			})
			return resp, err
		},
	)

	isRecommendedTask := task.GetNextTask()
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskRecommendedOr(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: OR, rules: [{metric: cpu, limit: 1, stat: AVG, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg/cpu/9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 4,
				"min": 0,
				"max": 1,
			})
			return resp, err
		},
	)

	isRecommendedTask := task.GetNextTask()
	t.Log(isRecommendedTask)
	assert.Equal(t, true, isRecommendedTask)
}

func TestTaskNotRecommendedAnd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: AND, rules: [{metric: cpu, limit: 2, stat: AVG, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg/cpu/9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"Avg": 1,
				"Min": 0,
				"Max": 1,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg/mem/9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"Avg": 1,
				"Min": 0,
				"Max": 1,
			})
			return resp, err
		},
	)

	isRecommendedTask := task.GetNextTask()
	t.Log(isRecommendedTask)
	assert.Equal(t, false, isRecommendedTask)
}

func TestTaskRecommendedAnd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yamlString := `{task_name: scale_up_by_1, operator: AND, rules: [{metric: cpu, limit: 1, stat: AVG, decision_period: 9}, {metric: mem, limit: 59, stat: AVG, decision_period: 9}]}`
	var task = new(Task)
	err := yaml.Unmarshal([]byte(yamlString), &task)
	if err != nil {
		t.Fail()
		t.Logf("failed to unmarshal yaml: %v", err.Error())
	}

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg/cpu/9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 4,
				"min": 0,
				"max": 1,
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://localhost:5000/stats/avg/mem/9",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"avg": 60,
				"min": 0,
				"max": 1,
			})
			return resp, err
		},
	)

	isRecommendedTask := task.GetNextTask()
	t.Log(isRecommendedTask)
	assert.Equal(t, true, isRecommendedTask)
}
