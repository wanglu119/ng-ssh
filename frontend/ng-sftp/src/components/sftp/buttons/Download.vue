<template>
  <v-btn color="blue" fab small icon @click="download">
    <v-badge color="green" :content="selectedCount">
      <v-icon>mdi-cloud-download</v-icon>
    </v-badge>
  </v-btn>
</template>

<script>
import { mapGetters, mapState, mapMutations } from "vuex";
import { files as api } from "@/api";

export default {
  name: "download-button",
  computed: {
    ...mapState("sftp", ["req", "selected"]),
    ...mapGetters("sftp", ["selectedCount"]),
    ...mapState(["currTabConf"])
  },
  methods: {
    ...mapMutations("sftp", ["showHover", "closeHovers"]),
    download() {
      if (this.selectedCount === 1 && !this.req[this.currTabConf.id].file_list[this.selected[0]].is_dir) {
        api.download(null, false, this.req[this.currTabConf.id].file_list[this.selected[0]].path);
        return;
      }
      if (this.selectedCount >= 1) {
        this.showHover("download");
      }
    }
  }
};
</script>
