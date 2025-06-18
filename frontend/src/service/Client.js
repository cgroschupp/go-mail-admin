import Api from './Api'

export default {
    ping() {
        return Api().get("/api/v1/ping")
    },
    getStatus() {
        return Api().get("/api/v1/status")
    },
    getStats() {
        return Api().get(`/api/v1/stats`)
    },
    triggerDomains() {
        return Api().post("/api/v1/domain/trigger")
    },
    getDomains() {
        return Api().get("/api/v1/domain")
    },
    getDomainStats() {
        return Api().get(`/api/v1/domain/stats`)
    },
    getDomainDetails(id) {
        return Api().get(`/api/v1/domain/${id}`)
    },
    addDomain(domainname) {
        return Api().post("/api/v1/domain", { "name": domainname })
    },
    removeDomain(id) {
        return Api().delete(`/api/v1/domain/${id}`)
    },
    getAliases() {
        return Api().get("/api/v1/alias");
    },
    getAlias(id) {
        return Api().get(`/api/v1/alias/${id}`)
    },
    getAliasStats() {
        return Api().get(`/api/v1/alias/stats`)
    },
    saveAlias(id, data) {
        return Api().patch(`/api/v1/alias/${id}`, data, {headers: {"content-type":"application/merge-patch+json"}})
    },
    createAlias(data) {
        return Api().post("/api/v1/alias", data)
    },
    removeAlias(id) {
        return Api().delete(`/api/v1/alias/${id}`)
    },
    getAccounts() {
        return Api().get("/api/v1/account")
    },
    getAccountStats() {
        return Api().get(`/api/v1/account/stats`)
    },
    getAccount(id) {
        return Api().get(`/api/v1/account/${id}`)
    },
    saveAccount(id, data) {
        return Api().patch(`/api/v1/account/${id}`, data, {headers: {"content-type":"application/merge-patch+json"}})
    },
    createAccount(data) {
        return Api().post("/api/v1/account", data)
    },
    deleteAccount(id) {
        return Api().delete(`/api/v1/account/${id}`)
    },
    listTLSPolicies() {
        return Api().get("/api/v1/tlspolicy")
    },
    getTLSPolicy(id) {
        return Api().get(`/api/v1/tlspolicy/${id}`)
    },
    createTLSPolicy(data) {
        return Api().post("/api/v1/tlspolicy", data)
    },
    updateTLSPolicy(id, data) {
        return Api().patch(`/api/v1/tlspolicy/${id}`, data, {headers: {"content-type":"application/merge-patch+json"}})
    },
    deleteTLSPolicy(id) {
        return Api().delete(`/api/v1/tlspolicy/${id}`)
    },
    changePassword(id, newpassword) {
        return Api().put(`/api/v1/account/${id}/password`, { "password": newpassword })
    },
    login(username, password) {
        return Api().post("/api/v1/login", { "username": username, "password": password })
    },
    logout() {
        return Api().post("/api/v1/logout")
    },
    getVersion() {
        return Api().get("/api/v1/version")
    }
}
