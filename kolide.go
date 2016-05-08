package main

import (
	"fmt"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/mephux/kolide/config"
	"github.com/mephux/kolide/model"
	"github.com/mephux/kolide/router"
	"github.com/mephux/kolide/server"
	"github.com/mephux/kolide/shared/formatter"
	"github.com/mephux/kolide/version"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug      = kingpin.Flag("debug", "Enable debug mode.").Bool()
	configPath = kingpin.Flag("config", "Configuration file").Short('c').Required().ExistingFile()
	dev        = kingpin.Flag("dev", "Run in dev mode. (serve assets from disk)").Bool()

	// build information
	build string
)

func main() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)

	v := version.Version

	if len(build) > 0 {
		v = fmt.Sprintf("%s+%s", version.Version, build)
		version.VersionBuild = build
	}

	kingpin.Version(v)
	kingpin.Parse()

	logrus.SetFormatter(new(prefixed.TextFormatter))
	logrus.SetLevel(logrus.InfoLevel)

	var c *config.Config
	var err error

	if len(*configPath) > 0 {

		c, err = config.Load(*configPath)

		if err != nil {
			logrus.Fatal(err)
		}

	} else {
		c = config.Default(*debug, !*dev)
	}

	// debug level if requested by user
	if *debug {
		c.Server.Debug = true
		logrus.SetLevel(logrus.DebugLevel)
	}

	c.Server.Production = !*dev

	if c.Server.Production {
		logrus.Info("ENV: Production")
		gin.SetMode(gin.ReleaseMode)
	} else {
		logrus.Info("ENV: Development")
		gin.SetMode(gin.DebugMode)
	}

	db, err := model.NewDatabase(c)

	if err != nil {
		logrus.Fatal("Database Error: ", err)
	}

	defer db.Close()

	s := server.Load(c)
	s.Run(router.Load(c))
}
