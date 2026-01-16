<script>
export default {
  data(){
    return{
      searchQuery: "",               // Stringa che conterrà il valore della barra di ricerca
      showSuggestions: false,        // Booleana che deciderà se mostrare o no la lista degli utenti simili nella barra di ricerca
      suggestedUsers: [],            // Array che conterrà (se esistono) gli usernames simili alla stringa immessa nella barra di ricerca
      highlightedUser: null,         // Permette l'illuminazione degli user al passaggio del mouse
      blurTimer: null,               // Variabile che conterrà il timer che ritarda la chiusura della lista degli usernames
    }
  },
  methods:{

    // ### SUGGESTED USERS LIST SECTION ###

    delayedCloseSuggestions() {
      // Ritarda l'attivazione di closeSuggestions di 200 millisecondi
      this.blurTimer = setTimeout(() => {
        this.closeSuggestions();
      }, 200);
    },

    closeSuggestions () {
      // Chiudo la lista degli utenti suggeriti
      this.showSuggestions = false
    },

    openSuggestions() {
      // Apre sempre la lista dei suggerimenti quando clicchi sulla barra di input
      this.showSuggestions = true;

      // Esegui la ricerca solo se ci sono elementi nella barra di ricerca
      if (this.searchQuery.trim() !== '') {
        this.searchUsers();
      }
    },
    
    async searchUsers() {
      // Richiedo la lista degli usernames simili alla stringa immessa nella barra di ricerca
      if(this.searchQuery.length > 0) {
        try {
          const response = await this.getSuggestedUsernames()
          if(response.status == 200) {
            this.suggestedUsers = response.data
            this.showSuggestions = true;
          }
        } 
        catch (e) {
          console.log(e)
        }
      } 
      else {
        this.suggestedUsers = []
        this.showSuggestions = false
      }
    },

    async getSuggestedUsernames() {
      // Invio una richiesta per ottenere la lista degli usernames simili
      try {
        let response = await this.$axios.get("/users", {
          params: {
            users: this.searchQuery
          }
        })
        return response
      }
      catch (e) {
        console.log(e)
      }
    },


    // ### PROFILE CHANGE LOGIC SECTION ###

    selectUser(user) {
      // Cambio la pagina corrente con quella dell'utente selezionato nella lista
      this.$emit("userProfile", user.userID)
      this.showSuggestions = false;
      clearTimeout(this.blurTimer);
    },

    viewAllUsers() {
      // Apro la pagina con tutti gli usernames simili alla stringa immessa nella barra di ricerca
      this.$emit("searchAllUsers", this.suggestedUsers)
      this.showSuggestions = false;
    },


    // ### HOME,PROFILE,LOGOUT BUTTONS SECTION ###


    logout() {
      // Effettuo il logout
      localStorage.removeItem("token")
      this.$emit("userLogout", false)
      this.$router.replace("/login")
    },

    home() {
      this.$router.replace("/home")
    },

    profile() {
      this.$emit("userProfile", localStorage.getItem("token"))
    },


    // ### STYLE SECTION ###


    highlightUser(user) {
      // Illumina l'elemento al passaggio del mouse
      this.highlightedUser = user;
    },

    unhighlightUser() {
      // Disattiva l'illuminazione al mouseleave
      this.highlightedUser = null;
    },
  },

  emits: ['userLogout', "userProfile", "searchAllUsers"],

  props: ["profileId"]
}
</script>

<template>
   <nav class="navbar navbar-expand-lg navbar-dark sticky-top">

    <!-- Logo NavBar-->
    <a class="navbar-brand" href="#">WasaPhoto</a>
    
    <div class="collapse navbar-collapse" id="navbarNav">
      
      <!-- Barra di ricerca e pulsante per la ricerca degli usernames-->
      <form class="form-inline mx-auto my-2 my-lg-0 position-relative" :disabled="searchQuery.length == 0">
      <input v-model="searchQuery" @input="searchUsers" @blur="delayedCloseSuggestions" class="form-control mr-sm-2" type="search" placeholder="Cerca utenti" aria-label="Search">
      <button class="btn btn-outline-light my-2 my-sm-0" type="submit" @click="openSuggestions(); viewAllUsers()" :disabled="searchQuery.length == 0">Cerca</button>
        <div v-if="showSuggestions && suggestedUsers != null" class="suggestions position-absolute w-100">
          <ul class="list-group">
            <li v-for="(user, index) in suggestedUsers.slice(0, 6)" :key="user.user_id" class="list-group-item" @click="selectUser(user)" @mouseover="highlightUser(index)" @mouseleave="unhighlightUser">
              {{ user.username }}
            </li>
          </ul>
          <button v-if="suggestedUsers != null" @click="viewAllUsers" class="btn btn-light w-100">Vedi tutti</button>
        </div>
      </form>

    <!-- Tasti HOME PROFILO LOGOUT -->
    <ul class="navbar-nav ml-auto">
      <li class="nav-item active">
        <i class="nav-link pointer" @click="home">Home</i>
      </li>
      <li class="nav-item">
        <i class="nav-link pointer" @click="profile">Profilo</i>
      </li>
      <li class="nav-item">
        <i class="nav-link pointer" @click="logout">Logout</i>
      </li>
    </ul>

    </div>

</nav>
</template>

<style scoped>
     body {
      background-size: cover !important;
      background-color: #e6f7ff;
    }

    .navbar {
      background-color: #1ebbcf; 
    }
    .navbar-dark .navbar-brand {
      color: #ffffff; 
      background-color: unset;
      box-shadow: none;
    }

    .navbar-dark .navbar-toggler-icon {
      background-color: #ffffff; 
    }

    .navbar-dark .navbar-nav .nav-link {
      color: #ffffff; 
      font-style: normal;
    }

    .navbar-dark .navbar-nav .nav-link:hover {
      color: #cce5ff; 
    }

    .navbar-dark .navbar-toggler {
      border-color: #ffffff; 
    }

    .navbar-dark .navbar-toggler:hover {
      background-color: #ffffff;
    }

    .suggestions {
      top: 100%; 
      max-height: 200px;
      overflow-y: auto;
    }

    .list-group-item {
      cursor: pointer;
    }
    .suggestions ul {
      list-style-type: none;
      padding: 0;
      margin: 0;
    }

    .suggestions li {
      padding: 8px;
      cursor: pointer;
    }

    .suggestions button {
      width: 100%;
      padding: 8px;
      border: none;
      background-color: #f0f0f0;
      cursor: pointer;
    }

    .list-group-item:hover {
      background-color: #f8f9fa; 
    }

    .list-group-item.active {
      background-color: #007bff;
      color: #fff;
    }
</style>