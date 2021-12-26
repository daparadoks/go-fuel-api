package rabbit

import "github.com/daparadoks/go-fuel-api/config"

var quesConfig []config.QueConfig

func (c *Client) getRegisteredQues() []config.QueConfig{
	if quesConfig !=nil{
		return quesConfig
	}

	configs:= make([]config.QueConfig, 0)

	configs = append(configs, c.quesConfig.Consumption.ConsumptionCreated)
	configs = append(configs, c.quesConfig.Consumption.ConsumptionUpdated
	
	quesConfig = configs
	return configs
}