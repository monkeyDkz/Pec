import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/auth/models/user.dart';

class AdminStats {
  final int totalUsers;
  final int totalStreams;
  final int liveStreams;
  final int totalTracks;
  final int totalPlaylists;

  const AdminStats({
    required this.totalUsers,
    required this.totalStreams,
    required this.liveStreams,
    required this.totalTracks,
    required this.totalPlaylists,
  });

  factory AdminStats.fromJson(Map<String, dynamic> json) => AdminStats(
        totalUsers: json['total_users'] as int? ?? 0,
        totalStreams: json['total_streams'] as int? ?? 0,
        liveStreams: json['live_streams'] as int? ?? 0,
        totalTracks: json['total_tracks'] as int? ?? 0,
        totalPlaylists: json['total_playlists'] as int? ?? 0,
      );
}

class AdminRepository {
  final ApiClient apiClient;

  AdminRepository({required this.apiClient});

  Future<AdminStats> stats() async {
    final res = await apiClient.dio.get('/admin/stats');
    return AdminStats.fromJson(res.data as Map<String, dynamic>);
  }

  Future<List<User>> listUsers() async {
    final res = await apiClient.dio.get('/admin/users');
    final data = res.data as List<dynamic>;
    return data.map((e) => User.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<User> updateRole(String userId, String role) async {
    final res = await apiClient.dio.put('/admin/users/$userId/role', data: {'role': role});
    return User.fromJson(res.data as Map<String, dynamic>);
  }
}
