package querycontrol

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/kolide/kolide/model"
	"github.com/kolide/kolide/shared/osquery"
)

// Control global
var Control *control

func init() {
	Control = New()
}

// New query control context
func New() *control {
	return &control{
		queries: make(map[string]*BatchQuery),
		pending: make(map[string]osquery.ReadResp),
	}
}

type control struct {
	queries map[string]*BatchQuery
	pending map[string]osquery.ReadResp
	mutex   sync.RWMutex
}

// Remove batch from control
func (c *control) Remove(batch *BatchQuery) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.queries, batch.ID)

	// remove the batch from all pending queries
	for _, resp := range c.pending {
		delete(resp.Queries, batch.ID)
	}
}

// PendingQueries lists all pending batch queries
func (c *control) PendingQueries(node *model.Node) *osquery.ReadResp {
	queries, ok := c.pending[node.Key]

	log.Debugf("Searching for pending queries for: %s", node.Key)

	if ok {
		log.Debugf("%s pending queries: %+v", node.Key, c.pending)

		// TODO: lock
		c.mutex.Lock()
		defer c.mutex.Unlock()

		delete(c.pending, node.Key)
		return &queries
	}
	return nil
}

// Submit batch query
func (c *control) Submit(batch *BatchQuery) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queries[batch.ID] = batch

	log.Infof("Submitting queries for batch: %s", batch.ID)

	for _, query := range batch.Queries {

		queries, ok := c.pending[query.Node.Key]

		if !ok {
			queries = osquery.ReadResp{
				Queries: make(osquery.QueryType),
				Invalid: false,
			}
		}

		queries.Queries[query.ID] = query.Request
		c.pending[query.Node.Key] = queries
	}
}

// AddResponse to control context
func (c *control) AddResponse(node *model.Node, response *osquery.WriteReq) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	nodeKey := response.Key

	for id, response := range response.Queries {
		batch, ok := c.queries[id]
		log.Debugf("Batch Results for %s (id=%s) ", nodeKey, id)

		if ok {
			err := batch.Done(nodeKey, response)

			if err != nil {
				log.Errorf("Batch Failure (id=%s): %s", id, err)
				return
			}

			log.Debugf("Batch Done (id=%d)", id)
		} else {
			log.Errorf("Missing Query Batch (id=%s)", id)
		}

	}
}
