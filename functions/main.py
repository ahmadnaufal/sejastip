import base64
import json
import os

from pyfcm import FCMNotification
from google.cloud import datastore
from datetime import datetime

ds_client = datastore.Client(namespace='Sejastip')

def send_push_notification_pubsub(event, context):
    """Triggered from a message on a Cloud Pub/Sub topic.
    Args:
         event (dict): Event payload.
         context (google.cloud.functions.Context): Metadata for the event.
    """
    encoded_payload = base64.b64decode(event['data']).decode('utf-8')
    notification_data = json.loads(encoded_payload)
    # ideally, notification data will contain
    # {
    #   device:
    #   user_id:
    #   data: {
    #     title:
    #     content:
    #   }
    # }
    api_key = os.getenv("FIREBASE_API_KEY")
    push_service = FCMNotification(api_key=api_key)
    result = push_service.notify_single_device(
      registration_id=notification_data['device'],
      message_title=notification_data['data']['title'],
      message_body=notification_data['data']['content']
    )

    now = datetime.now()
    key = ds_client.key('NotificationLogs')
    entity = datastore.Entity(key=key, exclude_from_indexes=['title', 'content'])
    entity.update({
      'device_target': notification_data['device'],
      'user_id': notification_data['user_id'],
      'title': notification_data['data']['title'],
      'content': notification_data['data']['content'],
      'status': 'sent',
      'created_at': now,
      'updated_at': now
    })
    ds_client.put(entity)

    print(result)
