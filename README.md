# bridge-cni

A very simple tool to generate a CNI config that uses the bridge CNI plugin with the PodCIDR assigned by Kubernetes.
It is usually deployed as a `DaemonSet` and generates the file `/etc/cni/net.d/bridge-cni.conflist` with something like this:

```json
{
  "cniVersion": "1.0.0",
  "name": "cbr0",
  "plugins": [
    {
      "type": "bridge",
      "ipam": {
        "type": "host-local",
        "subnet": "2001:db8::c:0/120"
      },
      "dns": {},
      "isDefaultGateway": true
    }
  ]
}
```

## Usage

### Kubeadm

```yaml
apiVersion: kubeadm.k8s.io/v1beta4
kind: ClusterConfiguration
kubernetesVersion: v1.31.0
controllerManager:
  extraArgs:
    "node-cidr-mask-size": "120"
networking:
  podSubnet: 2001:db8::c:0/112
  serviceSubnet: 2001:db8::b:0/112
```

Apply the following YAML after cluster initialization:

```bash
kubectl apply -f https://raw.githubusercontent.com/lion7/bridge-cni/refs/heads/main/deploy/bridge-cni.yaml
```

Note: make sure to set the node CIDR mask size to something smaller than the prefix size you use for pods.

### Talos Linux

```bash
cluster:
  controllerManager:
    extraArgs:
      node-cidr-mask-size: 120
  network:
    cni:
      name: custom
      urls:
        - https://raw.githubusercontent.com/lion7/bridge-cni/main/deploy/bridge-cni.yaml
    podSubnets:
      - 2001:db8::c:0/112
    serviceSubnets:
      - 2001:db8::b:0/112
```

Note: make sure to set the node CIDR mask size to something smaller than the prefix size that you use for pods.
