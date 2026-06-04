import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:streampulse/features/player/cubit/player_cubit.dart';
import 'package:streampulse/features/streams/cubit/streams_cubit.dart';
import 'package:streampulse/features/streams/models/live_stream.dart';
import 'package:streampulse/features/streams/repositories/stream_repository.dart';

class StreamsScreen extends StatelessWidget {
  const StreamsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (ctx) => StreamsCubit(ctx.read<StreamRepository>())..load(),
      child: const _StreamsView(),
    );
  }
}

class _StreamsView extends StatelessWidget {
  const _StreamsView();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Live Streams'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => context.read<StreamsCubit>().load(),
          ),
        ],
      ),
      body: BlocBuilder<StreamsCubit, StreamsState>(
        builder: (context, state) {
          if (state.loading && state.streams.isEmpty) {
            return const Center(child: CircularProgressIndicator());
          }
          if (state.error != null) {
            return _ErrorRetry(
              message: state.error!,
              onRetry: () => context.read<StreamsCubit>().load(),
            );
          }
          if (state.streams.isEmpty) {
            return const _Empty(message: 'No live streams right now');
          }
          return RefreshIndicator(
            onRefresh: () => context.read<StreamsCubit>().load(),
            child: ListView.builder(
              padding: const EdgeInsets.all(8),
              itemCount: state.streams.length,
              itemBuilder: (context, i) => _StreamCard(stream: state.streams[i]),
            ),
          );
        },
      ),
    );
  }
}

class _StreamCard extends StatelessWidget {
  final LiveStream stream;
  const _StreamCard({required this.stream});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      child: ListTile(
        leading: const CircleAvatar(child: Icon(Icons.radio)),
        title: Text(stream.title),
        subtitle: Text(
          '${stream.broadcaster.isEmpty ? 'unknown' : stream.broadcaster} • '
          '${stream.listenerCount} listening',
        ),
        trailing: Container(
          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
          decoration: BoxDecoration(
            color: theme.colorScheme.errorContainer,
            borderRadius: BorderRadius.circular(12),
          ),
          child: Text('LIVE',
              style: theme.textTheme.labelSmall
                  ?.copyWith(color: theme.colorScheme.onErrorContainer)),
        ),
        onTap: () {
          context.read<PlayerCubit>().playStream(stream);
          context.push('/player');
        },
      ),
    );
  }
}

class _Empty extends StatelessWidget {
  final String message;
  const _Empty({required this.message});

  @override
  Widget build(BuildContext context) {
    return ListView(
      children: [
        const SizedBox(height: 120),
        Icon(Icons.radio, size: 64, color: Theme.of(context).disabledColor),
        const SizedBox(height: 16),
        Center(child: Text(message)),
      ],
    );
  }
}

class _ErrorRetry extends StatelessWidget {
  final String message;
  final VoidCallback onRetry;
  const _ErrorRetry({required this.message, required this.onRetry});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(message),
          const SizedBox(height: 12),
          FilledButton(onPressed: onRetry, child: const Text('Retry')),
        ],
      ),
    );
  }
}
