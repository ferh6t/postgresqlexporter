// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package postgresqlexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		"postgresqlexporter",
		CreateDefaultConfig,
		exporter.WithMetrics(CreateMetricsExporter, component.StabilityLevelDevelopment))
}

func CreateDefaultConfig() component.Config {
	qs := exporterhelper.NewDefaultQueueSettings()
	qs.Enabled = false
	return &Config{}
}

func CreateMetricsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Metrics, error) {
	cf := cfg.(*Config)
	exporter := NewPostgreSqlExporter(cf)
	return exporterhelper.NewMetricsExporter(
		ctx,
		set,
		cfg,
		exporter.consumeMetrics,
		exporterhelper.WithStart(exporter.Start),
		exporterhelper.WithShutdown(exporter.Shutdown),
	)
}
