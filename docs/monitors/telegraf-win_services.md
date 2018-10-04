<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# telegraf/win_services

 This monitor reports metrics about Windows services.
This monitor is based on the Telegraf win_services plugin.  More information about the Telegraf plugin
can be found [here](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/win_services).

Sample YAML configuration:

```yaml
monitors:
 - type: telegraf/win_services  # monitor all services
```

```yaml
monitors:
 - type: telegraf/win_services
   serviceNames:
     - exampleService1  # only monitor exampleService1
```


Monitor Type: `telegraf/win_services`

[Monitor Source Code](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/telegraf/monitors/winservices)

**Accepts Endpoints**: No

**Multiple Instances Allowed**: Yes

## Configuration

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `serviceNames` | no | `list of string` | Names of services to monitor.  All services will be monitored if none are specified. |




## Metrics

The following table lists the metrics available for this monitor. Metrics that are not marked as Custom are standard metrics and are monitored by default.

| Name | Type | Custom | Description |
| ---  | ---  | ---    | ---         |
| `win_services.startup_mode` | gauge | X | The startup mode configured for the service: 0 `boot start`, 1 `system start`, 2 `auto start`, 3 `demand start`, or 4 `disabled` |
| `win_services.state` | gauge | X | The current state of the service: 1 `stopped`, 2 `start pending`, 3 `stop pending`, 4 `running`, 5 `continue pending`, 6 `pause pending`, or 7 `paused` |


To specify custom metrics you want to monitor, add a `metricsToInclude` filter
to the agent configuration, as shown in the code snippet below. The snippet
lists all available custom metrics. You can copy and paste the snippet into
your configuration file, then delete any custom metrics that you do not want
sent.

Note that some of the custom metrics require you to set a flag as well as add
them to the list. Check the monitor configuration file to see if a flag is
required for gathering additional metrics.

```yaml

metricsToInclude:
  - metricNames:
    - win_services.startup_mode
    - win_services.state
    monitorType: telegraf/win_services
```



