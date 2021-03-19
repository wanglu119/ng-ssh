<template>
  <v-btn color="blue" fab small icon @click="toggleMultipleSelection()">
    <v-icon>mdi-check-circle</v-icon>
  </v-btn>
</template>

<script>
import { mapMutations, mapState } from "vuex";

export default {
  name: "select-multiple-button",
  computed: {
    ...mapState("sftp", ["multiple", "show", "showMessage"])
  },
  methods: {
    ...mapMutations("sftp", ["showHover", "closeHovers", "setMultiple"]),
    toggleMultipleSelection() {
      if (this.multiple) {
        this.closeHovers();
      } else {
        this.setMultiple(true);
        this.showHover({
          show: "snack-bar",
          showMessage: "多选模式已开启"
        });
      }
    }
  }
};
</script>
