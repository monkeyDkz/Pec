import 'package:flutter/material.dart';

class AppTheme {
  static const _primaryColor = Color(0xFF6C5CE7);

  static ThemeData get light => ThemeData(
        useMaterial3: true,
        colorSchemeSeed: _primaryColor,
        brightness: Brightness.light,
      );

  static ThemeData get dark => ThemeData(
        useMaterial3: true,
        colorSchemeSeed: _primaryColor,
        brightness: Brightness.dark,
      );
}
