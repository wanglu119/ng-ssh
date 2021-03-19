<template>
  <div class="text-center">
    <v-dialog :value="dialog" width="300" @click:outside="closeHovers">
      <v-card>
        <v-card-title>Download files</v-card-title>
        <v-card-text>Choose the format you want to download.</v-card-text>
        <v-divider></v-divider>
        <v-card-text>
          <v-col>
            <v-row class="pt-3">
              <v-btn block color="primary" @click="download('zip')">zip</v-btn>
            </v-row>
            <v-row class="pt-3">
              <v-btn block color="primary" @click="download('tar')">tar</v-btn>
            </v-row>
            <v-row class="pt-3">
              <v-btn block color="primary" @click="download('targz')">tar.gz</v-btn>
            </v-row>
            <v-row class="pt-3">
              <v-btn block color="primary" @click="download('tarbz2')">tar.bz2</v-btn>
            </v-row>
            <v-row class="pt-3">
              <v-btn block color="primary" @click="download('tarxz')">tar.xz</v-btn>
            </v-row>
            <v-row class="pt-3">
              <v-btn block color="primary" @click="download('tarlz4')">tar.lz4</v-btn>
            </v-row>
            <v-row class="pt-3">
              <v-btn block color="primary" @click="download('tarsz')">tar.sz</v-btn>
            </v-row>
          </v-col>
        </v-card-text>
        <v-divider></v-divider>
      </v-card>
    </v-dialog>
    <div>{{show}}</div>
  </div>
</template>

<script>
import { mapState, mapMutations } from "vuex";

import { files as api } from "@/api";

export default {
  name: "download",
  computed: {
    ...mapState("sftp", ["show"]),
    ...mapState("sftp", ["req", "selected"]),
    ...mapState(["currTabConf"]),
    dialog() {
      if (this.show === "download") {
        return true;
      } else {
        return false;
      }
    }
  },
  methods: {
    ...mapMutations("sftp", ["closeHovers"]),
    download: function(format) {
      if (this.selectedCount === 0) {
        return;
      } else {
        const files = [];

        for (const i of this.selected) {
          files.push(this.req[this.currTabConf.id].file_list[i].path);
        }

        api.download(format, true, ...files);
      }

      this.closeHovers();
    }
  }
};
</script>