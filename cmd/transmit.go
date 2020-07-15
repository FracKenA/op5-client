package cmd

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	op5monitor "github.com/FracKenA/op5-client/op5"
)

type Check struct {
	Type string
	Data op5monitor.CheckResults
}

// transmitCmd represents the transmit command
var transmitCmd = &cobra.Command{
	Use:   "transmit",
	Short: "Transmists metrics to configured platforms.",
	Long:  ``,
	Run:   start,
}

func init() {
	rootCmd.AddCommand(transmitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transmitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transmitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func start(cmd *cobra.Command, args []string) {
	var checkList []time.Ticker

	cfgServer := op5monitor.Setup(
		cfgTree.Metrics.Server,
		cfgTree.Account.Name,
		cfgTree.Account.Pass,
		"json",
		"json",
	)

	log.Println("Start transmission of metrics...")

	// Channel to capture signals.
	flags := make(chan os.Signal)
	// Channel to control the goroutines.
	work := make(chan bool)
	// Channel to store the messages to send.
	txq := make(chan Check, 32)

	signal.Notify(flags, syscall.SIGINT, syscall.SIGTERM)

	// Start messages transmit schedule.
	tx := scheduleTX(
		msgSend,
		cfgServer,
		100*time.Millisecond,
		work,
		txq,
	)
	// Running the first checks manually since we want the first checks to fire
	// off as soon as we start.
	go hostAlive(txq)
	// Scheduling the host alive check to fire off at set intervals.
	host := scheduleCheck(
		hostAlive,
		time.Duration(cfgTree.Metrics.HostAlive.Interval)*time.Second,
		work,
		txq,
	)

	checkList = append(checkList, *host)

	for _, check := range cfgTree.Metrics.Enabled {
		check = strings.ToLower(check)

		switch check {
		case "cpu":
			go checkCPUPercent(txq)
			cpu := scheduleCheck(
				checkCPUPercent,
				time.Duration(cfgTree.Metrics.CPU.Interval)*time.Second,
				work,
				txq,
			)

			checkList = append(checkList, *cpu)
		}
	}

	sig := <-flags
	log.Printf("Caught %s signal.", sig)
	close(work)

	log.Println("Stopping checks...")
	for _, check := range checkList {
		check.Stop()
	}

	tx.Stop()
	log.Println("End transmission of metrics.")
}

func scheduleTX(
	send func(<-chan Check, op5monitor.ConfigTree),
	config op5monitor.ConfigTree,
	interval time.Duration,
	work <-chan bool,
	txq <-chan Check,
) (clock *time.Ticker) {
	clock = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-clock.C:
				send(txq, config)
			case <-work:
				return
			}
		}
	}()
	return
}

func scheduleCheck(
	check func(chan<- Check),
	interval time.Duration,
	work <-chan bool,
	txq chan<- Check,
) (clock *time.Ticker) {
	clock = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-clock.C:
				check(txq)
			case <-work:
				return
			}
		}
	}()
	return
}

func msgSend(txq <-chan Check, config op5monitor.ConfigTree) {
	for msg := range txq {
		log.Printf("Message: %+v", msg)
		op5monitor.SendCheck(config, msg.Type, msg.Data)
	}
}

func hostAlive(txq chan<- Check) {
	txq <- Check{
		Type: "host",
		Data: op5monitor.CheckResults{
			Hostname:     cfgTree.Host.Name,
			StatusCode:   0,
			PluginOutput: "I'm alive!",
		},
	}
}

func checkCPUPercent(txq chan<- Check) {
	txq <- Check{
		Type: "service",
		Data: op5monitor.CheckResults{
			Hostname:           cfgTree.Host.Name,
			StatusCode:         0,
			PluginOutput:       "CPU OK | percent=1%",
			ServiceDescription: cfgTree.Metrics.CPU.Description,
		},
	}
}
