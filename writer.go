package postgresqlexporter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/common/model"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
)

type writerToPostgresql struct {
	conn *conn
}

func newWriter(cfg *Config) *writerToPostgresql {
	default_link := "host=localhost user=your_username password=your_password  dbname=your_database_name  port=5432  sslmode=disable"
	var conn *conn
	if cfg.Endpoint != "" {
		conn = InitDatabase(cfg.Endpoint)
	} else {
		conn = InitDatabase(default_link)
	}

	return &writerToPostgresql{
		conn: conn,
	}
}

func (a *writerToPostgresql) WriteToDataBase(resourceMetrics pmetric.ResourceMetrics) {
	for j := 0; j < resourceMetrics.ScopeMetrics().Len(); j++ {
		resourceAttrs := resourceMetrics.Resource().Attributes()
		ilMetrics := resourceMetrics.ScopeMetrics().At(j)
		for k := 0; k < ilMetrics.Metrics().Len(); k++ {
			metricK := ilMetrics.Metrics().At(k)
			now := time.Now()
			switch metricK.Type() {
			case pmetric.MetricTypeGauge:
				a.accumulateGauge(metricK, ilMetrics.Scope(), resourceAttrs, now)
			case pmetric.MetricTypeHistogram:
				a.accumulateDoubleHistogram(metricK, ilMetrics.Scope(), resourceAttrs, now)
			case pmetric.MetricTypeSum:
				a.accumulateSum(metricK, ilMetrics.Scope(), resourceAttrs, now)
			case pmetric.MetricTypeSummary:
				a.accumulateSummary(metricK, ilMetrics.Scope(), resourceAttrs, now)
			}

		}
	}
}

func (a *writerToPostgresql) accumulateGauge(metric pmetric.Metric, il pcommon.InstrumentationScope, resourceAttrs pcommon.Map, now time.Time) (n int) {
	dps := metric.Gauge().DataPoints()
	fmt.Println("unit is " + metric.Unit())
	for i := 0; i < dps.Len(); i++ {
		ip := dps.At(i)
		signature := timeseriesSignature(il.Name(), metric, ip.Attributes(), resourceAttrs)

		myGauge := Gauge{
			Signature:  signature,
			Timestamp:  ip.Timestamp().AsTime(),
			Attributes: spew.Sdump(ip.Attributes()),
			IntVal:     ip.IntValue(),
			DoubleVal:  ip.DoubleValue(),
			IsInt:      ip.ValueType() != pmetric.NumberDataPointValueTypeDouble,
		}

		a.conn.SaveGauge(myGauge)

	}
	return
}

func (a *writerToPostgresql) accumulateSum(metric pmetric.Metric, il pcommon.InstrumentationScope, resourceAttrs pcommon.Map, now time.Time) (n int) {
	doubleSum := metric.Sum()

	if doubleSum.AggregationTemporality() == pmetric.AggregationTemporalityUnspecified {
		return
	}

	if doubleSum.AggregationTemporality() == pmetric.AggregationTemporalityDelta && !doubleSum.IsMonotonic() {
		return
	}

	dps := doubleSum.DataPoints()
	for i := 0; i < dps.Len(); i++ {
		ip := dps.At(i)

		signature := timeseriesSignature(il.Name(), metric, ip.Attributes(), resourceAttrs)
		mySum := Sum{
			Signature:  signature,
			Timestamp:  ip.Timestamp().AsTime(),
			Attributes: spew.Sdump(ip.Attributes()),
			IntVal:     ip.IntValue(),
			DoubleVal:  ip.DoubleValue(),
			IsInt:      ip.ValueType() != pmetric.NumberDataPointValueTypeDouble,
		}

		a.conn.SaveSum(mySum)
	}
	return
}

