param name string
param agentPoolProfiles array
param AdminAccount string
param K8SVersion string
param location string = 'west europe' //= resourceGroup().location
@secure()
param AdminSSHPublicKey string




resource aksCluster 'Microsoft.ContainerService/managedClusters@2021-03-01' = {
  name: name
  location: location
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    kubernetesVersion: K8SVersion
    dnsPrefix: 'dnsprefix'
    enableRBAC: true
    agentPoolProfiles: agentPoolProfiles
    linuxProfile: {
      adminUsername: AdminAccount
      ssh: {
        publicKeys: [
          {
            keyData: AdminSSHPublicKey
          }
        ]
      }
    }
  }
}

