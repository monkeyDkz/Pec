import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/features/player/cubit/player_cubit.dart';

class NowPlayingScreen extends StatelessWidget {
  const NowPlayingScreen({super.key});

  String _fmt(Duration d) {
    final m = d.inMinutes.remainder(60).toString().padLeft(2, '0');
    final s = d.inSeconds.remainder(60).toString().padLeft(2, '0');
    return d.inHours > 0 ? '${d.inHours}:$m:$s' : '$m:$s';
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Now Playing')),
      body: BlocBuilder<PlayerCubit, PlayerState>(
        builder: (context, state) {
          final stream = state.current;
          if (stream == null) {
            return const Center(child: Text('Nothing playing'));
          }
          final cubit = context.read<PlayerCubit>();

          return Padding(
            padding: const EdgeInsets.all(32),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Container(
                  width: 200,
                  height: 200,
                  decoration: BoxDecoration(
                    color: theme.colorScheme.primaryContainer,
                    borderRadius: BorderRadius.circular(24),
                  ),
                  child: Icon(Icons.graphic_eq,
                      size: 96, color: theme.colorScheme.onPrimaryContainer),
                ),
                const SizedBox(height: 32),
                Text(stream.title,
                    style: theme.textTheme.headlineSmall, textAlign: TextAlign.center),
                const SizedBox(height: 8),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Container(
                      width: 8,
                      height: 8,
                      decoration: BoxDecoration(
                          color: theme.colorScheme.error, shape: BoxShape.circle),
                    ),
                    const SizedBox(width: 6),
                    Text('LIVE • ${stream.broadcaster}', style: theme.textTheme.bodyMedium),
                  ],
                ),
                const SizedBox(height: 24),
                StreamBuilder<Duration>(
                  stream: cubit.audioPlayer.positionStream,
                  builder: (context, snap) {
                    final pos = snap.data ?? Duration.zero;
                    return Column(
                      children: [
                        LinearProgressIndicator(
                          value: state.isPlaying ? null : 0,
                        ),
                        const SizedBox(height: 4),
                        Text('Elapsed ${_fmt(pos)}', style: theme.textTheme.bodySmall),
                      ],
                    );
                  },
                ),
                const SizedBox(height: 24),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    IconButton.outlined(
                      iconSize: 32,
                      icon: const Icon(Icons.stop),
                      onPressed: () {
                        cubit.stop();
                        Navigator.of(context).maybePop();
                      },
                    ),
                    const SizedBox(width: 24),
                    state.buffering
                        ? const SizedBox(
                            width: 64,
                            height: 64,
                            child: Center(child: CircularProgressIndicator()))
                        : IconButton.filled(
                            iconSize: 48,
                            icon: Icon(state.isPlaying ? Icons.pause : Icons.play_arrow),
                            onPressed: () =>
                                state.isPlaying ? cubit.pause() : cubit.resume(),
                          ),
                  ],
                ),
                const SizedBox(height: 24),
                Row(
                  children: [
                    const Icon(Icons.volume_down),
                    Expanded(
                      child: Slider(
                        value: state.volume,
                        onChanged: cubit.setVolume,
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
    );
  }
}
