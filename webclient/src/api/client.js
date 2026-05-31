import { authStore } from '@/stores/authStore.vue'

export async function request(path, options = {}) {
  const store = authStore()
  const headers = { ...(options.headers || {}) }
  if (store.status) {
    headers.Authorization = `Bearer ${store.token}`
  }
  const r = await fetch(store.API + path, { ...options, headers })
  const data = await r.json().catch(() => ({}))
  if (r.status === 401) {
    store.breakAuth()
  }
  if (!r.ok) {
    throw new Error(data.error || 'Request failed')
  }
  return data
}

export async function requestBlob(path) {
  const store = authStore()
  const headers = {}
  if (store.status) {
    headers.Authorization = `Bearer ${store.token}`
  }
  const r = await fetch(store.API + path, { headers })
  if (r.status === 401) {
    store.breakAuth()
  }
  if (!r.ok) {
    throw new Error('Request failed')
  }
  return r.blob()
}
