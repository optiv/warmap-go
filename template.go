package main

import (
	"bytes"
	"fmt"

	"html/template"
)

//Page Holds the Values for html template
type Page struct {
	Lat        float64
	Lng        float64
	Heatmap    template.JS
	ConvexHull template.JS
	Drive      template.JS
	HighDB     template.JS
	LowDB      template.JS
	PathLength int
	Apikey     *string
}

const (
	//tpl defines the HTML template
	tpl = `
<!DOCTYPE HTML SYSTEM>
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">
  <head>
    <title>WarMap</title>
    <script type="text/javascript" src="https://optiv.github.io/warmap-go/src/jquery.min.js"></script>
    <script type="text/javascript" src="https://optiv.github.io/warmap-go/src/gmaps.js"></script>
    <link href="https://optiv.github.io/warmap-go/src/bootstrap.min.css" rel="stylesheet">
    <link href="https://optiv.github.io/warmap-go/src/fontawesome.all.css" rel="stylesheet">
    <style>
      #map {
        display: block;
        width: 100%;
        height: 700;
      }
    </style>
  </head>
  <body>
    <div class="container" style="margin-left: 30px; margin-right: 30px; width: 95%">
      <div class="row">
        <div class="col-xs-12">
				  <div class="row">
					  <p>
						  <div class="col-xs-2"></div>
						  <div class="col-xs-2" style="padding-bottom: 10px">Mapped Points: {{.PathLength}}</div>
						  <div class="col-xs-2">Strongest DB: {{.HighDB}}</div>
						  <div class="col-xs-2">Weakest DB: {{.LowDB}}</div>
					  </p>
          </div>
        </div>
      </div>
			<div class="row">
        <div class="col-xs-1">
            <div class="row">
              <div style="padding-top: 10px"></div>
              <p>
                <button title="Toggle Heatmap" onclick="toggleHeatmap()"><i class="fas fa-fire fa-2x" style="padding-top: 3px"></i></button>
                &nbsp;
                <button title="Add Ruler" onclick="addRuler()" style="width: 38px; height: 36; padding-top: 4; padding-left: 0"><i class="fas fa-ruler fa-2x"></i></i></button>
              </p>
            </div>
            <div class="row">
              <p>
                <button title="Toggle Overlay" onclick="toggleOverlay()" style="width: 38px; padding-left: 3"><i class="fab fa-battle-net fa-2x" style="padding-top: 3px"></i></button>
                &nbsp;
                <button title="Edit Overlay" onclick="overlayEditable()" style="width: 38px; height: 36; padding-top: 4; padding-left: 3"><i class="fas fa-edit fa-2x"></i></button>
              </p>
            </div>
            <div class="row">
              <p>
                <button title="Toggle Drive" onclick="toggleDrive()" style="width: 38px; height: 36; padding-top: 4; padding-left: 1"><i class="fas fa-car-crash fa-2x"></i></button>
                &nbsp;
                <button title="Edit Drive" onclick="driveEditable()" style="width: 38px; height: 36; padding-top: 4; padding-left: 3"><i class="fas fa-edit fa-2x"></i></button>
              </p>
            </div>
        </div>
				<div class="col-xs-10">
          <div style="height: 85%" id="map"></div>
        </div>
      </div>
    </div>
  </body>
  <script type="text/javascript" src="https://maps.googleapis.com/maps/api/js?key={{.Apikey}}&libraries=visualization"></script>
  <script type="text/javascript" src="https://optiv.github.io/warmap-go/src/labels.js"></script>
<script>
var heatMapData = {{.Heatmap}};
var overlayCoords = {{.ConvexHull}};
var overlayDrive = {{.Drive}};

var map = new google.maps.Map(document.getElementById('map'), {
  zoom: 16,
  center: {lat: {{.Lat}}, lng: {{.Lng}}},
  mapTypeId: 'satellite',
	mapTypeControlOptions: {style: google.maps.MapTypeControlStyle.DROPDOWN_MENU},
	controlSize: 30,
	streetViewControl: false
});

var heatmap = new google.maps.visualization.HeatmapLayer({
  data: heatMapData
});

var convexHull =  new google.maps.Polygon({
					paths: overlayCoords,
					editable: false,
          strokeColor: '#3366FF',
          strokeOpacity: 0.8,
          strokeWeight: 2,
          fillColor: '#3366FF',
          fillOpacity: 0.35
        });

var drive =  new google.maps.Polyline({
		path: overlayDrive,
		editable: false,
		geodesic: true,
		strokeColor: '#3366FF',
		strokeOpacity: 1.0,
	});

function toggleHeatmap() {
    heatmap.setMap(heatmap.getMap() ? null : map);
}
function toggleOverlay() {
	convexHull.setMap(convexHull.getMap() ? null : map)
}
function toggleDrive() {
	drive.setMap(drive.getMap() ? null : map)
}
function driveEditable() {
	if (drive.editable) {
		drive.setEditable(false);
	}
	else {
		drive.setEditable(true);
	}
}
function overlayEditable() {
	if (convexHull.editable) {
		convexHull.setEditable(false);
	}
	else {
		convexHull.setEditable(true);
	}
}
var lines = new Array();
function addRuler() {
  var ruler1 = new google.maps.Marker({
    position: map.getCenter(),
    map: map,
    draggable: true
  });
  var ruler2 = new google.maps.Marker({
    position: map.getCenter(),
    map: map,
    draggable: true
  });
  var ruler1label = new Label({
    map: map
  });
  var ruler2label = new Label({
    map: map
  });
  ruler1label.bindTo('position', ruler1, 'position');
  ruler2label.bindTo('position', ruler2, 'position');
  var rulerpoly = new google.maps.Polyline({
    path: [ruler1.position, ruler2.position],
    strokeColor: "#FFFF00",
    strokeOpacity: .7,
    strokeWeight: 7
  });
  rulerpoly.setMap(map);
  ruler1label.set('text', distance(ruler1.getPosition().lat(), ruler1.getPosition().lng(), ruler2.getPosition().lat(), ruler2.getPosition().lng()));
  ruler2label.set('text', distance(ruler1.getPosition().lat(), ruler1.getPosition().lng(), ruler2.getPosition().lat(), ruler2.getPosition().lng()));
  google.maps.event.addListener(ruler1, 'drag', function() {
    rulerpoly.setPath([ruler1.getPosition(), ruler2.getPosition()]);
    ruler1label.set('text', distance(ruler1.getPosition().lat(), ruler1.getPosition().lng(), ruler2.getPosition().lat(), ruler2.getPosition().lng()));
    ruler2label.set('text', distance(ruler1.getPosition().lat(), ruler1.getPosition().lng(), ruler2.getPosition().lat(), ruler2.getPosition().lng()));
  });
  google.maps.event.addListener(ruler2, 'drag', function() {
    rulerpoly.setPath([ruler1.getPosition(), ruler2.getPosition()]);
    ruler1label.set('text', distance(ruler1.getPosition().lat(), ruler1.getPosition().lng(), ruler2.getPosition().lat(), ruler2.getPosition().lng()));
    ruler2label.set('text', distance(ruler1.getPosition().lat(), ruler1.getPosition().lng(), ruler2.getPosition().lat(), ruler2.getPosition().lng()));
  });

  google.maps.event.addListener(ruler1, 'dblclick', function() {
    ruler1.setMap(null);
    ruler2.setMap(null);
    ruler1label.setMap(null);
    ruler2label.setMap(null);
    rulerpoly.setMap(null);
  });

  google.maps.event.addListener(ruler2, 'dblclick', function() {
    ruler1.setMap(null);
    ruler2.setMap(null);
    ruler1label.setMap(null);
    ruler2label.setMap(null);
    rulerpoly.setMap(null);
  });

  // Add our new ruler to an array for later reference
  lines.push([ruler1, ruler2, ruler1label, ruler2label, rulerpoly]);
}

function distance(lat1, lon1, lat2, lon2) {
  var R = 3959; // Here's the right settings for miles and feet
  var dLat = (lat2 - lat1) * Math.PI / 180;
  var dLon = (lon2 - lon1) * Math.PI / 180;
  var a = Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
    Math.sin(dLon / 2) * Math.sin(dLon / 2);
  var c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  var d = R * c;
  if (d > 1) return Math.round(d) + "mi";
  else if (d <= 1) return Math.round(d * 5280) + "ft";
  return d;
}
</script>
</html>
`
)

