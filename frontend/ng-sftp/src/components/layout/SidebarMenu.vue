<template>
  <div>
    <!-- header -->
    <v-app-bar app clipped-left>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
      <v-toolbar-title>
        <span>ng-</span>
        <span class="font-weight-light">sftp</span>
      </v-toolbar-title>
    </v-app-bar>
    <!-- Sidebar -->
    <v-navigation-drawer
      v-model="drawer"
      app
      clipped
    >
      <v-list dense>
        <!-- Multi-level menu list -->
        <template v-for="item in sidebar">
          <v-row
            v-if="item.heading"
            :key="item.heading"
            align="center"
          >
            <v-col cols="6">
              <v-subheader v-if="item.heading">
                {{ item.heading }}
              </v-subheader>
            </v-col>
            <v-col
              cols="6"
              class="text-center"
            >
              <a
                href="#!"
                class="body-2 black--text"
              >EDIT</a>
            </v-col>
          </v-row>
          <v-list-group
            v-else-if="item.children"
            :key="item.text"
            v-model="item.model"
            :prepend-icon="item.model ? item.icon : item['icon-alt']"
            append-icon=""
          >
            <template v-slot:activator>
              <v-list-item-content>
                <v-list-item-title>
                  {{ item.text }}
                </v-list-item-title>
              </v-list-item-content>
            </template>
            <v-list-item
              v-for="(child, i) in item.children"
              :key="i"
              link
              @click="to(child)"
            >
              <v-list-item-action v-if="child.icon">
                <v-icon>{{ child.icon }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>
                  {{child.text}}
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list-group>
          <v-list-item
            v-else
            :key="item.text"
            link
            @click="to(item)"
          >
            <v-list-item-action>
              <v-icon>{{ item.icon }}</v-icon>
            </v-list-item-action>
            <v-list-item-content>
              <v-list-item-title>
                {{item.text }}
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </template>
      </v-list>
    </v-navigation-drawer>
  </div>
</template>

<script>
// import { mapGetters, mapState, mapMutations } from 'vuex'
import {procURL} from '@/api/utils'
import { mapMutations } from 'vuex'

export default {
  name:"SidebarMenu",
  data: () => ({
    drawer: null,
    sidebar: [
      {
        link: 'sftp',
        icon: 'mdi-chevron-up',
        'icon-alt': 'mdi-chevron-down',
        text: 'sftp',
        model: false,
        children: [
          {
            type: '',
            text: 'sidebar.admin.users',
            sshConfigName: 'admin_users',
            icon: 'mdi-cog-outline',
          },
        ]
      },
      {
        link: 'scp',
        icon: 'mdi-chevron-up',
        'icon-alt': 'mdi-chevron-down',
        text: 'scp',
        model: false,
        children: [
          {
            type: '',
            text: 'sidebar.admin.users',
            sshConfigName: 'admin_users',
            icon: 'mdi-view-dashboard',
          },
        ]
      },
    ],
  }),
  mounted: async function() {
    await this.getSshConfig()
  },
  computed: {
  },
  methods: {
    ...mapMutations(['setTab']),
    to(conf) {
      const newConf = {...conf}
      newConf['id'] = Date.now()
      newConf['name'] = conf.type+':'+conf.sshConfigName
      this.setTab(newConf)
    },
    getSshConfig: async function() {
      const r = await procURL('/ngssh_api/ssh_config/getSftpConfigNames',{
        method: 'GET',
      })
      this.sidebar[0].children = []
      this.sidebar[1].children = []
      for(const i in r) {
        const val = r[i]
        this.sidebar[0].children.push({
          text: val,
          sshConfigName: val,
          type: 'sftp',
          icon: 'mdi-cog-outline',
        })
        this.sidebar[1].children.push({
          text: val,
          sshConfigName: val,
          type: 'scp',
          icon: 'mdi-view-dashboard',
        })
      }
    }
  }
}
</script>