resource "aws_dynamodb_table" "ConfigTable" {
  name           = "ConfigTable"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "org_id"
  range_key      = "config_id"

  attribute {
    name = "org_id"
    type = "S"
  }

  attribute {
    name = "config_id"
    type = "S"
  }

  attribute {
    name = "filters"
    type = "S"
  }

  global_secondary_index {
    name               = "GSI-Index"
    hash_key           = "org_id"
    range_key          = "filters"
    projection_type    = "ALL"
    write_capacity     = 1 
    read_capacity      = 1 
  }

  ttl {
    attribute_name = "" 
    enabled        = false
  }

  tags = {
    Name        = "ConfigTable"
    Environment = "Development"
    Project     = "LedgerSystem"
  }
}

output "dynamodb_table_arn" {
  value = aws_dynamodb_table.ConfigTable.arn
}

output "dynamodb_table_name" {
  value = aws_dynamodb_table.ConfigTable.name
}