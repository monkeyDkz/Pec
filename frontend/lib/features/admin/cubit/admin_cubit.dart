import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/features/auth/models/user.dart';
import 'package:streampulse/features/admin/repositories/admin_repository.dart';

class AdminState extends Equatable {
  final bool loading;
  final AdminStats? stats;
  final List<User> users;
  final String? error;

  const AdminState({this.loading = false, this.stats, this.users = const [], this.error});

  @override
  List<Object?> get props => [loading, stats?.totalUsers, users, error];
}

class AdminCubit extends Cubit<AdminState> {
  final AdminRepository repository;

  AdminCubit(this.repository) : super(const AdminState(loading: true));

  Future<void> load() async {
    emit(const AdminState(loading: true));
    try {
      final stats = await repository.stats();
      final users = await repository.listUsers();
      emit(AdminState(stats: stats, users: users));
    } catch (e) {
      emit(const AdminState(error: 'Failed to load admin data'));
    }
  }

  Future<void> updateRole(String userId, String role) async {
    try {
      await repository.updateRole(userId, role);
      await load();
    } catch (_) {
      emit(AdminState(stats: state.stats, users: state.users, error: 'Failed to update role'));
    }
  }
}
