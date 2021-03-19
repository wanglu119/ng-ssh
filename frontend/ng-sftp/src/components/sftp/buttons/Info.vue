<template>
  <v-btn
    color="blue"
    fab
    small
    icon
    @click="showInfo"
  >
    <v-icon> mdi-information </v-icon>
  </v-btn>
</template>

<script>
import {mapGetters, mapMutations } from "vuex"

export default {
  name: 'info-button',
  computed: {
    ...mapGetters("sftp", ["selectedCount"])
  },
  methods: {
      ...mapMutations("sftp", ["showHover"]),
    showInfo() {
      this.showHover("info");
    }
  }
}
</script>