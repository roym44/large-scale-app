package config

type CacheConfig struct {
	BaseConfig
	RegistryAddresses []string `yaml:"registry_addresses"`
}
