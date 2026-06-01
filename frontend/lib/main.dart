import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/core/api/api_client.dart';
import 'package:streampulse/core/router/app_router.dart';
import 'package:streampulse/core/theme/app_theme.dart';
import 'package:streampulse/features/auth/bloc/auth_bloc.dart';
import 'package:streampulse/features/auth/repositories/auth_repository.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  final apiClient = ApiClient();
  final authRepository = AuthRepository(apiClient: apiClient);

  runApp(StreamPulseApp(authRepository: authRepository));
}

class StreamPulseApp extends StatelessWidget {
  final AuthRepository authRepository;

  const StreamPulseApp({super.key, required this.authRepository});

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: [
        BlocProvider(
          create: (_) => AuthBloc(authRepository: authRepository)..add(AuthCheckRequested()),
        ),
      ],
      child: MaterialApp.router(
        title: 'StreamPulse',
        theme: AppTheme.light,
        darkTheme: AppTheme.dark,
        themeMode: ThemeMode.system,
        routerConfig: AppRouter.router,
        debugShowCheckedModeBanner: false,
      ),
    );
  }
}
