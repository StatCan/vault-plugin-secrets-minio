# vault-plugin-secrets-minio

This is a plugin for [HashiCorp Vault][vault] which will provision multi-user keys for the [Minio object storage][minio] server.

## Usage

Once the plugin is registered with your vault instance, you can enable it on a particular path:

```sh
vault secrets enable \
  -path=minio \
  -plugin-name=vault-plugin-secrets-minio \
  -description="Instance of the Minio plugin" \
  plugin
```

### Configuration

In order to configure the plugin instance, you must supply it with your Minio endpoint, the access key ID, and the secret access key for the Minio initial user.

```sh
vault write minio/config \
  endpoint=<minio ip>:<minio port> \
  accessKeyId=<minio access key ID> \
  secretAccessKey=<minio secret access key> \
  useSSL=<true|false>
```

You can read the current configuration:

```sh
vault read minio/config
```

### Roles

Before you can issue keys, you must define a role. A role defines the policy which will be applied to the newly created user, and a name prefix for the key.

```sh
vault write minio/roles/example-role \
  policy=<existing minio policy name>
  user_name_prefix=<user name prefix>
```

The `<user name prefix>` is prefixed to the Vault request id for a key request, and defaults to an empty string. Having the Vault request id as the latter part of the name allows you to trace the key issuer via the Vault audit log. You may also optionally supply a `default_ttl` and `max_ttl` which will apply to the lease created by this role.

```sh
vault read minio/roles/example-role
```

Returns the configuration for a particular role.

```sh
vault list minio/roles
```

Lists all configured roles.

### Provisioning keys

```sh
vault read minio/keys/example-role
```

Returns the accessKeyId, secretAccessKey, policy and account status for the newly generated key.

## Acknowledgements

Leverages and builds upon the amazing work done by [Kula][kula]

[kula]:       https://github.com/kula/vault-plugin-secrets-minio
[vault]:      https://vaultproject.io
[minio]:      https://minio.io
