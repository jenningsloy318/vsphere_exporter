# vsphere_exporter
A vsphere_exporter to get  metrics from vCenter 

## Runing
vsphere_exporter can be running in two mode:
- single vCenter mode
    ```yaml
    mode: single
    enabled_cluster: 10.36.51.11
    clusters:
        10.36.51.11:
            username: user
            password: pass
    ```
    then you can get the metrics via `http://<ip address>:9271/vsphere`

- - multi-vCenter mode
    ```yaml
    mode: single
    enabled_cluster: 10.36.51.11
    clusters:
        10.36.51.11:
            username: user
            password: pass
        10.36.51.12:
            username: user
            password: pass        
    ```
    when wanna to get the vCenter metrics, you should specify the target at the request,thus get the metrics via `http://localhost:9272/vsphere?target=10.36.51.11`


## prometheus job config

You can then setup [Prometheus] to scrape the target using something like this in your Prometheus configuration files:
- for single-vCenter mode
```yaml
  - job_name: 'vsphere-exporter'

    # metrics_path defaults to '/metrics'
    metrics_path: /vsphere
```

- for muti-vCenter mode
```yaml
  - job_name: 'vsphere-exporter'

    # metrics_path defaults to '/metrics'
    metrics_path: /vsphere

    # scheme defaults to 'http'.

    static_configs:
    - targets:
       - 10.36.48.24
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: <IP address of vsphere_exporter>:9272  ### the address of the redfish-exporter address
```

## Reference
- https://code.vmware.com/apis/358/vsphere/doc/index-mo_types.html
- https://raw.githubusercontent.com/vmware/govmomi/381aa00a0d03120e12e9a4fba08feaa757d24e5d/vim25/types/types.go

## Acknowledgement

- [1] thanks [govmomi](https://github.com/vmware/govmomi/) provide underlying library to manipulate the vcenter resources