import 'dart:io';

import 'package:dio/dio.dart';
import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/playlists/models/playlist_model.dart';

class TrackRepository {
  final ApiClient apiClient;
  TrackRepository({required this.apiClient});

  Future<List<TrackModel>> list({int offset = 0, int limit = 100}) async {
    final r = await apiClient.dio.get('/tracks', queryParameters: {
      'offset': offset,
      'limit': limit,
    });
    return (r.data as List)
        .cast<Map<String, dynamic>>()
        .map(TrackModel.fromJson)
        .toList(growable: false);
  }

  Future<TrackModel> upload({
    required File file,
    required String title,
    String? artist,
    void Function(int sent, int total)? onProgress,
  }) async {
    final form = FormData.fromMap({
      'title': title,
      if (artist != null && artist.isNotEmpty) 'artist': artist,
      'file': await MultipartFile.fromFile(file.path, filename: file.uri.pathSegments.last),
    });
    final r = await apiClient.dio.post(
      '/tracks',
      data: form,
      options: Options(contentType: 'multipart/form-data'),
      onSendProgress: onProgress,
    );
    return TrackModel.fromJson(r.data as Map<String, dynamic>);
  }

  Future<void> delete(String id) async {
    await apiClient.dio.delete('/tracks/$id');
  }
}
