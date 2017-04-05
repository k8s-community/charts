package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/takama/charts/version"
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
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", chartsPort), http.FileServer(http.Dir(servicePath)))
}
