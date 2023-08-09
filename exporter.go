// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package postgresqlexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type postgresqlExporter struct {
	path   string
	writer *writerToPostgresql
}

func NewPostgreSqlExporter(cfg *Config) *postgresqlExporter {
	return &postgresqlExporter{
		writer: newWriter(),
	}
}

func (e postgresqlExporter) Start(_ context.Context, _ component.Host) (err error) {
	return nil
}

func (e postgresqlExporter) Shutdown(_ context.Context) error {
	return nil
}

func (e postgresqlExporter) consumeMetrics(ctx context.Context, metric pmetric.Metrics, cfg component.Config) error {

	for i := 0; i < metric.ResourceMetrics().Len(); i++ {
		resourceMetrics := metric.ResourceMetrics().At(i)
		e.writer.WriteToDataBase(resourceMetrics)
	}

	return nil
}

// if necessary they will be implemeted
// func (e postgresqlExporter) consumeTraces(_ context.Context, _ ptrace.Traces) error {
// 	return nil
// }

// func (e postgresqlExporter) consumeLogs(_ context.Context, _ plog.Logs) error {
// 	return nil
// }
