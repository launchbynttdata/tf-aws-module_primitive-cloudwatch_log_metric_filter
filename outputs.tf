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

output "id" {
  description = "The ID of the metric filter (same as the name)."
  value       = aws_cloudwatch_log_metric_filter.metric_filter.id
}

output "name" {
  description = "The name of the metric filter."
  value       = aws_cloudwatch_log_metric_filter.metric_filter.name
}

output "log_group_name" {
  description = "The name of the CloudWatch log group."
  value       = aws_cloudwatch_log_metric_filter.metric_filter.log_group_name
}

output "pattern" {
  description = "The filter pattern."
  value       = aws_cloudwatch_log_metric_filter.metric_filter.pattern
}

output "metric_name" {
  description = "The CloudWatch metric name from the metric transformation."
  value       = aws_cloudwatch_log_metric_filter.metric_filter.metric_transformation[0].name
}

output "metric_namespace" {
  description = "The CloudWatch metric namespace from the metric transformation."
  value       = aws_cloudwatch_log_metric_filter.metric_filter.metric_transformation[0].namespace
}

output "metric_value" {
  description = "The value emitted by the metric transformation."
  value       = aws_cloudwatch_log_metric_filter.metric_filter.metric_transformation[0].value
}
