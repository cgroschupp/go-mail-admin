<template>
    <div>
        <v-container>
            <v-card>
                <v-card-text>
                    <v-text-field v-model="alias.source_username" placeholder="Source Username" label="Source Username"
                        v-on:keyup="checkatsymbole"></v-text-field>
                    <v-select v-model="alias.source_domain_id" :item-title="item=>item.name" :item-value="item=>item.id" :items="domains" label="Source Domain"></v-select>
                    <v-text-field v-model="alias.destination_username" placeholder="Destination Username"
                        label="Destination Username"></v-text-field>
                    <v-text-field v-model="alias.destination_domain" placeholder="Destination Domain"
                        label="Destination Domain"></v-text-field>
                    <v-checkbox v-model="alias.enabled" label="Enabled"></v-checkbox>
                </v-card-text>
                <v-card-actions>
                    <v-btn @click="saveAlias" prepend-icon="mdi-content-save">Save</v-btn>
                </v-card-actions>
            </v-card>
        </v-container>
    </div>
</template>

<script>
import Client from "../service/Client";
export default {
    name: 'AliasEditView',
    methods: {
        checkatsymbole: function (r) {
            if (r.key == "@") {
                this.alias.source_username = this.alias.source_username.substr(0, this.alias.source_username.length - 1);
                this.$refs.domain.focus();
            }
        },
        getAlias: function () {
            Client.getAlias(this.$route.params.id).then((res) => {
                this.alias = res.data
            });
        },
        getDomains: function () {
            Client.getDomains().then((res) => {
                this.domains = res.data.items
            });
        },
        saveAlias: function () {
            if (this.alias.id) {
                Client.saveAlias(this.alias.id, this.alias).then(() => {
                    this.$notify.info("Alias saved", { location: 'bottom center' });
                    this.$router.push("/alias")
                })
            } else {
                Client.createAlias(this.alias).then(() => {
                    this.$notify.info("Alias created", { location: 'bottom center' });
                    this.$router.push("/alias")
                }, (e) => {
                    var msg = e.response.data
                    if (msg == "Source Username can`t be empty string, only null or string is valid") {
                        msg = "Enter Source Username or enable catch all!"
                    }
                    this.$notify.error("Something go wrong", { location: 'bottom center' });
                })
            }

        }
    },

    mounted: function () {
        this.getDomains();
        if (this.$route.params.id != "new") {
            this.getAlias();
        }
    },
    components: {

    },
    data: () => ({
        alias: { enabled: true },
        domains: [],
    }),
}
</script>
