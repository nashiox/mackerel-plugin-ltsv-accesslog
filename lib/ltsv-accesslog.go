package mpltsvaccesslog

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Songmu/postailer"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/mackerelio/golib/pluginutil"
	"github.com/montanaflynn/stats"
	"github.com/najeira/ltsv"
)

// LTSVAccesslogPlugin mackerel plugin
type LTSVAccesslogPlugin struct {
	Prefix         string
	File           string
	PosFile        string
	NoPosFile      bool
	StatusKey      string
	ReqTimeKey     string
	CacheStatusKey string
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p *LTSVAccesslogPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "accesslog"
	}
	return p.Prefix
}

// GraphDefinition interface for mackerelplugin
func (p *LTSVAccesslogPlugin) GraphDefinition() map[string](mp.Graphs) {
	labelPrefix := strings.Title(p.Prefix)
	return map[string]mp.Graphs{
		"access_num": {
			Label: labelPrefix + " Access Num",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "total_count", Label: "Total Count"},
				{Name: "503_count", Label: "HTTP 503 Count", Stacked: true},
				{Name: "5xx_count", Label: "HTTP 5xx Count", Stacked: true},
				{Name: "404_count", Label: "HTTP 404 Count", Stacked: true},
				{Name: "4xx_count", Label: "HTTP 4xx Count", Stacked: true},
				{Name: "3xx_count", Label: "HTTP 3xx Count", Stacked: true},
				{Name: "2xx_count", Label: "HTTP 2xx Count", Stacked: true},
			},
		},
		"access_rate": {
			Label: labelPrefix + " Access Rate",
			Unit:  "percentage",
			Metrics: []mp.Metrics{
				{Name: "503_percentage", Label: "HTTP 503 Percentage", Stacked: true},
				{Name: "5xx_percentage", Label: "HTTP 5xx Percentage", Stacked: true},
				{Name: "404_percentage", Label: "HTTP 404 Percentage", Stacked: true},
				{Name: "4xx_percentage", Label: "HTTP 4xx Percentage", Stacked: true},
				{Name: "3xx_percentage", Label: "HTTP 3xx Percentage", Stacked: true},
				{Name: "2xx_percentage", Label: "HTTP 2xx Percentage", Stacked: true},
			},
		},
		"latency": {
			Label: labelPrefix + " Latency",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "99_percentile", Label: "99 Percentile"},
				{Name: "95_percentile", Label: "95 Percentile"},
				{Name: "90_percentile", Label: "90 Percentile"},
				{Name: "average", Label: "Average"},
				{Name: "min", Label: "Min"},
				{Name: "max", Label: "Max"},
			},
		},
		"cache_rate": {
			Label: labelPrefix + " Cache Status",
			Unit:  "percentage",
			Metrics: []mp.Metrics{
				{Name: "HIT_percentage", Label: "HIT Percentage", Stacked: true},
				{Name: "MISS_percentage", Label: "MISS Percentage", Stacked: true},
				{Name: "EXPIRED_percentage", Label: "EXPIRED Percentage", Stacked: true},
				{Name: "REVALIDATED_percentage", Label: "REVALIDATED Percentage", Stacked: true},
				{Name: "BYPASS_percentage", Label: "BYPASS Percentage", Stacked: true},
				{Name: "STALE_percentage", Label: "STALE Percentage", Stacked: true},
				{Name: "UPDATING_percentage", Label: "UPDATING Percentage", Stacked: true},
			},
		},
	}
}

var posRe = regexp.MustCompile(`^([a-zA-Z]):[/\\]`)

func (p *LTSVAccesslogPlugin) getPosPath() string {
	base := p.File + ".pos.json"
	if p.PosFile != "" {
		if filepath.IsAbs(p.PosFile) {
			return p.PosFile
		}
		base = p.PosFile
	}

	return filepath.Join(
		pluginutil.PluginWorkDir(),
		"mackerel-plugin-ltsv-accesslog.d",
		posRe.ReplaceAllString(base, `$1`+string(filepath.Separator)),
	)
}

