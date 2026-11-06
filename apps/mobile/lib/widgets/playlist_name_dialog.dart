import 'package:flutter/material.dart';

import '../constants/auth_config.dart';
import '../constants/theme.dart';
import '../i18n/dictionary.dart';

bool isValidPlaylistName(String name) {
  final length = name.trim().runes.length;
  return length >= 1 && length <= kPlaylistNameMaxLength;
}

Future<String?> showPlaylistNameDialog(
  BuildContext context, {
  required String title,
  required String confirmLabel,
  String initial = '',
}) {
  return showDialog<String>(
    context: context,
    builder: (context) => _PlaylistNameDialog(
      title: title,
      confirmLabel: confirmLabel,
      initial: initial,
    ),
  );
}

class _PlaylistNameDialog extends StatefulWidget {
  final String title;
  final String confirmLabel;
  final String initial;

  const _PlaylistNameDialog({
    required this.title,
    required this.confirmLabel,
    required this.initial,
  });

  @override
  State<_PlaylistNameDialog> createState() => _PlaylistNameDialogState();
}

class _PlaylistNameDialogState extends State<_PlaylistNameDialog> {
  late final _controller = TextEditingController(text: widget.initial);
  bool _touched = false;

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _submit() {
    setState(() => _touched = true);
    if (!isValidPlaylistName(_controller.text)) return;
    Navigator.of(context).pop(_controller.text.trim());
  }

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    final invalid = _touched && !isValidPlaylistName(_controller.text);
    return AlertDialog(
      backgroundColor: AppColors.surfaceRaised,
      title: Text(widget.title, style: TextStyle(color: AppColors.textPrimary)),
      content: TextField(
        controller: _controller,
        autofocus: true,
        textInputAction: TextInputAction.done,
        onSubmitted: (_) => _submit(),
        onChanged: (_) {
          if (_touched) setState(() {});
        },
        style: TextStyle(color: AppColors.textPrimary),
        decoration: InputDecoration(
          labelText: t.playlistNameLabel,
          errorText: invalid
              ? t.playlistNameInvalid(kPlaylistNameMaxLength)
              : null,
          counterText: '',
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: Text(
            t.cancel,
            style: TextStyle(color: AppColors.textSecondary),
          ),
        ),
        TextButton(
          onPressed: _submit,
          child: Text(
            widget.confirmLabel,
            style: const TextStyle(
              color: AppColors.accent,
              fontWeight: FontWeight.w700,
            ),
          ),
        ),
      ],
    );
  }
}

Future<bool> showConfirmDialog(
  BuildContext context, {
  required String title,
  required String body,
  required String confirmLabel,
}) async {
  final t = currentStrings;
  final ok = await showDialog<bool>(
    context: context,
    builder: (context) => AlertDialog(
      backgroundColor: AppColors.surfaceRaised,
      title: Text(title, style: TextStyle(color: AppColors.textPrimary)),
      content: Text(body, style: TextStyle(color: AppColors.textSecondary)),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(false),
          child: Text(
            t.cancel,
            style: TextStyle(color: AppColors.textSecondary),
          ),
        ),
        TextButton(
          onPressed: () => Navigator.of(context).pop(true),
          child: Text(
            confirmLabel,
            style: TextStyle(
              color: AppColors.errorText,
              fontWeight: FontWeight.w700,
            ),
          ),
        ),
      ],
    ),
  );
  return ok ?? false;
}
