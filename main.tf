resource "google_pubsub_topic" "new-instance-topic" {
  project = var.project_name
  name    = var.topic
}

resource "google_logging_project_sink" "new_instance_sink" {
  project                = var.project_name
  name                   = "new_instance_sink"
  destination            = "pubsub.googleapis.com/projects/${var.project_name}/topics/${var.topic}"
  filter                 = "protoPayload.methodName=v1.compute.instances.insert OR protoPayload.methodName=beta.compute.instances.insert"
  unique_writer_identity = true
  depends_on             = [google_pubsub_topic.new-instance-topic]
}

resource "google_pubsub_topic_iam_member" "new_instance_writer" {
  project = var.project_name
  role    = "roles/pubsub.publisher"
  topic   = "new-instance-event"
  member  = google_logging_project_sink.new_instance_sink.writer_identity
}

data "archive_file" "labelzip" {
  type        = "zip"
  source_dir  = "${path.module}/src/"
  output_path = "${path.module}/label-function.zip"
}

resource "google_storage_bucket_object" "gslabelzip" {
  name   = "label-costing/${data.archive_file.labelzip.output_base64sha256}.zip"
  source = data.archive_file.labelzip.output_path
  bucket = var.bucket
}

resource "google_cloudfunctions_function" "label" {
  project               = var.project_name
  name                  = var.function_name
  description           = "Labels VMs with instance-{name,id}"
  region                = var.region
  source_archive_bucket = var.bucket
  source_archive_object = google_storage_bucket_object.gslabelzip.name
  entry_point           = "Label"
  runtime               = "go111"
  event_trigger {
    event_type = "google.pubsub.topic.publish"
    resource   = var.topic
  }
  depends_on = [google_pubsub_topic.new-instance-topic]
}

// https://www.terraform.io/docs/providers/google/r/cloudfunctions_cloud_function_iam.html

resource "google_cloudfunctions_function_iam_member" "member" {
  project        = google_cloudfunctions_function.label.project
  region         = google_cloudfunctions_function.label.region
  cloud_function = google_cloudfunctions_function.label.name
  role           = "roles/compute.instanceAdmin.v1"
  member         = "allUsers"
}
