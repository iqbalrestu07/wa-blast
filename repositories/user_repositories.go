package repositories

import (
	"wa-blast/models"
)

// UserRepository ...
type UserRepository interface {
	GetAuth(where string, args ...interface{}) (models.CompaniesAuth, error)
	GetDevice(where string, args ...interface{}) (models.CompaniesDevice, error)
	GetDevices(where string, args ...interface{}) ([]models.CompaniesDevice, error)
	CreateDevice(Data models.CompaniesDevice) error
	UpdateDevice(ID string, Data models.CompaniesDevice) (models.CompaniesDevice, error)
	DeleteDevice(ID string) error
}

type userRepository struct{}

func initUserRepository() UserRepository {
	// Prepare statements
	var r userRepository
	return &r
}

func (s *userRepository) GetAuth(where string, args ...interface{}) (models.CompaniesAuth, error) {
	var data models.CompaniesAuth
	err := DBConnect.Table("companies_auth").Where(where, args...).Find(&data).Error
	return data, err
}

func (s *userRepository) GetDevice(where string, args ...interface{}) (models.CompaniesDevice, error) {
	var data models.CompaniesDevice
	err := DBConnect.Table("companies_device").Where(where, args...).Find(&data).Error
	return data, err
}

func (s *userRepository) GetDevices(where string, args ...interface{}) ([]models.CompaniesDevice, error) {
	var data []models.CompaniesDevice
	err := DBConnect.Table("companies_device").Where(where, args...).Find(&data).Error
	return data, err
}

func (s *userRepository) CreateDevice(Data models.CompaniesDevice) error {
	err := DBConnect.Table("companies_device").Create(&Data).Error
	return err
}

func (s *userRepository) UpdateDevice(ID string, Data models.CompaniesDevice) (models.CompaniesDevice, error) {

	err := DBConnect.Debug().Table("companies_device").Model(&Data).Where("id = ?", ID).Updates(&Data).Error

	return Data, err
}

func (s *userRepository) DeleteDevice(ID string) error {

	err := DBConnect.Debug().Table("companies_device").Where("id = ?", ID).Delete(models.CompaniesDevice{}).Error

	return err
}
