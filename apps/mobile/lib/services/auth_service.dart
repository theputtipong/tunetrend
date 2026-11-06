import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:google_sign_in/google_sign_in.dart';

import '../constants/auth_config.dart';
import '../models/user.dart';
import 'api_client.dart';

abstract class TokenStore {
  Future<String?> read();
  Future<void> write(String token);
  Future<void> delete();
}

class SecureTokenStore implements TokenStore {
  static const _key = 'tunetrend_session_jwt';
  final FlutterSecureStorage _storage = const FlutterSecureStorage();

  @override
  Future<String?> read() async {
    try {
      return await _storage.read(key: _key);
    } catch (_) {
      return null;
    }
  }

  @override
  Future<void> write(String token) => _storage.write(key: _key, value: token);

  @override
  Future<void> delete() => _storage.delete(key: _key);
}

abstract class GoogleIdTokenProvider {
  Future<String?> obtainIdToken();
  Future<String?> obtainIdTokenSilently();
  Future<void> signOut();
}

class GoogleSignInIdTokenProvider implements GoogleIdTokenProvider {
  Future<void>? _init;

  Future<void> _ensureInitialized() {
    return _init ??= GoogleSignIn.instance.initialize(
      clientId:
          defaultTargetPlatform == TargetPlatform.iOS &&
              kGoogleIosClientId.isNotEmpty
          ? kGoogleIosClientId
          : null,
      serverClientId: kGoogleServerClientId,
    );
  }

  @override
  Future<String?> obtainIdToken() async {
    await _ensureInitialized();
    if (!GoogleSignIn.instance.supportsAuthenticate()) return null;
    try {
      final account = await GoogleSignIn.instance.authenticate();
      return account.authentication.idToken;
    } on GoogleSignInException catch (e) {
      if (e.code == GoogleSignInExceptionCode.canceled ||
          e.code == GoogleSignInExceptionCode.interrupted) {
        return null;
      }
      rethrow;
    }
  }

  @override
  Future<String?> obtainIdTokenSilently() async {
    try {
      await _ensureInitialized();
      final account = await GoogleSignIn.instance
          .attemptLightweightAuthentication();
      return account?.authentication.idToken;
    } catch (_) {
      return null;
    }
  }

  @override
  Future<void> signOut() async {
    await _ensureInitialized();
    await GoogleSignIn.instance.signOut();
  }
}

class AuthService extends ChangeNotifier {
  AuthService({
    required this._api,
    required this._store,
    required this._google,
    required this.enabled,
  }) {
    _api.onUnauthorized = () => unawaited(_clear());
  }

  static final AuthService instance = AuthService(
    api: ApiClient(),
    store: SecureTokenStore(),
    google: GoogleSignInIdTokenProvider(),
    enabled: kLoginEnabled,
  );

  final ApiClient _api;
  final TokenStore _store;
  final GoogleIdTokenProvider _google;
  final bool enabled;

  AppUser? _user;
  bool _busy = false;

  AppUser? get user => _user;
  bool get isSignedIn => enabled && _api.hasAuthToken;
  bool get busy => _busy;

  Future<void> restore() async {
    if (!enabled) return;
    final token = await _store.read();
    if (token == null || token.isEmpty) return;
    _api.setAuthToken(token);
    notifyListeners();
    try {
      _user = await _api.fetchMe();
      notifyListeners();
    } on UnauthorizedException {
      await _clear();
      final idToken = await _google.obtainIdTokenSilently();
      if (idToken == null) return;
      try {
        await _exchange(idToken);
      } on ApiException {
        return;
      }
    } on ApiException {
      return;
    }
  }

  Future<bool> signIn() async {
    if (!enabled) return false;
    _busy = true;
    notifyListeners();
    try {
      final idToken = await _google.obtainIdToken();
      if (idToken == null) return false;
      await _exchange(idToken);
      return true;
    } finally {
      _busy = false;
      notifyListeners();
    }
  }

  Future<void> signOut() async {
    await _clear();
    try {
      await _google.signOut();
    } catch (_) {
      return;
    }
  }

  Future<void> deleteAccount() async {
    await _api.deleteMe();
    await signOut();
  }

  Future<void> _exchange(String idToken) async {
    final session = await _api.loginWithGoogle(idToken);
    if (session.token.isEmpty) throw UnauthorizedException();
    await _store.write(session.token);
    _api.setAuthToken(session.token);
    _user = session.user;
    notifyListeners();
  }

  Future<void> _clear() async {
    final hadSession = _api.hasAuthToken || _user != null;
    _api.setAuthToken(null);
    _user = null;
    await _store.delete();
    if (hadSession) notifyListeners();
  }
}
