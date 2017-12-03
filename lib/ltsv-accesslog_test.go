package mpltsvaccesslog

import (
	"reflect"
	"testing"
)

var fetchMetricsTests = []struct {
	Name   string
	InFile string
	Output map[string]float64
}{
	{
		Name:   "LTSV log",
		InFile: "testdata/sample-ltsv.tsv",
		Output: map[string]float64{
			"2xx_count":      7,
			"3xx_count":      1,
			"404_count":      1,
			"4xx_count":      0,
			"503_count":      0,
			"5xx_count":      1,
			"total_count":    10,
			"2xx_percentage": 70,
			"3xx_percentage": 10,
			"404_percentage": 10,
			"4xx_percentage": 0,
			"503_percentage": 0,
			"5xx_percentage": 10,
			"average":        0.7603999999999999,
			"min":            0.011,
			"max":            4.018,
			"90_percentile":  3.018,
			"95_percentile":  4.018,
			"99_percentile":  4.018,
		},
	},
	{
		Name:   "LTSV long line log",
		InFile: "testdata/sample-ltsv-long.tsv",
		Output: map[string]float64{
			"2xx_count":      2,
			"3xx_count":      0,
			"404_count":      0,
			"4xx_count":      0,
			"503_count":      0,
			"5xx_count":      0,
			"total_count":    2,
			"2xx_percentage": 100,
			"3xx_percentage": 0,
			"404_percentage": 0,
			"4xx_percentage": 0,
			"503_percentage": 0,
			"5xx_percentage": 0,
			"average":        0.015,
			"min":            0.01,
			"max":            0.02,
			"90_percentile":  0.02,
			"95_percentile":  0.02,
			"99_percentile":  0.02,
		},
	},
}

func TestFetchMetrics(t *testing.T) {
	for _, tt := range fetchMetricsTests {
		t.Logf("testing: %s", tt.Name)
		p := &LTSVAccesslogPlugin{
			File:       tt.InFile,
			NoPosFile:  true,
			StatusKey:  "status",
			ReqTimeKey: "reqtime",
		}
		out, err := p.FetchMetrics()
		if err != nil {
			t.Errorf("%s(err): error should be nil but: %+v", tt.Name, err)
			continue
		}
		if !reflect.DeepEqual(out, tt.Output) {
			t.Errorf("%s: \n out:  %#v\n want: %#v", tt.Name, out, tt.Output)
		}
	}
}

func TestFetchMetricsWithCustomKey(t *testing.T) {
	// OK case
	p := &LTSVAccesslogPlugin{
		File:       "testdata/sample-custom-ltsv.tsv",
		NoPosFile:  true,
		StatusKey:  "http_status",
		ReqTimeKey: "responsetime",
	}
	out, err := p.FetchMetrics()
	if err != nil {
		t.Errorf("error should be nil but: %+v", err)
		return
	}

	expected := map[string]float64{
		"2xx_count":      7,
		"3xx_count":      1,
		"404_count":      1,
		"4xx_count":      0,
		"503_count":      0,
		"5xx_count":      1,
		"total_count":    10,
		"2xx_percentage": 70,
		"3xx_percentage": 10,
		"404_percentage": 10,
		"4xx_percentage": 0,
		"503_percentage": 0,
		"5xx_percentage": 10,
		"average":        0.7603999999999999,
		"min":            0.011,
		"max":            4.018,
		"90_percentile":  3.018,
		"95_percentile":  4.018,
		"99_percentile":  4.018,
	}
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("out:  %#v\n want: %#v", out, expected)
	}
}
