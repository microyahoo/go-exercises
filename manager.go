package metrics

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"

	"xsky-demon/config"
	"xsky-demon/errors"
	"xsky-demon/etcd"
	"xsky-demon/log"
	"xsky-demon/models"
	"xsky-demon/utils"
)

var logger = log.NewLogger(log.AppKernel, "metric", "", log.Level.Info)

type queryResult struct {
	statMap   map[string]float64
	timestamp time.Time
}

// labelsResultMap's map key is the metric name
// labelsResultMap's map value's key is the labels string
type labelsResultMap map[string]map[string]model.Metric

type metricResultMap map[model.Time]model.SampleValue

// statsResultMap's map key is the metric name
// value is a map with sample timestamp as map key, sample value as map value
type statsResultMap map[string]metricResultMap

// statsResultWithLabelMap's map key is the metric name
// statsResultWithLabelMap's map value's key is the labels string
// statsResultWithLabelMap's map value's value is a map with sample timestamp as map key,
// sample value as map value
type statsResultWithLabelMap map[string]map[string]metricResultMap

// range query result is the result returned in range query
type rangeQeuryResult struct {
	statsMap   statsResultMap // sample result map
	start, end model.Time     // the earliest and latest timestamp of all metrics
	step       time.Duration  // time step of the query
}

// rangeQueryResultWithLabel is the result returned in range query
type rangeQueryResultWithLabel struct {
	labelsMap  labelsResultMap         // labels result map
	statsMap   statsResultWithLabelMap // sample result map
	start, end model.Time              // the earliest and latest timestamp of all metrics
	step       time.Duration           // time step of the query
}

var manager = new(prometheusManager)

// CurrentManager returns the singleton prometheus manager
var CurrentManager = currentManager

func currentManager() Manager {
	return manager
}

// Manager defines prometheus manager interface
type Manager interface {
	Init(host string, etm *etcd.Manager) error
	Host() string
	GetVolumeStats(volume *models.Volume, params models.PagingParams) ([]*models.VolumeStat, error)
	DeleteVolumeStat(volume *models.Volume) error
	GetCloudInstanceStats(instance *models.CloudInstance, params models.PagingParams) (
		[]*models.CloudInstanceStat, error)
	DeleteCloudInstanceStat(instance *models.CloudInstance) error
	GetHostStats(host *models.Host, params models.PagingParams) ([]*models.HostStat, error)
	DeleteHostStat(host *models.Host) error
	GetNetworkInterfaceStats(nic *models.NetworkInterface, params models.PagingParams) (
		[]*models.NetworkInterfaceStat, error)
	DeleteNetworkInterfaceStat(nic *models.NetworkInterface) error
	GetDiskStats(disk *models.Disk, params models.PagingParams) ([]*models.DiskStat, error)
	DeleteDiskStat(disk *models.Disk) (err error)
	GetOsdStats(osd *models.Osd, params models.PagingParams) ([]*models.OsdStat, error)
	DeleteOsdStat(osd *models.Osd) (err error)
	GetPartitionStats(partition *models.Partition, params models.PagingParams) ([]*models.PartitionStat, error)
	DeletePartitionStat(partition *models.Partition) (err error)
	GetObjectStorageUserStats(user *models.ObjectStorageUser, params models.PagingParams) (
		[]*models.ObjectStorageUserStat, error)
	DeleteObjectStorageUserStat(user *models.ObjectStorageUser) error
	GetObjectStorageBucketStats(bucket *models.ObjectStorageBucket, params models.PagingParams) (
		[]*models.ObjectStorageBucketStat, error)
	DeleteObjectStorageBucketStat(bucket *models.ObjectStorageBucket) error
	GetObjectStorageGatewayStats(gateway *models.ObjectStorageGateway, params models.PagingParams) (
		[]*models.ObjectStorageGatewayStat, error)
	DeleteObjectStorageGatewayStat(gateway *models.ObjectStorageGateway) error
	GetObjectStorageZoneStats(zone *models.ObjectStorageZone,
		params models.PagingParams) ([]*models.ObjectStorageZoneStat, error)
	GetOSSearchGatewayStats(osGateway *models.OSSearchGateway, params models.PagingParams) (
		[]*models.OSSearchGatewayStat, error)
	DeleteObjectStorageZoneStat(zone *models.ObjectStorageZone) error
	GetOSReplicationZoneStats(zone *models.OSReplicationZone,
		params models.PagingParams) ([]*models.OSReplicationZoneStat, error)
	DeleteOSReplicationZoneStat(zone *models.OSReplicationZone) error
	GetNFSGatewayStats(gateway *models.NFSGateway, params models.PagingParams) ([]*models.NFSGatewayStat, error)
	DeleteNFSGatewayStat(gateway *models.NFSGateway) error

	GetFSFolderStats(*models.FSFolder, models.PagingParams) ([]*models.FSFolderStat, error)
	DeleteFSFolderStat(*models.FSFolder) error

	GetPoolStats(pool *models.Pool, params models.PagingParams) ([]*models.PoolStat, error)
	DeletePoolStat(pool *models.Pool) error
	GetClusterStats(cluster *models.Cluster, params models.PagingParams) ([]*models.ClusterStat, error)
	DeleteClusterStat(cluster *models.Cluster) error
	GetOsdGroupStats(group *models.OsdGroup, params models.PagingParams) ([]*models.OsdGroupStat, error)
	DeleteOsdGroupStat(group *models.OsdGroup) error

	GetS3LoadBalancerStats(*models.S3LoadBalancer, models.PagingParams) ([]*models.S3LoadBalancerStat, error)
	DeleteS3LoadBalancerStat(*models.S3LoadBalancer) error

	GetOSSearchEngineStats(engine *models.OSSearchEngine, params models.PagingParams) ([]*models.OSSearchEngineStat, error)
	DeleteOSSearchEngineStat(engine *models.OSSearchEngine) error
	GetOSSamples(label map[string]string, ts ...time.Time) ([]map[string][]map[string]*models.OSSample, error)
}

