resource "aws_sqs_queue" "config-sqs-queue" {
  name = "script-sqs-queue"
}

resource "aws_sns_topic_subscription" "config-sns-sqs-subscription" {
  topic_arn = aws_sns_topic.config-sns-topic.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.config-sqs-queue.arn 
}

output "aws_sqs_queue_arn" {
  value = aws_sqs_queue.config-sqs-queue.arn
}

output "aws_sqs_queue_name" {
  value = aws_sqs_queue.config-sqs-queue.name
}

output "aws_sns_topic_subscription_confirmation" {
  value = aws_sns_topic_subscription.config-sns-sqs-subscription.confirmation_was_authenticated
}