package responses

import (
	"time"

	"github.com/daparadoks/go-fuel-api/internal/consumption"
)

type ConsumptionModel struct {
	Id             uint
	VehicleName    string
	Odometer       uint
	Distance       uint
	Price          float32
	FuelAmount     float32
	Avarage        float32
	CurrencyTypeId uint
	AmountTypeId   uint
	FuelupDate     time.Time
	CityPercent    uint
	IsPartial      bool
	FuelBrand      string
}

type ConsumptionListModel struct {
	TotalFuelUp uint
	BestAvarage float32
	LastAvarage float32
	Fuelups     []ConsumptionModel
}

func InitConsumptionListModel(consumptions []consumption.Consumption) *ConsumptionListModel {
	model := new(ConsumptionListModel)
	previousItem := new(ConsumptionModel)
	for i := 0; i < len(consumptions); i++ {
		previousOdometer := uint(0)
		previousFuelAmount := float32(0)
		if i > 0 {
			previousOdometer = previousItem.Odometer
		}
		modelItem := InitConsumptionModel(consumptions[i], previousOdometer, previousFuelAmount)
		model.Fuelups = append(model.Fuelups, *modelItem)
		if modelItem.IsPartial {
			previousFuelAmount += modelItem.FuelAmount
		} else {
			previousOdometer = modelItem.Odometer
			previousFuelAmount = 0
		}
	}

	return model
}

func InitConsumptionModel(consumption consumption.Consumption, previousOdometer uint, previousFlueAmount float32) *ConsumptionModel {
	model := new(ConsumptionModel)
	model.Id = consumption.ID
	model.Odometer = consumption.Odometer
	model.Distance = 0
	model.IsPartial = consumption.IsPartial
	if previousOdometer > 0 {
		model.Distance = previousOdometer - consumption.Odometer
	}
	model.Price = consumption.Price
	model.FuelAmount = consumption.FuelAmount
	model.Avarage = 0
	if model.Distance > 0 && !model.IsPartial {
		model.Avarage = (previousFlueAmount + model.FuelAmount) / (float32(model.Distance) / 100)
	}
	model.CurrencyTypeId = consumption.CurrencyTypeId
	model.AmountTypeId = consumption.AmountTypeId
	model.FuelupDate = consumption.FuelupDate
	model.CityPercent = consumption.CityPercent
	model.FuelBrand = consumption.FuelBrand

	return model
}
