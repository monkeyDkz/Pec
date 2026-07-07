import 'dart:async';
import 'dart:math';

import 'package:dio/dio.dart';
import 'package:equatable/equatable.dart';
import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:record/record.dart';
import 'package:streampulse/features/broadcaster/repositories/broadcaster_repository.dart';
import 'package:streampulse/features/streams/models/stream_model.dart';

// ── Events ────────────────────────────────────────────────────────────
abstract class BroadcasterEvent extends Equatable {
  const BroadcasterEvent();
  @override
  List<Object?> get props => [];
}

class BroadcasterCreateRequested extends BroadcasterEvent {
  final String title;
  final String description;
  const BroadcasterCreateRequested({required this.title, required this.description});
  @override
  List<Object?> get props => [title, description];
}

class BroadcasterStartRequested extends BroadcasterEvent {
  const BroadcasterStartRequested();
}

class BroadcasterStopRequested extends BroadcasterEvent {
  const BroadcasterStopRequested();
}

// ── States ────────────────────────────────────────────────────────────
abstract class BroadcasterState extends Equatable {
  const BroadcasterState();
  @override
  List<Object?> get props => [];
}

class BroadcasterIdle extends BroadcasterState {
  const BroadcasterIdle();
}

class BroadcasterCreating extends BroadcasterState {
  const BroadcasterCreating();
}

class BroadcasterReady extends BroadcasterState {
  final StreamModel stream;
  const BroadcasterReady(this.stream);
  @override
  List<Object?> get props => [stream];
}

class BroadcasterLive extends BroadcasterState {
  final StreamModel stream;
  final Duration elapsed;
  final bool usingMicrophone;
  const BroadcasterLive({
    required this.stream,
    required this.elapsed,
    required this.usingMicrophone,
  });
  @override
  List<Object?> get props => [stream, elapsed, usingMicrophone];
}

class BroadcasterError extends BroadcasterState {
  final String message;
  const BroadcasterError(this.message);
  @override
  List<Object?> get props => [message];
}

// ── BLoC ──────────────────────────────────────────────────────────────
class BroadcasterBloc extends Bloc<BroadcasterEvent, BroadcasterState> {
  final BroadcasterRepository repository;
  final AudioRecorder _recorder;

  StreamModel? _stream;
  CancelToken? _cancelToken;
  Timer? _elapsedTimer;
  DateTime? _startedAt;
  bool _usingMic = false;

  BroadcasterBloc({
    required this.repository,
    AudioRecorder? recorder,
  })  : _recorder = recorder ?? AudioRecorder(),
        super(const BroadcasterIdle()) {
    on<BroadcasterCreateRequested>(_onCreate);
    on<BroadcasterStartRequested>(_onStart);
    on<BroadcasterStopRequested>(_onStop);
  }

  Future<void> _onCreate(BroadcasterCreateRequested event, Emitter<BroadcasterState> emit) async {
    emit(const BroadcasterCreating());
    try {
      _stream = await repository.createStream(
        title: event.title,
        description: event.description,
      );
      emit(BroadcasterReady(_stream!));
    } catch (e) {
      emit(BroadcasterError(e.toString()));
    }
  }

  Future<void> _onStart(BroadcasterStartRequested event, Emitter<BroadcasterState> emit) async {
    final stream = _stream;
    if (stream == null) {
      emit(const BroadcasterError('No stream created'));
      return;
    }
    _cancelToken = CancelToken();
    _startedAt = DateTime.now();

    // Try real microphone capture; fall back to synthetic audio if unavailable
    // (e.g. simulator, web preview, permission refused).
    final audio = await _resolveAudioSource();
    _usingMic = audio.isMic;

    emit(BroadcasterLive(
      stream: stream,
      elapsed: Duration.zero,
      usingMicrophone: _usingMic,
    ));

    _elapsedTimer = Timer.periodic(const Duration(seconds: 1), (_) {
      if (state is BroadcasterLive) {
        emit(BroadcasterLive(
          stream: stream,
          elapsed: DateTime.now().difference(_startedAt!),
          usingMicrophone: _usingMic,
        ));
      }
    });

    unawaited(_publish(stream.id, audio.source, emit));
  }

  Future<_AudioSource> _resolveAudioSource() async {
    try {
      if (!await _recorder.hasPermission()) {
        return _AudioSource(_syntheticAudioSource(), false);
      }
      final stream = await _recorder.startStream(
        const RecordConfig(
          encoder: AudioEncoder.pcm16bits,
          numChannels: 1,
          sampleRate: 44100,
          bitRate: 128000,
        ),
      );
      return _AudioSource(stream, true);
    } catch (_) {
      return _AudioSource(_syntheticAudioSource(), false);
    }
  }

  Future<void> _publish(String streamId, Stream<List<int>> audio, Emitter<BroadcasterState> emit) async {
    try {
      // Browsers cannot stream an HTTP request body, so publish over a
      // WebSocket on the web; native platforms keep the chunked HTTP upload.
      if (kIsWeb) {
        await repository.publishWs(
          streamId: streamId,
          audio: audio,
          cancelToken: _cancelToken!,
        );
      } else {
        await repository.publish(
          streamId: streamId,
          audio: audio,
          cancelToken: _cancelToken!,
        );
      }
    } catch (e) {
      if (_cancelToken?.isCancelled == true) return;
      if (!isClosed) emit(BroadcasterError(e.toString()));
    } finally {
      _elapsedTimer?.cancel();
      _elapsedTimer = null;
      if (_usingMic) {
        try {
          await _recorder.stop();
        } catch (_) {}
      }
    }
  }

  Future<void> _onStop(BroadcasterStopRequested event, Emitter<BroadcasterState> emit) async {
    _elapsedTimer?.cancel();
    if (_usingMic) {
      try {
        await _recorder.stop();
      } catch (_) {}
    }
    _cancelToken?.cancel('user_stopped');
    _cancelToken = null;
    if (_stream != null) {
      emit(BroadcasterReady(_stream!));
    } else {
      emit(const BroadcasterIdle());
    }
  }

  @override
  Future<void> close() {
    _elapsedTimer?.cancel();
    _cancelToken?.cancel('bloc_closed');
    _recorder.dispose();
    return super.close();
  }
}

class _AudioSource {
  final Stream<List<int>> source;
  final bool isMic;
  _AudioSource(this.source, this.isMic);
}

/// Fallback synthetic audio stream used when microphone capture is
/// unavailable (simulator, denied permission, web). Generates random
/// bytes at a realistic-ish rate so listeners still see chunks flowing.
Stream<List<int>> _syntheticAudioSource() async* {
  final random = Random();
  while (true) {
    yield List<int>.generate(4096, (_) => random.nextInt(256));
    await Future<void>.delayed(const Duration(milliseconds: 100));
  }
}
