import 'package:equatable/equatable.dart';

class StreamModel extends Equatable {
  final String id;
  final String title;
  final String description;
  final String broadcasterId;
  final String broadcasterUsername;
  final String status; // 'live' | 'offline'
  final int listenerCount;
  final DateTime createdAt;
  final DateTime updatedAt;

  const StreamModel({
    required this.id,
    required this.title,
    required this.description,
    required this.broadcasterId,
    required this.broadcasterUsername,
    required this.status,
    required this.listenerCount,
    required this.createdAt,
    required this.updatedAt,
  });

  bool get isLive => status == 'live';

  factory StreamModel.fromJson(Map<String, dynamic> json) => StreamModel(
        id: json['id'] as String,
        title: json['title'] as String? ?? '',
        description: json['description'] as String? ?? '',
        broadcasterId: json['broadcaster_id'] as String? ?? '',
        broadcasterUsername: json['broadcaster_username'] as String? ?? '',
        status: json['status'] as String? ?? 'offline',
        listenerCount: (json['listener_count'] as num?)?.toInt() ?? 0,
        createdAt: DateTime.tryParse(json['created_at'] as String? ?? '') ??
            DateTime.now(),
        updatedAt: DateTime.tryParse(json['updated_at'] as String? ?? '') ??
            DateTime.now(),
      );

  @override
  List<Object?> get props => [id, title, status, listenerCount, updatedAt];
}
