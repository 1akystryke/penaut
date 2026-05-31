<template>
  <div class="msg-row" :class="{ 'msg-row--mine': message.isMine }">
    <v-avatar v-if="!message.isMine" size="30" rounded="lg" class="msg-avatar">
      <v-img :src="channelUrl" />
    </v-avatar>

    <div class="msg-body">
      <span v-if="!message.isMine" class="msg-sender">{{ message.sender }}</span>
      <div class="msg-bubble" :class="{ 'msg-bubble--mine': message.isMine }">
        <p class="msg-text">{{ message.text }}</p>
        <div v-if="message.attachments?.length" class="msg-attachments">
          <a
            v-for="attachment in message.attachments"
            :key="attachment.file_path"
            :href="attachment.url"
            target="_blank"
            rel="noopener"
            class="msg-attachment"
          >
            <img v-if="isImage(attachment)" :src="attachment.url" class="attach-img" />
            <div v-else class="attach-file">
              <v-icon size="14">mdi-file-outline</v-icon>
              <span>{{ attachment.file_path.split('/').pop() }}</span>
            </div>
          </a>
        </div>
        <span class="msg-time">{{ formatMessengerDate(message.time) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { formatMessengerDate } from '@/utils/date'

defineProps({ message: Object, channelUrl: String })

function isImage(attachment) {
  return attachment.file_type?.startsWith('image/')
}
</script>

<style scoped>
.msg-row {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  margin-bottom: 2px;
}

.msg-row--mine {
  flex-direction: row-reverse;
}

.msg-avatar {
  flex-shrink: 0;
  margin-bottom: 2px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.12);
}

.msg-body {
  max-width: 62%;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.msg-row--mine .msg-body {
  align-items: flex-end;
}

.msg-sender {
  font-size: 11px;
  font-weight: 600;
  color: #6366f1;
  padding: 0 12px;
  line-height: 1;
}

.msg-bubble {
  background: #ffffff;
  border-radius: 16px 16px 16px 4px;
  padding: 9px 13px 7px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06), 0 1px 2px rgba(0, 0, 0, 0.04);
  max-width: 100%;
}

.msg-bubble--mine {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  border-radius: 16px 16px 4px 16px;
  box-shadow: 0 2px 10px rgba(99, 102, 241, 0.32);
}

.msg-text {
  font-size: 14px;
  line-height: 1.5;
  color: #1f2937;
  margin: 0 0 5px;
  word-break: break-word;
  white-space: pre-wrap;
}

.msg-bubble--mine .msg-text {
  color: rgba(255, 255, 255, 0.96);
}

.msg-time {
  font-size: 10px;
  color: #b0b8c8;
  display: block;
  text-align: right;
  line-height: 1;
}

.msg-bubble--mine .msg-time {
  color: rgba(255, 255, 255, 0.5);
}

.msg-attachments {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 6px;
}

.msg-attachment {
  text-decoration: none;
  display: block;
}

.attach-img {
  width: 200px;
  height: 140px;
  object-fit: cover;
  border-radius: 10px;
  display: block;
}

.attach-file {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  background: rgba(0, 0, 0, 0.06);
  border-radius: 8px;
  font-size: 12px;
  color: #4b5563;
  max-width: 220px;
}

.msg-bubble--mine .attach-file {
  background: rgba(255, 255, 255, 0.16);
  color: rgba(255, 255, 255, 0.92);
}

.attach-file span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
