<template>
  <v-btn
    color="blue"
    fab
    small
    icon
    @click="showDelete"
  >
    <v-icon> mdi-delete </v-icon>
  </v-btn>
</template>

<script>
import {mapMutations, mapGetters} from "vuex"

export default {
  name: 'delete-button',
  computed: {
    ...mapGetters("sftp", ["selectedCount"])
  },
  methods: {
      ...mapMutations("sftp", ["showHover"]),
    showDelete() {
      if(this.selectedCount > 0) {
        this.showHover("delete");
      }
    }
  }
}
</script>