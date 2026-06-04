import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/features/auth/models/user.dart';
import 'package:streampulse/features/profile/repositories/profile_repository.dart';

class CurrentUserState extends Equatable {
  final bool loading;
  final User? user;
  final String? error;

  const CurrentUserState({this.loading = false, this.user, this.error});

  @override
  List<Object?> get props => [loading, user?.id, user?.role, error];
}

class CurrentUserCubit extends Cubit<CurrentUserState> {
  final ProfileRepository repository;

  CurrentUserCubit(this.repository) : super(const CurrentUserState());

  Future<void> load() async {
    emit(const CurrentUserState(loading: true));
    try {
      final user = await repository.me();
      emit(CurrentUserState(user: user));
    } catch (e) {
      emit(const CurrentUserState(error: 'Failed to load profile'));
    }
  }

  Future<void> updateUsername(String username) async {
    try {
      final user = await repository.updateUsername(username);
      emit(CurrentUserState(user: user));
    } catch (_) {
      // keep current state on failure
    }
  }

  void clear() => emit(const CurrentUserState());
}
