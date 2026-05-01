resource "aws_dynamodb_table_item" "pismo-101" {
  table_name = aws_dynamodb_table.ScriptsTable.name
  hash_key   = aws_dynamodb_table.ScriptsTable.hash_key
  range_key  = aws_dynamodb_table.ScriptsTable.range_key

  item = jsonencode({
    "script_id"   = { "S" =  "pismo-101" },
    "filters"     = { "S" = "TENANT#PISMO#PROGRAMID#0000#EVENTTYPEID#101" },
    "org_id"      = { "S" = "PISMO" },
    "description_id" = { "S" = "Parcelamento xpto" },
    "level"       = { "S" = "platform" },
    "enable"      = { "BOOL" = true },
    "version"     = { "N" = "2" },
    "event_type_id" = { "S" = "101" },
    "created_at"  = { "S" = timestamp() },
    "updated_at"  = { "S" = timestamp() },
    "entries"     = { "L" : [
      {
        "M" : {
          "amount_name"    : { "S" : "amount" },
          "entry_type_id" : { "N" : "101" },
          "description"   : { "S" : "Parcelamento" },
          "flow"          : { "S" : "regular" },
        }
      },
      {
        "M" : {
          "expression"     : { "S" : "Amount.amount + Fee.iof" },
          "entry_type_id" : { "N" : "102" },
          "description"   : { "S" : "Parcelamento" },
          "flow"          : { "S" : "migration" },
        }
      }
    ]}
  })
}