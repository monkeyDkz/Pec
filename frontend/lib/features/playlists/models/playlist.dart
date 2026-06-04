import 'package:streampulse/features/playlists/models/track.dart';

class Playlist {
  final String id;
  final String name;
  final String description;
  final List<Track> tracks;

  const Playlist({
    required this.id,
    required this.name,
    required this.description,
    required this.tracks,
  });

  factory Playlist.fromJson(Map<String, dynamic> json) => Playlist(
        id: json['id'] as String,
        name: json['name'] as String? ?? '',
        description: json['description'] as String? ?? '',
        tracks: (json['tracks'] as List<dynamic>? ?? [])
            .map((t) => Track.fromJson(t as Map<String, dynamic>))
            .toList(),
      );
}