func (a *writerToPostgresql) accumulateSummary(metric pmetric.Metric, il pcommon.InstrumentationScope, resourceAttrs pcommon.Map, now time.Time) (n int) {
	summary := metric.Summary()

	dps := summary.DataPoints()
	for i := 0; i < dps.Len(); i++ {
		ip := dps.At(i)
		fmt.Println(strconv.FormatUint(ip.Count(), 10))
		signature := timeseriesSignature(il.Name(), metric, ip.Attributes(), resourceAttrs)
		fmt.Println(signature)
		fmt.Println(strconv.FormatFloat(ip.Sum(), 'f', -1, 64))
		fmt.Println(spew.Sdump(ip.Attributes()))
		mySummary := Summary{
			Signature:  signature,
			Timestamp:  ip.Timestamp().AsTime(),
			Attributes: spew.Sdump(ip.Attributes()),
			Sum:        ip.Sum(),
		}

		a.conn.SaveSummary(mySummary)
	}
	return
}

func (a *writerToPostgresql) accumulateDoubleHistogram(metric pmetric.Metric, il pcommon.InstrumentationScope, resourceAttrs pcommon.Map, now time.Time) (n int) {
	doubleHistogram := metric.Histogram()

	// Drop metrics with non-cumulative aggregations
	if doubleHistogram.AggregationTemporality() != pmetric.AggregationTemporalityCumulative {
		return
	}

	dps := doubleHistogram.DataPoints()

	for i := 0; i < dps.Len(); i++ {
		ip := dps.At(i)
		signature := timeseriesSignature(il.Name(), metric, ip.Attributes(), resourceAttrs)
		fmt.Println(ip.Count())
		fmt.Println(ip.Sum())
		fmt.Println(spew.Sdump(ip.Attributes()))
		myhistgoram := Histogram{
			Signature:  signature,
			Timestamp:  ip.Timestamp().AsTime(),
			Attributes: spew.Sdump(ip.Attributes()),
			IntVal:     ip.Sum(),
			Count:      ip.Count(),
		}

		a.conn.SaveHistogram(myhistgoram)
	}

	return
}

func (a *writerToPostgresql) accumulateExponentialHistogram(metric pmetric.Metric, il pcommon.InstrumentationScope, resourceAttrs pcommon.Map, now time.Time) (n int) {
	doubleHistogram := metric.ExponentialHistogram()

	// Drop metrics with non-cumulative aggregations
	if doubleHistogram.AggregationTemporality() != pmetric.AggregationTemporalityCumulative {
		return
	}

	dps := doubleHistogram.DataPoints()
	for i := 0; i < dps.Len(); i++ {

		for i := 0; i < dps.Len(); i++ {
			fmt.Println("one")
			ip := dps.At(i)
			fmt.Println(ip.Max())
			fmt.Println(ip.Min())
		}
	}
	return
}

func timeseriesSignature(ilmName string, metric pmetric.Metric, attributes pcommon.Map, resourceAttrs pcommon.Map) string {
	var b strings.Builder
	b.WriteString(metric.Type().String())
	b.WriteString("*" + ilmName)
	b.WriteString("*" + metric.Name())
	attrs := make([]string, 0, attributes.Len())
	attributes.Range(func(k string, v pcommon.Value) bool {
		attrs = append(attrs, k+"*"+v.AsString())
		return true
	})
	sort.Strings(attrs)
	b.WriteString("*" + strings.Join(attrs, "*"))
	if job, ok := extractJob(resourceAttrs); ok {
		b.WriteString("*" + model.JobLabel + "*" + job)
	}
	if instance, ok := extractInstance(resourceAttrs); ok {
		b.WriteString("*" + model.InstanceLabel + "*" + instance)
	}
	return b.String()
}

func extractInstance(attributes pcommon.Map) (string, bool) {
	// Map service.instance.id to instance
	if inst, ok := attributes.Get(conventions.AttributeServiceInstanceID); ok {
		return inst.AsString(), true
	}
	return "", false
}

func extractJob(attributes pcommon.Map) (string, bool) {
	// Map service.namespace + service.name to job
	if serviceName, ok := attributes.Get(conventions.AttributeServiceName); ok {
		job := serviceName.AsString()
		if serviceNamespace, ok := attributes.Get(conventions.AttributeServiceNamespace); ok {
			job = fmt.Sprintf("%s/%s", serviceNamespace.AsString(), job)
		}
		return job, true
	}
	return "", false
}
