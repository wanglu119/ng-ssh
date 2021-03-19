<template>
  <div>
    <v-toolbar-items v-show="showSftpButton">
      <rename-button v-show="showRenameButton"/>
      <copy-button v-show="showCopyButton"/>
      <move-button v-show="showMoveButton"/>
      <delete-button v-show="showDeleteButton"/>
      <download-button v-show="showDownloadButton"/>
      <upload-button />
      <info-button />
      <select-multiple-button />
    </v-toolbar-items>
  </div>
</template>

<script>
import { mapGetters,mapState } from "vuex"

import InfoButton from "@/components/sftp/buttons/Info";
import RenameButton from "@/components/sftp/buttons/Rename";
import SelectMultipleButton from "@/components/sftp/buttons/SelectMultiple";
import UploadButton from "@/components/sftp/buttons/Upload";
import DownloadButton from "@/components/sftp/buttons/Download";
import DeleteButton from "@/components/sftp/buttons/Delete";
import MoveButton from "@/components/sftp/buttons/Move";
import CopyButton from "@/components/sftp/buttons/Copy";

export default {
  name: "site-header",
  components: {
    InfoButton,
    RenameButton,
    SelectMultipleButton,
    UploadButton,
    DownloadButton,
    DeleteButton,
    MoveButton,
    CopyButton
  },
  computed: {
    ...mapGetters("sftp", ["selectedCount"]),
    ...mapState(["currTabConf"]),
    showSftpButton() {
      if(this.currTabConf && this.currTabConf.type ==='sftp') {
        return true
      }
      return false
    },
    showRenameButton() {
      if(this.selectedCount === 1) {
        return true
      }
      return false
    },
    showCopyButton() {
      if(this.selectedCount > 0) {
        return true
      }
      return false
    },
    showMoveButton() {
      if(this.selectedCount > 0 ) {
        return true
      }
      return false
    },
    showDeleteButton() {
      if(this.selectedCount > 0) {
        return true
      }
      return false
    },
    showDownloadButton() {
      if(this.selectedCount > 0) {
        return true
      }
      return false
    },
  },
};
</script>
