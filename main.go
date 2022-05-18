package main

import (
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/client"
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

	// Create client
	storCycle, err := client.CreateClient(args)
	if err != nil {
		fmt.Printf("could not create client %v\n", err)
		return
	}


	// Run the command
	err = client.RunCommand(storCycle, args)
	if err != nil {
		fmt.Printf("error running command %s, %v\n", args.Command, err)
		return
	}

}