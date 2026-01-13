package main

import (
	"log"
	"os"
	"translateserver/config"
	"translateserver/db"
	"translateserver/worker"

	SherryServer "github.com/asccclass/sherryserver"
)

func main() {
	// Initialize Config
	config.Init()

	// Initialize Database
	if err := db.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	// Initialize Worker
	go worker.Start()

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	documentRoot := os.Getenv("DocumentRoot")
	if documentRoot == "" {
		documentRoot = "www"
	}
	templateRoot := os.Getenv("TemplateRoot")
	if templateRoot == "" {
		templateRoot = "www/html"
	}

	// Initialize Server
	// NewServer(listenAddr, documentRoot, templatePath)
	server, err := SherryServer.NewServer(":"+port, documentRoot, templateRoot)
	if err != nil {
		panic(err)
	}

	router := NewRouter(server, documentRoot)
	server.Server.Handler = server.CheckCROS(router)
	server.Start()

}
