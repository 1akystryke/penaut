<template>
  <v-app>
    <v-navigation-drawer
      v-model="drawer"
      width="320"
      permanent
      class="channel-drawer"
    >
      <!-- Заголовок панели каналов -->
      <v-list-item
        class="pa-4"
        title="Каналы"
        prepend-icon="mdi-message-outline"
      >
        <template #append>
          <v-btn
            icon="mdi-plus"
            variant="text"
            size="small"
            @click="addChannel"
          />
        </template>
      </v-list-item>

      <v-divider />

      <!-- Поиск каналов -->
      <div class="pa-2">
        <v-text-field
          v-model="searchQuery"
          density="compact"
          variant="outlined"
          placeholder="Поиск каналов..."
          prepend-inner-icon="mdi-magnify"
          hide-details
          clearable
        />
      </div>

      <!-- Список каналов -->
      <v-list nav density="compact">
        <v-list-item
          v-for="channel in filteredChannels"
          :key="channel.id"
          :value="channel.id"
          :title="channel.name"
          :subtitle="channel.lastMessage"
          :active="activeChannel === channel.id"
          @click="selectChannel(channel.id)"
          color="primary"
        >
          <template #prepend>
            <v-avatar size="40" color="surface-variant">
              <v-icon>{{ channel.icon }}</v-icon>
            </v-avatar>
          </template>

          <template #append>
            <div class="text-caption text-medium-emphasis">
              {{ channel.time }}
            </div>
          </template>
        </v-list-item>
      </v-list>

      <!-- Если каналов нет -->
      <div
        v-if="filteredChannels.length === 0"
        class="text-center pa-4 text-medium-emphasis"
      >
        Каналы не найдены
      </div>
    </v-navigation-drawer>

    <!-- Основная область сообщений -->
    <v-main class="chat-main">
      <!-- Верхняя панель чата -->
      <v-app-bar
        v-if="activeChannelData"
        flat
        border="bottom"
        class="chat-header"
      >
        <template #prepend>
          <v-avatar size="36" color="surface-variant" class="mr-2">
            <v-icon>{{ activeChannelData.icon }}</v-icon>
          </v-avatar>
          <div>
            <div class="text-h6">{{ activeChannelData.name }}</div>
            <div class="text-caption text-medium-emphasis">
              {{ activeChannelData.members }} участников
            </div>
          </div>
        </template>

        <template #append>
          <v-btn icon="mdi-phone" variant="text" />
          <v-btn icon="mdi-video" variant="text" />
          <v-btn icon="mdi-information-outline" variant="text" />
        </template>
      </v-app-bar>

      <!-- Область сообщений -->
      <div class="messages-container" ref="messagesContainer">
        <div v-if="!activeChannel" class="d-flex align-center justify-center fill-height">
          <div class="text-center text-medium-emphasis">
            <v-icon size="64" class="mb-2">mdi-message-text-outline</v-icon>
            <div class="text-h6">Выберите канал для общения</div>
          </div>
        </div>

        <div v-else class="pa-4">
          <div
            v-for="message in activeMessages"
            :key="message.id"
            :class="['message-wrapper', message.isMine ? 'message-mine' : 'message-other']"
          >
            <!-- Аватар (только для чужих сообщений) -->
            <v-avatar
              v-if="!message.isMine"
              size="32"
              class="mr-2"
              color="surface-variant"
            >
              <v-icon size="18">{{ message.avatar }}</v-icon>
            </v-avatar>

            <div>
              <!-- Имя отправителя -->
              <div
                v-if="!message.isMine"
                class="text-caption text-medium-emphasis mb-1"
              >
                {{ message.sender }}
              </div>

              <!-- Сообщение -->
              <div class="message-bubble">
                {{ message.text }}
                <div class="text-caption text-disabled text-right mt-1">
                  {{ message.time }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Поле ввода сообщения -->
      <v-footer
        v-if="activeChannel"
        app
        class="chat-footer pa-2"
        border="top"
      >
        <v-text-field
          v-model="newMessage"
          variant="outlined"
          density="compact"
          placeholder="Введите сообщение..."
          hide-details
          @keyup.enter="sendMessage"
        >
          <template #prepend-inner>
            <v-btn icon="mdi-emoticon-outline" variant="text" size="small" />
            <v-btn icon="mdi-paperclip" variant="text" size="small" />
          </template>
        </v-text-field>
        <v-btn
          icon="mdi-send"
          color="primary"
          variant="text"
          class="ml-2"
          @click="sendMessage"
          :disabled="!newMessage.trim()"
        />
      </v-footer>
    </v-main>
  </v-app>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue';

const drawer = ref(true);
const activeChannel = ref(null);
const searchQuery = ref('');
const newMessage = ref('');
const messagesContainer = ref(null);

// Тестовые данные каналов
const channels = ref([
  {
    id: 1,
    name: 'Общий',
    icon: 'mdi-pound',
    lastMessage: 'Привет всем!',
    time: '12:30',
    members: 128
  },
  {
    id: 2,
    name: 'Разработка',
    icon: 'mdi-code-tags',
    lastMessage: 'Нужно поправить баг в продакшене',
    time: '11:45',
    members: 45
  },
  {
    id: 3,
    name: 'Дизайн',
    icon: 'mdi-palette-outline',
    lastMessage: 'Новые макеты готовы',
    time: '10:20',
    members: 23
  },
  {
    id: 4,
    name: 'Случайное',
    icon: 'mdi-dice-5',
    lastMessage: 'Кто смотрел новый фильм?',
    time: '09:15',
    members: 56
  }
]);

// Тестовые сообщения для каналов
const messages = ref({
  1: [
    { id: 1, text: 'Всем привет!', sender: 'Анна', avatar: 'mdi-account', time: '12:00', isMine: false },
    { id: 2, text: 'Привет! Как дела?', sender: 'Вы', time: '12:05', isMine: true },
    { id: 3, text: 'Отлично! Работаю над проектом', sender: 'Анна', avatar: 'mdi-account', time: '12:10', isMine: false },
    { id: 4, text: 'Здорово, я тоже', sender: 'Вы', time: '12:15', isMine: true }
  ],
  2: [
    { id: 1, text: 'Кто-нибудь может помочь с багом?', sender: 'Петр', avatar: 'mdi-account', time: '11:30', isMine: false },
    { id: 2, text: 'Да, конечно. Что за проблема?', sender: 'Вы', time: '11:35', isMine: true }
  ],
  3: [
    { id: 1, text: 'Новые макеты в Figma', sender: 'Мария', avatar: 'mdi-account', time: '10:00', isMine: false }
  ],
  4: []
});

// Фильтрация каналов по поиску
const filteredChannels = computed(() => {
  if (!searchQuery.value) return channels.value;
  return channels.value.filter(channel =>
    channel.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  );
});

// Активный канал
const activeChannelData = computed(() => {
  return channels.value.find(c => c.id === activeChannel.value);
});

// Сообщения активного канала
const activeMessages = computed(() => {
  if (!activeChannel.value) return [];
  return messages.value[activeChannel.value] || [];
});

// Выбор канала
function selectChannel(channelId) {
  activeChannel.value = channelId;
  scrollToBottom();
}

// Отправка сообщения
function sendMessage() {
  if (!newMessage.value.trim() || !activeChannel.value) return;

  const now = new Date();
  const time = `${now.getHours()}:${String(now.getMinutes()).padStart(2, '0')}`;

  messages.value[activeChannel.value].push({
    id: Date.now(),
    text: newMessage.value,
    sender: 'Вы',
    time: time,
    isMine: true
  });

  // Обновляем последнее сообщение в канале
  const channel = channels.value.find(c => c.id === activeChannel.value);
  if (channel) {
    channel.lastMessage = newMessage.value;
    channel.time = time;
  }

  newMessage.value = '';
  scrollToBottom();
}

// Функция для добавления канала (заглушка)
function addChannel() {
  const name = prompt('Название канала:');
  if (name) {
    const newChannel = {
      id: Date.now(),
      name: name,
      icon: 'mdi-pound',
      lastMessage: 'Канал создан',
      time: new Date().toLocaleTimeString().slice(0, 5),
      members: 1
    };
    channels.value.push(newChannel);
    messages.value[newChannel.id] = [];
  }
}

// Прокрутка вниз
function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
    }
  });
}
</script>

<style scoped>
.chat-main {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.channel-drawer {
  height: 100vh;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  background-color: rgb(var(--v-theme-surface));
}

.chat-header {
  position: sticky;
  top: 0;
  z-index: 1;
}

.chat-footer {
  position: sticky;
  bottom: 0;
  background-color: rgb(var(--v-theme-surface));
}

.message-wrapper {
  display: flex;
  align-items: flex-start;
  margin-bottom: 16px;
}

.message-mine {
  justify-content: flex-end;
}

.message-bubble {
  max-width: 70%;
  padding: 8px 16px;
  border-radius: 16px;
  background-color: rgb(var(--v-theme-primary));
  color: white;
}

.message-other .message-bubble {
  background-color: rgb(var(--v-theme-surface-variant));
  color: inherit;
}

.message-mine .message-bubble {
  border-bottom-right-radius: 4px;
}

.message-other .message-bubble {
  border-bottom-left-radius: 4px;
}
</style>