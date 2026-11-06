import 'strings.dart';

String _minutesAgo(int n) => '$n นาทีที่แล้ว';
String _hoursAgo(int n) => '$n ชั่วโมงที่แล้ว';
String _daysAgo(int n) => '$n วันที่แล้ว';
String _weeksAgo(int n) => '$n สัปดาห์ที่แล้ว';

String _emptyDescription(String countryName) =>
    'ยังไม่มีข้อมูลชาร์ตของ$countryNameในหมวดนี้ ข้อมูลใหม่จะซิงค์ทุก 3 ชั่วโมง — กลับมาเช็คอีกครั้งเร็ว ๆ นี้';

String _closingIn(int n) => 'ปิดในอีก $n วินาที…';
String _autoPlayPromptCountdown(int n) =>
    'จะเล่นต่อเนื่องอัตโนมัติในอีก $n วินาที…';
String _autoAdvanceCountdown(int n) => 'เล่นเพลงถัดไปในอีก $n วินาที…';
String _shareMessage(String title, String url, String? appUrl) {
  final base = '$title\n$url\n\nแชร์จากแอป TuneTrend';
  return appUrl == null ? base : '$base\nดาวน์โหลดแอป: $appUrl';
}

String _searchHits(String n) => '$n ค้นหา';
String _youtubeRemaining(int remaining, int limit) =>
    'เหลือ $remaining/$limit ครั้งวันนี้';
String _youtubeQuotaExhausted(String? time) => time == null
    ? 'วันนี้มีคนค้นเยอะจนโควตาเต็ม ระหว่างนี้ยังค้นในชาร์ต TuneTrend ได้ตามปกติ'
    : 'วันนี้มีคนค้นเยอะจนโควตาเต็ม จะกลับมาประมาณ $time น. ระหว่างนี้ยังค้นในชาร์ต TuneTrend ได้ตามปกติ';
String _youtubeUnavailable(String? time) => time == null
    ? 'ลองใหม่อีกครั้งภายหลัง ระหว่างนี้ยังค้นในชาร์ต TuneTrend ได้ตามปกติ'
    : 'ลองใหม่อีกครั้งประมาณ $time น. ระหว่างนี้ยังค้นในชาร์ต TuneTrend ได้ตามปกติ';
String _youtubeUserLimit(int limit) =>
    'วันนี้ค้นบน YouTube ครบ $limit ครั้งแล้ว พรุ่งนี้ค้นได้อีก หรือลองค้นในชาร์ตของเราดูก่อน';
String _playlistSongCount(int n) => '$n เพลง';
String _playlistNameInvalid(int max) => 'ชื่อต้องยาว 1–$max ตัวอักษร';
String _deletePlaylistBody(String name) =>
    'ลบ "$name" และเพลงทั้งหมดในนั้น กู้คืนไม่ได้';
String _saveToPlaylistTitle(String title) => 'บันทึก "$title" ลงใน…';
String _savedToPlaylist(String name) => 'บันทึกลง $name แล้ว';
String _alreadyInPlaylist(String name) => 'มีเพลงนี้ใน $name แล้ว';

