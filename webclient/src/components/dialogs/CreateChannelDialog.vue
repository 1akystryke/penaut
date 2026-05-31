<template>
  <v-dialog
    :model-value="modelValue"
    max-width="440"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <div class="dialog-card">
      <div class="dialog-header">
        <div class="dialog-icon">
          <v-icon size="18" color="#a5b4fc">mdi-pound</v-icon>
        </div>
        <div>
          <h3 class="dialog-title">Создать канал</h3>
          <p class="dialog-sub">Настройте новый канал для команды</p>
        </div>
      </div>

      <div class="dialog-body">
        <label class="field-label">Название</label>
        <input v-model="name" class="field-input" placeholder="например, общий" />

        <label class="field-label" style="margin-top: 16px;">Тип</label>
        <div class="type-toggle">
          <button
            v-for="t in channelTypes"
            :key="t.value"
            class="type-btn"
            :class="{ active: type === t.value }"
            @click="type = t.value"
          >
            <v-icon size="15">{{ t.icon }}</v-icon>
            {{ t.label }}
          </button>
        </div>

        <label class="field-label" style="margin-top: 16px;">Изображение</label>
        <div v-if="!imagePreviewUrl" class="upload-zone" @click="triggerFileInput">
          <v-icon size="24" color="rgba(165,180,252,0.5)">mdi-image-plus</v-icon>
          <span class="upload-hint">Нажмите для выбора</span>
          <input ref="imgInput" type="file" accept="image/*" class="hidden-input" @change="onFileChange" />
        </div>
        <div v-else class="preview-card">
          <img :src="imagePreviewUrl" class="preview-img" />
          <div class="preview-info">
            <span class="preview-name">{{ selectedFileName }}</span>
            <span class="preview-size">{{ formatFileSize(selectedFileSize) }}</span>
          </div>
          <div class="preview-actions">
            <button class="preview-btn" @click="replaceImage">
              <v-icon size="14">mdi-reload</v-icon>
            </button>
            <button class="preview-btn preview-btn--danger" @click="removeImage">
              <v-icon size="14">mdi-delete</v-icon>
            </button>
          </div>
        </div>
      </div>

      <div class="dialog-footer">
        <button class="btn-cancel" @click="$emit('update:modelValue', false)">Отмена</button>
        <button class="btn-primary" @click="submit">Создать канал</button>
      </div>
    </div>
  </v-dialog>
</template>

<script setup>
import { ref } from 'vue'

defineProps({ modelValue: Boolean })
const emit = defineEmits(['update:modelValue', 'create'])

const channelTypes = [
  { value: 'public', label: 'Публичный', icon: 'mdi-earth' },
  { value: 'private', label: 'Приватный', icon: 'mdi-lock-outline' },
]

const name = ref('')
const type = ref('public')
const imageFile = ref(null)
const imagePreviewUrl = ref(null)
const selectedFileName = ref('')
const selectedFileSize = ref(0)
const imgInput = ref(null)

function triggerFileInput() {
  imgInput.value?.click()
}

function onFileChange(e) {
  const file = e.target.files?.[0]
  if (!file || !file.type.startsWith('image/')) return
  imageFile.value = file
  selectedFileName.value = file.name
  selectedFileSize.value = file.size
  if (imagePreviewUrl.value) URL.revokeObjectURL(imagePreviewUrl.value)
  imagePreviewUrl.value = URL.createObjectURL(file)
}

function replaceImage() {
  imagePreviewUrl.value = null
  imageFile.value = null
  triggerFileInput()
}

function removeImage() {
  if (imagePreviewUrl.value) URL.revokeObjectURL(imagePreviewUrl.value)
  imagePreviewUrl.value = null
  imageFile.value = null
  selectedFileName.value = ''
  selectedFileSize.value = 0
}

function formatFileSize(bytes) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

function submit() {
  emit('create', { name: name.value, type: type.value, imageFile: imageFile.value })
}
</script>

<style scoped>
.dialog-card {
  background: #13131f;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 18px;
  overflow: hidden;
}

.dialog-header {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 22px 22px 0;
}

.dialog-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: rgba(99, 102, 241, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.dialog-title {
  font-size: 16px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.92);
  margin: 0;
}

.dialog-sub {
  font-size: 12.5px;
  color: rgba(255, 255, 255, 0.38);
  margin: 2px 0 0;
}

.dialog-body {
  padding: 18px 22px;
  display: flex;
  flex-direction: column;
}

.field-label {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: rgba(255, 255, 255, 0.38);
  margin-bottom: 7px;
  display: block;
}

.field-input {
  width: 100%;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 9px;
  padding: 10px 13px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.88);
  font-family: inherit;
  outline: none;
  transition: border-color 0.14s;
  box-sizing: border-box;
}

.field-input::placeholder {
  color: rgba(255, 255, 255, 0.22);
}

.field-input:focus {
  border-color: rgba(99, 102, 241, 0.5);
}

.type-toggle {
  display: flex;
  gap: 6px;
}

.type-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 12px;
  border-radius: 9px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.45);
  font-size: 13px;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.14s ease;
}

.type-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: rgba(255, 255, 255, 0.72);
}

.type-btn.active {
  background: rgba(99, 102, 241, 0.2);
  border-color: rgba(99, 102, 241, 0.4);
  color: #a5b4fc;
  font-weight: 600;
}

.upload-zone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px dashed rgba(255, 255, 255, 0.12);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.14s ease;
}

.upload-zone:hover {
  background: rgba(99, 102, 241, 0.06);
  border-color: rgba(99, 102, 241, 0.3);
}

.upload-hint {
  font-size: 12.5px;
  color: rgba(255, 255, 255, 0.3);
}

.hidden-input {
  display: none;
}

.preview-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
}

.preview-img {
  width: 44px;
  height: 44px;
  border-radius: 8px;
  object-fit: cover;
  flex-shrink: 0;
}

.preview-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.preview-name {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.preview-size {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.32);
}

.preview-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.preview-btn {
  width: 28px;
  height: 28px;
  border-radius: 7px;
  border: none;
  background: rgba(255, 255, 255, 0.07);
  color: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.12s ease;
}

.preview-btn:hover {
  background: rgba(255, 255, 255, 0.12);
  color: rgba(255, 255, 255, 0.85);
}

.preview-btn--danger:hover {
  background: rgba(239, 68, 68, 0.18);
  color: #fca5a5;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 0 22px 20px;
}

.btn-cancel {
  padding: 9px 18px;
  border-radius: 9px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: transparent;
  color: rgba(255, 255, 255, 0.5);
  font-size: 13.5px;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.14s;
}

.btn-cancel:hover {
  background: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.75);
}

.btn-primary {
  padding: 9px 18px;
  border-radius: 9px;
  border: none;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #ffffff;
  font-size: 13.5px;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.16s ease;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.35);
}

.btn-primary:hover {
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.48);
  transform: translateY(-1px);
}
</style>
