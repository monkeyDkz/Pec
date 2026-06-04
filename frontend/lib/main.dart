import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/core/router/app_router.dart';
import 'package:streampulse/core/storage/secure_storage.dart';
import 'package:streampulse/core/theme/app_theme.dart';
import 'package:streampulse/features/admin/repositories/admin_repository.dart';
import 'package:streampulse/features/auth/bloc/auth_bloc.dart';
import 'package:streampulse/features/auth/repositories/auth_repository.dart';
import 'package:streampulse/features/player/cubit/player_cubit.dart';
import 'package:streampulse/features/playlists/repositories/playlist_repository.dart';
import 'package:streampulse/features/playlists/repositories/track_repository.dart';
import 'package:streampulse/features/profile/cubit/current_user_cubit.dart';
import 'package:streampulse/features/profile/repositories/profile_repository.dart';
import 'package:streampulse/features/streams/repositories/stream_repository.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(StreamPulseApp(apiClient: ApiClient(), storage: SecureStorage()));
}

class StreamPulseApp extends StatelessWidget {
  final ApiClient apiClient;
  final SecureStorage storage;

  const StreamPulseApp({super.key, required this.apiClient, required this.storage});

  @override
  Widget build(BuildContext context) {
    return MultiRepositoryProvider(
      providers: [
        RepositoryProvider.value(value: apiClient),
        RepositoryProvider(create: (_) => AuthRepository(apiClient: apiClient)),
        RepositoryProvider(create: (_) => ProfileRepository(apiClient: apiClient)),
        RepositoryProvider(create: (_) => StreamRepository(apiClient: apiClient)),
        RepositoryProvider(create: (_) => PlaylistRepository(apiClient: apiClient)),
        RepositoryProvider(create: (_) => TrackRepository(apiClient: apiClient)),
        RepositoryProvider(create: (_) => AdminRepository(apiClient: apiClient)),
      ],
      child: _AppRoot(storage: storage),
    );
  }
}

class _AppRoot extends StatefulWidget {
  final SecureStorage storage;
  const _AppRoot({required this.storage});

  @override
  State<_AppRoot> createState() => _AppRootState();
}

class _AppRootState extends State<_AppRoot> {
  late final AuthBloc _authBloc;
  late final CurrentUserCubit _currentUserCubit;
  late final PlayerCubit _playerCubit;
  late final GoRouter _router;

  @override
  void initState() {
    super.initState();
    _authBloc = AuthBloc(authRepository: context.read<AuthRepository>())
      ..add(AuthCheckRequested());
    _currentUserCubit = CurrentUserCubit(context.read<ProfileRepository>());
    _playerCubit = PlayerCubit(
      streamRepository: context.read<StreamRepository>(),
      storage: widget.storage,
    );
    _router = AppRouter.create(_authBloc);
  }

  @override
  void dispose() {
    _authBloc.close();
    _currentUserCubit.close();
    _playerCubit.close();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: [
        BlocProvider.value(value: _authBloc),
        BlocProvider.value(value: _currentUserCubit),
        BlocProvider.value(value: _playerCubit),
      ],
      child: MaterialApp.router(
        title: 'StreamPulse',
        theme: AppTheme.light,
        darkTheme: AppTheme.dark,
        themeMode: ThemeMode.system,
        routerConfig: _router,
        debugShowCheckedModeBanner: false,
      ),
    );
  }
}
