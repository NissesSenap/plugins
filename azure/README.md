# Azure k8s audit

The Azure k8s audit plugin is mainly supposed to use MSI to authenticate towards Azure but it's also possible to use environment variables.

## Develop

Assuming that your Azure account have access to the workspace.
Remember that you can use `AZURE_CONFIG_DIR` to point to your Azure config dir if you don't use the default.

```shell
export AZURE_RESOURCE_GROUP=example
export AZURE_WORKSPACE_NAME=example
export AZURE_SUBSCRIPTION_ID=11111-11111-1111-1111
```

If you don't use a SP

```shell
export AZURE_RESOURCE_GROUP=
export AZURE_SUBSCRIPTION_ID=
export AZURE_CLIENT_ID=
export AZURE_CLIENT_SECRET=
export AZURE_TENANT_ID=
```
