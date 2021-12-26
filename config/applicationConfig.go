package config

type ApplicationConfig struct {
	Rabbit RabbitConfig
}

type RabbitConfig struct {
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	VirtualHost    string `yaml:"virtualHost"`
	ConnectionName string `yaml:"connectionName"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
}

type QuesConfig struct {
	Consumption ConsumptionQueConfig `yaml:"consumption"`
}

type ConsumptionQueConfig struct {
	ConsumptionCreated QueConfig
	ConsumptionUpdated QueConfig
}

type QueConfig struct {
	PrefetchCount int    `yaml:"prefetchCount"`
	ChannelCount  int    `yaml:"channelCount"`
	Exchange      string `yaml:"exchange"`
	ExchangeType  string `yaml:"cxchangeType"`
	RoutingKey    string `yaml:"routingKey"`
	Queue         string `yaml:"queue"`
}
