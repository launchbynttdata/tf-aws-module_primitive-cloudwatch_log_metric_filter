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

func describeLogGroupWithRetry(t *testing.T, client *cloudwatchlogs.Client, logGroupName string) *cwltypes.LogGroup {
	t.Helper()
	var lastErr error
	for attempt := 1; attempt <= 12; attempt++ {
		output, err := client.DescribeLogGroups(context.TODO(), &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(logGroupName),
		})
		if err != nil {
			lastErr = err
			time.Sleep(5 * time.Second)
			continue
		}
		for i := range output.LogGroups {
			if aws.ToString(output.LogGroups[i].LogGroupName) == logGroupName {
				return &output.LogGroups[i]
			}
		}
		lastErr = fmt.Errorf("log group %s not found yet", logGroupName)
		time.Sleep(5 * time.Second)
	}
	require.Failf(t, "unable to describe log group", "%s after retries: %v", logGroupName, lastErr)
	return nil
}

func assertLogGroupKMSEncryption(t *testing.T, client *cloudwatchlogs.Client, logGroupName, expectedKMSArn string) {
	t.Helper()
	logGroup := describeLogGroupWithRetry(t, client, logGroupName)
	require.NotNil(t, logGroup.KmsKeyId, "log group should have a KMS key attached")
	assert.Equal(t, expectedKMSArn, aws.ToString(logGroup.KmsKeyId), "KMS key ARN should match")
}

func waitForMetricFilter(t *testing.T, client *cloudwatchlogs.Client, logGroupName, filterName string) (*cwltypes.MetricFilter, error) {
	t.Helper()
	var lastErr error
	for i := 0; i < 12; i++ {
		output, err := client.DescribeMetricFilters(context.TODO(), &cloudwatchlogs.DescribeMetricFiltersInput{
			LogGroupName:     aws.String(logGroupName),
			FilterNamePrefix: aws.String(filterName),
		})
		lastErr = err
		if err == nil {
			for i := range output.MetricFilters {
				if aws.ToString(output.MetricFilters[i].FilterName) == filterName {
					return &output.MetricFilters[i], nil
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
	return nil, lastErr
}

func assertMetricFilterConfig(t *testing.T, matched *cwltypes.MetricFilter, opts *terraform.Options, filterName string) {
	t.Helper()
	expectedPattern := terraform.Output(t, opts, "pattern")
	expectedMetricName := terraform.Output(t, opts, "metric_name")
	expectedNamespace := terraform.Output(t, opts, "metric_namespace")
	expectedMetricValue := terraform.Output(t, opts, "metric_value")

	require.NotNil(t, matched, "metric filter should be found by name")
	require.Equal(t, filterName, aws.ToString(matched.FilterName), "filter name should match")
	assert.Equal(t, expectedPattern, aws.ToString(matched.FilterPattern), "filter pattern should match")
	require.NotEmpty(t, matched.MetricTransformations, "metric transformations should be present")
	assert.Equal(t, expectedMetricName, aws.ToString(matched.MetricTransformations[0].MetricName), "metric name should match")
	assert.Equal(t, expectedNamespace, aws.ToString(matched.MetricTransformations[0].MetricNamespace), "metric namespace should match")
	assert.Equal(t, expectedMetricValue, aws.ToString(matched.MetricTransformations[0].MetricValue), "metric value should match")
}

func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	t.Run("VerifyTerraformOutputs", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		id := terraform.Output(t, opts, "id")
		name := terraform.Output(t, opts, "name")
		pattern := terraform.Output(t, opts, "pattern")
		metricName := terraform.Output(t, opts, "metric_name")
		metricNamespace := terraform.Output(t, opts, "metric_namespace")
		metricValue := terraform.Output(t, opts, "metric_value")

		assert.Equal(t, name, id, "id should equal name for metric filter")
		assert.NotEmpty(t, pattern, "pattern should be set")
		assert.NotEmpty(t, metricName, "metric name should be set")
		assert.NotEmpty(t, metricNamespace, "metric namespace should be set")
		assert.NotEmpty(t, metricValue, "metric value should be set")
	})

	t.Run("VerifyLogGroupKMSEncryption", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		logGroupName := terraform.Output(t, opts, "log_group_name")
		kmsKeyARN := terraform.Output(t, opts, "kms_key_arn")
		region := terraform.Output(t, opts, "region")

		client := getCloudWatchLogsClient(t, region)
		assertLogGroupKMSEncryption(t, client, logGroupName, kmsKeyARN)
	})

	t.Run("VerifyMetricFilterViaAPI", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		filterName := terraform.Output(t, opts, "name")
		logGroupName := terraform.Output(t, opts, "log_group_name")
		region := terraform.Output(t, opts, "region")

		client := getCloudWatchLogsClient(t, region)
		matched, err := waitForMetricFilter(t, client, logGroupName, filterName)
		require.NoError(t, err, "DescribeMetricFilters should succeed")
		assertMetricFilterConfig(t, matched, opts, filterName)
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

		// CloudWatch metric publication can lag several minutes after matching log events.
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

	t.Run("VerifyLogGroupKMSEncryption", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		logGroupName := terraform.Output(t, opts, "log_group_name")
		kmsKeyARN := terraform.Output(t, opts, "kms_key_arn")
		region := terraform.Output(t, opts, "region")

		client := getCloudWatchLogsClient(t, region)
		assertLogGroupKMSEncryption(t, client, logGroupName, kmsKeyARN)
	})

	t.Run("VerifyMetricFilterViaAPI", func(t *testing.T) {
		opts := ctx.TerratestTerraformOptions()
		filterName := terraform.Output(t, opts, "name")
		logGroupName := terraform.Output(t, opts, "log_group_name")
		region := terraform.Output(t, opts, "region")

		client := getCloudWatchLogsClient(t, region)
		matched, err := waitForMetricFilter(t, client, logGroupName, filterName)
		require.NoError(t, err, "DescribeMetricFilters should succeed")
		assertMetricFilterConfig(t, matched, opts, filterName)
	})
}
