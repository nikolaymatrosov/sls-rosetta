variable "folder_id" {
  description = "Yandex Cloud folder ID"
  type        = string
}

variable "function_name" {
  description = "Name of the WebSocket handler function"
  type        = string
  default     = "ws-go-handler"
}

variable "database_name" {
  description = "Name of the YDB database"
  type        = string
  default     = "ws-go-database"
}

variable "gateway_name" {
  description = "Name of the API Gateway"
  type        = string
  default     = "ws-go-gateway"
}

variable "service_account_name" {
  description = "Name of the service account"
  type        = string
  default     = "ws-go-function-sa"
}

variable "topic_name" {
  description = "Name of the YDB topic for broadcasting"
  type        = string
  default     = "broadcast-topic"
}

variable "topic_consumer_name" {
  description = "Name of the topic consumer"
  type        = string
  default     = "broadcast-consumer"
}
