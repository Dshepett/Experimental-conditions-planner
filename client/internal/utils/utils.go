package utils

import (
	"client/internal/models"
	"client/internal/validator"
	"errors"
	"strconv"
)

func ConvertToCharacteristics(data [][]string) ([]models.Characteristics, error) {
	if !validator.IsValidCsv(data) {
		return nil, errors.New("invalid csv file")
	}
	zoles := []models.Characteristics{}
	for _, zole := range data {
		size, _ := strconv.ParseFloat(zole[0], 64)
		consistance, _ := strconv.ParseFloat(zole[1], 64)
		stability, _ := strconv.ParseFloat(zole[2], 64)
		zoles = append(zoles, models.Characteristics{Size: size, Consistence: consistance, Stability: stability})
	}
	return zoles, nil
}
