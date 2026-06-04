import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:streampulse/features/admin/screens/admin_screen.dart';
import 'package:streampulse/features/auth/models/user.dart';
import 'package:streampulse/features/broadcaster/screens/broadcaster_screen.dart';
import 'package:streampulse/features/playlists/screens/playlists_screen.dart';
import 'package:streampulse/features/profile/cubit/current_user_cubit.dart';
import 'package:streampulse/features/profile/screens/profile_screen.dart';
import 'package:streampulse/features/streams/screens/streams_screen.dart';
import 'package:streampulse/features/player/widgets/mini_player.dart';

class HomeShell extends StatefulWidget {
  const HomeShell({super.key});

  @override
  State<HomeShell> createState() => _HomeShellState();
}

class _HomeShellState extends State<HomeShell> {
  int _index = 0;

  @override
  void initState() {
    super.initState();
    context.read<CurrentUserCubit>().load();
  }

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<CurrentUserCubit, CurrentUserState>(
      builder: (context, state) {
        if (state.loading) {
          return const Scaffold(body: Center(child: CircularProgressIndicator()));
        }
        final user = state.user;

        // Build role-gated tabs.
        final pages = <Widget>[const StreamsScreen(), const PlaylistsScreen()];
        final dests = <NavigationDestination>[
          const NavigationDestination(icon: Icon(Icons.radio), label: 'Streams'),
          const NavigationDestination(icon: Icon(Icons.queue_music), label: 'Playlists'),
        ];
        if (user != null && user.isBroadcaster) {
          pages.add(const BroadcasterScreen());
          dests.add(const NavigationDestination(icon: Icon(Icons.podcasts), label: 'Broadcast'));
        }
        if (user != null && user.isAdmin) {
          pages.add(const AdminScreen());
          dests.add(const NavigationDestination(
              icon: Icon(Icons.admin_panel_settings), label: 'Admin'));
        }
        pages.add(ProfileScreen(user: user ?? const User(id: '', email: '', username: '?', role: 'user')));
        dests.add(const NavigationDestination(icon: Icon(Icons.person), label: 'Profile'));

        final safeIndex = _index < pages.length ? _index : 0;

        return Scaffold(
          body: Column(
            children: [
              Expanded(child: IndexedStack(index: safeIndex, children: pages)),
              const MiniPlayer(),
            ],
          ),
          bottomNavigationBar: NavigationBar(
            selectedIndex: safeIndex,
            onDestinationSelected: (i) => setState(() => _index = i),
            destinations: dests,
          ),
        );
      },
    );
  }
}
