gcloud functions deploy send-push-notification-pubsub \
  --entry-point send_push_notification_pubsub \
  --runtime python37 \
  --region asia-east2 \
  --trigger-topic send-push-notification \
  --set-env-vars FIREBASE_API_KEY=<production-api-key-here>
