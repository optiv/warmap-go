package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

//filterBSSID returns all GPSXMLPoint structs that have a particular bssid field
func filterBSSID(points Points, bssid []string, filter *int) (filteredPoints Points) {
	for _, i := range points {
		for _, n := range bssid {
			if i.BSSID == n {
				if *filter > 0 && (i.Dbm*-1) < *filter {
					filteredPoints = append(filteredPoints, i)
				} else {
					filteredPoints = append(filteredPoints, i)
				}
			}
		}
	}
	if len(filteredPoints) == 0 {
		log.Fatal("Your BSSID was not found in the file")
	}
	return
}

//parseBssid takes a filename or comma seperated list of BSSIDs
//and outputs an array containing the parsed BSSIDs
func parseBssid(bssids *string) []string {
	var (
		bssidSlice     []string
		tempBssidSlice []string
	)
	r, err := regexp.Compile("(([a-zA-Z0-9]{2}:)){5}[a-zA-Z0-9]{2}")
	checkError(err)
	file, err := os.Open(*bssids)
	if err == nil {
		defer file.Close()
		var lines []string
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		bssidSlice = lines
	} else {
		bssidSlice = strings.Split(*bssids, ",")
	}
	for i := 0; i < len(bssidSlice); i++ {
		if r.MatchString(bssidSlice[i]) {
			tempBssidSlice = append(tempBssidSlice, strings.ToUpper(bssidSlice[i]))
		}
	}
	if len(tempBssidSlice) == 0 {
		log.Fatal("Looks like you didn't have any correctly formatted SSIDs")
	}
	return tempBssidSlice
}
