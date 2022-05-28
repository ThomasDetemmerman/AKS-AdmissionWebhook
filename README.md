# Admission webhook
An admission webhook for pod mutations in AKS.

## use case
You want to reduce cost by using spot VM's for your AKS cluster in non production evironments. Since spot vm's have a taint configured on them, you have to ask your developpers to add a toleration and node affinity to all of their deployments. Alternativally, you can enforce this using admission controllers. 

## Acknowledgement
The build-in admission controller "PodTolerationRestriction" can be used for the toleration part ([source](https://docs.microsoft.com/en-us/azure/aks/faq#what-kubernetes-admission-controllers-does-aks-support-can-admission-controllers-be-added-or-removed)). To achieve this, all your namespaces should have the proper labels ([more info](https://docs.microsoft.com/en-us/azure/aks/faq#what-kubernetes-admission-controllers-does-aks-support-can-admission-controllers-be-added-or-removed)). However, this does not cover the node affinity part. Of course, if you have only spot VM's, that issue is also resolved by itself. But what is the fun in that? So here is the solution with a custom admission webhook.


# Guide
## CI/CD pipeline
![Workflow](docs/img/workflow.png)
location: [.github/workflows](.github/workflows)

## Deploy AKS
First of all, we need to deploy our AKS cluster using spot VM.

location: [infrastructure/bicep](infrastructure/bicep)
sources:
- [Add a spot node pool to an Azure Kubernetes Service (AKS) cluster](https://docs.microsoft.com/en-us/azure/aks/spot-node-pool)
- [How to deploy bicep with Github actions](https://docs.microsoft.com/en-us/azure/azure-resource-manager/bicep/deploy-github-actions?tabs=CLI
)


## Admission webhook in GO
Here comes the hard part. Luckly, I had [this amazing tutorial](https://slack.engineering/simple-kubernetes-webhook/) wich basically covers everyting you need. So I cloned [their repository](https://github.com/slackhq/simple-kubernetes-webhook) and started working from there. I used [GoLand by Jetbrains](https://www.jetbrains.com/go/) for my IDE which basically did the coding for me. I only had to  point the IDE in the right direction ;).

The part that we need to alter is located in [Webhook/pkg/mutation](Webhook/pkg/mutation). If you have a look at the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/aks/spot-node-pool#verify-the-spot-node-pool), you see that the pod should have the following fields:

```yml
spec:
  containers:
  - name: spot-example
  tolerations:
  - key: "kubernetes.azure.com/scalesetpriority"
    operator: "Equal"
    value: "spot"
    effect: "NoSchedule"
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: "kubernetes.azure.com/scalesetpriority"
            operator: In
            values:
            - "spot"
   ...
```
So we will implement the `podMutator` interface twice. This interface is defined in [mutation.go](Webhook/pkg/mutation/mutation.go) as follows:
```go
type podMutator interface {
	Mutate(*corev1.Pod) (*corev1.Pod, error)
	Name() string
}
```
So in that directory, create two new files which I called `spotVM_affinity.go` and `spotVM_toleration.go`. The hard part is obviously constructing the toleration object and affinity object. Both objects should be created in the `Mutate` function. 
The following tips might help:
- Make exsessivly use of the GoLand autocompletion.
- Keep the desired YAML configuration nearby. As you can see from the snippet below, the GoLang code is very close to the YAML structure.

Below snippet is for illustrating purposes. For the full implementation of the interface, have a look at [spotVM_toleration.go](Webhook/pkg/mutation/spotVM_toleration.go) and [spotVM_affinity.go](Webhook/pkg/mutation/spotVM_affinity.go)
```go
spotvmAffinity := corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: nil
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{{
				Weight: 1,
				Preference: corev1.NodeSelectorTerm{
					MatchExpressions: []corev1.NodeSelectorRequirement{{
						Key:      "kubernetes.azure.com/scalesetpriority",
						Operator: corev1.NodeSelectorOpIn,
						Values:   []string{"Spot"},
					}},
					MatchFields: nil,
				},
			}},
		},
	}
```

Once you have implemented these interfaces, you should call them. This is done in [mutation.go](Webhook/pkg/mutation/mutation.go)
```go
// list of all mutations to be applied to the pod
	mutations := []podMutator{
		minLifespanTolerations{Logger: log},
		affinity{Logger: log},
	}
```


Sources: 
  - [simple kubernetes webhook](https://slack.engineering/simple-kubernetes-webhook/)
## Multi stage docker file
Hereafter, it is time to build and run the docker file. Building this API is best done using a multistage docker file. This implies that you have a docker image with all the necesairy Go libraries present. However, this makes the image large and all these files are not needed when running the application. So once we have build the application, we copy it to a small alpine image.

location: [Webook/Dockerfile](Webook/Dockerfile)
Sources: 
- [multistage build](https://docs.docker.com/develop/develop-images/multistage-build/)
- [multistage dockerfile for Go](https://slack.engineering/simple-kubernetes-webhook/)


## Configuring AKS
At last, it is time to deploy our application. Obviously, we need to deploy the deployment and put a service in front of it.
However, you should also tell AKS for wich namespaces you want to perform a mutation. This can be done using this [MutatingWebhookConfiguration](infrastructure/AKSConfig/MutatingWebhookConfiguration.yml).
There are two important remarks to this. The first one is the `namespaceSelector` which **SHOULD** exclude the namespace in which your admission webhook is deployed. Otherwise you have a vicous circle. If, due to a bug, your admission webhook needs to be redeployed, AKS will try to validate it agains himself. 
```
 - name: "customadmissionwebhook.example.com"
   namespaceSelector:
      matchLabels:
        admission-webhook: enabled
```
A second important remark are the SSL certificates. AKS requires you to perform the call using HTTPS. When using HTTPS, Golang enforces you to use SAN (and not only CN). This leads to the following command:
```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=admission-webhook-service.default.svc" -addext "subjectAltName = DNS:admission-webhook-service.default.svc" 
```
Automatically injecting the certificates in AKS was done with the following task in the [Github Workflow](.github/workflows/main.yml)
``` bash
-   name: Generate, inject and deploy TLS certificate
    shell: bash
    run: |  	
      kubectl delete secret admission-webhook-tls-secret
      openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=admission-webhook-service.default.svc" -addext "subjectAltName = DNS:admission-webhook-service.default.svc" 
      kubectl create secret tls admission-webhook-tls-secret --key="tls.key" --cert="tls.crt"
      tlscrt=$(cat tls.crt | base64 -w 0)
      sed -i -- "s/#{tobeoverriden-tlscrt}#/$(echo $tlscrt)/g" infrastructure/AKSConfig/mutatingAdmissionWebhookConfiguration.yml
     

```


location: [infrastructure/AKSConfig](infrastructure/AKSConfig)  
sources:
- [TLS secrets in AKS](https://shocksolution.com/2018/12/14/creating-kubernetes-secrets-using-tls-ssl-as-an-example/)
- [simple kubernetes webhook](https://slack.engineering/simple-kubernetes-webhook/)
- [admission controllers](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)