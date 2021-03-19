<template>
  <v-main>
    <v-navigation-drawer v-model="drawer" app clipped>
      <v-list dense v-for="item in ttys" :key="item.name">
        <v-list-item>
          <v-list-item-action>
            <deleteBtn v-show="item.name.length>0" :procData="item.name" :procMethod="deleteTty" />
          </v-list-item-action>
          <v-list-item-content @dblclick="createTab(item)">
            <v-list-item-title>
              <v-btn text class="text-none" @click="createTab(item)">{{item.name}}</v-btn>
            </v-list-item-title>
            <v-list-item-subtitle>{{item.cmd}}</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app clipped-left>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-btn icon @click="ttyDialog = true">
        <v-icon>mdi-view-grid-plus</v-icon>
      </v-btn>
      <!--
      <v-toolbar-title
        style="width: 300px"
        class="ml-0 pl-4"
      >
        <span class="hidden-sm-and-down">ng-ssh</span>
      </v-toolbar-title>
      -->
      <v-btn icon @click="openSftp()">
        <v-img :src="imgUrl('/img/sftp.png')" max-height="40" max-width="60" />
      </v-btn>
      <v-spacer></v-spacer>
      <v-btn icon @click="logout()">
        <v-icon large color="red">mdi-power</v-icon>
      </v-btn>
    </v-app-bar>

    <v-card height="100%">
      <v-tabs
        v-model="tab"
        dark
        next-icon="mdi-arrow-right-bold-box-outline"
        prev-icon="mdi-arrow-left-bold-box-outline"
        show-arrows
      >
        <v-tabs-slider color="yellow"></v-tabs-slider>
        <v-tab v-for="item in tabs" :key="item.id" class="text-none">
          {{ item.name }}
          <v-btn icon @click="closeTab(item.id)">
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-tab>
      </v-tabs>

      <v-tabs-items v-model="tab">
        <v-tab-item v-for="item in tabs" :key="item.id">
          <!-- color="red" -->
          <v-card height="calc(100vh - 150px)">
            <Term :ref="item.id" :ttyConf="item" />
          </v-card>
        </v-tab-item>
      </v-tabs-items>
    </v-card>

    <!-- add ssh config -->
    <v-row justify="center">
      <v-dialog v-model="ttyDialog" persistent max-width="600px">
        <v-card>
          <v-card-title>
            <span class="headline">New Ssh</span>
          </v-card-title>
          <v-card-text>
            <v-form ref="tty_from" v-model="ttyValid">
              <v-text-field
                v-model="name"
                :rules="[rules.required,existTyy,rules.notContainPoint]"
                label="*Name"
                required
              ></v-text-field>
              <v-row>
                <v-col cols="12" sm="9">
                  <v-text-field v-model="host" :rules="[rules.required]" label="*Host" required></v-text-field>
                </v-col>
                <v-col cols="12" sm="3">
                  <v-text-field v-model="port" :rules="[rules.required]" label="*Port" required></v-text-field>
                </v-col>
              </v-row>
              <v-text-field v-model="username" :rules="[rules.required]" label="*Username" required></v-text-field>
              <v-text-field v-model="password"  :rules="[rules.required]" label="*Password" required 
                :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'" 
                :type="showPassword ? 'text' : 'password'"
                @click:append="showPassword = !showPassword"
                ></v-text-field>
            </v-form>

            <small>*indicates required field</small>
          </v-card-text>
          <v-card-actions>
            <div class="flex-grow-1"></div>
            <v-btn color="blue darken-1" text @click="ttyDialog = false">Close</v-btn>
            <v-btn color="blue darken-1" text @click="createTty()">Save</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-row>
  </v-main>
</template>

<script>
import Term from "@/api/tty/vue-tty"
import DeleteBtn from "@/components/buttons/Delete"
import {ttyTabStore} from '@/api/tty/tty_store'
import store from '@/store'
import {procURL,getImageUrl} from '@/api/utils'
import {getQueryString} from '@/utils/auth'

export default {
  props: {
    source: String,
  },
  components: {Term,DeleteBtn},
  data: () => ({
    drawer: null,
    tab: null,
    ttyDialog: false,
    ttyValid: false,
    name: '',
    host: '',
    port: 22,
    username: '',
    password: '',
    tabs: [],
    ttys: [],
    rules: {
      required: value => !!value || 'Required.',
      notContainPoint: value => value && value.indexOf('.')<0 || 'Cannot contain point', 
    },
    showPassword: false,
  }),
  mounted: async function() {
    await this.getTtys()
  },
  methods: {
    getTtys: async function() {
      this.ttys = []
      const r = await procURL('/ngssh_api/ssh_config/getSshConfigs',{
        method: 'GET',
      })
      for(const k in r) {
        const t = r[k]
        t['key'] = k
        this.ttys.push(r[k])
      }
    },
    existTyy: function(name) {
      for(const t in this.ttys) {
        if(this.ttys[t].name === name) {
          return 'Has Exist'
        } 
      }
      return true
    },
    createTab: function(ttyConf) {
      const newTttConf = {...ttyConf}
      newTttConf['id'] = Date.now()
      newTttConf['ttyServer'] = this.getTtyServer()
      this.tabs.push(newTttConf)
      ttyTabStore.dispatch('addTty', newTttConf)
    },
    closeTab: function(tabId) {
      for(const t in this.tabs) {
        if(this.tabs[t].id === tabId) {
          this.tabs.splice(t,1)
          ttyTabStore.dispatch('removeTty',tabId)
          break
        }
      }
    },
    createTty: async function() {
      if(this.$refs.tty_from.validate()) {
        const tty = {
          name: this.name,
          host: this.host,
          port: parseInt(this.port),
          username: this.username,
          password: this.password,
        }

        await procURL('/ngssh_api/ssh_config/addSshConfig', {
          method: 'POST',
          data: tty
        })

        this.getTtys()

        this.ttyDialog = false
        this.$refs.tty_from.reset()
      }
    },
    deleteTty: async function(key) {
      if(key.length <= 0 ){
        return
      }
      console.log(key)
      await procURL('/ngssh_api/ssh_config/deleteSshConfig', {
        method: 'DELETE',
        data: {name:key}
      })
      await this.getTtys()
    },
    getTtyServer: function() {     
      const url = 'ws://'+store.state.serverUrl.replace('http://','')+'/ngssh_api/ssh/ws?subdomain='+store.state.subdomain+"&auth="+store.state.jwt
      return url
    },
    imgUrl(path) {
      return getImageUrl(path)
    },
    openSftp() {
      const debug = getQueryString('debug')
      let url = process.env.VUE_APP_SFTP_URL+'?subdomain='+store.state.subdomain+'&serverUrl='+store.state.serverUrl+'&auth='+store.state.jwt
      if(debug) {
        url = url+'&debug='+debug
      }
      window.open(url)
    },
    logout() {
      localStorage.setItem('jwt', '')
      if(process.env.VUE_APP_LOGIN_URL === '') {
        window.location = window.location.protocol+'//'+window.location.host+'/SignIn'
      } else {
        window.location = process.env.VUE_APP_LOGIN_URL
      }
    }
  },
  watch: {
    tab: function() {
      if(this.tab !== undefined) {
        ttyTabStore.dispatch('setCurrTtyId', this.tabs[this.tab].id)
      }
    }
  },
}
</script>
