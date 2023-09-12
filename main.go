package main

import (
	"github.com/SpectraLogic/ssc_go_client/client"
	"log"
	"os"
)

func main() {
	// Parse the arguments.
	args, err := client.ParseArgs()
	if err != nil {
		if err.Error() == "Must specify a command." {
			client.ListCommands(args)
			return
		}
	}

	if args.Command == "list_commands" {
		client.ListCommands(args)
		return
	}

	if len(args.LogFile) > 0 {
		wOut := os.Stdout
		if len(args.LogFile) > 0 {
			f, err := os.Create(args.LogFile)
			if err != nil {
				log.Printf("Could not create log file %s\n%v\n", args.LogFile, err)
			}
			defer f.Close()
			wOut = f
		}
		log.SetOutput(wOut)
	}

	// Create client if command requires it
	var storCycle *client.SscClient
	tokenRequired, err := client.CommandRequiresClientToken(args)
	if err != nil {
		log.Printf("unknown command %v\n", err)
		return
	}
	if tokenRequired {
		storCycle, err = client.CreateClient(args)
		if err != nil {
			log.Printf("could not create StorCycle client %v\n", err)
			return
		}
	}

	// Run the command
	err = client.RunCommand(storCycle, args)
	if err != nil {
		log.Printf("error running command %s, %v\n", args.Command, err)
		return
	}

}
