package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"
)

func DevServer(
	srvFs fs.FS,
	htmlFile string,
) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler(srvFs, htmlFile))
	mux.HandleFunc("/ws", wsHandler)
	return mux
}

func StartFileWatcher(filePath string) {
	watchFileInotify(filePath, Throttle(broadcastReload, 100*time.Millisecond))
}

func StartDevServer(
	host string,
	port int,
	srvFs fs.FS,
	htmlFile string,
) {
	go printServerState()
	ServerState.StartedAt = time.Now()
	ServerState.ServePath = htmlFile
	ServerState.Urls = []string{fmt.Sprintf("http://%s:%d", host, port)}
	notifyStateUpdate()
	go StartFileWatcher(htmlFile)
	server := DevServer(srvFs, htmlFile)
	addr := fmt.Sprintf("%s:%d", host, port)
	// log.Printf("Serving %s on http://%s", htmlFile, addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Printf("Unrecoverable error: %v", err)
		log.Fatal(err)
	}
}

// broadcastReload sends a "reload" message to all connected WebSocket clients.
func broadcastReload() {
	ServerState.NoUpdates += 1
	notifyStateUpdate()
	wsMutex.Lock()
	defer wsMutex.Unlock()
	// log.Printf("Broadcasting to %d clients\n", len(wsClients))
	for i := 0; i < len(wsClients); {
		conn := wsClients[i]
		err := sendPayload(conn, []byte("reload"))
		if err != nil {
			log.Printf("Error sending reload message: %v\n", err)
			// Remove the connection if there's an error.
			conn.Close()
			wsClients = append(wsClients[:i], wsClients[i+1:]...)
		} else {
			i++
		}
	}
}
