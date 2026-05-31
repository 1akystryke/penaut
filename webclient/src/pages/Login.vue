<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-logo">
        <div class="logo-icon">
          <v-icon size="24" color="white">mdi-peanut</v-icon>
        </div>
        <span class="logo-text">Peanut</span>
      </div>

      <h1 class="login-title">Добро пожаловать</h1>
      <p class="login-sub">Войдите в свою учётную запись</p>

      <div class="form">
        <div class="field-group">
          <label class="field-label">Email</label>
          <div class="field-wrap" :class="{ focused: emailFocused }">
            <v-icon size="16" class="field-icon">mdi-email-outline</v-icon>
            <input
              v-model="email"
              type="email"
              class="field-input"
              placeholder="you@example.com"
              @focus="emailFocused = true"
              @blur="emailFocused = false"
              @keyup.enter="tryLogin"
            />
          </div>
        </div>

        <div class="field-group">
          <label class="field-label">Пароль</label>
          <div class="field-wrap" :class="{ focused: pwFocused }">
            <v-icon size="16" class="field-icon">mdi-lock-outline</v-icon>
            <input
              v-model="password"
              :type="visible ? 'text' : 'password'"
              class="field-input"
              placeholder="••••••••"
              @focus="pwFocused = true"
              @blur="pwFocused = false"
              @keyup.enter="tryLogin"
            />
            <button class="toggle-btn" @click="visible = !visible">
              <v-icon size="16">{{ visible ? 'mdi-eye-off' : 'mdi-eye' }}</v-icon>
            </button>
          </div>
        </div>

        <button class="login-btn" :disabled="loading" @click="tryLogin">
          <span v-if="!loading">Войти</span>
          <v-progress-circular v-else size="18" width="2" indeterminate color="white" />
        </button>
      </div>
    </div>

    <div class="login-bg">
      <div class="bg-blob bg-blob-1" />
      <div class="bg-blob bg-blob-2" />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { authStore } from '@/stores/authStore.vue'

const store = authStore()
const DEFAULT_HOST = 'localhost:8080'
const API = `http://${DEFAULT_HOST}`

const email = ref('')
const password = ref('')
const visible = ref(false)
const loading = ref(false)
const emailFocused = ref(false)
const pwFocused = ref(false)

async function hashPassword(pwd) {
  const encoder = new TextEncoder()
  const data = encoder.encode(pwd)
  const hashBuffer = await crypto.subtle.digest('SHA-256', data)
  const hashArray = Array.from(new Uint8Array(hashBuffer))
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
}

async function tryLogin() {
  if (!email.value || !password.value) return
  try {
    loading.value = true
    const hash = await hashPassword(password.value)
    const r = await fetch(API + '/auth', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: email.value, pwd: hash })
    })
    const result = await r.json().catch(() => ({}))
    if (!r.ok) throw new Error(result.error || 'Ошибка входа')
    localStorage.setItem('token', result.token)
    localStorage.setItem('meName', result.name)
    localStorage.setItem('meId', result.id)
    store.getAuth()
  } catch (e) {
    alert(e.message)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #0a0a14;
  position: relative;
  overflow: hidden;
}

.login-card {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 400px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 20px;
  padding: 36px 36px 32px;
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.4);
}

.login-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 28px;
}

.logo-icon {
  width: 38px;
  height: 38px;
  border-radius: 12px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.45);
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.92);
  letter-spacing: 0.01em;
}

.login-title {
  font-size: 22px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.94);
  margin: 0 0 6px;
}

.login-sub {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.38);
  margin: 0 0 28px;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-size: 12px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.45);
  letter-spacing: 0.04em;
}

.field-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0 12px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  transition: border-color 0.15s ease;
}

.field-wrap.focused {
  border-color: rgba(99, 102, 241, 0.55);
  background: rgba(99, 102, 241, 0.05);
}

.field-icon {
  color: rgba(255, 255, 255, 0.28);
  flex-shrink: 0;
}

.field-wrap.focused .field-icon {
  color: rgba(165, 180, 252, 0.7);
}

.field-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  padding: 12px 0;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.88);
  font-family: inherit;
}

.field-input::placeholder {
  color: rgba(255, 255, 255, 0.2);
}

.toggle-btn {
  border: none;
  background: transparent;
  cursor: pointer;
  padding: 0;
  color: rgba(255, 255, 255, 0.28);
  display: flex;
  align-items: center;
  flex-shrink: 0;
  transition: color 0.12s;
}

.toggle-btn:hover {
  color: rgba(255, 255, 255, 0.6);
}

.login-btn {
  width: 100%;
  padding: 13px;
  margin-top: 6px;
  border-radius: 11px;
  border: none;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #ffffff;
  font-size: 15px;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.18s ease;
  box-shadow: 0 3px 12px rgba(99, 102, 241, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 48px;
}

.login-btn:not(:disabled):hover {
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.52);
  transform: translateY(-1px);
}

.login-btn:disabled {
  opacity: 0.7;
  cursor: default;
}

/* Background decoration */
.login-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.bg-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.25;
}

.bg-blob-1 {
  width: 400px;
  height: 400px;
  background: #6366f1;
  top: -100px;
  right: -100px;
}

.bg-blob-2 {
  width: 300px;
  height: 300px;
  background: #8b5cf6;
  bottom: -80px;
  left: -60px;
}
</style>
