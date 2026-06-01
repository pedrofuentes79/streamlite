// ------- Streamlite frontend logic (vanilla, no build step) -------
// Runs entirely in the browser on the phone. Talks to the Go server through:
//   GET   /api/media          -> catalog
//   GET   /api/audio/{id}     -> audio stream (background-listening friendly)
//   PATCH /api/progress/{id}  -> save listening position
//
// Background audio note: on iOS this only keeps playing when the page runs as a
// normal Safari TAB (not an installed/standalone PWA). The Media Session API
// below provides the lock-screen / control-center play/pause/scrub controls.

const catalogEl = document.getElementById("catalog");
const statusEl  = document.getElementById("catalog-status");
const catalogView = document.getElementById("catalog-view");
const playerView  = document.getElementById("player-view");
const playerTitle = document.getElementById("player-title");
const playerDate  = document.getElementById("player-date");
const audio       = document.getElementById("player");
const backBtn     = document.getElementById("back-btn");

let current = null;      // { id, title, totalSeconds }
let lastSaved = -1;      // last progress value sent, to avoid spamming

// ---- Helpers ----

function fmtDate(iso) {
  try {
    return new Date(iso).toLocaleDateString("es-ES", {
      day: "numeric", month: "short", year: "numeric",
    });
  } catch { return iso; }
}

function fmtClock(totalSec) {
  const s = Math.round(totalSec);
  const h = Math.floor(s / 3600);
  const m = Math.floor((s % 3600) / 60);
  return h > 0 ? `${h}h ${m}m` : `${m}m`;
}

// ---- Catalog ----

async function loadCatalog() {
  statusEl.textContent = "Cargando…";
  try {
    const res = await fetch("/api/media");
    if (!res.ok) throw new Error(res.status);
    renderCatalog(await res.json());
  } catch {
    statusEl.classList.add("status");
    statusEl.textContent = "No se pudo cargar el catálogo.";
    catalogEl.replaceChildren(statusEl);
  }
}

function renderCatalog(items) {
  if (!items.length) {
    statusEl.textContent = "No hay nada todavía.";
    catalogEl.replaceChildren(statusEl);
    return;
  }

  const cards = items.map((m) => {
    const pct = m.total_seconds > 0
      ? Math.min(100, (m.progress_seconds / m.total_seconds) * 100)
      : 0;

    const li = document.createElement("li");
    li.className = "card";
    li.innerHTML = `
      <div class="card-title"></div>
      <div class="card-meta">
        <span class="card-date"></span>
        <span class="card-dur"></span>
      </div>
      <div class="progress"><div class="progress-fill" style="width:${pct}%"></div></div>`;
    li.querySelector(".card-title").textContent = m.title;
    li.querySelector(".card-date").textContent = fmtDate(m.date);
    li.querySelector(".card-dur").textContent = fmtClock(m.total_seconds);

    li.addEventListener("click", () => openPlayer(m));
    return li;
  });

  catalogEl.replaceChildren(...cards);
}

// ---- Player ----

function openPlayer(m) {
  current = { id: m.id, title: m.title, totalSeconds: m.total_seconds };
  lastSaved = m.progress_seconds;

  playerTitle.textContent = m.title;
  playerDate.textContent = fmtDate(m.date);
  audio.src = `/api/audio/${m.id}`;

  catalogView.classList.add("hidden");
  playerView.classList.remove("hidden");

  setupMediaSession(m);

  // Resume where we left off, once the browser knows the duration.
  const resume = () => {
    if (m.progress_seconds > 0 && m.progress_seconds < audio.duration) {
      audio.currentTime = m.progress_seconds;
    }
    audio.removeEventListener("loadedmetadata", resume);
  };
  audio.addEventListener("loadedmetadata", resume);
  audio.play().catch(() => { /* iOS may require a tap on the play button */ });
}

function closePlayer() {
  saveProgress();
  audio.pause();
  audio.removeAttribute("src");
  audio.load();
  current = null;

  playerView.classList.add("hidden");
  catalogView.classList.remove("hidden");
  loadCatalog();
}

// ---- Media Session (lock screen + control center) ----

function setupMediaSession(m) {
  if (!("mediaSession" in navigator)) return;

  navigator.mediaSession.metadata = new MediaMetadata({
    title: m.title,
    artist: "Streamlite",
    album: fmtDate(m.date),
    artwork: [
      { src: "/icons/icon-192.png", sizes: "192x192", type: "image/png" },
      { src: "/icons/icon-512.png", sizes: "512x512", type: "image/png" },
    ],
  });

  const set = (action, handler) => {
    try { navigator.mediaSession.setActionHandler(action, handler); }
    catch { /* action not supported on this OS */ }
  };

  set("play",  () => audio.play());
  set("pause", () => audio.pause());
  set("seekbackward", (d) => {
    audio.currentTime = Math.max(0, audio.currentTime - (d.seekOffset || 15));
  });
  set("seekforward", (d) => {
    audio.currentTime = Math.min(audio.duration || Infinity, audio.currentTime + (d.seekOffset || 30));
  });
  set("seekto", (d) => {
    if (d.fastSeek && "fastSeek" in audio) audio.fastSeek(d.seekTime);
    else audio.currentTime = d.seekTime;
  });
}

function updatePositionState() {
  if (!("mediaSession" in navigator) || !navigator.mediaSession.setPositionState) return;
  if (!isFinite(audio.duration)) return;
  try {
    navigator.mediaSession.setPositionState({
      duration: audio.duration,
      playbackRate: audio.playbackRate,
      position: Math.min(audio.currentTime, audio.duration),
    });
  } catch { /* ignore invalid states */ }
}

// ---- Progress saving ----

function saveProgress(beacon = false) {
  if (!current) return;

  // Clamp to the catalog total so we never trip the DB's
  // CHECK(progress_seconds <= total_seconds) constraint.
  let pos = Math.floor(audio.currentTime || 0);
  pos = Math.min(pos, Math.floor(current.totalSeconds));
  if (pos < 0 || pos === lastSaved) return;
  lastSaved = pos;

  fetch(`/api/progress/${current.id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ new_progress: pos }),
    keepalive: beacon,
  }).catch(() => {});
}

// ---- Wiring ----

backBtn.addEventListener("click", closePlayer);

audio.addEventListener("pause", () => {
  if ("mediaSession" in navigator) navigator.mediaSession.playbackState = "paused";
  saveProgress();
});
audio.addEventListener("play", () => {
  if ("mediaSession" in navigator) navigator.mediaSession.playbackState = "playing";
});

let tick = 0;
audio.addEventListener("timeupdate", () => {
  updatePositionState();
  if (audio.currentTime - tick >= 5) {     // throttle saves to ~every 5s
    tick = audio.currentTime;
    saveProgress();
  }
});

// iOS fires these when backgrounded / closed — flush with keepalive.
window.addEventListener("pagehide", () => saveProgress(true));
document.addEventListener("visibilitychange", () => {
  if (document.visibilityState === "hidden") saveProgress(true);
});

// ---- PWA service worker (caches the UI shell only) ----
if ("serviceWorker" in navigator) {
  window.addEventListener("load", () => {
    navigator.serviceWorker.register("/sw.js").catch(() => {});
  });
}

loadCatalog();
