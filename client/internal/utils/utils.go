package utils

import (
	"client/internal/models"
	"client/internal/validator"
	"errors"
	"fmt"
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

func ConvertToCsv(data []models.SynthesisData) string {
	res := "temperature;time;c_acid;c_ti;acid;treatment;size;consistance;stability\n"
	for _, val := range data {
		res += fmt.Sprintf("%.2f;%.2f;%.2f;%.2f;%s;%.2f;%.2f;%.2f;%.2f\n", val.Conditions.Temperature, val.Conditions.Time,
			val.Conditions.CAcid, val.Conditions.CTi, val.Conditions.Acid, val.Conditions.Treatment, val.Characteristics.Size,
			val.Characteristics.Consistence, val.Characteristics.Stability)
	}
	return res
}
