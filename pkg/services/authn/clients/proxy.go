package clients

import (
	"context"
	"fmt"
	"net"
	"path"
	"strings"

	"github.com/grafana/grafana/pkg/services/authn"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
	"github.com/grafana/grafana/pkg/util/errutil"
)

var (
	errNotAcceptedIP = errutil.NewBase(errutil.StatusValidationFailed, "auth-proxy.invalid-ip")
)

var _ authn.Client = new(Proxy)

func ProvideProxy(cfg *setting.Cfg, clients ...authn.ProxyClient) (*Proxy, error) {
	list, err := parseAcceptList(cfg.AuthProxyWhitelist)
	if err != nil {
		return nil, err
	}
	return &Proxy{cfg, clients, list}, nil
}

type Proxy struct {
	cfg         *setting.Cfg
	clients     []authn.ProxyClient
	acceptedIPs []*net.IPNet
}

func (c *Proxy) Authenticate(ctx context.Context, r *authn.Request) (*authn.Identity, error) {
	if !c.isAllowedIP(r) {
		return nil, errNotAcceptedIP.Errorf("request ip is not in the configured accept list")
	}

	username := c.getHeader(r)

	var clientErr error
	for _, proxyClient := range c.clients {
		var identity *authn.Identity
		identity, clientErr = proxyClient.AuthenticateProxy(ctx, r, username)
		if identity != nil {
			return identity, nil
		}
	}

	return nil, clientErr
}

func (c *Proxy) Test(ctx context.Context, r *authn.Request) bool {
	return len(c.getHeader(r)) != 0
}

func (p *Proxy) getHeader(r *authn.Request) string {
	if r.HTTPRequest == nil {
		return ""
	}
	v := r.HTTPRequest.Header.Get(p.cfg.AuthProxyHeaderName)
	if p.cfg.AuthProxyHeadersEncoded {
		v = util.DecodeQuotedPrintable(v)
	}
	return v
}

func (p *Proxy) isAllowedIP(r *authn.Request) bool {
	if len(p.acceptedIPs) == 0 {
		return true
	}

	host, _, err := net.SplitHostPort(r.HTTPRequest.RemoteAddr)
	if err != nil {
		return false
	}

	ip := net.ParseIP(host)
	for _, v := range p.acceptedIPs {
		if v.Contains(ip) {
			return true
		}
	}

	return false
}

func parseAcceptList(s string) ([]*net.IPNet, error) {
	addresses := strings.Split(s, ",")
	list := make([]*net.IPNet, 0, len(addresses))
	for _, addr := range addresses {
		result, err := coerceProxyAddress(addr)
		if err != nil {
			return nil, err
		}
		list = append(list, result)
	}
	return list, nil
}

// coerceProxyAddress gets network of the presented CIDR notation
func coerceProxyAddress(proxyAddr string) (*net.IPNet, error) {
	proxyAddr = strings.TrimSpace(proxyAddr)
	if !strings.Contains(proxyAddr, "/") {
		proxyAddr = path.Join(proxyAddr, "32")
	}

	_, network, err := net.ParseCIDR(proxyAddr)
	if err != nil {
		return nil, fmt.Errorf("could not parse the network: %w", err)
	}
	return network, nil
}
