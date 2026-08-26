// Service worker แบบมือ ไม่พึ่ง library (ดูเหตุผลใน plan: next-pwa/@ducanh2912
// hook เข้า webpack config โดยตรง เข้ากันไม่แน่นอนกับ Turbopack ที่โปรเจกต์นี้ใช้)
//
// Scope ตั้งใจจำกัดแค่ "เปิดแอปได้แม้เน็ตหลุดชั่วครู่" (เห็น app shell แทน
// browser error page) ไม่ใช่ full offline data — จึงไม่แคช response จาก
// /trends*, /api/* เด็ดขาด เพราะข้อมูลชาร์ตเป็น live data ที่เปลี่ยนทุก
// 3 ชั่วโมง แคชไว้จะทำให้เห็นข้อมูลเก่าโดยไม่รู้ตัว

const CACHE_NAME = "tunetrend-shell-v1";
const SHELL_URLS = ["/", "/icon.svg", "/apple-icon.png"];

self.addEventListener("install", (event) => {
  event.waitUntil(caches.open(CACHE_NAME).then((cache) => cache.addAll(SHELL_URLS)));
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((keys) =>
        Promise.all(keys.filter((key) => key !== CACHE_NAME).map((key) => caches.delete(key))),
      ),
  );
});

self.addEventListener("fetch", (event) => {
  // เฉพาะ navigation request (เปิด/reload หน้า) เท่านั้น — ปล่อย request อื่น
  // ทั้งหมด (รวม /trends*, /api/*, รูปภาพ thumbnail จาก YouTube ฯลฯ) ให้ผ่าน
  // ไปตามปกติ ไม่ intercept
  if (event.request.mode !== "navigate") return;

  event.respondWith(
    fetch(event.request).catch(() => caches.match("/").then((res) => res ?? Response.error())),
  );
});
