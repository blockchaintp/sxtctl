# sxtctl

[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B9585%2Fgit%40github.com%3Acatenasys%2Fsxtctl.git.svg?type=shield)](https://app.fossa.com/projects/custom%2B9585%2Fgit%40github.com%3Acatenasys%2Fsxtctl.git?ref=badge_shield)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=catenasys_sxtctl&metric=alert_status&token=970d267a29a53f6e5f14cbcbb54e65a6001e90cd)](https://sonarcloud.io/dashboard?id=catenasys_sxtctl)

`sxtctl` is CLI to manage sextant clusters via CI or your terminal.

## Install

```bash
# Set ARCH appropriately for your platform
# Supported platforms:
#   darwin-amd64
#   linux-amd64
#   linux-arm64
#   windows-amd64

ARCH=linux-amd64
VERSION=0.5.3
sudo curl -o /usr/local/bin/sxtctl \
  https://sxtctl.s3.amazonaws.com/${VERSION}/sxtctl-${ARCH}
sudo chmod a+x /usr/local/bin/sxtctl
```

## Usage

```text
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
      --token string    Your Sextant access token
      --url string      The URL of the remote sextant api

Use "sxtctl [command] --help" for more information about a command.
```

## Configuration

You can use `sxtctl` right out of the box by setting the following environment
variables (useful if running in CI):

* `SEXTANT_URL` - the http(s) URL of the remote sextant API server including
  the port if necessary

* `SEXTANT_TOKEN` - the access token for the remote sextant API server

## Configuring Remotes

To help manage multiple sextant api servers - sxtctl supports `remotes` which
it stores in a config file in much the same way that `kubectl` stores its
cluster configs. This is normally stored in `~/.sextant` directory.

If you do not provide `--url` or `SEXTANT_URL` (plus token) variables - sextant
will use the currently active remote from your saved list.

### List Remotes

To view your current remotes:

```bash
$ sxtctl remote list

+------------+------------------+
|    NAME    |       URL        |
+------------+------------------+
| apples (*) | http://localhost |
+------------+------------------+
```

### Add a Remote

To add a new remote:

```bash
sxtctl remote add apples --url https://my.sextant.api.com --token XXXX
```

### Remove a Remote

To remove a remote:

```bash
sxtctl remote remove apples
```

### Switching the Remote

To switch to use a different remote by default:

```bash
sxtctl remote use apples
```

## Cluster Operations

### List Clusters

To view a list of clusters on the current remote:

```bash
$ sxtctl cluster list

+----+------+--------------------------+-------------+---------------------------------+-------------+
| ID | NAME |         CREATED          |   STATUS    |           API SERVER            | DEPLOYMENTS |
+----+------+--------------------------+-------------+---------------------------------+-------------+
|  3 | kind | 2021-01-07T10:36:06.742Z | provisioned | https://kind-control-plane:6443 |           1 |
+----+------+--------------------------+-------------+---------------------------------+-------------+
```

If you want to output the list of clusters in JSON format (useful when in CI):

```bash
sxtctl cluster list -o json
```

## Deployment Operations

### List Deployments

To view a list of deployments on a cluster, you must provide either the cluster
id or name.

This should be assigned to the `--cluster` command line argument or using the
`SEXTANT_CLUSTER` env variable:

For example - in the output above - we have a cluster called `kind` with an id
of `3`

Therefore - all of the following work the same way

```bash
$ sxtctl deployment list --cluster kind
$ sxtctl deployment list --cluster 3
$ SEXTANT_CLUSTER=kind sxtctl deployment list
$ SEXTANT_CLUSTER=3 sxtctl deployment list

+----+-------+--------------------------+-------------+------------------+-------------------+
| ID | NAME  |         CREATED          |   STATUS    |  DEPLOYMENTTYPE  | DEPLOYMENTVERSION |
+----+-------+--------------------------+-------------+------------------+-------------------+
| 15 | daml2 | 2021-01-07T16:02:47.269Z | provisioned | daml-on-sawtooth |               1.3 |
+----+-------+--------------------------+-------------+------------------+-------------------+
```

If you want to output the list of deployments in JSON format (useful when in
CI):

```bash
sxtctl deployment list --cluster 3 -o json
```

### Undeploy a deployment

If you want to undeploy a provisioned deployment on a cluster (i.e. terminate
all its pods but retain the state of the underlying blockchain) then in addition
to supplying details of cluster as above you need to provide either the id or
the name of the deployment using the `--deployment` command line argument or
using the `SEXTANT_DEPLOYMENT` env variable:

```bash
sxtctl deployment undeploy --cluster 3 --deployment 15
```

If you list the deployments you will see that the status of the deployment is
_deleted_ once it is undeployed.

_NOTE_ with DAML on Besu you need to ensure that persistence is enabled.

### Re-deploy a deployment

If you want to redeploy a deleted deployment i.e. re-provision all its pods:

```bash
sxtctl deployment redeploy --cluster 3 --deployment 15
```

If you list the deployments you will see that the status of the deployment is
_provisioned_ once it is redeployed.

## Examples

### Scale a cluster to zero worker nodes

In the following example we configure everything via env variables (typical for
a CI environment) and define functions that suspend and resume a deployment
respectively.

```bash
export SEXTANT_URL=https://my.sextant.com
export SEXTANT_TOKEN=XXX
export SEXTANT_CLUSTER=cluster-name
export SEXTANT_DEPLOYMENT=deployment-name

function suspend-deployment() {
  sxtctl deployment undeploy
}

function resume-deployment() {
  sxtctl deployment redeploy
}
```

These functions can be integrated into automated workflows that:

* suspend a provisioned deployment; confirm that all the pods associated with
  the deployment have been terminated then scale back the cluster to zero worker
  nodes
* scale up cluster to the appropriate number of worker nodes; confirm that these
  are ready then resume a deleted deployment

_NOTE_ that the Kubernetes namespace for a deployment is provided if you list
deployments using JSON format.
