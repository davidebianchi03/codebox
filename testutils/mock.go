package testutils

import "github.com/gocraft/work"

type MockEnqueuer struct {
	CalledJobs []string
}

func (m *MockEnqueuer) Enqueue(jobName string, args map[string]interface{}) (*work.Job, error) {
	m.CalledJobs = append(m.CalledJobs, jobName)
	return &work.Job{}, nil
}
