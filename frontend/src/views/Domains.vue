<template>
    <v-container>
        <v-card>
            <v-card-title>Add Domain</v-card-title>
            <v-card-text>
                <v-text-field v-model="newDomain" placeholder="example.com"></v-text-field>
                <v-card-actions>
                    <v-btn @click="addDomain()" prepend-icon="mdi-plus-circle-outline">Add Domain</v-btn><br><br>
                </v-card-actions>
            </v-card-text>
        </v-card>
    </v-container>
    <v-container>
        <v-card>
            <v-card-title>Domains</v-card-title>
            <v-card-text>
                <v-text-field v-model="search" append-icon="mdi-magnify" label="Search" single-line
                    hide-details></v-text-field>

                <v-card-actions>
                    <v-btn @click="triggerRefresh()" prepend-icon="mdi-refresh">Refresh</v-btn>
                </v-card-actions>

                <v-data-table :headers="headers" :items="domains" :search="search" :single-select=true :single-expand=true
                    v-model="selected" show-select show-expand :expanded.sync="expanded" return-object>
                    <template v-slot:expanded-row="{ columns, item }">
                        <tr>
                            <td :colspan="columns.length">
                                <h2>General</h2>
                                <ul>
                                    <li><b>Name:</b> {{ item.name }}</li>
                                    <li><b>Created At:</b> {{ item.created_at }}</li>
                                    <li><b>Updated At:</b> {{ item.updated_at }}</li>
                                </ul>
                            </td>
                        </tr>
                    </template>
                </v-data-table>
            </v-card-text>
            <v-card-actions>
                <v-btn @click="deleteDomain()" v-if="selected[0]" prepend-icon="mdi-delete">Delete Domain(s)</v-btn><br><br>
            </v-card-actions>
        </v-card>
    </v-container>
    <ConfirmDialog ref="confirmDialog" />
</template>
<script>

import Client from "../service/Client";
import ConfirmDialog from "../components/ConfirmDialog.vue";

export default {
    name: 'DomainView',
    methods: {
        triggerRefresh: function () {
            Client.triggerDomains()
        },
        getDomains: function () {
            Client.getDomains().then((res) => {
                this.domains = res.data.items;
            });
        },
        async deleteDomain() {
            const domain = { ...this.selected[0] };

            const confirmed = await this.$refs.confirmDialog.confirm({
                title: "Delete Domain",
                text: `Do you want to delete the Domain ${domain.name}?`,
                cancelText: "No",
                confirmText: "Yes",
            });

            if (!confirmed) {
                return;
            }

            await Client.removeDomain(domain.id);

            this.selected = []
            await this.getDomains();
        },
        addDomain: function () {
            Client.addDomain(this.newDomain).then(() => {
                this.getDomains();
                this.newDomain = "";

            }).catch((e) => {
                var msg = e.response.data.message
                window.$notify.notify({
                    text: "Oups, something go wrong",
                    color: "error",
                });
            });
        }
    },
    mounted: function () {
        this.getDomains();
    },
    components: {
        ConfirmDialog,
    },
    data: () => ({
        expanded: [],
        'headers': [
            {
                title: '#',
                sortable: true,
                value: 'id'
            },
            {
                title: 'Domain',
                sortable: true,
                value: 'name'
            }
        ],
        'search': '',
        'domains': [],
        'selected': [],
        'newDomain': ''
    }),
}
</script>
