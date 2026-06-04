class User {
  final String id;
  final String email;
  final String username;
  final String role;

  const User({
    required this.id,
    required this.email,
    required this.username,
    required this.role,
  });

  factory User.fromJson(Map<String, dynamic> json) => User(
        id: json['id'] as String,
        email: json['email'] as String? ?? '',
        username: json['username'] as String? ?? '',
        role: json['role'] as String? ?? 'user',
      );

  bool get isBroadcaster => role == 'broadcaster' || role == 'admin';
  bool get isAdmin => role == 'admin';
}
