<template>
  <v-app class="app-root">
    <CreateChannelDialog v-model="createChannelDialog" @create="onCreateChannel" />
    <ChannelInfoDialog
      v-model="channelInfoDialog"
      :channel="activeChannelData"
      :members="activeChannelMembers"
      @add-user="addUserDialogOpen"
    />
    <AddUserDialog
      v-model="addUserDialog"
      :channel="activeChannelData"
      :users="usersList"
      @submit="onAddUser"
    />
    <CreateDirectDialog v-model="createDirectDialog" :users="usersList" @create="onCreateDirect" />

    <v-navigation-drawer v-model="drawer" :width="268" permanent class="sidebar">
      <div class="sidebar-workspace">
        <div class="workspace-icon">
          <v-icon size="17" color="white">mdi-peanut</v-icon>
        </div>
        <span class="workspace-title">Peanut</span>
      </div>

      <div class="sidebar-tabs">
        <button
          class="sidebar-tab"
          :class="{ active: sidePanelMode === 'channels' }"
          @click="sidePanelMode = 'channels'"
        >Каналы</button>
        <button
          class="sidebar-tab"
          :class="{ active: sidePanelMode === 'directs' }"
          @click="sidePanelMode = 'directs'; getUsers()"
        >Сообщения</button>
      </div>

      <ChannelList
        v-if="sidePanelMode === 'channels'"
        :channels="channelsList.filter(c => c.type !== 'direct')"
        :active-channel="activeChannel"
        @select="selectChannel"
        @create="createChannelDialog = true"
      />
      <DirectList
        v-if="sidePanelMode === 'directs'"
        :directs="channelsList.filter(c => c.type === 'direct')"
        :active-channel="activeChannel"
        @select="selectChannel"
        @create="createDirectDialog = true"
      />
    </v-navigation-drawer>

    <v-main>
      <div class="chat-layout">
        <div class="chat-header" :class="{ 'chat-header--empty': !activeChannelData }">
          <template v-if="activeChannelData">
            <div class="chat-header-info">
              <v-avatar size="36" rounded="lg" class="chat-avatar">
                <v-img :src="activeChannelData.url" />
              </v-avatar>
              <div>
                <div class="chat-name">
                  {{ activeChannelData.type === 'direct' ? activeChannelData.directName : activeChannelData.name }}
                </div>
                <div class="chat-sub">{{ activeChannelMembers.length }} участников</div>
              </div>
            </div>
            <button class="header-action" @click="channelInfoDialog = true">
              <v-icon size="19">mdi-information-outline</v-icon>
            </button>
          </template>
        </div>

        <div class="messages-scroll" ref="messagesContainer">
          <div v-if="!activeChannel" class="empty-state">
            <div class="empty-icon">
              <v-icon size="36" color="#6366f1">mdi-message-text-outline</v-icon>
            </div>
            <p class="empty-title">Выберите канал</p>
            <p class="empty-sub">Начните общение, выбрав канал или переписку из списка</p>
          </div>
          <MessageList v-else :messages="messages" :channel-url="activeChannelData?.url" />
        </div>

        <MessageComposer v-if="activeChannel" :is-sending="isSendingMessage" @send="onSendMessage" />
      </div>
    </v-main>
  </v-app>
</template>

<script setup>
import { ref, computed, nextTick, onBeforeUnmount } from 'vue'
import { request } from '@/api/client'
import { useChannels } from '@/composables/useChannels'
import { useMessages } from '@/composables/useMessages'
import { useWebSocket } from '@/composables/useWebSocket'
import { useAttachments } from '@/composables/useAttachments'

import CreateChannelDialog from './dialogs/CreateChannelDialog.vue'
import CreateDirectDialog from './dialogs/CreateDirectDialog.vue'
import ChannelInfoDialog from './dialogs/ChannelInfoDialog.vue'
import AddUserDialog from './dialogs/AddUserDialog.vue'
import ChannelList from './sidebar/ChannelList.vue'
import DirectList from './sidebar/DirectList.vue'
import MessageList from './chat/MessageList.vue'
import MessageComposer from './chat/MessageComposer.vue'

const drawer = ref(true)
const sidePanelMode = ref('channels')
const activeChannel = ref(null)
const activeChannelMembers = ref([])
const usersList = ref([])
const createChannelDialog = ref(false)
const channelInfoDialog = ref(false)
const addUserDialog = ref(false)
const createDirectDialog = ref(false)
const messagesContainer = ref(null)

const { channelsList, loadChannels, addChannel, addDirect } = useChannels(usersList)
const { messages, isSendingMessage, getChannelMessages, sendMessage, normalizePost } = useMessages()
const { hydrateMessageAttachments, revokeAll } = useAttachments()
const ws = useWebSocket(messages, scrollToBottom)

