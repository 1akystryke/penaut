<template>
  <v-dialog
    :model-value="modelValue"
    max-width="400"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <div class="dialog-card">
      <div class="dialog-header">
        <div class="dialog-icon">
          <v-icon size="18" color="#a5b4fc">mdi-account-plus-outline</v-icon>
        </div>
        <div>
          <h3 class="dialog-title">Добавить участника</h3>
          <p class="dialog-sub">в канал <strong class="channel-name">{{ channel?.name }}</strong></p>
        </div>
      </div>

      <div class="dialog-body">
        <label class="field-label">Пользователь</label>
        <div class="select-wrap">
          <select v-model="selectedUserId" class="field-select">
            <option value="" disabled selected>Выберите пользователя...</option>
            <option v-for="user in users" :key="user.id" :value="user.id">
              {{ user.name }}
            </option>
          </select>
          <v-icon size="16" class="select-arrow">mdi-chevron-down</v-icon>
        </div>
      </div>

      <div class="dialog-footer">
        <button class="btn-cancel" @click="$emit('update:modelValue', false)">Отмена</button>
        <button class="btn-primary" :disabled="!selectedUserId" @click="submit">Добавить</button>
      </div>
    </div>
  </v-dialog>
</template>

<script setup>
import { ref } from 'vue'

defineProps({ modelValue: Boolean, channel: Object, users: Array })
const emit = defineEmits(['update:modelValue', 'submit'])

const selectedUserId = ref('')

function submit() {
  if (selectedUserId.value) emit('submit', selectedUserId.value)
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

.channel-name {
  color: #a5b4fc;
  font-weight: 600;
}

.dialog-body {
  padding: 18px 22px;
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

.select-wrap {
  position: relative;
}

.field-select {
  width: 100%;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 9px;
  padding: 10px 36px 10px 13px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.88);
  font-family: inherit;
  outline: none;
  appearance: none;
  cursor: pointer;
  transition: border-color 0.14s;
  box-sizing: border-box;
}

.field-select:focus {
  border-color: rgba(99, 102, 241, 0.5);
}

.field-select option {
  background: #1e1e30;
  color: rgba(255, 255, 255, 0.88);
}

.select-arrow {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  color: rgba(255, 255, 255, 0.35);
  pointer-events: none;
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

.btn-primary:not(:disabled):hover {
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.48);
  transform: translateY(-1px);
}

.btn-primary:disabled {
  opacity: 0.45;
  cursor: default;
  box-shadow: none;
}
</style>
