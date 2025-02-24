package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed assets/*
var assetsFs embed.FS

func indexHandler(
	htmlFile string,
) http.HandlerFunc {
	liveReloadScriptBytes, err := fs.ReadFile(assetsFs, "assets/live-reload.js")
	if err != nil {
		fmt.Printf("%sError reading live-reload.js. Live reload will not work.%s\n", Clr.Red, Clr.Reset)
	}
	liveReloadScript := strings.ReplaceAll(string(liveReloadScriptBytes), "{{dotdev::version}}", Version)
	errorResponseBytes, err := fs.ReadFile(assetsFs, "assets/error.html")
	if err != nil {
		fmt.Printf("%sError reading error.html. Error page will not work.%s\n", Clr.Red, Clr.Reset)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ServerState.NoRequests += 1
		notifyServerStateUpdate()
		if r.URL.Path != "/" {
			handleError(
				w, errorResponseBytes, http.StatusNotFound,
				"Not Found",
				fmt.Sprintf("The page %s does not exist.", r.URL.Path),
			)
			return
		}

		content, err := os.ReadFile(htmlFile)
		if err != nil {
			log.Printf("Error reading file: %v\n", err)
			handleError(
				w, errorResponseBytes, http.StatusInternalServerError,
				"Unexpected error",
				fmt.Sprintf("Error reading file: %v", err),
			)
			return
		}

		htmlContent := string(content)
		snippet := fmt.Sprintf("<script type=\"text/javascript\">\n%s\n</script>", liveReloadScript)
		if idx := strings.LastIndex(htmlContent, "</body>"); idx != -1 {
			htmlContent = htmlContent[:idx] + "\n" + snippet + "\n" + htmlContent[idx:]
		} else {
			htmlContent += snippet
		}
		w.Write([]byte(htmlContent))
	}
}

func handleError(w http.ResponseWriter, errorResponseBytes []byte, statusCode int, message string, description string) {
	ServerState.NoErrors += 1
	notifyServerStateUpdate()
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	materializedErorrData := strings.ReplaceAll(string(errorResponseBytes), "{{dotdev::error.statusCode}}", fmt.Sprintf("%d", statusCode))
	materializedErorrData = strings.ReplaceAll(materializedErorrData, "{{dotdev::error.message}}", message)
	materializedErorrData = strings.ReplaceAll(materializedErorrData, "{{dotdev::error.description}}", description)
	materializedErorrData = strings.ReplaceAll(materializedErorrData, "{{dotdev::version}}", Version)
	w.Write([]byte(materializedErorrData))
}
