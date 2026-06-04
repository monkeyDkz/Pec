class LiveStream {
  final String id;
  final String title;
  final String description;
  final String broadcasterId;
  final String broadcaster;
  final String status;
  final int listenerCount;

  const LiveStream({
    required this.id,
    required this.title,
    required this.description,
    required this.broadcasterId,
    required this.broadcaster,
    required this.status,
    required this.listenerCount,
  });

  bool get isLive => status == 'live';

  factory LiveStream.fromJson(Map<String, dynamic> json) => LiveStream(
        id: json['id'] as String,
        title: json['title'] as String? ?? '',
        description: json['description'] as String? ?? '',
        broadcasterId: json['broadcaster_id'] as String? ?? '',
        broadcaster: json['broadcaster'] as String? ?? '',
        status: json['status'] as String? ?? 'offline',
        listenerCount: json['listener_count'] as int? ?? 0,
      );
}
