<template>
    <div class="container">
      Logout...
    </div>
</template>

<script>
    // @ is an alias to /src
    export default {
        name: 'LogoutView',
        components: {
        },
        mounted() {
            localStorage.removeItem('token')
            this.$router.push({ name: 'Login' })
        }
    }
</script>
