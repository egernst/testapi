// Copyright 2018, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"net/http"
	"os"

	"go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

func main() {
	serviceName := "client"
	server := os.Getenv("SERVER_URI")
	if server == "" {
		server = "http://localhost:8888"
	}

	// Setup tracing
	// reporterURI: zipkin reporter URI
	reporterURI := os.Getenv("REPORTER_URI")
	if reporterURI == "" {
		reporterURI = "http://localhost:9411/api/v2/spans"
	}

	localEndpoint, err := openzipkin.NewEndpoint(serviceName, ":0")
	if err != nil {
		log.Fatalf("Failed to create Zipkin localEndpoint with URI %q error: %v", serviceName, err)
	}

	reporter := zipkinHTTP.NewReporter(reporterURI)
	ze := zipkin.NewExporter(reporter, localEndpoint)

	// And now finally register it as a Trace Exporter
	trace.RegisterExporter(ze)

	//TODO: Switch to trace.ProbabilitySampler if needed
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	client := &http.Client{Transport: &ochttp.Transport{}}

	for i := 0; i < 100000; i++ {
		resp, err := client.Get(server)
		if err != nil {
			log.Printf("Failed to get response: %v", err)
		} else {
			resp.Body.Close()
		}
	}

}
