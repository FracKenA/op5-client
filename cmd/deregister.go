package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	op5monitor "github.com/FracKenA/op5-client/op5"
)

// deregisterCmd represents the deregister command
var deregisterCmd = &cobra.Command{
	Use:   "deregister",
	Short: "Removes a host and any services from OP5 Monitor.",
	Long:  ``,
	Run:   deregister,
}

func init() {
	rootCmd.AddCommand(deregisterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deregisterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deregisterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deregister(cmd *cobra.Command, args []string) {
	log.Println("Removing host from servers...")

	for _, element := range cfgTree.Products.Enabled {
		log.Printf("Removing system from %s...", element)

		if element == "monitor" {
			err := deregMonitor()
			if err != nil {
				log.Fatalf("%s\n", err)
			}
		}
	}
}

func deregMonitor() error {
	var tmpServer string

	if len(cfgTree.Products.Monitor.Servers) < 1 {
		return errors.New("No Monitor server provided.")
	}

	cfgMonitor := op5monitor.Setup(
		"",
		cfgTree.Account.Name,
		cfgTree.Account.Pass,
		"json",
		"json",
	)

	for _, server := range cfgTree.Products.Monitor.Servers {
		log.Printf("Trying to remove system with server %s...", server)
		tmpServer = server

		cfgMonitor.URL = strings.Join([]string{"https:/", server, "api"}, "/")

		if op5monitor.HostDelete(cfgMonitor, cfgTree.Host.Name) {
			break
		} else {
			log.Fatalf(
				"Unable to remove %s from server %s",
				cfgTree.Host.Name,
				server,
			)
		}
	}

	if !op5monitor.QueueSave(cfgMonitor) {
		return fmt.Errorf(
			"Unable to save queue. Host %s is in limbo on server %s",
			cfgTree.Host.Name,
			tmpServer,
		)
	}

	return nil
}
