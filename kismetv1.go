package main

import (
	"bufio"
	"os"
	"strings"

	"encoding/xml"
)

//GPSXMLPoint defines a struct to hold the values
//of the kismet generated gps
type GPSXMLPoint struct {
	Bssid     string  `xml:"bssid,attr"`
	Lat       float64 `xml:"lat,attr"`
	Lon       float64 `xml:"lon,attr"`
	Source    string  `xml:"source,attr"`
	TimeSec   int     `xml:"time-sec,attr"`
	TimeUSec  int     `xml:"time-usec,attr"`
	Spd       float64 `xml:"spd,attr"`
	Heading   float64 `xml:"heading,attr"`
	Fix       int     `xml:"fix,attr"`
	Alt       float64 `xml:"alt,attr"`
	SignalDbm int     `xml:"signal_dbm,attr"`
	NoiseDbm  int     `xml:"noise_dbm,attr"`
}

//parseXML parses the specified XML file and returns a Points array with all the values
func parseXML(file *string, bssids []string, filter *int) (points Points) {
	xmlFile, err := os.Open(*file)
	if err != nil {
		panic("Ensure the GPSXML file exists")
	}
	defer xmlFile.Close()
	xmlScanner := bufio.NewScanner(xmlFile)
	for xmlScanner.Scan() {
		line := xmlScanner.Text()
		if strings.Contains(line, "<gps-point") {
			var gpsxml GPSXMLPoint
			xml.Unmarshal([]byte(line), &gpsxml)
			points = append(points, Point{gpsxml.Lon, gpsxml.Lat, gpsxml.SignalDbm, gpsxml.Bssid})
		}
	}
	points = filterBSSID(points, bssids, filter)
	return
}
