package main

import (
	http "karboosx/http"
)

func main() {
    server := http.CreateServer(8080)
	server.Listen()
}
