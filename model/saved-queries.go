package model

import (
	"errors"

	log "github.com/Sirupsen/logrus"
)

// SavedQuery database table.
type SavedQuery struct {
	Id    int64  `json:"id"`
	Name  string `xorm:"NOT NULL" from:"name" binding:"required"`
	Query string `xorm:"NOT NULL" from:"query" binding:"required"`
	Type  string `xorm:"NOT NULL" form:"type" binding:"required"`
}

// LoadDefaultSavedQueries Load all of the default saved queries. This will
// only run once on initial database creation.
func LoadDefaultSavedQueries() error {
	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	queries := []SavedQuery{}

	q1 := SavedQuery{
		Name:  "Process List",
		Query: "select * from processes;",
		Type:  "all",
	}

	queries = append(queries, q1)

	q2 := SavedQuery{
		Name:  "All listening ports joined with processes",
		Query: "select * from listening_ports join processes using (pid);",
		Type:  "all",
	}

	queries = append(queries, q2)

	q3 := SavedQuery{
		Name:  "All currently executing processes where the original binary no longer exists",
		Query: "SELECT name, path, pid FROM processes WHERE on_disk = 0;",
		Type:  "all",
	}

	queries = append(queries, q3)

	q4 := SavedQuery{
		Name: "All processes that are listening on network ports",
		Query: `SELECT DISTINCT process.name, listening.port, listening.address, process.pid
FROM processes AS process 
JOIN listening_ports 
AS listening ON process.pid = listening.pid;`,
		Type: "all",
	}

	queries = append(queries, q4)

	q5 := SavedQuery{
		Name:  "Third-party kernel extensions (OS X)",
		Query: "SELECT * FROM kernel_extensions WHERE name NOT LIKE 'com.apple.%' AND name != '__kernel__';",
		Type:  "all",
	}

	queries = append(queries, q5)

	q6 := SavedQuery{
		Name:  "Startup items (OS X / LaunchDaemons & LaunchAgents)",
		Query: `SELECT disabled, path, program FROM launchd;`,
		Type:  "all",
	}

	queries = append(queries, q6)

	q7 := SavedQuery{
		Name:  "Shell history",
		Query: `SELECT * FROM shell_history;`,
		Type:  "all",
	}

	queries = append(queries, q7)

	q8 := SavedQuery{
		Name:  "All users with group information",
		Query: `SELECT * FROM users u JOIN groups g where u.gid = g.gid;`,
		Type:  "all",
	}

	queries = append(queries, q8)

	q9 := SavedQuery{
		Name: "Interface information",
		Query: `SELECT address, mac, id.interface
FROM interface_details AS id, interface_addresses AS ia WHERE id.interface = ia.interface;`,
		Type: "all",
	}

	queries = append(queries, q9)

	// q8 := SavedQuery{
	// Name: "All empty groups",
	// Query: `SELECT groups.gid, groups.name FROM groups
	// LEFT JOIN users ON (groups.gid = users.gid) WHERE users.uid IS NULL;`,
	// Type: "all",
	// }

	// queries = append(queries, q8)

	if _, err := sess.Insert(&queries); err != nil {
		sess.Rollback()
		return err
	}

	err := sess.Commit()

	if err != nil {
		return err
	}

	return nil
}

// FindSavedQueryById a saved query by its id.
func FindSavedQueryById(id int64) (*SavedQuery, error) {
	log.Debugf("Looking for saved query with id: %d", id)

	query := &SavedQuery{Id: id}

	has, err := x.Get(query)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("Saved Query not found")
	}

	return query, nil
}

// Delete a saved query
func (s *SavedQuery) Delete() error {
	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	if _, err := sess.Delete(&SavedQuery{Id: s.Id}); err != nil {
		log.Debug("Saved Query Delete Error: ", err)
		return err
	}

	err := sess.Commit()

	if err != nil {
		return err
	}

	return nil
}

// AllSavedQueries returns all saved queries in the database.
func AllSavedQueries() ([]*SavedQuery, error) {
	var data []*SavedQuery
	err := x.Find(&data)

	return data, err
}

// NewSavedQuery will iadd a new saved query to the database.
func NewSavedQuery(s SavedQuery) error {
	sess := x.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	if _, err := sess.Insert(s); err != nil {
		sess.Rollback()
		return err
	}

	err := sess.Commit()

	if err != nil {
		return err
	}

	return nil
}
