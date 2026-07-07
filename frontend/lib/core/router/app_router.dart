import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/features/auth/bloc/auth_bloc.dart';
import 'package:streampulse/features/auth/screens/login_screen.dart';
import 'package:streampulse/features/auth/screens/register_screen.dart';
import 'package:streampulse/features/broadcaster/screens/broadcaster_screen.dart';
import 'package:streampulse/features/player/screens/player_screen.dart';
import 'package:streampulse/features/playlists/screens/playlists_screen.dart';
import 'package:streampulse/features/streams/bloc/streams_bloc.dart';
import 'package:streampulse/features/streams/repositories/stream_repository.dart';
import 'package:streampulse/features/streams/screens/streams_screen.dart';

class AppRouter {
  static GoRouter create({required ApiClient apiClient}) {
    final streamRepo = StreamRepository(apiClient: apiClient);

    return GoRouter(
      initialLocation: '/login',
      redirect: (context, state) {
        final auth = context.read<AuthBloc>().state;
        final loggingIn = state.matchedLocation == '/login' ||
            state.matchedLocation == '/register';

        if (auth is AuthUnauthenticated && !loggingIn) return '/login';
        if (auth is AuthAuthenticated && loggingIn) return '/streams';
        return null;
      },
      routes: [
        GoRoute(
          path: '/login',
          builder: (_, __) => const LoginScreen(),
        ),
        GoRoute(
          path: '/register',
          builder: (_, __) => const RegisterScreen(),
        ),
        ShellRoute(
          builder: (context, state, child) => _ShellScaffold(child: child),
          routes: [
            GoRoute(
              path: '/streams',
              builder: (_, __) => BlocProvider(
                create: (_) => StreamsBloc(repository: streamRepo)
                  ..add(const StreamsLoadRequested()),
                child: const StreamsScreen(),
              ),
            ),
            GoRoute(
              path: '/playlists',
              builder: (_, __) => PlaylistsScreen(apiClient: apiClient),
            ),
            GoRoute(
              path: '/broadcaster',
              builder: (_, __) => BroadcasterScreen(apiClient: apiClient),
            ),
          ],
        ),
        GoRoute(
          path: '/player/:id',
          builder: (_, state) => PlayerScreen(
            streamId: state.pathParameters['id']!,
            repository: streamRepo,
          ),
        ),
      ],
      errorBuilder: (_, state) => Scaffold(
        body: Center(child: Text('Route inconnue : ${state.uri}')),
      ),
    );
  }
}

class _ShellScaffold extends StatelessWidget {
  final Widget child;
  const _ShellScaffold({required this.child});

  static const _tabs = [
    ('/streams', Icons.radio, 'Streams'),
    ('/playlists', Icons.playlist_play, 'Playlists'),
    ('/broadcaster', Icons.mic, 'Diffuser'),
  ];

  int _indexFor(BuildContext context) {
    final location = GoRouterState.of(context).matchedLocation;
    for (var i = 0; i < _tabs.length; i++) {
      if (location.startsWith(_tabs[i].$1)) return i;
    }
    return 0;
  }

  @override
  Widget build(BuildContext context) {
    final index = _indexFor(context);
    return Scaffold(
      body: child,
      bottomNavigationBar: NavigationBar(
        selectedIndex: index,
        onDestinationSelected: (i) => context.go(_tabs[i].$1),
        destinations: [
          for (final t in _tabs)
            NavigationDestination(
              icon: Icon(t.$2),
              label: t.$3,
            ),
        ],
      ),
    );
  }
}
