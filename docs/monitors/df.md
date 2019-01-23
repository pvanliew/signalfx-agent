<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# df


This monitor reports metrics about free disk space.

```yaml
monitors:
 - type: df
```


Monitor Type: `df`

[Monitor Source Code](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/df)

**Accepts Endpoints**: No

**Multiple Instances Allowed**: **No**

## Configuration

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `procFSPath` | no | `string` | The path to the proc filesystem. Useful to override in containerized environments.  (Does not apply to windows) (**default:** `/proc`) |
| `ignoreSelected` | no | `bool` | If true, the filesystems selected by `fsTypes` and `mountPoints` will be excluded and all others included. |
| `fsTypes` | no | `list of string` | The filesystem types to include/exclude. (**default:** `[aufs overlay tmpfs proc sysfs nsfs cgroup devpts selinuxfs devtmpfs debugfs mqueue hugetlbfs securityfs pstore binfmt_misc autofs]`) |
| `mountPoints` | no | `list of string` | The mount paths to include/exclude, is interpreted as a regex if surrounded by `/`.  Note that you need to include the full path as the agent will see it, irrespective of the hostFSPath option. (**default:** `[/^/var/lib/docker/containers/ /^/var/lib/rkt/pods/ /^/net// /^/smb//]`) |
| `reportByDevice` | no | `bool` |  (**default:** `false`) |
| `reportInodes` | no | `bool` |  (**default:** `false`) |




## Metrics

The following table lists the metrics available for this monitor. Metrics that are not marked as Custom are standard metrics and are monitored by default.

| Name | Type | Custom | Description |
| ---  | ---  | ---    | ---         |
| `df_complex.free` | gauge |  | Free disk space in bytes |
| `df_complex.used` | gauge |  | Used disk space in bytes |
| `df_inodes.free` | gauge | X | (Linux Only) Number of inodes that are free. |
| `df_inodes.used` | gauge | X | (Linux Only) Number of inodes that are used. |
| `disk.summary_utilization` | gauge |  | Percent of disk space utilized on all volumes on this host. This metric reports with plugin dimension set to "system-utilization". |
| `disk.utilization` | gauge |  | Percent of disk used on this volume. This metric reports with plugin dimension set to "system-utilization". |
| `percent_bytes.free` | gauge | X | Free disk space on the file system, expressed as a percentage. |
| `percent_bytes.used` | gauge | X | Used disk space on the file system, expressed as a percentage. |
| `percent_inodes.free` | gauge | X | (Linux Only) Free inodes on the file system, expressed as a percentage. |
| `percent_inodes.used` | gauge | X | (Linux Only) Used inodes on the file system, expressed as a percentage. |


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
    - df_inodes.free
    - df_inodes.used
    - percent_bytes.free
    - percent_bytes.used
    - percent_inodes.free
    - percent_inodes.used
    monitorType: df
```



