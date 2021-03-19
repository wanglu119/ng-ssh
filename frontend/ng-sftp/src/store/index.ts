import Vue from "vue";
import Vuex from "vuex";

import sftp from "./modules/sftp";
// import upload from "./modules/upload";
const uploadlib = require("./modules/upload")
const upload = uploadlib.default

Vue.use(Vuex);

const state = {
  error: {
    code: 0,
    msg: ""
  },
  user: {},
  jwt: "",
  subdomain: "",
  serverUrl: "",
  tabs: [],
  currTabConf: null,
  sftpReload: false,
};

const mutations = {
  setError: (state: any, value: any) => {
    state.error = value;
  },
  setJwt: (state: any, value: any) => {
    state.jwt = value;
  },
  setUser: (state: any, value: any) => {
    state.user = value;
  },
  setSubdomain: (state: any, value: any) => {
    state.subdomain = value;
  },
  setServerUrl: (state: any, value: any) => {
    state.serverUrl = value;
  },
  setTab: (state: any, value: any) => {
    state.tabs.push(value);
    if (value.type === "sftp") {
      sftp.state.selected = [];
    }
  },
  removeTab: (state: any, id: any) => {
    for (const t in state.tabs) {
      if (state.tabs[t].id === id) {
        if (state.tabs[t].type === "sftp") {
          sftp.state.selected = [];
        }
        state.tabs.splice(t, 1);
        break;
      }
    }
  },
  setCurrTabConf: (state: any, conf: any) => {
    state.currTabConf = conf;
  },
  setSftpReload: (state: any, reload: any) => {
    state.sftpReload = reload
  }
};

export default new Vuex.Store({
  state: state,
  mutations: mutations,
  modules: {
    sftp,
    upload
  }
});
