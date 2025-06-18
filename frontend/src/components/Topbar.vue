<script setup>
import { ref, onMounted } from 'vue'
import Client from "../service/Client";

const color = ref("gray")

onMounted(() => {
    Client.getStatus().then((res) => {
        if (res.data.healthy) {
            color.value = "green";
        } else {
            color.value = "red";
        }
    }).catch(() => {
        color.value = "grey";
    })
})
</script>

<template>
    <div>
        <v-app-bar theme="dark">
            <v-app-bar-title>Go Mail Admin <v-icon :style="{ color: color }">mdi-heart</v-icon></v-app-bar-title>
            <template v-slot:append>
                <v-btn aria-label="Dashboard" to="/" icon="mdi-view-dashboard-variant"></v-btn>
                <v-btn aria-label="Domains" to="/domains" icon="mdi-dns"></v-btn>
                <v-btn aria-label="Aliases" to="/alias" icon="mdi-forwardburger"></v-btn>
                <v-btn aria-label="Accounts" to="/account" icon="mdi-account"></v-btn>
                <v-btn aria-label="TLS Policies" to="/tls" icon="mdi-security"></v-btn>
                <v-btn aria-label="Logout" to="/logout" icon="mdi-logout"></v-btn>
            </template>
        </v-app-bar>
    </div>
</template>
