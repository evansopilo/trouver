package main

import (
	"context"
	"fmt"
	"time"

	"github.com/evansopilo/trouver/internal/data"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// A config struct definition to hold all the configuration settings for our
// application. For now, the only configuration settings will be the network
// port that we want the server to listen on, and the name of the current
// operating environment for the application (development, staging, production,
// etc.). We will read in these
// configuration settings from config file when the application starts.
type Config struct {
	// Hold application version number.
	Version string
	// Hold application name.
	App string
	// Hold configuration settings for the server, which will be read from the
	// config file on application start-up.
	Server struct {
		Port string
		Env  string
	}
	// Hold the configuration settings for the database connection pool, which
	// we will read in from config file.
	DB struct {
		DSN string
	}
	// Struct contains fields for the requests-per-second and burst values, and
	// a boolean field which we can use to enable/disable rate limiting
	// altogether.
	Limiter struct {
		RPS     float64
		Burst   int
		Enabled bool
	}
}

// An application struct to hold the dependencies for our HTTP handlers,
// helpers, and middleware. At the moment this only contains a copy of the
// config struct and a logger, but it will grow to include a lot more as our
// build progresses.
type Application struct {
	Config Config
	Logger logrus.Logger
	Models data.Models
}

func main() {
	var cfg Config
	loadConfig(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DB.DSN))
	if err != nil {
		logrus.Fatal(err)
	}

	app := &Application{
		Config: cfg,
		Models: data.Models{
			Place: data.NewPlaceModel(client),
		},
	}

	app.Router().Listen(fmt.Sprintf(":%v", app.Config.Server.Port))
}
