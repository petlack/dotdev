function fetchAndReload() {
    console.log("[dotdev] Fetching updated content", location.href, "...");
    fetch(location.href)
        .then(response => response.text())
        .then(html => {
            var parser = new DOMParser();
            var doc = parser.parseFromString(html, 'text/html');
            document.head.innerHTML = doc.head.innerHTML;
            document.body.innerHTML = doc.body.innerHTML;
        })
        .catch(err => {
            console.error("[dotdev] Hot update failed:", err);
        });
}

function connectWs() {
    var ws = new WebSocket("ws://" + location.host + "/ws");

    ws.onopen = () => {
        console.log("[dotdev] WebSocket connected");
        fetchAndReload();
    };

    ws.onmessage = msg => {
        if (msg.data === "refresh") {
            location.reload();
        }
        if (msg.data === "reload") {
            console.log("[dotdev] Hot updating app content ...");
            fetchAndReload();
        }
    };

    ws.onclose = () => {
        console.log("[dotdev] WebSocket disconnected, reconnecting in 1s...");
        setTimeout(connectWs, 1000);
    };

    ws.onerror = err => {
        console.error("[dotdev] WebSocket encountered error: ", err);
        ws.close();
    };
}

function main() {
    console.log("[dotdev] Version {{dotdev::version}}");
    connectWs();
}

main()
