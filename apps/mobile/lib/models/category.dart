class Category {
  final String id;
  final String countryCode;
  final String title;
  final bool assignable;

  const Category({
    required this.id,
    required this.countryCode,
    required this.title,
    required this.assignable,
  });

  factory Category.fromJson(Map<String, dynamic> json) {
    return Category(
      id: json['id'] as String? ?? '',
      countryCode: json['countryCode'] as String? ?? '',
      title: json['title'] as String? ?? '',
      assignable: json['assignable'] as bool? ?? false,
    );
  }
}
