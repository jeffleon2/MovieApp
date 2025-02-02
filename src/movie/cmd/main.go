package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
	"movieexample.com/src/movie/internal/controller/movie"
	metadataGateway "movieexample.com/src/movie/internal/gateway/grpc/metadata"
	ratingGateway "movieexample.com/src/movie/internal/gateway/grpc/rating"
	httpHandler "movieexample.com/src/movie/internal/handler/http"
	"movieexample.com/src/pkg/discovery"
	"movieexample.com/src/pkg/discovery/consul"
)

const serviceName = "movie"

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
	f, err := os.Open("base.yaml")
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
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("0.0.0.0:%d", port)); err != nil {
		log.Println("PANIC REGISTRY unavilable to conect")
		panic(err)
	}
	log.Printf("Registring service in consul port %d", consulPort)
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state" + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)
	log.Printf("Starting the movie service in port %d", port)
	metadataGateway := metadataGateway.New(registry)
	ratingGateway := ratingGateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httpHandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
