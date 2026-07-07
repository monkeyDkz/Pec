import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/features/playlists/models/playlist_model.dart';
import 'package:streampulse/features/playlists/repositories/playlist_repository.dart';

abstract class PlaylistsEvent extends Equatable {
  const PlaylistsEvent();
  @override
  List<Object?> get props => [];
}

class PlaylistsLoadRequested extends PlaylistsEvent {
  const PlaylistsLoadRequested();
}

class PlaylistCreateRequested extends PlaylistsEvent {
  final String name;
  final String description;
  const PlaylistCreateRequested({required this.name, required this.description});
  @override
  List<Object?> get props => [name, description];
}

class PlaylistDeleteRequested extends PlaylistsEvent {
  final String id;
  const PlaylistDeleteRequested(this.id);
  @override
  List<Object?> get props => [id];
}

abstract class PlaylistsState extends Equatable {
  const PlaylistsState();
  @override
  List<Object?> get props => [];
}

class PlaylistsInitial extends PlaylistsState {
  const PlaylistsInitial();
}

class PlaylistsLoading extends PlaylistsState {
  const PlaylistsLoading();
}

class PlaylistsLoaded extends PlaylistsState {
  final List<PlaylistModel> playlists;
  const PlaylistsLoaded(this.playlists);
  @override
  List<Object?> get props => [playlists];
}

class PlaylistsError extends PlaylistsState {
  final String message;
  const PlaylistsError(this.message);
  @override
  List<Object?> get props => [message];
}

class PlaylistsBloc extends Bloc<PlaylistsEvent, PlaylistsState> {
  final PlaylistRepository repository;

  PlaylistsBloc({required this.repository}) : super(const PlaylistsInitial()) {
    on<PlaylistsLoadRequested>(_load);
    on<PlaylistCreateRequested>(_create);
    on<PlaylistDeleteRequested>(_delete);
  }

  Future<void> _load(PlaylistsEvent event, Emitter<PlaylistsState> emit) async {
    if (state is! PlaylistsLoaded) emit(const PlaylistsLoading());
    try {
      final playlists = await repository.listMine();
      emit(PlaylistsLoaded(playlists));
    } catch (e) {
      emit(PlaylistsError(e.toString()));
    }
  }

  Future<void> _create(PlaylistCreateRequested event, Emitter<PlaylistsState> emit) async {
    try {
      await repository.create(name: event.name, description: event.description);
      final updated = await repository.listMine();
      emit(PlaylistsLoaded(updated));
    } catch (e) {
      emit(PlaylistsError(e.toString()));
    }
  }

  Future<void> _delete(PlaylistDeleteRequested event, Emitter<PlaylistsState> emit) async {
    try {
      await repository.delete(event.id);
      final updated = await repository.listMine();
      emit(PlaylistsLoaded(updated));
    } catch (e) {
      emit(PlaylistsError(e.toString()));
    }
  }
}
