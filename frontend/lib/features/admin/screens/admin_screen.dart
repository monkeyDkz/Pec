import 'package:flutter/material.dart';
import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/admin/repositories/admin_repository.dart';

class AdminScreen extends StatelessWidget {
  final ApiClient apiClient;
  const AdminScreen({super.key, required this.apiClient});

  @override
  Widget build(BuildContext context) {
    final repo = AdminRepository(apiClient: apiClient);
    return DefaultTabController(
      length: 3,
      child: Scaffold(
        appBar: AppBar(
          title: const Text('Administration'),
          bottom: const TabBar(tabs: [
            Tab(icon: Icon(Icons.insights), text: 'Stats'),
            Tab(icon: Icon(Icons.people), text: 'Utilisateurs'),
            Tab(icon: Icon(Icons.feedback_outlined), text: 'Feedbacks'),
          ]),
        ),
        body: TabBarView(children: [
          _StatsTab(repo: repo),
          _UsersTab(repo: repo),
          _FeedbackTab(repo: repo),
        ]),
      ),
    );
  }
}

class _StatsTab extends StatefulWidget {
  final AdminRepository repo;
  const _StatsTab({required this.repo});
  @override
  State<_StatsTab> createState() => _StatsTabState();
}

class _StatsTabState extends State<_StatsTab> {
  late Future<AdminStats> _future;
  @override
  void initState() {
    super.initState();
    _future = widget.repo.stats();
  }

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
      onRefresh: () async => setState(() => _future = widget.repo.stats()),
      child: FutureBuilder<AdminStats>(
        future: _future,
        builder: (context, snap) {
          if (snap.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snap.hasError) {
            return ListView(children: [
              Padding(
                padding: const EdgeInsets.all(24),
                child: Text('Erreur : ${snap.error}'),
              ),
            ]);
          }
          final s = snap.data!;
          final items = [
            ('Utilisateurs', s.totalUsers, Icons.people),
            ('Diffuseurs', s.totalBroadcasters, Icons.mic),
            ('Streams', s.totalStreams, Icons.podcasts),
            ('En direct maintenant', s.liveStreams, Icons.fiber_manual_record),
            ('Playlists', s.totalPlaylists, Icons.queue_music),
            ('Tracks', s.totalTracks, Icons.music_note),
            ('Feedbacks', s.totalFeedbacks, Icons.feedback),
          ];
          return GridView.builder(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: const EdgeInsets.all(16),
            gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
              maxCrossAxisExtent: 220,
              childAspectRatio: 1.1,
              mainAxisSpacing: 12,
              crossAxisSpacing: 12,
            ),
            itemCount: items.length,
            itemBuilder: (_, i) {
              final (label, value, icon) = items[i];
              return Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(icon, size: 32),
                      const SizedBox(height: 8),
                      Text('$value',
                          style: Theme.of(context).textTheme.headlineMedium),
                      Text(label,
                          style: Theme.of(context).textTheme.bodyMedium,
                          textAlign: TextAlign.center),
                    ],
                  ),
                ),
              );
            },
          );
        },
      ),
    );
  }
}

class _UsersTab extends StatefulWidget {
  final AdminRepository repo;
  const _UsersTab({required this.repo});
  @override
  State<_UsersTab> createState() => _UsersTabState();
}

class _UsersTabState extends State<_UsersTab> {
  late Future<List<AdminUser>> _future;
  @override
  void initState() {
    super.initState();
    _reload();
  }

  void _reload() => setState(() => _future = widget.repo.listUsers());

  Future<void> _changeRole(AdminUser u) async {
    final role = await showDialog<String>(
      context: context,
      builder: (dialogCtx) => AlertDialog(
        title: Text('Rôle de ${u.username}'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            for (final r in ['user', 'broadcaster', 'admin'])
              ListTile(
                title: Text(r),
                leading: Radio<String>(
                  value: r,
                  groupValue: u.role,
                  onChanged: (v) => Navigator.of(dialogCtx).pop(v),
                ),
              ),
          ],
        ),
      ),
    );
    if (role != null && role != u.role) {
      await widget.repo.updateRole(u.id, role);
      _reload();
    }
  }

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
      onRefresh: () async => _reload(),
      child: FutureBuilder<List<AdminUser>>(
        future: _future,
        builder: (context, snap) {
          if (snap.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snap.hasError) {
            return ListView(children: [
              Padding(padding: const EdgeInsets.all(24), child: Text('${snap.error}'))
            ]);
          }
          final users = snap.data!;
          return ListView.separated(
            physics: const AlwaysScrollableScrollPhysics(),
            itemCount: users.length,
            separatorBuilder: (_, __) => const Divider(height: 1),
            itemBuilder: (_, i) {
              final u = users[i];
              return ListTile(
                leading: CircleAvatar(child: Text(u.username.isEmpty ? '?' : u.username[0])),
                title: Text(u.username),
                subtitle: Text(u.email),
                trailing: Chip(label: Text(u.role)),
                onTap: () => _changeRole(u),
              );
            },
          );
        },
      ),
    );
  }
}

class _FeedbackTab extends StatefulWidget {
  final AdminRepository repo;
  const _FeedbackTab({required this.repo});
  @override
  State<_FeedbackTab> createState() => _FeedbackTabState();
}

class _FeedbackTabState extends State<_FeedbackTab> {
  late Future<List<AdminFeedback>> _future;
  @override
  void initState() {
    super.initState();
    _future = widget.repo.listFeedback();
  }

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
      onRefresh: () async => setState(() => _future = widget.repo.listFeedback()),
      child: FutureBuilder<List<AdminFeedback>>(
        future: _future,
        builder: (context, snap) {
          if (snap.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snap.hasError) {
            return ListView(children: [
              Padding(padding: const EdgeInsets.all(24), child: Text('${snap.error}'))
            ]);
          }
          final fbs = snap.data!;
          if (fbs.isEmpty) {
            return const Center(child: Text('Aucun feedback pour l\'instant'));
          }
          return ListView.separated(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: const EdgeInsets.all(8),
            itemCount: fbs.length,
            separatorBuilder: (_, __) => const SizedBox(height: 4),
            itemBuilder: (_, i) {
              final f = fbs[i];
              return Card(
                child: ListTile(
                  leading: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Icon(Icons.star, color: Colors.amber),
                      Text('${f.rating}/5'),
                    ],
                  ),
                  title: Text(f.comment.isEmpty ? '(sans commentaire)' : f.comment),
                  subtitle: Text(f.createdAt.toLocal().toString()),
                ),
              );
            },
          );
        },
      ),
    );
  }
}
