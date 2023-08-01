// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package postgresqlexporter

import (
	"context"
	"fmt"
	"strconv"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type postgresqlExporter struct {
	path string
}

func (e postgresqlExporter) start(_ context.Context, _ component.Host) error {
	return nil
}

func (e postgresqlExporter) shutdown(_ context.Context) error {
	return nil
}

func (e postgresqlExporter) consumeMetrics(_ context.Context, metric pmetric.Metrics) error {
	fmt.Printf(strconv.Itoa(metric.MetricCount()) + "/n")
	return nil
}

// if necessary they will be implemeted
// func (e postgresqlExporter) consumeTraces(_ context.Context, _ ptrace.Traces) error {
// 	return nil
// }

// func (e postgresqlExporter) consumeLogs(_ context.Context, _ plog.Logs) error {
// 	return nil
// }
