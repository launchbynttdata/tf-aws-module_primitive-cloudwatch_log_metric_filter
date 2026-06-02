package testimpl

import (
	"regexp"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
)

// TestComposableComplete verifies the deployed resource and exercises it with write
// operations. When replacing this template with a real module, add at least one
// write operation below — for example, uploading an object to a bucket, sending a
// message to a queue, or invoking a function. The random_string resource used in this
// template has no writable cloud state, so no write operation is shown here.
func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	t.Run("TestOutputFormat", func(t *testing.T) {
		output := terraform.Output(t, ctx.TerratestTerraformOptions(), "string")

		// Verify output contains only alphanumeric characters and 🍰.
		assert.Regexp(t, regexp.MustCompile("^[A-Za-z🍰0-9]+$"), output)
	})
}

// TestComposableCompleteReadOnly verifies the deployed resource using only read
// operations. Do NOT add write operations (object uploads, message sends, API
// mutations, etc.) to this function — those belong in TestComposableComplete.
func TestComposableCompleteReadOnly(t *testing.T, ctx types.TestContext) {
	t.Run("TestOutputFormat", func(t *testing.T) {
		output := terraform.Output(t, ctx.TerratestTerraformOptions(), "string")

		// Verify output contains only alphanumeric characters and 🍰.
		assert.Regexp(t, regexp.MustCompile("^[A-Za-z🍰0-9]+$"), output)
	})
}
