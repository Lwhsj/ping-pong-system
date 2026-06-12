<template>
  <el-container class="layout-container">
    <el-aside width="260px" class="aside">
      <div class="logo">
        <div class="logo-mark">
          <el-icon><Trophy /></el-icon>
        </div>
        <div>
          <div class="logo-title">赛事辅助系统</div>
          <div class="logo-subtitle">Match Control</div>
        </div>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="el-menu-vertical"
        router
        background-color="transparent"
        text-color="#a7b5ca"
        active-text-color="#ffffff"
      >
        <el-menu-item index="/matches">
          <el-icon><List /></el-icon>
          <span>比赛历史</span>
        </el-menu-item>
        <el-menu-item index="/scoreboard">
          <el-icon><Monitor /></el-icon>
          <span>实时比分</span>
        </el-menu-item>
      </el-menu>
      <div class="aside-footer">
        <span class="status-dot"></span>
        <span>赛事数据中心</span>
      </div>
    </el-aside>

    <el-container class="workspace">
      <el-header class="header">
        <div class="header-content">
          <div>
            <p class="eyebrow">Competition Dashboard</p>
            <h3>乒乓球赛事辅助系统</h3>
          </div>
          <div class="header-badge">
            <span class="pulse"></span>
            实时监控
          </div>
        </div>
      </el-header>
      
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const activeMenu = computed(() => {
  if (route.path === '/match-detail') return '/matches'
  return route.path
})
</script>

<style scoped>
.layout-container {
  height: 100vh;
  background: #eef3f7;
}

.aside {
  position: relative;
  overflow: hidden;
  color: #fff;
  display: flex;
  flex-direction: column;
  background:
    linear-gradient(160deg, rgba(50, 211, 197, 0.16), transparent 34%),
    linear-gradient(340deg, rgba(255, 176, 0, 0.18), transparent 32%),
    #142033;
  box-shadow: 18px 0 42px rgba(25, 38, 58, 0.14);
}

.aside::before {
  position: absolute;
  inset: 0;
  pointer-events: none;
  content: '';
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.045) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.045) 1px, transparent 1px);
  background-size: 34px 34px;
  mask-image: linear-gradient(to bottom, #000, transparent 82%);
}

.logo {
  position: relative;
  z-index: 1;
  min-height: 92px;
  padding: 24px 22px 18px;
  display: flex;
  align-items: center;
  gap: 14px;
}

.logo-mark {
  display: grid;
  width: 46px;
  height: 46px;
  place-items: center;
  color: #142033;
  background: linear-gradient(135deg, #fff6cf, #33dacb);
  border-radius: 14px;
  box-shadow: 0 14px 26px rgba(51, 218, 203, 0.22);
}

.logo-mark .el-icon {
  font-size: 24px;
}

.logo-title {
  font-size: 18px;
  font-weight: 800;
  letter-spacing: 0;
}

.logo-subtitle {
  margin-top: 3px;
  color: #7fded7;
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
}

.el-menu-vertical {
  position: relative;
  z-index: 1;
  flex: 1;
  padding: 8px 14px;
  border-right: none;
}

.el-menu-vertical :deep(.el-menu-item) {
  height: 48px;
  margin: 8px 0;
  border-radius: 12px;
  font-weight: 700;
}

.el-menu-vertical :deep(.el-menu-item:hover) {
  color: #ffffff;
  background: rgba(255, 255, 255, 0.09);
}

.el-menu-vertical :deep(.el-menu-item.is-active) {
  background: linear-gradient(135deg, rgba(51, 218, 203, 0.95), rgba(255, 189, 64, 0.92));
  box-shadow: 0 14px 28px rgba(20, 32, 51, 0.3);
}

.el-menu-vertical :deep(.el-menu-item.is-active .el-icon) {
  color: #132033;
}

.aside-footer {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 16px;
  padding: 14px 16px;
  color: #d8e6f2;
  font-size: 13px;
  font-weight: 700;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 12px;
}

.status-dot,
.pulse {
  width: 9px;
  height: 9px;
  background: #33dacb;
  border-radius: 999px;
  box-shadow: 0 0 0 6px rgba(51, 218, 203, 0.14);
}

.workspace {
  min-width: 0;
}

.header {
  height: 84px;
  background: rgba(255, 255, 255, 0.86);
  border-bottom: 1px solid rgba(206, 217, 229, 0.72);
  display: flex;
  align-items: center;
  padding: 0 34px;
  backdrop-filter: blur(16px);
}

.header-content {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.header-content h3 {
  margin: 0;
  color: #172338;
  font-size: 24px;
  font-weight: 850;
  letter-spacing: 0;
}

.eyebrow {
  margin: 0 0 4px;
  color: #758397;
  font-size: 12px;
  font-weight: 800;
  text-transform: uppercase;
}

.header-badge {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  color: #1e6e67;
  font-size: 13px;
  font-weight: 800;
  background: #e6fbf8;
  border: 1px solid #bbf2eb;
  border-radius: 999px;
}

.main-content {
  min-width: 0;
  background:
    radial-gradient(circle at top right, rgba(51, 218, 203, 0.2), transparent 34%),
    linear-gradient(180deg, #f6f9fc 0%, #edf3f7 100%);
  padding: 30px;
  overflow: auto;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 900px) {
  .layout-container {
    height: auto;
    min-height: 100vh;
  }

  .aside {
    display: none;
  }

  .header {
    height: auto;
    min-height: 74px;
    padding: 18px 20px;
  }

  .header-content {
    align-items: flex-start;
    flex-direction: column;
  }

  .header-badge {
    display: none;
  }

  .main-content {
    padding: 18px;
  }
}
</style>
