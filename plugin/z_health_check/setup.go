package zone_health_check

import (
	"fmt"
	"strings"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/z_health_check/job"
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("z_health_check", caddy.Plugin{ServerType: "dns", Action: setup})
	worker = job.NewWorker()
}

func setup(c *caddy.Controller) (err error) {
	loaded, err = fileParse(c)
	if err != nil {
		return
	}

	for _, v := range loaded {
		api := v.Api
		if !strings.HasPrefix(api, "/") {
			api = "/" + api
		}

		for _, ip := range v.IPs {
			worker.AddJob(fmt.Sprintf("%s%s%s", v.Protocol, ip, api), v.Domain, ip)
		}
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		z := ZHealth{Next: next}
		return z
	})

	return
}
