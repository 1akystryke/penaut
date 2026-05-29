<template>
  <v-app>
    <v-dialog v-model="createChannelDialog">
      <v-col>
      <v-card>
        <v-card-title>Создать канал</v-card-title>
        <v-card-item>
          <v-text-field
          v-model="newChannelName"
          density="compact"
          variant="outlined"
          label="Название"
        ></v-text-field>
        <v-select label="Тип" v-model="newChannelType" :items="channelTypes"></v-select>
        <div class="image-uploader">
    <!-- Отображаем input только если изображение еще не выбрано -->
    <v-file-input
      v-if="!imagePreviewUrl"
      label="Выберите изображение"
      accept="image/*"
      prepend-icon="mdi-camera"
      variant="outlined"
      @update:model-value="handleImageSelect"
    ></v-file-input>

    <!-- Если изображение выбрано, показываем превью -->
    <div v-else>
      <v-card class="mx-auto" max-width="400">
        <!-- Блок с превью -->
        <div class="position-relative">
          <v-img
            :src="imagePreviewUrl"
            height="300"
            cover
            class="bg-grey-lighten-2"
          >
            <!-- Слот для overlay при желании -->
          </v-img>
          
          <!-- Кнопки управления поверх изображения -->
          <div class="image-actions">
            <v-btn
              icon="mdi-delete"
              color="error"
              size="small"
              class="mr-2"
              @click="removeImage"
            ></v-btn>
            <v-btn
              icon="mdi-reload"
              color="info"
              size="small"
              @click="replaceImage"
            ></v-btn>
          </div>
        </div>
        
        <!-- Информация о файле -->
        <v-card-text class="text-center">
          <div class="text-subtitle-1 text-truncate">
            {{ selectedFileName }}
          </div>
          <div class="text-caption text-grey">
            {{ formatFileSize(selectedFileSize) }}
          </div>
        </v-card-text>
      </v-card>
    </div>
  </div>
        <v-btn @click="addChannel">Создать</v-btn>
      </v-card-item>
      </v-card>
      </v-col>
    </v-dialog>

    <v-dialog v-model="channelInfoDialog">
      <v-col>
      <v-card>
        <v-card-title>{{  activeChannelData.name}}</v-card-title>
        <v-card-item>
          <v-list>
            <v-list-item>
              Тип канала: {{ activeChannelData["type"] }}
            </v-list-item>
            <v-list-item>
              ID канала: {{ activeChannelData["id"] }}
            </v-list-item>
            <v-list-item>
              Члены: 
              <div v-for="member in activeChannelMembers">{{ member["name"] }}</div>
              <v-btn v-if="activeChannelData['type']!='direct'" @click="addUserDialogOpen">Добавить</v-btn>
            </v-list-item>
          </v-list>
          
      </v-card-item>
      </v-card>
      </v-col>
    </v-dialog>
    <v-dialog v-model="addUserDialog">
      <v-card>
        <v-card-title>Добавить пользователя в канал {{ activeChannelData["name"] }}</v-card-title>
        <v-card-item><v-select 
          :items="usersList" 
          item-title="name"
          item-value="id"  
          v-model="selectedUserToAdd"
          >
        </v-select></v-card-item>
        <v-card-item ><v-btn @click="sumbitAddUser">apply</v-btn></v-card-item>
      </v-card>
    </v-dialog>
    
    <v-dialog v-model="createDirectDialog">
      <v-col>
      <v-card>
        <v-card-title>Создать переписку</v-card-title>
        <v-card-item>
          <v-select
          v-model="newDirectId"
          :items="usersList" 
          item-title="name"
          item-value="id"  
          density="compact"
          variant="outlined"
          label="С кем?"
        ></v-select>
        <v-btn @click="addDirect">Создать</v-btn>
      </v-card-item>
      </v-card>
      </v-col>
    </v-dialog>

    <v-navigation-drawer
      v-model="drawer"
      width="320"
      permanent
      class="channel-drawer"
    >
      <!-- Заголовок панели каналов -->
       <v-row>
        <v-btn-toggle>
          <v-col cols="6">
          <v-btn block @click="sidePanelMode=`channels`" >Каналы</v-btn>
          </v-col>
          <v-col cols="6">
          <v-btn block @click="sidePanelMode=`directs`;getUsers()">Личные сообщения</v-btn>
          </v-col>
       </v-btn-toggle>
       </v-row>
       <div v-if="sidePanelMode===`channels`">
        <v-list-item
          class="pa-4"
          title="Каналы"
          prepend-icon="mdi-account-group"
        >
          <template #append>
            <v-btn
              icon="mdi-plus"
              variant="text"
              size="small"
              @click="createChannelDialog=true"
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
                <v-img :src="channel.url"></v-img>
              </v-avatar>
            </template>

            <template #append>
              <div class="text-caption text-medium-emphasis">
                {{ channel.type }}
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
      </div>
      <div v-if="sidePanelMode===`directs`">
        <v-list-item
          class="pa-4"
          title="Личные сообщения"
          prepend-icon="mdi-account"
        >
          <template #append>
            <v-btn
              icon="mdi-plus"
              variant="text"
              size="small"
              @click="createDirectDialog=true"
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
            v-for="channel in filteredDirects"
            :key="channel.id"
            :value="channel.id"
            :title="channel.directName"
            :subtitle="channel.lastMessage"
            :active="activeChannel === channel.id"
            @click="selectChannel(channel.id)"
            color="primary"
          >
            <template #prepend>
              <v-avatar size="40" color="surface-variant">
                
                <v-img :src="channel.url"></v-img>
              </v-avatar>
            </template>

            <template #append>
              <div class="text-caption text-medium-emphasis">
                {{ channel.type }}
              </div>
            </template>
          </v-list-item>
        </v-list>

        <!-- Если каналов нет -->
        <div
          v-if="filteredDirects.length === 0"
          class="text-center pa-4 text-medium-emphasis"
        >
          Личных сообщений пока нет
        </div>
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
            <v-img :src="activeChannelData.url"></v-img>
          </v-avatar>
          <div>
            <div v-if="activeChannelData.type!=`direct`" class="text-h6">{{ activeChannelData.name }}</div>
            <div v-if="activeChannelData.type==`direct`" class="text-h6">{{ activeChannelData.directName }}</div>
            <div class="text-caption text-medium-emphasis">
              {{ activeChannelMembers.length }} участников
            </div>
          </div>
        </template>

        <template #append>
          
          <v-btn icon="mdi-information-outline" variant="text" @click="showChannelInfo" />
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
          <v-list class="overflow-y-auto">
            <v-list-item v-for="message in messages"
            :key="message.id">
            <div v-if="!message.isMine">
              <v-avatar
              size="32"
              class="mr-2"
              color="surface-variant"
              >
              <v-img :src="activeChannelData.url"></v-img>
            </v-avatar>
            {{ message.sender }}
            </div>
            
            <v-list-item-content :class="{ 'text-right': message.isMine }">
              <v-list-item-title>{{ message.text }} <br/>
                {{ formatMessengerDate(message.time) }}
              </v-list-item-title>
              <div
                v-if="message.attachments && message.attachments.length"
                class="message-attachments"
                :class="{ 'justify-end': message.isMine }"
              >
                <a
                  v-for="attachment in message.attachments"
                  :key="attachment.file_path"
                  :href="attachment.url"
                  target="_blank"
                  rel="noopener"
                  class="message-attachment"
                >
                  <v-img
                    v-if="isImageAttachment(attachment)"
                    :src="attachment.url"
                    width="180"
                    height="120"
                    cover
                    rounded="lg"
                  />
                  <v-chip
                    v-else
                    prepend-icon="mdi-paperclip"
                    variant="tonal"
                    size="small"
                  >
                    {{ attachment.file_path }}
                  </v-chip>
                </a>
              </div>
            </v-list-item-content>
          </v-list-item>
          </v-list>
        </div>
      </div>

      <!-- Поле ввода сообщения -->
      <v-footer
        v-if="activeChannel"
        app
        class="chat-footer pa-2"
        border="top"
      >
        <div class="composer">
          <div v-if="messageFiles.length" class="composer-attachments">
            <v-chip
              v-for="(file, index) in messageFiles"
              :key="file.name + file.size + index"
              closable
              prepend-icon="mdi-paperclip"
              variant="tonal"
              size="small"
              @click:close="removeMessageAttachment(index)"
            >
              {{ file.name }}
            </v-chip>
          </div>
          <div class="composer-input">
            <input
              ref="messageFileInput"
              type="file"
              multiple
              class="d-none"
              @change="handleMessageFilesSelect"
            />
            <v-btn
              icon="mdi-paperclip"
              variant="text"
              class="mr-1"
              @click="openMessageFilePicker"
            />
            <v-text-field
              v-model="newMessage"
              variant="outlined"
              density="compact"
              placeholder="Введите сообщение..."
              hide-details
              @keyup.enter="sendMessage2"
            >
              <template #prepend-inner>
                <v-btn icon="mdi-emoticon-outline" variant="text" size="small" />
              </template>
            </v-text-field>
            <v-btn
              icon="mdi-send"
              color="primary"
              variant="text"
              class="ml-2"
              @click="sendMessage2"
              :disabled="isSendingMessage || !newMessage.trim()"
            />
          </div>
        </div>
      </v-footer>
    </v-main>
  </v-app>
