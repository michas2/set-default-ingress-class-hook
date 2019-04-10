# What

This creates a webhook, which will set a default ingress class whenever you create an ingress without any class.

# Why

If you have a cluster with multiple ingress controller, you must specify for each ingress object, or multiple controllers will fight for that object.
This hook makes sure each object does specify a class. -> No more fighting.

# How

Using the [github plugin](https://github.com/sagansystems/helm-github) you can install it like this:
```
helm github install --repo git@github.com:michas2/set-default-ingress-class-hook.git --path helm/ingressInjectorWebhook --set class="nginx"
```

Simply specify which ingress class you want to set as default. Then whenever a classless ingress object is created, the hook will add the default ingress class.

* The `image` directory is automatically build by docker hub as `michas2/set-default-ingress-class-hook`.
* The `helm` directory contains the chart which will use that image.
* The `cert.sh` script creates a tls certificate for secure communication between api-server and webhook.

Currently you also need to `--set caBundle=...` to the certificate created by the script.
You currently probably want to use it like this:
```
$ ./cert.sh
[...]
Please set caBundle to 'LS0tXXX'.
$ helm install helm/ingressInjectorWebhook --set class="nginx" --set caBundle="LS0tXXX"
```
Looks like there is no good way to automate this in the helm chart...
