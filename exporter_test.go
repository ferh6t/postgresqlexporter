package postgresqlexporter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ferh6t/postgresqlexporter/test"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func TestLoggingTracesExporterNoErrors(t *testing.T) {
}

func TestLoggingMetricsExporterNoErrors(t *testing.T) {
	f := NewFactory()
	lme, err := f.CreateMetricsExporter(context.Background(), exportertest.NewNopCreateSettings(), f.CreateDefaultConfig())
	require.NotNil(t, lme)
	assert.NoError(t, err)

	assert.NoError(t, lme.ConsumeMetrics(context.Background(), pmetric.NewMetrics()))
	assert.NoError(t, lme.ConsumeMetrics(context.Background(), test.GenerateMetricsAllTypes()))
	assert.NoError(t, lme.ConsumeMetrics(context.Background(), test.GenerateMetricsAllTypesEmpty()))
	assert.NoError(t, lme.ConsumeMetrics(context.Background(), test.GenerateMetricsMetricTypeInvalid()))
	assert.NoError(t, lme.ConsumeMetrics(context.Background(), test.GenerateMetrics(10)))

	// assert.NoError(t, lme.Shutdown(context.Background()))
}

func TestLoggingLogsExporterNoErrors(t *testing.T) {
}

func TestLoggingExporterErrors(t *testing.T) {
}

type errMarshaler struct {
	err error
}

func (e errMarshaler) MarshalLogs(plog.Logs) ([]byte, error) {
	return nil, e.err
}

func (e errMarshaler) MarshalMetrics(pmetric.Metrics) ([]byte, error) {
	return nil, e.err
}

func (e errMarshaler) MarshalTraces(ptrace.Traces) ([]byte, error) {
	return nil, e.err
}
