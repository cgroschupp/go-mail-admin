<template>
    <div class="container">
      Logout...
    </div>
</template>

<script>
    // @ is an alias to /src
    import Client from "../service/Client";
    export default {
        name: 'LogoutView',
        components: {
        },
        mounted() {
            Client.logout();
        }

    }
</script>
