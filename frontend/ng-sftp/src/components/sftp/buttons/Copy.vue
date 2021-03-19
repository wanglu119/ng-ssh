<template>
  <v-btn
    color="blue"
    fab
    small
    icon
    @click="showCopy"
  >
    <v-icon>mdi-content-copy </v-icon>
  </v-btn>
</template>

<script>
import { mapGetters,mapMutations } from "vuex"

export default {
  name: 'copy-button',
  computed: {
    ...mapGetters("sftp", ["selectedCount"])
  },
  methods: {
      ...mapMutations("sftp", ["showHover"]),
    showCopy() {
      if(this.selectedCount > 0) {
        this.showHover("copy");
      }
    }
  }
}
</script>