<template>
  <div style="position: relative;">
    <mainMessangerComponent />
    <button class="logout-btn" @click="dropAuth">
      <v-icon size="15">mdi-logout</v-icon>
    </button>
  </div>
</template>

<script setup>
import mainMessangerComponent from '@/components/mainMessangerComponent.vue'
import { authStore } from '@/stores/authStore.vue'

const store = authStore()

function dropAuth() {
  store.breakAuth()
}
</script>

<style scoped>
.logout-btn {
  position: fixed;
  bottom: 20px;
  left: 14px;
  z-index: 100;
  width: 34px;
  height: 34px;
  border-radius: 9px;
  border: none;
  background: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.35);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.logout-btn:hover {
  background: rgba(239, 68, 68, 0.18);
  color: #fca5a5;
}
</style>
