<template>
  <v-btn
    color="blue"
    fab
    small
    icon
    @click="showRename"
  >
    <v-icon>mdi-pencil</v-icon>
  </v-btn>
</template>

<script>
import { mapGetters,mapMutations } from "vuex"

export default {
  name: 'rename-button',
  computed: {
    ...mapGetters("sftp", ["selectedCount"])
  },
  methods: {
      ...mapMutations("sftp", ["showHover"]),
    showRename() {
      if(this.selectedCount === 1) {
        this.showHover("rename");
      }
    }
  }
}
</script>