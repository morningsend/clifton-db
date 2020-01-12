package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("cliftondb", "CliftonDB command line")

	serverCmd         = app.Command("server", "CliftonDB server")
	clusterManagerCmd = app.Command("cluster", "CliftonDB cluster")
	cliCmd            = app.Command("cli", "CliftonDB CLI")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case serverCmd.FullCommand():
		runServer()
	case clusterManagerCmd.FullCommand():
	case cliCmd.FullCommand():
	}
}
