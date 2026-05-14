<template>
  <v-snackbar
    v-model="show"
    :color="color"
    :timeout="timeout"
    location="top right"
  >
    {{ text }}

    <template #actions>
      <v-btn
        variant="text"
        @click="show = false"
      >
        Close
      </v-btn>
    </template>
  </v-snackbar>
</template>

<script>
export default {
  name: "NotificationBar",

  data() {
    return {
      show: false,
      text: "",
      color: "success",
      timeout: 3000,
    };
  },

  methods: {
    notify(options) {
      this.text = options.text || "";
      this.color = options.color || "success";
      this.timeout = options.timeout || 3000;

      this.show = true;
    },
  },
};
</script>