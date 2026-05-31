import { ref } from 'vue'
import { request } from '@/api/client'
import { authStore } from '@/stores/authStore.vue'
import { useAttachments } from './useAttachments'
import { parseCustomDate } from '@/utils/date'

export function useMessages() {
  const store = authStore()
  const { hydrateMessageAttachments } = useAttachments()

  const messages = ref([])
  const isSendingMessage = ref(false)

  function normalizePost(post) {
    post.time = parseCustomDate(post.created_at)
    post.sender = post.user_name
    post.avatar = 'peanut-outline'
    post.isMine = post.author === store.meId
  }

  async function getChannelMessages(channelId) {
    try {
      const posts = await request(`/channels/${channelId}/posts`)
      await Promise.all(posts.map(async post => {
        normalizePost(post)
        await hydrateMessageAttachments(post)
      }))
      messages.value = posts
    } catch (e) {
      alert(e.message)
    }
  }

  async function sendMessage(channelId, text, files) {
    if (!text) return null
    try {
      isSendingMessage.value = true
      const formData = new FormData()
      formData.append('text', text)
      files.forEach(file => formData.append('attachments[]', file))
      return await request(`/channels/${channelId}/posts`, { method: 'POST', body: formData })
    } catch (e) {
      alert(e.message)
      return null
    } finally {
      isSendingMessage.value = false
    }
  }

  return { messages, isSendingMessage, getChannelMessages, sendMessage, normalizePost }
}
