// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

output "region" {
  description = "The AWS region where resources are deployed."
  value       = data.aws_region.current.id
}

output "log_group_name" {
  description = "The example log group name."
  value       = aws_cloudwatch_log_group.example.name
}

output "kms_key_arn" {
  description = "The KMS key ARN used for log group encryption."
  value       = aws_kms_key.logs.arn
}

output "id" {
  description = "The metric filter ID."
  value       = module.metric_filter.id
}

output "name" {
  description = "The metric filter name."
  value       = module.metric_filter.name
}

output "pattern" {
  description = "The filter pattern."
  value       = module.metric_filter.pattern
}

output "metric_name" {
  description = "The CloudWatch metric name."
  value       = module.metric_filter.metric_name
}

output "metric_namespace" {
  description = "The CloudWatch metric namespace."
  value       = module.metric_filter.metric_namespace
}

output "metric_value" {
  description = "The metric transformation value."
  value       = module.metric_filter.metric_value
}
