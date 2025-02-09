package main

import (
	"fmt"

	"github.com/chrisatcho/external-dns-technitiumdns-webhook/cmd/webhook/init/configuration"
	"github.com/chrisatcho/external-dns-technitiumdns-webhook/cmd/webhook/init/dnsprovider"
	"github.com/chrisatcho/external-dns-technitiumdns-webhook/cmd/webhook/init/logging"
	"github.com/chrisatcho/external-dns-technitiumdns-webhook/cmd/webhook/init/server"
	"github.com/chrisatcho/external-dns-technitiumdns-webhook/pkg/webhook"
	log "github.com/sirupsen/logrus"
)

const banner = `
 external-dns-technitium-webhook
 version: %s (%s)

`

var (
	Version = "local"
	Gitsha  = "?"
)

func main() {
	fmt.Printf(banner, Version, Gitsha)
	logging.Init()
	config := configuration.Init()
	provider, err := dnsprovider.Init(config)
	if err != nil {
		log.Fatalf("Failed to initialize DNS provider: %v", err)
	}
	srv := server.Init(config, webhook.New(provider))
	server.ShutdownGracefully(srv)
}
