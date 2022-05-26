# AKS-AdmissionWebhook
An admission webhook for pod mutations in AKS

# Structure
## CI/CD pipeline
.github/workflows

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