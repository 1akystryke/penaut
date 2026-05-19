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
       <v-btn @click="sidePanelMode=`channels`" >Каналы</v-btn><v-btn @click="sidePanelMode=`directs`;getUsers()">Личные сообщения</v-btn>
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
                <v-icon>{{"mdi-"+channel.icon }}</v-icon>
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
                <v-icon>{{"mdi-"+channel.icon }}</v-icon>
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
            <v-icon>{{'mdi-'+ activeChannelData.icon }}</v-icon>
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
              <v-icon size="18" :icon="'mdi-'+message.avatar"></v-icon>
            </v-avatar>
            {{ message.sender }}
            </div>
            
            <v-list-item-content :class="{ 'text-right': message.isMine }">
              <v-list-item-title>{{ message.text }} <br/>
                {{ formatMessengerDate(message.time) }}
              </v-list-item-title>
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
          :disabled="!newMessage.trim()"
        />
      </v-footer>
    </v-main>
  </v-app>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue';
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

const createDirectDialog = ref(false)
const directList = ref([])
const newDirectId = ref(null)

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

function addUserDialogOpen(){
  getUsers()
  addUserDialog.value=true
}
// Выбор канала
function selectChannel(channelId) {
  activeChannel.value = channelId;
  messages.value = getChannelMessages(channelId);
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


    messages.value.push(data)
    scrollToBottom()
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
    'Content-Type':'application/json',
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

async function loadChannels(){
  try{
    let channels = await request('/channels') 
    
    channelsList.value = []
    channels.forEach(elem => {
      elem.icon = 'peanut'
      if (elem.type==="direct") elem.directName = extractDirectName(elem.name)
      else elem.directName = ""
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

    posts.forEach(element => {
      element.time = parseCustomDate(element.created_at)
      element.sender = element.user_name
      element.avatar = 'peanut-outline'
      if (element.author==store.meId){
        element.isMine = true
      }
      else{
        element.isMine = false

      }
    });

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

    // websocket first
    if(socket && socket.readyState===1){

      socket.send(JSON.stringify({text}))

    }else{

      // fallback REST
      await request(
        `/channels/${activeChannel.value.id}/posts`,
        {
          method:'POST',
          body:JSON.stringify({text})
        }
      )

      await loadMessages(activeChannel.value.id)
    }

    newMessage.value=''

  }catch(e){
    alert(e.message)
  }
}

// Функция для добавления канала (заглушка)
async function addChannel() {
  if (newChannelName.value==""){
    alert("заполните имя сначала")
    return
  }
  try{

    await request('/channels',{
      method:'POST',
      body:JSON.stringify({
        "type":newChannelType.value,
        "name":newChannelName.value
      })
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
</script>
