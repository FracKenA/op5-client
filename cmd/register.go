package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	op5monitor "github.com/FracKenA/op5-client/op5"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registers host with OP5 Monitor.",
	Long:  ``,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("register called")
	//},
	Run: register,
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func register(cmd *cobra.Command, args []string) {
	log.Println("Registering with servers...")

	for _, element := range cfgTree.Products.Enabled {
		log.Printf("Registering system with %s...", element)

		if element == "monitor" {
			err := regMonitor()
			if err != nil {
				log.Fatalf("%s\n", err)
			}
		}
	}

}

func regMonitor() error {
	var ok bool
	var tmpServer string
	saveLimit := 80

	if len(cfgTree.Products.Monitor.Servers) < 1 {
		return errors.New("No OP5 Monitor servers provided")
	}

	cfgMonitor := op5monitor.Setup("", cfgTree.Account.Name, cfgTree.Account.Pass, "json", "json")

	data := make(map[string]interface{})
	data["host_name"] = cfgTree.Host.Name
	data["template"] = cfgTree.Host.Template

	for _, server := range cfgTree.Products.Monitor.Servers {
		log.Printf("Trying to register system with server %s...", server)

		// Creating a slice of strings to join to form the complete URL.
		//s := []string{"https:/", server, "api"}
		tmpServer = server
		cfgMonitor.URL = strings.Join([]string{"https:/", server, "api"}, "/")

		ok = op5monitor.HostCreate(cfgMonitor, data, false)

		if ok {
			break
		}
	}

	ok = op5monitor.QueueSave(cfgMonitor)

	if !ok {
		return fmt.Errorf("Queue not saved. Host %s is in limbo.", cfgTree.Host.Name)
	} else {
		log.Printf("Queue saved. %s created on %s.\n", cfgTree.Host.Name, tmpServer)
	}

	for index, metric := range cfgTree.Metrics.Enabled {
		log.Printf("Metrics config tree: %+v\n", cfgTree.Metrics)
		// This isn't really dynamic, and there should be a better way.
		if metric == "cpu" {
			data["template"] = cfgTree.Metrics.CPU.Template
			data["service_description"] = cfgTree.Metrics.CPU.Description
		}

		log.Printf("Trying to register %s service on host %s with server %s",
			strings.ToUpper(metric),
			cfgTree.Host.Name,
			tmpServer)

		ok = op5monitor.ServiceCreate(cfgMonitor, data)

		if !ok {
			log.Fatalf("Unable to register %s service on host %s with server %s",
				strings.ToUpper(metric),
				cfgTree.Host.Name,
				tmpServer)
		} else {
			log.Printf("Registered %s service on host %s with server %s",
				strings.ToUpper(metric),
				cfgTree.Host.Name,
				tmpServer)
		}

		if index >= saveLimit {
			ok = op5monitor.QueueSave(cfgMonitor)
			if !ok {
				log.Fatalf("Queue not saved. %s service on host %s is in limbo.", cfgTree.Host.Name)
			}
		}
	}

	// TODO: Write a queue check command to check if there is anything in the save queue.
	ok = op5monitor.QueueSave(cfgMonitor)
	if !ok {
		return fmt.Errorf("Queue not saved. Services for host %s are in limbo.", cfgTree.Host.Name)
	} else {
		log.Printf("Queue saved. Services for %s created on %s.\n", cfgTree.Host.Name, tmpServer)
	}

	return nil
}
