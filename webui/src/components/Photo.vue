<script>

export default {
  data(){
    return{ 
      imagePath: "",                                   // Path dell'immagine della photo
      likeClass: "far fa-heart pointer likebtn",               // Variabile per gestire il cambio di stile del pulsante di like
      numLikes: Object.keys(this.likes).length,        // Numero dei likes
      isLiked: false,                                  // Booleana per controllare se è stato messo il like alla photo
      showModal: false,                                // Booleana per gestire la finestra modale dei commenti
      commentUsernames: {},                            // Array che contiene gli username dei vari commenti
      likeUsernames: {},
      textComment: "",                                 // Testo del commento
      loggedUserID: localStorage.getItem("token"),     // User_id dell'utente loggato
      authorIconPath: null,
      iconIdModal: {},
      commentPlaceholder: "Scrivi un commento"
    }
  },
  methods:{

    // ### PHOTO SECTION ###

    getImage() {
      // Richiedo il path dell'immagine della photo
      this.imagePath = __API_URL__+"/users/" + this.author_id + "/photos/" + this.photo_id
    },

    deletephoto() {
    // Rimuove la photo
      this.$emit("DeletePhoto", this.photo_id)
    },


    // ### STYLE SECTION ###

    fillHeart() {
      // Cambio lo style del pulsante dei like se premuto
      this.likeClass = this.likeClass === "far fa-heart pointer likebtn" ? "fa-solid fa-heart pointer color" : "far fa-heart pointer likebtn";
    },


    // ### LIKE SECTION ###


    checkIfLiked() {
      // Controllo se l'utente loggato ha messo like alla photo
      let check = Object.values(this.likes).some(item => item.user_id == this.loggedUserID);

      if(check) {
        // Nel caso l'utente ha messo like alla photo cambio lo style del pulsante dei like
        this.likeClass = "fa-solid fa-heart pointer color";
        this.isLiked = true;
      }
    },

    async putLike() {
      // Mando una richiesta al backend per aggiungere il like alla photo
      try {
        let response = await this.$axios.put("/users/" + this.author_id + "/photos/" + this.photo_id + "/likes/" + this.loggedUserID)
        return response
      } 
      catch(e) {
        console.log(e)
      }
    },

    async deleteLike() {
      // Mando una richiesta al backend per rimuovere il like alla photo
      try {
        let response = await this.$axios.delete("/users/" + this.author_id + "/photos/" + this.photo_id + "/likes/" + this.loggedUserID)
        return response
      } 
      catch(e) {
        console.log(e)
      }
    },

    async updateLike() {
    // Gestisce il like alla photo

      // Questo blocco di codice gestisce il like nel caso si vuole mettere like alla photo
      if(!this.isLiked) {
        try {
          // Mando una richiesta al backend per aggiungere il like alla photo
          const response = await this.putLike() 

          // Se tutto è andato correttamente eseguo questo sottoblocco
          if(response.status == 204) {
            this.numLikes++
            this.isLiked = true
            this.$emit("updateComments")
          }
        }
        catch (e) {
          console.log(e)
        }
      }
      // Questo blocco di codice gestisce il like nel caso si vuole rimuovere il like alla photo
      else {
        try {
          // Mando una richiesta al backend per mettere like alla photo
          const response = await this.deleteLike()

          // Se tutto è andato correttamente eseguo questo sottoblocco
          if(response.status == 204) {
            this.numLikes--
            this.isLiked = false
            this.$emit("updateComments")
          }
        } 
        catch (e) {
          console.log(e)
        }
      }
    },

    async getLikeUsernames() {
      // Invio delle richieste al backend per farmi ritornare gli username di chi ha messo like alla photo
      Object.values(this.likes).forEach(async (like) => {  
        try {
            let response = await this.$axios.get("/users/" + like.user_id)
            this.likeUsernames[like.user_id] = response.data;
        } 
        catch(e) {
            console.log(e)
        }
      })
    },


    // ### COMMENT SECTION ###


    async getCommentUsernames() {
      // Invio delle richieste al backend per farmi ritornare gli username degli autori dei commenti della photo
      Object.values(this.comments).forEach(async (comment) => {  
        try {
            let response = await this.$axios.get("/users/" + comment.author_id)
            this.commentUsernames[comment.author_id] = response.data;
        } 
        catch(e) {
            console.log(e)
        }
      })
    },

    async getIconIDModal() {
      // Invio delle richieste al backend per farmi ritornare gli id delle icone dei commenti e dei likes
      Object.values(this.comments).forEach(async (comment) => {  
        try {
            let response = await this.$axios.get("/users/" + comment.author_id + "/icon")
            this.iconIdModal[comment.author_id] = response.data;
        } 
        catch(e) {
            console.log(e)
        }
      })

      Object.values(this.likes).forEach(async (like) => {  
        try {
            let response = await this.$axios.get("/users/" + like.user_id + "/icon")
            this.iconIdModal[like.user_id] = response.data;
        } 
        catch(e) {
            console.log(e)
        }
      })
    },

    async postComment() {
      // Invio una richiesta al backend per inserire nel db il commento nella photo
      try {
          let response = await this.$axios.post("/users/" + this.author_id + "/photos/" + this.photo_id + "/comments", this.textComment)
          return response
        } 
      catch(e) {
        console.log(e)
        this.handleResponseError(e)
      }
    },

    async commentPhoto() {
      // Aggiungo un commento 

      try {
        // Invio una richiesta al backend per inserire nel db il commento nella photo
        const response = await this.postComment()

        // Se tutto è andato correttamente eseguo questo sottoblocco
        if(response.status == 201) {
          this.$emit("updateComments")
          this.textComment = ""
        }
      } 
      catch (e) {
        console.log(e)
      }
    },

    async deleteComment(id) {
    // Invio una richiesta al backend per rimuovere dal db il commento dalla photo
      try {
        let response = await this.$axios.delete("/users/" + this.author_id + "/photos/" + this.photo_id + "/comments/" + id)
        return response
      } 
      catch(e) {
        console.log(e)
      }
    },

    async removeComment(id) {
    // Rimuove un commento
      
      try {
        // Invio una richiesta al backend per rimuovere dal db il commento dalla photo
        const response = await this.deleteComment(id)

        // Se tutto va bene eseguo questo sottoblocco
        if(response.status == 204) {
          this.$emit("updateComments")
        }
      } 
      catch (e) {
        console.log(e)
      }
    },

    async getUserIconID(userID) {
      // Invio una richiesta al backend per ottenere l'id dell'icona di uno user
      try {
        let response = await this.$axios.get("/users/" + userID + "/icon")
        return response
      }
      catch(e) {

      }
    },

    getIcon(iconID) {
      // Creo il path per richiedere un'icona tramite il suo id
      return __API_URL__+"/icons/"+iconID+"/"
    },
    
    async getAuthorIcon() {
      // Richiedo l'icona di uno user al server e se tutto va bene la cambio
      try {
        const response = await this.getUserIconID(this.author_id)
        
        if(response.status == 200) {
          this.authorIconPath = this.getIcon(response.data)
          
        }
      } 
      catch (e) {
        console.log(e)        
      }
    },

    toggleFullScreen() {
      // Metodo che gestisce il fullscreen
      const elem = document.getElementById(this.imagePath);
      if (!document.fullscreenElement) {
        elem.requestFullscreen().catch(err => {
          console.error(`Impossibile entrare in modalità fullscreen: ${err.message}`);
        });
      } else {
        document.exitFullscreen();
      }
    },

    handleResponseError(e) {
        // Gestisce l'errore per farlo visualizzare all'utente
        this.textComment = ""
        this.commentPlaceholder = e.response.data
        console.log(this.commentPlaceholder)
      },
  },
  
  async mounted() {
    this.getImage();                  // Ottengo l'immagine della photo
    this.checkIfLiked();              // Controllo se l'utente loggato ha messo like alla photo
    await this.getCommentUsernames()        // Ottengo gli usernames dei commenti nella photo
    await this.getLikeUsernames()        // Ottengo gli usernames di chi ha messo like nella photo
    await this.getAuthorIcon()
    await this.getIconIDModal()
  },
  watch: {
    async comments() {
      await this.getCommentUsernames()
      await this.getIconIDModal()
    },
    async likes() {
      await this.getLikeUsernames()        // Ottengo gli usernames di chi ha messo like nella photo
    },
    async change_icon() {
      await this.getAuthorIcon()
    },
    async author_id() {
      await this.getAuthorIcon()
      await this.getIconIDModal()
    },
    
  },

  props: ["author_id","likes","comments","datetime","photo_id", "username", "change_icon"],

  emits: ['updateComments', "DeletePhoto"],
}
</script>

