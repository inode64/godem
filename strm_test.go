package godem_test

import (
	"log"
	"testing"

	"github.com/inode64/godem"
	"github.com/stretchr/testify/assert"
)

func checkSrtmFileName(t *testing.T, lat, lon float64, expectedZip, expectedFile, expectedDem string) {
	dem, zip, file := godem.GetSrtm(lat, lon)
	log.Printf("Checking Lat: %f Lon: %f", lat, lon)
	if file != expectedFile {
		t.Errorf("SRTM FILE for (%v, %v) should be %s but is %s", lat, lon, expectedFile, file)
	}
	if zip != expectedZip {
		t.Errorf("SRTM ZIP for (%v, %v) should be %s but is %s", lat, lon, expectedZip, zip)
	}
	if dem != expectedDem {
		t.Errorf("DEM for (%v, %v) should be %s but is %s", lat, lon, expectedDem, dem)
	}
}

func TestFindSrtmFileName(t *testing.T) {
	checkSrtmFileName(t, 45, 13, "L33", "N45E013.hgt", godem.DEM1)
	checkSrtmFileName(t, 45.1, 13, "L33", "N45E013.hgt", godem.DEM1)
	checkSrtmFileName(t, 44.9, 13, "L33", "N44E013.hgt", godem.DEM1)
	checkSrtmFileName(t, 45, 13.1, "L33", "N45E013.hgt", godem.DEM1)
	checkSrtmFileName(t, 45, 12.9, "L33", "N45E012.hgt", godem.DEM1)
	checkSrtmFileName(t, 25, -80, "G17", "N25W080.hgt", godem.DEM3)
	checkSrtmFileName(t, 25, -80.1, "G17", "N25W081.hgt", godem.DEM1)
	checkSrtmFileName(t, 25, -79.9, "G17", "N25W080.hgt", godem.DEM3)
	checkSrtmFileName(t, 25.1, -80, "G17", "N25W080.hgt", godem.DEM3)
	checkSrtmFileName(t, -32, 152, "SH56", "S32E152.hgt", godem.DEM3)
	checkSrtmFileName(t, 72.2342, -55.0033, "S21", "n72w056.hgt", godem.DEM1)

	// This file don't exists but the get_file_name is expected to return the supposed file:
	checkSrtmFileName(t, 0, 0, "", "", "")
}

func TestElevation(t *testing.T) {
	strm, err := godem.NewSrtm(nil)
	assert.NoError(t, err)

	ele, dem, err := strm.GetElevation(43.37012643, -8.39114853)
	assert.NoError(t, err)
	assert.Equal(t, ele, 20.0)
	assert.Equal(t, dem, godem.DEM1)
}
