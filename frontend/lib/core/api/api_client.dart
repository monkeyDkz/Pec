import 'package:dio/dio.dart';
import 'package:streampulse/core/storage/secure_storage.dart';

class ApiClient {
  late final Dio dio;
  final SecureStorage _storage = SecureStorage();

  ApiClient() {
    dio = Dio(BaseOptions(
      // Configure via environment or build config
      baseUrl: const String.fromEnvironment('API_URL', defaultValue: 'http://localhost:8080/api/v1'),
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
      headers: {'Content-Type': 'application/json'},
    ));

    dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        final token = await _storage.getToken();
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        handler.next(options);
      },
      onError: (error, handler) {
        if (error.response?.statusCode == 401) {
          _storage.deleteToken();
        }
        handler.next(error);
      },
    ));
  }
}
