<template>
    <v-container>
        <v-card>
            <v-card-title>Accounts</v-card-title>
            <v-card-text>
                <v-text-field v-model="search" append-icon="mdi-magnify" label="Search" single-line
                    hide-details></v-text-field>
                <v-data-table :headers="headers" :items="accounts" :search="search" v-model="selected" show-select
                    return-object>
                    <template #item.print="{ item }">{{ item.username }}@{{ item.domain.name }}</template>
                    <template v-slot:item.enabled="{ item }">
                        <v-chip color="green" v-if="item.enabled">Yes</v-chip>
                        <v-chip color="red" v-if="!item.enabled">No</v-chip>
                    </template>
                    <template v-slot:item.sendonly="{ item }">
                        <v-chip color="green" v-if="item.sendonly">Yes</v-chip>
                        <v-chip color="red" v-if="!item.sendonly">No</v-chip>
                    </template>
                </v-data-table>
            </v-card-text>
            <v-card-actions>
                <v-btn to="/account/new" prepend-icon="mdi-plus-circle-outline">Add</v-btn>
                <v-btn @click="deleteAccount()" v-if="selected.length > 0"
                    prepend-icon="mdi-close-circle-outline">Delete</v-btn>
                <v-btn @click="editAlias()" v-if="selected.length == 1" prepend-icon="mdi-circle-edit-outline">Edit</v-btn>
            </v-card-actions>
        </v-card>
    </v-container>
</template>

<script>

import Client from "../service/Client";

export default {
    name: 'AccountsView',
    methods: {
        getAccounts: function () {
            Client.getAccounts().then((res) => {
                this.accounts = res.data.items;
            });
        },
        editAlias: function () {
            this.$router.push("/account/" + this.selected[0].id)
        },
        deleteAccount: function () {
            this.$dialog.create({
                title: "Delete Account",
                text: "Do you want to delete the Account(s) " + this.selected.map(entry => entry.username + "@" + entry.domain.name).join(', ') + "?",
                cancelText: "No",
                confirmationText: "Yes",
                buttons: [
                    { key: 'cancel', title: 'Cancel', value: 'cancel', color: 'grey', variant: 'text' },
                    { key: 'ok', title: 'OK', value: 'ok', color: 'info', variant: 'tonal' }
                ],
            }).then((anwser) => {
                if (anwser === "ok") {
                    for (var i = 0; i <= this.selected.length; i++) {
                        Client.deleteAccount(this.selected[i].id).then(() => {
                            this.getAccounts();
                        })
                    }
                }
            })
        }

    },
    mounted: function () {
        this.getAccounts();

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
                title: 'E-Mail',
                sortable: true,
                key: "email",
                value: item => `${item.username}@${item.domain.name}`,
            },
            {
                title: 'Quota (MB)',
                sortable: true,
                value: 'quota'
            },
            {
                title: 'Send Only',
                sortable: true,
                value: 'sendonly'
            },
            {
                title: 'Enabled',
                sortable: true,
                value: 'enabled'
            }
        ],
        'search': '',
        'accounts': [],
        'selected': [],

    }),
}
</script>
