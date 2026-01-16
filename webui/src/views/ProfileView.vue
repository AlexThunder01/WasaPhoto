<script>
export default {
	data: function() {
		return {
			profileData:[],            // Array che conterrà i dati richiesti per la visualizzazione del profilo
      username: "",              // Nome utente del profilo
      newUsername: "",           // Nuovo nuovo utente
      numFollowers: 0,           // Numero dei followers
      numFollowed: 0,            // Numero dei following
      numPhoto: 0,               // Numero dei Photo
      isMyProfile: true,         // Booleana per controllare se stiamo sul profilo dell'utente loggato
      isBan: false,              // Booleana per controllare se siamo sul profilo di un utente bannato
      isFollowed: false,         // Booleana per controllare se siamo sul profilo di un utente seguito
      profileIcon: "",           // Path dell'icona del profilo 
      change_icon: true,         // Booleana usata per accorgersi quando viene cambiata l'icona del profilo
      error: false,              // Booleana gestione degli errori
      msg: ""                    // Testo di ritorno dell'errore
		}
	},
  props: ["profileId", "usersList"],

  emits: ['userIsLogged', 'userProfile'],

	methods: {

      // ### BAN SECTION ###

      async banUser() {
        // Mando una richiesta al backend per bannare l'utente
        try {
          let response = await this.$axios.put("/users/" + localStorage.getItem("token") + "/ban/" + this.profileId)
          return response
        }
        catch(e) {
          console.log(e)
        }
      },


      async unbanUser() {
         // Mando una richiesta al backend per sbannare l'utente
        try {
          let response = await this.$axios.delete("/users/" + localStorage.getItem("token") + "/ban/" + this.profileId)
          return response
        }
        catch(e) {
          console.log(e)
        }
      },


      async updateBan() {
        // Si occupa di aggiornare il ban
        if(!this.isBan) {
          try {
            const response = await this.banUser()
            if(response.status == 204) {
              this.isBan = true
              this.numFollowers -= this.isFollowed ? 1 : 0;
              this.isFollowed = false
              await this.updateProfileData() 
            }
          } 
          catch (e) {
            console.log(e)
          }
        }
        else {
          try {
            const response = await this.unbanUser()
            if(response.status == 204) {
              this.isBan = false
              await this.updateProfileData() 
            }
          } 
          catch (e) {
            
          }
        }
      },


      async checkIfBanned() {
      // Controllo se l'utente del profilo attuale è bannato
        if (!this.isMyProfile) {
          try {
            const response = await this.getProfileData()
            this.profileData = response.data
            
            this.isBan = response.status == 206 ? true : false

          } 
          catch (e) {
            console.log(e)
          }
        } 
      },


      // ### FOLLOW SECTION ###


      async followUser() {
        // Mando una richiesta al backend per seguire l'utente
        try {
            let response = await this.$axios.put("/users/" + localStorage.getItem("token") + "/followed/" + this.profileId)
            return response
          }
          catch(e) {
            console.log(e)
          }
      },


      async unfollowUser() {
        // Mando una richiesta al backend per smettere di seguire l'utente
        try {
            let response = await this.$axios.delete("/users/" + localStorage.getItem("token") + "/followed/" + this.profileId)
            return response
          }
          catch(e) {
            console.log(e)
          }
      },


      async updateFollow() {
        // Si occupa di aggiornare il follow

        // Eseguiamo questo blocco di codice se vogliamo seguire l'utente del profilo attuale
        if(!this.isFollowed) {
          try {
            const response = await this.followUser()
            if(response.status == 204) {
              this.numFollowers++
              this.isFollowed = true
            }
          } 
          catch (e) {
            console.log(e)
          }
        }

        // Eseguiamo questo blocco di codice se vogliamo smettere di seguire l'utente del profilo attuale
        else {
          try {
            const response = await this.unfollowUser()
            if(response.status == 204) {
              this.numFollowers--
              this.isFollowed = false
            }
          } 
          catch (e) {
            console.log(e)
          }
        }
      },


      checkIfFollowed() {
        // Controllo se l'utente del profilo attuale è seguito
        if (!this.isMyProfile && this.profileData.followers != null) {
          if(this.profileData.followers.includes(+localStorage.getItem("token"))) {
            this.isFollowed = true
          }
        } 
        else {
          this.isFollowed = false
        }
      },



      // ### USERNAME SECTION ###



      async getUsername() {
        // Richiedo lo username associato all'id del profilo
        try {
          let response = await this.$axios.get('/users/' + this.profileId)
          return response
        } 
        catch(e) {
          console.log(e)
        }
      },


      async setUsername() {
        // Eseguo getUsername e se tutto va bene metto il contenuto della risposta in username
        try {
          const response = await this.getUsername()
          if(response.status == 200) {
            this.username = response.data
          }
        } 
        catch(e) {
          console.log(e)
        }
      },


      async changeUsername() {
        /* Cambio lo username solo se il nuovo username ha un lunghezza compresa tra 3 e 16
           ed il nuovo username non esiste gia nel db */
        if(this.newUsername >= "3") { 
            try {
              let response = await this.$axios.put("/users/" + localStorage.getItem("token"), {
                  username:this.newUsername
              })
              this.username = this.newUsername
              this.newUsername = ""
            } 
            catch (e) {
              console.log(e)
              this.handleResponseError(e)
            }
        } 
      },


      // ### PROFILE SECTION ###



	    async getProfileData() {
        // Richiedo al backend i dati necessari per la visualizzazione del profilo
     	 	try {
          let response = await this.$axios.get('/users/' + this.profileId + '/profile')
          return response
        } 
        catch(e) {
        	console.log(e)
      	}
	    },


      async setProfileData() {
        // Eseguo getProfileData e se tutto va bene metto il contenuto della risponse in profileData
     	 	try {
          const response = await this.getProfileData()
          if(response.status == 200) {
            this.profileData = response.data
          }
        } 
        catch(e) {
        	console.log(e)
      	}
	    },


      async updateProfileData() {
        // Aggiorno i dati del profilo
        await this.setProfileData()
        await this.setUsername()
      },



      // ### PHOTO SECTION ###


      async makePhoto() {
        // Si occupa di creare un Photo

        let inputFile = document.getElementById('inputFile');

        // Metto in questa costante l'immagine caricata dall'utente
        const image = inputFile.files[0];

        const reader = new FileReader();
        reader.readAsArrayBuffer(image);

        // Questo blocco di codice verrà eseguito solo se la lettura ha successo
        reader.onload = async ()=> {
          // Mando una richiesta al backend per creare un Photo
          try {
            let response = await this.$axios.post("/users/" + localStorage.getItem("token") + "/photos", reader.result, {
                Headers: {
                  "Content-Type": image.type
                }
            })
          } catch(e) {

          }

          // Rimuovo l'immagine caricata
          inputFile.value = "";

          // Chiudo la modale
          $('#fileModal').modal('hide');

          // Aggiorno i dati del profilo dopo aver creato una Photo
          await this.updateProfileData()
        };

        // Questo blocco di codice verrà eseguito in caso di errore durante la lettura
        reader.onerror = (error)=> {
          console.error('Errore durante la lettura:', error);
        };

      },


      async DeletePhoto(photoId) {
        try {
          let response = await this.$axios.delete("/users/" + localStorage.getItem("token") + "/photos/" + photoId)
          if(response.status == 204) {
            await this.updateProfileData()
          }
        } 
        catch (e) {
          console.log(e)
        }
      },

       
      
      // ### ICON SECTION ###



      getProfileIcon() {
        this.profileIcon = this.getIcon(this.profileData.icon_id)
      },


      async setIcon(iconID) {
        // Metodo che manda una richiesta di tipo PUT per modificare l'icona di un utente
        try {
          let response = await this.$axios.put("/users/" + this.profileId + "/icon", iconID)
          return response
        } 
        catch(e) {
          console.log(e)
        }
      },


      async changeIcon(iconID) {
        // Eseguo setIcon e se tutto va bene aggiorno i dati correnti della pagina
        try {
          let response = await this.setIcon(iconID)
          if(response.status == 204) {
            await this.setProfileData()
            this.change_icon = !this.change_icon
          }
        } 
        catch(e) {
          console.log(e)
        }
      },


      getIcon(iconID) {
        // Restistuice l'indirizzo per fare richiesta di una certa icona
        return __API_URL__+"/icons/"+iconID+"/"
      },


      // ### ERROR SECTION ###

      handleResponseError(e) {
        // Gestisce l'errore per farlo visualizzare all'utente
        // Fa scomparire la visualizzazione dell'errore dopo 10 secondi
        this.error = true;
        this.msg = e.response.data;

        setTimeout(() => {
          this.error = false;
          this.msg = "";
        }, 10000); // 10 secondi
      },

      customError(errorMessage) {
        // Crea uno specifico messaggio di errore per farlo visualizzare all'utente
        // Fa scomparire la visualizzazione dell'errore dopo 10 secondi
        this.error = true;
        this.msg = errorMessage;

        setTimeout(() => {
          this.error = false;
          this.msg = "";
        }, 10000); // 10 secondi
      }
  },

  async mounted() {

    this.isMyProfile = localStorage.getItem("token") == this.profileId                        // Controllo se il profilo che stiamo visualizzando è dell'utente loggato
    await this.checkIfBanned()                                                                // Controllo se il profilo che stiamo visualizzando è di un utente bannato
	  await this.updateProfileData()                                                            // Ottengo i dati del profilo

	},

  watch: {

    profileData() {
      this.numFollowers = this.profileData.followers ? this.profileData.followers.length : 0;   // Controllo il numero dei followers
      this.numFollowed = this.profileData.followed ? this.profileData.followed.length : 0;      // Controllo il numero dei following
      this.numPhoto = this.profileData.photos ? this.profileData.photos.length : 0;             // Controllo il numero dei photo
      this.checkIfFollowed()                                                                    // Controllo se il profilo che stiamo visualizzando è di un utente seguito
      this.getProfileIcon()                                                                     // Ottengo l'icona del profilo
    },
       
    async profileId() {
      await this.updateProfileData()                                                            // Riaggiorno i dati del profilo
      this.isMyProfile = localStorage.getItem("token") == this.profileId                        // Controllo se il profilo che stiamo visualizzando è dell'utente loggato
      await this.checkIfBanned()                                                                // Controllo se il profilo che stiamo visualizzando è di un utente bannato
    }
  }

}
</script>

