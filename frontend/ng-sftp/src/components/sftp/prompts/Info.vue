<template>
  <div class="text-center">
    <v-dialog :value="dialog" width="400" @click:outside="closeHovers">
      <v-card
        class="mx-auto"
        tile
      >
        <v-card-title class="headline">
          File information
        </v-card-title>
        <v-divider></v-divider>

        <v-card-text>
          <p v-if="selected.length > 1">{{ selectedCount }} files selected.</p>

          <p class="break-word" v-if="selected.length < 2"><strong>Display Name:</strong> {{ name }}</p>
          <p v-if="!dir || selected.length > 1"><strong>Size:</strong> <span id="content_length"></span> {{ humanSize }}</p>
          <p v-if="selected.length < 2"><strong>Last Modified:</strong> {{ humanTime }}</p>

          <template v-if="dir && selectedCount === 0">
            <p><strong>Number of files:</strong> {{ req[currTabConf.id].ls_meta.file_count }}</p>
            <p><strong>Number of directories:</strong> {{ req[currTabConf.id].ls_meta.dir_count}}</p>
          </template>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="green darken-1"
            text
            @click="closeHovers"
          >
            Ok
          </v-btn>
         
        </v-card-actions>
        
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { mapMutations, mapState, mapGetters } from "vuex";
import filesize from 'filesize'
import moment from 'moment'

import { files as api } from '@/api'

export default {
  name: "info",
  data() {
    return {
      dest: '',
    }
  },
  computed: {
    ...mapState("sftp", ["show","req", "selected"]),
    ...mapGetters("sftp", ["selectedCount"]),
    ...mapState(["currTabConf"]),
    dialog() {
      if (this.show === "info") {
        return true;
      } else {
        return false;
      }
    },
    dir: function () {
      if(this.selectedCount === 1) {
        return this.req[this.currTabConf.id].file_list[0].is_dir
      } 
      return true
    },
    humanSize: function () {
      if (this.selectedCount === 0) {
        return 0
      }
      if (this.selectedCount === 1) {
        return filesize(this.req[this.currTabConf.id].file_list[0].size)
      }

      let sum = 0

      for (const selected of this.selected) {
        sum += this.req[this.currTabConf.id].file_list[selected].size
      }

      return filesize(sum)
    },
    humanTime: function () {
      if(this.selectedCount === 0) {
        return moment(this.req[this.currTabConf.id].ls_meta.time).fromNow()
      }
      return moment(this.req[this.currTabConf.id].file_list[0].time).fromNow()
    },
    name: function () {
      if(this.selectedCount === 0) {
        return this.req[this.currTabConf.id].ls_meta.name
      } else {
        return this.req[this.currTabConf.id].file_list[0].name
      }
    },
  },
  methods: {
    ...mapMutations("sftp", ["closeHovers","showHover"]),
    ...mapMutations(["setSftpReload"]),
  }
};
</script>
