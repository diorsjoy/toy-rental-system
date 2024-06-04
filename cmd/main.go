package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"toy-rental-system/internal/api/handler"
	"toy-rental-system/internal/config"
	"toy-rental-system/internal/repository/postgres"
	"toy-rental-system/internal/service"
	pkg "toy-rental-system/pkg/jsonlog"
)

type configuration struct {

	port int
	env  string
	db   struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
	}
	limiter struct {
	rps     float64
	burst   int
	enabled bool
	}
	}



type application struct {
	config configuration
	subscriptionHandler *handler.SubscriptionHandler
	logger pkg.Logger
	wg sync.WaitGroup
}

func main() {
	var cfg configuration

	env, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configuration:", err)
	}

	// stripeKey := env.StripeSecret
	dbHost := env.DBSource
	smtpUser := env.SMTPUsername
	smtpPass := env.SMTPPassword
	rabbitMQHost := env.RabbitMQSource

	flag.IntVar(&cfg.port, "port", serverAddress, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", dbHost, "PostgreSQL DSN")

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", dbHost, "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

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

	subscriptionRepo := postgres.NewSubscriptionRepository(db)
	subscriptionService := service.NewSubscriptionService(env, subscriptionRepo)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	app := &application{
		config:              cfg,
		logger:              pkg.Logger{},
		subscriptionHandler: subscriptionHandler,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// Again, we use the PrintInfo() method to write a "starting server" message at the
	// INFO level. But this time we pass a map containing additional properties (the
	// operating environment and server address) as the final parameter.
	logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})
	err = srv.ListenAndServe()
	// Use the PrintFatal() method to log the error and exit.
	logger.PrintFatal(err, nil)

	//// Call app.serve() to start the server.
	//err = app.serve()
	//if err != nil {
	//	logger.PrintFatal(err, nil)
	//}

}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {

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
