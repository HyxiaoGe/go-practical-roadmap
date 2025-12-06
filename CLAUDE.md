# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This repository is a Go learning roadmap designed for Java backend developers who want to learn Go through hands-on projects. It follows a practical approach with 4 progressive stages:

1. Web API template using Gin + Gorm
2. Concurrency worker system using goroutines and worker pools
3. gRPC microservices with service discovery using Etcd
4. Cross-platform CLI toolkit using Cobra

## Repository Structure

The repository is organized into 4 main project directories:
```
├── 01-web-api-template/        # Web API project (Gin + Gorm)
├── 02-concurrency-worker/      # Concurrent task system (goroutine + worker pool)
├── 03-grpc-microservices/      # Microservice project (gRPC + Etcd)
└── 04-cli-toolkit/             # CLI tool (Cobra)
```

Each directory represents a learning stage with its own README documentation.

## Language Policy

All communications, code comments, documentation, and any written content in this repository should be in Chinese (中文). This includes:
- All responses and explanations
- Code comments and documentation strings
- README files and other documentation
- Commit messages

## Common Development Tasks

When working with Go projects in this repository, you'll typically be dealing with:

1. Building Go applications with `go build`
2. Running applications with `go run`
3. Managing dependencies with Go modules (`go mod tidy`)
4. Testing with `go test`
5. Working with Docker for containerization

## Key Technologies Used

- Web Framework: Gin
- ORM: Gorm
- Logging: Zap
- Configuration: Viper
- CLI Framework: Cobra
- RPC: gRPC with Protobuf
- Service Discovery: Etcd
- Concurrency: goroutines, channels, worker pools

## Environment Requirements

- Go 1.21+
- Docker & Docker Compose
- Protobuf Compiler (`protoc`)
- Git

Follow the progression through the four stages in numerical order to build up Go engineering skills systematically.