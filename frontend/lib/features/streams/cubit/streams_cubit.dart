import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/features/streams/models/live_stream.dart';
import 'package:streampulse/features/streams/repositories/stream_repository.dart';

class StreamsState extends Equatable {
  final bool loading;
  final List<LiveStream> streams;
  final String? error;

  const StreamsState({this.loading = false, this.streams = const [], this.error});

  StreamsState copyWith({bool? loading, List<LiveStream>? streams, String? error}) =>
      StreamsState(
        loading: loading ?? this.loading,
        streams: streams ?? this.streams,
        error: error,
      );

  @override
  List<Object?> get props => [loading, streams, error];
}

class StreamsCubit extends Cubit<StreamsState> {
  final StreamRepository repository;

  StreamsCubit(this.repository) : super(const StreamsState(loading: true));

  Future<void> load() async {
    emit(state.copyWith(loading: true, error: null));
    try {
      final streams = await repository.listLive();
      emit(StreamsState(loading: false, streams: streams));
    } catch (e) {
      emit(StreamsState(loading: false, error: 'Failed to load streams'));
    }
  }
}
