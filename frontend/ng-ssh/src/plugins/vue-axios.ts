import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'

axios.defaults.baseURL = process.env.VUE_APP_AXIOS_BASE_URL;

Vue.use(VueAxios, axios)