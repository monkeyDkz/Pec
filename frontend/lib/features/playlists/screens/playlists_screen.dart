import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:streampulse/features/playlists/cubit/playlists_cubit.dart';
import 'package:streampulse/features/playlists/repositories/playlist_repository.dart';

class PlaylistsScreen extends StatelessWidget {
  const PlaylistsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (ctx) => PlaylistsCubit(ctx.read<PlaylistRepository>())..load(),
      child: const _PlaylistsView(),
    );
  }
}

class _PlaylistsView extends StatelessWidget {
  const _PlaylistsView();

  Future<void> _showCreateDialog(BuildContext context) async {
    final cubit = context.read<PlaylistsCubit>();
    final name = TextEditingController();
    final desc = TextEditingController();
    await showDialog<void>(
      context: context,
      builder: (dialogCtx) => AlertDialog(
        title: const Text('New playlist'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(controller: name, decoration: const InputDecoration(labelText: 'Name')),
            TextField(
                controller: desc, decoration: const InputDecoration(labelText: 'Description')),
          ],
        ),
        actions: [
          TextButton(
              onPressed: () => Navigator.of(dialogCtx).pop(), child: const Text('Cancel')),
          FilledButton(
            onPressed: () {
              if (name.text.trim().isNotEmpty) {
                cubit.create(name.text.trim(), desc.text.trim());
              }
              Navigator.of(dialogCtx).pop();
            },
            child: const Text('Create'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Playlists')),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showCreateDialog(context),
        child: const Icon(Icons.add),
      ),
      body: BlocBuilder<PlaylistsCubit, PlaylistsState>(
        builder: (context, state) {
          if (state.loading) {
            return const Center(child: CircularProgressIndicator());
          }
          if (state.playlists.isEmpty) {
            return const Center(child: Text('No playlists yet. Tap + to create one.'));
          }
          return ListView(
            padding: const EdgeInsets.all(8),
            children: [
              for (final p in state.playlists)
                Card(
                  child: ListTile(
                    leading: const CircleAvatar(child: Icon(Icons.queue_music)),
                    title: Text(p.name),
                    subtitle: Text('${p.tracks.length} tracks'),
                    trailing: IconButton(
                      icon: const Icon(Icons.delete_outline),
                      onPressed: () => context.read<PlaylistsCubit>().delete(p.id),
                    ),
                    onTap: () async {
                      await context.push('/playlist/${p.id}');
                      if (context.mounted) context.read<PlaylistsCubit>().load();
                    },
                  ),
                ),
            ],
          );
        },
      ),
    );
  }
}
