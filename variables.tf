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

# -----------------------------------------------------------------------------
# Required
# -----------------------------------------------------------------------------

variable "name" {
  description = "Name of the CloudWatch log metric filter. Must be unique within the log group."
  type        = string

  validation {
    condition     = length(var.name) >= 1 && length(var.name) <= 512
    error_message = "Metric filter name must be between 1 and 512 characters."
  }
}

variable "pattern" {
  description = "Filter pattern for extracting metric data from log events. See CloudWatch Logs filter and pattern syntax."
  type        = string

  validation {
    condition     = length(var.pattern) >= 1
    error_message = "Pattern must not be empty."
  }
}

variable "log_group_name" {
  description = "Name of the CloudWatch log group to associate with the metric filter."
  type        = string
}

variable "metric_transformation" {
  description = <<-EOT
    name = CloudWatch metric name created from the filter
    namespace = CloudWatch metric namespace
    value = Value to emit (default "1")
    default_value = Value when no match (optional)
    unit = Standard unit for the metric (optional)
    dimensions = Map of dimension names to log field names (optional)
  EOT
  type = object({
    name          = string
    namespace     = string
    value         = optional(string, "1")
    default_value = optional(string)
    unit          = optional(string)
    dimensions    = optional(map(string))
  })
}
