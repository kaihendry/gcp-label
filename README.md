Need: Difficult to identity which instances are having the biggest impact on
your GCP bill.

Feature: Function to label compute instances by instance {name,id} so they can
be viewed inside a [billing
export](https://cloud.google.com/blog/products/gcp/use-labels-to-gain-visibility-into-gcp-resource-usage-and-spending)
or primitively from the group spending by labels in the Billing > Reports page.

Benefit: Now you can better identify costly instances.

# How to deploy

Install [Terraform](https://www.terraform.io/)

Edit `input.tf` to reflect your project name aka PROJECT_ID from `gcloud
projects list` and bucket to stage the function.

	terraform {init,plan,apply}
