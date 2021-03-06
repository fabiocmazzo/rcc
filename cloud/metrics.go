package cloud

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/xviper"
)

var (
	telemetryBarrier = sync.WaitGroup{}
)

const (
	trackingUrl = `/metric-v1/%v/%v/%v/%v/%v`
	metricsHost = `https://telemetry.robocorp.com`
)

func sendMetric(kind, name, value string) {
	common.Timeline("%s:%s = %s", kind, name, value)
	defer func() {
		status := recover()
		if status != nil {
			common.Debug("Telemetry panic recovered: %v", status)
		}
		telemetryBarrier.Done()
	}()
	client, err := NewClient(metricsHost)
	if err != nil {
		common.Debug("ERROR: %v", err)
		return
	}
	timestamp := time.Now().UnixNano()
	url := fmt.Sprintf(trackingUrl, url.PathEscape(kind), timestamp, url.PathEscape(xviper.TrackingIdentity()), url.PathEscape(name), url.PathEscape(value))
	common.Debug("DEBUG: Sending metric as %v%v", metricsHost, url)
	client.Put(client.NewRequest(url))
}

func BackgroundMetric(kind, name, value string) {
	common.Debug("DEBUG: BackgroundMetric kind:%v name:%v value:%v send:%v", kind, name, value, xviper.CanTrack())
	if xviper.CanTrack() {
		telemetryBarrier.Add(1)
		go sendMetric(kind, name, value)
	}
}

func WaitTelemetry() {
	common.Debug("DEBUG: wait telemetry to complete")
	telemetryBarrier.Wait()
	common.Debug("DEBUG: telemetry sending completed")
}
