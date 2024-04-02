package cloudflare

import (
	"cloudflare-policy-manager/config"

	"github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
)

type CloudflareCli struct {
	Api *cloudflare.API
}

// NewCloudflareCli creates a new CloudflareCli
func NewCloudflareCli(cfg config.CloudflareConfig) *CloudflareCli {
	api, err := cloudflare.New(cfg.CloudflareApiKey, cfg.CloudflareEmail)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create Cloudflare API client")
	}

	return &CloudflareCli{Api: api}
}
