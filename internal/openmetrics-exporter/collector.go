package collectors

import (
	"context"
	client "purestorage/fb-openmetrics-exporter/internal/rest-client"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Collector(ctx context.Context, metrics string, registry *prometheus.Registry, fbclient *client.FBClient) bool {
	span, ctx := tracer.StartSpanFromContext(ctx, "collector", tracer.ResourceName("metrics.collector"))
	defer span.Finish()

	filesystems := fbclient.GetFileSystems()
	buckets := fbclient.GetBuckets()
	registry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)
	if metrics == "all" || metrics == "array" {
		arrayCollector := NewArraysCollector(fbclient)
		perfCollector := NewPerfCollector(fbclient)
		s3perfCollector := NewS3PerfCollector(fbclient)
		httpPerfCollector := NewHttpPerfCollector(fbclient)
		nfsPerfCollector := NewNfsPerfCollector(fbclient)
		perfReplCollector := NewPerfReplicationCollector(fbclient)
		bucketsPerfCollector := NewBucketsPerfCollector(fbclient, buckets)
		buckestS3PerfCollector := NewBucketsS3PerfCollector(fbclient, buckets)
		filesystemsPerfCollector := NewFileSystemsPerfCollector(fbclient, filesystems)
		filesystemsSpaceCollector := NewFileSystemsSpaceCollector(filesystems)
		arraySpaceCollector := NewArraySpaceCollector(fbclient)
		bucketsSpaceCollector := NewBucketsSpaceCollector(buckets)
		alertsCollector := NewAlertsCollector(fbclient)
		hardwareCollector := NewHardwareCollector(fbclient)
		hwPerfConnectorsCollector := NewHwConnectorsPerfCollector(fbclient)
		registry.MustRegister(
			arrayCollector,
			perfCollector,
			s3perfCollector,
			httpPerfCollector,
			nfsPerfCollector,
			perfReplCollector,
			bucketsPerfCollector,
			buckestS3PerfCollector,
			filesystemsPerfCollector,
			arraySpaceCollector,
			bucketsSpaceCollector,
			filesystemsSpaceCollector,
			alertsCollector,
			hardwareCollector,
			hwPerfConnectorsCollector,
		)
	}
	if metrics == "all" || metrics == "clients" {
		clientsPerfCollector := NewClientsPerfCollector(fbclient)
		registry.MustRegister(clientsPerfCollector)
	}
	if metrics == "all" || metrics == "usage" {
		usageCollector := NewUsageCollector(ctx, fbclient, filesystems)
		registry.MustRegister(usageCollector)
	}
	if metrics == "all" || metrics == "policies" {
		policiesCollector := NewNfsPoliciesCollector(fbclient)
		registry.MustRegister(policiesCollector)
	}
	return true
}
