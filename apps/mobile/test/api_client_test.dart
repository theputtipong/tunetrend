import 'package:flutter_test/flutter_test.dart';
import 'package:tunetrend_mobile/services/api_client.dart';

import 'helpers/fake_http.dart';

void main() {
  late FakeAdapter adapter;
  late ApiClient api;

  setUp(() {
    adapter = FakeAdapter();
    api = ApiClient.withDio(fakeDio(adapter));
  });

  group('local search', () {
    test(
      'parses items and suggestYouTube, defaults category to music',
      () async {
        adapter.on('GET', '/search', 200, {
          'success': true,
          'data': {
            'items': [
              {
                'id': 'v1',
                'title': 'Rain',
                'channelTitle': 'Band',
                'thumbnailUrl': 'https://i.ytimg.com/vi/v1/hq.jpg',
                'viewCount': '12',
                'publishedAt': '2026-11-01T00:00:00Z',
                'videoType': 'MV',
                'categoryId': '10',
                'source': 'chart',
                'chartTab': 'trending',
                'chartRank': 3,
              },
            ],
            'suggestYouTube': true,
          },
        });

        final result = await api.search('rain', 'TH');

        expect(result.suggestYouTube, isTrue);
        expect(result.items.single.id, 'v1');
        expect(result.items.single.chartRank, 3);
        expect(result.items.single.countryCode, 'TH');
        expect(adapter.requests.single.query, {
          'q': 'rain',
          'country': 'TH',
          'category': '10',
        });
      },
    );

    test('400 surfaces server message as ApiException', () async {
      adapter.on('GET', '/search', 400, {
        'success': false,
        'error': 'q must be 2-100 characters',
      });

      expect(
        () => api.search('r', 'TH'),
        throwsA(
          isA<ApiException>()
              .having((e) => e.statusCode, 'status', 400)
              .having((e) => e.message, 'message', contains('2-100')),
        ),
      );
    });
  });

  group('YouTube search error mapping', () {
    setUp(() => api.setAuthToken('jwt-1'));

    test('sends bearer token and parses remaining', () async {
      adapter.on('GET', '/search/youtube', 200, {
        'success': true,
        'data': {'items': [], 'remaining': 7, 'cached': false},
      });

      final result = await api.searchYouTube('rain', 'TH');

      expect(result.remaining, 7);
      expect(adapter.requests.single.headers['Authorization'], 'Bearer jwt-1');
    });

    for (final c in [
      (429, 'user_limit', YouTubeBlockReason.userLimit),
      (503, 'quota_exhausted', YouTubeBlockReason.quotaExhausted),
      (503, 'search_unavailable', YouTubeBlockReason.unavailable),
    ]) {
      test('${c.$2} -> YouTubeSearchBlockedException', () async {
        adapter.on('GET', '/search/youtube', c.$1, {
          'success': false,
          'error': c.$2,
          'retryAt': '2026-11-07T08:00:00Z',
          'remaining': 0,
        });

        await expectLater(
          api.searchYouTube('rain', 'TH'),
          throwsA(
            isA<YouTubeSearchBlockedException>()
                .having((e) => e.reason, 'reason', c.$3)
                .having((e) => e.statusCode, 'status', c.$1)
                .having(
                  (e) => e.retryAt,
                  'retryAt',
                  DateTime.utc(2026, 11, 7, 8),
                )
                .having((e) => e.remaining, 'remaining', 0),
          ),
        );
      });
    }

    test('youtube_error -> YouTubeErrorException', () async {
      adapter.on('GET', '/search/youtube', 502, {
        'success': false,
        'error': 'youtube_error',
      });

      await expectLater(
        api.searchYouTube('rain', 'TH'),
        throwsA(isA<YouTubeErrorException>()),
      );
    });

    test('401 -> UnauthorizedException and notifies listener', () async {
      var notified = 0;
      api.onUnauthorized = () => notified++;
      adapter.on('GET', '/search/youtube', 401, {
        'success': false,
        'error': 'unauthorized',
      });

      await expectLater(
        api.searchYouTube('rain', 'TH'),
        throwsA(isA<UnauthorizedException>()),
      );
      expect(notified, 1);
    });
  });

  test(
    'authed call without token fails fast with UnauthorizedException',
    () async {
      await expectLater(
        api.fetchPlaylists(),
        throwsA(isA<UnauthorizedException>()),
      );
      expect(adapter.requests, isEmpty);
    },
  );

  test('auth_unavailable on login -> AuthUnavailableException', () async {
    adapter.on('POST', '/auth/google', 503, {
      'success': false,
      'error': 'auth_unavailable',
    });

    await expectLater(
      api.loginWithGoogle('id-token'),
      throwsA(
        isA<AuthUnavailableException>().having(
          (e) => e.statusCode,
          'status',
          503,
        ),
      ),
    );
  });

  test('invalid Google token on login stays UnauthorizedException', () async {
    adapter.on('POST', '/auth/google', 401, {
      'success': false,
      'error': 'invalid token',
    });

    await expectLater(
      api.loginWithGoogle('id-token'),
      throwsA(isA<UnauthorizedException>()),
    );
  });

  test('network failure -> plain ApiException without status', () async {
    await expectLater(
      api.search('rain', 'TH'),
      throwsA(
        isA<ApiException>().having((e) => e.statusCode, 'status', isNull),
      ),
    );
  });

  group('playlists', () {
    setUp(() => api.setAuthToken('jwt-1'));

    test('409 on add item -> AlreadyInPlaylistException', () async {
      adapter.on('POST', '/me/playlists/p1/items', 409, {
        'success': false,
        'error': 'video already in playlist',
      });

      await expectLater(
        api.addPlaylistItem(
          'p1',
          videoId: 'v1',
          title: 't',
          channelTitle: 'c',
          thumbnailUrl: 'https://i.ytimg.com/vi/v1/hqdefault.jpg',
        ),
        throwsA(isA<AlreadyInPlaylistException>()),
      );
      expect(adapter.requests.single.data, {
        'videoId': 'v1',
        'title': 't',
        'channelTitle': 'c',
        'thumbnailUrl': 'https://i.ytimg.com/vi/v1/hqdefault.jpg',
      });
    });

    test('reorder sends videoIds and parses playlist with items', () async {
      adapter.on('PUT', '/me/playlists/p1/order', 200, {
        'success': true,
        'data': {
          'id': 'p1',
          'name': 'Drive',
          'itemCount': 0,
          'items': [
            {'videoId': 'b', 'position': 0, 'title': 'B'},
            {'videoId': 'a', 'position': 1, 'title': 'A'},
          ],
        },
      });

      final p = await api.reorderPlaylist('p1', ['b', 'a']);

      expect(adapter.requests.single.data, {
        'videoIds': ['b', 'a'],
      });
      expect(p.items.map((i) => i.videoId), ['b', 'a']);
      expect(p.itemCount, 2);
    });

    test('list parses playlists', () async {
      adapter.on('GET', '/me/playlists', 200, {
        'success': true,
        'data': [
          {'id': 'p1', 'name': 'Drive', 'itemCount': 24},
        ],
      });

      final list = await api.fetchPlaylists();

      expect(list.single.name, 'Drive');
      expect(list.single.itemCount, 24);
    });
  });

  test('searched trends tab fills country on songs', () async {
    adapter.on('GET', '/trends/searched', 200, {
      'success': true,
      'data': [
        {'id': 'v1', 'title': 'x', 'rank': 1, 'rankChange': 2, 'isNew': false},
      ],
    });

    final songs = await api.fetchSearchedTrends('TH');

    expect(songs.single.countryCode, 'TH');
    expect(songs.single.rankChange, 2);
  });
}
