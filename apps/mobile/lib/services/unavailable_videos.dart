import 'package:flutter/foundation.dart';

class UnavailableVideos extends ValueNotifier<Set<String>> {
  UnavailableVideos._() : super(const {});

  static final instance = UnavailableVideos._();

  bool contains(String videoId) => value.contains(videoId);

  void mark(String videoId) {
    if (videoId.isEmpty || value.contains(videoId)) return;
    value = {...value, videoId};
  }
}
