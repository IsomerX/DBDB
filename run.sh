#!/bin/bash
trap 'kill %1; exit' SIGINT      # Set up trap for Ctrl+C
(cd backend && go run main.go) & # Start backend and run in background (&)
(cd frontend && npm run dev)     # Start frontend in foreground
