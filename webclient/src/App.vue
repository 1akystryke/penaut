<template>
  <v-app>
    <v-main>
      
      <login v-if="!store.status"/>
      <mainpage v-else/>
    </v-main>
    <v-btn
      class="ma-2"
      icon="mdi-theme-light-dark"
      location="top right"
      position="absolute"
      @click="$vuetify.theme.cycle()"
    />
  </v-app>
</template>

<script setup>
  import { ref, reactive, computed, watch, onMounted } from 'vue'
  //import HelloWorld from '@/components/HelloWorld.vue'
  import mainpage from '@/pages/Mainpage.vue'
  import login from '@/pages/Login.vue'
  import { authStore } from '@/stores/authStore.vue'
  const store = authStore()
  store.getAuth()
</script>

