<template>
  <div>
    <component ref="currentComponent" :is="currentComponent"></component>
    <snack-bar />
  </div>
</template>

<script>
import { mapState } from "vuex";
import SnackBar from "./SnackBar";
import Download from "./Download";
import Upload from './Upload';
import Replace from './Replace'
import Delete from './Delete'
import Move from "./Move"
import ReplaceRename from "./ReplaceRename"
import Copy from "./Copy"
import Info from "./Info"
import Rename from "./Rename"

export default {
  name: "prompts",
  components: {
    SnackBar,
    Download,
    Upload,
    Replace,
    Delete,
    Move,
    ReplaceRename,
    Copy,
    Info,
    Rename,
  },
  computed: {
    ...mapState("sftp", ["show"]),
    currentComponent: function() {
      console.log('------------currentComponent', this.show)
      const matched =
        [
          "info",
          "help",
          "delete",
          "rename",
          "move",
          "copy",
          "newFile",
          "newDir",
          "download",
          "replace",
          "replace-rename",
          "share",
          "upload"
        ].indexOf(this.show) >= 0;

      return (matched && this.show) || null;
    }
  }
};
</script>
