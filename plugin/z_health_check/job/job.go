package job

import (
	"net/http"
	"sync"
	"time"

	"github.com/mojo-zd/zone_health_check/httplib"
)

var (
	DefaultHealthCheckInterval = time.Second * 5
	locker                     sync.Mutex
)

type HealthCheckJob struct {
	URL    string
	domain string
	ip     string
	worker *Worker
}

func (j HealthCheckJob) HealthCheck() {
	go func() {
		if r := j.check(); r {
			add(j.domain, j.ip)
		} else {
			delete(j.domain, j.ip)
		}

		ticker := time.NewTicker(DefaultHealthCheckInterval)
		for {
			select {
			case <-ticker.C:
				if r := j.check(); r {
					add(j.domain, j.ip)
				} else {
					delete(j.domain, j.ip)
				}
			}
		}
	}()
}

func (j HealthCheckJob) check() (r bool) {
	r = true
	response, err := httplib.NewHttpRestTemplate().Get(j.URL)
	if err != nil || response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		r = false
		j.worker.stop <- j.URL
	}
	return
}

func add(domain, ip string) {
	locker.Lock()
	defer locker.Unlock()
	if _, ok := HealtherMap[domain]; !ok {
		HealtherMap[domain] = []string{ip}
		return
	}
	exist := false
	for _, aip := range HealtherMap[domain] {
		if aip == ip {
			exist = true
		}
	}

	if exist {
		return
	}
	HealtherMap[domain] = append(HealtherMap[domain], ip)
}

func delete(domain, ip string) {
	locker.Lock()
	defer locker.Unlock()
	for k, ips := range HealtherMap {
		if domain != k {
			continue
		}

		for i, aip := range ips {
			if ip == aip {
				HealtherMap[k] = append(ips[0:i], ips[i+1:]...)
			}
		}
	}
}
