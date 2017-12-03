# mackerel-plugin-ltsv-accesslog

LTSV format accesslog custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-ltsv-accesslog [-metric-key-prefix=<prefix>] [-no-posfile] [-posfile=<posfile>] [-request-time-key=<reqtime key> (default reqtime)] [-status-key=<status key> (default status)] [-tempfile=<tempfile>] /path/to/access.log
```

## Example of mackerel-agent.conf

```
[plugin.metrics.accesslog]
command = "/path/to/mackerel-plugin-ltsv-accesslog -request-time-key=reqtime -status-key=status /path/to/access.log"
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

Latency (Available only with LTSV format)

- accesslog.average
- accesslog.min
- accesslog.max
- accesslog.90_percentile
- accesslog.95_percentile
- accesslog.99_percentile