</template>

<script setup>
import { ref, computed, watch, nextTick, onBeforeUnmount } from 'vue';
import { authStore } from '@/stores/authStore.vue'
const store = authStore()

const addUserDialog = ref(false);
const usersList = ref([])
const selectedUserToAdd = ref(null)

const createChannelDialog = ref(false);
const drawer = ref(true);
const activeChannel = ref(null);
const activeChannelMembers = ref([]);
const searchQuery = ref('');
const newMessage = ref('');
const messagesContainer = ref(null);
const newChannelName = ref("");
const newChannelType = ref("public")
const channelTypes = ref(["public","private"])
const channelsList = ref([])
const messages = ref([])
const channelInfoDialog = ref(false);
const sidePanelMode = ref('channels')
const directPictures = ref({})
const createDirectDialog = ref(false)
const newDirectId = ref(null)
const messageFiles = ref([])
const messageFileInput = ref(null)
const isSendingMessage = ref(false)
const attachmentObjectUrls = new Set()

var socket = null



// Фильтрация каналов по поиску
const filteredChannels = computed(() => {
  if (!searchQuery.value) return channelsList.value.filter(channel => channel.type != 'direct');
  return channelsList.value.filter(channel =>
    channel.name.toLowerCase().includes(searchQuery.value.toLowerCase()) && channel.type != 'direct'
  );
});



