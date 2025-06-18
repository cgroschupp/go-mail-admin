import "@typespec/http";
import "@typespec/openapi";

using Http;
@service(#{ title: "Mail Admin Service"})
@server("http://localhost/api/v1")
namespace MailAdminService;

model Domain {
  @visibility(Lifecycle.Read)
  id: int32;
  @format("hostname")
  @OpenAPI.extension("x-oapi-codegen-extra-tags", #{validate: "hostname"})
  name: string;
  @visibility(Lifecycle.Read)
  created_at: utcDateTime;
  @visibility(Lifecycle.Read)
  updated_at: utcDateTime;
}

model DomainList {
  items: Domain[];
}

model Account {
  @visibility(Lifecycle.Read)
  id: int32,
  @visibility(Lifecycle.Read)
  domain: Domain,
  @visibility(Lifecycle.Read,Lifecycle.Create)
  domain_id: int32,
  username: string,
  @visibility(Lifecycle.Create, Lifecycle.Update)
  password: string,
  quota?: int32,
  sendonly?: boolean,
  enabled?: boolean,
  @visibility(Lifecycle.Read)
  created_at: utcDateTime,
  @visibility(Lifecycle.Read)
  updated_at: utcDateTime,
}

model AccountList {
  items: Account[];
}

model Alias {
  @visibility(Lifecycle.Read)
  id: int32,
  enabled: boolean,
  source_domain_id: int32,
  @visibility(Lifecycle.Read)
  source_domain: Domain,
  source_username?: string,
  @visibility(Lifecycle.Read)
  source_display: string,
  destination_domain: string,
  destination_username: string,
  @visibility(Lifecycle.Read)
  destination_display: string,
  @visibility(Lifecycle.Read)
  created_at: utcDateTime,
  @visibility(Lifecycle.Read)
  updated_at: utcDateTime,
}

model AliasList {
  items: Alias[];
}

model TLSPolicy {
  @visibility(Lifecycle.Read)
  id: int32,
  @visibility(Lifecycle.Read,Lifecycle.Create)
  domain_id: int32,
  @visibility(Lifecycle.Read)
  domain: Domain,
  params?: string,
  policy: "none" | "may" | "encrypt" | "dane" | "dane-only" | "fingerprint" | "verify" | "secure",
}

model TLSPolicyList {
  items: TLSPolicy[];
}

@error
model Error {
  code: int32;
  message: string;
  error: string;
}

model ChangePasswordRequest {
  password: string;
}

@route("/domain")
@tag("Domains")
interface Domains {
  /** List domains */
  @get list(): DomainList | Error;
  /** Read domains */
  @get read(@path id: int32): Domain | Error;
  /** Create a domain */
  @post create(@body body: Domain): Domain | Error;
  /** Update a domain */
  @patch update(@path id: int32, @body body: MergePatchUpdate<Domain>): Domain | Error;
  /** Delete a domain */
  @delete delete(@path id: int32): void | Error;
}

@route("/account")
@tag("Accounts")
interface Accounts {
    /** List accounts */
  @get list(): AccountList | Error;
    /** Read account */
  @get read(@path id: int32): Account | Error;
    /** Create an account */
  @post create(@body body: Account): Account | Error;
  /** Update an account */
  @patch update(@path id: int32, @body body: MergePatchUpdate<Account>): Account | Error;
  /** Delete an account */
  @delete delete(@path id: int32): void | Error;

  @route("{id}/password") @put change_password(@path id: int32, @body body: ChangePasswordRequest): null | Error;
}

@route("/alias")
@tag("Aliases")
interface Aliases {
    /** List aliases */
  @get list(): AliasList | Error;
    /** Read alias */
  @get read(@path id: int32): Alias | Error;
    /** Create an alias */
  @post create(@body body: Alias): Alias | Error;
  /** Update an alias */
  @patch update(@path id: int32, @body body: MergePatchUpdate<Alias>): Alias | Error;
  /** Delete an alias */
  @delete delete(@path id: int32): void | Error;
}
@route("/tlspolicy")
@tag("Tlspolicy")
interface TLSPolicies {
    /** List aliases */
  @get list(): TLSPolicyList | Error;
    /** Read alias */
  @get read(@path id: int32): TLSPolicy | Error;
    /** Create an alias */
  @post create(@body body: TLSPolicy): TLSPolicy | Error;
  /** Update an alias */
  @patch update(@path id: int32, @body body: MergePatchUpdate<TLSPolicy>): TLSPolicy | Error;
  /** Delete an alias */
  @delete delete(@path id: int32): void | Error;
}

model LoginParameter {
  username: string;
  password: string
}

model LoginResponse {
  login: boolean;
  token: string;
}

@tag("User")
interface UserOperations {
  @tag("Public")
  @route("login")
  @post login(@body body: LoginParameter): LoginResponse | Error;
  @route("logout")
  @post logout(): boolean | Error;
}

model DashboardStatsItem {
  labels: string[];
  datasets: DashboardDataset[];
}

model DashboardDataset {
  data: int32[];
  backgroundColor: string[];
}

model DashboardStats {
  ...Record<DashboardStatsItem>;
}

model DashboardStatus {
  healthy: boolean
}

model DashboardVersion {
  version: string
}

@tag("Dashboard")
interface DashboardOperations {
  @route("stats")
  @get stats(): DashboardStats | Error;
  @tag("Public")
  @route("status")
  @get status(): DashboardStatus | Error;
  @tag("Public")
  @route("version")
  @get version(): DashboardVersion | Error;
}