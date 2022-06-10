# Preflight Controller

This is a very-alpha controller that runs [Openshift
Preflight](https://github.com/redhat-openshift-ecosystem/openshift-preflight) as
a job in a Kubernetes Cluster on creation of a
[PreflightCheck](./config/samples/sample-default.yaml) resource.

## Quick Run

Assuming KUBECONFIG is set, and a running cluster with admin privs (because the
controller runs on your machine)

```shell
# Clone the repo
git clone https://github.com/opdev/preflight-controller && cd preflight-controller

# Install the CRD
make install && make run

# Apply the samples (or configure your own)
make install-samples

# Watch the samples do their thing (they don't all succeed)
watch oc get preflightcheck
# sample output is below

# When done, clean up the samples.
make uninstall-samples
```

Sample output of watch command:

```
$ oc get preflightcheck
NAME                       IMAGE                                      TYPE        SUCCESSFUL   JOB
sample-custom-dockercfg    quay.io/opdev/simple-demo-operator:0.0.6   container   true         sample-custom-dockercfg-gfs4q
sample-custom-pflt-img     quay.io/opdev/simple-demo-operator:0.0.6   container   true         sample-custom-pflt-img-pg6fl
sample-default             quay.io/opdev/simple-demo-operator:0.0.6   container   true         sample-default-kqjvq
sample-operator-checks     quay.io/opdev/simple-demo-operator:0.0.6   operator    false        sample-operator-checks-m6kh7
sample-with-cert-proj-id   quay.io/opdev/simple-demo-operator:0.0.6   container   false        sample-with-cert-proj-id-ngfq2
```

## Notes on behavior

- This spins up a single job for a given `PreflightCheck`. If an job already exists for an instance, no other jobs are created.
- This will attempt the job a single time.
- The type of Preflight checks to apply depend on which `PreflightCheck.Spec.CheckOptions` exist in your manifest.
- Credentials aren't fully wired up, so they will not function as expected as of the time of this writing.
