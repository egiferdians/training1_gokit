package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net"

	"training1_gokit/modules/account/delivery/protobuf/account_grpc"

	_ "github.com/lib/pq"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/log/level"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"training1_gokit/modules/account"
	"training1_gokit/modules/account/delivery"
	grpcdelivery "training1_gokit/modules/account/delivery/grpc"
	httpdelivery "training1_gokit/modules/account/delivery/http"
)

const (
	hostname      = "localhost"
	host_port     = 5432
	username      = "postgres"
	password      = "admin123"
	database_name = "training1_gokit"
)

func restMode(
	ctx context.Context,
	logger log.Logger,
	endpoints delivery.Endpoints,
) {
	/*
		http addres flag and geting port from env
		Perpare HTTP Handler
	*/
	port := os.Getenv("HTTP_PORT")
	httpAddr := flag.String("http", port, "http listen address")
	handler := httpdelivery.NewHTTPServe(ctx, endpoints, logger)
	flag.Parse()

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: handler,
		}
		errs <- server.ListenAndServe()
	}()
	level.Error(logger).Log("exit", <-errs)
}

func grpcMode(
	_ context.Context,
	logger log.Logger,
	endpoints delivery.Endpoints,
) {
	port := ":5051"
	var (
		accountService  = grpcdelivery.NewGRPCServer(endpoints, logger)
		grpcListener, _ = net.Listen("tcp", port)
		grpcServer      = grpc.NewServer()
		g               group.Group
	)

	g.Add(func() error {
		logger.Log("transport", "gRPC", "addr", port)
		account_grpc.RegisterAccountServiceServer(grpcServer, accountService)
		return grpcServer.Serve(grpcListener)
	}, func(error) {
		grpcListener.Close()
	})

	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
	level.Error(logger).Log("exit", g.Run())
}

func main() {
	dbsource := fmt.Sprintf("port=%d host=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host_port, hostname, username, password, database_name)

	// var httpAddr = flag.String("http", ":8080", "http listen address")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "account",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var db *sql.DB
	{
		var err error

		db, err = sql.Open("postgres", dbsource)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}

	}

	flag.Parse()
	ctx := context.Background()
	var srv account.Service
	{
		repository := account.NewRepo(db, logger)

		srv = account.NewService(repository, logger)
	}

	// errs := make(chan error)

	// go func() {
	// 	c := make(chan os.Signal, 1)
	// 	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	// 	errs <- fmt.Errorf("%s", <-c)
	// }()

	endpoints := delivery.MakeEndpoints(srv)
	grpcMode(ctx, logger, endpoints)

	// go func() {
	// 	fmt.Println("listening on port", *httpAddr)
	// 	handler := account.NewHTTPServer(ctx, endpoints)
	// 	errs <- http.ListenAndServe(*httpAddr, handler)
	// }()

	// level.Error(logger).Log("exit", <-errs)
}
