package main

import (
	"errors"
	"os"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/theonlyjohnny/go-logger/logger"
	"github.com/theonlyjohnny/phoenix/common"
	"github.com/theonlyjohnny/phoenix/logic"
)

var log *logger.Logger

func run(args common.RunArgs) {
	log.Debugf("%#v", args)
	controller := logic.NewController(args, log)
	err := controller.Act()
	if err != nil {
		panic(err)
	}
}

func main() {
	parser := argparse.NewParser("phoenix", "Creates Phoenix Infrastructure Systems")
	action := parser.String("", "action", &argparse.Options{
		Help:    "Action to commit: [create|destroy]",
		Default: "create",
		Validate: func(args []string) error {
			if args[0] != "create" && args[0] != "destroy" {
				return errors.New("action must be create or destroy")
			}
			return nil
		},
	})
	region := parser.String("r", "region", &argparse.Options{
		Required: true,
		Help:     "AWS Region to operate in",
	})

	passedVPCName := parser.String("", "name", &argparse.Options{
		Required: false,
		Help:     "AWS VPC Name to operate on",
	})

	loggerConfig := logger.Config{
		AppName:    "phoenix",
		LogLevel:   "debug",
		LogConsole: true,
		LogSyslog:  nil,
	}

	var err error
	log, err = logger.CreateLogger(loggerConfig)
	if err != nil {
		panic(err)
	}

	if err := parser.Parse(os.Args); err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		for _, line := range strings.Split(parser.Usage(err), "\n") {
			log.Error(line)
		}
		// fmt.Print(parser.Usage(err))
	} else {
		args := common.RunArgs{
			*action,
			*region,
			*passedVPCName,
		}
		run(args)
	}
}
