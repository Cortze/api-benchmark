package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/cortze/api-benchmark/cmd"
	"github.com/cortze/api-benchmark/pkg/utils"
)

var (
	Version = "v0.0.1"
	CliName = "Api-Benchmark"
	log     = logrus.WithField(
		"cli", "CliName",
	)
)

func main() {
	fmt.Println(CliName, Version, "\n")

	//ctx, cancel := context.WithCancel(context.Background())

	// Set the general log configurations for the entire tool
	logrus.SetFormatter(utils.ParseLogFormatter("text"))
	logrus.SetOutput(utils.ParseLogOutput("terminal"))
	logrus.SetLevel(utils.ParseLogLevel("debug"))

	app := &cli.App{
		Name:      CliName,
		Usage:     "Tinny client that performs an API Benchmark on the given endpoint.",
		UsageText: "api-benchmark [commands] [arguments...]",
		Authors: []*cli.Author{
			{
				Name:  "Cortze",
				Email: "cortze@protonmail.com",
			},
		},
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			cmd.RunCommand,
		},
	}

	// generate the crawler
	if err := app.RunContext(context.Background(), os.Args); err != nil {
		log.Errorf("error: %v\n", err)
		os.Exit(1)
	}
	/*
		// only leave the app up running if the command was empty or help
		if len(os.Args) <= 1 || helpInArgs(os.Args) {
			os.Exit(0)
		} else {
			// check the shutdown signal
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)

			// keep the app running until syscall.SIGTERM
			sig := <-sigs
			log.Printf("Received %s signal - Stopping...\n", sig.String())
			signal.Stop(sigs)
			cancel()
		}
	*/
}

func helpInArgs(args []string) bool {
	help := false
	for _, b := range args {
		switch b {
		case "--help":
			help = true
			return help
		case "-h":
			help = true
			return help
		case "h":
			help = true
			return help
		case "help":
			help = true
			return help
		}
	}
	return help
}
