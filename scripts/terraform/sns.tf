variable "account-id" {
  type = string
  default = "000000000000"
}

resource "aws_sns_topic" "config-sns-topic" {
  name = "script-sns-topic"
}

resource "aws_sns_topic_policy" "default" {
  arn = aws_sns_topic.config-sns-topic.arn

  policy = data.aws_iam_policy_document.sns_topic_policy.json
}

data "aws_iam_policy_document" "sns_topic_policy" {
  policy_id = "__default_policy_ID"

  statement {
    actions = [
      "SNS:Subscribe",
      "SNS:SetTopicAttributes",
      "SNS:RemovePermission",
      "SNS:Receive",
      "SNS:Publish",
      "SNS:ListSubscriptionsByTopic",
      "SNS:GetTopicAttributes",
      "SNS:DeleteTopic",
      "SNS:AddPermission",
    ]

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceOwner"

      values = [
        var.account-id,
      ]
    }

    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    resources = [
      aws_sns_topic.config-sns-topic.arn,
    ]

    sid = "__default_statement_ID"
  }
}

output "aws_sns_topic_arn" {
  value = aws_sns_topic.config-sns-topic.arn
}

output "aws_sns_topic_name" {
  value = aws_sns_topic.config-sns-topic.name
}
