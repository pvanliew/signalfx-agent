<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# disk


This monitor reports metrics about free disk space.

```yaml
monitors:
 - type: disk
```


Monitor Type: `disk`

[Monitor Source Code](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/disk)

**Accepts Endpoints**: No

**Multiple Instances Allowed**: **No**

## Configuration

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `procFSPath` | no | `string` | The path to the proc filesystem. Useful to override in containerized environments.  (Does not apply to windows) (**default:** `/proc`) |
| `disks` | no | `list of string` | Which devices to include/exclude (**default:** `[/^loop[0-9]+$/ /^dm-[0-9]+$/]`) |
| `ignoreSelected` | no | `bool` | If true, the disks selected by `disks` will be excluded and all others included. |




## Metrics

The following table lists the metrics available for this monitor. Metrics that are not marked as Custom are standard metrics and are monitored by default.

| Name | Type | Custom | Description |
| ---  | ---  | ---    | ---         |
| `disk_merged.read` | cumulative | X | (Linux Only) The number of disk reads merged into single physical disk access operations. |
| `disk_merged.write` | cumulative | X | (Linux Only) The number of disk writes merged into single physical disk access operations. |
| `disk_octets.avg_read` | gauge | X | (Windows Only) The average number of octets (bytes) read. |
| `disk_octets.avg_write` | gauge | X | (Windows Only) The average number of octets (bytes) written. |
| `disk_octets.read` | cumulative | X | (Linux Only) The number of bytes (octets) read from a disk. |
| `disk_octets.write` | cumulative | X | (Linux Only) The number of bytes (octets) written to a disk. |
| `disk_ops.avg_read` | gauge | X | (Windows Only) The average disk read queue length. |
| `disk_ops.avg_write` | gauge | X | (Windows Only) The average disk write queue length. |
| `disk_ops.read` | cumulative |  | (Linux Only) The number of disk read operations. |
| `disk_ops.write` | cumulative |  | (Linux Only) The number of disk write operations. |
| `disk_time.avg_read` | gauge | X | (Windows Only) The average time spent reading from the disk. |
| `disk_time.avg_write` | gauge | X | (Windows Only) The average time spent writing to the disk |
| `disk_time.read` | cumulative | X | (Linux Only) The average amount of time it took to do a read operation. |
| `disk_time.write` | cumulative | X | (Linux Only) The average amount of time it took to do a write operation. |


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
    - disk_merged.read
    - disk_merged.write
    - disk_octets.avg_read
    - disk_octets.avg_write
    - disk_octets.read
    - disk_octets.write
    - disk_ops.avg_read
    - disk_ops.avg_write
    - disk_time.avg_read
    - disk_time.avg_write
    - disk_time.read
    - disk_time.write
    monitorType: disk
```



