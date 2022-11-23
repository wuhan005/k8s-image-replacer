package main

import (
	"flag"
	"net/http"

	"github.com/flamego/flamego"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"

	"github.com/wuhan005/k8s-image-replacer/internal/conf"
	"github.com/wuhan005/k8s-image-replacer/internal/webhook"
)

func main() {
	configFile := flag.String("config", "config.yaml", "configuration file")
	flag.Parse()

	if err := conf.Init(*configFile); err != nil {
		logrus.WithError(err).Fatal("Failed to init config")
	}

	wh, err := webhook.NewHandler()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create webhook handler")
	}
	whHandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{Webhook: wh})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create webhook HTTP handler")
	}

	f := flamego.Classic()
	f.Any("/webhook", func(ctx flamego.Context) {
		whHandler.ServeHTTP(ctx.ResponseWriter(), ctx.Request().Request)
	})
	f.Any("/metrics", promhttp.Handler())
	f.Get("/healthz")

	tlsCertFile := conf.ImageReplacer.TlsCrtFile
	if tlsCertFile == "" {
		tlsCertFile = "crt/tls.crt"
	}
	tlsKeyFile := conf.ImageReplacer.TlsKeyFile
	if tlsKeyFile == "" {
		tlsKeyFile = "crt/tls.key"
	}

	if err := http.ListenAndServeTLS(":443", tlsCertFile, tlsKeyFile, f); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}
}
