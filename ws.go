package main

import (
	"crypto/sha1"
	"encoding/base64"
	"net"
	"net/http"
	"strings"
	"sync"
)

var (
	wsClients = make([]net.Conn, 0)
	wsMutex   sync.Mutex
)

// wsHandler handles the WebSocket handshake and upgrades the connection.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Header.Get("Upgrade")) != "websocket" {
		http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
		return
	}
	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		http.Error(w, "Missing Sec-WebSocket-Key", http.StatusBadRequest)
		return
	}
	acceptKey := computeAcceptKey(key)

	h := w.Header()
	h.Add("Upgrade", "websocket")
	h.Add("Connection", "Upgrade")
	h.Add("Sec-WebSocket-Accept", acceptKey)
	w.WriteHeader(http.StatusSwitchingProtocols)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		return
	}

	wsMutex.Lock()
	wsClients = append(wsClients, conn)
	wsMutex.Unlock()

	ServerState.ConnectedClients = len(wsClients)
	notifyStateUpdate()

	go func() {
		buf := make([]byte, 1024)
		for {
			_, err := conn.Read(buf)
			if err != nil {
				wsMutex.Lock()
				for i, c := range wsClients {
					if c == conn {
						wsClients = append(wsClients[:i], wsClients[i+1:]...)
						ServerState.ConnectedClients = len(wsClients)
						notifyStateUpdate()
						break
					}
				}
				wsMutex.Unlock()
				conn.Close()
				return
			}
		}
	}()
}

// computeAcceptKey computes the Sec-WebSocket-Accept key as specified in RFC 6455.
func computeAcceptKey(key string) string {
	const magicGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	h := sha1.New()
	h.Write([]byte(key + magicGUID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// sendReload sends a minimal WebSocket text frame with the "reload" command.
func sendPayload(conn net.Conn, payload []byte) error {
	frame := []byte{0x81} // 0x81 means FIN set and opcode 0x1 (text)

	payloadLen := len(payload)
	if payloadLen < 126 {
		frame = append(frame, byte(payloadLen))
	} else if payloadLen < 65536 {
		frame = append(frame, 126)
		frame = append(frame, byte(payloadLen>>8), byte(payloadLen&0xff))
	} else {
		frame = append(frame, 127)
		// Append 8 bytes for the length (big endian). Not needed for "reload" but shown for completeness.
		for i := 7; i >= 0; i-- {
			frame = append(frame, byte((uint64(payloadLen)>>(8*i))&0xff))
		}
	}
	frame = append(frame, payload...)
	_, err := conn.Write(frame)
	return err
}
