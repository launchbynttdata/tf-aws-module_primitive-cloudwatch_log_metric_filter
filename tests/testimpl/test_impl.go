package testimpl

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cwltypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getCloudWatchLogsClient(t *testing.T, region string) *cloudwatchlogs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	require.NoError(t, err, "unable to load AWS config")
	return cloudwatchlogs.NewFromConfig(cfg)
}

func getCloudWatchClient(t *testing.T, region string) *cloudwatch.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	require.NoError(t, err, "unable to load AWS config")
	return cloudwatch.NewFromConfig(cfg)
}

func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	t.Run("VerifyTerraformOutputs", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		id := terraform.Output(t, opts, "id")
		name := terraform.Output(t, opts, "name")
		pattern := terraform.Output(t, opts, "pattern")
		metricName := terraform.Output(t, opts, "metric_name")
		metricNamespace := terraform.Output(t, opts, "metric_namespace")

		assert.Equal(t, name, id, "id should equal name for metric filter")
		assert.Equal(t, "ERROR", pattern, "pattern should match example")
		assert.Equal(t, "ExampleErrorCount", metricName, "metric name should match example")
		assert.Equal(t, "Launch/CloudWatchMetricFilterTest", metricNamespace, "metric namespace should match example")
	})

	t.Run("VerifyMetricFilterViaAPI", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		filterName := terraform.Output(t, opts, "name")
		logGroupName := terraform.Output(t, opts, "log_group_name")
		region := terraform.Output(t, opts, "region")
		expectedMetricName := terraform.Output(t, opts, "metric_name")
		expectedNamespace := terraform.Output(t, opts, "metric_namespace")

		client := getCloudWatchLogsClient(t, region)

		var filters []cwltypes.MetricFilter
		var err error
		for i := 0; i < 12; i++ {
			output, describeErr := client.DescribeMetricFilters(context.TODO(), &cloudwatchlogs.DescribeMetricFiltersInput{
				LogGroupName: aws.String(logGroupName),
				FilterNamePrefix: aws.String(filterName),
			})
			err = describeErr
			if err == nil && len(output.MetricFilters) > 0 {
				filters = output.MetricFilters
				break
			}
			time.Sleep(5 * time.Second)
		}
		require.NoError(t, err, "DescribeMetricFilters should succeed")
		require.NotEmpty(t, filters, "metric filter should exist")

		var matched *cwltypes.MetricFilter
		for i := range filters {
			if aws.ToString(filters[i].FilterName) == filterName {
				matched = &filters[i]
				break
			}
		}
		require.NotNil(t, matched, "metric filter should be found by name")
		require.NotEmpty(t, matched.MetricTransformations, "metric transformations should be present")
		assert.Equal(t, expectedMetricName, aws.ToString(matched.MetricTransformations[0].MetricName), "metric name should match")
		assert.Equal(t, expectedNamespace, aws.ToString(matched.MetricTransformations[0].MetricNamespace), "metric namespace should match")
	})

	t.Run("PutLogEventsAndVerifyMetric", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		logGroupName := terraform.Output(t, opts, "log_group_name")
		region := terraform.Output(t, opts, "region")
		metricName := terraform.Output(t, opts, "metric_name")
		metricNamespace := terraform.Output(t, opts, "metric_namespace")

		logsClient := getCloudWatchLogsClient(t, region)
		cwClient := getCloudWatchClient(t, region)
		streamName := "terratest-stream"

		_, err := logsClient.CreateLogStream(context.TODO(), &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(streamName),
		})
		if err != nil && !strings.Contains(err.Error(), "ResourceAlreadyExistsException") {
			require.NoError(t, err, "CreateLogStream should succeed")
		}

		_, err = logsClient.PutLogEvents(context.TODO(), &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(streamName),
			LogEvents: []cwltypes.InputLogEvent{
				{
					Message:   aws.String("ERROR example event for metric filter test"),
					Timestamp: aws.Int64(time.Now().UnixMilli()),
				},
			},
		})
		require.NoError(t, err, "PutLogEvents should succeed")

		endTime := time.Now()
		startTime := endTime.Add(-10 * time.Minute)
		var datapoints []cwtypes.Datapoint
		for i := 0; i < 18; i++ {
			stats, err := cwClient.GetMetricStatistics(context.TODO(), &cloudwatch.GetMetricStatisticsInput{
				Namespace:  aws.String(metricNamespace),
				MetricName: aws.String(metricName),
				StartTime:  aws.Time(startTime),
				EndTime:    aws.Time(endTime),
				Period:     aws.Int32(60),
				Statistics: []cwtypes.Statistic{cwtypes.StatisticSum},
			})
			require.NoError(t, err, "GetMetricStatistics should succeed")
			if len(stats.Datapoints) > 0 {
				datapoints = stats.Datapoints
				break
			}
			time.Sleep(10 * time.Second)
		}
		require.NotEmpty(t, datapoints, "metric datapoints should be published after matching log event")
		assert.GreaterOrEqual(t, aws.ToFloat64(datapoints[0].Sum), 1.0, "metric sum should be at least 1")
	})
}

func TestComposableCompleteReadOnly(t *testing.T, ctx types.TestContext) {
	t.Run("VerifyTerraformOutputs", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		id := terraform.Output(t, opts, "id")
		name := terraform.Output(t, opts, "name")

		assert.Equal(t, name, id, "id should equal name for metric filter")
	})

	t.Run("VerifyMetricFilterExistsViaAPI", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		filterName := terraform.Output(t, opts, "name")
		logGroupName := terraform.Output(t, opts, "log_group_name")
		region := terraform.Output(t, opts, "region")

		client := getCloudWatchLogsClient(t, region)

		output, err := client.DescribeMetricFilters(context.TODO(), &cloudwatchlogs.DescribeMetricFiltersInput{
			LogGroupName:     aws.String(logGroupName),
			FilterNamePrefix: aws.String(filterName),
		})
		require.NoError(t, err, "DescribeMetricFilters should succeed")

		found := false
		for _, filter := range output.MetricFilters {
			if aws.ToString(filter.FilterName) == filterName {
				found = true
				break
			}
		}
		assert.True(t, found, "metric filter should exist")
	})
}
