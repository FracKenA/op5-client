# OP5 Satellite
Passive check agent for OP5 Monitor with automatic registration and deregistration.

## Installation
### Prereqs
* Go 1.11.4+ installed

### Instructions

1. Clone git repo or unpacked tarball.
2. Move into the directory.
3. Run the go build command.
    * `go build`
4. Run the resulting binary.
    * `./satellite`

## Commands

* register
    Registers OP5 Satellite with OP5 Monitor and sets up passive checks.

    `-t`, `--temporal` Sets Satellite to deregister when the host is shutdown.

* deregister
    Removes OP5 Satellite from OP5 Monitor and associated passive checks.

* transmit
    Sends data to OP5 Monitor.

    `-c </path/to/config>`, `--config=</path/to/config>` Sets the path to the config file. Default: `./config.yml`, `/etc/opt/satellite`

    `-l </path/to/file.log>`, `--logfile=</path/to/file.log>` Sets the location of The log file. Default: `./log/event.log`, `/var/log/op5/satellite.log`

    `-r`, `--reload-config` Watch the config files for changes and reload when changed.

## Config File

Config file is in YAML format, but it will accept TOML, JSON, and some Java config. It's whatever the cobra and viper packages will accept. Whatever format is used, it needs to have the format below.

* reload: bool (Watch the config file for changes and reload it.)
* logfile: string (Path to the log file. (Future))

* host: (Host specific settings.)
  * name: string (Name of host.)
  * template: string (Template to use for the host.)

* account: (Credentials to log into the server.)
  * name: string
  * pass: string

* products: (Which products Satellite should report to.)
  * enabled: (List of enabled products.)
    - monitor
  * monitor: (OP5 Monitor server specific options.)
    * servers: (List of OP5 Monitor servers for registration.)
      - master0.domain.tld
      - master1.domain.tld
    * secure: bool (Bool turning cert validation on or off.)

* metrics: (Options for checks.)
  * enabled: (List of enabled checks.)
    - metricname (Name of check to run.)
  * server: string (Server to report metrics to.)
  * hostalive: (Settings for the host-alive check. This is not optional.)
    * interval: int (Seconds between checks.)
  * metricname: (Check specific settings)
    * servicedescription: string (Name of the service in monitor.)
    * template: string (Name of template to use.)
    * interval: int (Seconds between checks.)
    * warning: string (Nagios range string.)
    * critical: string (Nagios range string.)

## Metrics
* cpupercenttotal
  * Reports the percent utilization of all CPUs in aggregate.
* cpupercenteach
  * Reports the percent utilization of each CPU.
* cputimetotal
  * Reports the amount of time all CPUs has spent performing different kinds of work.
  * CPU       string  `json:"cpu"`
  * User      float64 `json:"user"`
  * System    float64 `json:"system"`
  * Idle      float64 `json:"idle"`
  * Nice      float64 `json:"nice"`
  * Iowait    float64 `json:"iowait"`
  * Irq       float64 `json:"irq"`
  * Softirq   float64 `json:"softirq"`
  * Steal     float64 `json:"steal"`
  * Guest     float64 `json:"guest"`
  * GuestNice float64 `json:"guestNice"`
  * Stolen    float64 `json:"stolen"`
* cputimeeach
  * Reports the amount of time each CPU has spent performing different kinds of work.
  * CPU       string  `json:"cpu"`
  * User      float64 `json:"user"`
  * System    float64 `json:"system"`
  * Idle      float64 `json:"idle"`
  * Nice      float64 `json:"nice"`
  * Iowait    float64 `json:"iowait"`
  * Irq       float64 `json:"irq"`
  * Softirq   float64 `json:"softirq"`
  * Steal     float64 `json:"steal"`
  * Guest     float64 `json:"guest"`
  * GuestNice float64 `json:"guestNice"`
  * Stolen    float64 `json:"stolen"`
