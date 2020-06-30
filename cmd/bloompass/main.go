package main

import (
	"bloompass/pkg/bloompass"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const SEED uint32 = 0xd5c48bfc // static seed for demo purposes, taken from murmur3 tests

var (
	directory = flag.String("dir", "", "Directory to scan for password lists (*.txt)")
	host      = flag.String("host", "localhost", "Host to use for API server")
	port      = flag.String("port", ":3030", "Port to use for API server")
)

func isExistingDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsDir()
}

func main() {
	flag.Parse()
	if len(*directory) == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if !isExistingDir(*directory) {
		log.Fatalf("could not access specified directory")
	}

	filter := bloompass.NewBloom(SEED)
	if err := bloompass.ParseFiles(*directory, func(s string) {
		filter.Add(s)
	}, func(m string) { log.Printf("Parsing file %s", m) }); err != nil {
		log.Fatalf("Error parsing files: %v", err)
	}

	log.Println("Starting API server...")
	server := bloompass.NewApiServer(*host, *port, filter)

	cancelChan := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := server.Stop(ctx); err != nil {
			log.Fatal(err)
		}
		close(cancelChan)
	}()

	err := server.Start()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-cancelChan
}
