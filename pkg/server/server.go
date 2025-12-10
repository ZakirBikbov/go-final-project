package server

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"final-project/pkg/api"
)

func Start(webDir string, defaultPort int) error {
	port := resolvePort(defaultPort)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	api.Init(mux)

	addr := ":" + strconv.Itoa(port)
	log.Printf("serving %s on http://localhost%s/", webDir, addr)
	return http.ListenAndServe(addr, mux)
}

func resolvePort(defaultPort int) int {
	envPort := os.Getenv("TODO_PORT")
	if envPort == "" {
		return defaultPort
	}

	if parsed, err := strconv.Atoi(envPort); err == nil && parsed > 0 && parsed < 65536 {
		return parsed
	}

	log.Printf("invalid TODO_PORT %q, using default %d", envPort, defaultPort)
	return defaultPort
}
