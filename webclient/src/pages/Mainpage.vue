<template>
    <mainMessangerComponent/>
    <v-btn @click="dropAuth">Logout</v-btn>
    <v-img :src="imageUrl" height="200" class="mb-4"></v-img>
</template>
<script setup>
    import { ref, reactive, computed, watch, onMounted } from 'vue'
    import mainMessangerComponent from '@/components/mainMessangerComponent.vue'
    import { authStore } from '@/stores/authStore.vue'
    const store = authStore()
    function dropAuth(){
        store.breakAuth()
    }

    const imageUrl = ref('')

    const loadImage = async () => {
    try {
        const response = await fetch(store.API + `/users/` + store.meId + `/pic`, {
        headers: {
            'Authorization': `Bearer ${store.token}`, // ваш заголовок
            // или другие кастомные заголовки
        }
        })
        
        if (!response.ok) throw new Error('Ошибка загрузки')
        
        const blob = await response.blob()
        imageUrl.value = URL.createObjectURL(blob)
    } catch (error) {
        console.error('Не удалось загрузить аватар:', error)
        imageUrl.value = '' // или путь к заглушке
    }
    }

    loadImage()
</script>