<template>
    <div>
        <v-container>
            <v-card>
                <v-card-title v-if="account.id">Account {{account.username}}@{{account.domain.name}} </v-card-title>
                <v-card-text>
                    <span v-if="!account.id">
                        <v-text-field v-model="account.username" v-on:keyup="checkatsymbole" label="Username" placeholder="Username"></v-text-field>
                        <v-select v-model="account.domain_id" data-app="true" :item-title="item=>item.name" :item-value="item=>item.id" :items="domainNames" label="Destination-Domain"></v-select>
                        <v-text-field v-on:keyup="passwordFieldChanged()" v-on:change="passwordFieldChanged()" v-model="account.password" label="Password" :type="passwordFieldType" placeholder="Password"></v-text-field>
                        <v-btn @click="generateRandomPassword()" size="x-small">Random Password</v-btn>
                    </span>
                    <v-text-field v-model.number="account.quota" label="Quota" placeholder="Quota" type="number"></v-text-field>
                    <v-checkbox v-model="account.enabled" label="Enabled"></v-checkbox>
                    <v-checkbox v-model="account.sendonly" label="Send Only"></v-checkbox>
                    <v-card-actions>
                        <v-btn @click="saveAccount()"  prepend-icon="mdi-content-save">Save</v-btn>
                    </v-card-actions>
                </v-card-text>
            </v-card>
        </v-container>
        <v-container v-if="account.id">
            <v-card>
                <v-card-title>Change Password</v-card-title>
                <v-card-text>
                    <v-text-field v-model="password" :type="passwordFieldType" v-on:keyup="passwordFieldChanged()" v-on:change="passwordFieldChanged()" label="New Password" placeholder="New Password"></v-text-field>
                    <v-btn @click="generateRandomPassword()" size="x-small">Random Password</v-btn>
                </v-card-text>
                <v-card-actions>
                    <v-btn @click="changePassword()"  prepend-icon="mdi-content-save">Change Password</v-btn>
                </v-card-actions>
            </v-card>
        </v-container>
    </div>
</template>

<script>
    import Client from "../service/Client";
    export default {
        name: 'AccountEditView',
        methods: {
            checkatsymbole: function (r) {
                if(r.key == "@") {
                    this.account.username = this.account.username.substr(0, this.account.username.length -1 );
                    this.$refs.domain.focus();
                }

            },
            generateRandomPassword: function () {
                this.passwordFieldType = "password"
                this.passwordFieldType = "text"
                this.generate();
            },
            passwordFieldChanged: function() {
                if(this.account.password != this.lastRandomPassword) {
                    this.passwordFieldType = "password";
                }
            },
            generate () {
                let CharacterSet = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789![]{}()%&*<>~.';
                let password = '';

                for(let i=0; i < 25; i++) {
                    password += CharacterSet.charAt(Math.floor(Math.random() * CharacterSet.length));
                }
                this.account.password = password;
                this.password = password;
                this.lastRandomPassword = password;
            },
            getAccounts: function () {
                Client.getAccount(this.$route.params.id).then((res) => {
                    this.account = res.data
                });
            },
            getDomains: function () {
                Client.getDomains().then((res) => {
                    this.domainNames = res.data.items
                });

            },
            saveAccount: function () {
                if (this.account.quota === '') {
                   delete this.account.quota
                }
                if(this.account.id) {
                    Client.saveAccount(this.account.id, this.account).then(() => {
                        this.getAccounts();
                        this.$notify.info("Account saved");
                        this.$router.push("/account")
                    })
                } else {
                    Client.createAccount(this.account).then(() => {
                        this.$notify.info("Account created");
                        this.$router.push("/account")
                    })
                }
            },
            changePassword: function () {
                Client.changePassword(this.account.id, this.password).then(()=> {
                    this.$notify.info("Password changed");
                }).catch(() => {
                    this.$notify.error("Oups, something go wrong")
                })
            }
        },

        mounted: function() {
            this.getDomains();
            if (this.$route.params.id) {
                this.getAccounts();
            }
        },
        components: {

        },
        data: () => ({
            account: { "quota": 1024, "enabled": true },
            password: '',
            domainNames: [],
            passwordFieldType: "password",
            lastRandomPassword: ""
        }),
    }
</script>
