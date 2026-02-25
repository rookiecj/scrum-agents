package main

// Version is set at build time via:
//
//	go build -ldflags "-X main.Version=$(cat VERSION)" ./cmd/server
var Version = "dev"
