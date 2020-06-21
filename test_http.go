package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"xsky-demon/errors"
)

var defaultRoundTripper http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	TLSHandshakeTimeout: 10 * time.Second,
}

type queryMap map[string]string

var (
	metricDeleteAPITimeout = 2 * time.Minute
	metricDeleteAdminPath  = "/api/v1/admin/tsdb/delete_series"
	metricPort             = 9090
	namespace              = "xms"
	volumeSubsystem        = "volume"
	metricUser             = "XA2uEHQLSO"
	metricPassword         = "fMlIANshhrQB8O6IFifO"
)

// FromTime returns a new millisecond timestamp from a time.
func FromTime(t time.Time) int64 {
	return t.Unix()*1000 + int64(t.Nanosecond())/int64(time.Millisecond)
}

func main() {
	log.Println(FromTime(time.Time{}))
	log.Println(time.Time{}.Format(time.RFC3339Nano))
	s := time.Time{}.Format(time.RFC3339Nano)
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		log.Println(t)
	}
	// roundTripper := &http.Transport{
	// 	Proxy: http.ProxyFromEnvironment,
	// 	DialContext: (&net.Dialer{
	// 		Timeout:   30 * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 	}).DialContext,
	// 	TLSHandshakeTimeout: 10 * time.Second,
	// }
	// roundTripper.RoundTrip(req)
	client := http.Client{Transport: defaultRoundTripper}
	manager := &prometheusManager{client: client}
	volumeID := 1
	host := "10.255.101.74"
	metricQueries := toQueryMap(namespace, volumeSubsystem)
	manager.deleteMetrics(metricQueries, host, int64(volumeID))
}

func joinHostPort(host string, port int64) string {
	return net.JoinHostPort(host, strconv.FormatInt(port, 10))
}

func deleteURL(host string) (*url.URL, error) {
	return url.Parse(fmt.Sprintf("http://%s"+metricDeleteAdminPath,
		joinHostPort(host, int64(metricPort))))
}

type prometheusManager struct {
	client http.Client
}

func (p *prometheusManager) deleteMetric(host string, query url.Values) (err error) {
	u, err := deleteURL(host)
	if err != nil {
		return errors.Trace(err)
	}
	u.RawQuery = query.Encode()

	log.Printf("Delete series: %s", u.String())

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return errors.Trace(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(metricUser, metricPassword)

	ctx, cancel := context.WithTimeout(context.Background(), metricDeleteAPITimeout)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := p.client.Do(req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return errors.Trace(err)
	}

	log.Printf("delete series resp: %#v", resp)
	return nil
}

func (p *prometheusManager) deleteMetrics(metricQueries queryMap, host string, id int64) (err error) {
	query := url.Values{}
	for _, metric := range metricQueries {
		matcher := fmt.Sprintf(`%s{id="%d"}`, metric, id)
		query.Add("match[]", matcher)
	}
	if err = p.deleteMetric(host, query); err != nil {
		log.Printf("Failed to delete stat of id %d: %s, query: %v", id, err, query)
	}
	return nil
}

func toQueryMap(namespace, subsystem string) (metricQueries queryMap) {
	metricQueries = make(queryMap, len(metricMap))
	for metric := range metricMap {
		metricQueries[metric] = fmt.Sprintf("%s_%s_%s", namespace, subsystem, metric)
	}
	return
}

var metricMap = map[string]struct{}{
	"read_iops":                     struct{}{},
	"read_bandwidth_kbyte":          struct{}{},
	"read_latency_us":               struct{}{},
	"read_wait_us":                  struct{}{},
	"write_iops":                    struct{}{},
	"write_bandwidth_kbyte":         struct{}{},
	"write_latency_us":              struct{}{},
	"write_wait_us":                 struct{}{},
	"total_iops":                    struct{}{},
	"total_bandwidth_kbyte":         struct{}{},
	"queue_depth":                   struct{}{},
	"io_size_0_4kbyte":              struct{}{},
	"io_size_4_8kbyte":              struct{}{},
	"io_size_8_32kbyte":             struct{}{},
	"io_size_32_64kbyte":            struct{}{},
	"io_size_64_512kbyte":           struct{}{},
	"io_size_above_512kbyte":        struct{}{},
	"non_io_task_xcopy":             struct{}{},
	"non_io_task_unmap":             struct{}{},
	"non_io_task_write_same":        struct{}{},
	"non_io_task_ats":               struct{}{},
	"non_io_task_other":             struct{}{},
	"failed_task_check_cond":        struct{}{},
	"failed_task_busy":              struct{}{},
	"failed_task_resv_conflict":     struct{}{},
	"failed_task_abort":             struct{}{},
	"migrate_write_iops":            struct{}{},
	"migrate_write_bandwidth_kbyte": struct{}{},
	"migrate_write_latency_us":      struct{}{},
}
