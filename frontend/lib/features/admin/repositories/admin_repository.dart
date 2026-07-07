import 'package:streampulse/core/api/api_client.dart';

class AdminUser {
  final String id;
  final String email;
  final String username;
  final String role;
  AdminUser({
    required this.id,
    required this.email,
    required this.username,
    required this.role,
  });
  factory AdminUser.fromJson(Map<String, dynamic> j) => AdminUser(
        id: j['id'] as String,
        email: j['email'] as String? ?? '',
        username: j['username'] as String? ?? '',
        role: j['role'] as String? ?? 'user',
      );
}

class AdminStats {
  final int totalUsers;
  final int totalBroadcasters;
  final int totalStreams;
  final int liveStreams;
  final int totalPlaylists;
  final int totalTracks;
  final int totalFeedbacks;

  AdminStats({
    required this.totalUsers,
    required this.totalBroadcasters,
    required this.totalStreams,
    required this.liveStreams,
    required this.totalPlaylists,
    required this.totalTracks,
    required this.totalFeedbacks,
  });

  factory AdminStats.fromJson(Map<String, dynamic> j) => AdminStats(
        totalUsers: (j['total_users'] as num?)?.toInt() ?? 0,
        totalBroadcasters: (j['total_broadcasters'] as num?)?.toInt() ?? 0,
        totalStreams: (j['total_streams'] as num?)?.toInt() ?? 0,
        liveStreams: (j['live_streams'] as num?)?.toInt() ?? 0,
        totalPlaylists: (j['total_playlists'] as num?)?.toInt() ?? 0,
        totalTracks: (j['total_tracks'] as num?)?.toInt() ?? 0,
        totalFeedbacks: (j['total_feedbacks'] as num?)?.toInt() ?? 0,
      );
}

class AdminFeedback {
  final String id;
  final String userId;
  final int rating;
  final String comment;
  final DateTime createdAt;
  AdminFeedback({
    required this.id,
    required this.userId,
    required this.rating,
    required this.comment,
    required this.createdAt,
  });
  factory AdminFeedback.fromJson(Map<String, dynamic> j) => AdminFeedback(
        id: j['id'] as String,
        userId: j['user_id'] as String? ?? '',
        rating: (j['rating'] as num?)?.toInt() ?? 0,
        comment: j['comment'] as String? ?? '',
        createdAt: DateTime.tryParse(j['created_at'] as String? ?? '') ??
            DateTime.now(),
      );
}

class AdminRepository {
  final ApiClient apiClient;
  AdminRepository({required this.apiClient});

  Future<List<AdminUser>> listUsers({int offset = 0, int limit = 100}) async {
    final r = await apiClient.dio.get('/admin/users', queryParameters: {
      'offset': offset,
      'limit': limit,
    });
    return (r.data as List)
        .cast<Map<String, dynamic>>()
        .map(AdminUser.fromJson)
        .toList(growable: false);
  }

  Future<AdminUser> updateRole(String userId, String role) async {
    final r = await apiClient.dio.put('/admin/users/$userId/role', data: {'role': role});
    return AdminUser.fromJson(r.data as Map<String, dynamic>);
  }

  Future<AdminStats> stats() async {
    final r = await apiClient.dio.get('/admin/stats');
    return AdminStats.fromJson(r.data as Map<String, dynamic>);
  }

  Future<List<AdminFeedback>> listFeedback({int offset = 0, int limit = 50}) async {
    final r = await apiClient.dio.get('/admin/feedback', queryParameters: {
      'offset': offset,
      'limit': limit,
    });
    return (r.data as List)
        .cast<Map<String, dynamic>>()
        .map(AdminFeedback.fromJson)
        .toList(growable: false);
  }
}
