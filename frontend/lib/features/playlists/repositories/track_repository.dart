import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/playlists/models/track.dart';

class TrackRepository {
  final ApiClient apiClient;

  TrackRepository({required this.apiClient});

  Future<List<Track>> list() async {
    final res = await apiClient.dio.get('/tracks');
    final data = res.data as List<dynamic>;
    return data.map((e) => Track.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<Track> create({
    required String title,
    required String fileUrl,
    String artist = '',
    int duration = 0,
  }) async {
    final res = await apiClient.dio.post('/tracks', data: {
      'title': title,
      'artist': artist,
      'duration': duration,
      'file_url': fileUrl,
    });
    return Track.fromJson(res.data as Map<String, dynamic>);
  }
}
