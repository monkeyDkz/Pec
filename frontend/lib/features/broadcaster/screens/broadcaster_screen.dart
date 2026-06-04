import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/features/broadcaster/cubit/broadcaster_cubit.dart';
import 'package:streampulse/features/profile/cubit/current_user_cubit.dart';
import 'package:streampulse/features/streams/repositories/stream_repository.dart';

class BroadcasterScreen extends StatelessWidget {
  const BroadcasterScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (ctx) => BroadcasterCubit(
        repository: ctx.read<StreamRepository>(),
        currentUserId: ctx.read<CurrentUserCubit>().state.user?.id ?? '',
      )..load(),
      child: const _BroadcasterView(),
    );
  }
}

class _BroadcasterView extends StatelessWidget {
  const _BroadcasterView();

  Future<void> _showCreateDialog(BuildContext context) async {
    final cubit = context.read<BroadcasterCubit>();
    final title = TextEditingController();
    final desc = TextEditingController();
    await showDialog<void>(
      context: context,
      builder: (dialogCtx) => AlertDialog(
        title: const Text('Go live'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: title,
              decoration: const InputDecoration(labelText: 'Title'),
            ),
            TextField(
              controller: desc,
              decoration: const InputDecoration(labelText: 'Description'),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogCtx).pop(),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () {
              if (title.text.trim().isNotEmpty) {
                cubit.create(title.text.trim(), desc.text.trim());
              }
              Navigator.of(dialogCtx).pop();
            },
            child: const Text('Start'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Broadcast'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => context.read<BroadcasterCubit>().load(),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _showCreateDialog(context),
        icon: const Icon(Icons.podcasts),
        label: const Text('Go live'),
      ),
      body: BlocBuilder<BroadcasterCubit, BroadcasterState>(
        builder: (context, state) {
          if (state.loading) {
            return const Center(child: CircularProgressIndicator());
          }
          if (state.mine.isEmpty) {
            return const Center(
              child: Padding(
                padding: EdgeInsets.all(24),
                child: Text('You have no live streams.\nTap "Go live" to start broadcasting.',
                    textAlign: TextAlign.center),
              ),
            );
          }
          return ListView(
            padding: const EdgeInsets.all(8),
            children: [
              for (final s in state.mine)
                Card(
                  child: ListTile(
                    leading: const CircleAvatar(child: Icon(Icons.sensors)),
                    title: Text(s.title),
                    subtitle: Text('${s.listenerCount} listening • LIVE'),
                    trailing: IconButton(
                      icon: const Icon(Icons.stop_circle, color: Colors.red),
                      onPressed: () => context.read<BroadcasterCubit>().stop(s.id),
                    ),
                  ),
                ),
            ],
          );
        },
      ),
    );
  }
}
