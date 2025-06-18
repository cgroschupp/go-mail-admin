<template>
    <v-container>
        <v-card>
            <v-card-title>TLS Policy</v-card-title>
            <v-card-text>
                <v-text-field v-model="search" append-icon="mdi-magnify" label="Search" single-line
                    hide-details></v-text-field>
                <v-data-table :headers="headers" :items="tlspolicys" :search="search" :single-select=true v-model="selected"
                    show-select return-object></v-data-table>
            </v-card-text>
            <v-card-actions>
                <v-btn :to="{ name: 'TLSNew' }" icon><v-icon>mdi-plus-circle-outline</v-icon></v-btn>
                <v-btn @click="removePolicy()" v-if="selected[0]" icon><v-icon>mdi-close-circle-outline</v-icon></v-btn>
                <v-btn :to="{ name: 'TLSEdit', params: { id: this.selected[0].id } }" v-if="selected[0]"
                    icon><v-icon>mdi-circle-edit-outline</v-icon></v-btn>
            </v-card-actions>
        </v-card>
    </v-container>
</template>

<script>

import Client from "../service/Client";

export default {
    name: 'TLSPolicyView',
    methods: {
        getPolicys: function () {
            Client.listTLSPolicies().then((res) => {
                this.tlspolicys = res.data.items;
            });
        },
        removePolicy: function () {
            this.$dialog.create({
                title: "Delete TLS Policy",
                text: "Do you want to delete the TLS Policy for domain " + this.selected[0].domain_id + "?",
                cancelText: "No",
                confirmationText: "Yes",
                buttons: [
                    { key: 'cancel', title: 'Cancel', value: 'cancel', color: 'grey', variant: 'text' },
                    { key: 'ok', title: 'OK', value: 'ok', color: 'info', variant: 'tonal' }
                ],
            }).then((anwser) => {
                if (anwser === "ok") {
                    Client.deleteTLSPolicy(this.selected[0].id).then(() => {
                        this.getPolicys();
                    })
                }
            })
        }
    },
    mounted: function () {
        this.getPolicys();
    },
    components: {

    },
    data: () => ({
        'headers': [
            {
                title: '#',
                sortable: true,
                value: 'id',
            },
            {
                title: 'Domain',
                sortable: true,
                value: 'domain.name'
            },
            {
                title: 'Policy',
                sortable: true,
                value: 'policy'
            },
            {
                title: 'Params',
                sortable: true,
                value: 'params'
            }
        ],
        'search': '',
        'tlspolicys': [],
        'selected': [],

    }),
}
</script>
