<template>
  <div>
    <v-snackbar
      v-model="isShow"
      :timeout="-1"
      :vertical="true"
      shaped
      bottom
      center
      color="blue"
    >
      {{ showMessage }}
      <template v-slot:action="{ attrs }">
        <v-btn color="yellow" text v-bind="attrs" @click="closeHovers"
          >Close
        </v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script>
import { mapMutations, mapState } from "vuex";

export default {
  name: "snackbar",
  computed: {
    ...mapState("sftp", ["show", "showMessage"]),
    isShow() {
      if (this.show === "snack-bar") {
        return true;
      } else {
        return false;
      }
    }
  },
  methods: {
    ...mapMutations("sftp", ["closeHovers"])
  }
};
</script>
