package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestLiveReload performs an end-to-end test of the live-reload server.
func TestLiveReload(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "live-reload-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "index.html")
	originalHTML := `<html><head><title>Test</title></head><body>Hello World</body></html>`
	if err := os.WriteFile(filePath, []byte(originalHTML), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	srvFS := os.DirFS(tmpDir)
	handler := DevServer(srvFS, "index.html")
	go StartFileWatcher(filePath)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	htmlContent := getHtmlContent(t, ts.URL)
	if !strings.Contains(htmlContent, "WebSocket") {
		t.Fatalf("Expected injected javascript snippet in html, got: %s", htmlContent)
	}
	if !strings.Contains(htmlContent, Version) {
		t.Fatalf("Expected version to be injected, got: %s", htmlContent)
	}

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("Failed to parse test server URL: %v", err)
	}

	// Connect to the /ws endpoint using a raw TCP connection and perform a minimal WebSocket handshake.
	wsConn := dialWebSocket(t, u.Host)
	defer wsConn.Close()

	// Simulate a file change.
	updatedHTML := `<html><head><title>Test</title></head><body>Hello Reload</body></html>`
	if err := os.WriteFile(filePath, []byte(updatedHTML), 0644); err != nil {
		t.Fatalf("Failed to update file: %v", err)
	}
	newTime := time.Now().Add(2 * time.Second)
	if err := os.Chtimes(filePath, newTime, newTime); err != nil {
		t.Fatalf("Failed to change file times: %v", err)
	}

	// Wait for the "reload" message from the websocket.
	done := make(chan string, 1)
	go func() {
		msg := readWebSocketMessage(t, wsConn)
		done <- msg
	}()

	select {
	case msg := <-done:
		if msg != "reload" {
			t.Fatalf("Expected 'reload' message, got: %q", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timed out waiting for reload message")
	}

	newHtmlContent := getHtmlContent(t, ts.URL)
	if !strings.Contains(newHtmlContent, "Hello Reload") {
		t.Fatalf("Expected server to serve update version of the HTML file")
	}
}

func getHtmlContent(t *testing.T, url string) string {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	htmlContent := string(body)
	return htmlContent
}

// dialWebSocket performs a minimal WebSocket handshake on the given host for the /ws endpoint.
func dialWebSocket(t *testing.T, host string) net.Conn {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		t.Fatalf("Failed to connect to %s: %v", host, err)
	}

	// Prepare a minimal handshake request.
	key := "x3JJHMbDL1EzLkh9GBhXDw=="
	handshake := fmt.Sprintf("GET /ws HTTP/1.1\r\n"+
		"Host: %s\r\n"+
		"Upgrade: websocket\r\n"+
		"Connection: Upgrade\r\n"+
		"Sec-WebSocket-Key: %s\r\n"+
		"Sec-WebSocket-Version: 13\r\n\r\n", host, key)

	_, err = conn.Write([]byte(handshake))
	if err != nil {
		t.Fatalf("Failed to write handshake: %v", err)
	}

	// Read the response using http.ReadResponse.
	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		t.Fatalf("Failed to read handshake response: %v", err)
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("Expected 101 Switching Protocols, got %d", resp.StatusCode)
	}

	// Optionally, verify the Sec-WebSocket-Accept header.
	expectedAccept := computeAcceptKey(key)
	if resp.Header.Get("Sec-WebSocket-Accept") != expectedAccept {
		t.Fatalf("Invalid Sec-WebSocket-Accept, expected %s, got %s",
			expectedAccept, resp.Header.Get("Sec-WebSocket-Accept"))
	}

	return conn
}

// readWebSocketMessage reads a single text frame from the websocket connection.
// This simple implementation assumes that the message payload length is < 126 bytes.
func readWebSocketMessage(t *testing.T, conn net.Conn) string {
	// Read the 2-byte header.
	header := make([]byte, 2)
	if _, err := io.ReadFull(conn, header); err != nil {
		t.Fatalf("Failed to read ws header: %v", err)
	}

	// Check that FIN is set and the opcode is 1 (text frame).
	if header[0]&0x80 == 0 || header[0]&0x0F != 1 {
		t.Fatalf("Unexpected frame header: %v", header)
	}

	payloadLen := int(header[1] & 0x7F)
	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(conn, payload); err != nil {
		t.Fatalf("Failed to read ws payload: %v", err)
	}

	return string(payload)
}
