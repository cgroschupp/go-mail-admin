<script setup>
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth.js'
import { storeToRefs } from 'pinia'
const username = ref("")
const password = ref("")

const store = useAuthStore()
const {loginFailed} = storeToRefs(store)

function login() {
    store.login(username.value, password.value)
}
</script>
<template>
    <v-container>
        <v-sheet width="300" class="mx-auto">
            <v-alert color="error" v-if="loginFailed" text="Login failed"></v-alert>
            <v-form fast-fail @submit.prevent>
                <v-text-field v-model="username" placeholder="Username" label="Username"></v-text-field>
                <v-text-field v-model="password" placeholder="Password" label="Password" type="password"></v-text-field>
                <v-btn type="submit" block color="success" class="mt-2" @click="login">Submit</v-btn>
            </v-form>
        </v-sheet>
    </v-container>
</template>
