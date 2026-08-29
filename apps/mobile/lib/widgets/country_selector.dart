import 'package:flutter/material.dart';

import '../constants/countries.dart';
import 'pill_button.dart';

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
          return PillButton(
            label: country.code,
            active: active,
            onTap: () => onSelect(country.code),
          );
        },
      ),
    );
  }
}
