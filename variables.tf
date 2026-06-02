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

variable "length" {
  description = "Length of the random string to generate."
  type        = number
  default     = 24

  validation {
    condition     = var.length > 0 && var.length < 100
    error_message = "Length must be a positive integer less than 100."
  }
}

variable "number" {
  description = "Whether the random string should include numbers. Defaults to true."
  type        = bool
  default     = true
}

variable "special" {
  description = "Whether the random string should include special characters. Defaults to false."
  type        = bool
  default     = false
}
