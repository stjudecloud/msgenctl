package internal

import (
	"fmt"

	"github.com/spf13/pflag"
)

type ServiceConfig struct {
	BaseURL   string
	AccessKey string
}

type StorageConfig struct {
	AccountName   string
	AccountKey    string
	ContainerName string
}

type InputConfig struct {
	Storage  StorageConfig
	BlobName string
}

type OutputConfig struct {
	Storage    StorageConfig
	Basename   string
	Overwrite  bool
	IncludeLog bool
}

type ProcessConfig struct {
	Name string
	Args string
}

type OptionalArgsConfig struct {
	EmitRefConfidence ReferenceConfidenceMode
	BgzipOutput       bool
}

type SubmitConfig struct {
	Service ServiceConfig

	Input       InputConfig
	Process     ProcessConfig
	Description string

	Output OutputConfig

	OptionalArgs OptionalArgsConfig
}

func SubmitConfigFromFlags(flags *pflag.FlagSet) (SubmitConfig, error) {
	config := SubmitConfig{}

	serviceConfig, err := ServiceConfigFromFlags(flags)

	if err != nil {
		return config, err
	}

	config.Service = serviceConfig

	inputConfig, err := inputConfigFromFlags(flags)

	if err != nil {
		return config, err
	}

	config.Input = inputConfig

	processConfig, err := processConfigFromFlags(flags)

	if err != nil {
		return config, err
	}

	config.Process = processConfig

	description, err := flags.GetString("description")

	if err != nil {
		return config, err
	}

	config.Description = description

	outputConfig, err := outputConfigFromFlags(flags)

	if err != nil {
		return config, err
	}

	config.Output = outputConfig

	optionalArgsConfig, err := optionalArgsConfigFromFlags(flags)

	if err != nil {
		return config, err
	}

	config.OptionalArgs = optionalArgsConfig

	return config, nil
}

func ServiceConfigFromFlags(flags *pflag.FlagSet) (ServiceConfig, error) {
	config := ServiceConfig{}

	baseURL, err := flags.GetString("base-url")

	if err != nil {
		return config, err
	}

	config.BaseURL = baseURL

	accessKey, err := flags.GetString("access-key")

	if err != nil {
		return config, err
	}

	config.AccessKey = accessKey

	return config, nil
}

func inputConfigFromFlags(flags *pflag.FlagSet) (InputConfig, error) {
	config := InputConfig{}

	storageConfig, err := storageConfigFromFlags(flags, "input")

	if err != nil {
		return config, err
	}

	config.Storage = storageConfig

	blobName, err := flags.GetString("input-blob-name")

	if err != nil {
		return config, err
	}

	config.BlobName = blobName

	return config, nil
}

func processConfigFromFlags(flags *pflag.FlagSet) (ProcessConfig, error) {
	config := ProcessConfig{}

	processName, err := flags.GetString("process-name")

	if err != nil {
		return config, err
	}

	config.Name = processName

	processArgs, err := flags.GetString("process-args")

	if err != nil {
		return config, err
	}

	config.Args = processArgs

	return config, nil
}

func outputConfigFromFlags(flags *pflag.FlagSet) (OutputConfig, error) {
	config := OutputConfig{}

	storageConfig, err := storageConfigFromFlags(flags, "output")

	if err != nil {
		return config, err
	}

	config.Storage = storageConfig

	basename, err := flags.GetString("output-basename")

	if err != nil {
		return config, err
	}

	config.Basename = basename

	overwrite, err := flags.GetBool("output-overwrite")

	if err != nil {
		return config, err
	}

	config.Overwrite = overwrite

	includeLog, err := flags.GetBool("output-include-log")

	if err != nil {
		return config, err
	}

	config.IncludeLog = includeLog

	return config, nil
}

func storageConfigFromFlags(flags *pflag.FlagSet, prefix string) (StorageConfig, error) {
	var key string

	config := StorageConfig{}

	key = fmt.Sprintf("%v-storage-account-name", prefix)
	accountName, err := flags.GetString(key)

	if err != nil {
		return config, err
	}

	config.AccountName = accountName

	key = fmt.Sprintf("%v-storage-account-key", prefix)
	accountKey, err := flags.GetString(key)

	if err != nil {
		return config, err
	}

	config.AccountKey = accountKey

	key = fmt.Sprintf("%v-storage-container-name", prefix)
	containerName, err := flags.GetString(key)

	if err != nil {
		return config, err
	}

	config.ContainerName = containerName

	return config, nil
}

func optionalArgsConfigFromFlags(flags *pflag.FlagSet) (OptionalArgsConfig, error) {
	config := OptionalArgsConfig{}

	emitRefConfidence, err := flags.GetString("emit-ref-confidence")

	if err != nil {
		return config, err
	}

	referenceConfidenceMode, err := ParseReferenceConfidenceMode(emitRefConfidence)

	if err != nil {
		return config, err
	}

	config.EmitRefConfidence = referenceConfidenceMode

	bgzipOutput, err := flags.GetBool("bgzip-output")

	if err != nil {
		return config, err
	}

	config.BgzipOutput = bgzipOutput

	return config, nil
}
