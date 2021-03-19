<template>
  <div class="text-center">
    <v-dialog :value="dialog" width="400" @click:outside="closeHovers">
      <v-card>
        <v-card-title>Upload</v-card-title>
        <v-card-text>Select an option to upload.</v-card-text>
        <v-divider></v-divider>
        <v-card-text align="center" class="mt-3">
          <v-btn color="success" x-large dark outlined max-width="100" max-height="100" @click="uploadFile">
            <span class="text-none">File</span>
            <v-icon>mdi-file</v-icon>
          </v-btn>

          <v-btn
            color="success"
            x-large
            dark
            outlined
            max-width="100"
            max-height="100"
            class="ml-10"
            @click="uploadFolder"
          >
            <span class="text-none">Folder</span>
            <v-icon>mdi-folder</v-icon>
          </v-btn>
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { mapMutations, mapState } from "vuex";

export default {
  name: "upload",
  computed: {
    ...mapState("sftp", ["show"]),
    dialog() {
      if (this.show === "upload") {
        return true;
      } else {
        return false;
      }
    }
  },
  methods: {
    ...mapMutations("sftp", ["closeHovers"]),
    uploadFile: function () {
      document.getElementById('upload-input').value = ''
      document.getElementById('upload-input').click()
    },
    uploadFolder: function () {
      document.getElementById('upload-folder-input').value = ''
      document.getElementById('upload-folder-input').click()
    }
  }
};
</script>
