<script>
  export default {
    
    data() {
      return {
        username: "",   // Stringa dove inseriremo lo username che l'utente inserisce nella loginView
        placeholder: "Inserisci uno username"
      };
    },
    
    methods: {

      async login() {
        // Richiama session e se tutto va bene cambio lo stato della pagina
			  try {
				  const response = await this.session()

          if(response.status == 201) {
            localStorage.setItem('token', response.data.userID)
            this.$router.replace("/home")
            this.$emit('userIsLogged', true)
          }
				
			  } 
        catch (e) {
				  console.log(e)
			  }
		  },

      async session() {
        // Metodo per richiedere al backend l'id dello username con cui si fa richiesta
        try {
          let response = await this.$axios.post("/session",{
					  Username: this.username.toLowerCase().trim()
				  });

          return response
        } 
        catch (e) {
          console.log(e)
          this.handleResponseError(e)
        }
      },

      handleResponseError(e) {
        // Gestisce l'errore per farlo visualizzare all'utente
        // Fa scomparire la visualizzazione dell'errore dopo 10 secondi
        this.username = ""
        this.placeholder = e.response.data;
      },

    },

    mounted() {
      if (localStorage.getItem('token')) {
        this.$router.replace("/home")
      }
    },

    emits: ['userIsLogged'],
    
  };
</script>

<template>
  <div class="container login-container">
    <div class="card shadow p-3 mb-5 bg-body rounded border border-primary">
      <h2 class="text-center mb-4">📸 WasaPhoto</h2>
      <form @submit.prevent="username.length >= 3 ? login() : null" class="text-center">
        <div class="form-group">
          <label for="username" class="form-label">Login</label>
          <div class="input-group">
            <span class="input-group-text"><i class="bi bi-person"></i></span>
            <input type="text" class="form-control" v-model="username" @input="username = username.toLowerCase().trim()" id="username" maxlength="16" :placeholder=placeholder required>
          </div>
        </div>
        <button type="submit" class="btn btn-primary btn-block" :disabled="username.length < 3">Accedi</button>
      </form>
    </div>
  </div>

</template>
  
<style>
  
  .login-container {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100vh;
  }
  
  .card {
    width: 300px;
  }
  
  .card-header {
    background-color: #007bff;
    color: white;
  }
  
  .form-group {
    margin-bottom: 20px;
  }
  
  .input-group-text {
    background-color: #0bc6e3;
    color: white;
  }
  
  .btn-primary {
    background-color: #0bc6e3;
    border-color: #0bc6e3;
  }
  
  .btn-primary:hover {
    background-color: #0978af;
    border-color: #0978af;
  }
  
  .link-primary {
    color: #007bff;
    text-decoration: none;
    font-weight: bold;
  }
  
  .link-primary:hover {
    text-decoration: underline;
  }
  </style>
  