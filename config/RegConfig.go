package config

type RegConfig struct {
	BaseConfig
	ListenPort int `yaml:"listen_port"`
}