func (p *LTSVAccesslogPlugin) getReadCloser() (io.ReadCloser, bool, error) {
	if p.NoPosFile {
		rc, err := os.Open(p.File)
		return rc, true, err
	}

	posfile := p.getPosPath()
	fi, err := os.Stat(posfile)

	takeCount := err == nil && fi.ModTime().After(time.Now().Add(-2*time.Minute))
	rc, err := postailer.Open(p.File, posfile)
	return rc, takeCount, err
}

// FetchMetrics interface for mackerelplugin
func (p *LTSVAccesslogPlugin) FetchMetrics() (map[string]float64, error) {
	rc, takeCount, err := p.getReadCloser()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	countMetrics := []string{"total_count", "2xx_count", "3xx_count", "4xx_count", "404_count", "5xx_count", "503_count"}
	ret := make(map[string]float64)
	cacheCount := make(map[string]float64)
	for _, k := range countMetrics {
		ret[k] = 0
	}
	var reqtimes []float64
	r := ltsv.NewReader(rc)

	for {
		record, err := r.Read()
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}

		if record[p.StatusKey] == "404" {
			ret[string(record[p.StatusKey])+"_count"]++
		} else if record[p.StatusKey] == "503" {
			ret[string(record[p.StatusKey])+"_count"]++
		} else {
			ret[string(record[p.StatusKey][0])+"xx_count"]++
		}

		if record[p.CacheStatusKey] != "" && record[p.CacheStatusKey] != "-" {
			cacheCount[string(record[p.CacheStatusKey])+"_count"]++
			cacheCount["cache_total_count"]++
		}

		ret["total_count"]++

		v, err := strconv.ParseFloat(record[p.ReqTimeKey], 64)
		if err != nil {
			log.Println(err)
		}
		reqtimes = append(reqtimes, v)
	}

	if cacheCount["cache_total_count"] > 0 {
		for _, v := range []string{"HIT", "MISS", "EXPIRED", "REVALIDATED", "BYPASS", "STALE", "UPDATING"} {
			ret[v+"_percentage"] = cacheCount[v+"_count"] * 100 / cacheCount["cache_total_count"]
		}
	}

	if ret["total_count"] > 0 {
		for _, v := range []string{"2xx", "3xx", "4xx", "404", "5xx", "503"} {
			ret[v+"_percentage"] = ret[v+"_count"] * 100 / ret["total_count"]
		}
	}

	if len(reqtimes) > 0 {
		ret["average"], _ = stats.Mean(reqtimes)
		ret["min"], _ = stats.Min(reqtimes)
		ret["max"], _ = stats.Max(reqtimes)

		for _, v := range []int{90, 95, 99} {
			ret[fmt.Sprintf("%d_percentile", v)], _ = stats.Percentile(reqtimes, float64(v))
		}
	}
	if !takeCount {
		for _, k := range countMetrics {
			delete(ret, k)
		}
	}
	return ret, nil
}

// main function
func Do() {
	var (
		optPrefix         = flag.String("metric-key-prefix", "", "Metric key prefix")
		optPosFile        = flag.String("posfile", "", "(Not necessary to specify it in the usual use case) posfile")
		optNoPosFile      = flag.Bool("no-posfile", false, "No position file")
		optStatusKey      = flag.String("status-key", "status", "Status key name in log format")
		optReqTimeKey     = flag.String("request-time-key", "reqtime", "Request time key name in log format")
		optCacheStatusKey = flag.String("cache-status-key", "upstream_cache_status", "Cache Status key name in log format")
		optTempfile       = flag.String("tempfile", "", "Temp file name")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION] /path/to/access.log\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	helper := mp.NewMackerelPlugin(&LTSVAccesslogPlugin{
		Prefix:         *optPrefix,
		File:           flag.Args()[0],
		PosFile:        *optPosFile,
		NoPosFile:      *optNoPosFile,
		StatusKey:      *optStatusKey,
		ReqTimeKey:     *optReqTimeKey,
		CacheStatusKey: *optCacheStatusKey,
	})
	helper.Tempfile = *optTempfile

	helper.Run()
}
