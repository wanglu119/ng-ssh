<template>
  <v-main>
    <SidebarMenu />

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
            <sftpTab :conf="item" :ref="item.id" />
          </v-card>
          <v-card height="calc(100vh - 120px)" v-else-if="item.type == 'scp'">
            <scpTab :conf="item" />
          </v-card>
        </v-tab-item>
      </v-tabs-items>
    </v-card>

    <input
      style="display:none"
      type="file"
      id="upload-input"
      @change="uploadInput($event)"
      multiple
    />
    <input
      style="display:none"
      type="file"
      id="upload-folder-input"
      @change="uploadInput($event)"
      webkitdirectory
      multiple
    />
  </v-main>
</template>

<script>
// import { mapGetters, mapState, mapMutations } from 'vuex'
import SidebarMenu from "@/components/layout/SidebarMenu";
import sftpTab from "@/components/sftp/sftpTab";
import scpTab from "@/components/scp/scpTab";
import { mapState, mapMutations } from "vuex";

import * as upload  from '@/utils/upload'

export default {
  props: {
    source: String
  },
  components: { SidebarMenu, sftpTab, scpTab },
  data: () => ({
    title: process.env.VUE_APP_HOME_TITLE,
    drawer: null,
    tab: null,
    ttys: [],
    rules: {
      required: value => !!value || "Required."
    }
  }),
  created() {
    this.$vuetify.theme.dark = false;
  },
  computed: {
    ...mapState(["tabs","currTabConf","sftpReload"]),
    ...mapState("sftp",["req"])
  },
  watch: {
    tab: function() {
      if (this.tab !== undefined) {
        this.setCurrTabConf(this.tabs[this.tab]);
        this.resetSelected();
      }
    },
    sftpReload(val) {
      if(val) {
        this.$refs[this.currTabConf.id][0].getSftpFileList()
        this.setSftpReload(!val)
      }
    }
  },
  methods: {
    ...mapMutations(["removeTab", "setCurrTabConf","setSftpReload"]),
    ...mapMutations("sftp", ["resetSelected","removeReq","closeHovers","showHover"]),
    closeTab: function(tabId) {
      this.removeTab(tabId);
      this.removeReq(tabId)
    },

    getTtyServer: function() {
      // const url = 'ws://'+store.state.user.serverUrl.replace('http://','')+'/ngssh_api/ssh/ws?subdomain='+store.state.user.subdomain+"&auth="+store.state.jwt
      // return url
    },
    uploadInput(event) {
      this.closeHovers()
      const files = event.currentTarget.files
      const folder_upload = files[0].webkitRelativePath !== undefined && files[0].webkitRelativePath !== ''
      if (folder_upload) {
        for (let i = 0; i < files.length; i++) {
          const file = files[i]
          files[i].path = file.webkitRelativePath
        }
      }

      const conflict = upload.checkConflict(files, this.req[this.currTabConf.id].file_list)
      if (conflict) {
        this.showHover({
          prompt: 'replace',
          confirm: (event) => {
            event.preventDefault()
            this.closeHovers()
            upload.handleFiles(files, this.currTabConf.currDir, true)
          }
        })

        return
      }

      upload.handleFiles(files, this.currTabConf.currDir)
    }
  },
};
</script>
