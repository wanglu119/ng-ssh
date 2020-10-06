<template>
  <v-app>
    <router-view/>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  name: 'App',
  data: () => ({
  }),
  mounted () {
    const loading: HTMLElement | null = document.getElementById('loading')
    if(loading !== null) {
      loading.classList.add('done')

      setTimeout(function () {
        if(loading.parentNode!== null) {
          loading.parentNode.removeChild(loading)
        }
      }, 200)
    }
  }
});
</script>
