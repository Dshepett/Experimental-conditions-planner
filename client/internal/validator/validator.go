package validator

import "strconv"

func IsValidCsv(data [][]string) bool {
	if len(data[0]) != 3 {
		return false
	}
	for _, i := range data {
		if val, err := strconv.ParseFloat(i[0], 64); err != nil {
			return false
		} else {
			if val < 0 {
				return false
			}
		}
		if val, err := strconv.ParseFloat(i[1], 64); err != nil {
			return false
		} else {
			if val < 0 || val > 100 {
				return false
			}
		}
		if val, err := strconv.Atoi(i[2]); err != nil {
			return false
		} else {
			if val < 0 {
				return false
			}
		}
	}
	return true
}
