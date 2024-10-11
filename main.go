package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/t3201v/seat-arrangement/controller"
	"github.com/t3201v/seat-arrangement/gen/cinema"
	"github.com/t3201v/seat-arrangement/internal/libs/assert"
	"github.com/t3201v/seat-arrangement/internal/libs/logger"
	"github.com/t3201v/seat-arrangement/repository"
	"github.com/t3201v/seat-arrangement/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	urlGRPC, urlHTTP   string
	grpcServerEndpoint *string
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	assert.NoError(err, "read config file failed")
	urlGRPC = ":" + viper.GetString("port_grpc")
	urlHTTP = ":" + viper.GetString("port_http")
	grpcServerEndpoint = flag.String("grpc-server-endpoint", urlGRPC, "gRPC server endpoint")
	assert.NotNilf(grpcServerEndpoint, "malformed gRPC server endpoint")
}

func main() {
	l := log.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(log.DebugLevel)
	l.SetFormatter(&logger.CallerFormatter{
		TextFormatter: &log.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		},
	})
	go startGRPC(l)
	go startHTTP(l)

	// loop forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func startGRPC(l *log.Logger) {
	l.Info("Starting gRPC server...")
	lis, err := net.Listen("tcp", urlGRPC)
	if err != nil {
		l.Fatal("failed to listen :", err)
	}

	// server
	s := grpc.NewServer()

	repo := repository.NewCinema(l)
	svc := service.NewCinema(l, repo)
	impl := controller.NewCinema(l, svc)
	cinema.RegisterCinemaServiceServer(s, impl)

	l.Infof("gRPC server started on %v", urlGRPC)
	err = s.Serve(lis)
	if err != nil {
		l.Fatal("failed to serv :", err)
	}
}

func startHTTP(l *log.Logger) {
	l.Info("Starting http server")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := cinema.RegisterCinemaServiceHandlerFromEndpoint(ctx, rmux, *grpcServerEndpoint, opts)
	if err != nil {
		l.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", rmux)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		l.Debug("health check http ok")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// swagger
	mux.HandleFunc("/swagger.json", serveSwagger)
	fs := http.FileServer(http.Dir("./www"))
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))

	l.Infof("http server started on %v", urlHTTP)
	http.ListenAndServe(urlHTTP, mux)
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./gen/cinema/cinema.swagger.json")
}
