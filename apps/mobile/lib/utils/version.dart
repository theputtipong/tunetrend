bool isVersionLower(String current, String minimum) {
  final c = _parseVersion(current);
  final m = _parseVersion(minimum);
  for (var i = 0; i < 3; i++) {
    if (c[i] != m[i]) return c[i] < m[i];
  }
  return false;
}

List<int> _parseVersion(String version) {
  final parts = version.split('.');
  return List.generate(
    3,
    (i) => i < parts.length ? int.tryParse(parts[i]) ?? 0 : 0,
  );
}