type prometheusManager struct {
	sync.Mutex
	host       string
	client     api.Client
	api        v1.API
	etm        *etcd.Manager
	workerPool *utils.WorkerPool
}

func (p *prometheusManager) Init(host string, etm *etcd.Manager) (err error) {
	p.Lock()
	defer p.Unlock()
	cfg := api.Config{Address: p.metricAddr(host)}
	client, err := newMetricClient(cfg)
	if err != nil {
		return errors.Trace(err)
	}
	p.host = host
	p.client = client
	p.api = v1.NewAPI(p.client)
	p.etm = etm

	p.workerPool = &utils.WorkerPool{
		WorkerFunc: p.doRequest,
		// TODO(zhenliang):update it
		// invoke Stop()
		MaxWorkersCount: 12000,
	}
	p.workerPool.Start()

	return nil
}

func (p *prometheusManager) doRequest(req *http.Request) error {
	ctx, cancel := utils.TimeoutContext(config.C.Metric.MetricDeleteAPITimeout)
	defer cancel()
	_, _, err := p.client.Do(ctx, req)
	if err != nil {
		log.Warnf("Failed to delete series of %s: %s", req.URL.String(), err)
		return errors.Trace(err)
	}
	return nil
}

func (p *prometheusManager) Host() string {
	return p.host
}

func (p *prometheusManager) query(query string, timestamp time.Time) (result model.Vector, err error) {
	cxt, cancel := utils.TimeoutContext(config.C.Metric.MetricAPITimeout)
	defer cancel()
	v, err := p.api.Query(cxt, query, timestamp)
	if err != nil {
		err = errors.Annotatef(err, "error in query request")
		return
	}
	result, ok := v.(model.Vector)
	if !ok {
		err = errors.Errorf("failed to convert query result to Vector type")
		return
	}
	return
}

