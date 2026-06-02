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

resource "aws_cloudwatch_log_metric_filter" "metric_filter" {
  name           = var.name
  pattern        = var.pattern
  log_group_name = var.log_group_name

  metric_transformation {
    name          = var.metric_transformation.name
    namespace     = var.metric_transformation.namespace
    value         = var.metric_transformation.value
    default_value = var.metric_transformation.default_value
    unit          = var.metric_transformation.unit
    dimensions    = var.metric_transformation.dimensions
  }
}