const thStrings = AppStrings(
  aboutTooltip: 'เกี่ยวกับ TuneTrend',
  themeToggleTooltip: 'สลับธีมสว่าง/มืด',
  languageToggleTooltip: 'สลับภาษา',
  backToTrends: '← กลับไปหน้าชาร์ต',
  tabTrending: 'กำลังฮิต',
  tabNew: 'เพลงใหม่',
  tabNewGeneric: 'วิดีโอใหม่',
  tabMv: 'มิวสิควิดีโอ',
  musicCategory: 'เพลง',
  autoPlayPromptTitle: 'เล่นต่อเนื่องเลยไหม?',
  autoPlayPromptDescription:
      'เล่นวิดีโอถัดไปในลำดับนี้ให้อัตโนมัติทุกครั้งที่วิดีโอจบ',
  autoPlayPromptAccept: 'เล่นต่อเนื่อง',
  autoPlayPromptDecline: 'ไม่ ขอบคุณ',
  autoPlayPromptCountdown: _autoPlayPromptCountdown,
  autoAdvanceCountdown: _autoAdvanceCountdown,
  autoPlayStoppedMessage: 'หยุดเล่นต่อเนื่องแล้ว เนื่องจากเปลี่ยนแท็บ',
  viewCountNoticeTooltip: 'เกี่ยวกับยอดวิว',
  viewCountNoticeTitle: 'เกี่ยวกับยอดวิว',
  viewCountNoticeBody:
      'วิดีโอนี้เล่นผ่านตัวเล่นวิดีโอ (embedded player) อย่างเป็นทางการของ YouTube '
      'ยอดวิว เวลาในการรับชม และรายได้จากโฆษณาของผู้สร้างคอนเทนต์ จะถูกนับตามปกติ '
      'เหมือนรับชมบน youtube.com ทุกประการ TuneTrend เป็นเพียงช่องทางรวบรวมข้อมูลเทรนด์สาธารณะ '
      'ผ่าน YouTube Data API เท่านั้น ไม่มีการดาวน์โหลด นำไปเผยแพร่ซ้ำ '
      'หรือแทรกแซงคอนเทนต์ของท่านแต่อย่างใด',
  viewCountNoticeDismiss: 'เข้าใจแล้ว',
  maintenanceTitle: 'TuneTrend กำลังปิดปรับปรุงระบบ',
  maintenanceDescription: 'เรากำลังปรับปรุงระบบเพื่อประสบการณ์ที่ดีขึ้น จะกลับมาให้บริการเร็วๆ นี้ ขอบคุณที่รอครับ',
  updateRequiredTitle: 'ต้องอัปเดตแอป',
  updateRequiredDescription:
      'TuneTrend เวอร์ชันใหม่พร้อมใช้งานแล้ว กรุณาอัปเดตเพื่อใช้งานแอปต่อ',
  updateNowButton: 'อัปเดตเลย',
  shareTooltip: 'แชร์',
  shareMessage: _shareMessage,
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
  onboardingTabsTitle: 'สลับมุมมอง',
  onboardingTabsDescription:
      'สลับดูกำลังฮิต เพลงใหม่ หรือมิวสิควิดีโอได้ที่นี่',
  onboardingCategoryTitle: 'กรองตามหมวดหมู่',
  onboardingCategoryDescription: 'เลือกดูเฉพาะหมวดหมู่ YouTube ที่สนใจ เช่น เกม หรือกีฬา ปัดซ้าย-ขวาเพื่อดูหมวดหมู่เพิ่มเติมได้',
  onboardingLanguageTitle: 'ภาษา',
  onboardingLanguageDescription: 'สลับภาษาทั้งแอประหว่างไทยกับอังกฤษ',
  onboardingThemeTitle: 'ธีมสว่าง/มืด',
  onboardingThemeDescription:
      'สลับธีมได้ตามใจ — ปกติ TuneTrend จะตามธีมของระบบเครื่องคุณอยู่แล้ว',
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
  privacyPolicy: 'นโยบายความเป็นส่วนตัว',
  navCharts: 'ชาร์ต',
  navSearch: 'ค้นหา',
  navPlaylists: 'Playlist',
  tabSearched: 'มาแรงจากการค้นหา',
  searchedTabCaption: 'เพลงที่คนในประเทศนี้ค้นหาแล้วกดเล่นมากที่สุดใน 7 วัน',
  newBadge: 'ใหม่',
  searchHits: _searchHits,
  searchHint: 'ค้นหาเพลง ศิลปิน…',
  recentSearches: 'ค้นหาล่าสุด:',
  clearSearchTooltip: 'ล้าง',
  searchPromptTitle: 'ค้นหาเพลงที่อยากฟัง',
  searchPromptDescription: 'ค้นในชาร์ต TuneTrend ได้ทันทีโดยไม่ต้องเข้าสู่ระบบ',
  searchQueryTooShort: 'พิมพ์อย่างน้อย 2 ตัวอักษร',
  localSectionTitle: 'ในชาร์ต',
  localEmpty: 'ไม่พบในชาร์ต TuneTrend',
  searchedBefore: 'เคยถูกค้นหา',
  youtubeSectionTitle: 'จาก YouTube',
  youtubeMoreHint: 'ยังไม่เจอเพลงที่หา? ค้นต่อบน YouTube ได้ ผลที่เจอจะถูกเก็บไว้ให้คนอื่นค้นเจอเร็วขึ้นด้วย',
  youtubeSearchButton: 'ค้นหาเพิ่มบน YouTube',
  youtubeRemaining: _youtubeRemaining,
  youtubeLoading: 'กำลังโหลด…',
  youtubeEmpty: 'ไม่พบผลบน YouTube',
  youtubeLoginTitle: 'เข้าสู่ระบบเพื่อค้นต่อบน YouTube',
  youtubeLoginDescription:
      'ผลในระบบยังค้นได้ตามปกติ การค้นบน YouTube เปิดให้เฉพาะสมาชิก',
  signInWithGoogle: 'เข้าสู่ระบบด้วย Google',
  loginUnavailable: 'ระบบ login ขัดข้องชั่วคราว ลองใหม่อีกครั้ง',
  videoUnavailable: 'วิดีโอนี้ไม่พร้อมใช้งานแล้ว',
  signInFailed: 'เข้าสู่ระบบไม่สำเร็จ ลองใหม่อีกครั้ง',
  youtubeUnavailableTitle: 'ค้นบน YouTube ใช้ไม่ได้ชั่วคราว',
  youtubeQuotaExhausted: _youtubeQuotaExhausted,
  youtubeUnavailable: _youtubeUnavailable,
  youtubeUserLimitTitle: 'วันนี้ค้นบน YouTube ครบแล้ว',
  youtubeUserLimit: _youtubeUserLimit,
  youtubeErrorTitle: 'ค้นบน YouTube ไม่สำเร็จ',
  youtubeErrorDescription: 'เกิดข้อผิดพลาด ลองใหม่อีกครั้ง',
  playlistsTitle: 'Playlist ของฉัน',
  playlistSongCount: _playlistSongCount,
  newPlaylist: 'Playlist ใหม่',
  createNewPlaylist: 'สร้าง playlist ใหม่',
  playlistNameLabel: 'ชื่อ playlist',
  playlistNameInvalid: _playlistNameInvalid,
  create: 'สร้าง',
  save: 'บันทึก',
  cancel: 'ยกเลิก',
  delete: 'ลบ',
  playlistsEmptyTitle: 'ยังไม่มี playlist',
  playlistsEmptyDescription: 'กด + บนเพลงที่ชอบเพื่อเก็บไว้ฟังทีหลัง',
  playlistItemsEmpty: 'ยังไม่มีเพลงใน playlist นี้',
  playlistPrivate: 'Playlist ส่วนตัว',
  playAll: 'เล่นทั้งหมด',
  rename: 'แก้ชื่อ',
  deletePlaylistTitle: 'ลบ playlist นี้?',
  deletePlaylistBody: _deletePlaylistBody,
  removeFromPlaylist: 'เอาออก',
  addToPlaylistTooltip: 'บันทึกลง playlist',
  saveToPlaylistTitle: _saveToPlaylistTitle,
  savedToPlaylist: _savedToPlaylist,
  alreadyInPlaylist: _alreadyInPlaylist,
  actionFailed: 'ทำรายการไม่สำเร็จ ลองใหม่อีกครั้ง',
  loginPromptTitle: 'เก็บเพลงไว้ฟังทีหลัง',
  loginPromptDescription: 'เข้าสู่ระบบด้วย Google เพื่อสร้าง playlist ของคุณเอง ใช้ได้ทั้งบนเว็บและแอป',
  continueWithGoogle: 'ดำเนินการต่อด้วย Google',
  loginPrivacyNote:
      'เราเก็บแค่ชื่อ อีเมล และรูปโปรไฟล์ ไม่มีการโพสต์อะไรลง YouTube ของคุณ',
  accountMenuTooltip: 'บัญชี',
  signOut: 'ออกจากระบบ',
  deleteAccount: 'ลบบัญชี',
  deleteAccountTitle: 'ลบบัญชีถาวร?',
  deleteAccountBody:
      'บัญชีและ playlist ทั้งหมดของคุณจะถูกลบถาวรและกู้คืนไม่ได้',
  accountDeleted: 'ลบบัญชีแล้ว',
);
