import 'package:go_router/go_router.dart';
import 'package:streampulse/features/auth/screens/login_screen.dart';
import 'package:streampulse/features/streams/screens/streams_screen.dart';

class AppRouter {
  static final router = GoRouter(
    initialLocation: '/login',
    routes: [
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/streams',
        builder: (context, state) => const StreamsScreen(),
      ),
    ],
  );
}