//populateTemplate populates the html template
func populateTemplate(points *Points, apikey *string) []byte {
	// func populateTemplate(points *Points, apikey *string) []byte {
	var page Page
	var convexData string
	var heatmap string
	var driveData string
	high, low := 0, 0
	if apikey != nil {
		page.Apikey = apikey
	} else {
		page.Apikey = nil
	}
	var tplBuffer bytes.Buffer
	convexPoints := findConvexHull(*points)
	for _, point := range convexPoints {
		convexData += fmt.Sprintf("(new google.maps.LatLng(%g, %g)), ", point.Y, point.X)
	}
	for _, point := range *points {
		heatmap += fmt.Sprintf("{location: new google.maps.LatLng(%g, %g), weight: %f}, ", point.Y, point.X, (float64(point.Dbm)/10.0)+9.0)
		sigCheck(&high, &low, point.Dbm)
	}
	for _, point := range *points {
		driveData += fmt.Sprintf("(new google.maps.LatLng(%g, %g)), ", point.Y, point.X)
	}
	page.Lat = (*points)[0].Y
	page.Lng = (*points)[0].X
	page.PathLength = len(*points)
	page.ConvexHull = template.JS("[" + convexData[:len(convexData)-2] + "]")
	page.Heatmap = template.JS("[" + heatmap[:len(heatmap)-2] + "]")
	page.Drive = template.JS("[" + driveData[:len(driveData)-2] + "]")
	page.HighDB = template.JS(fmt.Sprintf("%v", high))
	page.LowDB = template.JS(fmt.Sprintf("%v", low))
	t, err := template.New("webpage").Parse(tpl)
	checkError(err)
	err = t.Execute(&tplBuffer, page)
	checkError(err)
	return tplBuffer.Bytes()
}
