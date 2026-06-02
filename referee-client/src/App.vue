<template>
  <div class="app-shell">
    <header class="topbar">
      <div class="brand">
        <span class="brand-mark" aria-hidden="true"></span>
        <div>
          <p class="eyebrow">Referee Console</p>
          <h1>乒乓球裁判系统</h1>
        </div>
      </div>
      <div class="topbar-meta">
        <span>现场记分</span>
        <span>录像取证</span>
      </div>
    </header>

    <main class="app-main">
      <router-view></router-view>
    </main>
  </div>
</template>

<script setup>
</script>

<style>
:root {
  color-scheme: light;
  --ink: #17201c;
  --muted: #66726d;
  --panel: #ffffff;
  --panel-soft: #f6f7f2;
  --line: #dfe5da;
  --court: #0f5f52;
  --court-dark: #123f3b;
  --lime: #cfe85f;
  --amber: #f2b84b;
  --danger: #d74b4b;
  --shadow: 0 24px 70px rgba(21, 42, 35, 0.14);
}

* {
  box-sizing: border-box;
}

html,
body,
#app {
  min-height: 100%;
}

body {
  margin: 0;
  font-family: "Microsoft YaHei", "PingFang SC", "Segoe UI", sans-serif;
  color: var(--ink);
  background:
    linear-gradient(135deg, rgba(15, 95, 82, 0.08), transparent 36%),
    linear-gradient(315deg, rgba(207, 232, 95, 0.16), transparent 30%),
    #eef2eb;
}

button,
input,
select {
  font: inherit;
}

button {
  transition: transform 0.18s ease, box-shadow 0.18s ease, background 0.18s ease, border-color 0.18s ease;
}

button:not(:disabled):hover {
  transform: translateY(-1px);
}

button:disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

.app-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  padding: 18px clamp(18px, 4vw, 48px);
  color: #fff;
  background:
    linear-gradient(120deg, rgba(207, 232, 95, 0.16), transparent 42%),
    linear-gradient(135deg, #123f3b 0%, #0f5f52 56%, #17201c 100%);
  border-bottom: 4px solid var(--lime);
  box-shadow: 0 14px 40px rgba(18, 63, 59, 0.22);
}

.brand {
  display: flex;
  align-items: center;
  gap: 14px;
}

.brand-mark {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  background:
    radial-gradient(circle at 35% 35%, #fff 0 13%, transparent 14%),
    linear-gradient(145deg, var(--lime), var(--amber));
  box-shadow: inset -8px -9px 0 rgba(23, 32, 28, 0.14), 0 8px 22px rgba(0, 0, 0, 0.2);
}

.eyebrow {
  margin: 0 0 3px;
  color: rgba(255, 255, 255, 0.72);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0;
  text-transform: uppercase;
}

h1 {
  margin: 0;
  font-size: clamp(22px, 3vw, 32px);
  line-height: 1.15;
  letter-spacing: 0;
}

.topbar-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.topbar-meta span {
  display: inline-flex;
  align-items: center;
  min-height: 32px;
  padding: 0 12px;
  border: 1px solid rgba(255, 255, 255, 0.22);
  border-radius: 999px;
  color: rgba(255, 255, 255, 0.82);
  background: rgba(255, 255, 255, 0.08);
  font-size: 13px;
}

.app-main {
  flex: 1;
  width: 100%;
  padding: clamp(18px, 4vw, 42px);
}

@media (max-width: 720px) {
  .topbar {
    align-items: flex-start;
    flex-direction: column;
  }

  .topbar-meta {
    width: 100%;
  }
}
</style>
