import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

import '../constants/app_links.dart';
import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import '../services/api_client.dart';
import '../services/auth_service.dart';

class LoginPrompt extends StatelessWidget {
  final String title;
  final String description;
  final String buttonLabel;
  final bool showPrivacyNote;
  final IconData icon;
  final VoidCallback? onSignedIn;

  const LoginPrompt({
    super.key,
    required this.title,
    required this.description,
    required this.buttonLabel,
    this.showPrivacyNote = true,
    this.icon = Icons.music_note_rounded,
    this.onSignedIn,
  });

  Future<void> _signIn(BuildContext context) async {
    final messenger = ScaffoldMessenger.maybeOf(context);
    try {
      final ok = await AuthService.instance.signIn();
      if (ok) onSignedIn?.call();
    } catch (e) {
      final t = currentStrings;
      messenger?.showSnackBar(
        SnackBar(
          content: Text(
            e is AuthUnavailableException ? t.loginUnavailable : t.signInFailed,
          ),
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;
    return ListenableBuilder(
      listenable: AuthService.instance,
      builder: (context, _) {
        final busy = AuthService.instance.busy;
        return Padding(
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 56,
                height: 56,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: AppColors.accent.withValues(alpha: 0.15),
                ),
                child: Icon(icon, color: AppColors.accent, size: 28),
              ),
              const SizedBox(height: 14),
              Text(
                title,
                textAlign: TextAlign.center,
                style: displayFont(fontSize: 18),
              ),
              const SizedBox(height: 6),
              Text(
                description,
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 13.5,
                  height: 1.5,
                  color: AppColors.textSecondary,
                ),
              ),
              const SizedBox(height: 16),
              FilledButton.icon(
                onPressed: busy ? null : () => _signIn(context),
                style: FilledButton.styleFrom(
                  backgroundColor: AppColors.accent,
                  foregroundColor: AppColors.accentInk,
                  padding: const EdgeInsets.symmetric(
                    horizontal: 20,
                    vertical: 12,
                  ),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(999),
                  ),
                ),
                icon: busy
                    ? const SizedBox(
                        width: 16,
                        height: 16,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          color: AppColors.accentInk,
                        ),
                      )
                    : const Icon(Icons.login, size: 18),
                label: Text(
                  buttonLabel,
                  style: const TextStyle(fontWeight: FontWeight.w700),
                ),
              ),
              if (showPrivacyNote) ...[
                const SizedBox(height: 12),
                Text(
                  t.loginPrivacyNote,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 11.5,
                    color: AppColors.textTertiary,
                  ),
                ),
                TextButton(
                  onPressed: () => launchUrl(
                    Uri.parse(kPrivacyPolicyUrl),
                    mode: LaunchMode.externalApplication,
                  ),
                  child: Text(
                    t.privacyPolicy,
                    style: const TextStyle(
                      fontSize: 12,
                      color: AppColors.accent,
                    ),
                  ),
                ),
              ],
            ],
          ),
        );
      },
    );
  }
}
