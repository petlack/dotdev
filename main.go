package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed version.txt
var versionFile embed.FS
var Version string

const (
	DEFAULT_PORT = 4774
	DEFAULT_HOST = "127.0.0.1"
)

func init() {
	versionData, err := fs.ReadFile(versionFile, "version.txt")
	if err != nil {
		log.Fatal(err)
	}
	Version = strings.TrimSpace(string(versionData))
}

func main() {
	args := os.Args[1:]
	action := ""
	serveFile := ""
	if len(args) < 1 {
		log.Printf("Please provide a file to watch\n")
		log.Printf("Usage: dotdev <file-to-watch> [--host <host>] [--port <port>]\n")
		action = "help"
	} else if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		action = "help"
	} else if args[0] == "version" || args[0] == "--version" || args[0] == "-v" {
		action = "version"
	} else {
		serveFile = os.Args[1]
		action = "serve"
		args = args[1:]
	}

	switch action {

	case "serve":
		defaultPort, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			defaultPort = DEFAULT_PORT
		}
		defaultHost := os.Getenv("HOST")
		if defaultHost == "" {
			defaultHost = DEFAULT_HOST
		}
		serveFileParentDir := filepath.Dir(serveFile)
		if _, err := os.Stat(serveFile); err != nil {
			log.Printf("Serve file not found: %s\n", serveFile)
			log.Fatal(err)
		}
		go monitorServerState()
		ServerState.ServeFsDir = serveFileParentDir
		notifyServerStateUpdate()
		configFlagSet := flag.NewFlagSet("dotdev", flag.ContinueOnError)
		host := configFlagSet.String("host", defaultHost, "Host of the dev server")
		port := configFlagSet.Int("port", defaultPort, "Port of the dev server")
		configFlagSet.Parse(args)
		StartDevServer(*host, *port, serveFile)
		os.Exit(0)

	case "help":
		printHelp()
		os.Exit(0)

	case "version":
		fmt.Printf("%s", Version)
		os.Exit(0)

	default:
		printHelp()
		log.Printf("Unknown action: %s\n", action)
		os.Exit(1)

	}
}

func printHelp() {
	fmt.Fprintf(os.Stderr, "%sdotdev %sv%s%s\n", Clr.Bold, Clr.Neutral, Version, Clr.Reset)
	fmt.Fprintf(os.Stderr, "    Simple HTTP server with live reload\n\n")
	fmt.Fprintf(os.Stderr, "%sUSAGE:%s\n", Clr.Bold, Clr.Reset)
	fmt.Fprintf(os.Stderr, "    dotdev <file> [options]\n")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "    %sOPTIONS%s:\n", Clr.Underline, Clr.Reset)
	fmt.Fprintf(os.Stderr, "    %s--port <PORT>%s\n", Clr.Bold, Clr.Reset)
	fmt.Fprintf(os.Stderr, "        Port of the dev server\n")
	fmt.Fprintf(os.Stderr, "    %s--host <HOST>%s\n", Clr.Bold, Clr.Reset)
	fmt.Fprintf(os.Stderr, "        Host of the dev server\n")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "%sEXAMPLE%s:\n", Clr.Bold, Clr.Reset)
	fmt.Fprintf(os.Stderr, "echo \"<html><body>Hello World</body></html>\" > ./index.html\n")
	fmt.Fprintf(os.Stderr, "dotdev ./index.html --host localhost --port 4774\n")
	fmt.Fprintln(os.Stderr)
}

type Colors struct {
	Bold, Green, Neutral, Red, Reset, ResetUnderline, Underline, Yellow string
}

var Clr = Colors{
	Bold:           "\033[1;39m",
	Green:          "\033[1;92m",
	Neutral:        "\033[0;97m",
	Red:            "\033[1;91m",
	Reset:          "\033[0;39m",
	ResetUnderline: "\033[24m",
	Underline:      "\033[4m",
	Yellow:         "\033[1;93m",
}