<template>
  <!-- Titolo del profilo -->
  <div class="text-center mb-4 mt-4">
    <span v-if="isMyProfile" class="d-inline-block align-middle ml-2">
        <button class="btn btn-link p-0" data-toggle="modal" data-target="#changeIconModal" @click="getIcons">
          <i class="fa-solid fa-plus"></i>
        </button>
    </span>
    <span class="d-inline-block align-middle rounded-circle ml-2" style="overflow: hidden; width: 40px; height: 40px;">
        <img :src=profileIcon alt="Profilo utente" class="img-fluid" style="width: 100%;">
    </span>
    <h2 class="d-inline-block align-middle ml-2">{{ username }}</h2>
</div>

  <div class="row mb-4 custom-container">

    <!-- Informazioni utente -->
    <div class="col-md-6">
      <h4>Statistiche Utente</h4>
      <p><strong>photos:</strong> {{ numPhoto }}</p>
      <p><strong>Followers:</strong> {{ numFollowers }}</p>
      <p><strong>Following:</strong> {{ numFollowed }}</p>
    </div>

    <!-- Cambio dello username -->
    <div v-if="isMyProfile" class="col-md-6">
    <h4>Cambia Username</h4>
      <div class="input-group mb-3">
        <input type="text" v-model="newUsername" @input="newUsername = newUsername.toLowerCase().trim()" class="form-control" placeholder="Nuovo Username" @keyup.enter="newUsername.length >= 3 ? changeUsername() : null" maxlength="16">
        <div class="input-group-append">
          <button class="btn" @click="changeUsername" :class="{ 'btn-primary': newUsername.length >= 3, 'btn-secondary': newUsername.length < 3 }" :disabled="newUsername.length < 3">Cambia</button>
        </div>
      </div>
    </div>

    <!-- Pulsanti per follow/unfollow e ban/unban -->
    <div v-else>
      <button v-if="!isBan" class="btn btn-primary ml-5" style="width: 8rem;" @click="updateFollow">{{ isFollowed? "unfollow" : "follow" }}</button>
      <button class="btn btn-primary ml-5" style="background-color: orangered; border-color: orangered; width: 8rem;" @click=updateBan>{{ isBan? "unban" : "ban"}}</button>
    </div>
  </div>

  <!-- finestra di errore -->
  <div class="container-sm">
    <ErrorMsg v-if="error" :msg="msg">

    </ErrorMsg>
  </div>

  <!-- Lista delle photo -->
  <div class="container mt-5 text-center">
  <h3>Ultime photo</h3>

  <!-- Pulsante per selezionare il file -->
  <button v-if="isMyProfile" class="btn btn-primary" data-toggle="modal" data-target="#fileModal">Carica una photo</button>
