# mackerel-plugin-ltsv-accesslog

LTSV format accesslog custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-ltsv-accesslog [-metric-key-prefix=<prefix>] [-no-posfile] [-posfile=<posfile>] [-request-time-key=<reqtime key> (default reqtime)] [-status-key=<status key> (default status)] [-cache-status-key=<cache status key> (default upstream_cache_status)] [-tempfile=<tempfile>] /path/to/access.log
```

## Example of mackerel-agent.conf

```
[plugin.metrics.accesslog]
command = "/path/to/mackerel-plugin-ltsv-accesslog -request-time-key=reqtime -status-key=status -cache-status-key=upstream_cache_status /path/to/access.log"
```

## Graphs and Metrics

### accesslog.access_num

- accesslog.access_num.total_count
- accesslog.access_num.2xx_count
- accesslog.access_num.3xx_count
- accesslog.access_num.400_count
- accesslog.access_num.4xx_count
- accesslog.access_num.503_count
- accesslog.access_num.5xx_count

### accesslog.access_rate

- accesslog.access_rate.2xx_percentage
- accesslog.access_rate.3xx_percentage
- accesslog.access_rate.404_percentage
- accesslog.access_rate.4xx_percentage
- accesslog.access_rate.503_percentage
- accesslog.access_rate.5xx_percentage

## accesslog.latency

- accesslog.latency.average
- accesslog.latency.min
- accesslog.latency.max
- accesslog.latency.90_percentile
- accesslog.latency.95_percentile
- accesslog.latency.99_percentile

## accesslog.cache_rate

- accesslog.cache_rate.HIT_percentage
- accesslog.cache_rate.MISS_percentage
- accesslog.cache_rate.EXPIRED_percentage
- accesslog.cache_rate.REVALIDATED_percentage
- accesslog.cache_rate.BYPASS_percentage
- accesslog.cache_rate.STALE_percentage
- accesslog.cache_rate.UPDATING_percentage