const activeChannelData = computed(() => channelsList.value.find(c => c.id === activeChannel.value))

async function getUsers() {
  try { usersList.value = await request('/users') } catch (e) { alert(e.message) }
}

async function getChannelMembers(channelId) {
  try { activeChannelMembers.value = await request(`/channels/${channelId}/members`) } catch (e) { alert(e.message) }
}

function selectChannel(channelId) {
  activeChannel.value = channelId
  getChannelMessages(channelId)
  ws.subscribe(channelId)
  getChannelMembers(channelId)
  scrollToBottom()
}

async function onCreateChannel({ name, type, imageFile }) {
  const ok = await addChannel(name, type, imageFile)
  if (ok) createChannelDialog.value = false
}

async function onCreateDirect(userId) {
  const ok = await addDirect(userId)
  if (ok) { createDirectDialog.value = false; await loadChannels() }
}

async function onAddUser(userId) {
  try {
    await request(`/channels/${activeChannel.value}/members`, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId })
    })
    await getChannelMembers(activeChannel.value)
    addUserDialog.value = false
  } catch (e) { alert(e.message) }
}

function addUserDialogOpen() {
  getUsers()
  addUserDialog.value = true
}

async function onSendMessage({ text, files }) {
  const post = await sendMessage(activeChannel.value, text, files)
  if (post && !ws.isConnected()) {
    normalizePost(post)
    await hydrateMessageAttachments(post)
    messages.value.push(post)
    scrollToBottom()
  }
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value)
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  })
}

onBeforeUnmount(() => { ws.close(); revokeAll() })

getUsers()
loadChannels()
</script>

<style>
.app-root .v-navigation-drawer.sidebar {
  background: #0d0d1c !important;
  border-right: 1px solid rgba(255, 255, 255, 0.05) !important;
}
.app-root .v-navigation-drawer.sidebar .v-navigation-drawer__content {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.app-root .v-main {
  background: #eef0f6;
}
.app-root .v-main .v-main__wrap {
  display: flex;
  flex-direction: column;
  height: 100%;
}
</style>

<style scoped>
.sidebar-workspace {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 18px 16px 15px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  flex-shrink: 0;
}

.workspace-icon {
  width: 31px;
  height: 31px;
  border-radius: 9px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 3px 10px rgba(99, 102, 241, 0.45);
}

.workspace-title {
  font-size: 15px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.92);
  letter-spacing: 0.02em;
}

.sidebar-tabs {
  display: flex;
  gap: 3px;
  padding: 10px 10px 6px;
  flex-shrink: 0;
}

.sidebar-tab {
  flex: 1;
  padding: 6px 8px;
  border-radius: 7px;
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.38);
  font-size: 12.5px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.14s ease;
  font-family: inherit;
}

.sidebar-tab:hover {
  background: rgba(255, 255, 255, 0.07);
  color: rgba(255, 255, 255, 0.7);
}

.sidebar-tab.active {
  background: rgba(99, 102, 241, 0.2);
  color: #a5b4fc;
  font-weight: 600;
}

/* Chat layout */
.chat-layout {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 14px 0 20px;
  height: 58px;
  background: #ffffff;
  border-bottom: 1px solid rgba(0, 0, 0, 0.07);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
  flex-shrink: 0;
}

.chat-header--empty {
  background: #ffffff;
}

.chat-header-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.chat-avatar {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
}

.chat-name {
  font-size: 14.5px;
  font-weight: 600;
  color: #111827;
  line-height: 1.25;
}

.chat-sub {
  font-size: 11.5px;
  color: #9ca3af;
  line-height: 1.25;
  margin-top: 1px;
}

.header-action {
  width: 33px;
  height: 33px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: #9ca3af;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.14s ease;
}

.header-action:hover {
  background: #f3f4f6;
  color: #6366f1;
}

.messages-scroll {
  flex: 1;
  overflow-y: auto;
  padding: 24px 28px 12px;
  min-height: 0;
  scroll-behavior: smooth;
}

.messages-scroll::-webkit-scrollbar { width: 4px; }
.messages-scroll::-webkit-scrollbar-track { background: transparent; }
.messages-scroll::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 2px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  padding-bottom: 80px;
}

.empty-icon {
  width: 72px;
  height: 72px;
  border-radius: 22px;
  background: rgba(99, 102, 241, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 18px;
}

.empty-title {
  font-size: 17px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 7px;
}

.empty-sub {
  font-size: 13px;
  color: #9ca3af;
  margin: 0;
  max-width: 260px;
  line-height: 1.55;
}
</style>
