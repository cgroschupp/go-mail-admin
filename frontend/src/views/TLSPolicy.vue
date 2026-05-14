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
    <ConfirmDialog ref="confirmDialog" />
</template>

<script>

import Client from "../service/Client";
import ConfirmDialog from "../components/ConfirmDialog.vue";

export default {
    name: 'TLSPolicyView',
    methods: {
        getPolicys: function () {
            Client.listTLSPolicies().then((res) => {
                this.tlspolicys = res.data.items;
            });
        },
        async removePolicy() {
            const policy = { ...this.selected[0] };

            const confirmed = await this.$refs.confirmDialog.confirm({
                title: "Delete TLS Policy",
                text: `Do you want to delete the TLS Policy for domain ${policy.domain_id}?`,
                cancelText: "No",
                confirmText: "Yes",
            });

            if (!confirmed) {
                return;
            }

            await Client.deleteTLSPolicy(policy.id);
            this.selected = [];
            await this.getPolicys();
        }
    },
    mounted: function () {
        this.getPolicys();
    },
    components: {
        ConfirmDialog,
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
