import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/playlists/bloc/playlists_bloc.dart';
import 'package:streampulse/features/playlists/models/playlist_model.dart';
import 'package:streampulse/features/playlists/repositories/playlist_repository.dart';
import 'package:streampulse/features/playlists/screens/playlist_detail_screen.dart';

class PlaylistsScreen extends StatelessWidget {
  final ApiClient apiClient;
  const PlaylistsScreen({super.key, required this.apiClient});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => PlaylistsBloc(
        repository: PlaylistRepository(apiClient: apiClient),
      )..add(const PlaylistsLoadRequested()),
      child: const _PlaylistsView(),
    );
  }
}

class _PlaylistsView extends StatelessWidget {
  const _PlaylistsView();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Mes playlists')),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _showCreateDialog(context),
        icon: const Icon(Icons.playlist_add),
        label: const Text('Nouvelle'),
      ),
      body: BlocBuilder<PlaylistsBloc, PlaylistsState>(
        builder: (context, state) {
          if (state is PlaylistsLoading || state is PlaylistsInitial) {
            return const Center(child: CircularProgressIndicator());
          }
          if (state is PlaylistsError) {
            return Center(
              child: Padding(
                padding: const EdgeInsets.all(24),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Icon(Icons.error_outline, size: 48),
                    const SizedBox(height: 8),
                    Text(state.message, textAlign: TextAlign.center),
                    const SizedBox(height: 16),
                    FilledButton(
                      onPressed: () => context
                          .read<PlaylistsBloc>()
                          .add(const PlaylistsLoadRequested()),
                      child: const Text('Réessayer'),
                    ),
                  ],
                ),
              ),
            );
          }
          final playlists = (state as PlaylistsLoaded).playlists;
          if (playlists.isEmpty) {
            return Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Icon(Icons.playlist_play, size: 64),
                  const SizedBox(height: 8),
                  Text("Aucune playlist pour l'instant",
                      style: Theme.of(context).textTheme.bodyLarge),
                ],
              ),
            );
          }
          return RefreshIndicator(
            onRefresh: () async => context
                .read<PlaylistsBloc>()
                .add(const PlaylistsLoadRequested()),
            child: ListView.separated(
              physics: const AlwaysScrollableScrollPhysics(),
              padding: const EdgeInsets.all(16),
              itemCount: playlists.length,
              separatorBuilder: (_, __) => const SizedBox(height: 8),
              itemBuilder: (_, i) => _PlaylistTile(playlist: playlists[i]),
            ),
          );
        },
      ),
    );
  }

  void _showCreateDialog(BuildContext context) {
    final formKey = GlobalKey<FormState>();
    final nameCtrl = TextEditingController();
    final descCtrl = TextEditingController();
    final bloc = context.read<PlaylistsBloc>();

    showDialog<void>(
      context: context,
      builder: (dialogCtx) => AlertDialog(
        title: const Text('Nouvelle playlist'),
        content: Form(
          key: formKey,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextFormField(
                controller: nameCtrl,
                decoration: const InputDecoration(labelText: 'Nom'),
                validator: (v) => (v == null || v.trim().isEmpty) ? 'Nom requis' : null,
              ),
              const SizedBox(height: 8),
              TextFormField(
                controller: descCtrl,
                decoration: const InputDecoration(labelText: 'Description'),
              ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogCtx).pop(),
            child: const Text('Annuler'),
          ),
          FilledButton(
            onPressed: () {
              if (formKey.currentState?.validate() ?? false) {
                bloc.add(PlaylistCreateRequested(
                  name: nameCtrl.text.trim(),
                  description: descCtrl.text.trim(),
                ));
                Navigator.of(dialogCtx).pop();
              }
            },
            child: const Text('Créer'),
          ),
        ],
      ),
    );
  }
}

class _PlaylistTile extends StatelessWidget {
  final PlaylistModel playlist;
  const _PlaylistTile({required this.playlist});

  @override
  Widget build(BuildContext context) {
    final apiClient = context.findAncestorWidgetOfExactType<PlaylistsScreen>()!.apiClient;
    return Card(
      child: ListTile(
        leading: const CircleAvatar(child: Icon(Icons.queue_music)),
        title: Text(playlist.name),
        subtitle: Text('${playlist.tracks.length} morceau(x)'),
        onTap: () async {
          final bloc = context.read<PlaylistsBloc>();
          await Navigator.of(context).push(MaterialPageRoute(
            builder: (_) => PlaylistDetailScreen(
              apiClient: apiClient,
              playlistId: playlist.id,
            ),
          ));
          bloc.add(const PlaylistsLoadRequested());
        },
        trailing: IconButton(
          icon: const Icon(Icons.delete_outline),
          tooltip: 'Supprimer',
          onPressed: () => _confirmDelete(context),
        ),
      ),
    );
  }

  void _confirmDelete(BuildContext context) {
    final bloc = context.read<PlaylistsBloc>();
    showDialog<void>(
      context: context,
      builder: (dialogCtx) => AlertDialog(
        title: const Text('Supprimer cette playlist ?'),
        content: Text('« ${playlist.name} » sera supprimée définitivement.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogCtx).pop(),
            child: const Text('Annuler'),
          ),
          FilledButton(
            style: FilledButton.styleFrom(backgroundColor: Colors.red),
            onPressed: () {
              bloc.add(PlaylistDeleteRequested(playlist.id));
              Navigator.of(dialogCtx).pop();
            },
            child: const Text('Supprimer'),
          ),
        ],
      ),
    );
  }
}
