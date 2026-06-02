# with_cake

This demonstrates a minimum example module that produces some output for the tests to validate.

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.0 |
| <a name="requirement_random"></a> [random](#requirement\_random) | ~> 3.6 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_cake_prefix"></a> [cake\_prefix](#module\_cake\_prefix) | ../.. | n/a |
| <a name="module_cake_suffix"></a> [cake\_suffix](#module\_cake\_suffix) | ../.. | n/a |

## Resources

| Name | Type |
|------|------|
| [random_integer.cake_pos](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/integer) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_length"></a> [length](#input\_length) | Length of the random string to generate. | `number` | `24` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_string"></a> [string](#output\_string) | The random string to be generated (with cake). |
<!-- END_TF_DOCS -->
