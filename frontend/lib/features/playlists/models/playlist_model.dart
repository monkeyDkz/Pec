import 'package:equatable/equatable.dart';

class TrackModel extends Equatable {
  final String id;
  final String title;
  final String artist;
  final int durationSeconds;
  final String fileUrl;

  const TrackModel({
    required this.id,
    required this.title,
    required this.artist,
    required this.durationSeconds,
    required this.fileUrl,
  });

  factory TrackModel.fromJson(Map<String, dynamic> json) => TrackModel(
        id: json['id'] as String,
        title: json['title'] as String? ?? '',
        artist: json['artist'] as String? ?? '',
        durationSeconds: (json['duration_seconds'] as num?)?.toInt() ?? 0,
        fileUrl: json['file_url'] as String? ?? '',
      );

  @override
  List<Object?> get props => [id, title, artist];
}

class PlaylistModel extends Equatable {
  final String id;
  final String name;
  final String description;
  final String ownerId;
  final List<TrackModel> tracks;
  final DateTime createdAt;

  const PlaylistModel({
    required this.id,
    required this.name,
    required this.description,
    required this.ownerId,
    required this.tracks,
    required this.createdAt,
  });

  factory PlaylistModel.fromJson(Map<String, dynamic> json) => PlaylistModel(
        id: json['id'] as String,
        name: json['name'] as String? ?? '',
        description: json['description'] as String? ?? '',
        ownerId: json['owner_id'] as String? ?? '',
        tracks: ((json['tracks'] as List?) ?? const [])
            .cast<Map<String, dynamic>>()
            .map(TrackModel.fromJson)
            .toList(growable: false),
        createdAt: DateTime.tryParse(json['created_at'] as String? ?? '') ??
            DateTime.now(),
      );

  @override
  List<Object?> get props => [id, name, tracks];
}
