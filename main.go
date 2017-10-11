/*
<!--
Copyright (c) 2017 Christoph Berger. Some rights reserved.

Use of the text in this file is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

Use of the code in this file is governed by a BSD 3-clause license that can be found
in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = ""
description = ""
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2017-00-00"
draft = "true"
domains = [""]
tags = ["", "", ""]
categories = ["Tutorial"]
+++

### Summary goes here

<!--more-->

## Intro goes here

<!-- material

* Install and run Grafana

-> Using Docker as this can be done on all OSes

Step 1: Create volume

	docker create volume grafana-storage

Step 2: Run Grafana
* Attach volume
* Install plugin "grafana-simple-json-datasource"

    docker run -d -p 3000:3000 --name grafana --mount src=grafana-storage,dst=/var/lib/grafana --network one -e "GF_INSTALL_PLUGINS=grafana-simple-json-datasource" grafana/grafana

Step 3: Login

admin/admin by default

Step 4: Create the datasource

Menu > Data Sources > Add datasource

Name: My DIY Datasource
Default: Check
Type: Simplejson
URL: http://grafanago:3001
URL: http://docker.for.mac.localhost:3001
Access: proxy

Click Add

-> A warning may appear as the backend is not running yet.


Step 5: Create Dashboard

Click "Create your first dashboard"

Add a Singlestat panel.

Metrics > Data Source: Keep the "default" selection.



Save as "My DIY Dashboard"

Run the Go code

env GOOS=linux GOARCH=amd64 go build .; and docker build -t grafanago .; and docker run --rm -it -p 3001:3001 --network one --name grafanago grafanago




-->

## The code
*/

// ## Imports and globals
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Query is a `/query` request from Grafana.
//
// All JSON-related structs were generated from the JSON examples
// of the "SimpleJson" data source documentation
// using [JSON-to-Go](https://mholt.github.io/json-to-go/).
type Query struct {
	PanelID int `json:"panelId"`
	Range   struct {
		From time.Time `json:"from"`
		To   time.Time `json:"to"`
		Raw  struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"raw"`
	} `json:"range"`
	RangeRaw struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"rangeRaw"`
	Interval   string `json:"interval"`
	IntervalMs int    `json:"intervalMs"`
	Targets    []struct {
		Target string `json:"target"`
		RefID  string `json:"refId"`
		Type   string `json:"type"`
	} `json:"targets"`
	Format        string `json:"format"`
	MaxDataPoints int    `json:"maxDataPoints"`
}

// TimeseriesResponse is the response to a `/query` request
// if "Type" is set to "timeserie".
// It sends time series data back to Grafana.
type TimeseriesResponse struct {
	Target     string      `json:"target"`
	Datapoints [][]float64 `json:"datapoints"`
}

// TableResponse is the response to send when "Type" is "table".
type Column struct {
	Text string `json:"text"`
	Type string `json:"type"`
}
type Row []interface{}
type TableResponse struct {
	Columns []Column `json:"columns"`
	Rows    []Row    `json:"rows"`
	Type    string   `json:"type"`
}

// ## The server

func writeError(w http.ResponseWriter, e error, m string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"error\": \"" + m + ": " + e.Error() + "\"}"))

}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	var q bytes.Buffer

	n, err := q.ReadFrom(r.Body)
	if err != nil {
		writeError(w, err, "Cannot read request body")
		return
	}

	log.Printf("Request body (%d bytes read): %s", n, string(q.Bytes()))

	query := &Query{}
	err = json.Unmarshal(q.Bytes(), query)
	if err != nil {
		writeError(w, err, "cannot unmarshal request body")
		return
	}

	// Our example should contain exactly one target.
	target := query.Targets[0].Target

	log.Println("Sending response for target " + target)

	// Depending on the type, we need to send either a timeseries response
	// or a table response.
	switch query.Targets[0].Type {
	case "timeserie":
		sendTimeseries(w, query)
	case "table":
		sendTable(w, query)
	}
}

func sendTimeseries(w http.ResponseWriter, q *Query) {

	log.Println("Sending time series data")

	// from := q.Range.From
	// to := q.Range.To
	// interval := q.IntervalMs

	response := []TimeseriesResponse{
		{
			Target: q.Targets[0].Target,
			Datapoints: [][]float64{
				[]float64{68.0, float64(int64(time.Now().UnixNano() / 1000000))},
				[]float64{49.0, float64(int64(time.Now().UnixNano() / 1000000))},
				[]float64{2.0, float64(int64(time.Now().UnixNano() / 1000000))},
				[]float64{11.0, float64(int64(time.Now().UnixNano() / 1000000))},
			},
		},
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		writeError(w, err, "cannot marshal response")
	}

	log.Println("Response: ", string(jsonResp))

	w.Write(jsonResp)

}

func sendTable(w http.ResponseWriter, q *Query) {

	log.Println("Sending table data")

	// from := q.Range.From
	// to := q.Range.To
	// interval := q.IntervalMs

	response := []TableResponse{
		{
			Columns: []Column{
				{Text: "Name", Type: "string"},
				{Text: "Value", Type: "number"},
				{Text: "Time", Type: "time"},
			},
			Rows: []Row{
				{"Alpha", 68, float64(int64(time.Now().UnixNano() / 1000000))},
				{"Bravo", 49, float64(int64(time.Now().UnixNano() / 1000000))},
				{"Charlie", 2, float64(int64(time.Now().UnixNano() / 1000000))},
				{"Delta", 11, float64(int64(time.Now().UnixNano() / 1000000))},
			},
			Type: "table",
		},
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		writeError(w, err, "cannot marshal response")
	}

	log.Println("Response: ", string(jsonResp))

	w.Write(jsonResp)

}

func main() {
	// Grafana expects a "200 OK" status for "/" when testing the connection.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received: ", r.URL.Path)
		log.Println("Request body: ")
		io.Copy(os.Stderr, r.Body)

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/query", queryHandler)

	// Start the server.
	log.Println("start grafanago")
	defer log.Println("stop grafanago")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

/*
## How to get and run the code

Step 1: `go get` the code. Note the `-d` flag that prevents auto-installing
the binary into `$GOPATH/bin`.

    go get -d github.com/appliedgo/TODO:

Step 2: `cd` to the source code directory.

    cd $GOPATH/src/github.com/appliedgo/TODO:

Step 3. Run the binary.

    go run TODO:.go


## Odds and ends
## Some remarks
## Tips
## Links


**Happy coding!**

*/
