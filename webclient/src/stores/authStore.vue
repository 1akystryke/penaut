<script>
  import { defineStore } from 'pinia'
  export const authStore = defineStore('auth', {
    state: () => ({
      status: false,
      token: "",
      meName: "",
      meId: "",
      API:"",
      WS:""
    }),
    
    actions: {
      getAuth() {
        let token = localStorage.getItem('token')
        let meName = localStorage.getItem("meName")
        let meId = localStorage.getItem("meId")
        const DEFAULT_HOST='localhost:8080'
        this.API=`http://${DEFAULT_HOST}`
        this.WS=`ws://${DEFAULT_HOST}/ws`
        if (token){
          this.status = true
          this.token = token
          this.meName = meName
          this.meId = meId
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