import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:streampulse/features/player/cubit/player_cubit.dart';

class MiniPlayer extends StatelessWidget {
  const MiniPlayer({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<PlayerCubit, PlayerState>(
      builder: (context, state) {
        final stream = state.current;
        if (stream == null) return const SizedBox.shrink();
        final cubit = context.read<PlayerCubit>();
        final theme = Theme.of(context);

        return Material(
          color: theme.colorScheme.surfaceContainerHighest,
          child: InkWell(
            onTap: () => context.push('/player'),
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              child: Row(
                children: [
                  Icon(Icons.graphic_eq, color: theme.colorScheme.primary),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(stream.title,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: theme.textTheme.titleSmall),
                        Text('LIVE • ${stream.broadcaster}',
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                            style: theme.textTheme.bodySmall),
                      ],
                    ),
                  ),
                  if (state.buffering)
                    const SizedBox(
                        width: 24, height: 24, child: CircularProgressIndicator(strokeWidth: 2))
                  else
                    IconButton(
                      icon: Icon(state.isPlaying ? Icons.pause : Icons.play_arrow),
                      onPressed: () => state.isPlaying ? cubit.pause() : cubit.resume(),
                    ),
                  IconButton(
                    icon: const Icon(Icons.stop),
                    onPressed: cubit.stop,
                  ),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}
