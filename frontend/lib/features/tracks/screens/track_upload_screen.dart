import 'dart:io';

import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/tracks/repositories/track_repository.dart';

class TrackUploadScreen extends StatefulWidget {
  final ApiClient apiClient;
  const TrackUploadScreen({super.key, required this.apiClient});

  @override
  State<TrackUploadScreen> createState() => _TrackUploadScreenState();
}

class _TrackUploadScreenState extends State<TrackUploadScreen> {
  final _titleCtrl = TextEditingController();
  final _artistCtrl = TextEditingController();
  File? _pickedFile;
  double _progress = 0;
  bool _uploading = false;

  late final TrackRepository _repo = TrackRepository(apiClient: widget.apiClient);

  @override
  void dispose() {
    _titleCtrl.dispose();
    _artistCtrl.dispose();
    super.dispose();
  }

  Future<void> _pickFile() async {
    final result = await FilePicker.pickFiles(
      type: FileType.audio,
      allowMultiple: false,
    );
    if (result != null && result.files.single.path != null) {
      setState(() {
        _pickedFile = File(result.files.single.path!);
        if (_titleCtrl.text.isEmpty) {
          _titleCtrl.text = result.files.single.name.replaceAll(
              RegExp(r'\.(mp3|aac|ogg|wav)$', caseSensitive: false), '');
        }
      });
    }
  }

  Future<void> _upload() async {
    if (_pickedFile == null || _titleCtrl.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Fichier et titre requis')),
      );
      return;
    }
    setState(() {
      _uploading = true;
      _progress = 0;
    });
    try {
      await _repo.upload(
        file: _pickedFile!,
        title: _titleCtrl.text.trim(),
        artist: _artistCtrl.text.trim(),
        onProgress: (sent, total) {
          if (total > 0) setState(() => _progress = sent / total);
        },
      );
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Track téléversé avec succès')),
      );
      Navigator.of(context).pop();
    } catch (e) {
      if (mounted) {
        setState(() => _uploading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Échec : $e')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Téléverser un track')),
      body: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            OutlinedButton.icon(
              onPressed: _uploading ? null : _pickFile,
              icon: const Icon(Icons.audio_file),
              label: Text(_pickedFile == null
                  ? 'Choisir un fichier audio'
                  : _pickedFile!.uri.pathSegments.last),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _titleCtrl,
              decoration: const InputDecoration(labelText: 'Titre'),
              enabled: !_uploading,
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _artistCtrl,
              decoration: const InputDecoration(labelText: 'Artiste (facultatif)'),
              enabled: !_uploading,
            ),
            const SizedBox(height: 24),
            if (_uploading) ...[
              LinearProgressIndicator(value: _progress),
              const SizedBox(height: 8),
              Text('${(_progress * 100).toStringAsFixed(0)} %',
                  textAlign: TextAlign.center),
              const SizedBox(height: 16),
            ],
            FilledButton.icon(
              onPressed: _uploading ? null : _upload,
              icon: const Icon(Icons.cloud_upload),
              label: const Text('Téléverser'),
            ),
            const SizedBox(height: 24),
            Card(
              color: Theme.of(context).colorScheme.surfaceContainerHighest,
              child: const Padding(
                padding: EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(children: [
                      Icon(Icons.info_outline, size: 16),
                      SizedBox(width: 8),
                      Text('Formats acceptés', style: TextStyle(fontWeight: FontWeight.bold)),
                    ]),
                    SizedBox(height: 4),
                    Text('• MP3, AAC, OGG, WAV — 50 Mo maximum'),
                    Text('• Rôle broadcaster ou admin requis'),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
