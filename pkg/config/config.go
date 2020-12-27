package config

import (
	"fmt"
	"os"

	"github.com/namsral/flag"
)

// Config provides all the configuration needed
// for the Cloud::1 server.
type Config struct {
	FileSystem     *string
	DataDirectory  *string
	RunOnHost      *bool
	ServerIP       *string
	HostsPath      *string
	AWSServices    *string
	GCloudServices *string
	AzureServices  *string
	Debug          *bool
}

// Load deals with loading configuration from
// either environment variables, a config file or command line options.
func Load() (*Config, error) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	args := os.Args[1:]
	if flagSet.Lookup(flag.DefaultConfigFlagname) == nil {
		flagSet.String(flag.DefaultConfigFlagname, "", "path to config file")
	}
	config := prepare(flagSet)
	err := flagSet.Parse(args)
	if err != nil {
		return nil, err
	}
	err = validate(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func prepare(flagSet *flag.FlagSet) *Config {
	var fileSystem string
	flagSet.StringVar(
		&fileSystem,
		"cloud_one_file_system",
		"",
		"The file system to use for Cloud::1 custom emulators, the only choice is memory."+
			" Any other value will mean the os is used.",
	)

	var dataDirectory string
	flagSet.StringVar(
		&dataDirectory,
		"cloud_one_data_dir",
		"/lib/data",
		"The location in the docker container to store all the data for the cloud service emulators. "+
			"(If your file system is set to memory, then not all data for services will be persisted to disk)",
	)

	var runOnHost bool
	flagSet.BoolVar(
		&runOnHost,
		"cloud_one_run_on_host",
		false,
		"If set, this will enable in-process features that are only available to privileged host applications. "+
			"It will also embed the functionality to interact with the os hosts file in-process as opposed to in "+
			"docker mode where a host agent that communicates over a unix socket is needed.",
	)

	var serverIP string
	flagSet.StringVar(
		&serverIP,
		"cloud_one_ip",
		"172.18.0.22",
		"The IP Address the cloud one server is running on, this is ignored when running the server directly on the host.",
	)

	var hostsPath string
	flagSet.StringVar(
		&hostsPath,
		"cloud_one_hosts_path",
		"",
		"A custom path to the hosts file on the host machine,"+
			" otherwise defaults to the correct hosts file for the OS the host agent/server directly on the host is running on.",
	)

	var awsServices string
	flagSet.StringVar(
		&awsServices,
		"cloud_one_aws_services",
		"",
		"AWS Services to run emulations for.",
	)

	var gcloudServices string
	flagSet.StringVar(
		&awsServices,
		"cloud_one_gcloud_services",
		"",
		"Google Cloud Services to run emulations for.",
	)

	var azureServices string
	flagSet.StringVar(
		&azureServices,
		"cloud_one_azure_services",
		"",
		"Azure Services to run emulations for.",
	)

	var debug bool
	flagSet.BoolVar(
		&debug,
		"debug",
		false,
		"Whether or not to run the application in debug mode, where debug logs will be written to stdout.",
	)

	return &Config{
		FileSystem:     &fileSystem,
		DataDirectory:  &dataDirectory,
		RunOnHost:      &runOnHost,
		ServerIP:       &serverIP,
		HostsPath:      &hostsPath,
		AWSServices:    &awsServices,
		GCloudServices: &gcloudServices,
		AzureServices:  &azureServices,
		Debug:          &debug,
	}
}

func validate(config *Config) error {
	noAWSServices := *config.AWSServices == ""
	noAzureServices := *config.AzureServices == ""
	noGCloudServices := *config.GCloudServices == ""
	if noAWSServices && noAzureServices && noGCloudServices {
		return fmt.Errorf("You must select some services to run for at least one cloud provider to emulate")
	}
	return nil
}
