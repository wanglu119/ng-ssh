<template>
  <div>
    <v-card>
      <v-breadcrumbs :items="folders">
        <template v-slot:item="{ item }">
          <v-breadcrumbs-item
            href="#"
            @click="setDirPath(item.href)"
            :disabled="item.disabled"
          >{{ item.text }}</v-breadcrumbs-item>
        </template>
      </v-breadcrumbs>

      <v-divider></v-divider>

      <div v-if="dirCount > 0">
        <v-subheader>Folders:</v-subheader>
        <v-row class="fill-height">
          <v-col v-for="fileMeta in dirList" :key="fileMeta.index" cols="auto">
            <file-item
              v-if="fileMeta.is_dir"
              :sshConfigName="conf.sshConfigName"
              :index="fileMeta.index"
              :fileMeta="fileMeta"
              @dblclick="setDirPath"
            />
          </v-col>
        </v-row>
      </div>

      <v-divider></v-divider>
      <div v-if="fileCount > 0">
        <v-subheader>Files:</v-subheader>
        <v-row class="fill-height">
          <v-col v-for="fileMeta in fileList" :key="fileMeta.index" cols="auto">
            <file-item
              v-if="!fileMeta.is_dir"
              :sshConfigName="conf.sshConfigName"
              :index="fileMeta.index"
              :fileMeta="fileMeta"
            />
          </v-col>
        </v-row>
      </div>
    </v-card>
    <prompts />
  </div>
</template>

<script>
import { mapMutations,mapState } from "vuex";

import { procURL } from "@/api/utils";
import FileItem from "@/components/sftp/files/fileItem";
import Prompts from "@/components/sftp/prompts/Prompts";

export default {
  components: { FileItem, Prompts },
  props: ["conf"],
  name: "sftp-file-list",
  data() {
    return {
      benched: 0,
      folders: [],
      fileList: [],
      dirList: [],
      fileCount: 0,
      dirCount: 0,
    };
  },
  computed: {
    ...mapState("sftp", ["req","currPath"]),
    ...mapState(["currTabConf"])
  },
  mounted() {
    this.getSftpFileList();
  },
  methods: {
    ...mapMutations("sftp", ["resetSelected", "setReq","closeHovers","showHover"]),
    ...mapMutations(["setSftpReload"]),
    setDirPath(dirPath) {
      this.currTabConf.currDir = dirPath
      this.getSftpFileList();
    },
    getSftpFileList: async function() {
      this.resetSelected();

      const r = await procURL("/ngssh_api/sftp/listFiles", {
        method: "POST",
        data: {
          name: this.conf.sshConfigName,
          dir_path: this.currTabConf.currDir
        }
      });

      this.fileCount = r.ls_meta.file_count;
      this.dirCount = r.ls_meta.dir_count;

      // parser folders
      this.folders = [];
      let href = "";
      this.currTabConf.currDir = r.ls_meta.dir_path
      const dirs = r.ls_meta.dir_path.split("/");
      for (const i in dirs) {
        const name = dirs[i];
        if (name === "") {
          this.folders.push({
            text: "root",
            disabled: false,
            href: "/"
          });
        } else {
          href += "/" + name;
          this.folders.push({
            text: name,
            disabled: false,
            href: href
          });
        }
      }
      this.folders[this.folders.length - 1].disabled = true;
      this.fileList = [];
      this.dirList = [];
      this.setReq({tabId:this.conf.id, data:r});
      for (const i in r.file_list) {
        const meta = r.file_list[i];
        if (meta.is_dir) {
          this.dirList.push(meta);
        } else {
          this.fileList.push(meta);
        }
      }
    },
    
  }
};
</script>