const filteredDirects = computed(() => {
  if (!searchQuery.value) return channelsList.value.filter(channel => channel.type === 'direct');
  return channelsList.value.filter(channel =>
    channel.directName.toLowerCase().includes(searchQuery.value.toLowerCase()) && channel.type === 'direct'
  );
});



// Активный канал
const activeChannelData = computed(() => {
  return channelsList.value.find(c => c.id === activeChannel.value);
});

// Сообщения активного канала
const activeMessages = computed(() => {
  if (!activeChannel.value) return [];
  return messages.value[activeChannel.value] || [];
});
function extractDirectID(directRawName){
  const users = directRawName.split("_");

  for (const user of users) {
      
      if (user !== store.meId) {
        return user
      }
  }

}

function extractDirectName(directRawName){
  const users = directRawName.split("_");

  for (const user of users) {
      
      if (user !== store.meId) {
        for (const searchUser of usersList.value){
          if (searchUser.id==user){
            return searchUser.name
          }
        }
      }
  }

}
// users = directRawName.split("_")
// for user in users:
//   if user.ID != store.meId:
//     return user.Name
    



async function getUsers(){
  try{
    usersList.value = await request('/users')
  }catch(e){
    alert(e.message)
  }
}

async function addDirect() {
  try{
    await request('/users/'+newDirectId.value+'/direct',{
      method:'POST',
      body:JSON.stringify({"abobe":"obeba"})
    })
    createDirectDialog.value=false
  }catch(e){
    alert(e.message)
  }
}

async function sumbitAddUser(){
  
  try{
    await request('/channels/'+activeChannel.value+'/members',{
      method:'POST',
      body:JSON.stringify({
        "user_id":selectedUserToAdd.value
      })
    })
    getChannelMembers(activeChannel.value)
    addUserDialog.value=false
  }catch(e){
    alert(e.message)
  }


}

async function loadImage(userID) {
        try {
            const response = await fetch(store.API + `/users/` + userID + `/pic`, {
            headers: {
                'Authorization': `Bearer ${store.token}`, // ваш заголовок
                // или другие кастомные заголовки
            }
            })
            
            if (!response.ok) throw new Error('Ошибка загрузки')
            
            const blob = await response.blob()
            return URL.createObjectURL(blob)
        } catch (error) {
            console.error('Не удалось загрузить аватар:', error)
            return '' // или путь к заглушке
        }
    }

function addUserDialogOpen(){
  getUsers()
  addUserDialog.value=true
}
// Выбор канала
function selectChannel(channelId) {
  activeChannel.value = channelId;
  getChannelMessages(channelId);
  subscribeWebSocket(channelId);
  getChannelMembers(channelId)
  scrollToBottom();
  
}

function showChannelInfo(){
  channelInfoDialog.value = true
}

