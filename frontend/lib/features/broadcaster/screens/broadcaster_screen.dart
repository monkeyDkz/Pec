import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/broadcaster/bloc/broadcaster_bloc.dart';
import 'package:streampulse/features/broadcaster/repositories/broadcaster_repository.dart';

class BroadcasterScreen extends StatelessWidget {
  final ApiClient apiClient;
  const BroadcasterScreen({super.key, required this.apiClient});

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (_) => BroadcasterBloc(
        repository: BroadcasterRepository(apiClient: apiClient),
      ),
      child: const _BroadcasterView(),
    );
  }
}

class _BroadcasterView extends StatefulWidget {
  const _BroadcasterView();

  @override
  State<_BroadcasterView> createState() => _BroadcasterViewState();
}

class _BroadcasterViewState extends State<_BroadcasterView> {
  final _formKey = GlobalKey<FormState>();
  final _titleCtrl = TextEditingController();
  final _descCtrl = TextEditingController();

  @override
  void dispose() {
    _titleCtrl.dispose();
    _descCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Diffuser un live')),
      body: BlocConsumer<BroadcasterBloc, BroadcasterState>(
        listener: (context, state) {
          if (state is BroadcasterError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text(state.message)),
            );
          }
        },
        builder: (context, state) {
          if (state is BroadcasterLive) {
            return _LivePanel(state: state);
          }
          if (state is BroadcasterReady) {
            return _ReadyPanel(stream: state.stream);
          }
          return _CreatePanel(
            formKey: _formKey,
            titleCtrl: _titleCtrl,
            descCtrl: _descCtrl,
            loading: state is BroadcasterCreating,
            onSubmit: () {
              if (_formKey.currentState?.validate() ?? false) {
                context.read<BroadcasterBloc>().add(BroadcasterCreateRequested(
                      title: _titleCtrl.text.trim(),
                      description: _descCtrl.text.trim(),
                    ));
              }
            },
          );
        },
      ),
    );
  }
}

class _CreatePanel extends StatelessWidget {
  final GlobalKey<FormState> formKey;
  final TextEditingController titleCtrl;
  final TextEditingController descCtrl;
  final bool loading;
  final VoidCallback onSubmit;
  const _CreatePanel({
    required this.formKey,
    required this.titleCtrl,
    required this.descCtrl,
    required this.loading,
    required this.onSubmit,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(24),
      child: Form(
        key: formKey,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text("Préparez votre session de diffusion",
                style: Theme.of(context).textTheme.titleMedium),
            const SizedBox(height: 16),
            TextFormField(
              controller: titleCtrl,
              decoration: const InputDecoration(labelText: 'Titre'),
              validator: (v) => (v == null || v.trim().isEmpty) ? 'Titre requis' : null,
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: descCtrl,
              decoration: const InputDecoration(labelText: 'Description'),
              maxLines: 3,
            ),
            const SizedBox(height: 24),
            FilledButton.icon(
              onPressed: loading ? null : onSubmit,
              icon: const Icon(Icons.add),
              label: Text(loading ? "Création..." : "Créer le stream"),
            ),
          ],
        ),
      ),
    );
  }
}

class _ReadyPanel extends StatelessWidget {
  final dynamic stream;
  const _ReadyPanel({required this.stream});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.mic_none, size: 96, color: theme.colorScheme.primary),
          const SizedBox(height: 16),
          Text(stream.title, style: theme.textTheme.titleLarge),
          const SizedBox(height: 8),
          Text(stream.description, style: theme.textTheme.bodyMedium),
          const SizedBox(height: 32),
          FilledButton.icon(
            style: FilledButton.styleFrom(
              backgroundColor: Colors.red,
              minimumSize: const Size(200, 56),
            ),
            onPressed: () => context.read<BroadcasterBloc>().add(const BroadcasterStartRequested()),
            icon: const Icon(Icons.fiber_manual_record),
            label: const Text('Démarrer la diffusion'),
          ),
        ],
      ),
    );
  }
}

class _LivePanel extends StatelessWidget {
  final BroadcasterLive state;
  const _LivePanel({required this.state});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final h = state.elapsed.inHours.toString().padLeft(2, '0');
    final m = (state.elapsed.inMinutes % 60).toString().padLeft(2, '0');
    final s = (state.elapsed.inSeconds % 60).toString().padLeft(2, '0');

    return Padding(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
            decoration: BoxDecoration(
              color: Colors.red,
              borderRadius: BorderRadius.circular(4),
            ),
            child: const Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(Icons.fiber_manual_record, color: Colors.white, size: 12),
                SizedBox(width: 6),
                Text('LIVE',
                    style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
              ],
            ),
          ),
          const SizedBox(height: 12),
          Chip(
            avatar: Icon(
              state.usingMicrophone ? Icons.mic : Icons.science_outlined,
              size: 16,
            ),
            label: Text(
              state.usingMicrophone ? 'Capture micro' : 'Audio synthétique (démo)',
            ),
          ),
          const SizedBox(height: 24),
          Text(state.stream.title, style: theme.textTheme.titleLarge),
          const SizedBox(height: 8),
          Text('$h:$m:$s', style: theme.textTheme.displaySmall),
          const SizedBox(height: 8),
          Text('${state.stream.listenerCount} auditeur(s)',
              style: theme.textTheme.bodyLarge),
          const SizedBox(height: 48),
          FilledButton.icon(
            style: FilledButton.styleFrom(
              backgroundColor: theme.colorScheme.error,
              minimumSize: const Size(200, 56),
            ),
            onPressed: () => context.read<BroadcasterBloc>().add(const BroadcasterStopRequested()),
            icon: const Icon(Icons.stop),
            label: const Text('Arrêter la diffusion'),
          ),
        ],
      ),
    );
  }
}
