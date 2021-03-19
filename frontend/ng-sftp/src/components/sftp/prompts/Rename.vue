<template>
  <div class="text-center">
    <v-dialog :value="dialog" width="400" @click:outside="closeHovers">
      <v-card
        class="mx-auto"
        tile
      >
        <v-card-title class="headline">
          Rename
        </v-card-title>
        <v-divider></v-divider>

        <v-card-text>
        <v-text-field
            v-model="name"
            :rules="rules.noEmpty"
            counter
            label=""
          >
            <template v-slot:label>
              Insert a new name for <code>{{ oldName }}</code>:
            </template>
          </v-text-field>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="green darken-1"
            text
            @click="closeHovers"
          >
            Cancle
          </v-btn>
          <v-btn
            color="green darken-1"
            text
            @click="rename"
          >
            Rename
          </v-btn>
         
        </v-card-actions>
        
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { mapMutations, mapState, mapGetters } from "vuex";

import { files as api } from '@/api'
import url from '@/utils/url'

export default {
  name: "rename",
  data() {
    return {
      dest: '',
      name: this.oldName,
      rules: {
        noEmpty: [val => (val || '').length > 0 || 'This field is required'],
      }
    }
  },
  computed: {
    ...mapState("sftp", ["show","selectedCount","req", "selected"]),
    ...mapGetters("sftp", ["selectedCount"]),
    ...mapState(["currTabConf"]),
    dialog() {
      if (this.show === "rename") {
        return true;
      } else {
        return false;
      }
    },
    oldName() {
      if(this.selectedCount === 1) {
        this.setName(this.req[this.currTabConf.id].file_list[this.selected[0]].name)
        return this.req[this.currTabConf.id].file_list[this.selected[0]].name
      } 
      return ""
    }
  },
  methods: {
    ...mapMutations("sftp", ["closeHovers","showHover"]),
    ...mapMutations(["setSftpReload"]),
    setName(oldName) {
      this.name = oldName
    },
    rename: async function() {
      let oldLink = ''
      let newLink = ''

      oldLink = this.req[this.currTabConf.id].file_list[this.selected[0]].path
      newLink = url.removeLastDir(oldLink) + '/' + encodeURIComponent(this.name)

      try {
        await api.move([{ from: oldLink, to: newLink }])

        this.setSftpReload(true)
      } catch (e) {
        console.log("rename error:", e)
        this.setSftpReload(true)
      }

      this.closeHovers()
    }
  }
};
</script>
