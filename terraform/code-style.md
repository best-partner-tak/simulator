# Terraform Style Guide

**Table of Contents**

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
- [Introduction](#introduction)
- [Directory Structure](#directory-structure)
- [Syntax](#syntax)
  - [Spacing](#spacing)
  - [Resource Block Alignment](#resource-block-alignment)
  - [Comments](#comments)
  - [Organizing Variables](#organizing-variables)
  - [Naming Conventions](#naming-conventions)
    - [File Names](#file-names)
    - [Parameter, Meta-parameter and User Variable Naming](#parameter-meta-parameter-and-user-variable-naming)
    - [Resource Naming](#resource-naming)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Introduction

This outlines coding conventions for Terraform's HashiCorp Configuration Language (HCL). Terraform allows infrastructure to be described as code. As such, we should adhere to a style guide to ensure readable and high quality code.

# Directory Structure

Directories are split into __deployments__ and __modules__

__deployments__ contains subdirectories for each terraform deployment
__modules__ contains subdirectories for each cloud vendor as appropriate (i.e AWS, Azure, GCP)

Taking __AWS__ as an example __modules__ subdirectory, the __AWS__ directory will then contains modules logically segregated into their core function (Networking, SecurityGroups etc).

# Syntax

- Strings are in double-quotes.

## Spacing

Use 2 spaces when defining resources except when defining inline policies or other inline resources.

```
resource "aws_iam_role" "iam_role" {
  name = "${var.resource_name}-role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}
```

## Resource Block Alignment

Parameter definitions in a resource block should be aligned. The `terraform fmt` command can do this for you.

```
provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region     = "us-east-1"
}
```


## Comments

When commenting use two "//" and a space in front of the comment.

```
// CREATE ELK IAM ROLE 
...
```

## Organizing Variables

The `variables.tf` file should be broken down into three sections with each section arranged alphabetically. Starting at the top of the file:

1. Variables that have no defaults defined
2. Variables that contain defaults
3. All locals blocks 

For example:

```
variable "image_tag" {}

variable "desired_count" {
  default = "2"
}

locals {
  domain_name = "${data.terraform_remote_state.account.domain_name}"
}
```

## Naming Conventions

### File Names

Create a separate resource file for each type of AWS resource. Similar resources should be defined in the same file.

```
main.tf
providers.tf
variables.tf
```

### Parameter, Meta-parameter and User Variable Naming

 __Only use an underscore (`_`) when naming Terraform resources like TYPE/NAME parameters and user provided variables.__
 
```
resource "aws_security_group" "security_group" {
...
```

__Variables provided as output from modules should following CamelCase convention to make module provided variables easy to identify__

```
output "ControlPlaneSecurityGroupID" {
...
```

### Resource Naming

__Only use a hyphen (`-`) when naming the component being created.__

```
resource "aws_security_group" "security_group" {
  name = "${var.resource_name}-security-group"
...
```

__A resource's NAME should describe TYPE pre-pended with simulator, unique ident and minus the provider.  Common shorthand for TYPE is fine here__

```
resource "aws_security_group" "simulator_controlplane_sg" {
...
```



