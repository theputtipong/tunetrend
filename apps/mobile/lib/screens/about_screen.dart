import 'package:flutter/material.dart';
import '../constants/theme.dart';
import '../i18n/dictionary.dart';
import '../i18n/lang.dart';
import '../i18n/strings.dart';
import '../widgets/buy_me_coffee_button.dart';
import '../widgets/contact_form.dart';
import '../widgets/logo_mark.dart';

enum _AboutMenuAction { language, theme }

class AboutScreen extends StatefulWidget {
  const AboutScreen({super.key});

  @override
  State<AboutScreen> createState() => _AboutScreenState();
}

class _AboutScreenState extends State<AboutScreen> {
  @override
  void initState() {
    super.initState();
    LangController.instance.addListener(_onChanged);
    ThemeController.instance.addListener(_onChanged);
  }

  @override
  void dispose() {
    LangController.instance.removeListener(_onChanged);
    ThemeController.instance.removeListener(_onChanged);
    super.dispose();
  }

  void _onChanged() => setState(() {});

  @override
  Widget build(BuildContext context) {
    final t = currentStrings;

    return Scaffold(
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => showDialog(
          context: context,
          builder: (_) => Dialog(
            backgroundColor: AppColors.surfaceRaised,
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: SingleChildScrollView(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(t.contactOpenButton, style: displayFont(fontSize: 17)),
                        IconButton(
                          onPressed: () => Navigator.of(context).pop(),
                          icon: Icon(Icons.close, color: AppColors.textSecondary, size: 20),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    ContactForm(onClose: () => Navigator.of(context).pop()),
                  ],
                ),
              ),
            ),
          ),
        ),
        backgroundColor: AppColors.accent,
        foregroundColor: AppColors.accentInk,
        icon: const Icon(Icons.chat_bubble_outline),
        label: Text(t.contactOpenButton),
      ),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.fromLTRB(16, 16, 16, 32),
          children: [
            Row(
              children: [
                Expanded(
                  child: GestureDetector(
                    onTap: () => Navigator.of(context).pop(),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        LogoMark(size: 26),
                        const SizedBox(width: 8),
                        Text('TuneTrend', style: displayFont(fontSize: 17)),
                      ],
                    ),
                  ),
                ),
                TextButton(
                  onPressed: () => Navigator.of(context).pop(),
                  style: TextButton.styleFrom(
                    backgroundColor: AppColors.surface,
                    foregroundColor: AppColors.textSecondary,
                    side: BorderSide(color: AppColors.border),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(999)),
                    padding: const EdgeInsets.symmetric(horizontal: 14),
                    minimumSize: const Size(0, 36),
                  ),
                  child: Text(
                    t.backToTrends,
                    style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 13),
                  ),
                ),
                const SizedBox(width: 8),
                PopupMenuButton<_AboutMenuAction>(
                  tooltip: t.menuTooltip,
                  icon: Icon(Icons.more_vert, color: AppColors.textSecondary, size: 20),
                  onSelected: (action) {
                    switch (action) {
                      case _AboutMenuAction.language:
                        LangController.instance.toggle();
                      case _AboutMenuAction.theme:
                        ThemeController.instance.toggle(context);
                    }
                  },
                  itemBuilder: (context) => [
                    PopupMenuItem(
                      value: _AboutMenuAction.language,
                      child: ListTile(
                        leading: Icon(Icons.language, size: 20, color: AppColors.textSecondary),
                        title: Text(
                          LangController.instance.resolve() == AppLang.en ? 'ภาษาไทย' : 'English',
                          style: TextStyle(color: AppColors.textPrimary, fontSize: 14),
                        ),
                        contentPadding: EdgeInsets.zero,
                      ),
                    ),
                    PopupMenuItem(
                      value: _AboutMenuAction.theme,
                      child: ListTile(
                        leading: Icon(
                          ThemeController.instance.resolve(context) == Brightness.dark
                              ? Icons.light_mode_outlined
                              : Icons.dark_mode_outlined,
                          size: 20,
                          color: AppColors.textSecondary,
                        ),
                        title: Text(
                          t.themeToggleTooltip,
                          style: TextStyle(color: AppColors.textPrimary, fontSize: 14),
                        ),
                        contentPadding: EdgeInsets.zero,
                      ),
                    ),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 28),
            Text(
              t.aboutEyebrow,
              style: TextStyle(
                fontSize: 12,
                fontWeight: FontWeight.w700,
                letterSpacing: 0.6,
                color: AppColors.accent,
              ),
            ),
            const SizedBox(height: 10),
            Text(t.aboutHeading, style: displayFont(fontSize: 24)),
            const SizedBox(height: 12),
            Text(
              t.aboutLead,
              style: TextStyle(fontSize: 14.5, height: 1.5, color: AppColors.textSecondary),
            ),
            const SizedBox(height: 24),
            GridView.count(
              crossAxisCount: 2,
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              crossAxisSpacing: 10,
              mainAxisSpacing: 10,
              childAspectRatio: 1.5,
              children: t.aboutStats.map((stat) => _StatCard(stat: stat)).toList(),
            ),
            const SizedBox(height: 24),
            Text(
              t.aboutBodyP1,
              style: TextStyle(fontSize: 14.5, height: 1.5, color: AppColors.textSecondary),
            ),
            const SizedBox(height: 12),
            Text(
              t.aboutBodyP2,
              style: TextStyle(fontSize: 14.5, height: 1.5, color: AppColors.textSecondary),
            ),
            const SizedBox(height: 28),
            Text(t.aboutStackHeading, style: displayFont(fontSize: 17)),
            const SizedBox(height: 4),
            Text(
              t.aboutStackCaption,
              style: TextStyle(fontSize: 13, color: AppColors.textTertiary),
            ),
            const SizedBox(height: 12),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: t.aboutStack.map((item) => _StackChip(label: item)).toList(),
            ),
            const SizedBox(height: 32),
            Container(height: 1, color: AppColors.border),
            const SizedBox(height: 16),
            BuyMeCoffeeButton(label: t.supportDevelopment, size: 26),
            const SizedBox(height: 56),
          ],
        ),
      ),
    );
  }
}

class _StatCard extends StatelessWidget {
  final AboutStat stat;

  const _StatCard({required this.stat});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: AppColors.surface,
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: AppColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(
            stat.value,
            style: displayFont(fontSize: 22, color: AppColors.accent),
          ),
          const SizedBox(height: 4),
          Text(
            stat.label,
            style: TextStyle(fontSize: 11.5, height: 1.3, color: AppColors.textTertiary),
          ),
        ],
      ),
    );
  }
}

class _StackChip extends StatelessWidget {
  final String label;

  const _StackChip({required this.label});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 11, vertical: 6),
      decoration: BoxDecoration(
        color: AppColors.surfaceRaised,
        borderRadius: BorderRadius.circular(999),
        border: Border.all(color: AppColors.border),
      ),
      child: Text(
        label,
        style: TextStyle(
          fontSize: 11,
          fontWeight: FontWeight.w600,
          letterSpacing: 0.4,
          color: AppColors.textSecondary,
        ),
      ),
    );
  }
}
