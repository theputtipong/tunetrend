import 'dart:async';

import 'package:flutter/material.dart';
import '../constants/theme.dart';
import '../constants/validation.dart';
import '../i18n/dictionary.dart';
import '../services/api_client.dart';

enum _ContactMethod { email, phone }

const _autoCloseSeconds = 3;

class ContactForm extends StatefulWidget {
  final VoidCallback? onClose;

  const ContactForm({super.key, this.onClose});

  @override
  State<ContactForm> createState() => _ContactFormState();
}

class _ContactFormState extends State<ContactForm> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _messageController = TextEditingController();
  final _contactController = TextEditingController();
  final _api = ApiClient();

  _ContactMethod _method = _ContactMethod.email;
  bool _submitting = false;
  String? _successMessage;
  String? _errorMessage;
  int _secondsLeft = _autoCloseSeconds;
  Timer? _closeTimer;

  @override
  void dispose() {
    _closeTimer?.cancel();
    _nameController.dispose();
    _messageController.dispose();
    _contactController.dispose();
    super.dispose();
  }

  void _startCloseCountdown() {
    _secondsLeft = _autoCloseSeconds;
    _closeTimer?.cancel();
    _closeTimer = Timer.periodic(const Duration(seconds: 1), (timer) {
      if (!mounted) {
        timer.cancel();
        return;
      }
      setState(() => _secondsLeft--);
      if (_secondsLeft <= 0) {
        timer.cancel();
        widget.onClose?.call();
      }
    });
  }

  Future<void> _submit() async {
    final t = currentStrings;
    setState(() {
      _successMessage = null;
      _errorMessage = null;
    });

    if (!_formKey.currentState!.validate()) return;

    setState(() => _submitting = true);
    try {
      await _api.submitContact(
        name: _nameController.text.trim(),
        message: _messageController.text.trim(),
        contactEmail: _method == _ContactMethod.email ? _contactController.text.trim() : null,
        contactPhone: _method == _ContactMethod.phone ? _contactController.text.trim() : null,
      );
      setState(() {
        _successMessage = t.contactSuccessMessage;
        _nameController.clear();
        _messageController.clear();
        _contactController.clear();
      });
      _startCloseCountdown();
    } on ApiException catch (e) {
      setState(() {
        _errorMessage = e.statusCode == 429 ? t.contactErrorRateLimited : t.contactErrorGeneric;
      });
    } finally {
      if (mounted) setState(() => _submitting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;

    if (_successMessage != null) {
      return Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 56,
            height: 56,
            decoration: BoxDecoration(shape: BoxShape.circle, color: AppColors.accent),
            child: Icon(Icons.check, color: AppColors.accentInk, size: 28),
          ),
          const SizedBox(height: 12),
          Text(
            _successMessage!,
            textAlign: TextAlign.center,
            style: TextStyle(color: AppColors.textPrimary, fontWeight: FontWeight.w600, fontSize: 15),
          ),
          const SizedBox(height: 8),
          Text(
            t.closingIn(_secondsLeft),
            style: TextStyle(color: AppColors.textTertiary, fontSize: 12.5),
          ),
        ],
      );
    }

    return Form(
      key: _formKey,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(t.contactNameLabel, style: TextStyle(fontSize: 12.5, color: AppColors.textSecondary)),
          const SizedBox(height: 6),
          TextFormField(
            controller: _nameController,
            decoration: InputDecoration(hintText: t.contactNamePlaceholder),
          ),
          const SizedBox(height: 14),
          Text(t.contactMessageLabel, style: TextStyle(fontSize: 12.5, color: AppColors.textSecondary)),
          const SizedBox(height: 6),
          TextFormField(
            controller: _messageController,
            maxLines: 4,
            decoration: InputDecoration(hintText: t.contactMessagePlaceholder),
            validator: (value) {
              final v = value?.trim() ?? '';
              if (v.isEmpty) return t.contactErrorMessageRequired;
              if (v.length > maxMessageLen) return t.contactErrorTooLong;
              return null;
            },
          ),
          const SizedBox(height: 14),
          Text(t.contactMethodLabel, style: TextStyle(fontSize: 12.5, color: AppColors.textSecondary)),
          const SizedBox(height: 6),
          SegmentedButton<_ContactMethod>(
            segments: [
              ButtonSegment(value: _ContactMethod.email, label: Text(t.contactMethodEmail)),
              ButtonSegment(value: _ContactMethod.phone, label: Text(t.contactMethodPhone)),
            ],
            selected: {_method},
            onSelectionChanged: (selection) {
              setState(() {
                _method = selection.first;
                _contactController.clear();
              });
            },
          ),
          const SizedBox(height: 10),
          TextFormField(
            controller: _contactController,
            keyboardType: _method == _ContactMethod.email
                ? TextInputType.emailAddress
                : TextInputType.phone,
            decoration: InputDecoration(
              hintText: _method == _ContactMethod.email
                  ? t.contactEmailPlaceholder
                  : t.contactPhonePlaceholder,
            ),
            validator: (value) {
              final v = value?.trim() ?? '';
              if (_method == _ContactMethod.email && !isValidEmail(v)) {
                return t.contactErrorInvalidEmail;
              }
              if (_method == _ContactMethod.phone && !isValidThaiPhone(v)) {
                return t.contactErrorInvalidPhone;
              }
              return null;
            },
          ),
          const SizedBox(height: 16),
          ElevatedButton(
            onPressed: _submitting ? null : _submit,
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.accent,
              foregroundColor: AppColors.accentInk,
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(999)),
              padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
            ),
            child: _submitting
                ? SizedBox(
                    width: 16,
                    height: 16,
                    child: CircularProgressIndicator(strokeWidth: 2, color: AppColors.accentInk),
                  )
                : Text(t.contactSubmit, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 14)),
          ),
          if (_errorMessage != null) ...[
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: AppColors.errorBg,
                borderRadius: BorderRadius.circular(10),
              ),
              child: Text(_errorMessage!, style: TextStyle(color: AppColors.errorText, fontSize: 13)),
            ),
          ],
        ],
      ),
    );
  }
}
