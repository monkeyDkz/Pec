import 'package:flutter_test/flutter_test.dart';
import 'package:streampulse/features/auth/models/user.dart';
import 'package:streampulse/features/playlists/models/playlist.dart';
import 'package:streampulse/features/streams/models/live_stream.dart';

void main() {
  group('User', () {
    test('parses json and resolves role helpers', () {
      final admin = User.fromJson({'id': '1', 'email': 'a@b.c', 'username': 'a', 'role': 'admin'});
      expect(admin.isAdmin, isTrue);
      expect(admin.isBroadcaster, isTrue);

      final broadcaster = User.fromJson({'id': '2', 'username': 'b', 'role': 'broadcaster'});
      expect(broadcaster.isAdmin, isFalse);
      expect(broadcaster.isBroadcaster, isTrue);

      final user = User.fromJson({'id': '3'});
      expect(user.role, 'user');
      expect(user.isBroadcaster, isFalse);
    });
  });

  group('LiveStream', () {
    test('parses json and computes isLive', () {
      final s = LiveStream.fromJson({
        'id': 's1',
        'title': 'Radio',
        'broadcaster': 'dj',
        'status': 'live',
        'listener_count': 5,
      });
      expect(s.isLive, isTrue);
      expect(s.listenerCount, 5);
      expect(s.broadcaster, 'dj');
    });
  });

  group('Playlist', () {
    test('parses nested tracks', () {
      final p = Playlist.fromJson({
        'id': 'p1',
        'name': 'Chill',
        'tracks': [
          {'id': 't1', 'title': 'Song', 'artist': 'A', 'file_url': 'https://x'},
        ],
      });
      expect(p.name, 'Chill');
      expect(p.tracks, hasLength(1));
      expect(p.tracks.first.title, 'Song');
    });

    test('handles missing tracks gracefully', () {
      final p = Playlist.fromJson({'id': 'p2', 'name': 'Empty'});
      expect(p.tracks, isEmpty);
    });
  });
}
