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

package test

import (
	"testing"

	"github.com/launchbynttdata/lcaf-component-terratest/lib"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/launchbynttdata/tf-aws-module_primitive-cloudwatch_log_metric_filter/tests/testimpl"
)

const (
	// Terratest only exercises this example folder; other examples need their own test package or constant.
	testConfigsExamplesFolderDefault = "../../examples/complete"
	infraTFVarFileNameDefault        = "test.tfvars"
)

func TestCloudWatchLogMetricFilter(t *testing.T) {
	ctx := types.CreateTestContextBuilder().
		SetTestConfig(&testimpl.ThisTFModuleConfig{}).
		SetTestConfigFolderName(testConfigsExamplesFolderDefault).
		SetTestConfigFileName(infraTFVarFileNameDefault).
		Build()

	lib.RunSetupTestTeardown(t, *ctx, testimpl.TestComposableComplete)
}
