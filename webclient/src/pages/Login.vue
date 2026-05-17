<template>
  <v-container class="fill-height d-flex flex-column justify-center" max-width="1100">
    <v-col>
    <v-card
      class="mx-auto pa-12 pb-8"
      elevation="3"
      max-width="448"
      rounded="lg"
    >
      <div class="text-body-large text-medium-emphasis">Account</div>

      <v-text-field
        v-model="email"
        density="compact"
        placeholder="Email address"
        prepend-inner-icon="mdi-email-outline"
        variant="outlined"
      ></v-text-field>

      <div class="text-body-large text-medium-emphasis d-flex align-center justify-space-between">
        Password
      </div>

      <v-text-field
        v-model="password"
        :append-inner-icon="visible ? 'mdi-eye-off' : 'mdi-eye'"
        :type="visible ? 'text' : 'password'"
        density="compact"
        placeholder="Enter your password"
        prepend-inner-icon="mdi-lock-outline"
        variant="outlined"
        @click:append-inner="visible = !visible"
      ></v-text-field>

      <v-card
        class="mb-12"
        color="surface-variant"
        variant="tonal"
      >
      </v-card>

      <v-btn
        class="mb-8"
        color="blue"
        size="large"
        variant="tonal"
        @click="tryLogin"
        block
      >
        Log In
      </v-btn>

    </v-card>
    </v-col>
  </v-container>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'

const DEFAULT_HOST='localhost:8080'
const API=`http://${DEFAULT_HOST}`
const WS=`ws://${DEFAULT_HOST}/ws`

const email = ref("")
const password = ref("")


async function hashPassword(password) {
    const encoder = new TextEncoder();
    const data = encoder.encode(password);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
    return hashHex;
}


async function request(path,options={}){
  const headers={
    'Content-Type':'application/json',
    ...(options.headers||{})
  }
  let token = localStorage.getItem("token")
  if(token){
    headers.Authorization=`Bearer ${token}`
  }

  const r=await fetch(API+path,{
    ...options,
    headers
  })

  const data=await r.json().catch(()=>({}))

  if(!r.ok){
    throw new Error(data.error||'Request failed')
  }

  return data
}


async function tryLogin(){
    try{
        const hash = await hashPassword(password.value)
        
        const result=await request('/auth',{
            method:'POST',
            body:JSON.stringify({
                email:email.value,
                pwd:hash
            })
        })
        console.log("result",result)
        let token=typeof result==='string' ? result : result
        localStorage.setItem('token',token)
        alert('Logged in')
        

    }
    catch(e){
        alert(e.message)
    }
}


</script>