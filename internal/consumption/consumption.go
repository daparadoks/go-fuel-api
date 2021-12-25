package consumption

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Service struct {
	DB *gorm.DB
}

type ConsumptionService interface {
	Get(id uint) (Consumption, error)
	GetList(memberId uint) ([]Consumption, error)
	Add(consumption Consumption) (Consumption, error)
	Update(id uint, newConsumption Consumption) (Consumption, error)
	Delete(id uint) error
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}

type Consumption struct {
	gorm.Model
	MemberId       uint
	VehicleId      uint
	Odometer       uint
	Price          float32
	FuelAmount     float32
	CurrencyTypeId uint
	AmountTypeId   uint
	FuelupDate     time.Time
	CityPercent    uint
	IsPartial      bool
	FuelBrand      string
}

func (s *Service) Get(id uint) (Consumption, error) {
	var consumption Consumption
	if result := s.DB.First(&consumption).Where("Id=?", id); result.Error != nil {
		return Consumption{}, result.Error
	}
	return consumption, nil
}

func (s *Service) GetList(memberId uint) ([]Consumption, error) {
	var consumptions []Consumption
	if result := s.DB.Find(&consumptions).Where("MemberId=?", memberId); result.Error != nil {
		return []Consumption{}, result.Error
	}
	return consumptions, nil
}

func (s *Service) Add(consumption Consumption) (Consumption, error) {
	if result := s.DB.Save(&consumption); result.Error != nil {
		return Consumption{}, result.Error
	}
	return consumption, nil
}

func (s *Service) Update(id uint, newConsumption Consumption) (Consumption, error) {
	consumption, err := s.Get(id)
	if err != nil {
		return Consumption{}, err
	}

	if result := s.DB.Model(&consumption).Updates(newConsumption); result.Error != nil {
		return Consumption{}, result.Error
	}

	return consumption, nil
}

func (s *Service) Delete(id uint) error {
	if result := s.DB.Delete(&Consumption{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}
