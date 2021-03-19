<template>
  <v-btn
    color="blue"
    fab
    small
    icon
    @click="showMove"
  >
    <v-icon>mdi-call-made </v-icon>
  </v-btn>
</template>

<script>
import { mapGetters,mapMutations } from "vuex"

export default {
  name: 'move-button',
  computed: {
    ...mapGetters("sftp", ["selectedCount"])
  },
  methods: {
      ...mapMutations("sftp", ["showHover"]),
    showMove() {
      if(this.selectedCount > 0) {
        this.showHover("move");
      }
    }
  }
}
</script>