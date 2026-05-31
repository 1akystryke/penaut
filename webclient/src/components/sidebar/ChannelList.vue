<template>
  <div class="section">
    <div class="section-header">
      <span class="section-label">Каналы</span>
      <button class="add-btn" title="Создать канал" @click="$emit('create')">
        <v-icon size="15">mdi-plus</v-icon>
      </button>
    </div>

    <div class="search-bar">
      <v-icon size="14" class="search-icon">mdi-magnify</v-icon>
      <input v-model="searchQuery" class="search-input" placeholder="Поиск..." />
    </div>

    <div class="item-list">
      <button
        v-for="channel in filtered"
        :key="channel.id"
        class="list-item"
        :class="{ active: activeChannel === channel.id }"
        @click="$emit('select', channel.id)"
      >
        <span class="active-bar" />
        <v-avatar size="33" rounded="md" class="item-avatar">
          <v-img :src="channel.url" />
        </v-avatar>
        <div class="item-text">
          <span class="item-name">{{ channel.name }}</span>
          <span class="item-meta">{{ channel.type }}</span>
        </div>
      </button>

      <p v-if="filtered.length === 0" class="empty-hint">Каналы не найдены</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({ channels: Array, activeChannel: String })
defineEmits(['select', 'create'])

const searchQuery = ref('')

const filtered = computed(() => {
  if (!searchQuery.value) return props.channels
  return props.channels.filter(c =>
    c.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})
</script>

<style scoped>
.section {
  display: flex;
  flex-direction: column;
  min-height: 0;
  flex: 1;
  overflow: hidden;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px 6px;
  flex-shrink: 0;
}

.section-label {
  font-size: 10.5px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.09em;
  color: rgba(255, 255, 255, 0.32);
}

.add-btn {
  width: 22px;
  height: 22px;
  border-radius: 6px;
  border: none;
  background: rgba(255, 255, 255, 0.07);
  color: rgba(255, 255, 255, 0.4);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.14s ease;
}

.add-btn:hover {
  background: rgba(99, 102, 241, 0.28);
  color: #a5b4fc;
}

.search-bar {
  display: flex;
  align-items: center;
  gap: 7px;
  margin: 0 10px 7px;
  padding: 6px 10px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.05);
  flex-shrink: 0;
}

.search-icon {
  color: rgba(255, 255, 255, 0.28);
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: rgba(255, 255, 255, 0.78);
  font-size: 12.5px;
  font-family: inherit;
}

.search-input::placeholder {
  color: rgba(255, 255, 255, 0.22);
}

.item-list {
  overflow-y: auto;
  flex: 1;
  padding: 0 7px 14px;
}

.item-list::-webkit-scrollbar { width: 3px; }
.item-list::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.08);
  border-radius: 2px;
}

.list-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 6px 9px 6px 10px;
  border-radius: 9px;
  border: none;
  background: transparent;
  cursor: pointer;
  text-align: left;
  transition: background 0.12s ease;
  position: relative;
}

.list-item:hover { background: rgba(255, 255, 255, 0.06); }

.list-item.active { background: rgba(99, 102, 241, 0.16); }

.active-bar {
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%) scaleY(0);
  width: 3px;
  height: 58%;
  background: #6366f1;
  border-radius: 0 2px 2px 0;
  transition: transform 0.15s ease;
}

.list-item.active .active-bar {
  transform: translateY(-50%) scaleY(1);
}

.item-avatar {
  flex-shrink: 0;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.28);
}

.item-text {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.item-name {
  font-size: 13px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.72);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.3;
  transition: color 0.12s;
}

.list-item.active .item-name {
  color: #c7d2fe;
  font-weight: 600;
}

.item-meta {
  font-size: 10.5px;
  color: rgba(255, 255, 255, 0.28);
  line-height: 1.2;
}

.list-item.active .item-meta {
  color: rgba(165, 180, 252, 0.55);
}

.empty-hint {
  text-align: center;
  padding: 20px 16px;
  font-size: 12.5px;
  color: rgba(255, 255, 255, 0.22);
  margin: 0;
}
</style>
