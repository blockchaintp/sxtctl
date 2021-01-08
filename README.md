## sxtctl

CLI to manage sextant clusters via CI or your terminal.

### install

```bash

sudo curl -o /usr/local/bin/sxtctl https://github.com/catenasys/sxtctl/releases/latest/download/sxtctl-linux-amd64
sudo chmod a+x /usr/local/bin/sxtctl
```

### usage

```bash
sxtctl
======

Manage your sextant installation.

Usage:
  sxtctl [command]

Available Commands:
  cluster     Manage k8s clusters
  deployment  Manage sextant deployments
  help        Help about any command
  remote      Manage your sextant credentials

Flags:
  -h, --help            help for sxtctl
      --remote string   The name of the remote sextant api
      --token string    Your testfaster access token
      --url string      The URL of the remote sextant api

Use "sxtctl [command] --help" for more information about a command.
```

### configure

You can use `sxtctl` right out of the box by setting the following environment variables (useful if running in CI):

 * `SEXTANT_URL` - the http(s) URL of the remote sextant API server
 * `SEXTANT_TOKEN` - the access token for the remote sextant API server

### remotes

To help manage multiple sextant api servers - sxtctl supports `remotes` which is will remember for next time.

If you do not provide `--url` or `SEXTANT_URL` (plus token) variables - sextant will use the currently active remote from your saved list.

To view your current remotes:

```bash
$ sxtctl remote list
```

Example output:

```
+------------+------------------+
|    NAME    |       URL        |
+------------+------------------+
| apples (*) | http://localhost |
+------------+------------------+
```

To add a new remote:

```bash
$ sxtctl remote add apples --url https://my.sextant.api.com --token XXXX
```

To remove a remote:

```bash
$ sxtctl remote remove apples
```

To switch to use a different remote by default:

```bash
$ sxtctl remote use apples
```

### clusters

To view a list of clusters on the current remote:

```bash
$ sxtctl cluster list
```

Example output:

```
+----+------+--------------------------+-------------+---------------------------------+-------------+
| ID | NAME |         CREATED          |   STATUS    |           API SERVER            | DEPLOYMENTS |
+----+------+--------------------------+-------------+---------------------------------+-------------+
|  3 | kind | 2021-01-07T10:36:06.742Z | provisioned | https://kind-control-plane:6443 |           1 |
+----+------+--------------------------+-------------+---------------------------------+-------------+
```

If you want to output the list of clusters in JSON format (useful when in CI):

```bash
$ sxtctl cluster list -o json
```

### deployments

To view a list of deployments on a cluster, you must provide either the cluster id or name.

This should be assigned to the `--cluster` command line argument or `SEXTANT_CLUSTER` env variable:

For example - in the output above - we have a cluster called `kind` with an id of `3`

Therefore - all of the following work the same way

```
$ sxtctl deployment list --cluster kind
$ sxtctl deployment list --cluster 3
$ SEXTANT_CLUSTER=kind sxtctl deployment list
$ SEXTANT_CLUSTER=3 sxtctl deployment list
```

Example output:

```
+----+-------+--------------------------+-------------+------------------+-------------------+
| ID | NAME  |         CREATED          |   STATUS    |  DEPLOYMENTTYPE  | DEPLOYMENTVERSION |
+----+-------+--------------------------+-------------+------------------+-------------------+
| 15 | daml2 | 2021-01-07T16:02:47.269Z | provisioned | daml-on-sawtooth |               1.3 |
+----+-------+--------------------------+-------------+------------------+-------------------+
```

If you want to output the list of deployments in JSON format (useful when in CI):

```bash
$ sxtctl deployment list --cluster 3 -o json
```

### pause & restart deployments

If you want to "pause" a deployment (remove all containers but keep the persistent volumes):

```bash
$ sxtctl deployment undeploy --cluster 3 --deployment 15
```

Then - later, if you want to reactivate a dedployment (i.e. reinstantiate all containers)

```bash
$ sxtctl deployment redeploy --cluster 3 --deployment 15
```

### examples

#### scale to zero nodes

In the following example - we configure everything via env variables (typical for a CI environment).

We can then run our two functions in order to pause and re-activate a deployment.

```
export SEXTANT_URL=https://my.sextant.com
export SEXTANT_TOKEN=XXX
export SEXTANT_CLUSTER=cluster-name
export SEXTANT_DEPLOYMENT=deployment-name

function pause-deployment() {
  sxtctl deployment undeploy
}

function restart-deployment() {
  sxtctl deployment redeploy
}
```