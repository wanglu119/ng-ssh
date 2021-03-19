<template>
  <div>
    <v-list dense>
      <v-list-item-group
        v-model="item"
        color="primary"
      >
        <v-list-item
          v-for="(item,i) in items"
          :key="i"
          @dblclick="next(item)"
        >
          <v-list-item-icon>
            <v-icon>mdi-folder</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title v-text="item.name"></v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list-item-group>
      <v-subheader><b>Currently navigating on:</b> {{current}}</v-subheader>
    </v-list>
  </div>
</template>

<script>
import { mapState } from "vuex"

import url from '@/utils/url'
import { files } from '@/api'

export default {
  name: "file-list",
  data() {
    return {
      items: [],
      item: null,
      current: "",
      selected: null,
    }
  },
  computed: {
    ...mapState("sftp", ["req"]),
    ...mapState(["currTabConf"])
  },
  watch: {
    item(val) {
      if(val !== undefined) {
        this.selected = this.items[val]
      } else {
        this.selected = null
      }
    }
  },
  mounted() {
    this.fillOptions(this.req[this.currTabConf.id])
  },
  methods: {
    fillOptions (req) {
      this.current = req.ls_meta.dir_path
      this.items = []

      // If the path isn't the root path,
      // show a button to navigate to the previous
      // directory.
      if (this.current !== '/') {
        this.items.push({
          name: '..',
          url: url.removeLastDir(this.current) + '/'
        })
      }

      // If this folder is empty, finish here.
      if (req === null) return

      // Otherwise we add every directory to the
      // move options.
      for (const item of req.file_list) {
        if (!item.is_dir) continue

        this.items.push({
          name: item.name,
          url: item.path
        })
      }
    },
    next: async function (item) {
      // Retrieves the URL of the directory the user
      // just clicked in and fill the options with its
      // content.
      const r = await files.fetch(item.url)
      this.fillOptions(r)
    },
  }
}
</script>