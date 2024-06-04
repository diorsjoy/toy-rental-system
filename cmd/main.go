package main

import (
	"context"
	"database/sql"
	"flag"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"sync"
	"time"
	pkg "toy-rental-system/pkg/jsonlog"
	"toy-rental-system/pkg/utils"
)

// Add a db struct field to hold the configuration settings for our database connection
// pool. For now this only holds the DSN, which we will read in from a command-line flag.
// config struct to hold all the configuration settings for our application.
// Add maxOpenConns, maxIdleConns and maxIdleTime fields to hold the configuration
// settings for the connection pool.
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	// Add a new limiter struct containing fields for the requests-per-second and burst
	// values, and a boolean field which we can use to enable/disable rate limiting
	// altogether.
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	// Update the config struct to hold the SMTP server settings.
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware
// Add a models field to hold our new Models struct.

// Include a sync.WaitGroup in the application struct. The zero-value for a
// sync.WaitGroup type is a valid, useable, sync.WaitGroup with a 'counter' value of 0,
// so we don't need to do anything else to initialize it before we can use it.
type application struct {
	config config
	logger pkg.Logger
	//models data.Models
	//mailer mailer.Mailer
	wg sync.WaitGroup
	//middleware
}

func main() {
	// Declare an instance of the config struct.
	var cfg config

	env, err := utils.LoadConfig("toy-rental-system/config/app.env")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	stripeKey := env.StripeSecret
	dbHost := env.DBSource
	smtpUser := env.SMTPUsername
	smtpPass := env.SMTPPassword
	rabbitMQHost := env.RabbitMQSource

	//r := mux.NewRouter()

	// Read the value of the port and env command-line flags into the config struct. We
	// default to using the port number 4000 and the environment "development"
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// Read the DSN value from the db-dsn command-line flag into the config struct. We
	// default to using our development DSN if no flag is provided.
	flag.StringVar(&cfg.db.dsn, "db-dsn", dbHost, "PostgreSQL DSN")

	// Read the connection pool settings from command-line flags into the config struct.
	// Notice the default values that we're using?
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	// Create command line flags to read the setting values into the config struct.
	// Notice that we use true as the default for the 'enabled' setting?
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// Read the SMTP server configuration settings into the config struct, using the
	// Mailtrap settings as the default values.
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", smtpUser, "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", smtpPass, "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Oyna <no-reply@oyna.com>", "SMTP sender")

	flag.Parse()

	logger := pkg.New(os.Stdout, pkg.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	rabbitConn, err := amqp.Dial(rabbitMQHost)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer rabbitConn.Close()

	logger.PrintInfo("RabbitMQ connection established", nil)

	app := &application{
		config: cfg,
		logger: pkg.Logger{},
		//models: data.NewModels(db),
		//mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	// Call app.serve() to start the server.
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	//userRepo := repository.NewUserRepository()
	//toyRepo := repository.NewToyRepository()
	//subscriptionRepo := postgres.NewSubscriptionRepository()

	//userService := service.NewUserService(userRepo)
	//toyService := service.NewToyService(toyRepo)
	//subscriptionService := service.NewSubscriptionService(subscriptionRepo, userService)

	//handler.NewUserHandler(r, userService)
	//handler.NewToyHandler(r, toyService)
	//handler.NewSubscriptionHandler(r, subscriptionService)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	// Use the time.ParseDuration() function to convert the idle timeout duration string
	// to a time.Duration type.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil
}
