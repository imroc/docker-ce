package main

import (
	"fmt"
	"net"
	"net/http"

	metrics "github.com/docker/go-metrics"
	"github.com/sirupsen/logrus"
)

func startMetricsServer(hasExperimental bool, addr string) error {
	if addr == "" {
		return nil
	}
	if !hasExperimental {
		return fmt.Errorf("metrics-addr is only supported when experimental is enabled")
	}
	if err := allocateDaemonPort(addr); err != nil {
		return err
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/metrics", metrics.Handler())
	go func() {
		if err := http.Serve(l, mux); err != nil {
			logrus.Errorf("serve metrics api: %s", err)
		}
	}()
	return nil
}
