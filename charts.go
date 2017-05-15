package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/k8s-community/charts/version"
	"github.com/takama/router"
)

const (
	// DEFAULTHEALTHPORT returns default health port number
	DEFAULTHEALTHPORT = "8082"
	// DEFAULTCHARTSPORT returns default port number
	DEFAULTCHARTSPORT = "8080"
	// DEFAULTPATH contains path to charts repo
	DEFAULTPATH = "/var/lib/charts"
)

// simplest logger, which initialized during starts of the application
var (
	stdlog = log.New(os.Stdout, "[CHARTS]: ", log.LstdFlags)
	errlog = log.New(os.Stderr, "[CHARTS:ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func shutdown() (string, error) {
	return "Shutdown", nil
}

func logger(c *router.Control) {
	remoteAddr := c.Request.Header.Get("X-Forwarded-For")
	if remoteAddr == "" {
		remoteAddr = c.Request.RemoteAddr
	}
	stdlog.Println(remoteAddr, c.Request.Method, c.Request.URL.Path)
}

func healthz(c *router.Control) {
	c.Code(http.StatusOK).Body("Ok")
}

func info(c *router.Control) {
	c.Code(http.StatusOK).Body(
		map[string]string{
			"version": version.RELEASE,
			"commit":  version.COMMIT,
			"repo":    version.REPO,
		},
	)
}

func root(c *router.Control) {
	c.Code(http.StatusOK).Body(fmt.Sprintf("Helm charts v%s", version.RELEASE))
}

func main() {
	chartsPort := os.Getenv("CHARTS_SERVICE_PORT")
	if len(chartsPort) == 0 {
		chartsPort = DEFAULTCHARTSPORT
	}
	healthPort := os.Getenv("CHARTS_SERVICE_HEALTH_PORT")
	if len(healthPort) == 0 {
		healthPort = DEFAULTHEALTHPORT
	}
	servicePath := os.Getenv("CHARTS_SERVICE_PATH")
	if len(servicePath) == 0 {
		servicePath = DEFAULTPATH
	}
	r := router.New()
	r.Logger = logger
	r.GET("/info", info)
	r.GET("/healthz", healthz)
	r.GET("/", root)
	go r.Listen(fmt.Sprintf("0.0.0.0:%s", healthPort))
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", chartsPort), http.FileServer(http.Dir(servicePath)))

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	killSignal := <-interrupt
	stdlog.Println("Got signal:", killSignal)
	status, err := shutdown()
	if err != nil {
		errlog.Printf("Error: %s Status: %s\n", err.Error(), status)
		os.Exit(1)
	}
	if killSignal == os.Kill {
		stdlog.Println("Service was killed")
	} else {
		stdlog.Println("Service was terminated by system signal")
	}
	stdlog.Println(status)
}
