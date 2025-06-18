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
                <v-btn @click="removeDomain()" v-if="selected[0]" prepend-icon="mdi-delete">Delete Domain(s)</v-btn><br><br>
            </v-card-actions>
        </v-card>
    </v-container>
</template>

<script>

import Client from "../service/Client";

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
        removeDomain: function () {
            this.$dialog.create({
                title: "Delete Domain",
                text: "Do you want to delete the Domain " + this.selected[0].name + "?",
                cancelText: "No",
                confirmationText: "Yes",
                buttons: [
                    { key: 'cancel', title: 'Cancel', value: 'cancel', color: 'grey', variant: 'text' },
                    { key: 'ok', title: 'OK', value: 'ok', color: 'info', variant: 'tonal' }
                ],
            }).then((anwser) => {
                if (anwser === "ok") {
                    Client.removeDomain(this.selected[0].id).then(() => {
                        this.getDomains();
                    })
                }
            })
        },
        addDomain: function () {
            Client.addDomain(this.newDomain).then(() => {
                this.getDomains();
                this.newDomain = "";

            }).catch((e) => {
                var msg = e.response.data.message
                this.$notify.error("Something go wrong: " +msg, { location: 'bottom center' });
            });
        }
    },
    mounted: function () {
        this.getDomains();
    },
    components: {

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
