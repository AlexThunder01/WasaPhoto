<script setup>
import { RouterLink, RouterView } from 'vue-router'
import Navbar from './components/Navbar.vue';
</script>
<script>
export default {
  data() {
	return {
		login: false,
		id: localStorage.getItem("token"),
		list: null
	}
  },
  methods: {
	updateUserLogin(status) {
		this.login = status
	},
	getProfile(profileId) {
		this.id = profileId
		this.$router.replace("/profile")
	},
	searchAllUsers(users) {
		this.list = users
		this.$router.replace("/search")
	}
  },
  mounted() {
	if (localStorage.getItem('token')){
		this.login = true	
	}
  },
  emits: ['userIsLogged', 'userLogout'],
  watch: {
	login(newvalue) {
		if(newvalue) {
			document.body.style.background = "white"
		} else {
			document.body.style.background = "linear-gradient(to bottom,#00d9ff,#22d9f9, #ea991e)"
		}
	}
  }
}
</script>

<template>
  <div>
	<div>
	  <main>
		<Navbar v-if="login" 
		@userLogout = "updateUserLogin" 
		@userProfile="getProfile"
		@searchAllUsers="searchAllUsers"
		:profileId="id"
		></Navbar>
        <RouterView @userIsLogged = "updateUserLogin"
		:profileId="id"
		:usersList="list" 
		@userProfile="getProfile"
		></RouterView>
					
	  </main>
	</div>
  </div>

</template>

<style>
	body {
    margin: 0;
    padding: 0;
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
  }
</style>

