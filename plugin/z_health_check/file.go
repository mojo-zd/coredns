package zone_health_check

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/mholt/caddy"
)

var (
	directive = "z_health_check"
	filePath  string
)

type HealthMap map[string]*Health

// fileParse ...
func fileParse(c *caddy.Controller) (healthMap HealthMap, err error) {

	healthMap = map[string]*Health{}
	config := dnsserver.GetConfig(c)

	for c.Next() {
		if !c.NextArg() {
			err = c.ArgErr()
			return
		}
		fileName := c.Val()
		if directive == fileName {
			continue
		}

		if !path.IsAbs(fileName) && config.Root != "" {
			fileName = path.Join(config.Root, fileName)
		}
		filePath = fileName
		reader, err := os.Open(fileName)
		if err != nil {
			logrus.Errorf("filed open failed %s", err.Error())
			return nil, err
		}

		z, err := parse(reader)
		if err != nil {
			return nil, err
		}

		for k, v := range z {
			healthMap[plugin.Host(k).Normalize()] = v
		}
	}
	return
}

// parse parse to zHealth obj
func parse(f io.Reader) (z map[string]*Health, err error) {
	var bytes []byte

	if bytes, err = ioutil.ReadAll(f); err != nil {
		return
	}

	if err = json.Unmarshal(bytes, &z); err == nil {
		for key, zm := range z {
			if len(zm.Protocol) == 0 {
				zm.Protocol = "http://"
			}
			zm.Domain = key
		}
	}
	return
}
