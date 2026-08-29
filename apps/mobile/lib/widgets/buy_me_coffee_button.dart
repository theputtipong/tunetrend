import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

import 'cached_thumbnail.dart';

const _bmacUrl = 'https://buymeacoffee.com/theputtipong';

const _bmacGifUrl =
    'https://media2.giphy.com/media/v1.Y2lkPTc5MGI3NjExOHJ1aHdnNjlrcXowa3h5YWp2YzljbHpkNzN1ZHQ0Ym5rdWNkcmsyNSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9cw/TDQOtnWgsBx99cNoyH/giphy.gif';

class BuyMeCoffeeButton extends StatelessWidget {
  final String label;
  final double size;

  const BuyMeCoffeeButton({super.key, required this.label, this.size = 40});

  @override
  Widget build(BuildContext context) {
    return InkWell(
      borderRadius: BorderRadius.circular(999),
      onTap: () =>
          launchUrl(Uri.parse(_bmacUrl), mode: LaunchMode.externalApplication),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          ClipOval(
            child: CachedThumbnail(
              imageUrl: _bmacGifUrl,
              width: size,
              height: size,
            ),
          ),
          const SizedBox(width: 8),
          Text(
            label,
            style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 13.5),
          ),
        ],
      ),
    );
  }
}
