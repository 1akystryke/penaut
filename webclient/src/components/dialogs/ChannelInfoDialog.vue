<template>
  <v-dialog
    :model-value="modelValue"
    max-width="400"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <div class="dialog-card">
      <div class="dialog-header">
        <div class="dialog-icon">
          <v-icon size="18" color="#a5b4fc">mdi-information-outline</v-icon>
        </div>
        <div class="header-text">
          <h3 class="dialog-title">{{ channel?.name }}</h3>
          <span class="type-badge" :class="`type-badge--${channel?.type}`">{{ channel?.type }}</span>
        </div>
      </div>

      <div class="dialog-body">
        <div class="meta-row">
          <v-icon size="14" class="meta-icon">mdi-identifier</v-icon>
          <span class="meta-value">{{ channel?.id }}</span>
        </div>

        <div class="members-section">
          <div class="members-header">
            <span class="section-label">Участники</span>
            <button
              v-if="channel?.type !== 'direct'"
              class="add-member-btn"
              @click="$emit('addUser')"
            >
              <v-icon size="13">mdi-plus</v-icon>
              Добавить
            </button>
          </div>

          <div class="members-list">
            <div v-for="member in members" :key="member.id" class="member-row">
              <div class="member-avatar">
                <v-icon size="14" color="rgba(165,180,252,0.7)">mdi-account</v-icon>
              </div>
              <span class="member-name">{{ member.name }}</span>
            </div>
            <p v-if="!members?.length" class="no-members">Нет участников</p>
          </div>
        </div>
      </div>

      <div class="dialog-footer">
        <button class="btn-close" @click="$emit('update:modelValue', false)">Закрыть</button>
      </div>
    </div>
  </v-dialog>
</template>

<script setup>
defineProps({ modelValue: Boolean, channel: Object, members: Array })
defineEmits(['update:modelValue', 'addUser'])
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

.header-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.dialog-title {
  font-size: 16px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.92);
  margin: 0;
}

.type-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  width: fit-content;
}

.type-badge--public {
  background: rgba(52, 211, 153, 0.14);
  color: #6ee7b7;
}

.type-badge--private {
  background: rgba(251, 191, 36, 0.14);
  color: #fcd34d;
}

.type-badge--direct {
  background: rgba(99, 102, 241, 0.15);
  color: #a5b4fc;
}

.dialog-body {
  padding: 18px 22px;
}

.meta-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.04);
  border-radius: 8px;
  margin-bottom: 16px;
}

.meta-icon {
  color: rgba(255, 255, 255, 0.28);
  flex-shrink: 0;
}

.meta-value {
  font-size: 12.5px;
  color: rgba(255, 255, 255, 0.45);
  font-family: 'Courier New', monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.members-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.members-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.section-label {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: rgba(255, 255, 255, 0.32);
}

.add-member-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 7px;
  border: 1px solid rgba(99, 102, 241, 0.3);
  background: rgba(99, 102, 241, 0.1);
  color: #a5b4fc;
  font-size: 12px;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  transition: all 0.14s;
}

.add-member-btn:hover {
  background: rgba(99, 102, 241, 0.2);
  border-color: rgba(99, 102, 241, 0.5);
}

.members-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 180px;
  overflow-y: auto;
}

.members-list::-webkit-scrollbar { width: 3px; }
.members-list::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.08);
  border-radius: 2px;
}

.member-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 7px 10px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.03);
}

.member-avatar {
  width: 26px;
  height: 26px;
  border-radius: 7px;
  background: rgba(99, 102, 241, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.member-name {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.72);
}

.no-members {
  font-size: 12.5px;
  color: rgba(255, 255, 255, 0.22);
  text-align: center;
  padding: 12px;
  margin: 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  padding: 0 22px 20px;
}

.btn-close {
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

.btn-close:hover {
  background: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.75);
}
</style>
