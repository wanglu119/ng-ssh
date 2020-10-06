<template>
<div>
  <v-card>
    <v-breadcrumbs :items="folders">
      <template v-slot:item="{ item }">
        <v-breadcrumbs-item
          href="#"
          @click="setDirPath(item.href)"
          :disabled="item.disabled"
        >
          {{ item.text }}
        </v-breadcrumbs-item>
      </template>
    </v-breadcrumbs>

    <v-divider></v-divider>

    <div v-if="dirCount > 0">
       <v-subheader>Folders:</v-subheader>
      <v-row class="fill-height" >
        <v-col
          v-for="fileMeta in dirList"
          :key="fileMeta.index"
          cols="auto"
        > 
          <fileItem v-if="fileMeta.is_dir" :sshConfigName="conf.sshConfigName" :index="fileMeta.index" :fileMeta="fileMeta"
            @dblclick="setDirPath"/>
        </v-col>
      </v-row>
    </div>

    <v-divider></v-divider>
    <div v-if="fileCount > 0">
       <v-subheader>Files:</v-subheader>
      <v-row  class="fill-height" >
        <v-col
          v-for="fileMeta in fileList"
          :key="fileMeta.index"
          cols="auto"
        >
          <fileItem v-if="!fileMeta.is_dir" :sshConfigName="conf.sshConfigName" :index="fileMeta.index" :fileMeta="fileMeta"/>
        </v-col>
      </v-row>
    </div>
  </v-card>
</div>
  
</template>

<script>
import { mapMutations, mapGetters, mapState,mapActions } from 'vuex'
import {procURL} from '@/api/utils'
import fileItem from '@/components/files/fileItem'

export default {
  components:{fileItem,},
  props: ['conf'],
  name:'sftp-file-list',
  data() {
    return {
      benched: 0,
      folders: [],
      fileList: [],
      dirList: [],
      fileCount: 0,
      dirCount:0,
      dirPath: '',
    }
  },
  mounted() {
    this.getSftpFileList()
  },
  methods: {
    ...mapMutations(['resetSftpSelected']),
    setDirPath(dirPath) {
      this.dirPath = dirPath
      this.getSftpFileList()
    },
    getSftpFileList: async function() {
      this.resetSftpSelected()
      
      const r = await procURL('/ngssh_api/sftp/listFiles',{
        method: 'POST',
        data: {
          name: this.conf.sshConfigName,
          'dir_path': this.dirPath,
        }
      })

      this.fileCount = r.ls_meta.file_count
      this.dirCount = r.ls_meta.dir_count

      // parser folders
      this.folders = []
      let href = ''
      const dirs = r.ls_meta.dir_path.split('/')
      for(const i in dirs) {
        const name = dirs[i]
        if(name === '') {
          this.folders.push(
            {
              text: 'root',
              disabled: false,
              href: '/',
            }
          )
        } else {
          href += '/'+name
          this.folders.push(
            {
              text: name,
              disabled: false,
              href: href,
            }
          )
        }
      }
      this.folders[this.folders.length-1].disabled=true
      this.fileList = []
      this.dirList = []
      for(const i in r.file_list) {
        const meta = r.file_list[i]
        if(meta.is_dir) {
          this.dirList.push(meta)
        } else {
          this.fileList.push(meta)
        }
      }
    },
  }
}
</script>