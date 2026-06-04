import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/features/auth/bloc/auth_bloc.dart';
import 'package:streampulse/features/auth/models/user.dart';
import 'package:streampulse/features/player/cubit/player_cubit.dart';
import 'package:streampulse/features/profile/cubit/current_user_cubit.dart';

class ProfileScreen extends StatelessWidget {
  final User user;
  const ProfileScreen({super.key, required this.user});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Profile')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Center(
            child: CircleAvatar(
              radius: 40,
              child: Text(
                user.username.isNotEmpty ? user.username[0].toUpperCase() : '?',
                style: theme.textTheme.headlineMedium,
              ),
            ),
          ),
          const SizedBox(height: 16),
          Center(child: Text(user.username, style: theme.textTheme.titleLarge)),
          Center(child: Text(user.email, style: theme.textTheme.bodyMedium)),
          const SizedBox(height: 8),
          Center(child: Chip(label: Text(user.role.toUpperCase()))),
          const SizedBox(height: 32),
          FilledButton.tonalIcon(
            icon: const Icon(Icons.logout),
            label: const Text('Log out'),
            onPressed: () {
              context.read<PlayerCubit>().stop();
              context.read<CurrentUserCubit>().clear();
              context.read<AuthBloc>().add(AuthLogoutRequested());
            },
          ),
        ],
      ),
    );
  }
}
