<template>

<v-card
  :elevation="5"
  height="75"
  width="260"
  link
  @click="click"
  @dblclick="open"
  :color="isSelected?'blue':''"
>
  <v-list-item three-line>
    <v-list-item-avatar tile>
      <v-icon left x-large >{{fileMeta.is_dir?"mdi-folder":"mdi-file"}}</v-icon>
    </v-list-item-avatar>
    <v-list-item-content>
      
      <v-row
        align="center"
        class="mx-0"
      >
        <p><strong>{{fileMeta.name}}</strong> <br> 
          {{fileMeta.is_dir?"&mdash;":humanSize()}} <br>
          <time>{{ humanTime() }}</time>
        </p>
      </v-row>
    </v-list-item-content>
  </v-list-item>
</v-card>
  
</template>

<script>
import filesize from 'filesize'
import moment from 'moment'
import { mapMutations, mapGetters, mapState,mapActions } from 'vuex'

export default {
  props: ['sshConfigName','index','fileMeta'],
  name: 'fileItem',
  data() {
    return {
      color: '',
    }
  },
  computed: {
    ...mapState(['sftpSelected',]),
    selectCount() {
      return this.sftpSelected.length
    },
    isSelected() {
      return (this.sftpSelected.indexOf(this.index) !== -1)
    },
  },
  methods: {
    ...mapMutations(['removeSftpSelected','resetSftpSelected','addSftpSelected']),
    humanSize: function () {
      return filesize(this.fileMeta.size)
    },
    humanTime: function () {
      return moment(this.fileMeta.time).fromNow()
    },
    click(event) {
      if (this.selectedCount !== 0){
        event.preventDefault()
      }

      if (this.isSelected) {
        this.removeSftpSelected(this.index)
        return
      }

      if (event.shiftKey && this.sftpSelected.length > 0) {
        let fi = 0
        let la = 0

        if (this.index > this.sftpSelected[0]) {
          fi = this.sftpSelected[0] + 1
          la = this.index
        } else {
          fi = this.index
          la = this.sftpSelected[0] - 1
        }

        for (; fi <= la; fi++) {
          if (this.sftpSelected.indexOf(fi) == -1) {
            this.addSftpSelected(fi)
          }
        }
        return
      }

      if (!event.ctrlKey) {
        this.resetSftpSelected()
      }
      this.addSftpSelected(this.index)
    },
    open() {
      if(this.fileMeta.is_dir) {
        this.$emit('dblclick',this.fileMeta.path)
      }
    },
  },
}
</script>

