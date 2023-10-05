package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
)

/*
	Uptime Monitor ToolBox
	Author: tobychui

	This is a generic tool for hosting basic information required for my server daily operations

*/

//go:embed web/*
var content embed.FS

// Default configs
var configFilepath = "config.json"
var logFilepath = "uptime.log"
var exampleTarget = Target{
	ID:       "example",
	Name:     "Example",
	URL:      "example.com",
	Protocol: "https",
}
var usingConfig *Config
var onlineStatusLog = map[string][]*Record{}

// TOTP configs
var totpConfig *TotpConfig

// Set this to true to use file system based UI files
var uiDebugMode = false

// Flags
var listeningPort = flag.String("p", ":8089", "Listening endpoint for http server")

func main() {
	flag.Parse()

	//Start the uptime monitor
	err := UptimeMonitorInit()
	if err != nil {
		panic(err)
	}

	//Start TOTP resolver
	err = totpInit()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/utm/update", HandleUptimeLogRead)
	http.HandleFunc("/totp/update", HandleTOTPUpdate)

	var webfs, _ = fs.Sub(content, "web")
	if uiDebugMode {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.FileServer(http.Dir("./web")).ServeHTTP(w, r)
		})
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.FileServer(http.FS(webfs)).ServeHTTP(w, r)
		})
	}

	//Start the web server interface
	log.Println("Listening on :8089")
	http.ListenAndServe(*listeningPort, nil)
}
