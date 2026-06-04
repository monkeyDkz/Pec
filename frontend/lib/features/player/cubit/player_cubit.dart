import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:just_audio/just_audio.dart';
import 'package:streampulse/core/storage/secure_storage.dart';
import 'package:streampulse/features/streams/models/live_stream.dart';
import 'package:streampulse/features/streams/repositories/stream_repository.dart';

class PlayerState extends Equatable {
  final LiveStream? current;
  final bool isPlaying;
  final double volume;
  final bool buffering;

  const PlayerState({
    this.current,
    this.isPlaying = false,
    this.volume = 1.0,
    this.buffering = false,
  });

  PlayerState copyWith({
    LiveStream? current,
    bool? isPlaying,
    double? volume,
    bool? buffering,
  }) =>
      PlayerState(
        current: current ?? this.current,
        isPlaying: isPlaying ?? this.isPlaying,
        volume: volume ?? this.volume,
        buffering: buffering ?? this.buffering,
      );

  @override
  List<Object?> get props => [current?.id, isPlaying, volume, buffering];
}

/// PlayerCubit wraps a just_audio [AudioPlayer] to play live audio streams.
/// The JWT is attached as an Authorization header on the audio source.
class PlayerCubit extends Cubit<PlayerState> {
  final StreamRepository streamRepository;
  final SecureStorage storage;
  final AudioPlayer player = AudioPlayer();

  PlayerCubit({required this.streamRepository, required this.storage})
      : super(const PlayerState()) {
    player.playingStream.listen((playing) {
      emit(state.copyWith(isPlaying: playing));
    });
    player.processingStateStream.listen((ps) {
      emit(state.copyWith(
        buffering: ps == ProcessingState.loading || ps == ProcessingState.buffering,
      ));
    });
  }

  /// Expose the underlying player so UI can subscribe to position streams.
  AudioPlayer get audioPlayer => player;

  Future<void> playStream(LiveStream stream) async {
    final token = await storage.getToken();
    emit(state.copyWith(current: stream, buffering: true));
    await player.setAudioSource(
      AudioSource.uri(
        Uri.parse(streamRepository.listenUrl(stream.id)),
        headers: token != null ? {'Authorization': 'Bearer $token'} : null,
      ),
    );
    await player.play();
  }

  Future<void> resume() => player.play();
  Future<void> pause() => player.pause();

  Future<void> stop() async {
    await player.stop();
    emit(const PlayerState());
  }

  Future<void> setVolume(double volume) async {
    await player.setVolume(volume);
    emit(state.copyWith(volume: volume));
  }

  @override
  Future<void> close() {
    player.dispose();
    return super.close();
  }
}
