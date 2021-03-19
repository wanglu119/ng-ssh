<template>
  <v-card
    :elevation="5"
    height="75"
    width="260"
    link
    @click="click"
    @dblclick="open"
    :color="isSelected ? 'blue' : ''"
  >
    <v-list-item three-line>
      <v-list-item-avatar tile>
        <v-icon left x-large>
          {{ fileMeta.is_dir ? "mdi-folder" : "mdi-file" }}
        </v-icon>
      </v-list-item-avatar>
      <v-list-item-content>
        <v-row align="center" class="mx-0">
          <p style="user-select:none">
            <strong>{{ fileMeta.name }}</strong>
            <br />
            {{ fileMeta.is_dir ? "&mdash;" : humanSize() }}
            <br />
            <time>{{ humanTime() }}</time>
          </p>
        </v-row>
      </v-list-item-content>
    </v-list-item>
  </v-card>
</template>

<script>
import filesize from "filesize";
import moment from "moment";
import { mapMutations, mapGetters, mapState } from "vuex";

export default {
  props: ["sshConfigName", "index", "fileMeta"],
  name: "file-item",
  data() {
    return {
      color: ""
    };
  },
  computed: {
    ...mapState("sftp", ["selected", "multiple"]),
    ...mapGetters("sftp", ["selectedCount"]),
    isSelected() {
      return this.selected.indexOf(this.index) !== -1;
    }
  },
  methods: {
    ...mapMutations("sftp", ["removeSelected", "resetSelected", "addSelected"]),
    humanSize: function() {
      return filesize(this.fileMeta.size);
    },
    humanTime: function() {
      return moment(this.fileMeta.time).fromNow();
    },
    click(event) {
      if (this.selectedCount !== 0) {
        event.preventDefault();
      }

      if (this.isSelected) {
        this.removeSelected(this.index);
        return;
      }

      if (event.shiftKey && this.selectedCount > 0) {
        let fi = 0;
        let la = 0;

        if (this.index > this.selected[0]) {
          fi = this.selected[0] + 1;
          la = this.index;
        } else {
          fi = this.index;
          la = this.selected[0] - 1;
        }

        for (; fi <= la; fi++) {
          if (this.selected.indexOf(fi) == -1) {
            this.addSelected(fi);
          }
        }
        return;
      }

      if (!event.ctrlKey && !this.multiple) {
        this.resetSelected();
      }
      this.addSelected(this.index);
    },
    open() {
      if (this.fileMeta.is_dir) {
        this.$emit("dblclick", this.fileMeta.path);
      }
    }
  }
};
</script>