func (p *prometheusManager) queryRange(query string, queryRange v1.Range) (result model.Matrix, err error) {

	cxt, cancel := utils.TimeoutContext(config.C.Metric.MetricAPITimeout)
	defer cancel()
	v, err := p.api.QueryRange(cxt, query, queryRange)
	if err != nil {
		err = errors.Annotatef(err, "error in range query request")
		return
	}
	result, ok := v.(model.Matrix)
	if !ok {
		err = errors.Errorf("failed to convert range query result to Matrix type")
		return
	}
	return
}

func (p *prometheusManager) queryMetrics(metricQueries queryMap, labels prometheus.Labels, ts ...time.Time) (
	result queryResult, err error) {

	var timestamp time.Time
	if len(ts) > 0 {
		timestamp = ts[0]
	} else {
		timestamp = time.Now()
	}

	result.statMap = make(map[string]float64, len(metricQueries))
	var metricResult model.Vector
	for key, metric := range metricQueries {
		query := &Query{metric: metric}
		query.AddLabels(labels)
		queryStr := query.String()

		metricResult, err = p.query(queryStr, timestamp)
		logger.Debugf("start query for metric %s using query string: %s", metric, queryStr)
		if err != nil {
			err = errors.Annotatef(err, "failed to query stat of metric: %s", metric)
			return
		}
		if len(metricResult) == 0 {
			logger.Warnf("no result found for metric %s", metric)
			return
		}
		value := float64(metricResult[0].Value)
		result.statMap[key] = value
	}
	result.timestamp = timestamp

	return
}

func (p *prometheusManager) queryRangeMetrics(metricQueries queryMap, labels prometheus.Labels,
	params models.PagingParams) (result rangeQeuryResult, err error) {

	begin, end, step := ParsePagingParams(params)
	logger.Debugf("query params: begin: %s, end: %s, step: %s", begin, end, step)
	queryRange := v1.Range{Start: begin, End: end, Step: step}

	result.statsMap = make(statsResultMap)
	var resultMatrix model.Matrix
	for key, metric := range metricQueries {
		queryStr := queryRangeStr(metric, labels, step)
		resultMatrix, err = p.queryRange(queryStr, queryRange)
		logger.Debugf("start range query for metric %s using query string: %s", metric, queryStr)
		if err != nil {
			err = errors.Annotatef(err, "failed to query stat of metric: %s", metric)
			return
		}

		stats := make(map[model.Time]model.SampleValue)
		if resultMatrix.Len() == 0 {
			logger.Debugf("no result found for query %s", queryStr)
		} else {
			values := resultMatrix[0].Values
			start, end := values[0].Timestamp, values[len(values)-1].Timestamp

			// find the start and end timestamp of results
			if result.start == 0 || start.Before(result.start) {
				result.start = start
			}
			if result.end == 0 || end.After(result.end) {
				result.end = end
			}

			for _, s := range values {
				stats[s.Timestamp] = s.Value
			}
		}
		result.statsMap[key] = stats
		result.step = step
	}
	return
}

