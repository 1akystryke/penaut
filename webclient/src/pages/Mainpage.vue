<template>
    <mainMessangerComponent/>
    <v-btn @click="dropAuth">drop auth</v-btn>

</template>
<script setup>
    import { ref, reactive, computed, watch, onMounted } from 'vue'
    import { authStore } from '@/stores/authStore.vue'
    import mainMessangerComponent from '@/components/mainMessangerComponent.vue'
    const store = authStore()
    function dropAuth(){
        store.breakAuth()
    }
</script>