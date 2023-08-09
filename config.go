package postgresqlexporter

import (
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines configuration for Elastic exporter.
type Config struct {
	exporterhelper.QueueSettings `mapstructure:"sending_queue"`

	Endpoint string `mapstructure:"endpoint"`
}
