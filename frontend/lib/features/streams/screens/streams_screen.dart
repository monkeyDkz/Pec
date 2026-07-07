import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:shimmer/shimmer.dart';
import 'package:streampulse/features/auth/bloc/auth_bloc.dart';
import 'package:streampulse/features/streams/bloc/streams_bloc.dart';
import 'package:streampulse/features/streams/models/stream_model.dart';

class StreamsScreen extends StatelessWidget {
  const StreamsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Streams en direct'),
        actions: [
          IconButton(
            tooltip: 'Se déconnecter',
            icon: const Icon(Icons.logout),
            onPressed: () => context.read<AuthBloc>().add(AuthLogoutRequested()),
          ),
        ],
      ),
      body: BlocBuilder<StreamsBloc, StreamsState>(
        builder: (context, state) {
          if (state is StreamsLoading || state is StreamsInitial) {
            return const _StreamsShimmer();
          }
          if (state is StreamsError) {
            return _ErrorRetry(
              message: state.message,
              onRetry: () => context.read<StreamsBloc>().add(const StreamsLoadRequested()),
            );
          }
          final streams = (state as StreamsLoaded).streams;
          if (streams.isEmpty) {
            return _EmptyState(
              onRetry: () => context.read<StreamsBloc>().add(const StreamsRefreshRequested()),
            );
          }
          return RefreshIndicator(
            onRefresh: () async {
              context.read<StreamsBloc>().add(const StreamsRefreshRequested());
            },
            child: ListView.separated(
              physics: const AlwaysScrollableScrollPhysics(),
              padding: const EdgeInsets.all(16),
              itemCount: streams.length,
              separatorBuilder: (_, __) => const SizedBox(height: 12),
              itemBuilder: (_, i) => _StreamTile(stream: streams[i]),
            ),
          );
        },
      ),
    );
  }
}

class _StreamTile extends StatelessWidget {
  final StreamModel stream;
  const _StreamTile({required this.stream});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Semantics(
      label: 'Stream ${stream.title}, animé par ${stream.broadcasterUsername}, '
          '${stream.listenerCount} auditeurs, ${stream.isLive ? 'en direct' : 'hors ligne'}',
      button: true,
      child: Card(
        child: ListTile(
          onTap: stream.isLive ? () => context.go('/player/${stream.id}') : null,
          leading: CircleAvatar(
            backgroundColor: stream.isLive ? Colors.red : theme.disabledColor,
            child: const Icon(Icons.podcasts, color: Colors.white),
          ),
          title: Text(stream.title, style: theme.textTheme.titleMedium),
          subtitle: Text(
            stream.broadcasterUsername.isEmpty
                ? stream.description
                : '@${stream.broadcasterUsername}',
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          trailing: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              if (stream.isLive)
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color: Colors.red,
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: const Text(
                    'LIVE',
                    style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 11),
                  ),
                ),
              const SizedBox(height: 4),
              Text('${stream.listenerCount} 👤', style: theme.textTheme.bodySmall),
            ],
          ),
        ),
      ),
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
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline, size: 48),
            const SizedBox(height: 8),
            Text('Erreur : $message', textAlign: TextAlign.center),
            const SizedBox(height: 16),
            FilledButton(onPressed: onRetry, child: const Text('Réessayer')),
          ],
        ),
      ),
    );
  }
}

class _EmptyState extends StatelessWidget {
  final VoidCallback onRetry;
  const _EmptyState({required this.onRetry});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Icon(Icons.radio_outlined, size: 64),
          const SizedBox(height: 8),
          const Text('Aucun stream en direct pour le moment'),
          const SizedBox(height: 16),
          OutlinedButton.icon(
            onPressed: onRetry,
            icon: const Icon(Icons.refresh),
            label: const Text('Rafraîchir'),
          ),
        ],
      ),
    );
  }
}

class _StreamsShimmer extends StatelessWidget {
  const _StreamsShimmer();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final base = theme.colorScheme.surfaceContainerHighest;
    final highlight = theme.colorScheme.surface;
    return ListView.separated(
      padding: const EdgeInsets.all(16),
      itemCount: 6,
      separatorBuilder: (_, __) => const SizedBox(height: 12),
      itemBuilder: (_, __) => Shimmer.fromColors(
        baseColor: base,
        highlightColor: highlight,
        child: Card(
          child: SizedBox(
            height: 76,
            child: Padding(
              padding: const EdgeInsets.all(12),
              child: Row(
                children: [
                  CircleAvatar(backgroundColor: base, radius: 24),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Container(height: 14, width: double.infinity, color: base),
                        const SizedBox(height: 8),
                        Container(height: 12, width: 120, color: base),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
