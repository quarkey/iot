package entities

import (
	"regexp"
	"strconv"
)

// GetSpecificSensorDataPoint takes a datasource string and extracts dataset id and column.
// e.g. d0c1 will return dataset_id=0, column=1
func GetSpecificSensorDataPoint(datasource string) (dataset_id, column int64) {
	re := regexp.MustCompile("[0-9]+")
	slice := re.FindAllString(datasource, -1)
	s0, err := strconv.ParseInt(slice[0], 10, 64)
	if err != nil {
		return 0, 0
	}
	s1, err := strconv.ParseInt(slice[1], 10, 64)
	if err != nil {
		return 0, 0
	}
	return s0, s1
}
