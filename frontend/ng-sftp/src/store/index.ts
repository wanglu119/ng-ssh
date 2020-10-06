import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

const state = {
  error: {
    code: 0,
    msg: '',
  },
  user: {},
  jwt: '',
  subdomain: "",
  serverUrl:'',
  sftpSelected:[],
  tabs: []
};

const getters = {
  // sidebar: state => {
  //   let r = [];
  //   state.sidebar.forEach(val => {
  //     if (val.role) {
  //       if (val.role === state.user.role) {
  //         r.push(val);
  //       }
  //     } else {
  //       r.push(val);
  //     }
  //   });
  //   return r;
  // }
};

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
  setTab: (state: any, value: any) => {
    state.tabs.push(value)
    if(value.type === 'sftp') {
      state.sftpSelected = []
    }
  },
  removeTab: (state: any, id: any) => {
    for(const t in state.tabs) {
      if(state.tabs[t].id === id) {
        if(state.tabs[t].type === 'sftp') {
          state.sftpSelected = []
        }
        state.tabs.splice(t,1)
        break
      }
    }
  },
  resetSftpSelected(state: any, sshConfigName: string) {
    state.sftpSelected = []
  },
  removeSftpSelected(state: any, index: any) {
    const i = state.sftpSelected.indexOf(index)
    state.sftpSelected.splice(i,1)
  },
  addSftpSelected(state: any, index:any) {
    state.sftpSelected.push(index)
  },
}


export default new Vuex.Store({
  state: state,
  mutations: mutations,
  modules: {}
});
