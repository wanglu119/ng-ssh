<template>
  <div class="text-center">
    <v-dialog :value="dialog" width="400" @click:outside="closeHovers">
      <v-card>
        <v-card-title class="headline">
          Replace
        </v-card-title>
        <v-card-text>One of the files you're trying to upload is conflicting because of its name. Do you wish to replace the existing one?</v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="green darken-1"
            text
            @click="closeHovers"
          >
            Cancel
          </v-btn>
          <v-btn
            color="green darken-1"
            text
            @click="showConfirm"
          >
            Ok
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>


<script>
import { mapMutations, mapState } from "vuex";

export default {
  name: "replace",
  computed: {
    ...mapState("sftp", ["show","showConfirm"]),
    dialog() {
      if (this.show === "replace") {
        return true;
      } else {
        return false;
      }
    }
  },
  methods: {
    ...mapMutations("sftp", ["closeHovers"]),

  }
};
</script>
