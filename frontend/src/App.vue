<script setup>
import { ref, onMounted } from 'vue'
import Client from "./service/Client";
import Topbar from "./components/Topbar";

const version = ref("unknown")

onMounted(() => {
    Client.getVersion().then((res) => {
      if (res !== undefined) {
        version.value = res.data.version
      }
    }).catch((error) => {
        version.value = "unknown"
    });
})
</script>

<template>
  <v-app>
    <Topbar/>
    <v-main>
      <router-view />
    </v-main>
    <v-footer>
      <v-row justify="center" no-gutters>
        Mail Server Admin GUI (Version: {{version}})
      </v-row>
    </v-footer>
  </v-app>
</template>
