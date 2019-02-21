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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"

	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	// Report latency seen by this server on the downstream call
	mLatencyUs = stats.Float64("latency", "The downstream latency in microseconds", "ms")
)

func main() {

	// e.g: http://localhost:8888/
	upstreamURI := os.Getenv("UPSTREAM_URI")
	if upstreamURI == "" {
		fmt.Println("Error: UPSTREAM_URI not present")
		os.Exit(1)
	}

	// URI of the downstream service
	// e.g: http://localhost:8889/
	downstreamURI := os.Getenv("DOWNSTREAM_URI")

	// Setup tracing
	// reporterURI: zipkin reporter URI
	reporterURI := os.Getenv("REPORTER_URI")
	if reporterURI == "" {
		reporterURI = "http://localhost:9411/api/v2/spans"
	}

	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "8887"
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		var err error
		if serviceName, err = os.Hostname(); err != nil {
			serviceName = "service"
		}
	}

	localEndpoint, err := openzipkin.NewEndpoint(serviceName, upstreamURI)
	if err != nil {
		log.Fatalf("Failed to create Zipkin localEndpoint with URI %q error: %v", upstreamURI, err)
	}

	reporter := zipkinHTTP.NewReporter(reporterURI)
	ze := zipkin.NewExporter(reporter, localEndpoint)

	// And now finally register it as a Trace Exporter
	trace.RegisterExporter(ze)

	//TODO: Switch to trace.ProbabilitySampler if needed
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// Setup metrics
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: serviceName,
	})
	if err != nil {
		log.Fatalf("Failed to create the Prometheus exporter: %v", err)
	}

	// register it as a stats exporter.
	view.RegisterExporter(pe)
	// Report stats at every second.
	view.SetReportingPeriod(1 * time.Second)

	latencyView := &view.View{
		Name:        serviceName + "/latency",
		Measure:     mLatencyUs,
		Description: "The distribution of the latencies",
		Aggregation: view.Distribution(0, 100, 500, 600, 700, 800, 900, 1000),
	}

	/*
		if err := view.Register(ochttp.DefaultServerViews...); err != nil {
			log.Fatalf("Failed to register metrics view: %v", err)
		}
	*/

	if err := view.Register(latencyView); err != nil {
		log.Fatalf("Failed to register metrics view: %v", err)
	}

	//Start the metrics server on the user specified port
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", pe)
		if err := http.ListenAndServe(":"+metricsPort, mux); err != nil {
			log.Fatalf("Failed to run Prometheus /metrics endpoint: %v", err)
		}
	}()

	// Run the actual workload
	client := &http.Client{Transport: &ochttp.Transport{}}
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		fmt.Fprintf(w, "hello:%s", serviceName)

		if downstreamURI != "" {
			_, span := trace.StartSpan(req.Context(), "downstream")
			defer span.End()

			span.Annotate([]trace.Attribute{trace.StringAttribute("key", "value")}, "something happened")
			span.AddAttributes(trace.StringAttribute("hello", serviceName))

			r, _ := http.NewRequest("GET", downstreamURI, nil)

			// Propagate the trace header info in the outgoing requests.
			r = r.WithContext(req.Context())

			startTime := time.Now()
			resp, err := client.Do(r)
			t := float64(time.Since(startTime).Nanoseconds()) / 1e3
			stats.Record(req.Context(), mLatencyUs.M(t))

			if err != nil {
				log.Println(err)
			} else {
				if body, err := ioutil.ReadAll(resp.Body); err != nil {
					fmt.Fprintf(w, ":%s", string(body))
				}
				resp.Body.Close()
			}
		} else {
			//This is time to work
			t := float64(time.Since(startTime).Nanoseconds()) / 1e3
			stats.Record(req.Context(), mLatencyUs.M(t))
		}
	})
	log.Fatal("Server", http.ListenAndServe(upstreamURI, &ochttp.Handler{}))

}
