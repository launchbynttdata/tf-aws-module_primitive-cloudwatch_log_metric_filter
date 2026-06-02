# TF AWS Module Primitive - CloudWatch Log Metric Filter

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)

## Overview

This Terraform module creates an [AWS CloudWatch log metric filter](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_metric_filter) that publishes custom metrics from log events.

## Pre-Commit Hooks

[.pre-commit-config.yaml](.pre-commit-config.yaml) defines pre-commit hooks for Terraform, Go, and common linting. The `commitlint` hook enforces conventional commit format. The `detect-secrets-hook` prevents new secrets from being introduced. See [pre-commit](https://pre-commit.com/#install) for installation. Install the commit-msg hook manually:

```
pre-commit install --hook-type commit-msg
```

## Usage

See [examples/complete](examples/complete) for a full working example.

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.10 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 5.100, < 7.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 6.47.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_cloudwatch_log_metric_filter.metric_filter](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_metric_filter) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_name"></a> [name](#input\_name) | Name of the CloudWatch log metric filter. Must be unique within the log group. | `string` | n/a | yes |
| <a name="input_pattern"></a> [pattern](#input\_pattern) | Filter pattern for extracting metric data from log events. See CloudWatch Logs filter and pattern syntax. | `string` | n/a | yes |
| <a name="input_log_group_name"></a> [log\_group\_name](#input\_log\_group\_name) | Name of the CloudWatch log group to associate with the metric filter. | `string` | n/a | yes |
| <a name="input_metric_transformation"></a> [metric\_transformation](#input\_metric\_transformation) | name = CloudWatch metric name created from the filter<br/>namespace = CloudWatch metric namespace<br/>value = Value to emit (default "1")<br/>default\_value = Value when no match (optional)<br/>unit = Standard unit for the metric (optional)<br/>dimensions = Map of dimension names to log field names (optional) | <pre>object({<br/>    name          = string<br/>    namespace     = string<br/>    value         = optional(string, "1")<br/>    default_value = optional(string)<br/>    unit          = optional(string)<br/>    dimensions    = optional(map(string))<br/>  })</pre> | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_id"></a> [id](#output\_id) | The ID of the metric filter (same as the name). |
| <a name="output_name"></a> [name](#output\_name) | The name of the metric filter. |
| <a name="output_log_group_name"></a> [log\_group\_name](#output\_log\_group\_name) | The name of the CloudWatch log group. |
| <a name="output_pattern"></a> [pattern](#output\_pattern) | The filter pattern. |
| <a name="output_metric_name"></a> [metric\_name](#output\_metric\_name) | The CloudWatch metric name from the metric transformation. |
| <a name="output_metric_namespace"></a> [metric\_namespace](#output\_metric\_namespace) | The CloudWatch metric namespace from the metric transformation. |
| <a name="output_metric_value"></a> [metric\_value](#output\_metric\_value) | The value emitted by the metric transformation. |
<!-- END_TF_DOCS -->
