<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="elevation-12">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title>Login form</v-toolbar-title>
            <v-spacer />
          </v-toolbar>
          <v-card-text>
            <v-form ref="form" v-model="valid" lazy-validation>
              <v-text-field
                label="Username"
                prepend-icon="mdi-account"
                type="text"
                :rules="[rules.required]"
                v-model="username"
              />
              <v-text-field
                label="Password"
                prepend-icon="mdi-lock"
                type="password"
                :rules="[rules.required]"
                v-model="password"
                @keyup.enter="login"
              />
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer />
            <v-btn color="primary" :disabled="!valid" @click="login">Sign In</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapMutations } from 'vuex'
import { parseToken } from '@/utils/auth'
import { procURL } from '@/api/utils'

export default {
    name: 'login',
    data() {
        return {
            valid: true,
            username: '',
            password: '',
            rules: {
                required: value => !!value || 'Required',
            },
        }
    },
 
    methods: {
        ...mapMutations(['setError']),
    
        getServerUrl: async function(mainTunnelInfos) {
            console.log(mainTunnelInfos)
        },
        login: async function() {
            // const debug = getQueryString('debug')
            if (this.$refs.form.validate()) {
                try {
                    const r = await procURL('/ngssh_api/auth/login',{
                        method: 'POST',
                        data: {
                            username: this.username,
                            password: this.password,
                        }
                    })
                    parseToken(r)
                    this.$router.push({
                        name: 'Home'
                    })
                }catch(err) {
                    console.log(err)
                }
            }
        },
    },
}
</script>
