import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/streams/models/live_stream.dart';

class StreamRepository {
  final ApiClient apiClient;

  StreamRepository({required this.apiClient});

  Future<List<LiveStream>> listLive() async {
    final res = await apiClient.dio.get('/streams');
    final data = res.data as List<dynamic>;
    return data.map((e) => LiveStream.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<LiveStream> get(String id) async {
    final res = await apiClient.dio.get('/streams/$id');
    return LiveStream.fromJson(res.data as Map<String, dynamic>);
  }

  Future<LiveStream> create({required String title, String description = ''}) async {
    final res = await apiClient.dio.post('/streams', data: {
      'title': title,
      'description': description,
    });
    return LiveStream.fromJson(res.data as Map<String, dynamic>);
  }

  Future<void> stop(String id) async {
    await apiClient.dio.post('/streams/$id/stop');
  }

  Future<void> delete(String id) async {
    await apiClient.dio.delete('/streams/$id');
  }

  /// Absolute URL of the live audio endpoint, consumed by the audio player.
  String listenUrl(String id) => '${apiClient.dio.options.baseUrl}/streams/$id/listen';
}