<template>
<div class="spacer"></div>

<div class="container" style="max-width: 30%;">

  <!-- photo-->
  <div class="photo-container fixed-height">

  <!-- Header della photo -->
  <div class="photo-header d-flex align-items-center justify-content-between">
    <div class="d-flex align-items-center">
      <img :src=authorIconPath alt="Profilo utente" class="rounded-circle">
      <div class="ml-2">
        <h6 class="mb-0">{{ username }}</h6>
        <small class="text-muted">{{ datetime }}</small>
      </div>
    </div>
    <i v-if="loggedUserID == author_id" class="far fa-light fa-trash-can pointer deletephotobtn" @click="deletephoto"></i>
  </div>


    <!-- Contenuto della photo -->
    <div class="photo-content">
      <img :id="imagePath" :src=imagePath alt="Immagine della photo" class="img-fluid">
    </div>

    <!-- Icone di Interazione -->
    <div class="d-flex justify-content-between mt-3">
      <div>
        <i :class="likeClass" @click="fillHeart();updateLike();"></i> <span class="pointer likesnumber" data-toggle="modal" :data-target="'#n'+this.photo_id">{{ numLikes }}</span>
      </div>
      <div>
        <i class="far fa-comment pointer commentbtn" data-toggle="modal" :data-target="'#m'+this.photo_id"></i> {{ Object.keys(this.comments).length }}
      </div>
      <div>
        <i class="fas fa-expand pointer fullscreen" @click="toggleFullScreen"></i>
      </div>
    </div>
 
  <!-- Finestra modale dei commenti -->
  <div class="modal fade" :id="'m'+this.photo_id" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel" aria-hidden="true">
    <div class="modal-dialog" role="document">
      <div class="modal-content">

        <!-- Header della modale -->
        <div class="modal-header">
          <h5 class="modal-title" id="exampleModalLabel">Commenti</h5>
          <button type="button" class="close" data-dismiss="modal" aria-label="Chiudi">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>

        <!-- Body della modale -->
        <div class="modal-body">
          <div v-for="(comment,index) in this.comments" class="photo-header" :key="index">
            <img :src=getIcon(iconIdModal[comment.author_id]) alt="Profilo utente" class="rounded-circle">
            <div>
              <h6 class="mb-0">{{ commentUsernames[comment.author_id]}}</h6>
              <small class="text-muted">{{ comment.datetime }}<i v-if="loggedUserID == comment.author_id" class="far fa-light fa-trash-can mx-2 pointer trashcommentbtn" @click="removeComment(comment.comment_id)"></i></small>
              <p>{{ comment.text}}</p>
            </div>
          </div>
        </div>

        <!-- Footer della modale -->
        <div class="modal-footer">
          <div class="input-group mb-3">
            <input type="text" v-model="textComment" class="form-control" :placeholder="commentPlaceholder" aria-label="Recipient's username" aria-describedby="basic-addon2" @keyup.enter="textComment != '' ? commentPhoto() : null">
            <div class="input-group-append">
              <button class="btn btn-outline-secondary" type="button" @click="commentPhoto" :disabled="textComment == ''">Invia commento</button>
            </div>
          </div>
          <button type="button" class="btn btn-secondary" data-dismiss="modal">Chiudi</button>
        </div>


      </div>
    </div>
  </div>

   <!-- Finestra modale dei likes -->
   <div class="modal fade" :id="'n'+this.photo_id" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel" aria-hidden="true">
    <div class="modal-dialog" role="document">
      <div class="modal-content">

        <!-- Header della modale -->
        <div class="modal-header">
          <h5 class="modal-title" id="exampleModalLabel">Likes</h5>
          <button type="button" class="close" data-dismiss="modal" aria-label="Chiudi">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>

        <!-- Body della modale -->
        <div class="modal-body">
          <div v-for="(like,index) in this.likes" class="photo-header" :key="index">
            <img :src=getIcon(iconIdModal[like.user_id]) alt="Profilo utente" class="rounded-circle">
            <div>
              <h6 class="mb-0">{{ likeUsernames[like.user_id]}}</h6>
            </div>
          </div>
        </div>


      </div>
    </div>
  </div>

  </div>
</div>

</template>

<style>

  .spacer {
      height: 20px; /* Altezza dello spazio prima del div della photo */
  }

  .photo-container {
    border: 1px solid #ddd;
    border-radius: 8px;
    margin: 20px 0; /* Margini superiori e inferiori a 20px, 0 ai lati */
    padding: 15px;
    background-color: #f8f9fa; /* Cambia il colore di sfondo della photo */
  }

  .photo-header {
    display: flex;
    align-items: center;
    margin-bottom: 10px;
  }

  .photo-header img {
    border-radius: 50%;
    margin-right: 10px;
  }

  .photo-content img {
    width: 100%;
    border-radius: 8px;
    margin-bottom: 10px;
  }

  .pointer {
    cursor: pointer;
  }

  .color, .color::after, .color::before {
    color: red;
  }

  .photo {
    margin-bottom: 20px;
  }

  .commentbtn:hover {
    color: green;
  }

  .likesnumber:hover {
    color: rgb(192, 192, 192);
  }

  .deletephotobtn:hover, .likebtn:hover, .trashcommentbtn:hover {
    color:red;
  }
  
  .fullscreen:hover {
    color: #acdfff;
  }



</style>