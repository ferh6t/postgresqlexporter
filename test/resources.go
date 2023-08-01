package test

import "go.opentelemetry.io/collector/pdata/pcommon"

func initResource(r pcommon.Resource) {
	r.Attributes().PutStr("resource-attr", "resource-attr-val-1")
}
