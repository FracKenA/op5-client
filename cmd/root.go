package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfgTree ConfigTreeRoot
var verbose bool

type ConfigTreeHost struct {
	Name     string `json:"host_name"`
	Template string `json:"template"`
}

type ConfigTreeAccount struct {
	Name string
	Pass string
}

type ConfigTreeMonitor struct {
	Servers []string
	Secure  bool
}

type ConfigTreeProducts struct {
	Enabled []string
	Monitor ConfigTreeMonitor
}

type ConfigTreeCheck struct {
	Type        string
	Description string
	Template    string
	Interval    int
	Warning     string
	Critical    string
}

// I'm not sure explicitly listing the checks is the best way to do this, but it's what I'm going with.
type ConfigTreeMetrics struct {
	Enabled   []string
	Server    string
	HostAlive ConfigTreeCheck
	CPU       ConfigTreeCheck
}

type ConfigTreeRoot struct {
	Reload   bool
	Logfile  string
	Host     ConfigTreeHost
	Account  ConfigTreeAccount
	Products ConfigTreeProducts
	Metrics  ConfigTreeMetrics
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "satellite",
	Short: "Sends passive checks to OP5 Monitor.",
	Long: `OP5 Satellite transmits check results from the host to OP5 Monitor.

It also has functions to register and deregister itself. Register will
setup the host in OP5 Monitor and add configured passive checks, and
deregister will remove the host from OP5 Monitor.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {},
	//Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//TODO: Create alias to -c for --config.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Set the location of the config file.")
	//TODO: Create alias to -v for --verbose.
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Increase the chattyness of the application.")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "", false, "Show the version.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config named "satellite" (without extension) in directories.
		viper.AddConfigPath("/usr/local/etc/op5")
		viper.AddConfigPath("/etc/op5")
		viper.AddConfigPath("./config")
		viper.SetConfigName("satellite")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		os.Exit(3)
	}

	err := viper.Unmarshal(&cfgTree)
	if err != nil {
		log.Printf("Unable to decode config into struct, %v\n", err)
		os.Exit(3)
	}
}

func run(cmd *cobra.Command, args []string) {
	log.Printf("Config Tree:\n%+v\n", cfgTree)
}
