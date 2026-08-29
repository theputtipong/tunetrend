import 'package:flutter/material.dart';

import '../i18n/lang.dart';
import '../models/category.dart';
import '../utils/categories.dart';
import 'pill_button.dart';
import 'scroll_hint_chevron.dart';

class CategoryFilter extends StatefulWidget {
  final String active;
  final List<Category> categories;
  final String musicLabel;
  final ValueChanged<String> onSelect;

  const CategoryFilter({
    super.key,
    required this.active,
    required this.categories,
    required this.musicLabel,
    required this.onSelect,
  });

  @override
  State<CategoryFilter> createState() => _CategoryFilterState();
}

class _CategoryFilterState extends State<CategoryFilter> {
  final _scrollController = ScrollController();
  bool _showScrollHint = false;

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
    WidgetsBinding.instance.addPostFrameCallback((_) => _checkOverflow());
  }

  @override
  void didUpdateWidget(CategoryFilter oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.categories != widget.categories) {
      WidgetsBinding.instance.addPostFrameCallback((_) => _checkOverflow());
    }
  }

  void _checkOverflow() {
    if (!mounted || !_scrollController.hasClients) return;
    final overflows = _scrollController.position.maxScrollExtent > 0;
    if (overflows != _showScrollHint) {
      setState(() => _showScrollHint = overflows);
    }
  }

  void _onScroll() {
    if (_scrollController.offset > 4 && _showScrollHint) {
      setState(() => _showScrollHint = false);
    }
  }

  @override
  void dispose() {
    _scrollController.removeListener(_onScroll);
    _scrollController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (widget.categories.isEmpty) return const SizedBox.shrink();

    final lang = LangController.instance.resolve();
    return SizedBox(
      height: 36,
      child: Stack(
        alignment: Alignment.centerRight,
        children: [
          ListView.separated(
            controller: _scrollController,
            scrollDirection: Axis.horizontal,
            itemCount: widget.categories.length + 1,
            separatorBuilder: (_, _) => const SizedBox(width: 8),
            itemBuilder: (context, index) {
              if (index == 0) {
                return PillButton(
                  label: widget.musicLabel,
                  active: widget.active.isEmpty,
                  onTap: () => widget.onSelect(''),
                );
              }
              final category = widget.categories[index - 1];
              return PillButton(
                label: categoryLabel(category.id, category.title, lang),
                active: category.id == widget.active,
                onTap: () => widget.onSelect(category.id),
              );
            },
          ),
          if (_showScrollHint) const ScrollHintChevron(),
        ],
      ),
    );
  }
}
