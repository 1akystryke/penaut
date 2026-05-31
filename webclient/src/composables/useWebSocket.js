import { authStore } from '@/stores/authStore.vue'
import { useAttachments } from './useAttachments'
import { parseCustomDate } from '@/utils/date'

export function useWebSocket(messages, scrollToBottom) {
  const store = authStore()
  const { hydrateMessageAttachments } = useAttachments()
  let socket = null

  function subscribe(channelId) {
    if (socket) socket.close()
    socket = new WebSocket(`${store.WS}?channel_id=${channelId}&token=${store.token}`)

    socket.onmessage = (e) => {
      const data = JSON.parse(e.data)
      if (data.error) { alert(data.error); return }

      data.time = parseCustomDate(data.created_at)
      data.sender = data.user_name
      data.avatar = 'peanut-outline'
      data.isMine = data.author === store.meId

      hydrateMessageAttachments(data).then(() => {
        messages.value.push(data)
        scrollToBottom()
      })
    }

    socket.onerror = () => console.error('WebSocket error')
  }

  function isConnected() {
    return socket !== null && socket.readyState === WebSocket.OPEN
  }

  function close() {
    if (socket) socket.close()
  }

  return { subscribe, isConnected, close }
}
