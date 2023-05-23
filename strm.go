package godem

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/lukeroth/gdal"
)

const (
	DEM1          = "dem1"
	DEM3          = "dem3"
	DEMURL0       = "http://viewfinderpanoramas.org/"
	DEMURL1       = "http://dem.gpxsee.org/"
	DEMURL2       = "https://step.esa.int/auxdata/dem/SRTMGL1/"
	SOURCE_VIEW   = 0
	SOURCE_GPXSEE = 1
	SOURCE_ESA    = 2
)

type Srtm struct {
	storage SrtmLocalStorage
	source  int
}

func NewSrtm(source int) (*Srtm, error) {
	if source != SOURCE_VIEW && source != SOURCE_GPXSEE && source != SOURCE_ESA {
		return nil, fmt.Errorf("invalid source")
	}
	storage, err := NewLocalFileSrtmStorage(source)
	if err != nil {
		return nil, err
	}

	return &Srtm{storage: storage, source: source}, nil
}

func (self *Srtm) GetSrtm(lat, lon float64) (string, string, string, string) {
	if self.source == SOURCE_GPXSEE {
		file := getSrtmFileNameAndCoordinates(lat, lon)
		return DEM1, file, file, DEMURL1 + getLat2Coordinates(lat) + "/" + file + ".zip"
	}
	if self.source == SOURCE_ESA {
		file := getSrtmFileNameAndCoordinates(lat, lon)
		zip := fmt.Sprintf("%s%s.SRTMGL1.hgt", getLat2Coordinates(lat), getLon2Coordinates(lon))
		return DEM1, file, file, DEMURL2 + "/" + zip + ".zip"
	}
	zip, file := getDem(dem1, lat, lon)
	if len(zip) != 0 {
		return DEM1, zip, file, DEMURL0 + DEM1 + "/" + zip + ".zip"
	}

	zip, file = getDem(dem3, lat, lon)
	if len(zip) != 0 {
		return DEM3, zip, file, DEMURL0 + DEM3 + "/" + zip + ".zip"
	}
	return "", "", "", ""
}

func getDem(data string, lat, lon float64) (string, string) {
	var fileStructure map[string][]string
	_ = json.Unmarshal([]byte(data), &fileStructure)

	lookupFile := getSrtmFileNameAndCoordinates(lat, lon)

	for zip, files := range fileStructure {
		for _, file := range files {
			if strings.EqualFold(file, lookupFile) {
				return zip, file
			}
		}
	}

	return "", ""
}

func getSrtmFileNameAndCoordinates(lat, lon float64) string {
	return fmt.Sprintf("%s%s.hgt", getLat2Coordinates(lat), getLon2Coordinates(lon))
}

func getLat2Coordinates(lat float64) string {
	northSouth := 'S'
	if lat >= 0 {
		northSouth = 'N'
	}

	latPart := int(math.Abs(math.Floor(lat)))

	return fmt.Sprintf("%s%02d", string(northSouth), latPart)
}

func getLon2Coordinates(lon float64) string {
	eastWest := 'W'
	if lon >= 0 {
		eastWest = 'E'
	}

	lonPart := int(math.Abs(math.Floor(lon)))

	return fmt.Sprintf("%s%03d", string(eastWest), lonPart)
}

func (self *Srtm) loadContents(dem, zip, file, url string) error {
	_, err := self.storage.FileExists(dem, zip, file)
	if err != nil {
		client := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}

		resp, err := client.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		fileInArchiveBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err = self.storage.Unzip(dem, zip, fileInArchiveBytes); err != nil {
			return err
		}
	}

	return nil
}

func (self *Srtm) GetElevation(lat, lon float64) (float64, string, error) {
	dem, zip, file, url := self.GetSrtm(lat, lon)
	if len(dem) == 0 {
		return 0, "", fmt.Errorf("no dem found")
	}
	err := self.loadContents(dem, zip, file, url)
	if err != nil {
		return 0, "", err
	}
	path, err := self.storage.FileExists(dem, zip, file)
	if err != nil {
		return 0, "", err
	}

	ele, err := GetElevationFromLocalFile(path, lat, lon)
	return ele, dem, err
}

func GetElevationFromLocalFile(path string, lat, lon float64) (float64, error) {
	ds, err := gdal.Open(path, gdal.ReadOnly)
	if err != nil {
		return 0, err
	}
	defer ds.Close()

	rb := ds.RasterBand(1)

	// Convert lat/lon to pixel coordinates
	geoTransform := ds.GeoTransform()
	if err != nil {
		return 0, err
	}

	pixelX := int((lon - geoTransform[0]) / geoTransform[1])
	pixelY := int((lat - geoTransform[3]) / geoTransform[5])

	// Read the pixel value
	buffer := make([]int32, 1)
	err = rb.IO(gdal.Read, pixelX, pixelY, 1, 1, buffer, 1, 1, 0, 0)
	if err != nil {
		return 0, err
	}

	return float64(buffer[0]), nil
}