// queryRangeMetricsWithLabels return all metric labels and corresponding time series data
func (p *prometheusManager) queryRangeMetricsWithLabels(metricQueries queryMap, labels prometheus.Labels,
	params models.PagingParams) (result *rangeQueryResultWithLabel, err error) {

	begin, end, step := ParsePagingParams(params)
	logger.Debugf("query params: begin: %s, end: %s, step: %s", begin, end, step)
	queryRange := v1.Range{Start: begin, End: end, Step: step}

	result = new(rangeQueryResultWithLabel)
	result.labelsMap = make(labelsResultMap)
	result.statsMap = make(statsResultWithLabelMap)
	result.step = step

	var resultMatrix model.Matrix
	for key, metric := range metricQueries {
		queryStr := queryRangeStrWithoutMax(metric, labels, step)
		resultMatrix, err = p.queryRange(queryStr, queryRange)
		logger.Debugf("start range query for metric %s using query string: %s", metric, queryStr)
		if err != nil {
			err = errors.Annotatef(err, "query stat of metric: %s", metric)
			return nil, err
		}

		labelsMap := make(map[string]model.Metric)
		result.labelsMap[key] = labelsMap
		statsMap := make(map[string]metricResultMap)
		result.statsMap[key] = statsMap
		for _, matrix := range resultMatrix {
			metric := matrix.Metric
			// when prometheus scrapes a target, it attaches instance and job labels
			// automatically to the scraped time series which serve to identify the scraped target：
			// 	   job: The configured job name that the target belongs to.
			//     instance: The <host>:<port> part of the target's URL that was scraped.
			// matrix.Metric contains all labels of time series, include instance and job, and for
			// time series that scraped from different node, its instance labels is different, so
			// we need delete instance label, by the way, delete useless job label.
			delete(metric, "instance")
			delete(metric, "job")
			labelStr := metric.String()
			labelsMap[labelStr] = metric
			stats := make(map[model.Time]model.SampleValue)
			statsMap[labelStr] = stats

			values := matrix.Values
			start, end := values[0].Timestamp, values[len(values)-1].Timestamp

			// find the start and end timestamp of results
			if result.start == 0 || start.Before(result.start) {
				result.start = start
			}
			if result.end == 0 || end.After(result.end) {
				result.end = end
			}

			for _, s := range values {
				stats[s.Timestamp] = s.Value
			}
		}
	}

	return result, nil
}

func (p *prometheusManager) queryMetricByLabels(metricQueries queryMap, labels prometheus.Labels,
	ts ...time.Time) (resultMap map[string]model.Vector, err error) {

	var timestamp time.Time
	if len(ts) > 0 {
		timestamp = ts[0]
	} else {
		timestamp = time.Now()
	}

	var metricResult model.Vector
	resultMap = map[string]model.Vector{}
	for key, metric := range metricQueries {
		query := &Query{metric: metric}
		query.AddLabels(labels)
		queryStr := query.String()

		metricResult, err = p.query(queryStr, timestamp)
		logger.Debugf("start query for metric %s using query string: %s", metric, queryStr)
		if err != nil {
			return nil, errors.Annotatef(err, "failed to query stat of metric: %s", metric)
		}
		if len(metricResult) == 0 {
			logger.Warnf("no result found for metric %s", metric)
			return nil, nil
		}
		resultMap[key] = metricResult
	}
	return resultMap, nil
}

func (p *prometheusManager) delete(host string, query url.Values) (err error) {
	u, err := deleteURL(host)
	if err != nil {
		return errors.Trace(err)
	}
	u.RawQuery = query.Encode()

	logger.Debugf("Delete series: %s", u.String())

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return errors.Trace(err)
	}
	req.Header.Add("Content-Type", "application/json")

	ok := p.workerPool.Serve(req)
	if !ok {
		logger.Warnf("Failed to fetch worker to serve the request of %s", u.String())
	}

	return nil
}

// deleteAll deletes metrics on all prometheus servers(candidate role)
func (p *prometheusManager) deleteAll(query url.Values) (err error) {
	candidates, _ := p.etm.ListCandidates()
	for _, ip := range candidates {
		e := p.delete(ip, query)
		if e != nil {
			logger.Debugf("failed to delete metric on host %s:%s", ip, e)
			if err == nil {
				err = e
			}
		}
	}
	return
}

func (p *prometheusManager) deleteMetrics(metricQueries queryMap, id int64) (err error) {
	query := url.Values{}
	for _, metric := range metricQueries {
		matcher := fmt.Sprintf(`%s{id="%d"}`, metric, id)
		query.Add("match[]", matcher)
	}
	if err = p.deleteAll(query); err != nil {
		logger.Warnf("Failed to delete stat of id %d: %s, query: %v", id, err, query)
	}
	return nil
}

func (p *prometheusManager) metricAddr(host string) string {
	return fmt.Sprintf("http://%s/", utils.JoinHostPort(host, int64(config.C.Metric.MetricPort)))
}
