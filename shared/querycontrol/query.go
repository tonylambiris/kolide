package querycontrol

import (
	"fmt"
	"sync"
	"time"

	"github.com/mephux/common/uuid"
	"github.com/mephux/kolide/model"
)

var (
	// TIMEOUT for query/batch processing
	TIMEOUT = 10 * time.Second
)

type QueryResult struct {
	Node     *model.Node `json:"node"`
	Response interface{} `json:"results"`
	TimedOut bool        `json:"timeout"`
}

// BatchQuery is a query
type BatchQuery struct {
	ID string
	// map of node key -> Query
	Queries map[string]*Query

	wg sync.WaitGroup
}

func NewBatchQuery(request string, nodes []*model.Node) *BatchQuery {
	id := uuid.NewV4().String()
	queries := make(map[string]*Query)

	for _, node := range nodes {
		query := NewQuery(id, node, request)
		queries[node.Key] = query
	}

	return &BatchQuery{
		ID:      id,
		Queries: queries,
	}
}

func (q *BatchQuery) Run(timeout time.Duration) map[string]*QueryResult {
	results := make(map[string]*QueryResult)

	for _, query := range q.Queries {
		q.wg.Add(1)
		go func(query *Query) {
			queryResults := query.WaitForResults(timeout)
			results[query.Node.Key] = queryResults
			q.wg.Done()
		}(query)
	}

	Control.Submit(q)

	q.wg.Wait()

	Control.Remove(q)

	return results
}

func (q *BatchQuery) Done(nodeKey string, result interface{}) error {

	query, ok := q.Queries[nodeKey]

	if ok {
		query.Finish(result)
		return nil
	}

	return fmt.Errorf("No node key found for id=%s key=%s", q.ID, nodeKey)
}

type Query struct {
	ID      string
	Node    *model.Node
	Request string
	Result  interface{}

	Done chan bool
}

func NewQuery(id string, node *model.Node, request string) *Query {
	return &Query{
		ID:      id,
		Node:    node,
		Request: request,
		Done:    make(chan bool, 1),
	}
}

func (q *Query) WaitForResults(timeout time.Duration) *QueryResult {

	timedOut := false
	select {
	case <-q.Done:
		break
	case <-time.After(timeout):
		timedOut = true
		break
	}

	if timedOut {
		return &QueryResult{
			Node:     q.Node,
			TimedOut: true,
		}
	}

	return &QueryResult{
		Response: q.Result,
		Node:     q.Node,
	}
}

func (q *Query) Finish(result interface{}) {
	q.Result = result
	q.Done <- true
}
