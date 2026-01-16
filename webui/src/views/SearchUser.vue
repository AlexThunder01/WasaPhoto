  <script>
  export default {
    data() {
        return {
            usersIconID: {}     //Conterrà l'iconID degli user che sono stati trovati nella ricerca
        };
    },

    props: ["usersList"],

    methods: {
        getProfile(userID) {
            // Chiedo al componente padre di cambiare la pagina col profilo dell'utente richiesto
            this.$emit("userProfile", userID)
        },

        getIcon(iconID) {
            // Path per la richiesta di un icona avendo l'iconID
            return __API_URL__+"/icons/"+iconID+"/"
        },

        async getUsersIconID() {
         // Invio delle richieste al backend per farmi ritornare gli IconID degli user che ho trovato con la ricerca
            if(this.usersList != null) {
                Object.values(this.usersList).forEach(async (user) => {  
                    try {
                        let response = await this.$axios.get("/users/" + user.userID + "/icon")
                        this.usersIconID[user.userID] = response.data;
                    } 
                    catch(e) {
                        console.log(e)
                    }
                })
            }
        },
    },

    async mounted(){
        await this.getUsersIconID()
    },
    watch: {
        async usersList() {
            await this.getUsersIconID()
        }
    }

  };

  </script>


<template>

    <div class="page-container">

        <!-- Colonna sinistra -->
        <div class="column"></div>

        <!-- Colonna centrale-->
        <div class="column">
            <div class="myContainer" v-if="usersList != null">
                <h1 class="title">Utenti Trovati</h1>
                <ul class="name-list">
                    <li v-for="(user, index) in usersList" class="photo-header" @click="getProfile(user.userID)" :key="index">
                        <img :src=getIcon(usersIconID[user.userID]) alt="Profilo utente" class="rounded-circle">
                        <div>
                            <h6 class="mb-2 ml-2">{{ user.username }}</h6>
                        </div>
                    </li>
                </ul>
            </div>
            <div v-else class="myContainer">
                <h1 class="title">Nessun Utente è stato trovato</h1>
            </div>
        </div>

        <!-- Colonna destra -->
        <div class="column"></div> 
        
    </div>

</template>
  
  
  <style scoped>

    .column:nth-child(2) {
        background-color: #ffffff; 
    }
    .title {
        font-size: 24px;
        margin-bottom: 20px;
    }
    .photo-header img {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        margin-right: 10px;
    }
    .photo-header:hover {
        background-color: whitesmoke; 
    }
    .page-container {
        display: flex;
    }
    .column {
        flex: 1;
    }
    .myContainer {
        height: 100%;
        overflow-y: auto;
    }
    .name-list {
        list-style: none;
        padding: 0;
    }
    .photo-header {
        display: flex;
        align-items: center;
        padding: 10px;
        border-bottom: 1px solid #ccc;
        border: 1px solid #ddd;
    }
    .rounded-circle {
        border-radius: 50%;
        margin-right: 10px;
    }
      
  </style>
  