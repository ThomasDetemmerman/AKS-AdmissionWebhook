# AKS-AdmissionWebhook
An admission webhook for pod mutations in AKS

# Structure
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