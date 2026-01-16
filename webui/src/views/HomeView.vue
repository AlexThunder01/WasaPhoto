<script>
export default {
	
	data: function() {
		return {
			photos: []    // Array che conterrà gli oggetti photos
		}
	},

	methods: {
		async loadPhotos() {
			// Faccio una richiesta per ottenere le photos
     	 	try {
       			let response = await this.$axios.get('/users/' + localStorage.getItem('token') + '/home')
        		this.photos = response.data
      		} 
			catch(e) {
        		console.log(e)
      		}
    	},
	},

	async mounted() {
		await this.loadPhotos()
	},

	emits: ['userIsLogged', 'userProfile'],

	props: ['profileId', 'usersList'],
}
</script>

<template>
	<div v-if=!photos>
		<h1 class="text-center mt-3">Start Following Someone</h1>	
	</div>

	<Photo v-else v-for="(photo,index) in photos" 
	    @updateComments = loadPhotos
		:key="index"
		:author_id="photo.author_id"
		:photo_id="photo.photo_id"
		:comments="photo.comments != null ? photo.comments : []"
		:likes="photo.likes != null ? photo.likes : []"
		:datetime="photo.datetime"
		:username="photo.username"
		:change_icon="null">
	</Photo>
	
</template>

<style scoped>
    
</style>
