<template>
    <v-container>
        <v-card>
            <v-card-title>Aliases </v-card-title>
            <v-card-text>
                <v-text-field v-model="search" append-icon="mdi-magnify" label="Search" single-line
                    hide-details></v-text-field>

                <v-data-table :headers="headers" :items="aliases" :search="search" v-model="selected" show-select return-object>
                    <template v-slot:item.enabled="{ item }">
                        <v-chip color="green" v-if="item.enabled">Yes</v-chip>
                        <v-chip color="red" v-if="!item.enabled">No</v-chip>
                    </template>
                </v-data-table>
            </v-card-text>
            <v-card-actions>
                <v-btn to="/alias/new" prepend-icon="mdi-plus-circle-outline">Add</v-btn>
                <v-btn @click="removeAlias()" v-if="selected.length > 0"
                    prepend-icon="mdi-delete">Delete</v-btn>
                <v-btn @click="editAlias()" v-if="selected.length == 1"
                    prepend-icon="mdi-pencil">Edit</v-btn>
            </v-card-actions>
        </v-card>
    </v-container>
</template>
<script>

import Client from "../service/Client";

export default {
    name: 'AliasView',
    methods: {
        getAliases: function () {
            Client.getAliases().then((res) => {
                this.aliases = res.data.items;
            });
        },
        editAlias: function () {
            this.$router.push("/alias/" + this.selected[0].id)
        },
        removeAlias: function () {
            this.$dialog.create({
                title: "Delete Alias",
                text: "Do you want to delete the Alias(es) " + this.selected.map(entry => entry.source_username + "@" + entry.source_domain.name).join(', ') + "?",
                cancelText: "No",
                confirmationText: "Yes",
                buttons: [
                    { key: 'cancel', title: 'Cancel', value: 'cancel', color: 'grey', variant: 'text' },
                    { key: 'ok', title: 'OK', value: 'ok', color: 'info', variant: 'tonal' }
                ],
            }).then((anwser) => {
                if (anwser === "ok") {
                    for (var i = 0; i < this.selected.length; i++) {
                        Client.removeAlias(this.selected[i].id).then(() => {
                            this.getAliases();
                        })
                    }
                }
            })
        }

    },
    mounted: function () {
        this.getAliases();

    },
    components: {

    },
    data: () => ({
        'headers': [
            {
                title: '#',
                sortable: true,
                value: 'id'
            },
            {
                title: 'Source',
                sortable: true,
                value: 'source_display',
            },
            {
                title: 'Destination',
                sortable: true,
                value: 'destination_display',
            },
            {
                title: 'Enabled',
                sortable: true,
                value: 'enabled'
            }
        ],
        'search': '',
        'aliases': [],
        'selected': [],

    }),
}
</script>
