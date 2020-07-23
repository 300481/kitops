# kitops

Kubernetes Git-Ops

This Operator will do nearly the same as [WeaveWorks Flux](https://github.com/fluxcd/flux)

## Key differences to Flux

* triggered by API endpoints and not by time

* send notifications for all deployments including Helm

## TODOs

A Lot ;-)

Currently it is very very basic


## Install Operator

```bash
kubectl apply -f deploy/all.yaml
```
