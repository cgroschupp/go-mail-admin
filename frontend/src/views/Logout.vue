<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth';

onMounted(() => {
    const store = useAuthStore()
    store.logout()
})
</script>

<template>
    <div class="container">
        Logout...
    </div>
</template>