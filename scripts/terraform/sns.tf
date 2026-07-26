variable "account-id" {
  type = string
  default = "000000000000"
}

resource "aws_sns_topic" "config-sns-topic" {
  name = "config-sns-topic"
}

output "aws_sns_topic_arn" {
  value = aws_sns_topic.config-sns-topic.arn
}

output "aws_sns_topic_name" {
  value = aws_sns_topic.config-sns-topic.name
}
