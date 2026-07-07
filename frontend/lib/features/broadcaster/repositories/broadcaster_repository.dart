import 'dart:async';
import 'dart:typed_data';

import 'package:dio/dio.dart';
import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/core/storage/secure_storage.dart';
import 'package:streampulse/features/streams/models/stream_model.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class BroadcasterRepository {
  final ApiClient apiClient;
  final SecureStorage _storage = SecureStorage();

  BroadcasterRepository({required this.apiClient});

  Future<StreamModel> createStream({
    required String title,
    required String description,
  }) async {
    final response = await apiClient.dio.post(
      '/streams',
      data: {'title': title, 'description': description},
    );
    return StreamModel.fromJson(response.data as Map<String, dynamic>);
  }

  Future<void> deleteStream(String id) async {
    await apiClient.dio.delete('/streams/$id');
  }

  /// Publishes a continuous audio stream of bytes.
  /// Used for live broadcasting: the API consumes chunks as fast as they
  /// arrive and fans them out to listeners through the in-memory hub.
  Future<void> publish({
    required String streamId,
    required Stream<List<int>> audio,
    required CancelToken cancelToken,
  }) async {
    await apiClient.dio.post<void>(
      '/streams/$streamId/publish',
      data: audio,
      options: Options(
        contentType: 'application/octet-stream',
        responseType: ResponseType.bytes,
        sendTimeout: null,
        receiveTimeout: null,
      ),
      cancelToken: cancelToken,
    );
  }

  /// Publishes a continuous audio stream over a WebSocket. Required on the web,
  /// where browsers cannot stream an HTTP request body. Auth is passed as a
  /// `token` query parameter (WebSockets can't set an Authorization header).
  ///
  /// Completes when the audio source ends, the socket closes, or [cancelToken]
  /// is cancelled (broadcaster pressed "stop").
  Future<void> publishWs({
    required String streamId,
    required Stream<List<int>> audio,
    required CancelToken cancelToken,
  }) async {
    final token = await _storage.getToken();
    // http://host/api/v1 -> ws://host/api/v1 ; https -> wss
    final wsBase = apiClient.dio.options.baseUrl.replaceFirst('http', 'ws');
    final uri = Uri.parse(
      '$wsBase/streams/$streamId/publish/ws?token=${Uri.encodeQueryComponent(token ?? '')}',
    );

    final channel = WebSocketChannel.connect(uri);
    await channel.ready;

    final done = Completer<void>();
    void finish() {
      if (!done.isCompleted) done.complete();
    }

    // Stop when the broadcaster cancels.
    cancelToken.whenCancel.then((_) async {
      await channel.sink.close();
      finish();
    });
    // Stop if the server closes the socket.
    channel.stream.listen((_) {}, onDone: finish, onError: (_) => finish());

    final sub = audio.listen(
      (chunk) => channel.sink.add(Uint8List.fromList(chunk)),
      onError: (_) => finish(),
      onDone: () async {
        await channel.sink.close();
        finish();
      },
      cancelOnError: true,
    );

    await done.future;
    await sub.cancel();
  }

  Future<List<StreamModel>> listMine() async {
    final response = await apiClient.dio.get('/streams');
    final list = (response.data as List).cast<Map<String, dynamic>>();
    return list.map(StreamModel.fromJson).toList(growable: false);
  }
}
