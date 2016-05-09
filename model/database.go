package model

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

	"github.com/mephux/kolide/config"

	// Loaded to cover deps for xorm postgres support
	_ "github.com/lib/pq"
)

var (
	x      *xorm.Engine
	tables []interface{}
)

func init() {
	tables = append(tables, new(Node), new(User),
		new(SavedQuery))
}

// NewDatabase will create, connect return a postgres connection pool
func NewDatabase(c *config.Config) (*xorm.Engine, error) {

	dbC := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		c.Database.Username, c.Database.Password,
		c.Database.Address, c.Database.Database, c.Database.SSL)

	engine, err := xorm.NewEngine("postgres", dbC)

	if err != nil {
		return nil, err
	}

	if err := engine.Ping(); err != nil {
		return nil, err
	}

	x = engine

	x.SetMapper(core.GonicMapper{})

	x.ShowSQL = true
	x.ShowInfo = true
	x.ShowDebug = true
	x.ShowErr = true
	x.ShowWarn = true

	w := log.StandardLogger().Writer()
	l := xorm.NewSimpleLogger(w)

	if c.Server.Debug {
		l.SetLevel(core.LOG_DEBUG)
	} else {
		l.SetLevel(core.LOG_ERR)
	}

	x.SetLogger(l)

	if err = x.Sync2(tables...); err != nil {
		return nil, fmt.Errorf("sync database struct error: %v\n", err)
	}

	if _, err := FindSavedQueryById(1); err != nil {
		if err := LoadDefaultSavedQueries(); err != nil {
			log.Error("Unable to load default saved queries.")
			log.Errorf("Error: %s", err)
		}
	}

	if err := CreateUser(&User{
		Enabled:  true,
		Id:       1,
		Name:     "kolide",
		Email:    "example@example.com",
		Password: "password",
		Admin:    true,
	}); err != nil {
		log.Print("Admin account already exists")
	}

	ticker := time.NewTicker(20 * time.Second)
	quit := make(chan struct{})

	log.Info("Starting database workers")

	go func() {
		for {
			select {
			case <-ticker.C:
				go nodeUpdateStatus()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return engine, nil
}
