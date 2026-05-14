<script setup>
import { ref, onMounted } from 'vue'
import Client from "./service/Client.js";
import Topbar from "./components/Topbar.vue";
import NotificationBar from "./components/NotificationBar.vue";

const version = ref("unknown")
const notifyRef = ref(null);

onMounted(() => {
  window.$notify = notifyRef.value;
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
      <NotificationBar ref="notifyRef" />
    </v-main>
    <v-footer>
      <v-row justify="center" no-gutters>
        Mail Server Admin GUI (Version: {{version}})
      </v-row>
    </v-footer>
  </v-app>
</template>
