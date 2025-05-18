package main

import (
	"context"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/plsort"
)

var version string

func init() {
	log.SetFlags(0)
}

func parseArgs() *plsort.Options {
	var cli struct {
		plsort.Options
		Version kong.VersionFlag
	}

	parser := kong.Must(&cli, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."
	_, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	return &cli.Options
}

func main() {
	options := parseArgs()
	ctx := context.Background()
	client, err := plsort.NewClient(ctx, options.Creds)

	if err != nil {
		log.Fatalf("error: %s", err)
	}

	err = client.Sort(ctx, &options.SortOptions)

	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
