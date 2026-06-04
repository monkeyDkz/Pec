import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/features/admin/cubit/admin_cubit.dart';
import 'package:streampulse/features/admin/repositories/admin_repository.dart';

class AdminScreen extends StatelessWidget {
  const AdminScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (ctx) => AdminCubit(ctx.read<AdminRepository>())..load(),
      child: const _AdminView(),
    );
  }
}

class _AdminView extends StatelessWidget {
  const _AdminView();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Admin'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => context.read<AdminCubit>().load(),
          ),
        ],
      ),
      body: BlocBuilder<AdminCubit, AdminState>(
        builder: (context, state) {
          if (state.loading) {
            return const Center(child: CircularProgressIndicator());
          }
          if (state.error != null) {
            return Center(child: Text(state.error!));
          }
          final stats = state.stats;
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              if (stats != null)
                Wrap(
                  spacing: 12,
                  runSpacing: 12,
                  children: [
                    _StatCard(label: 'Users', value: stats.totalUsers, icon: Icons.people),
                    _StatCard(label: 'Live', value: stats.liveStreams, icon: Icons.sensors),
                    _StatCard(label: 'Streams', value: stats.totalStreams, icon: Icons.radio),
                    _StatCard(label: 'Tracks', value: stats.totalTracks, icon: Icons.music_note),
                    _StatCard(
                        label: 'Playlists',
                        value: stats.totalPlaylists,
                        icon: Icons.queue_music),
                  ],
                ),
              const SizedBox(height: 24),
              Text('Users', style: Theme.of(context).textTheme.titleLarge),
              const Divider(),
              for (final u in state.users)
                ListTile(
                  leading: CircleAvatar(
                      child: Text(u.username.isNotEmpty ? u.username[0].toUpperCase() : '?')),
                  title: Text(u.username),
                  subtitle: Text(u.email),
                  trailing: DropdownButton<String>(
                    value: ['user', 'broadcaster', 'admin'].contains(u.role) ? u.role : 'user',
                    items: const [
                      DropdownMenuItem(value: 'user', child: Text('user')),
                      DropdownMenuItem(value: 'broadcaster', child: Text('broadcaster')),
                      DropdownMenuItem(value: 'admin', child: Text('admin')),
                    ],
                    onChanged: (role) {
                      if (role != null && role != u.role) {
                        context.read<AdminCubit>().updateRole(u.id, role);
                      }
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

class _StatCard extends StatelessWidget {
  final String label;
  final int value;
  final IconData icon;
  const _StatCard({required this.label, required this.value, required this.icon});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      width: 100,
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerHighest,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        children: [
          Icon(icon, color: theme.colorScheme.primary),
          const SizedBox(height: 8),
          Text('$value', style: theme.textTheme.headlineSmall),
          Text(label, style: theme.textTheme.bodySmall),
        ],
      ),
    );
  }
}
