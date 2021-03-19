<template>
  <div class="text-center">
    <v-dialog :value="dialog" width="400" @click:outside="closeHovers">
      <v-card
        class="mx-auto"
        tile
      >
        <v-card-title class="headline">
          Move
        </v-card-title>

        <file-list ref="fileList"/>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="green darken-1"
            text
            @click="closeHovers"
          >
            Cancel
          </v-btn>
          <v-btn
            color="green darken-1"
            text
            @click="move"
          >
            Move
          </v-btn>
        </v-card-actions>
        
      </v-card>
    </v-dialog>
  </div>
</template>


<script>
import { mapMutations, mapState, mapGetters } from "vuex";

import { files as api } from '@/api'
import FileList from './FileList'
import * as upload  from '@/utils/upload'

export default {
  name: "move",
  components: {
    FileList,
  },
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
      if (this.show === "move") {
        return true;
      } else {
        return false;
      }
    }
  },
  methods: {
    ...mapMutations("sftp", ["closeHovers","showHover"]),
    ...mapMutations(["setSftpReload"]),
    move: async function () {
      try {
        // this.closeHovers()
        this.dest = this.$refs["fileList"].selected?this.$refs["fileList"].selected.url:this.$refs["fileList"].current
        if(!this.dest){
          console.log("error, dest empty:",this.dest)
          return
        } 
        const items = []

        for (const item of this.selected) {
          const fileInfo = this.req[this.currTabConf.id].file_list[item]
          items.push({
            from: fileInfo.path,
            to: this.dest +"/"+ encodeURIComponent(fileInfo.name),
            name: fileInfo.name
          })
        }

        const action = async (overwrite, rename) => {
          // buttons.loading('move')
          
          await api.move(items, overwrite, rename).then(() => {
            // buttons.success('move')
            this.currTabConf.currDir = this.dest
            this.setSftpReload(true)
            this.closeHovers()
          }).catch((e) => {
            // buttons.done('move')
            // this.$showError(e)
            console.log(e)
            this.closeHovers()
          })
        }
        
        const dstItems = (await api.fetch(this.dest)).file_list
        const conflict = upload.checkConflict(items, dstItems)

        let overwrite = false
        let rename = false

        if (conflict) {
          this.showHover({
            prompt: 'replace-rename',
            confirm: (event, option) => {
              overwrite = option == 'overwrite'
              rename = option == 'rename'

              action(overwrite, rename)
            }
          })
          return
        }

        action(overwrite, rename)
      
      } catch (e) {
        console.log(e)
        this.closeHovers()
        this.setSftpReload(true)
      }
    }
  }
};
</script>
