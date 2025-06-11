package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func DevServer(
	htmlFile string,
) http.Handler {
	mux := http.NewServeMux()
	idxHandler := indexHandler(htmlFile)
	baseDir := filepath.Dir(htmlFile)
	fileServer := http.FileServer(http.Dir(baseDir))

	mux.HandleFunc("/ws", wsHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			idxHandler(w, r)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
	return mux
}

func StartFileWatcher(filePath string) {
	lower := strings.ToLower(filePath)
	if strings.HasSuffix(lower, ".html") || strings.HasSuffix(lower, ".htm") {
		if assets, err := GetIncludedAssets(filePath); err == nil {
			for _, a := range assets {
				go StartFileWatcher(a)
			}
		}
	}
	if runtime.GOOS == "linux" {
		watchFileInotify(filePath, Throttle(broadcastReload, 100*time.Millisecond))
	}
	watchFilePoll(filePath, broadcastReload)
}

func StartDevServer(
	host string,
	port int,
	htmlFile string,
) {
	ServerState.StartedAt = time.Now()
	ServerState.ServePath = htmlFile
	ServerState.Urls = []string{fmt.Sprintf("http://%s:%d", host, port)}
	notifyServerStateUpdate()
	go StartFileWatcher(htmlFile)
	server := DevServer(htmlFile)
	addr := fmt.Sprintf("%s:%d", host, port)
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Printf("Unrecoverable error: %v", err)
		log.Fatal(err)
	}
}

// broadcastReload sends a "reload" message to all connected WebSocket clients.
func broadcastReload() {
	ServerState.NoUpdates += 1
	notifyServerStateUpdate()
	wsMutex.Lock()
	defer wsMutex.Unlock()
	for i := 0; i < len(wsClients); {
		conn := wsClients[i]
		err := sendPayload(conn, []byte("reload"))
		if err != nil {
			log.Printf("Error sending reload message: %v\n", err)
			conn.Close()
			wsClients = append(wsClients[:i], wsClients[i+1:]...)
		} else {
			i++
		}
	}
}
