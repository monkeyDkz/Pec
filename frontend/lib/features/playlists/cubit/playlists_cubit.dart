import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/features/playlists/models/playlist.dart';
import 'package:streampulse/features/playlists/repositories/playlist_repository.dart';

class PlaylistsState extends Equatable {
  final bool loading;
  final List<Playlist> playlists;
  final String? error;

  const PlaylistsState({this.loading = false, this.playlists = const [], this.error});

  @override
  List<Object?> get props => [loading, playlists, error];
}

class PlaylistsCubit extends Cubit<PlaylistsState> {
  final PlaylistRepository repository;

  PlaylistsCubit(this.repository) : super(const PlaylistsState(loading: true));

  Future<void> load() async {
    emit(const PlaylistsState(loading: true));
    try {
      final playlists = await repository.list();
      emit(PlaylistsState(playlists: playlists));
    } catch (e) {
      emit(const PlaylistsState(error: 'Failed to load playlists'));
    }
  }

  Future<void> create(String name, String description) async {
    try {
      await repository.create(name: name, description: description);
      await load();
    } catch (_) {
      emit(PlaylistsState(playlists: state.playlists, error: 'Failed to create playlist'));
    }
  }

  Future<void> delete(String id) async {
    try {
      await repository.delete(id);
      await load();
    } catch (_) {
      emit(PlaylistsState(playlists: state.playlists, error: 'Failed to delete playlist'));
    }
  }
}
