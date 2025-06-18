<template>
  <v-container>
    <v-card>
      <v-card-title>TLS Policy</v-card-title>
      <v-card-text>
        <v-select ref="domain" :items="domains" data-app="true" label="Domain" v-model="domain_id"
          :item-title="item => item.name" :item-value="item => item.id" :disabled="this.$route.params.id"></v-select>
        <label>Policy</label>
        <v-select :items="policyOptions" label="Policy" v-model="policy" data-app="true"></v-select>
        <v-text-field type="text" v-model="params" label="Params"></v-text-field>
      </v-card-text>
      <v-card-actions>
        <v-btn @click="savePolicy()">Save Policy</v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script>
// @ is an alias to /src
import Client from "../service/Client";

export default {
  name: 'TLSPolicyEditView',
  components: {
  },
  methods: {
    'savePolicy': function () {
      if (this.$route.params.id) {
        Client.updateTLSPolicy(this.$route.params.id, { "domain": this.domain, "policy": this.policy, "params": this.params }).then(() => {
          this.$notify.info("Policy saved", { location: 'bottom center' });
        }).catch((res) => {
          this.$notify.error("Oups, something go wrong" + res)
        });
      } else {
        Client.createTLSPolicy({ "domain_id": this.domain_id, "policy": this.policy, "params": this.params }).then(() => {
          this.$notify.info("Policy saved", { location: 'bottom center' });
        }).catch((res) => {
          this.$notify.error("Oups, something go wrong" + res)
        });
      }
      this.$router.push("/tls")

    },
    getDomains: function () {
      Client.getDomains().then((res) => {
        this.domains = res.data.items
      });
    },
    getPolicy: function () {
      Client.getTLSPolicy(this.$route.params.id).then((res) => {
        this.domain_id = res.data.domain_id
        this.policy = res.data.policy
        this.params = res.data.params
      });
    }
  },
  mounted: function () {
    this.getDomains();
    if (this.$route.params.id) {
      this.getPolicy();
    }
  },
  data: () => ({
    domain_id: null,
    domains: [],
    policy: "",
    params: "",
    policyOptions: ['none', 'may', 'encrypt', 'dane', 'dane-only', 'fingerprint', 'verify', 'secure']
  })
}
</script>
