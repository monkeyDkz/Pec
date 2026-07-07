import 'package:flutter/material.dart';
import 'package:streampulse/core/api/api_client.dart';

Future<void> showFeedbackDialog(BuildContext context, {required ApiClient apiClient}) async {
  await showDialog<void>(
    context: context,
    builder: (_) => _FeedbackDialog(apiClient: apiClient),
  );
}

class _FeedbackDialog extends StatefulWidget {
  final ApiClient apiClient;
  const _FeedbackDialog({required this.apiClient});

  @override
  State<_FeedbackDialog> createState() => _FeedbackDialogState();
}

class _FeedbackDialogState extends State<_FeedbackDialog> {
  int _rating = 5;
  final _commentCtrl = TextEditingController();
  bool _submitting = false;

  @override
  void dispose() {
    _commentCtrl.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    setState(() => _submitting = true);
    try {
      await widget.apiClient.dio.post('/feedback', data: {
        'rating': _rating,
        'comment': _commentCtrl.text.trim(),
      });
      if (mounted) {
        Navigator.of(context).pop();
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Merci pour votre retour !')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Échec : $e')),
        );
        setState(() => _submitting = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Votre avis'),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('Quelle note donneriez-vous à StreamPulse ?'),
          const SizedBox(height: 12),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              for (var i = 1; i <= 5; i++)
                IconButton(
                  icon: Icon(
                    i <= _rating ? Icons.star : Icons.star_border,
                    color: Colors.amber,
                    size: 32,
                  ),
                  onPressed: () => setState(() => _rating = i),
                  tooltip: '$i étoile(s)',
                ),
            ],
          ),
          const SizedBox(height: 8),
          TextField(
            controller: _commentCtrl,
            decoration: const InputDecoration(
              labelText: 'Commentaire (facultatif)',
              border: OutlineInputBorder(),
            ),
            maxLines: 3,
            maxLength: 2000,
          ),
        ],
      ),
      actions: [
        TextButton(
          onPressed: _submitting ? null : () => Navigator.of(context).pop(),
          child: const Text('Annuler'),
        ),
        FilledButton(
          onPressed: _submitting ? null : _submit,
          child: _submitting
              ? const SizedBox(
                  height: 18,
                  width: 18,
                  child: CircularProgressIndicator(strokeWidth: 2),
                )
              : const Text('Envoyer'),
        ),
      ],
    );
  }
}
