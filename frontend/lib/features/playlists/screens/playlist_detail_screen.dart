import 'package:flutter/material.dart';
import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/playlists/models/playlist_model.dart';
import 'package:streampulse/features/playlists/repositories/playlist_repository.dart';
import 'package:streampulse/features/tracks/screens/track_browser_screen.dart';

class PlaylistDetailScreen extends StatefulWidget {
  final ApiClient apiClient;
  final String playlistId;

  const PlaylistDetailScreen({
    super.key,
    required this.apiClient,
    required this.playlistId,
  });

  @override
  State<PlaylistDetailScreen> createState() => _PlaylistDetailScreenState();
}

class _PlaylistDetailScreenState extends State<PlaylistDetailScreen> {
  late final _repo = PlaylistRepository(apiClient: widget.apiClient);
  PlaylistModel? _playlist;
  String? _error;
  bool _loading = true;
  List<TrackModel> _tracks = [];

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final p = await _repo.get(widget.playlistId);
      setState(() {
        _playlist = p;
        _tracks = List<TrackModel>.from(p.tracks);
        _loading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  Future<void> _reorder(int oldIndex, int newIndex) async {
    setState(() {
      if (newIndex > oldIndex) newIndex--;
      final item = _tracks.removeAt(oldIndex);
      _tracks.insert(newIndex, item);
    });
    try {
      await _repo.reorder(widget.playlistId, _tracks.map((t) => t.id).toList());
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Échec du réordonnancement : $e')),
        );
        _load();
      }
    }
  }

  Future<void> _removeTrack(TrackModel t) async {
    try {
      await _repo.removeTrack(widget.playlistId, t.id);
      setState(() => _tracks.removeWhere((x) => x.id == t.id));
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Échec : $e')),
        );
      }
    }
  }

  Future<void> _openTrackBrowser() async {
    final added = await Navigator.of(context).push<bool>(
      MaterialPageRoute(
        builder: (_) => TrackBrowserScreen(
          apiClient: widget.apiClient,
          playlistId: widget.playlistId,
          alreadyInPlaylist: _tracks.map((t) => t.id).toList(),
        ),
      ),
    );
    if (added == true) {
      _load();
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_loading) {
      return Scaffold(
        appBar: AppBar(title: const Text('Playlist')),
        body: const Center(child: CircularProgressIndicator()),
      );
    }
    if (_error != null) {
      return Scaffold(
        appBar: AppBar(title: const Text('Playlist')),
        body: Center(child: Text(_error!)),
      );
    }
    final p = _playlist!;
    return Scaffold(
      appBar: AppBar(title: Text(p.name)),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: _openTrackBrowser,
        icon: const Icon(Icons.add),
        label: const Text('Ajouter un track'),
      ),
      body: _tracks.isEmpty
          ? Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Icon(Icons.music_note, size: 64),
                  const SizedBox(height: 8),
                  const Text('Cette playlist est vide'),
                  const SizedBox(height: 16),
                  OutlinedButton.icon(
                    onPressed: _openTrackBrowser,
                    icon: const Icon(Icons.add),
                    label: const Text('Ajouter un premier track'),
                  ),
                ],
              ),
            )
          : ReorderableListView.builder(
              padding: const EdgeInsets.only(bottom: 80),
              itemCount: _tracks.length,
              onReorder: _reorder,
              itemBuilder: (_, i) {
                final t = _tracks[i];
                return ListTile(
                  key: ValueKey(t.id),
                  leading: const Icon(Icons.drag_handle),
                  title: Text(t.title),
                  subtitle: Text(t.artist.isEmpty ? '—' : t.artist),
                  trailing: IconButton(
                    icon: const Icon(Icons.remove_circle_outline),
                    onPressed: () => _removeTrack(t),
                  ),
                );
              },
            ),
    );
  }
}
