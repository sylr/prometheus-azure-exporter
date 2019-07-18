prometheus-azure-exporter
=========================

[![Build Status](https://travis-ci.com/sylr/prometheus-azure-exporter.svg?branch=master)](https://travis-ci.com/sylr/prometheus-azure-exporter)
[![Docker Repository on Quay](https://quay.io/repository/sylr/prometheus-azure-exporter/status "Docker Repository on Quay")](https://quay.io/repository/sylr/prometheus-azure-exporter)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fsylr%2Fprometheus-azure-exporter.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fsylr%2Fprometheus-azure-exporter?ref=badge_shield)


This is a daemon which calls Azure API to fetch resources metrics and expose them
with HTTP using the prometheus format.

History
-------

After several incidents in Production with Azure Batch we decided that we needed something better
in terms of monitoring than what Microsoft is currently proposing.

Disclaimer
----------

This is my 2nd Go project so It is far from being perfect in terms of design and implementation.

You are very welcome to open issues and pull requests if you want to improve it.

Azure resources
---------------

| Namespaces              | Metrics                                         | Labels
|-------------------------|-------------------------------------------------|--------------------------------------------------------
| Azure                   | azure_api_calls_total                           |
|                         | azure_api_calls_failed_total                    |
|                         | azure_api_calls_duration_seconds                |
|                         | azure_api_calls_duration_sum                    |
|                         | azure_api_calls_duration_count                  |
|                         | azure_api_calls_failed_total                    |
|                         | azure_api_batch_calls_total                     | subscription, resource_group, account
|                         | azure_api_batch_calls_failed_total              | subscription, resource_group, account
|                         | azure_api_batch_calls_duration_seconds_bucket   | subscription, resource_group, account
|                         | azure_api_batch_calls_duration_seconds_sum      | subscription, resource_group, account
|                         | azure_api_batch_calls_duration_seconds_count    | subscription, resource_group, account
|                         | azure_api_graph_calls_total                     |
|                         | azure_api_graph_calls_failed_total              |
|                         | azure_api_graph_calls_duration_seconds_bucket   |
|                         | azure_api_graph_calls_duration_seconds_sum      |
|                         | azure_api_graph_calls_duration_seconds_count    |
|                         | azure_api_read_rate_limit_remaining             | subscription
|                         | azure_api_storage_calls_total                   | subscription, resource_group, account
|                         | azure_api_storage_calls_failed_total            | subscription, resource_group, account
|                         | azure_api_storage_calls_duration_seconds_bucket | subscription, resource_group, account
|                         | azure_api_storage_calls_duration_seconds_sum    | subscription, resource_group, account
|                         | azure_api_storage_calls_duration_seconds_count  | subscription, resource_group, account
| Batch                   | azure_batch_pool_quota                          | subscription, resource_group, account
|                         | azure_batch_dedicated_core_quota                | subscription, resource_group, account
|                         | azure_batch_pool_dedicated_nodes                | subscription, resource_group, account, pool
|                         | azure_batch_job_tasks_active                    | subscription, resource_group, account, job_id, job_name
|                         | azure_batch_job_tasks_running                   | subscription, resource_group, account, job_id, job_name
|                         | azure_batch_job_tasks_completed_total           | subscription, resource_group, account, job_id, job_name
|                         | azure_batch_job_tasks_succeeded_total           | subscription, resource_group, account, job_id, job_name
|                         | azure_batch_job_tasks_failed_total              | subscription, resource_group, account, job_id, job_name
| Graph                   | azure_graph_application_key_expire_time         | application, key
|                         | azure_graph_application_password_expire_time    | application, password
| Storage                 | azure_storage_blob_size_bytes_bucket            | subscription, resource_group, account, container
|                         | azure_storage_blob_size_bytes_sum               | subscription, resource_group, account, container
|                         | azure_storage_blob_size_bytes_count             | subscription, resource_group, account, container

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fsylr%2Fprometheus-azure-exporter.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fsylr%2Fprometheus-azure-exporter?ref=badge_large)