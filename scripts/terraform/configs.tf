resource "aws_dynamodb_table_item" "ledger-101" {
  table_name = aws_dynamodb_table.ConfigTable.name
  hash_key   = aws_dynamodb_table.ConfigTable.hash_key
  range_key  = aws_dynamodb_table.ConfigTable.range_key

  item = jsonencode({
    "config_id"   = { "S" =  "ledger-101" },
    "filters"     = { "S" = "TENANT#LEDGER#PROGRAMID#0000#EVENTTYPEID#101" },
    "org_id"      = { "S" = "LEDGER" },
    "description" = { "S" = "Compra a vista" },
    "level"       = { "S" = "platform" },
    "enable"      = { "BOOL" = true },
    "version"     = { "N" = "2" },
    "process_code" = { "S" = "101" },
    "created_at"  = { "S" = timestamp() },
    "updated_at"  = { "S" = timestamp() },
    "scripts"     = { "L" : [
      {
        "M" : {
          "expression"    : { "S" : "Amount.amount + Fee.iof" },
          "script_id"     : { "N" : "101" },
          "description"   : { "S" : "Compra a vista - Cartão" },
          "flow"          : { "S" : "regular" },
        }
      },
      {
        "M" : {
          "expression"     : { "S" : "Amount.amount" },
          "script_id"      : { "N" : "101" },
          "description"    : { "S" : "Compra a vista - PIX" },
          "flow"           : { "S" : "migration" },
        }
      }
    ]}
  })
}