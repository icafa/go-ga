package google_analytics

import (
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
)



func TestGoogleAnalyticsReport(t *testing.T) {

	ga := GoogleAnlayticsReport{
         KeyFile: "./plugins/inputs/google_analytics/clrty-secret.json",
         ViewID: "ga:155849743",
	}

	acc := testutil.Accumulator{}
	err := ga.Gather(&acc)
	assert.NoError(t, err)
}