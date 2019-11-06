import base64
import json
import os

from pyfcm import FCMNotification

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

    print(result)
