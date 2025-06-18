<template>
  <v-container>
    <v-row no-gutters>
      <v-col cols="12" sm="6">
        
        <v-card>
          <v-card-title>Go Mail Admin</v-card-title>
          <v-card-subtitle>Fast Access</v-card-subtitle>
          <v-card-text>
            <v-list color="primary">
              <v-list-item to="/account/new">
                <v-list-item-title>New Account</v-list-item-title>
              </v-list-item>
              <v-list-item to="/alias/new">
                <v-list-item-title>New Aliases</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col v-for="(value, key) in stats" :key="key" cols="12" sm="2">
        <v-card style="margin-right: 10px;">
          <v-card-title>{{ key }}</v-card-title>
          <v-card-text>
            <Chart :chartData="value" />
          </v-card-text>
          </v-card>
        </v-col>
    </v-row>
    <v-card style="margin-top: 10px;">
      <v-card-title>Icons</v-card-title>
      <v-card-text>
        <v-table>
          <template v-slot:default>
            <thead>
            <tr>
              <th>Icon</th>
              <th>Description</th>
            </tr>
            </thead>
            <tbody>
            <tr>
              <td><v-icon>mdi-plus-circle-outline</v-icon></td>
              <td>Add a new entry to the current list (e.g. on the Account Page)</td>
            </tr>
            <tr>
              <td><v-icon>mdi-circle-edit-outline</v-icon></td>
              <td>Edit the current selected entry</td>
            </tr>
            <tr>
              <td><v-icon>mdi-close-circle-outline</v-icon></td>
              <td>Remove the current selected entry</td>
            </tr>
            <tr>
              <td><v-icon>mdi-dns</v-icon></td>
              <td>Domain List</td>
            </tr>
            <tr>
              <td><v-icon>mdi-forwardburger</v-icon></td>
              <td>Aliases</td>
            </tr>
            <tr>
              <td><v-icon>mdi-account</v-icon></td>
              <td>Accounts</td>
            </tr>
            <tr>
              <td><v-icon>mdi-security</v-icon></td>
              <td>TLS Policys</td>
            </tr>
            <tr>
              <td><v-icon>mdi-view-dashboard-variant</v-icon></td>
              <td>Dashboard</td>
            </tr>
            </tbody>
          </template>
        </v-table>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import Chart from "../components/Chart";
import { ref, onMounted } from 'vue'
import Client from "../service/Client";

const stats = ref([])

onMounted(async () => {
  Client.getStats().then((res) => {
    stats.value = res.data
  })
})
</script>
