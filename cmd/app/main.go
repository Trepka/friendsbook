package main

import (
	"flag"
	"friendsbook/internal/platform/server"
	"strconv"
)

var port *int

func init() {
	port = flag.Int("port", 8080, "port number")
}

func main() {
	flag.Parse()
	server.StartApp(strconv.Itoa(*port))
}
