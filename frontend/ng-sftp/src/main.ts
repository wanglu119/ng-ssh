import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import vuetify from "./plugins/vuetify";

import './plugins/vue-axios'
import axios from 'axios'

Vue.config.productionTip = false

const authlib = require('@/utils/auth')

async function start () {

  let subdomain  = ''
  if(process.env.VUE_APP_MODE === 'ng') {
    const debug = authlib.getQueryString('debug')
    try{
      
      let auth = authlib.getQueryString('auth')
      console.log('main.js ========>', auth)
      if(auth) {
        authlib.parseToken(auth)
      } else {
        auth = localStorage.getItem('jwt')
        authlib.parseToken(auth)
      }

      let serverUrl = authlib.getQueryString('serverUrl')
      console.log('main.js ->', serverUrl)
      if(serverUrl) {
        try{
          await authlib.healthcheck(serverUrl,auth)
        }catch(err) {
          console.log(err)
          serverUrl = process.env.VUE_APP_SERVER_URL
        }
      } else {
        serverUrl = process.env.VUE_APP_SERVER_URL
      }

      axios.defaults.baseURL = serverUrl

      subdomain = authlib.getQueryString('subdomain')
      
      await authlib.login(serverUrl,subdomain)
    } catch(err) {
      console.log('main.js =======> ',err)
      const url = authlib.getLoginUrl()
      console.log('main.js =======> ', url)
      if(!debug) {
        window.location = url
      }
    }
  } else {
    if(process.env.VUE_APP_SERVER_URL === '') {
      axios.defaults.baseURL = window.location.protocol+'//'+window.location.host
    }
  }

  try{
    const auth = localStorage.getItem('jwt')
    console.log('auth:',auth)
    if(auth) {
      authlib.parseToken(auth)
      await authlib.renew()
    }
  }catch(err) {
    console.log('main.js =======> ',err)
  }

  store.commit('setSubdomain', subdomain)
  store.commit('setServerUrl', axios.defaults.baseURL)

  new Vue({
    router,
    store,
    vuetify,
    render: h => h(App)
  }).$mount('#app')
}

start()
