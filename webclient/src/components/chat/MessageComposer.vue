<template>
  <div class="composer-wrap">
    <div class="composer-card">
      <div v-if="files.length" class="file-list">
        <div
          v-for="(file, index) in files"
          :key="file.name + file.size + index"
          class="file-chip"
        >
          <v-icon size="12">mdi-paperclip</v-icon>
          <span class="file-name">{{ file.name }}</span>
          <button class="file-remove" @click="removeFile(index)">
            <v-icon size="11">mdi-close</v-icon>
          </button>
        </div>
      </div>

      <div class="composer-row">
        <input ref="fileInput" type="file" multiple class="file-input" @change="handleFilesSelect" />
        <button class="icon-btn" title="Прикрепить файл" @click="fileInput?.click()">
          <v-icon size="17">mdi-paperclip</v-icon>
        </button>
        <input
          v-model="text"
          class="text-input"
          placeholder="Написать сообщение..."
          @keyup.enter="submit"
        />
        <button
          class="send-btn"
          :class="{ 'send-btn--active': text.trim() }"
          :disabled="isSending || !text.trim()"
          @click="submit"
        >
          <v-icon size="16">mdi-send</v-icon>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

defineProps({ isSending: Boolean })
const emit = defineEmits(['send'])

const text = ref('')
const files = ref([])
const fileInput = ref(null)

function handleFilesSelect(event) {
  files.value = [...files.value, ...Array.from(event.target.files || [])]
  event.target.value = ''
}

function removeFile(index) {
  files.value.splice(index, 1)
}

function submit() {
  if (!text.value.trim()) return
  emit('send', { text: text.value.trim(), files: [...files.value] })
  text.value = ''
  files.value = []
  if (fileInput.value) fileInput.value.value = ''
}
</script>

<style scoped>
.composer-wrap {
  padding: 8px 20px 18px;
  flex-shrink: 0;
}

.composer-card {
  background: #ffffff;
  border-radius: 14px;
  border: 1px solid rgba(0, 0, 0, 0.07);
  box-shadow: 0 2px 16px rgba(0, 0, 0, 0.07);
  overflow: hidden;
}

.file-list {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  padding: 10px 12px 4px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.file-chip {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 3px 8px 3px 7px;
  background: rgba(99, 102, 241, 0.08);
  border-radius: 20px;
  font-size: 12px;
  color: #6366f1;
  max-width: 180px;
}

.file-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-remove {
  border: none;
  background: transparent;
  cursor: pointer;
  padding: 0;
  color: #6366f1;
  display: flex;
  align-items: center;
  opacity: 0.6;
  flex-shrink: 0;
  transition: opacity 0.12s;
}

.file-remove:hover { opacity: 1; }

.composer-row {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 5px 6px;
}

.file-input {
  display: none;
}

.icon-btn {
  width: 34px;
  height: 34px;
  border-radius: 9px;
  border: none;
  background: transparent;
  color: #b0b8c8;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.14s ease;
  flex-shrink: 0;
}

.icon-btn:hover {
  background: rgba(99, 102, 241, 0.08);
  color: #6366f1;
}

.text-input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 14px;
  color: #1f2937;
  background: transparent;
  padding: 8px 6px;
  font-family: inherit;
  line-height: 1.4;
}

.text-input::placeholder {
  color: #b0b8c8;
}

.send-btn {
  width: 34px;
  height: 34px;
  border-radius: 9px;
  border: none;
  background: #f3f4f6;
  color: #c4cdd6;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.18s ease;
  flex-shrink: 0;
}

.send-btn--active {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #ffffff;
  box-shadow: 0 2px 10px rgba(99, 102, 241, 0.38);
}

.send-btn--active:not(:disabled):hover {
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.48);
  transform: scale(1.06);
}

.send-btn:disabled {
  cursor: default;
}
</style>
