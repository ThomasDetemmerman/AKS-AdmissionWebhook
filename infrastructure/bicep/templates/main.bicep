param name string
param agentPoolProfiles array
param AdminAccount string
param K8SVersion string
@secure()
param AdminSSHPublicKey string



module AKS '../modules/aks.bicep' = {
  name: 
  params: {
    AdminAccount: AdminAccount
    AdminSSHPublicKey: AdminSSHPublicKey
    agentPoolProfiles: agentPoolProfiles
    K8SVersion: K8SVersion
    name: name
  }
}
