<template>
  <v-app>
    <router-view/>
    <Snackbar/>
  </v-app>
</template>

<script>

import {mapState} from 'vuex'
import Snackbar from '@/components/msgs/Snackbar'

export default {
  name: 'App',
  components: {Snackbar,},
  data: () => ({
  }),
  created() {
    this.$vuetify.theme.dark = true
  },
  computed: {
    ...mapState(['jwt']),
  },
  mounted () {
    const loading = document.getElementById('loading')
    if(loading !== null) {
      loading.classList.add('done')

      setTimeout(function () {
        if(loading.parentNode!== null) {
          loading.parentNode.removeChild(loading)
        }
      }, 200)
    }
    
    if(this.jwt.length > 0) {
      console.log(this.$router)
      if(this.$router.currentRoute.name !== 'Home'){
        this.$router.push({
          name: 'Home'
        })
      }
    } else {
      if(this.$router.currentRoute.name && this.$router.currentRoute.name !== 'SignIn'){
        if(process.env.VUE_APP_LOGIN_URL === '') {
          window.location = window.location.protocol+'//'+window.location.host+'/SignIn'
        } else {
          window.location = process.env.VUE_APP_LOGIN_URL
        }
      }
    }
  }
}
</script>
