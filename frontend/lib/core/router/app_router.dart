import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:go_router/go_router.dart';
import 'package:streampulse/features/auth/bloc/auth_bloc.dart';
import 'package:streampulse/features/auth/screens/login_screen.dart';
import 'package:streampulse/features/auth/screens/register_screen.dart';
import 'package:streampulse/features/home/home_shell.dart';
import 'package:streampulse/features/player/screens/now_playing_screen.dart';
import 'package:streampulse/features/playlists/screens/playlist_detail_screen.dart';

class AppRouter {
  static GoRouter create(AuthBloc authBloc) {
    return GoRouter(
      initialLocation: '/login',
      refreshListenable: GoRouterRefreshStream(authBloc.stream),
      redirect: (context, state) {
        final authState = authBloc.state;
        // Until the initial auth check resolves, stay put.
        if (authState is AuthInitial || authState is AuthLoading) return null;

        final loggedIn = authState is AuthAuthenticated;
        final loc = state.matchedLocation;
        final onAuthPage = loc == '/login' || loc == '/register';

        if (!loggedIn && !onAuthPage) return '/login';
        if (loggedIn && onAuthPage) return '/home';
        return null;
      },
      routes: [
        GoRoute(path: '/login', builder: (c, s) => const LoginScreen()),
        GoRoute(path: '/register', builder: (c, s) => const RegisterScreen()),
        GoRoute(path: '/home', builder: (c, s) => const HomeShell()),
        GoRoute(path: '/player', builder: (c, s) => const NowPlayingScreen()),
        GoRoute(
          path: '/playlist/:id',
          builder: (c, s) => PlaylistDetailScreen(id: s.pathParameters['id']!),
        ),
      ],
    );
  }
}

/// Bridges a BLoC [Stream] to a [Listenable] so GoRouter re-evaluates redirects
/// whenever the auth state changes.
class GoRouterRefreshStream extends ChangeNotifier {
  GoRouterRefreshStream(Stream<dynamic> stream) {
    notifyListeners();
    _subscription = stream.asBroadcastStream().listen((_) => notifyListeners());
  }

  late final StreamSubscription<dynamic> _subscription;

  @override
  void dispose() {
    _subscription.cancel();
    super.dispose();
  }
}
