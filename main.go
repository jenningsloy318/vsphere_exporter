package main

import (
	"context"
	"github.com/jenningsloy318/vsphere_exporter/collector"
	"github.com/jenningsloy318/vsphere_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	configFile = kingpin.Flag(
		"config.file",
		"Path to configuration file.",
	).String()
	listenAddress = kingpin.Flag(
		"web.listen-address",
		"Address to listen on for web interface and telemetry.",
	).Default(":9272").String()
	sc = &config.SafeConfig{
		C: &config.Config{},
	}
	reloadCh chan chan error
)

// define new http handleer
func metricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()
		var target string
		var ctx context.Context
		var clusterConfig *config.ClusterConfig
		var err error
		if sc.C.Mode == "single" {
			target = sc.C.EnabledCluster
			if clusterConfig, err = sc.SetSingleModeClusterCredential(); err != nil {
				log.Errorf("Error getting credential %s", target, err)
				return
			}
		} else {
			target = r.URL.Query().Get("target")
			if target == "" {
				http.Error(w, "'target' parameter must be specified in multi scrape mode", 400)
				return
			}

			if clusterConfig, err = sc.ClusterConfigForTarget(target); err != nil {
				log.Errorf("Error getting credential for target %s,%s", target, err)
				return
			}
			ctx = r.Context()

		}
		if ctx == nil {
			ctx = context.Background()
		}
		log.Infof("starting scraping target %s", target)
		collector := collector.NewVshpereCollector(ctx, target, clusterConfig.Username, clusterConfig.Password)
		registry.MustRegister(collector)
		gatherers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			registry,
		}
		// Delegate http serving to Prometheus client library, which will call collector.Collect.
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)

	}
}

func main() {
	log.AddFlags(kingpin.CommandLine)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	log.Infoln("Starting vsphere_exporter")
	// load config  first time
	if err := sc.ReloadConfig(*configFile); err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}

	// load config in background to wathc config changes
	hup := make(chan os.Signal)
	reloadCh = make(chan chan error)
	signal.Notify(hup, syscall.SIGHUP)

	go func() {
		for {
			select {
			case <-hup:
				if err := sc.ReloadConfig(*configFile); err != nil {
					log.Errorf("Error reloading config: %s", err)
				}
			case rc := <-reloadCh:
				if err := sc.ReloadConfig(*configFile); err != nil {
					log.Errorf("Error reloading config: %s", err)
					rc <- err
				} else {
					rc <- nil
				}
			}
		}
	}()

	http.Handle("/vsphere", metricsHandler()) // Regular metrics endpoint for local vsphere metrics.
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head>
            <title>vSphere Exporter</title>
            </head>
			<body>
            <h1>vSphere Exporter</h1>
            <form action="/vsphere">
            <label>Target:</label> <input type="text" name="target" placeholder="xxxx" value="IP or dns name"><br>
            <input type="submit" value="Submit">
			</form>
			<p><a href="/metrics">Local metrics</a></p>
            </body>
            </html>`))
	})

	log.Infof("Listening on %s", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