</div>

<!-- Finestra modale per la selezione del file -->
<div class="modal" id="fileModal" tabindex="-1" role="dialog">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title">Seleziona un file</h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <input id="inputFile" type="file" accept=".jpg, .jpeg, .png">
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal">Chiudi</button>
        <button type="button" class="btn btn-primary" @click="makePhoto">Carica</button>
      </div>
    </div>
  </div>
</div>

<!-- Finestra modale per la selezione delle icone -->
<div class="modal" id="changeIconModal" tabindex="-1" role="dialog">
    <div class="modal-dialog" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title">Seleziona un'icona</h5>
                <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>
            <div class="modal-body">
                <div class="row mb-3">
                    <div v-for="(iconID, index) in 4" class="col-3 text-center icon-container" :key="index">
                        <img :src=getIcon(iconID+4) alt="icon" @click="changeIcon(iconID + 4)">
                    </div>
                </div>
                <div class="row mb-3">
                    <div v-for="(iconID, index) in 4" class="col-3 text-center icon-container" :key="index">
                        <img :src=getIcon(iconID+8) alt="icon" @click="changeIcon(iconID + 8)">
                    </div>
                </div>
                <div class="row mb-3">
                    <div v-for="(iconID, index) in 4" class="col-3 text-center icon-container" :key="index">
                        <img :src=getIcon(iconID+12) alt="icon" @click="changeIcon(iconID + 12)">
                    </div>
                </div>
                <div class="row mb-3">
                    <div v-for="(iconID, index) in 4" class="col-3 text-center icon-container" :key="index">
                        <img :src=getIcon(iconID+16) alt="icon" @click="changeIcon(iconID + 16)">
                    </div>
                </div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" data-dismiss="modal">Chiudi</button>
            </div>
        </div>
    </div>
</div>


<!-- photo -->
<div v-if="!isBan">
  <Photo ref="photo" v-for="photo in profileData.photos" 
	  @updateComments = updateProfileData
    @DeletePhoto = DeletePhoto
    :key="photo.photo_id"
		:author_id="photo.author_id"
		:photo_id="photo.photo_id"
		:comments="photo.comments != null ? photo.comments : []"
		:likes="photo.likes != null ? photo.likes : []"
		:datetime="photo.datetime"
		:username="username"
    :change_icon="change_icon">
	</Photo>
</div>
</template>

<style>
  .custom-container {
      max-width: 800px;
      margin: 0 auto;
      padding: 20px;
      background-color: #fff;
      border-radius: 8px;
      box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
      margin-top: 50px;
    }

    .background-light {
      background-color: #f8f9fa;
    }

    .icon-container {
        position: relative;
        overflow: hidden;
        width: 40px;
        height: 40px;
        border-radius: 50%;
        transition: transform 0.2s ease-in-out;
    }

    .icon-container:hover {
        transform: scale(1.2);
    }
</style>
