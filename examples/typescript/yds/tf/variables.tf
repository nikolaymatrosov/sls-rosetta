variable "ydb_connection_string" {
  description = "YDB connection string for the function"
  type        = string
}

variable "yds_topic_id" {
  description = "YDS topic name for the function"
  type        = string
}

variable "function_zip" {
  description = "Path to the zipped function code"
  type        = string
} 