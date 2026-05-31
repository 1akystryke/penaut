<template>
  <div class="msg-list">
    <MessageItem
      v-for="message in messages"
      :key="message.id"
      :message="message"
      :channel-url="channelUrl"
    />
  </div>
</template>

<script setup>
import MessageItem from './MessageItem.vue'

defineProps({ messages: Array, channelUrl: String })
</script>

<style scoped>
.msg-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding-bottom: 4px;
}
</style>
