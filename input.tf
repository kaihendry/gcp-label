variable project_name {
  default = "jsonunroller"
}

variable "bucket" {
  default     = "artifacts.jsonunroller.appspot.com"
  description = "The bucket where the function source is kept"
}

variable "region" {
  default     = "asia-east2"
  description = "Closest region which supports Cloud functions"
}

variable "topic" {
  default     = "new-instance-event"
  description = "Name of topic of when a new GCP VM comes online"
}

variable "function_name" {
  default     = "label"
  description = "Name of function that is subscribed to the topic above"
}
