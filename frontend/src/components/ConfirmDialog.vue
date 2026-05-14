<template>
  <v-dialog v-model="open" max-width="500">
    <v-card>
      <v-card-title>{{ title }}</v-card-title>

      <v-card-text>
        {{ text }}
      </v-card-text>

      <v-card-actions>
        <v-spacer />

        <v-btn variant="text" color="grey" @click="resolve(false)">
          {{ cancelText }}
        </v-btn>

        <v-btn variant="tonal" color="info" @click="resolve(true)">
          {{ confirmText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref } from "vue";

const open = ref(false);

const title = ref("");
const text = ref("");
const cancelText = ref("No");
const confirmText = ref("Yes");

let resolver = null;

function confirm(options) {
  title.value = options.title;
  text.value = options.text;
  cancelText.value = options.cancelText ?? "No";
  confirmText.value = options.confirmText ?? "Yes";

  open.value = true;

  return new Promise(resolve => {
    resolver = resolve;
  });
}

function resolve(value) {
  open.value = false;

  if (resolver) {
    resolver(value);
    resolver = null;
  }
}

defineExpose({
  confirm,
});
</script>