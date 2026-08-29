import '../i18n/lang.dart';

const Map<String, String> _categoryLabelsTh = {
  '1': 'ภาพยนตร์และแอนิเมชัน',
  '2': 'รถยนต์และยานพาหนะ',
  '15': 'สัตว์เลี้ยงและสัตว์',
  '17': 'กีฬา',
  '19': 'การเดินทางและกิจกรรม',
  '20': 'เกม',
  '22': 'ผู้คนและบล็อก',
  '23': 'ตลก',
  '24': 'บันเทิง',
  '25': 'ข่าวและการเมือง',
  '26': 'วิธีทำและสไตล์',
  '27': 'การศึกษา',
  '28': 'วิทยาศาสตร์และเทคโนโลยี',
};

String categoryLabel(String id, String title, AppLang lang) {
  if (lang != AppLang.th) return title;
  return _categoryLabelsTh[id] ?? title;
}
