import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';

const generalTopic = 'general';

const _channel = AndroidNotificationChannel(
  'general_channel',
  'General',
  description: 'ประกาศทั่วไปจาก TuneTrend',
  importance: Importance.high,
);

final _localNotifications = FlutterLocalNotificationsPlugin();

/// Runs in a separate isolate when a data/notification message arrives while
/// the app is backgrounded or terminated — must stay a top-level function.
@pragma('vm:entry-point')
Future<void> firebaseMessagingBackgroundHandler(RemoteMessage message) async {
  debugPrint('🔔 [Messaging] background message: ${message.messageId}');
}

class MessagingService {
  MessagingService._();

  static Future<void> initialize() async {
    await FirebaseMessaging.instance.requestPermission();
    await FirebaseMessaging.instance.subscribeToTopic(generalTopic);

    await _localNotifications.initialize(
      settings: const InitializationSettings(
        android: AndroidInitializationSettings('@mipmap/ic_launcher'),
        iOS: DarwinInitializationSettings(),
      ),
    );
    await _localNotifications
        .resolvePlatformSpecificImplementation<
          AndroidFlutterLocalNotificationsPlugin
        >()
        ?.createNotificationChannel(_channel);

    FirebaseMessaging.onMessage.listen(_showForegroundNotification);
  }

  static Future<void> _showForegroundNotification(RemoteMessage message) {
    final notification = message.notification;
    if (notification == null) return Future.value();

    return _localNotifications.show(
      id: notification.hashCode,
      title: notification.title,
      body: notification.body,
      notificationDetails: NotificationDetails(
        android: AndroidNotificationDetails(
          _channel.id,
          _channel.name,
          channelDescription: _channel.description,
          importance: Importance.high,
          priority: Priority.high,
        ),
        iOS: const DarwinNotificationDetails(),
      ),
    );
  }
}
