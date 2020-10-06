<template>
    <div>
        <v-snackbar v-model="show" :color="color">
            {{ error.msg }}
        <v-btn color="yellow" text @click="close">
            Close
        </v-btn>
        </v-snackbar>
    </div>
</template>

<script>
// import store from '@/store'
import { mapState,mapMutations } from 'vuex'

export default {
    name:'snackbar',
    data() {
        return {
            show: false,
            color: 'pink',
        }
    },
    computed: {
        ...mapState(['error']),
    },
    watch: {
        error(val) {
            if(val.code > 0) {
                this.show = true
                if(val.code > 200) {
                    this.color = 'red'
                } else {
                    this.color = 'green'
                }
            } 
        }
    },
    methods: {
        ...mapMutations(['setError']),
        close() {
            // state.commit('setError',{msg:'',code:0,show:false})
            this.setError({msg:'',code:0})
            this.show = false
        }
    }
}
</script>