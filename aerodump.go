package main

import (
	"os"
	"regexp"
	"strings"

	"encoding/json"
	"io/ioutil"
)

type GPSAeroPoint struct {
	Class  string
	Tag    string
	Device string
	Mode   int
	Time   string
	Ept    float64
	Lat    float64
	Lon    float64
	Alt    float64
	Epx    float64
	Epy    float64
	Epv    float64
	Track  float64
	Speed  float64
	Climb  float64
	Eps    float64
	Epc    float64
}

// parse the json file for Aerodump
func parseAeroGPS(file *string) (points Points) {
	gpsFile, err := os.Open(*file)
	if err != nil {
		panic("Ensure the Aero GPS file exists")
	}
	defer gpsFile.Close()
	jsonBytes, err := ioutil.ReadAll(gpsFile)
	if err != nil {
		panic("Error reading gpsFile")
	}

	// Aerodump leverages broken JSON objects on output files
	// regex to split file lines and deliniate between multiple GPS
	re := regexp.MustCompile(`(?:}[\s\t]+|}){"class`)

	for _, line := range strings.Split(string(re.ReplaceAll(jsonBytes, []byte("}\n{\"class"))), "\n") {
		var gpsaero GPSAeroPoint
		json.Unmarshal([]byte(line), &gpsaero)

		// Class TPV @ mode 3 are the Aerodump tracking of valid GPS hits
		// Mode 0 would indicate a false hit and result in 0/0 tracking
		if gpsaero.Class == "TPV" && gpsaero.Mode == 3 {
			points = append(points, Point{gpsaero.Lon, gpsaero.Lat, 0, ""})
		}

	}
	return
}
