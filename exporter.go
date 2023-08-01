// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package myapiserver // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/parquetexporter"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
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

	return nil
}

func (e postgresqlExporter) consumeTraces(_ context.Context, _ ptrace.Traces) error {
	return nil
}

func (e postgresqlExporter) consumeLogs(_ context.Context, _ plog.Logs) error {
	return nil
}
