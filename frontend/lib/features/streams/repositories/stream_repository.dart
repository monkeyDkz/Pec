import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/streams/models/stream_model.dart';

class StreamRepository {
  final ApiClient apiClient;

  StreamRepository({required this.apiClient});

  Future<List<StreamModel>> listLive() async {
    final response = await apiClient.dio.get('/streams');
    final list = (response.data as List).cast<Map<String, dynamic>>();
    return list.map(StreamModel.fromJson).toList(growable: false);
  }

  Future<StreamModel> getById(String id) async {
    final response = await apiClient.dio.get('/streams/$id');
    return StreamModel.fromJson(response.data as Map<String, dynamic>);
  }

  /// Returns the absolute URL to consume the live audio flow.
  /// The mobile player connects directly to this URL using just_audio.
  String listenUrl(String streamId) {
    final base = apiClient.dio.options.baseUrl.replaceAll(RegExp(r'/+$'), '');
    return '$base/streams/$streamId/listen';
  }
}
