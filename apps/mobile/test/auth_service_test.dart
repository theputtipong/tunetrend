import 'package:flutter_test/flutter_test.dart';
import 'package:tunetrend_mobile/services/api_client.dart';
import 'package:tunetrend_mobile/services/auth_service.dart';

import 'helpers/fake_http.dart';

class FakeTokenStore implements TokenStore {
  String? value;
  FakeTokenStore([this.value]);

  @override
  Future<String?> read() async => value;

  @override
  Future<void> write(String token) async => value = token;

  @override
  Future<void> delete() async => value = null;
}

class FakeGoogle implements GoogleIdTokenProvider {
  String? interactiveToken;
  String? silentToken;
  int signOuts = 0;

  @override
  Future<String?> obtainIdToken() async => interactiveToken;

  @override
  Future<String?> obtainIdTokenSilently() async => silentToken;

  @override
  Future<void> signOut() async => signOuts++;
}

const _user = {'id': 'u1', 'email': 'a@b.c', 'name': 'Alice', 'avatarUrl': ''};

void main() {
  late FakeAdapter adapter;
  late ApiClient api;
  late FakeTokenStore store;
  late FakeGoogle google;
  late AuthService auth;

  setUp(() {
    adapter = FakeAdapter();
    api = ApiClient.withDio(fakeDio(adapter));
    store = FakeTokenStore();
    google = FakeGoogle();
    auth = AuthService(api: api, store: store, google: google, enabled: true);
  });

  test('sign in exchanges Google ID token and persists session JWT', () async {
    google.interactiveToken = 'google-id-token';
    adapter.on('POST', '/auth/google', 200, {
      'success': true,
      'data': {'token': 'jwt-1', 'user': _user},
    });

    final ok = await auth.signIn();

    expect(ok, isTrue);
    expect(adapter.requests.single.data, {'idToken': 'google-id-token'});
    expect(store.value, 'jwt-1');
    expect(api.hasAuthToken, isTrue);
    expect(auth.isSignedIn, isTrue);
    expect(auth.user?.name, 'Alice');
  });

  test(
    'auth_unavailable during sign in surfaces typed error and keeps state',
    () async {
      google.interactiveToken = 'google-id-token';
      adapter.on('POST', '/auth/google', 503, {
        'success': false,
        'error': 'auth_unavailable',
      });

      await expectLater(
        auth.signIn(),
        throwsA(isA<AuthUnavailableException>()),
      );
      expect(store.value, isNull);
      expect(auth.isSignedIn, isFalse);
      expect(auth.busy, isFalse);
    },
  );

  test('cancelled Google sign in leaves state untouched', () async {
    google.interactiveToken = null;

    expect(await auth.signIn(), isFalse);
    expect(adapter.requests, isEmpty);
    expect(store.value, isNull);
    expect(auth.isSignedIn, isFalse);
  });

  test('restore with stored token loads user via GET /me', () async {
    store.value = 'jwt-1';
    adapter.on('GET', '/me', 200, {'success': true, 'data': _user});

    await auth.restore();

    expect(auth.isSignedIn, isTrue);
    expect(auth.user?.email, 'a@b.c');
    expect(adapter.requests.single.headers['Authorization'], 'Bearer jwt-1');
  });

  test(
    'restore with expired token clears storage when silent sign in fails',
    () async {
      store.value = 'jwt-old';
      adapter.on('GET', '/me', 401, {
        'success': false,
        'error': 'unauthorized',
      });

      await auth.restore();

      expect(store.value, isNull);
      expect(api.hasAuthToken, isFalse);
      expect(auth.isSignedIn, isFalse);
    },
  );

  test(
    'restore with expired token re-logs in silently when Google allows',
    () async {
      store.value = 'jwt-old';
      google.silentToken = 'silent-id-token';
      adapter.on('GET', '/me', 401, {
        'success': false,
        'error': 'unauthorized',
      });
      adapter.on('POST', '/auth/google', 200, {
        'success': true,
        'data': {'token': 'jwt-new', 'user': _user},
      });

      await auth.restore();

      expect(store.value, 'jwt-new');
      expect(auth.isSignedIn, isTrue);
    },
  );

  test('restore keeps token on network failure', () async {
    store.value = 'jwt-1';

    await auth.restore();

    expect(store.value, 'jwt-1');
    expect(auth.isSignedIn, isTrue);
    expect(auth.user, isNull);
  });

  test('restore without token does nothing', () async {
    await auth.restore();

    expect(adapter.requests, isEmpty);
    expect(auth.isSignedIn, isFalse);
  });

  test('401 from any authed call signs the user out', () async {
    store.value = 'jwt-1';
    adapter.on('GET', '/me', 200, {'success': true, 'data': _user});
    await auth.restore();
    adapter.on('GET', '/me/playlists', 401, {
      'success': false,
      'error': 'unauthorized',
    });

    await expectLater(
      api.fetchPlaylists(),
      throwsA(isA<UnauthorizedException>()),
    );
    await Future<void>.delayed(Duration.zero);

    expect(auth.isSignedIn, isFalse);
    expect(store.value, isNull);
  });

  test('sign out clears token and Google session', () async {
    store.value = 'jwt-1';
    adapter.on('GET', '/me', 200, {'success': true, 'data': _user});
    await auth.restore();

    await auth.signOut();

    expect(store.value, isNull);
    expect(api.hasAuthToken, isFalse);
    expect(google.signOuts, 1);
    expect(auth.isSignedIn, isFalse);
  });

  test('delete account calls DELETE /me then clears session', () async {
    store.value = 'jwt-1';
    adapter.on('GET', '/me', 200, {'success': true, 'data': _user});
    adapter.on('DELETE', '/me', 200, {
      'success': true,
      'data': {'deleted': true},
    });
    await auth.restore();

    await auth.deleteAccount();

    expect(adapter.requests.last.method, 'DELETE');
    expect(store.value, isNull);
    expect(auth.isSignedIn, isFalse);
  });

  test('disabled auth never touches storage', () async {
    final disabled = AuthService(
      api: api,
      store: FakeTokenStore('jwt-1'),
      google: google,
      enabled: false,
    );

    await disabled.restore();

    expect(disabled.isSignedIn, isFalse);
    expect(adapter.requests, isEmpty);
  });
}
