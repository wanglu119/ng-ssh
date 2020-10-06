<template>
  <v-main>
    <SidebarMenu/>

    <v-card height="100%">
      <v-tabs
        v-model="tab"
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
          
          <v-card height="calc(100vh - 120px)" v-if="item.type == 'sftp'">
            <sftpTab :conf="item"/>
          </v-card>
          <v-card height="calc(100vh - 120px)" v-else-if="item.type == 'scp'">
            <scpTab :conf="item"/>
          </v-card>

        </v-tab-item>
      </v-tabs-items>
    </v-card>

  </v-main>
</template>

<script>
// import { mapGetters, mapState, mapMutations } from 'vuex'
import {procURL} from '@/api/utils'
import SidebarMenu from '@/components/layout/SidebarMenu'
import sftpTab from '@/components/sftp/sftpTab'
import scpTab from '@/components/scp/scpTab'
import { mapState,mapMutations } from 'vuex'

export default {
  props: {
    source: String,
  },
  components: {SidebarMenu,sftpTab,scpTab},
  data: () => ({
    title: process.env.VUE_APP_HOME_TITLE,
    drawer: null,
    tab: null,
    ttys: [],
    rules: {
      required: value => !!value || 'Required.',
    },
  }),
  created() {
    this.$vuetify.theme.dark = false
  },
  computed: {
    ...mapState(['tabs']),
  },
  methods: {
    ...mapMutations(['removeTab','resetSftpSelected']),
    closeTab: function(tabId) {
      this.removeTab(tabId)
    },

    getTtyServer: function() {     
      // const url = 'ws://'+store.state.user.serverUrl.replace('http://','')+'/ngssh_api/ssh/ws?subdomain='+store.state.user.subdomain+"&auth="+store.state.jwt
      // return url
    }
  },
  watch: {
    tab: function() {
      if(this.tab !== undefined) {
        // ttyTabStore.dispatch('setCurrTtyId', this.tabs[this.tab].id)
        this.resetSftpSelected()
      }
    }
  },
}
</script>
