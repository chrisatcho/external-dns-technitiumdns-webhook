package dnsprovider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v8"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/provider"

	"github.com/chrisatcho/external-dns-technitiumdns-webhook/cmd/webhook/init/configuration"
	"github.com/chrisatcho/external-dns-technitiumdns-webhook/internal/technitium"
)

func Init(config configuration.Config) (provider.Provider, error) {
	var domainFilter endpoint.DomainFilter
	var filters []string

	if config.RegexDomainFilter != "" {
		filters = append(filters, fmt.Sprintf("regexp domain filter: '%s'", config.RegexDomainFilter))
		if config.RegexDomainExclusion != "" {
			filters = append(filters, fmt.Sprintf("regexp exclusion: '%s'", config.RegexDomainExclusion))
		}
		domainFilter = endpoint.NewRegexDomainFilter(
			regexp.MustCompile(config.RegexDomainFilter),
			regexp.MustCompile(config.RegexDomainExclusion),
		)
	} else {
		if len(config.DomainFilter) > 0 {
			filters = append(filters, fmt.Sprintf("domain filter: '%s'", strings.Join(config.DomainFilter, ",")))
		}
		if len(config.ExcludeDomains) > 0 {
			filters = append(filters, fmt.Sprintf("exclude domain filter: '%s'", strings.Join(config.ExcludeDomains, ",")))
		}
		domainFilter = endpoint.NewDomainFilterWithExclusions(config.DomainFilter, config.ExcludeDomains)
	}

	if len(filters) == 0 {
		log.Info("creating TechnitiumDNS provider with no domain filters")
	} else {
		log.Infof("creating TechnitiumDNS provider with %s", strings.Join(filters, ", "))
	}

	technitiumConfig := technitium.Configuration{}
	if err := env.Parse(&technitiumConfig); err != nil {
		return nil, fmt.Errorf("reading technitium configuration failed: %w", err)
	}

	return technitium.NewProvider(domainFilter, &technitiumConfig), nil
}
