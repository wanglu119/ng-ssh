<template>
  <div>
    <v-btn icon color="red" @click="toDelete()">
        <v-icon>mdi-delete</v-icon>
    </v-btn>
    
    <v-dialog v-model="delDialog" max-width="290">
        <v-card>
            <v-card-title class="headline">Sure delete?</v-card-title>
            <v-card-text>
            Are your sure delete  ?
            </v-card-text>

            <v-card-actions>
            <v-spacer></v-spacer>

            <v-btn color="green darken-1" text  @click="procDelete()">
                Yes
            </v-btn>

            <v-btn color="green darken-1" text @click="delDialog = false">
                No
            </v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
  </div>
</template>

<script>
export default {
    name: 'deleteBtn',
    props: ['procData','procMethod'],
    data() {
        return {
            delDialog: false,
        }
    },
    methods: {
        toDelete() {
            this.delDialog = true
        },
        procDelete: async function() {
            await this.procMethod(this.procData)
            this.delDialog = false
        }
    }
}
</script>