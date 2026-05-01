terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  access_key = "test"
  secret_key = "test"
  region = "us-east-1"
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  
  endpoints {
    apigateway     = "http://192.168.49.2:30002"
    apigatewayv2   = "http://192.168.49.2:30002"
    cloudformation = "http://192.168.49.2:30002"
    cloudwatch     = "http://192.168.49.2:30002"
    dynamodb       = "http://192.168.49.2:30002"
    ec2            = "http://192.168.49.2:30002"
    es             = "http://192.168.49.2:30002"
    elasticache    = "http://192.168.49.2:30002"
    firehose       = "http://192.168.49.2:30002"
    iam            = "http://192.168.49.2:30002"
    kinesis        = "http://192.168.49.2:30002"
    lambda         = "http://192.168.49.2:30002"
    rds            = "http://192.168.49.2:30002"
    redshift       = "http://192.168.49.2:30002"
    route53        = "http://192.168.49.2:30002"
    s3             = "http://s3.192.168.49.2.localstack.cloud:30002"
    secretsmanager = "http://192.168.49.2:30002"
    ses            = "http://192.168.49.2:30002"
    sns            = "http://192.168.49.2:30002"
    sqs            = "http://192.168.49.2:30002"
    ssm            = "http://192.168.49.2:30002"
    stepfunctions  = "http://192.168.49.2:30002"
    sts            = "http://192.168.49.2:30002"
  }
}