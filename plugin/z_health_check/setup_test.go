package zone_health_check

import (
	"testing"
	"time"

	"github.com/astaxie/beego/utils"
	"github.com/mholt/caddy"
	"github.com/mojo-zd/zone_health_check/job"
)

var (
	z_health_check = `z_health_check
z_health.org`

	test_block = `registry.com {
    z_health_check z_health.org
    file demo.registry.org
    log
}`
)

func TestSetup(t *testing.T) {
	c := caddy.NewTestController("dns", z_health_check)
	setup(c)
	time.Sleep(time.Second * 6)
	utils.Display("", job.HealtherMap)
	time.Sleep(time.Second * 5)
	utils.Display("second", job.HealtherMap)
	time.Sleep(time.Minute)
}