async function getChannelMembers(channelId){
  try{
    const members=await request(`/channels/${channelId}/members`)
    activeChannelMembers.value = members
  }catch(e){
    alert(e.message)
  }
}

function subscribeWebSocket(channelId){

  if(socket){
    socket.close()
  }
  
  socket=new WebSocket(
    `${store.WS}?channel_id=${channelId}&token=${store.token}`
  )

  socket.onmessage=(e)=>{
    const data=JSON.parse(e.data)

    if(data.error){
      alert(data.error)
      return
    }

    data.time = parseCustomDate(data.created_at)
    data.sender = data.user_name
    data.avatar = 'peanut-outline'
    if (data.author==store.meId){
      data.isMine = true
    }
    else{
      data.isMine = false

    }


    hydrateMessageAttachments(data).then(() => {
      messages.value.push(data)
      scrollToBottom()
    })
  }

  socket.onerror=()=>{
    showError('WebSocket error')
  }
}

function parseCustomDate(dateString) {
    // Просто используем встроенный парсер, так как формат ISO 8601 с микросекундами
    // JS Date API понимает этот формат (микросекунды просто игнорируются)
    const date = new Date(dateString);
    
    // Проверка на валидность
    if (isNaN(date.getTime())) {
        throw new Error('Invalid date format');
    }
    
    return date;
}

async function request(path,options={}){
  const headers={
    ...(options.headers||{})
  }
  
  if(store.status){
    headers.Authorization=`Bearer ${store.token}`
  }
  const r=await fetch(store.API+path,{
    ...options,
    headers
  })

  const data=await r.json().catch(()=>({}))
  if (r.status===401){
    store.breakAuth()
  }
  if(!r.ok){
    throw new Error(data.error||'Request failed')
  }

  return data
}

async function requestBlob(path){
  const headers={}

  if(store.status){
    headers.Authorization=`Bearer ${store.token}`
  }

  const r=await fetch(store.API+path,{
    headers
  })

  if (r.status===401){
    store.breakAuth()
  }
  if(!r.ok){
    throw new Error('Request failed')
  }

  return r.blob()
}

async function loadChannels(){
  try{
    let channels = await request('/channels') 
    
    channelsList.value = []
    channels.forEach(elem => {
      elem.icon = 'peanut'
      if (elem.type==="direct") { 
        elem.directName = extractDirectName(elem.name)
        elem.url = store.API+'/users/'+extractDirectID(elem.name)+'/pic?token='+store.token

      }
      else {elem.directName = ""
        elem.url = store.API+'/channels/'+elem.id+'/pic?token='+store.token}
    })
    channelsList.value = channels
  }catch(e){
    alert(e.message)
  }
}
function formatMessengerDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    
    // Проверка на валидность
    if (isNaN(date.getTime())) {
        throw new Error('Invalid date');
    }
    
    // Разница в миллисекундах
    const diffMs = now - date;
    const diffSeconds = Math.floor(diffMs / 1000);
    const diffMinutes = Math.floor(diffSeconds / 60);
    const diffHours = Math.floor(diffMinutes / 60);
    const diffDays = Math.floor(diffHours / 24);
    
    // Если прошло меньше 24 часов
    if (diffDays < 1) {
        // Форматируем время: HH:MM
        const hours = date.getHours().toString().padStart(2, '0');
        const minutes = date.getMinutes().toString().padStart(2, '0');
        return `${hours}:${minutes}`;
    }
    
    // Если прошло 1 день или больше - краткая дата + время
    // Формат: ДД.ММ ГГГГ, ЧЧ:ММ (как в WhatsApp/Telegram)
    const day = date.getDate().toString().padStart(2, '0');
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const year = date.getFullYear();
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    
    // Для текущего года можно не показывать год (опционально)
    if (year === now.getFullYear()) {
        return `${day}.${month} в ${hours}:${minutes}`;
    }
    
    return `${day}.${month}.${year} в ${hours}:${minutes}`;
}
async function getChannelMessages(channelId){
  try{
    const posts=await request(`/channels/${channelId}/posts`)

    await Promise.all(posts.map(async element => {
      element.time = parseCustomDate(element.created_at)
      element.sender = element.user_name
      element.avatar = 'peanut-outline'
      if (element.author==store.meId){
        element.isMine = true
      }
      else{
        element.isMine = false

      }
      await hydrateMessageAttachments(element)
    }));

    messages.value = posts
  }catch(e){
    alert(e.message)
  }
}

