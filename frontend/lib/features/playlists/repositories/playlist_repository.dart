import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/playlists/models/playlist.dart';

class PlaylistRepository {
  final ApiClient apiClient;

  PlaylistRepository({required this.apiClient});

  Future<List<Playlist>> list() async {
    final res = await apiClient.dio.get('/playlists');
    final data = res.data as List<dynamic>;
    return data.map((e) => Playlist.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<Playlist> get(String id) async {
    final res = await apiClient.dio.get('/playlists/$id');
    return Playlist.fromJson(res.data as Map<String, dynamic>);
  }

  Future<Playlist> create({required String name, String description = ''}) async {
    final res = await apiClient.dio.post('/playlists', data: {
      'name': name,
      'description': description,
    });
    return Playlist.fromJson(res.data as Map<String, dynamic>);
  }

  Future<void> delete(String id) async {
    await apiClient.dio.delete('/playlists/$id');
  }

  Future<void> addTrack(String playlistId, String trackId) async {
    await apiClient.dio.post('/playlists/$playlistId/tracks', data: {'track_id': trackId});
  }

  Future<void> removeTrack(String playlistId, String trackId) async {
    await apiClient.dio.delete('/playlists/$playlistId/tracks/$trackId');
  }
}
