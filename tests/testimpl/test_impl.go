package testimpl

import (
	"context"
	"fmt"
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
		pattern := terraform.Output(t, opts, "pattern")

		logsClient := getCloudWatchLogsClient(t, region)
		cwClient := getCloudWatchClient(t, region)
		streamName := fmt.Sprintf("terratest-%d", time.Now().UnixNano())
		logMessage := "ERROR example event for metric filter test"

		_, err := logsClient.CreateLogStream(context.TODO(), &cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(streamName),
		})
		require.NoError(t, err, "CreateLogStream should succeed")

		eventTime := time.Now()
		_, err = logsClient.PutLogEvents(context.TODO(), &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(streamName),
			LogEvents: []cwltypes.InputLogEvent{
				{
					Message:   aws.String(logMessage),
					Timestamp: aws.Int64(eventTime.UnixMilli()),
				},
			},
		})
		require.NoError(t, err, "PutLogEvents should succeed")

		filterStart := eventTime.Add(-2 * time.Minute).UnixMilli()
		var foundLog bool
		for i := 0; i < 24; i++ {
			output, filterErr := logsClient.FilterLogEvents(context.TODO(), &cloudwatchlogs.FilterLogEventsInput{
				LogGroupName:   aws.String(logGroupName),
				LogStreamNames: []string{streamName},
				FilterPattern:  aws.String(pattern),
				StartTime:      aws.Int64(filterStart),
			})
			require.NoError(t, filterErr, "FilterLogEvents should succeed")
			for _, event := range output.Events {
				if strings.Contains(aws.ToString(event.Message), logMessage) {
					foundLog = true
					break
				}
			}
			if foundLog {
				break
			}
			time.Sleep(5 * time.Second)
		}
		require.True(t, foundLog, "matching log event should be searchable after PutLogEvents")

		assert.Eventually(t, func() bool {
			endTime := time.Now()
			startTime := endTime.Add(-30 * time.Minute)
			stats, statsErr := cwClient.GetMetricStatistics(context.TODO(), &cloudwatch.GetMetricStatisticsInput{
				Namespace:  aws.String(metricNamespace),
				MetricName: aws.String(metricName),
				StartTime:  aws.Time(startTime),
				EndTime:    aws.Time(endTime),
				Period:     aws.Int32(60),
				Statistics: []cwtypes.Statistic{cwtypes.StatisticSum},
			})
			if statsErr != nil {
				return false
			}
			for _, dp := range stats.Datapoints {
				if aws.ToFloat64(dp.Sum) >= 1.0 {
					return true
				}
			}
			return false
		}, 12*time.Minute, 20*time.Second, "metric sum should be published after matching log event")
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
