import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:just_audio/just_audio.dart';

abstract class PlayerEvent extends Equatable {
  const PlayerEvent();
  @override
  List<Object?> get props => [];
}

class PlayerLoadRequested extends PlayerEvent {
  final String url;
  final String title;
  const PlayerLoadRequested({required this.url, required this.title});
  @override
  List<Object?> get props => [url, title];
}

class PlayerPlayRequested extends PlayerEvent {
  const PlayerPlayRequested();
}

class PlayerPauseRequested extends PlayerEvent {
  const PlayerPauseRequested();
}

class PlayerStopRequested extends PlayerEvent {
  const PlayerStopRequested();
}

class PlayerVolumeChanged extends PlayerEvent {
  final double volume;
  const PlayerVolumeChanged(this.volume);
  @override
  List<Object?> get props => [volume];
}

abstract class PlayerState extends Equatable {
  const PlayerState();
  @override
  List<Object?> get props => [];
}

class PlayerIdle extends PlayerState {
  const PlayerIdle();
}

class PlayerLoading extends PlayerState {
  final String title;
  const PlayerLoading(this.title);
  @override
  List<Object?> get props => [title];
}

class PlayerPlaying extends PlayerState {
  final String title;
  final double volume;
  const PlayerPlaying({required this.title, required this.volume});
  @override
  List<Object?> get props => [title, volume];
}

class PlayerPaused extends PlayerState {
  final String title;
  final double volume;
  const PlayerPaused({required this.title, required this.volume});
  @override
  List<Object?> get props => [title, volume];
}

class PlayerError extends PlayerState {
  final String message;
  const PlayerError(this.message);
  @override
  List<Object?> get props => [message];
}

class PlayerBloc extends Bloc<PlayerEvent, PlayerState> {
  final AudioPlayer _player;
  String _currentTitle = '';

  PlayerBloc({AudioPlayer? player})
      : _player = player ?? AudioPlayer(),
        super(const PlayerIdle()) {
    on<PlayerLoadRequested>(_onLoad);
    on<PlayerPlayRequested>(_onPlay);
    on<PlayerPauseRequested>(_onPause);
    on<PlayerStopRequested>(_onStop);
    on<PlayerVolumeChanged>(_onVolume);
  }

  Future<void> _onLoad(PlayerLoadRequested event, Emitter<PlayerState> emit) async {
    emit(PlayerLoading(event.title));
    _currentTitle = event.title;
    try {
      await _player.setUrl(event.url);
      await _player.play();
      emit(PlayerPlaying(title: event.title, volume: _player.volume));
    } catch (e) {
      emit(PlayerError(e.toString()));
    }
  }

  Future<void> _onPlay(PlayerPlayRequested event, Emitter<PlayerState> emit) async {
    await _player.play();
    emit(PlayerPlaying(title: _currentTitle, volume: _player.volume));
  }

  Future<void> _onPause(PlayerPauseRequested event, Emitter<PlayerState> emit) async {
    await _player.pause();
    emit(PlayerPaused(title: _currentTitle, volume: _player.volume));
  }

  Future<void> _onStop(PlayerStopRequested event, Emitter<PlayerState> emit) async {
    await _player.stop();
    emit(const PlayerIdle());
  }

  Future<void> _onVolume(PlayerVolumeChanged event, Emitter<PlayerState> emit) async {
    await _player.setVolume(event.volume.clamp(0.0, 1.0));
    if (state is PlayerPlaying) {
      emit(PlayerPlaying(title: _currentTitle, volume: _player.volume));
    } else if (state is PlayerPaused) {
      emit(PlayerPaused(title: _currentTitle, volume: _player.volume));
    }
  }

  @override
  Future<void> close() async {
    await _player.dispose();
    return super.close();
  }
}
