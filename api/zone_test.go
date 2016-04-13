package api

import (
	"testing"
)

func NewZoneAPI(t *testing.T) ZoneAPI {
	c := NewTestIaaSConnection(t)
	z := ZoneAPI{IaaSAPI{
		API{c},
	}}
	z.connector = c
	return z
}

func TestDescribeZones(t *testing.T) {
	a := NewZoneAPI(t)

	zones, err := a.DescribeZones()
	if err != nil {
		t.Error(err)
	}

	CheckIaaSAPIResponse(t, zones, "DescribeZonesResponse")
}
