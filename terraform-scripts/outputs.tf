
output "sqs_queue_url" {
  value = aws_sqs_queue.my_queue.id
}

output "dynamodb_table_name" {
  value = aws_dynamodb_table.my_table.name
}