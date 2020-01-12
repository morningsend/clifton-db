package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("cliftondb", "CliftonDB command line")

	server         = app.Command("server", "CliftonDB server")
	clusterManager = app.Command("cluster", "CliftonDB cluster")
	cli            = app.Command("cli", "CliftonDB CLI")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case server.FullCommand():
		runServer()
	}
}
