package entities_test

import (
	"testing"

	"github.com/quarkey/iot/entities"
)

func TestGetSpecificSensorDataPoint(t *testing.T) {
	tests := []struct {
		name           string
		args           string
		wantDataset_id int64
		wantColumn     int64
	}{
		{
			"GetSpecificSensorDataPoint",
			"d0c1",
			0,
			1,
		},
		{
			"GetSpecificSensorDataPoint",
			"d112c112",
			112,
			112,
		},
		{
			"GetSpecificSensorDataPoint",
			"d12345c12345",
			12345,
			12345,
		},
		{
			"GetSpecificSensorDataPoint",
			"d01c12345",
			01,
			12345,
		},
		{
			"GetSpecificSensorDataPoint",
			"d12345c01",
			12345,
			01,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDataset_id, gotColumn := entities.GetSpecificSensorDataPoint(tt.args)
			if gotDataset_id != tt.wantDataset_id {
				t.Errorf("GetSpecificSensorDataPoint() gotDataset_id = %v, want %v", gotDataset_id, tt.wantDataset_id)
			}
			if gotColumn != tt.wantColumn {
				t.Errorf("GetSpecificSensorDataPoint() gotColumn = %v, want %v", gotColumn, tt.wantColumn)
			}
		})
	}
}
