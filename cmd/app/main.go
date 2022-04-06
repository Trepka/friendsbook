package main

import (
	"flag"
	"friendsbook/internal/platform/server"
	"strconv"
)

var port *int

func init() {
	port = flag.Int("port", 3000, "port number")
}

func main() {
	flag.Parse()
	server.StartApp(strconv.Itoa(*port))
}
