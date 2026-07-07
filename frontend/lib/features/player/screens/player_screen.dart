import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:streampulse/features/player/bloc/player_bloc.dart';
import 'package:streampulse/features/streams/repositories/stream_repository.dart';

class PlayerScreen extends StatefulWidget {
  final String streamId;
  final StreamRepository repository;

  const PlayerScreen({super.key, required this.streamId, required this.repository});

  @override
  State<PlayerScreen> createState() => _PlayerScreenState();
}

class _PlayerScreenState extends State<PlayerScreen> {
  late final PlayerBloc _bloc;
  String _title = '';

  @override
  void initState() {
    super.initState();
    _bloc = PlayerBloc();
    _load();
  }

  Future<void> _load() async {
    try {
      final stream = await widget.repository.getById(widget.streamId);
      if (!mounted) return;
      setState(() => _title = stream.title);
      _bloc.add(PlayerLoadRequested(
        url: widget.repository.listenUrl(widget.streamId),
        title: stream.title,
      ));
    } catch (_) {
      // Even without metadata, try to play.
      _bloc.add(PlayerLoadRequested(
        url: widget.repository.listenUrl(widget.streamId),
        title: 'Stream',
      ));
    }
  }

  @override
  void dispose() {
    _bloc.close();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return BlocProvider.value(
      value: _bloc,
      child: Scaffold(
        appBar: AppBar(
          title: Text(_title.isEmpty ? 'Lecture' : _title),
          leading: IconButton(
            icon: const Icon(Icons.arrow_back),
            onPressed: () {
              _bloc.add(const PlayerStopRequested());
              context.pop();
            },
          ),
        ),
        body: BlocBuilder<PlayerBloc, PlayerState>(
          builder: (context, state) {
            if (state is PlayerError) {
              return _Error(message: state.message, onRetry: _load);
            }
            final isPlaying = state is PlayerPlaying;
            final volume = state is PlayerPlaying
                ? state.volume
                : state is PlayerPaused
                    ? state.volume
                    : 1.0;
            return Padding(
              padding: const EdgeInsets.all(24),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.podcasts, size: 96, color: Colors.red),
                  const SizedBox(height: 16),
                  Text(
                    state is PlayerLoading
                        ? 'Chargement…'
                        : isPlaying
                            ? '🔴 En direct'
                            : 'En pause',
                    style: Theme.of(context).textTheme.titleLarge,
                  ),
                  const SizedBox(height: 32),
                  IconButton.filled(
                    iconSize: 64,
                    onPressed: state is PlayerLoading
                        ? null
                        : () => isPlaying
                            ? _bloc.add(const PlayerPauseRequested())
                            : _bloc.add(const PlayerPlayRequested()),
                    tooltip: isPlaying ? 'Pause' : 'Lecture',
                    icon: Icon(isPlaying ? Icons.pause : Icons.play_arrow),
                  ),
                  const SizedBox(height: 32),
                  Row(
                    children: [
                      const Icon(Icons.volume_down),
                      Expanded(
                        child: Slider(
                          value: volume,
                          onChanged: (v) => _bloc.add(PlayerVolumeChanged(v)),
                        ),
                      ),
                      const Icon(Icons.volume_up),
                    ],
                  ),
                ],
              ),
            );
          },
        ),
      ),
    );
  }
}

class _Error extends StatelessWidget {
  final String message;
  final VoidCallback onRetry;
  const _Error({required this.message, required this.onRetry});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline, size: 64),
            const SizedBox(height: 8),
            Text(message, textAlign: TextAlign.center),
            const SizedBox(height: 16),
            FilledButton(onPressed: onRetry, child: const Text('Réessayer')),
          ],
        ),
      ),
    );
  }
}
