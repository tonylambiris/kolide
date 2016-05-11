package main

import (
	"fmt"
	"os"
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
	app = kingpin.New(version.Name, version.Description)

	debug = app.Flag("debug", "Enable debug mode.").OverrideDefaultFromEnvar("KOLIDE_DEBUG").Bool()
	quiet = app.Flag("quiet", "Remove all output logging").Short('q').Bool()

	// server sub-commands
	// cmdServer = app.Command("server", fmt.Sprintf("Run and control the %s web server.", version.Name))

	dev                   = app.Flag("dev", "enable dev mode (serve assets from disk)").OverrideDefaultFromEnvar("KOLIDE_DEV").Bool()
	configPath            = app.Flag("config", "configuration file").Short('c').Required().OverrideDefaultFromEnvar("KOLIDE_CONFIG_PATH").ExistingFile()
	cmdServerProduction   = app.Flag("production", "enable production mode").OverrideDefaultFromEnvar("KOLIDE_PRODUCTION").Bool()
	cmdServerAddress      = app.Flag("address", "web server network address").PlaceHolder(":8000").OverrideDefaultFromEnvar("KOLIDE_SERVER_ADDRESS").String()
	cmdServerEnrollSecret = app.Flag("enroll-secret", "osquery enroll secret").PlaceHolder("secret").OverrideDefaultFromEnvar("KOLIDE_ENROLL_SECRET").String()
	// cmdServerQueryTimeout = app.Flag("query-timeout", "Query timeout duration (10s)").OverrideDefaultFromEnvar("KOLIDE_QUERY_TIMEOUT").String()

	cmdDatabaseAddress  = app.Flag("db-address", "database network address").PlaceHolder(":5432").OverrideDefaultFromEnvar("KOLIDE_DATABASE_ADDRESS").String()
	cmdDatabaseUsername = app.Flag("db-username", "database username").PlaceHolder(version.Name).OverrideDefaultFromEnvar("KOLIDE_DATABASE_USERNAME").String()
	cmdDatabasePassword = app.Flag("db-password", "database password").PlaceHolder("secret").OverrideDefaultFromEnvar("KOLIDE_DATABASE_PASSWORD").String()
	cmdDatabaseDatabase = app.Flag("db-database", "database database").PlaceHolder(version.Name).OverrideDefaultFromEnvar("KOLIDE_DATABASE_DATABASE").String()

	cmdRedisAddress   = app.Flag("redis-address", "redis network address").PlaceHolder(":6379").OverrideDefaultFromEnvar("KOLIDE_REDIS_ADDRESS").String()
	cmdRedisProtocol  = app.Flag("redis-protocol", "redis network protocol").OverrideDefaultFromEnvar("KOLIDE_REDIS_PROTOCOL").Default("tcp").String()
	cmdRedisSize      = app.Flag("redis-size", "redis maximum number of idle connections").PlaceHolder("10").OverrideDefaultFromEnvar("KOLIDE_REDIS_SIZE").Int()
	cmdRedisPassword  = app.Flag("redis-password", "redis password").PlaceHolder("secret").OverrideDefaultFromEnvar("KOLIDE_REDIS_PASSWORD").String()
	cmdRedisSecretKey = app.Flag("redis-secret-key", "redis secret key").PlaceHolder("secret").OverrideDefaultFromEnvar("KOLIDE_REDIS_SECRET_KEY").String()
	// cmdRedisEncryptionKey = app.Flag("redis-encryption-key", "redis encryption key").PlaceHolder("secret").OverrideDefaultFromEnvar("KOLIDE_REDIS_ENCRYPTION_KEY").String()

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

	app.Version(v)
	args, err := app.Parse(os.Args[1:])

	// i may add sub-commands for user management
	switch kingpin.MustParse(args, err) {
	default:
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

		// merge with cli args
		if *cmdServerProduction {
			c.Server.Production = true
		}

		if len(*cmdServerAddress) > 0 {
			c.Server.Address = *cmdServerAddress
		}

		if len(*cmdServerEnrollSecret) > 0 {
			c.Server.EnrollSecret = *cmdServerEnrollSecret
		}

		// if len(*cmdServerQueryTimeout) > 0 {
		// c.Server.QueryTimeout = config.Duration(*cmdServerQueryTimeout)
		// }

		if len(*cmdDatabaseAddress) > 0 {
			c.Database.Address = *cmdDatabaseAddress
		}

		if len(*cmdDatabaseUsername) > 0 {
			c.Database.Username = *cmdDatabaseUsername
		}

		if len(*cmdDatabasePassword) > 0 {
			c.Database.Password = *cmdDatabasePassword
		}

		if len(*cmdDatabaseDatabase) > 0 {
			c.Database.Database = *cmdDatabaseDatabase
		}

		// REDIS

		if len(*cmdRedisAddress) > 0 {
			c.Session.Address = *cmdRedisAddress
		}

		if *cmdRedisSize > 0 {
			c.Session.Size = *cmdRedisSize
		}

		if len(*cmdRedisProtocol) > 0 {
			c.Session.Network = *cmdRedisProtocol
		}

		if len(*cmdRedisPassword) > 0 {
			c.Session.Password = *cmdRedisPassword
		}

		if len(*cmdRedisSecretKey) > 0 {
			c.Session.Key = *cmdRedisSecretKey
		}

		// if len(*cmdRedisEncryptionKey) > 0 {
		// c.Session.EncryptionKey = *cmdRedisEncryptionKey
		// }
		// merge complete

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

}
