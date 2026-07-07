import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/auth/bloc/auth_bloc.dart';
import 'package:streampulse/features/feedback/feedback_dialog.dart';
import 'package:streampulse/features/settings/screens/profile_screen.dart';
import 'package:streampulse/features/tracks/screens/track_upload_screen.dart';

class SettingsScreen extends StatelessWidget {
  final ApiClient apiClient;
  const SettingsScreen({super.key, required this.apiClient});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Paramètres')),
      body: ListView(
        children: [
          ListTile(
            leading: const Icon(Icons.account_circle),
            title: const Text('Mon profil'),
            subtitle: const Text('Email, nom d\'utilisateur, rôle'),
            onTap: () => Navigator.of(context).push(MaterialPageRoute(
              builder: (_) => ProfileScreen(apiClient: apiClient),
            )),
          ),
          ListTile(
            leading: const Icon(Icons.upload_file),
            title: const Text('Téléverser un track'),
            subtitle: const Text('Pour broadcasters et admins'),
            onTap: () => Navigator.of(context).push(MaterialPageRoute(
              builder: (_) => TrackUploadScreen(apiClient: apiClient),
            )),
          ),
          const Divider(height: 1),
          ListTile(
            leading: const Icon(Icons.feedback_outlined),
            title: const Text('Envoyer un feedback'),
            onTap: () => showFeedbackDialog(context, apiClient: apiClient),
          ),
          const Divider(height: 1),
          _SectionHeader(title: 'Vie privée (RGPD)', icon: Icons.privacy_tip),
          ListTile(
            leading: const Icon(Icons.download),
            title: const Text('Exporter mes données'),
            subtitle: const Text('Droit d\'accès — JSON complet de votre compte'),
            onTap: () => _exportData(context, apiClient),
          ),
          ListTile(
            leading: const Icon(Icons.delete_forever, color: Colors.red),
            title: const Text('Supprimer mon compte',
                style: TextStyle(color: Colors.red)),
            subtitle: const Text('Droit à l\'oubli — action irréversible'),
            onTap: () => _confirmDeleteAccount(context),
          ),
          const Divider(height: 1),
          _SectionHeader(title: 'Compte', icon: Icons.security),
          ListTile(
            leading: const Icon(Icons.logout),
            title: const Text('Se déconnecter'),
            onTap: () => context.read<AuthBloc>().add(AuthLogoutRequested()),
          ),
          const Divider(height: 1),
          const _SectionHeader(title: 'À propos', icon: Icons.info_outline),
          const ListTile(
            title: Text('StreamPulse'),
            subtitle: Text('v1.0.0 — École .decode — RNCP 38822'),
          ),
        ],
      ),
    );
  }

  Future<void> _exportData(BuildContext context, ApiClient apiClient) async {
    try {
      final response = await apiClient.dio.get('/users/me/data');
      final pretty = const JsonEncoder.withIndent('  ').convert(response.data);
      if (!context.mounted) return;
      await showDialog<void>(
        context: context,
        builder: (dialogCtx) => AlertDialog(
          title: const Text('Vos données'),
          content: SizedBox(
            width: 600,
            height: 400,
            child: SingleChildScrollView(
              child: SelectableText(pretty, style: const TextStyle(fontFamily: 'monospace', fontSize: 12)),
            ),
          ),
          actions: [
            TextButton.icon(
              icon: const Icon(Icons.copy),
              label: const Text('Copier'),
              onPressed: () async {
                await Clipboard.setData(ClipboardData(text: pretty));
                if (dialogCtx.mounted) {
                  ScaffoldMessenger.of(dialogCtx).showSnackBar(
                    const SnackBar(content: Text('Données copiées dans le presse-papier')),
                  );
                }
              },
            ),
            TextButton(
              onPressed: () => Navigator.of(dialogCtx).pop(),
              child: const Text('Fermer'),
            ),
          ],
        ),
      );
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erreur d\'export : $e')),
        );
      }
    }
  }

  Future<void> _confirmDeleteAccount(BuildContext context) async {
    final authBloc = context.read<AuthBloc>();
    final apiClient = this.apiClient;

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (dialogCtx) => AlertDialog(
        title: const Text('Supprimer votre compte ?'),
        content: const Text(
          "Toutes vos données (compte, streams, playlists, tracks, feedbacks) "
          "seront supprimées définitivement.\n\nCette action est irréversible.",
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogCtx).pop(false),
            child: const Text('Annuler'),
          ),
          FilledButton(
            style: FilledButton.styleFrom(backgroundColor: Colors.red),
            onPressed: () => Navigator.of(dialogCtx).pop(true),
            child: const Text('Supprimer définitivement'),
          ),
        ],
      ),
    );
    if (confirmed != true) return;

    try {
      await apiClient.dio.delete('/users/me');
      authBloc.add(AuthLogoutRequested());
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Échec de la suppression : $e')),
        );
      }
    }
  }
}

class _SectionHeader extends StatelessWidget {
  final String title;
  final IconData icon;
  const _SectionHeader({required this.title, required this.icon});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 24, 16, 8),
      child: Row(
        children: [
          Icon(icon, size: 18, color: Theme.of(context).colorScheme.primary),
          const SizedBox(width: 8),
          Text(title,
              style: Theme.of(context).textTheme.titleSmall?.copyWith(
                    color: Theme.of(context).colorScheme.primary,
                  )),
        ],
      ),
    );
  }
}
