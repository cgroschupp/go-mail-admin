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
                <v-btn @click="deleteAccounts()" v-if="selected.length > 0"
                    prepend-icon="mdi-close-circle-outline">Delete</v-btn>
                <v-btn @click="editAlias()" v-if="selected.length == 1" prepend-icon="mdi-circle-edit-outline">Edit</v-btn>
            </v-card-actions>
        </v-card>
    </v-container>
    <ConfirmDialog ref="confirmDialog" />
</template>

<script>

import Client from "../service/Client";
import ConfirmDialog from "@/components/ConfirmDialog.vue";

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
        async deleteAccounts() {
            const selected = [...this.selected];
            const accounts = selected
                .map(entry => `${entry.username}@${entry.domain.name}`)
                .join(", ");

            const confirmed = await this.$refs.confirmDialog.confirm({
                title: "Delete Account",
                text: `Do you want to delete the Account(s) ${accounts}?`,
                cancelText: "No",
                confirmText: "Yes",
            });

            if (!confirmed) {
                return;
            }

            await Promise.all(
                selected.map(account =>
                Client.deleteAccount(account.id)
                )
            );
            this.selected = [];
            await this.getAccounts();
        },
    },
    mounted: function () {
        this.getAccounts();

    },
    components: {
        ConfirmDialog,
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
