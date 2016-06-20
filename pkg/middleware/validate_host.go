package middleware

import (
	"strings"

	"github.com/ganting/grafana/pkg/setting"
	"gopkg.in/macaron.v1"
)

func ValidateHostHeader(domain string) macaron.Handler {
	return func(c *macaron.Context) {
		h := c.Req.Host
		if i := strings.Index(h, ":"); i >= 0 {
			h = h[:i]
		}

		if !strings.EqualFold(h, domain) {
			c.Redirect(strings.TrimSuffix(setting.AppUrl, "/")+c.Req.RequestURI, 301)
			return
		}
	}
}
