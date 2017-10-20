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
	"math/rand"

	"github.com/christophberger/grada"
)

// ## The data generator
//
// This function creates a stream of data that resembles a graph
// similar to stock data.
func newFakeDataFunc(max int, volatility float64) func() float64 {
	value := rand.Float64()
	return func() float64 {
		rnd := 2 * (rand.Float64() - 0.5)
		change := volatility * rnd
		change += (0.5 - value) * 0.1
		value += change
		return value * float64(max)
	}
}

func main() {
	server := grada.StartServer()
	server.Metrics.CreateMetric("GOGL", 100)
	server.Metrics.CreateMetric("AAPL", 100)

	gogl := newFakeDataFunc(100, 0.2)
	aapl := newFakeDataFunc(100, 0.1)

	go func() {
		for {
			(*server.Metrics)["GOGL"].Add(gogl())
		}
	}()
	for {
		(*server.Metrics)["AAPL"].Add(aapl())
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
