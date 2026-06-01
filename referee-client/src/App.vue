<template>
  <div class="app-container">
    <header>
      <h1>网球裁判系统</h1>
    </header>
    <main>
      <router-view></router-view>
    </main>
  </div>
</template>

<script setup>
</script>

<style>
body {
  margin: 0;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f0f2f5;
  color: #333;
}

.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

header {
  background-color: #2c3e50;
  color: white;
  padding: 10px 20px;
  text-align: center;
}

main {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}
</style>
