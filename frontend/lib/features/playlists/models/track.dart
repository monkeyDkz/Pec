class Track {
  final String id;
  final String title;
  final String artist;
  final int duration;
  final String fileUrl;

  const Track({
    required this.id,
    required this.title,
    required this.artist,
    required this.duration,
    required this.fileUrl,
  });

  factory Track.fromJson(Map<String, dynamic> json) => Track(
        id: json['id'] as String,
        title: json['title'] as String? ?? '',
        artist: json['artist'] as String? ?? '',
        duration: json['duration'] as int? ?? 0,
        fileUrl: json['file_url'] as String? ?? '',
      );
}
