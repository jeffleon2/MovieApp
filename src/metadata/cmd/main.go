package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"

	"movieexample.com/src/gen"
	"movieexample.com/src/metadata/internal/controller/metadata"
	grpcHandler "movieexample.com/src/metadata/internal/handler/grpc"
	"movieexample.com/src/metadata/internal/repository/mysql"
	"movieexample.com/src/pkg/discovery"
	"movieexample.com/src/pkg/discovery/consul"
)

const serviceName = "metadata"

type Config struct {
	API    apiConfig    `yaml:"api"`
	Consul consulConfig `yaml:"consul"`
}

type apiConfig struct {
	Port int `yaml:"port"`
}

type consulConfig struct {
	Port int `yaml:"port"`
}

func main() {
	f, err := os.Open("../configs/base.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	port := cfg.API.Port
	consulPort := cfg.Consul.Port
	ctx := context.Background()
	registry, err := consul.NewRegistry(fmt.Sprintf("localhost:%d", consulPort))
	if err != nil {
		panic(err)
	}
	log.Printf("Registring service in consul port %d", consulPort)
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	log.Printf("Starting the movie metadata service in port %d", port)
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state" + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)
	repo, err := mysql.New()
	if err != nil {
		panic(err)
	}
	ctrl := metadata.New(repo)
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMetadataServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
	// http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	// if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
	// 	panic(err)
	// }

}
