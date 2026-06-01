import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/core/storage/secure_storage.dart';

class AuthRepository {
  final ApiClient apiClient;
  final SecureStorage _storage = SecureStorage();

  AuthRepository({required this.apiClient});

  Future<void> login({required String email, required String password}) async {
    final response = await apiClient.dio.post('/auth/login', data: {
      'email': email,
      'password': password,
    });
    await _storage.saveToken(response.data['token']);
  }

  Future<void> register({
    required String email,
    required String username,
    required String password,
  }) async {
    final response = await apiClient.dio.post('/auth/register', data: {
      'email': email,
      'username': username,
      'password': password,
    });
    await _storage.saveToken(response.data['token']);
  }

  Future<void> logout() async {
    await _storage.deleteToken();
  }

  Future<bool> isAuthenticated() async {
    final token = await _storage.getToken();
    return token != null;
  }
}
