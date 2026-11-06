class AppUser {
  final String id;
  final String email;
  final String name;
  final String avatarUrl;

  const AppUser({
    required this.id,
    required this.email,
    required this.name,
    required this.avatarUrl,
  });

  factory AppUser.fromJson(Map<String, dynamic> json) {
    return AppUser(
      id: json['id'] as String? ?? '',
      email: json['email'] as String? ?? '',
      name: json['name'] as String? ?? '',
      avatarUrl: json['avatarUrl'] as String? ?? '',
    );
  }
}

class AuthSession {
  final String token;
  final AppUser? user;

  const AuthSession({required this.token, this.user});

  factory AuthSession.fromJson(Map<String, dynamic> json) {
    final user = json['user'];
    return AuthSession(
      token: json['token'] as String? ?? '',
      user: user is Map<String, dynamic> ? AppUser.fromJson(user) : null,
    );
  }
}
