package model

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/mephux/kolide/shared/hub"
	"github.com/mephux/kolide/shared/osquery"
)

// nodeUpdateStatus loop to keep node information current
// this loop is used in database.go
func nodeUpdateStatus() error {
	nodes, err := AllNodes()

	if err != nil {
		return err
	}

	for _, node := range nodes {
		if node.Updated.Unix() <= (time.Now().Unix() - 30) {
			node.Enabled = false
			node.Update()
		}
	}

	return nil
}

// Node database table schema
type Node struct {
	Id int64 `json:"id"`

	Key     string `xorm:"UNIQUE NOT NULL" json:"key"`
	Address string `json:"address"`
	Name    string `xorm:"INDEX NOT NULL" json:"name"`
	Enabled bool   `xorm:"INDEX NOT NULL" json:"enabled"`

	Created time.Time `xorm:"CREATED" json:"created"`
	Updated time.Time `xorm:"UPDATED" json:"updated"`
}

// AllNodes Return all nodes in the database
func AllNodes() ([]*Node, error) {
	var nodes []*Node
	err := x.Find(&nodes)

	return nodes, err
}

// FindNodeByRequest takes the osquery request and returns a node
// from the database
func FindNodeByRequest(c *gin.Context, req *osquery.KeyReq) (*Node, error) {
	req.Address = strings.Split(c.ClientIP(), ":")[0]
	return FindAndUpdateNode(req)
}

// CreateOrUpdateNode will create or update a node from
// a enroll request
func CreateOrUpdateNode(req *osquery.EnrollReq) (*Node, error) {
	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return nil, err
	}

	node, err := FindNodeByNodeKey(req.Key)

	if node != nil {
		log.Debug("Found existing node record.")

		node.Address = req.Address
		node.Enabled = true

		if _, err := sess.Id(node.Id).AllCols().Update(node); err != nil {
			sess.Rollback()
			return node, err
		}

		err := sess.Commit()

		if err != nil {
			return nil, err
		}

		return node, nil
	}

	log.Debugf("Error: %s", err)
	log.Debugf("Creating a new node record.")

	node = &Node{
		Key:     req.Key,
		Name:    req.Key,
		Address: req.Address,
		Enabled: true,
	}

	if _, err := sess.Insert(node); err != nil {
		sess.Rollback()
		return nil, err
	}

	err = sess.Commit()

	if err != nil {
		return nil, err
	}

	return node, nil

}

// FindAndUpdateNode node by node key which is also the osquery host id
func FindAndUpdateNode(req *osquery.KeyReq) (*Node, error) {
	log.Debugf("Looking for node with key: %s", req.Key)

	node := &Node{Key: req.Key}

	has, err := x.Get(node)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("node not found")
	}

	node.Address = req.Address
	// node.Updated = time.Now()
	err = node.Update()

	msg := hub.Message{
		Type: "node",
		Data: node,
	}

	hub.Websocket.Broadcast <- msg.JSON()

	return node, err
}

// JSON node format
func (n *Node) JSON() []byte {
	b, _ := json.Marshal(n)
	return b
}

// FindNodeByNodeKey node by node key which is also the osquery host id
func FindNodeByNodeKey(key string) (*Node, error) {
	log.Debugf("Looking for node with key: %s", key)

	node := &Node{Key: key}

	has, err := x.Get(node)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("node not found")
	}

	return node, nil
}

// Update node information
func (n *Node) Update() error {
	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	if _, err := sess.Id(n.Id).AllCols().Update(n); err != nil {
		sess.Rollback()
		return err
	}

	err := sess.Commit()

	if err != nil {
		return err
	}

	return nil
}

// Delete node from the database.
func (n *Node) Delete() error {
	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	if _, err := sess.Id(n.Id).Delete(n); err != nil {
		sess.Rollback()
		return err
	}

	err := sess.Commit()

	if err != nil {
		return err
	}

	return nil
}
