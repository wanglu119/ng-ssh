<template>
  <div class="text-center">
    <v-dialog :value="dialog" width="400" @click:outside="closeHovers">
      <v-card>
        <v-card-title class="headline">
          Delete
        </v-card-title>
        <v-card-text>
          <p v-if="selectedCount == 1">Are you sure you want to delete this file/folder?</p>
          <p v-else> Are you sure you want to delete {{selectedCount}} file(s)?</p>
        </v-card-text>
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
            @click="submit"
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
import { files as api } from '@/api'

export default {
  name: "delete",
  computed: {
    ...mapState("sftp", ["show","req", "selected"]),
    ...mapGetters("sftp", ["selectedCount"]),
    ...mapState(["currTabConf"]),
    dialog() {
      if (this.show === "delete") {
        return true;
      } else {
        return false;
      }
    }
  },
  methods: {
    ...mapMutations("sftp", ["closeHovers"]),
    ...mapMutations(["setSftpReload"]),
    submit: async function () {
      this.closeHovers()
      try {
        if (this.selectedCount === 0) {
          return
        }

        const promises = []
        for (const index of this.selected) {
          if(this.req[this.currTabConf.id].file_list[index].is_dir) {
            promises.push(api.removeDir(this.req[this.currTabConf.id].file_list[index].path))
          } else {
            promises.push(api.removeFile(this.req[this.currTabConf.id].file_list[index].path))
          }
        }
        await Promise.all(promises)
      
        this.setSftpReload(true)
      } catch (e) {
        console.log(e)
        this.setSftpReload(true)
      }
    }
  }
};
</script>
