package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//parseKismet parses the specified Kismet file and returns a Points array with all the values
func parseKismet(file *string, bssids []string, filter *int) (points Points) {
	db, err := sql.Open("sqlite3", *file)
	if err != nil {
		panic("Ensure the Kismit file exists")
	}

	rows, err := db.Query("SELECT sourcemac,lat,lon,signal FROM packets")
	if err != nil {
		panic("Kismet database corrupt")
	}
	defer rows.Close()
	for rows.Next() {
		var (
			mac      string
			lat, lon float64
			sig      int
		)

		if err := rows.Scan(&mac, &lat, &lon, &sig); err != nil {
			fmt.Println(err)
		}

		points = append(points, Point{lon, lat, sig, mac})
	}
	rows.Close()
	points = filterBSSID(points, bssids, filter)
	return
}
