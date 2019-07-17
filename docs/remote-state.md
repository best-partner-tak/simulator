# Simulator Terraform remote State

It is highly recommended to store Terraform state in a S3 bucket to avoid
having to maually tidy up the environment should you exit from the Simulator
container.

# Creating S3 bucket

You can either create a bucket via the AWS console, or use the aws cli, or use
the supplied helper [script](../scripts/create-terraform-s3-bucket)

The script takes the following command line arguments
```
-b [name for the bucket, which must be globally unique]
-p [name of aws profile - only required if not using the default profile]
```

# Configure Terraform backend storage

Once this has been created you then need to edit the
[provider.tf](../terraform/deployments/AWS/providers.tf) file and modify the
terraform section as below:

```
terraform {
  backend "s3" {
    key = "simulator.tfstate"
    bucket = "add bucket name"
    encrypt = false # Optional, S3 Bucket Server Side Encryption
  }
}
```

