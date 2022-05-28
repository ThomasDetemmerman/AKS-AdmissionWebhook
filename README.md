# Admission webhook
An admission webhook for pod mutations in AKS.

## use case
You want to reduce cost by using spot VM's for your AKS cluster in non production evironments. Since spot vm's have a taint configured on them, you have to ask your developpers to add a toleration and node affinity to all of their deployments. Alternativally, you can enforce this using admission controllers. 

## Acknowledgement
The build-in admission controller "PodTolerationRestriction" can be used for the toleration part ([source](https://docs.microsoft.com/en-us/azure/aks/faq#what-kubernetes-admission-controllers-does-aks-support-can-admission-controllers-be-added-or-removed)). To achieve this, all your namespaces should have the proper labels ([more info](https://docs.microsoft.com/en-us/azure/aks/faq#what-kubernetes-admission-controllers-does-aks-support-can-admission-controllers-be-added-or-removed)). However, this does not cover the node affinity part. Of course, if you have only spot VM's, that issue is also resolved by itself. But what is the fun in that? So here is the solution with a custom admission webhook.


# Guide
## CI/CD pipeline
.github/workflows

## Deploy AKS
https://docs.microsoft.com/en-us/azure/azure-resource-manager/bicep/deploy-github-actions?tabs=CLI

## Admission webhook in GO
Webhook
Sources: 
- https://slack.engineering/simple-kubernetes-webhook/
## Multi stage docker file
Webook/Dockerfile. 
Sources: 
- https://docs.docker.com/develop/develop-images/multistage-build/
- https://slack.engineering/simple-kubernetes-webhook/

## Deployment of infrastructure
infrastructure/bicep


## Configuring AKS
infrastructure/AKSConfig
sources:
- https://shocksolution.com/2018/12/14/creating-kubernetes-secrets-using-tls-ssl-as-an-example/
- https://slack.engineering/simple-kubernetes-webhook/
- https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/