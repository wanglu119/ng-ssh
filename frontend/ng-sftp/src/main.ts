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
  const debug = authlib.getQueryString('debug')
  let auth = authlib.getQueryString('auth')

  try {
    if(auth) {
      authlib.parseToken(auth)
    } else {
      auth = localStorage.getItem('jwt')
      authlib.parseToken(auth)
    }
  } catch(err) {
    console.log('main.ts auth error: ',err)
    const url = authlib.getLoginUrl()
    console.log('main.js =======> ', url)
    if(!debug) {
      window.location = url
    }
  }

  if(process.env.VUE_APP_MODE === 'ng') {
    try{
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
      console.log('ng mode error: ',err)
      const url = authlib.getLoginUrl()
      console.log('main.js =======> ', url)
      if(!debug) {
        window.location = url
      }
    }
  } else {
    if(process.env.VUE_APP_SERVER_URL === '') {
      axios.defaults.baseURL = window.location.protocol+'//'+window.location.host
    } else {
      axios.defaults.baseURL = process.env.VUE_APP_SERVER_URL
    }
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
