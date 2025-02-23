package main

import (
	"fmt"
	"os"
	"time"
)

type State struct {
	ConnectedClients int
	StartedAt        time.Time
	NoRequests       int
	NoErrors         int
	NoUpdates        int
	ServeFsDir       string
	ServePath        string
	Status           string
	Urls             []string
}

var ServerState = State{
	ConnectedClients: 0,
	StartedAt:        time.Now(),
	NoRequests:       0,
	NoErrors:         0,
	ServeFsDir:       "",
	ServePath:        "",
	Urls:             []string{},
}

var stateUpdateCh = make(chan struct{}, 1)

// notifyStateUpdate is called when ServerState is modified.
func unthrottledNotifyServerStateUpdate() {
	// Non-blocking send: if the channel already has a signal, we don't block.
	select {
	case stateUpdateCh <- struct{}{}:
	default:
	}
}

var notifyServerStateUpdate = Throttle(unthrottledNotifyServerStateUpdate, 100*time.Millisecond)

func monitorServerState() {
	renders := 0
	for {
		<-stateUpdateCh

		// Drain any additional pending updates.
		for len(stateUpdateCh) > 0 {
			<-stateUpdateCh
		}

		if renders > 0 {
			fmt.Fprintf(os.Stderr, "\033[6A")
		}
		renders += 1
		var url string
		if len(ServerState.Urls) > 0 {
			url = ServerState.Urls[0]
		} else {
			url = "<empty>"
		}
		fmt.Fprintf(os.Stderr, "\r\033[K%s%s%s%s serving %s%s%s from %s%s%s on\n",
			Clr.Bold, Clr.Green, "dotdev", Clr.Reset,
			Clr.Bold, ServerState.ServePath, Clr.Reset,
			Clr.Bold, ServerState.ServeFsDir, Clr.Reset,
		)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "\r\033[K    %s%s%s\n", Clr.Bold, url, Clr.Reset)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "\r\033[KRequests: %d, Updates: %d, Errors: %d, WS clients: %d\n", ServerState.NoRequests, ServerState.NoUpdates, ServerState.NoErrors, ServerState.ConnectedClients)
		fmt.Fprintf(os.Stderr, "\r\033[K%sRuntime: %s, Renders: %d%s\n", Clr.Neutral, time.Since(ServerState.StartedAt).Round(time.Second), renders, Clr.Reset)
	}
}
