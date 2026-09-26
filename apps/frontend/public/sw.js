// Only the app shell is cached. Chart data must never be served from cache, so API and image
// requests are not intercepted.

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
  if (event.request.mode !== "navigate") return;

  event.respondWith(
    fetch(event.request).catch(() => caches.match("/").then((res) => res ?? Response.error())),
  );
});
