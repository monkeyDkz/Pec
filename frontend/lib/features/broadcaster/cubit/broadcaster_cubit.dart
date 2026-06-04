import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/features/streams/models/live_stream.dart';
import 'package:streampulse/features/streams/repositories/stream_repository.dart';

class BroadcasterState extends Equatable {
  final bool loading;
  final List<LiveStream> mine;
  final String? error;

  const BroadcasterState({this.loading = false, this.mine = const [], this.error});

  @override
  List<Object?> get props => [loading, mine, error];
}

/// Manages the broadcaster dashboard: start/stop the diffuser's own live streams.
class BroadcasterCubit extends Cubit<BroadcasterState> {
  final StreamRepository repository;
  final String currentUserId;

  BroadcasterCubit({required this.repository, required this.currentUserId})
      : super(const BroadcasterState(loading: true));

  Future<void> load() async {
    emit(const BroadcasterState(loading: true));
    try {
      final all = await repository.listLive();
      final mine = all.where((s) => s.broadcasterId == currentUserId).toList();
      emit(BroadcasterState(mine: mine));
    } catch (e) {
      emit(const BroadcasterState(error: 'Failed to load your streams'));
    }
  }

  Future<void> create(String title, String description) async {
    try {
      await repository.create(title: title, description: description);
      await load();
    } catch (_) {
      emit(BroadcasterState(mine: state.mine, error: 'Failed to start stream'));
    }
  }

  Future<void> stop(String id) async {
    try {
      await repository.stop(id);
      await load();
    } catch (_) {
      emit(BroadcasterState(mine: state.mine, error: 'Failed to stop stream'));
    }
  }
}
