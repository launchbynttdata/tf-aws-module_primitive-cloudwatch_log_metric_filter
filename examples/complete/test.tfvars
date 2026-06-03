resource_names_map = {
  "kms_key" = {
    name       = "kmskey1"
    max_length = 64
  }
  "log_group" = {
    name       = "loggroup1"
    max_length = 64
  }
  "metric_filter" = {
    name       = "metricfilter1"
    max_length = 512
  }
}

logical_product_family  = "launch"
logical_product_service = "cloudwatch"
class_env               = "dev"
instance_env            = 1
instance_resource       = 1

pattern = "ERROR"

metric_transformation = {
  name      = "ExampleErrorCount"
  namespace = "Launch/CloudWatchMetricFilterTest"
  value     = "1"
}
