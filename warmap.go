package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"io/ioutil"
)

//To-Do
//Add Signal bleed based signal strength

//Type Definitions and variables
const (
	version = "2.0"
	tool    = "warmap-go"
	usage   = `
 Usage:
   warmap-go [OPTIONS] -f <GPS FILE>
   warmap-go -h | -help
   warmap-go -v

 Options:
   -h, -help              Show usage
   -a                     Aerodump specification switch
   -k                     Kismet database specification switch
   -b <[ list | file ]>   File or comma seperated list of bssids
   -o <file>              HTML Output file
   -p <file>              CSV Output for recorded BSSID values
   -v                     Print version

   -api <key>             Google Maps API key
   -sig <int>             HTLM output signal filtering
`
)

////////////////////
//Type Definitions//
///////////////////

// Points defines a []Point array
type Points []Point

//Point holds X, Y coordinates
type Point struct {
	X, Y  float64
	Dbm   int
	BSSID string
}

/////////////
//Functions//
////////////

//checkError is a generic error check function
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

//sigCheck checks if signal has changed and modifies the base
func sigCheck(baseH *int, baseL *int, c int) {
	if *baseH == 0 || *baseL == 0 {
		*baseH = c
		*baseL = c
	}
	if c > *baseH {
		*baseH = c
	}
	if c < *baseL {
		*baseL = c
	}
}

func (points *Points) printPoints(file *string) {
	// func printPoints(file *string, points *Points) {
	list := make(map[string]int)
	var data bytes.Buffer

	for _, p := range *points {
		if _, ok := list[p.BSSID]; !ok {
			list[p.BSSID] = p.Dbm
		} else if db, ok := list[p.BSSID]; ok && db < p.Dbm {
			list[p.BSSID] = p.Dbm
		}
	}
	for k, v := range list {
		data.WriteString(fmt.Sprintf("%s,%v\n", k, v))
	}
	ioutil.WriteFile(*file, data.Bytes(), 0644)
}

func main() {
	//Parse command line arguments
	var (
		flGPSFile   = flag.String("f", "", "")
		flBSSID     = flag.String("b", "", "")
		flOutFile   = flag.String("o", "", "")
		flAerodump  = flag.Bool("a", false, "")
		flKismet    = flag.Bool("k", false, "")
		flGAPI      = flag.String("api", "", "")
		flPoints    = flag.String("p", "", "")
		flSigFilter = flag.Int("sig", 0, "")
		flVersion   = flag.Bool("v", false, "")
	)

	flag.Usage = func() {
		fmt.Println(usage)
	}

	flag.Parse()
	if *flVersion {
		fmt.Printf("version: %s\n", version)
		os.Exit(0)
	}

	var gpsPoints Points
	if *flAerodump {
		gpsPoints = parseAeroGPS(flGPSFile)
	} else if *flKismet {
		gpsPoints = parseKismet(flGPSFile, parseBssid(flBSSID), flSigFilter)
	} else {
		gpsPoints = parseXML(flGPSFile, parseBssid(flBSSID), flSigFilter)
	}

	if *flPoints != "" {
		gpsPoints.printPoints(flPoints)
	}

	templateBuffer := populateTemplate(&gpsPoints, flGAPI)
	ioutil.WriteFile(*flOutFile, templateBuffer, 0644)
}
