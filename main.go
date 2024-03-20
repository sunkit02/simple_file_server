package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Default values
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current working directory. Aborting...")
	}
	port := "3000"

	// Get cli arguments
	args := os.Args

	if len(args) == 2 && (args[1] == "-h" || args[1] == "--help") {
		showHelp()
		os.Exit(0)
	}

	if len(args) >= 2 {
		rootDir = args[1]
	}
	if len(args) >= 3 {
		port = args[2]
	}

	allowedOrigins := ""
	if len(args) == 4 {
		allowedOrigins = args[3]
	}

	fsHandler := http.FileServer(http.Dir(rootDir))

	if allowedOrigins != "" {
		fsHandler = cors(fsHandler, allowedOrigins)
	}

	http.Handle("/", fsHandler)

	log.Printf("Root directory: %s", rootDir)
	log.Printf("Listening on port: %s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	if err != nil {
		log.Fatal(err)
	}
}

func cors(fs http.Handler, allowedOrigins string) http.HandlerFunc {
	log.Printf("Allowing traffic from origin(s): %s", allowedOrigins)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)

		fs.ServeHTTP(w, r)
	}
}

func showHelp() {
	fmt.Println("Usage: simple_file_server [rootDir=PATH] [port] [ALLOWED_ORIGIN_1[,ALLOWED_ORIGIN_2]]")
	fmt.Println("Defaults: [rootDir=CWD] [port=3000]")
}
