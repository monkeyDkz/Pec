import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/playlists/models/playlist_model.dart';

class PlaylistRepository {
  final ApiClient apiClient;
  PlaylistRepository({required this.apiClient});

  Future<List<PlaylistModel>> listMine() async {
    final response = await apiClient.dio.get('/playlists');
    final list = (response.data as List).cast<Map<String, dynamic>>();
    return list.map(PlaylistModel.fromJson).toList(growable: false);
  }

  Future<PlaylistModel> get(String id) async {
    final response = await apiClient.dio.get('/playlists/$id');
    return PlaylistModel.fromJson(response.data as Map<String, dynamic>);
  }

  Future<PlaylistModel> create({
    required String name,
    required String description,
  }) async {
    final response = await apiClient.dio.post('/playlists', data: {
      'name': name,
      'description': description,
    });
    return PlaylistModel.fromJson(response.data as Map<String, dynamic>);
  }

  Future<void> delete(String id) async {
    await apiClient.dio.delete('/playlists/$id');
  }

  Future<void> addTrack(String playlistId, String trackId) async {
    await apiClient.dio.post('/playlists/$playlistId/tracks/$trackId');
  }

  Future<void> removeTrack(String playlistId, String trackId) async {
    await apiClient.dio.delete('/playlists/$playlistId/tracks/$trackId');
  }

  Future<void> reorder(String playlistId, List<String> trackIds) async {
    await apiClient.dio.put('/playlists/$playlistId/tracks/reorder', data: {
      'track_ids': trackIds,
    });
  }

  Future<List<TrackModel>> listTracks() async {
    final response = await apiClient.dio.get('/tracks');
    final list = (response.data as List).cast<Map<String, dynamic>>();
    return list.map(TrackModel.fromJson).toList(growable: false);
  }
}
