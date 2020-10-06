import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const state = {
  error: {
    code: 0,
    msg: '',
  },
  user: {},
  jwt: '',
  subdomain:'',
  serverUrl:'',
}

const mutations = {
  setError: (state: any, value: any) => {
    state.error = value
  },
  setJwt: (state: any, value: any) => {
    state.jwt = value
  },
  setUser: (state: any, value: any) => {
    state.user = value
  },
  setSubdomain: (state: any, value: any) => {
    state.subdomain = value
  },
  setServerUrl: (state: any, value: any) => {
    state.serverUrl = value
  },
}

export default new Vuex.Store({
  state: state,
  mutations: mutations,
  actions: {
  },
  modules: {
  }
})
