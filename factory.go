// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package postgresqlexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

type Config struct {
	Path string `mapstructure:"path"`
}

// NewFactory creates a factory for the Parquet exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		"postgresqlexporter",
		CreateDefaultConfig,
		exporter.WithMetrics(CreateMetricsExporter, component.StabilityLevelDevelopment))
}

func CreateDefaultConfig() component.Config {
	return &Config{}
}

func CreateMetricsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Metrics, error) {
	fe := &postgresqlExporter{path: cfg.(*Config).Path}
	return exporterhelper.NewMetricsExporter(
		ctx,
		set,
		cfg,
		fe.consumeMetrics,
		exporterhelper.WithStart(fe.start),
		exporterhelper.WithShutdown(fe.shutdown),
	)
}
