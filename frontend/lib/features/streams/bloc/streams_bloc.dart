import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/features/streams/models/stream_model.dart';
import 'package:streampulse/features/streams/repositories/stream_repository.dart';

// ── Events ────────────────────────────────────────────────────────────
abstract class StreamsEvent extends Equatable {
  const StreamsEvent();
  @override
  List<Object?> get props => [];
}

class StreamsLoadRequested extends StreamsEvent {
  const StreamsLoadRequested();
}

class StreamsRefreshRequested extends StreamsEvent {
  const StreamsRefreshRequested();
}

// ── States ────────────────────────────────────────────────────────────
abstract class StreamsState extends Equatable {
  const StreamsState();
  @override
  List<Object?> get props => [];
}

class StreamsInitial extends StreamsState {
  const StreamsInitial();
}

class StreamsLoading extends StreamsState {
  const StreamsLoading();
}

class StreamsLoaded extends StreamsState {
  final List<StreamModel> streams;
  const StreamsLoaded(this.streams);
  @override
  List<Object?> get props => [streams];
}

class StreamsError extends StreamsState {
  final String message;
  const StreamsError(this.message);
  @override
  List<Object?> get props => [message];
}

// ── BLoC ──────────────────────────────────────────────────────────────
class StreamsBloc extends Bloc<StreamsEvent, StreamsState> {
  final StreamRepository repository;

  StreamsBloc({required this.repository}) : super(const StreamsInitial()) {
    on<StreamsLoadRequested>(_load);
    on<StreamsRefreshRequested>(_load);
  }

  Future<void> _load(StreamsEvent event, Emitter<StreamsState> emit) async {
    if (state is! StreamsLoaded) {
      emit(const StreamsLoading());
    }
    try {
      final streams = await repository.listLive();
      emit(StreamsLoaded(streams));
    } catch (e) {
      emit(StreamsError(e.toString()));
    }
  }
}
