import 'package:flutter/material.dart';
import '../constants/countries.dart';
import '../constants/theme.dart';

class CountrySelector extends StatelessWidget {
  final String activeCountry;
  final ValueChanged<String> onSelect;

  const CountrySelector({
    super.key,
    required this.activeCountry,
    required this.onSelect,
  });

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 36,
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        itemCount: kCountries.length,
        separatorBuilder: (_, _) => const SizedBox(width: 8),
        itemBuilder: (context, index) {
          final country = kCountries[index];
          final active = country.code == activeCountry;
          return _CountryPill(
            label: country.code,
            active: active,
            onTap: () => onSelect(country.code),
          );
        },
      ),
    );
  }
}

class _CountryPill extends StatelessWidget {
  final String label;
  final bool active;
  final VoidCallback onTap;

  const _CountryPill({required this.label, required this.active, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return InkWell(
      borderRadius: BorderRadius.circular(999),
      onTap: onTap,
      child: Container(
        alignment: Alignment.center,
        padding: const EdgeInsets.symmetric(horizontal: 14),
        decoration: BoxDecoration(
          color: active ? AppColors.accent : AppColors.surface,
          borderRadius: BorderRadius.circular(999),
          border: active ? null : Border.all(color: AppColors.border),
        ),
        child: Text(
          label,
          style: TextStyle(
            fontWeight: FontWeight.w600,
            fontSize: 13,
            color: active ? AppColors.accentInk : AppColors.textSecondary,
          ),
        ),
      ),
    );
  }
}
