final RegExp emailPattern = RegExp(r'^[^\s@]+@[^\s@]+\.[^\s@]+$');

final RegExp thaiPhonePattern = RegExp(r'^0(?:[689]\d{8}|[2-57]\d{7})$');

const int maxMessageLen = 2000;
const int maxNameLen = 100;

String normalizeThaiPhone(String raw) => raw.replaceAll(RegExp(r'[\s()-]'), '');

bool isValidEmail(String v) => emailPattern.hasMatch(v.trim());

bool isValidThaiPhone(String v) => thaiPhonePattern.hasMatch(normalizeThaiPhone(v));
