<script>
  import { defineStore } from 'pinia'
  export const authStore = defineStore('auth', {
    state: () => ({
      status: false,
      token: ""
    }),
    
    actions: {
      getAuth() {
        let token = localStorage.getItem('token')
        if (token){
          this.status = true
          this.token = token
        }else{
          this.status = false
          this.token = ""
        }
        
      },
      breakAuth(){
        localStorage.setItem("token",'')
        this.status = false
        this.token = ""
      }
    }
  })
</script>