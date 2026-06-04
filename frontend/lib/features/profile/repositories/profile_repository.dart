import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/auth/models/user.dart';

class ProfileRepository {
  final ApiClient apiClient;

  ProfileRepository({required this.apiClient});

  Future<User> me() async {
    final res = await apiClient.dio.get('/users/me');
    return User.fromJson(res.data as Map<String, dynamic>);
  }

  Future<User> updateUsername(String username) async {
    final res = await apiClient.dio.put('/users/me', data: {'username': username});
    return User.fromJson(res.data as Map<String, dynamic>);
  }
}
