<template>
  <v-app>
    <v-main>
      <login v-if="!store.status" />
      <mainpage v-else />
    </v-main>
  </v-app>
</template>

<script setup>
import mainpage from '@/pages/Mainpage.vue'
import login from '@/pages/Login.vue'
import { authStore } from '@/stores/authStore.vue'

const store = authStore()
store.getAuth()
</script>
