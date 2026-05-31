package main

import (
	"log"
	"os"

	"github.com/codemityio/gotpl/internal/app"
	"github.com/urfave/cli/v2"
)

func main() {
	application := app.New(
		app.WithValues(
			name,
			``,
			version,
			copyright,
			authorName,
			authorEmail,
			buildTime,
		),
	)

	application.Commands = []*cli.Command{}

	application.CommandNotFound = func(_ *cli.Context, cmd string) {
		log.Fatalf("error: unknown command %q", cmd)
	}

	if e := application.Run(os.Args); e != nil {
		log.Fatalf("error: %v", e)
	}
}
