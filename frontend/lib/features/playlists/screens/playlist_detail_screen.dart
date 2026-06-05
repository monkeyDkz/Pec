import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/features/playlists/models/playlist.dart';
import 'package:streampulse/features/playlists/repositories/playlist_repository.dart';
import 'package:streampulse/features/playlists/repositories/track_repository.dart';

class PlaylistDetailScreen extends StatefulWidget {
  final String id;
  const PlaylistDetailScreen({super.key, required this.id});

  @override
  State<PlaylistDetailScreen> createState() => _PlaylistDetailScreenState();
}

class _PlaylistDetailScreenState extends State<PlaylistDetailScreen> {
  Playlist? _playlist;
  bool _loading = true;
  String? _error;

  late final PlaylistRepository _playlistRepo = context.read<PlaylistRepository>();
  late final TrackRepository _trackRepo = context.read<TrackRepository>();

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() => _loading = true);
    try {
      final p = await _playlistRepo.get(widget.id);
      if (mounted) {
        setState(() {
          _playlist = p;
          _loading = false;
        });
      }
    } catch (_) {
      if (mounted) {
        setState(() {
          _error = 'Failed to load playlist';
          _loading = false;
        });
      }
    }
  }

  Future<void> _addTrackSheet() async {
    final tracks = await _trackRepo.list();
    if (!mounted) return;
    await showModalBottomSheet<void>(
      context: context,
      builder: (sheetCtx) {
        if (tracks.isEmpty) {
          return const SizedBox(
            height: 160,
            child: Center(child: Text('No tracks available')),
          );
        }
        return ListView(
          children: [
            for (final t in tracks)
              ListTile(
                leading: const Icon(Icons.music_note),
                title: Text(t.title),
                subtitle: Text(t.artist),
                onTap: () async {
                  await _playlistRepo.addTrack(widget.id, t.id);
                  if (sheetCtx.mounted) Navigator.of(sheetCtx).pop();
                  await _load();
                },
              ),
          ],
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(_playlist?.name ?? 'Playlist')),
      floatingActionButton: FloatingActionButton(
        onPressed: _addTrackSheet,
        child: const Icon(Icons.add),
      ),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : _error != null
              ? Center(child: Text(_error!))
              : (_playlist!.tracks.isEmpty
                  ? const Center(child: Text('No tracks. Tap + to add one.'))
                  : ListView(
                      children: [
                        for (final t in _playlist!.tracks)
                          ListTile(
                            leading: const Icon(Icons.music_note),
                            title: Text(t.title),
                            subtitle: Text(t.artist),
                            trailing: IconButton(
                              icon: const Icon(Icons.remove_circle_outline),
                              onPressed: () async {
                                await _playlistRepo.removeTrack(widget.id, t.id);
                                await _load();
                              },
                            ),
                          ),
                      ],
                    )),
    );
  }
}
