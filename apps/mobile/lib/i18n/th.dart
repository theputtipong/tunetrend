import 'strings.dart';

String _minutesAgo(int n) => '$n นาทีที่แล้ว';
String _hoursAgo(int n) => '$n ชั่วโมงที่แล้ว';
String _daysAgo(int n) => '$n วันที่แล้ว';
String _weeksAgo(int n) => '$n สัปดาห์ที่แล้ว';

String _emptyDescription(String countryName) =>
    'ยังไม่มีข้อมูลชาร์ตของ$countryNameในหมวดนี้ ข้อมูลใหม่จะซิงค์ทุก 3 ชั่วโมง — กลับมาเช็คอีกครั้งเร็ว ๆ นี้';

String _closingIn(int n) => 'ปิดในอีก $n วินาที…';

const thStrings = AppStrings(
  aboutTooltip: 'เกี่ยวกับ TuneTrend',
  themeToggleTooltip: 'สลับธีมสว่าง/มืด',
  languageToggleTooltip: 'สลับภาษา',
  backToTrends: '← กลับไปหน้าชาร์ต',
  tabTrending: 'กำลังฮิต',
  tabNew: 'เพลงใหม่',
  tabMv: 'มิวสิควิดีโอ',
  errorTitle: 'โหลดเพลงกำลังฮิตไม่สำเร็จ',
  errorDescription: 'เราไม่สามารถเชื่อมต่อกับ TuneTrend service ได้ ตรวจสอบการเชื่อมต่อแล้วลองใหม่อีกครั้ง',
  retry: 'ลองใหม่',
  emptyTitle: 'ยังไม่มีเพลงกำลังฮิต',
  emptyDescription: _emptyDescription,
  viewsSuffix: 'วิว',
  minutesAgo: _minutesAgo,
  hoursAgo: _hoursAgo,
  daysAgo: _daysAgo,
  weeksAgo: _weeksAgo,
  aboutEyebrow: 'เกี่ยวกับ — TuneTrend',
  aboutHeading: 'สถาปัตยกรรมระดับธนาคาร ผสานข้อมูลชาร์ตเพลงยอดฮิต',
  aboutLead:
      'TuneTrend ไม่ได้เริ่มจากความอยากรู้เล่น ๆ ในวันหยุดว่าโลกกำลังฟังเพลงอะไร แต่ถูกสร้างขึ้นเป็นสนามทดลอง เพื่อนำประสบการณ์ full-stack กว่าทศวรรษ '
      'สถาปัตยกรรม microservices และระบบ deploy อัตโนมัติ มาประกอบเป็นผลิตภัณฑ์ multi-platform ที่ขัดเกลาจนเนียบ',
  aboutStats: [
    AboutStat('10+', 'ปีที่คลุกคลีกับระบบสเกลใหญ่'),
    AboutStat('7', 'ปีที่เขียนภาษาสาย mobile ยุคใหม่'),
    AboutStat('90%+', 'มาตรฐาน test coverage ที่รักษาไว้'),
    AboutStat('100%', 'CI/CD deploy และ monitoring แบบอัตโนมัติ'),
  ],
  aboutBodyP1:
      'ตอนกลางวัน โฟกัสหลักคือโครงสร้างพื้นฐานด้านการเงินตัวจริง — ออกแบบ microservices '
      'สำหรับระบบชำระเงินที่ซับซ้อน จัดการโค้ดให้พร้อมผ่าน penetration test '
      'และดูแลความเสถียรของระบบให้กับหนึ่งในธนาคารที่ใหญ่ที่สุดของไทย',
  aboutBodyP2:
      'ตอนกลางคืน ความเข้มงวดแบบเดียวกันนี้ขับเคลื่อน TuneTrend ระบบถูกออกแบบให้ดึงข้อมูลชาร์ต YouTube '
      'แบบสดจาก 5 ประเทศ รีเฟรชทุก 3 ชั่วโมง วิ่งผ่าน Go backend ที่แข็งแรง ให้บริการทั้งเว็บแอป Next.js '
      'และแอปมือถือ Flutter — ใช้ design system เดียวกันและ pipeline CI/CD เดียวกันทั้งหมด',
  aboutStackHeading: 'สถาปัตยกรรมระบบ',
  aboutStackCaption:
      'ออกแบบมาเพื่อสเกล ความเร็ว และความสอดคล้องกันทุกแพลตฟอร์ม',
  aboutStack: [
    'GO · FIBER',
    'NEXT.JS',
    'FLUTTER',
    'POSTGRESQL',
    'CI/CD PIPELINES',
  ],
  onboardingCountryTitle: 'เลือกประเทศ',
  onboardingCountryDescription: 'สลับดูชาร์ตของ TH, KR, JP, US, GB ได้ที่นี่',
  onboardingTabsTitle: 'เลือกหมวดหมู่',
  onboardingTabsDescription: 'สลับดูกำลังฮิต เพลงใหม่ หรือมิวสิควิดีโอได้ที่นี่',
  onboardingLanguageTitle: 'ภาษา',
  onboardingLanguageDescription: 'สลับภาษาทั้งแอประหว่างไทยกับอังกฤษ',
  onboardingThemeTitle: 'ธีมสว่าง/มืด',
  onboardingThemeDescription: 'สลับธีมได้ตามใจ — ปกติ TuneTrend จะตามธีมของระบบเครื่องคุณอยู่แล้ว',
  onboardingAboutTitle: 'เกี่ยวกับ TuneTrend',
  onboardingAboutDescription: 'ดูเบื้องหลังการพัฒนาแอปนี้ และช่องทางติดต่อ',
  onboardingMenuTitle: 'ตัวเลือกเพิ่มเติม',
  onboardingMenuDescription: 'กดที่นี่เพื่อดูแนะนำการใช้งานระบบ สลับภาษา สลับธีม หรือดูเกี่ยวกับแอปนี้',
  menuTooltip: 'ตัวเลือกเพิ่มเติม',
  contactOpenButton: 'ติดต่อฉัน',
  contactNameLabel: 'ชื่อ (ไม่บังคับ)',
  contactNamePlaceholder: 'ชื่อของคุณ',
  contactMessageLabel: 'ข้อความ',
  contactMessagePlaceholder: 'อยากบอกอะไรกับผมบ้าง?',
  contactMethodLabel: 'อยากให้ติดต่อกลับทางไหน?',
  contactMethodEmail: 'อีเมล',
  contactMethodPhone: 'เบอร์โทร',
  contactEmailPlaceholder: 'you@example.com',
  contactPhonePlaceholder: '08xxxxxxxx',
  contactSubmit: 'ส่งข้อความ',
  contactSubmitting: 'กำลังส่ง…',
  contactSuccessMessage: 'ขอบคุณครับ ข้อความถูกส่งเรียบร้อยแล้ว',
  closingIn: _closingIn,
  contactErrorMessageRequired: 'กรุณากรอกข้อความ',
  contactErrorInvalidEmail: 'กรุณากรอกอีเมลให้ถูกต้อง',
  contactErrorInvalidPhone: 'กรุณากรอกเบอร์โทรไทยให้ถูกต้อง (เช่น 0812345678)',
  contactErrorTooLong: 'ข้อความยาวเกินไป (ไม่เกิน 2000 ตัวอักษร)',
  contactErrorRateLimited: 'ส่งบ่อยเกินไป กรุณาลองใหม่อีกครั้งในอีกสักครู่',
  contactErrorGeneric: 'เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้งภายหลัง',
  supportDevelopment: 'สนับสนุนการพัฒนาแอปนี้',
  replayTourTooltip: 'ดูแนะนำการใช้งานระบบอีกครั้ง',
);