async function sendMessage2(){

  const text=newMessage.value.trim()

  if(!text){
    return
  }

  try{
    isSendingMessage.value = true

    const formData = new FormData();
    formData.append('text', text);
    messageFiles.value.forEach(file => {
      formData.append('attachments[]', file);
    })

    const post = await request(
        `/channels/${activeChannel.value}/posts`,
        {
          method:'POST',
          body:formData
        }
      )
    if(!socket || !socket.readyState===1){
      post.time = parseCustomDate(post.created_at)
      post.sender = post.user_name
      post.avatar = 'peanut-outline'
      post.isMine = post.author==store.meId
      await hydrateMessageAttachments(post)
      messages.value.push(post)
      scrollToBottom()
    }

    newMessage.value=''
    clearMessageAttachments()

  }catch(e){
    alert(e.message)
  }finally{
    isSendingMessage.value = false
  }
}

function openMessageFilePicker(){
  messageFileInput.value?.click()
}

function handleMessageFilesSelect(event){
  const files = Array.from(event.target.files || [])
  messageFiles.value = [...messageFiles.value, ...files]
  event.target.value = ''
}

function removeMessageAttachment(index){
  messageFiles.value.splice(index, 1)
}

function clearMessageAttachments(){
  messageFiles.value = []
  if(messageFileInput.value){
    messageFileInput.value.value = ''
  }
}

async function hydrateMessageAttachments(message){
  try{
    const attachments = await request(`/posts/${message.id}/attachments`)
    message.attachments = await Promise.all(attachments.map(async attachment => {
      const blob = await requestBlob(`/attachment/${encodeURIComponent(attachment.file_path)}`)
      const url = URL.createObjectURL(blob)
      attachmentObjectUrls.add(url)
      return {
        ...attachment,
        url,
        file_type: attachment.file_type || blob.type,
      }
    }))
  }catch(e){
    message.attachments = []
  }
}

function isImageAttachment(attachment){
  return attachment.file_type?.startsWith('image/')
}

// Функция для добавления канала (заглушка)
async function addChannel() {
  if (newChannelName.value==""){
    alert("заполните имя сначала")
    return
  }
  try{

    

    const formData = new FormData();
    formData.append('type', newChannelType.value);
    formData.append("name", newChannelName.value);
    if (imageFile.value) {
        formData.append('pic', imageFile.value);
    }
    await request('/channels',{
      
      method:'POST',
      body: formData
    })

    await loadChannels()
  createChannelDialog.value = false

  }catch(e){
    alert(e.message)
  }


}
getUsers()
loadChannels()
// Прокрутка вниз
function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
    }
  });
}




// Хранилище данных компонента
const imageFile = ref(null)        // Объект File
const imagePreviewUrl = ref(null)  // Blob URL для отображения
const selectedFileName = ref('')
const selectedFileSize = ref(0)

// Обработка выбора файла
const handleImageSelect = (file) => {
  if (!file) return
  
  // Проверка типа файла
  if (!file.type.startsWith('image/')) {
    console.error('Пожалуйста, выберите изображение')
    return
  }
  
  // Сохраняем метаданные
  imageFile.value = file
  selectedFileName.value = file.name
  selectedFileSize.value = file.size
  
  // Создаем URL для предпросмотра
  if (imagePreviewUrl.value) {
    URL.revokeObjectURL(imagePreviewUrl.value)
  }
  imagePreviewUrl.value = URL.createObjectURL(file)
}

// Замена изображения
const replaceImage = () => {
  // Очищаем и показываем input снова
  imagePreviewUrl.value = null
  imageFile.value = null
}

// Удаление изображения
const removeImage = () => {
  if (imagePreviewUrl.value) {
    URL.revokeObjectURL(imagePreviewUrl.value)
  }
  imagePreviewUrl.value = null
  imageFile.value = null
  selectedFileName.value = ''
  selectedFileSize.value = 0
}

// Форматирование размера файла (байты -> человекочитаемый формат)
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}



const cleanup = () => {
  if (imagePreviewUrl.value) {
    URL.revokeObjectURL(imagePreviewUrl.value)
  }
  attachmentObjectUrls.forEach(url => URL.revokeObjectURL(url))
  attachmentObjectUrls.clear()
}

onBeforeUnmount(() => {
  cleanup()
})
</script>

<style scoped>
.composer {
  width: 100%;
}

.composer-input {
  align-items: center;
  display: flex;
  gap: 4px;
  width: 100%;
}

.composer-attachments {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 0 0 8px 48px;
}

.message-attachments {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 8px;
}

.message-attachment {
  color: inherit;
  display: inline-flex;
  max-width: 220px;
  text-decoration: none;
}
</style>
