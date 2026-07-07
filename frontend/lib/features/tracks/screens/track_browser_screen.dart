import 'package:flutter/material.dart';
import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/playlists/models/playlist_model.dart';
import 'package:streampulse/features/playlists/repositories/playlist_repository.dart';
import 'package:streampulse/features/tracks/repositories/track_repository.dart';

class TrackBrowserScreen extends StatefulWidget {
  final ApiClient apiClient;
  final String playlistId;
  final List<String> alreadyInPlaylist;

  const TrackBrowserScreen({
    super.key,
    required this.apiClient,
    required this.playlistId,
    required this.alreadyInPlaylist,
  });

  @override
  State<TrackBrowserScreen> createState() => _TrackBrowserScreenState();
}

class _TrackBrowserScreenState extends State<TrackBrowserScreen> {
  late final _trackRepo = TrackRepository(apiClient: widget.apiClient);
  late final _playlistRepo = PlaylistRepository(apiClient: widget.apiClient);
  late Future<List<TrackModel>> _future;
  final Set<String> _adding = {};

  @override
  void initState() {
    super.initState();
    _future = _trackRepo.list();
  }

  Future<void> _addToPlaylist(TrackModel t) async {
    setState(() => _adding.add(t.id));
    try {
      await _playlistRepo.addTrack(widget.playlistId, t.id);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('« ${t.title} » ajouté')),
        );
        Navigator.of(context).pop(true);
      }
    } catch (e) {
      if (mounted) {
        setState(() => _adding.remove(t.id));
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Échec : $e')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Ajouter un track')),
      body: FutureBuilder<List<TrackModel>>(
        future: _future,
        builder: (context, snap) {
          if (snap.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snap.hasError) {
            return Center(child: Text('${snap.error}'));
          }
          final tracks = snap.data!
              .where((t) => !widget.alreadyInPlaylist.contains(t.id))
              .toList();
          if (tracks.isEmpty) {
            return const Center(child: Text('Aucun track disponible'));
          }
          return ListView.separated(
            itemCount: tracks.length,
            separatorBuilder: (_, __) => const Divider(height: 1),
            itemBuilder: (_, i) {
              final t = tracks[i];
              final loading = _adding.contains(t.id);
              return ListTile(
                leading: const CircleAvatar(child: Icon(Icons.music_note)),
                title: Text(t.title),
                subtitle: Text(t.artist.isEmpty ? '—' : t.artist),
                trailing: loading
                    ? const SizedBox(
                        width: 24,
                        height: 24,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : IconButton(
                        icon: const Icon(Icons.add_circle_outline),
                        onPressed: () => _addToPlaylist(t),
                      ),
              );
            },
          );
        },
      ),
    );
  }
}
