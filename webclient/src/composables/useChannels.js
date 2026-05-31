import { ref } from 'vue'
import { request } from '@/api/client'
import { authStore } from '@/stores/authStore.vue'

export function useChannels(usersList) {
  const store = authStore()
  const channelsList = ref([])

  function extractDirectID(name) {
    return name.split('_').find(id => id !== store.meId)
  }

  function extractDirectName(name) {
    const otherId = extractDirectID(name)
    return usersList.value.find(u => u.id === otherId)?.name ?? ''
  }

  async function loadChannels() {
    try {
      const channels = await request('/channels')
      channels.forEach(elem => {
        elem.icon = 'peanut'
        if (elem.type === 'direct') {
          elem.directName = extractDirectName(elem.name)
          elem.url = `${store.API}/users/${extractDirectID(elem.name)}/pic?token=${store.token}`
        } else {
          elem.directName = ''
          elem.url = `${store.API}/channels/${elem.id}/pic?token=${store.token}`
        }
      })
      channelsList.value = channels
    } catch (e) {
      alert(e.message)
    }
  }

  async function addChannel(name, type, imageFile) {
    if (!name) { alert('заполните имя сначала'); return false }
    try {
      const formData = new FormData()
      formData.append('type', type)
      formData.append('name', name)
      if (imageFile) formData.append('pic', imageFile)
      await request('/channels', { method: 'POST', body: formData })
      await loadChannels()
      return true
    } catch (e) {
      alert(e.message)
      return false
    }
  }

  async function addDirect(userId) {
    try {
      await request(`/users/${userId}/direct`, {
        method: 'POST',
        body: JSON.stringify({ abobe: 'obeba' })
      })
      return true
    } catch (e) {
      alert(e.message)
      return false
    }
  }

  return { channelsList, loadChannels, addChannel, addDirect }
}
