package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Version number constant.
const Version = "0.0.8"

// Homepage url.
const Homepage = "https://github.com/tristanmorgan/https-echo"

var (
	httpAddr  = flag.String("listen", ":80", "Listen address")
	destPort  = flag.Int("port", -1, "Destination port")
	versDisp  = flag.Bool("version", false, "Display version")
	stsEnable = flag.Bool("sts", true, "Strict-Transport-Security header enable")
)

func redirect(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Server", "Https-echo/"+Version+" (+"+Homepage+")")
	if *stsEnable {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000")
	}
	hostname := strings.Split(req.Host, ":")
	dps := ""
	if *destPort > 0 {
		dps = fmt.Sprintf(":%d", *destPort)
	}
	target := "https://" + hostname[0] + dps + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		log.Printf("redirect to: %s from: %s xff %s", target, req.RemoteAddr, xff)
	} else {
		log.Printf("redirect to: %s from: %s", target, req.RemoteAddr)
	}
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func health(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Server", "Https-echo/"+Version+" (+"+Homepage+")")
	io.WriteString(w, "Healthy.\n")
}

func main() {
	flag.Parse()

	prometheus.MustRegister(version.NewCollector("https-echo"))
	if *versDisp {
		fmt.Printf("Version: v%s %s\n", Version, runtime.Version())
		fmt.Printf("Home Page: %s\n", Homepage)
		os.Exit(0)
	}

	log.Printf("Listening for incoming requests on TCP port '%s'...", *httpAddr)
	http.HandleFunc("/", redirect)
	http.HandleFunc("/health", health)
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
