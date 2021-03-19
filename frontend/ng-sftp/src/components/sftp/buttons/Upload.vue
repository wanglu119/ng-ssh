<template>
  <v-btn color="blue" fab small icon @click="showSelectUpload">
    <v-icon>mdi-cloud-upload</v-icon>
  </v-btn>
</template>

<script>
import { mapMutations } from "vuex";
export default {
  name: "upload-button",
  methods: {
    ...mapMutations("sftp", ["showHover"]),
    showSelectUpload() {
      this.showHover("upload");
    }
  }
};
</script>
