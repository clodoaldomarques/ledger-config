resource "aws_dynamodb_table" "ScriptsTable" {
  name           = "ScriptsTable"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "org_id"
  range_key      = "script_id"

  attribute {
    name = "org_id"
    type = "S"
  }

  attribute {
    name = "script_id"
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
    Name        = "ScriptsTable"
    Environment = "Development"
    Project     = "AccountingSystem"
  }
}

output "dynamodb_table_arn" {
  value = aws_dynamodb_table.ScriptsTable.arn
}

output "dynamodb_table_name" {
  value = aws_dynamodb_table.ScriptsTable.name
}