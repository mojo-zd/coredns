package zone_health_check

import (
	"context"
	"net"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/z_health_check/job"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

var (
	loaded map[string]*Health
	worker *job.Worker
)

type ZHealth struct {
	Next plugin.Handler
}

type Health struct {
	Domain   string
	IPs      []string
	Api      string
	Protocol string
	Timeout  int
}

func (z ZHealth) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	var zones []string
	var answers []dns.RR
	for k := range job.HealtherMap {
		zones = append(zones, plugin.Host(k).Normalize())
	}
	qname := state.QName()
	zone := plugin.Zones(zones).Matches(qname)
	if zone == "" {
		// PTR zones don't need to be specified in Origins
		if state.Type() != "PTR" {
			// If this doesn't match we need to fall through regardless of h.Fallthrough
			return plugin.NextOrFailure(z.Name(), z.Next, ctx, w, r)
		}
	}

	switch state.QType() {
	case dns.TypePTR:
		//names := z.LookupStaticAddr(dnsutil.ExtractAddressFromReverse(qname))
		//if len(names) == 0 {
		//	// If this doesn't match we need to fall through regardless of h.Fallthrough
		//	return plugin.NextOrFailure(z.Name(), z.Next, ctx, w, r)
		//}
		//answers = z.ptr(qname, names)
	case dns.TypeA:
		ips := z.LookupStaticHostVx(qname, 4)
		answers = a(qname, ips)
	case dns.TypeAAAA:
		ips := z.LookupStaticHostVx(qname, 6)
		answers = a(qname, ips)
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative, m.RecursionAvailable = true, true
	m.Answer = answers

	state.SizeAndDo(m)
	m, _ = state.Scrub(m)
	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

// Name implements the plugin.Handle interface.
func (z ZHealth) Name() string { return "z_health_check" }

func (z ZHealth) LookupStaticHostVx(qname string, version int) (ips []net.IP) {
	for addr, values := range job.HealtherMap {
		if plugin.Host(addr).Normalize() == qname {
			for _, ip := range values {
				ver := ipVersion(ip)
				if ver == version {
					ips = append(ips, parseLiteralIp(ip))
				}
			}
		}
	}
	return
}

// a takes aslice ofnet.IPs and returns a slice of A RRs.
func a(zone string, ips []net.IP) []dns.RR {
	var answers []dns.RR
	for _, ip := range ips {
		r := new(dns.A)
		r.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600}
		r.A = ip
		answers = append(answers, r)
	}
	return answers
}

// aaaa takes a slice of net.IPs and returns a slice of AAAA RRs.
func aaaa(zone string, ips []net.IP) []dns.RR {
	var answers []dns.RR
	for _, ip := range ips {
		r := new(dns.AAAA)
		r.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypeAAAA,
			Class: dns.ClassINET, Ttl: 3600}
		r.AAAA = ip
		answers = append(answers, r)
	}
	return answers
}

// parseLiteralIp parse string to ip
func parseLiteralIp(addr string) net.IP {
	if i := strings.Index(addr, "%"); i > 0 {
		addr = addr[0:i]
	}
	return net.ParseIP(addr)
}

// LookupStaticAddr ...
func (z ZHealth) LookupStaticAddr(addr string) []string {
	net.ParseIP(addr)
	addr = parseLiteralIP(addr).String()
	if addr == "" {
		return nil
	}

	return nil
}

func parseLiteralIP(addr string) net.IP {
	if i := strings.Index(addr, "%"); i >= 0 {
		// discard ipv6 zone
		addr = addr[0:i]
	}

	return net.ParseIP(addr)
}

// ipVersion returns what IP version was used textually
func ipVersion(s string) int {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return 4
		case ':':
			return 6
		}
	}
	return 0
}
