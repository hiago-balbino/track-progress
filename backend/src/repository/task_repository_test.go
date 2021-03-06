package repository

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

type stubSession struct {
}

type stubQuery struct {
}

type stubIter struct {
}

func (i stubIter) Close() error {
	return nil
}

var called = false
var uuid, _ = gocql.RandomUUID()

func (i stubIter) MapScan(m map[string]interface{}) bool {
	if !called {
		m["id"] = uuid
		m["title"] = "nice_title"
		m["description"] = "some_description"
		m["status"] = "active"
		called = true
		return true
	}
	return false
}

func (q stubQuery) Exec() error {
	return nil
}
func (q stubQuery) Iter() IterInterface {
	return stubIter{}
}

func (s stubSession) Query(string, ...interface{}) QueryInterface {
	return stubQuery{}
}

func (s stubSession) Close() {
}

func TestGetAll(t *testing.T) {
	logger, _ := test.NewNullLogger()
	session := stubSession{}
	repo := NewTaskRepository(logger, session)
	tasks := repo.GetAll()

	assert.Equal(t, uuid, tasks[0].ID)
	assert.Equal(t, "nice_title", tasks[0].Title)
	assert.Equal(t, "some_description", tasks[0].Description)
	assert.Equal(t, "active", tasks[0].Status)
}
