# TF AWS Module Primitive - CloudWatch Log Metric Filter

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)

## Overview

This Terraform module creates an [AWS CloudWatch log metric filter](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_metric_filter) that publishes custom metrics from log events.

## Usage

See [examples/complete](examples/complete) for a full working example.

## Module Development

### Pre-Requisites

The following commands should be available on your system:

- `asdf` or `mise`
- `make`
- `python3` (for pre-commit)

Additionally, your `git` user and email must be configured. Run the `make configure` command from the root of the repository to ensure that you meet these requirements.

### Pre-Commit hooks

The [.pre-commit-config.yaml](.pre-commit-config.yaml) file defines `pre-commit` hooks for Terraform formatting, validation, documentation generation, and detect-secrets. Hooks are installed when you run `make configure`. Go linting runs via `make lint` in local development and CI, not via pre-commit.

### Terratest examples

Post-deploy tests in `tests/post_deploy_functional/` and `tests/post_deploy_functional_readonly/` target `examples/complete` via an explicit folder constant in each `main_test.go`. Adding another example (for example `examples/minimal`) requires a new test entry point or updating that constant; it is not picked up automatically.

### Local Validation

You should validate the changes you make to any module locally, prior to pushing your changes in a branch to GitHub.

1. Ensure that you have run `make configure` successfully.
2. Ensure you are signed into the appropriate cloud provider (e.g. AWS) for the module under test in your current console session.
3. Run the Terraform and Golang linters:

```
make lint
```

4. Once linters pass, run integration tests (apply, test, destroy):

```
make test
```

The pre-commit validations, as well as the `make lint` and `make test` targets, are performed in CI. Running them locally before opening a PR helps ensure a smooth review.

### Review & Merge Process

Open a Pull Request to the default (`main`) branch. The PR title must follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#specification) format to merge and to drive semantic versioning.

Ensure CI workflows pass, address review feedback, and obtain approvals required by `CODEOWNERS`.

### Automatic Updates

Shared configuration and workflow files are largely managed through [launch-terraform-skeleton](https://github.com/launchbynttdata/launch-terraform-skeleton). Avoid one-off edits to copied skeleton files in this repository unless necessary (for example `.gitignore` entries for generated artifacts). Use `copier check-update` / `copier update` when refreshing from the skeleton.

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.10 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 5.100, < 7.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 6.53.0 |

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
