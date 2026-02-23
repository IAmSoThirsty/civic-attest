# Public Truth Verification Service (PTVS) - Kubernetes Deployment

This directory contains the complete Kubernetes deployment configuration for the Public Truth Verification Service (PTVS), a production-grade deployment leveraging modern cloud-native technologies.

## Architecture Overview

The PTVS deployment is built with the following components:

- **Kubernetes 1.28+** - Container orchestration
- **Istio Service Mesh** - Strict mTLS, traffic management, and observability
- **ArgoCD** - GitOps-based continuous delivery
- **Argo Rollouts** - Canary deployment strategy
- **OPA Gatekeeper** - Policy enforcement
- **Prometheus** - Metrics collection
- **HPA** - Horizontal Pod Autoscaling

## Prerequisites

Before deploying PTVS, ensure your cluster has the following components installed:

1. **Kubernetes 1.28+**
   ```bash
   kubectl version --short
   ```

2. **Istio Service Mesh**
   ```bash
   istioctl version
   ```

3. **ArgoCD**
   ```bash
   kubectl get pods -n argocd
   ```

4. **Argo Rollouts Controller**
   ```bash
   kubectl get pods -n argo-rollouts
   ```

5. **OPA Gatekeeper**
   ```bash
   kubectl get pods -n gatekeeper-system
   ```

6. **Prometheus** (for metrics scraping)
   ```bash
   kubectl get pods -n monitoring
   ```

## Directory Structure

```
ptvs-platform/
├── bootstrap/               # Namespace and PodSecurity configuration
│   ├── namespace.yaml       # mankind namespace with Istio injection
│   └── podsecurity.yaml     # PodDisruptionBudget for high availability
├── mesh/                    # Istio service mesh configuration
│   ├── peer-auth.yaml       # Strict mTLS enforcement
│   ├── destination-rule.yaml # Traffic policy configuration
│   └── authorization-policy.yaml # Service-to-service authorization
├── policy/                  # OPA Gatekeeper policies
│   └── gatekeeper-constraints.yaml # Required labels constraint
├── helm/                    # Helm chart for PTVS
│   └── ptvs/
│       ├── Chart.yaml       # Helm chart metadata
│       ├── values.yaml      # Configuration values
│       └── templates/       # Kubernetes resource templates
│           ├── deployment.yaml
│           ├── service.yaml
│           ├── serviceaccount.yaml
│           ├── hpa.yaml
│           ├── pdb.yaml
│           ├── networkpolicy.yaml
│           ├── virtualservice.yaml
│           └── rollout.yaml
└── argocd/                  # ArgoCD application definition
    └── application.yaml     # GitOps application manifest
```

## Deployment Steps

### Step 1: Bootstrap the Namespace

Create the `mankind` namespace with Istio injection and PodSecurity enforcement:

```bash
kubectl apply -f ptvs-platform/bootstrap/namespace.yaml
kubectl apply -f ptvs-platform/bootstrap/podsecurity.yaml
```

Verify the namespace:
```bash
kubectl get namespace mankind -o yaml
```

### Step 2: Configure Istio Service Mesh

Apply Istio security policies for strict mTLS and authorization:

```bash
kubectl apply -f ptvs-platform/mesh/peer-auth.yaml
kubectl apply -f ptvs-platform/mesh/destination-rule.yaml
kubectl apply -f ptvs-platform/mesh/authorization-policy.yaml
```

Verify Istio configuration:
```bash
kubectl get peerauthentication -n mankind
kubectl get destinationrule -n mankind
kubectl get authorizationpolicy -n mankind
```

### Step 3: Apply OPA Gatekeeper Policies

Enforce required labels on all pods:

```bash
kubectl apply -f ptvs-platform/policy/gatekeeper-constraints.yaml
```

Verify the constraint:
```bash
kubectl get k8srequiredlabels
```

### Step 4: Deploy via ArgoCD (Recommended)

Deploy using GitOps:

```bash
kubectl apply -f ptvs-platform/argocd/application.yaml
```

Monitor the deployment:
```bash
kubectl get application -n argocd
kubectl describe application ptvs -n argocd
```

ArgoCD will automatically:
- Deploy the Helm chart from the repository
- Monitor for changes in the Git repository
- Automatically sync changes (self-heal enabled)
- Prune resources that are removed from Git

### Step 5: Manual Helm Deployment (Alternative)

If you prefer to deploy manually using Helm:

```bash
helm install ptvs ptvs-platform/helm/ptvs \
  --namespace mankind \
  --create-namespace
```

Or upgrade an existing deployment:
```bash
helm upgrade ptvs ptvs-platform/helm/ptvs \
  --namespace mankind
```

## Configuration

### Image Configuration

Update the image repository and tag in `ptvs-platform/helm/ptvs/values.yaml`:

```yaml
image:
  repository: your-registry/ptvs
  tag: "1.0.0"
  pullPolicy: IfNotPresent
```

### Resource Limits

Adjust resource requests and limits:

```yaml
resources:
  requests:
    cpu: 500m
    memory: 512Mi
  limits:
    cpu: 1500m
    memory: 1Gi
```

### Autoscaling

Configure HPA parameters:

```yaml
autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 15
  cpuUtilization: 65
```

## Security Features

### 1. Strict mTLS

All service-to-service communication is encrypted with mutual TLS:
- Enforced via `PeerAuthentication` in STRICT mode
- Certificate management handled by Istio

### 2. Pod Security Standards

The `mankind` namespace enforces the `restricted` Pod Security Standard:
- Containers run as non-root
- Privilege escalation is disabled
- All capabilities are dropped
- Seccomp profile is applied

### 3. Network Policies

NetworkPolicy restricts pod-to-pod communication:
- Only allows ingress from namespaces with Istio injection
- Egress rules can be further restricted as needed

### 4. Authorization Policies

Istio AuthorizationPolicy controls which services can communicate with PTVS:
- Currently allows traffic from `ocee` service account
- Modify `authorization-policy.yaml` to add more authorized services

### 5. OPA Gatekeeper

Enforces organizational policies:
- Requires `app` label on all pods
- Can be extended with additional constraints

## Image Signing with Cosign

To enable image verification with Cosign:

1. Sign your images in CI/CD:
```bash
cosign sign --key cosign.key your-registry/ptvs:1.0.0
```

2. Install Kyverno policy controller:
```bash
kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/main/config/install.yaml
```

3. Create a Kyverno ClusterPolicy to verify signed images:
```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: verify-ptvs-image
spec:
  validationFailureAction: enforce
  rules:
  - name: verify-signature
    match:
      resources:
        kinds:
        - Pod
    verifyImages:
    - imageReferences:
      - "your-registry/ptvs:*"
      attestors:
      - count: 1
        entries:
        - keys:
            publicKeys: |-
              -----BEGIN PUBLIC KEY-----
              <your-cosign-public-key>
              -----END PUBLIC KEY-----
```

## Canary Deployments

PTVS uses Argo Rollouts for progressive canary deployments:

1. **20% traffic** for 60 seconds
2. **50% traffic** for 120 seconds
3. **100% traffic** (full rollout)

To trigger a canary deployment:

```bash
# Update the image tag in values.yaml or via Helm
helm upgrade ptvs ptvs-platform/helm/ptvs \
  --set image.tag=1.1.0 \
  --namespace mankind
```

Monitor the rollout:
```bash
kubectl argo rollouts get rollout ptvs -n mankind --watch
```

Promote manually (if auto-promotion is disabled):
```bash
kubectl argo rollouts promote ptvs -n mankind
```

Abort a rollout:
```bash
kubectl argo rollouts abort ptvs -n mankind
```

## Observability

### Health Checks

PTVS exposes the following endpoints:

- `/healthz` - Liveness probe
- `/readyz` - Readiness probe
- `/metrics` - Prometheus metrics

### Prometheus Metrics

The deployment is annotated for Prometheus scraping:

```yaml
annotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "8080"
```

Metrics are automatically collected if Prometheus is configured with pod annotation discovery.

### Viewing Logs

View pod logs:
```bash
kubectl logs -n mankind -l app=ptvs --tail=100 -f
```

View logs with Istio sidecar:
```bash
kubectl logs -n mankind -l app=ptvs -c ptvs --tail=100 -f
kubectl logs -n mankind -l app=ptvs -c istio-proxy --tail=100 -f
```

### Istio Observability

View service mesh traffic:
```bash
istioctl dashboard kiali
```

View Grafana dashboards:
```bash
istioctl dashboard grafana
```

## High Availability

PTVS is configured for high availability:

1. **Minimum 3 replicas** - Ensures availability during rolling updates
2. **PodDisruptionBudget** - Maintains at least 2 pods during disruptions
3. **HPA** - Scales from 3 to 15 replicas based on CPU utilization
4. **Multi-zone deployment** - Use node affinity for zone distribution

Example zone distribution (add to deployment template):

```yaml
affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchExpressions:
          - key: app
            operator: In
            values:
            - ptvs
        topologyKey: topology.kubernetes.io/zone
```

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -n mankind -l app=ptvs
kubectl describe pod -n mankind -l app=ptvs
```

### Check Istio Sidecar Injection

```bash
kubectl get pod -n mankind -l app=ptvs -o jsonpath='{.items[*].spec.containers[*].name}'
```

Should show both `ptvs` and `istio-proxy`.

### Verify mTLS

```bash
istioctl authn tls-check -n mankind ptvs.mankind.svc.cluster.local
```

### Check NetworkPolicy

```bash
kubectl get networkpolicy -n mankind
kubectl describe networkpolicy ptvs -n mankind
```

### View Events

```bash
kubectl get events -n mankind --sort-by='.lastTimestamp'
```

### Debug Authorization

If requests are denied:

```bash
kubectl logs -n mankind -l app=ptvs -c istio-proxy | grep RBAC
```

## Cleanup

To remove the entire PTVS deployment:

```bash
# If deployed via ArgoCD
kubectl delete application ptvs -n argocd

# If deployed via Helm
helm uninstall ptvs -n mankind

# Remove policies and namespace
kubectl delete -f ptvs-platform/mesh/
kubectl delete -f ptvs-platform/policy/
kubectl delete -f ptvs-platform/bootstrap/
```

## Production Considerations

### 1. Image Registry

- Use a private container registry
- Implement image scanning in CI/CD
- Enable vulnerability scanning (e.g., Trivy, Clair)

### 2. Secrets Management

- Use Kubernetes Secrets or external secret managers (e.g., Vault, AWS Secrets Manager)
- Never commit secrets to Git
- Rotate secrets regularly

### 3. Backup and Disaster Recovery

- Backup Istio configuration
- Backup ArgoCD application definitions
- Document rollback procedures

### 4. Monitoring and Alerting

- Set up alerts for pod crashes
- Monitor HPA scaling events
- Alert on canary rollout failures
- Monitor certificate expiration

### 5. Cost Optimization

- Right-size resource requests and limits
- Use cluster autoscaler for node scaling
- Monitor resource utilization and adjust HPA settings

### 6. Multi-Cluster Deployment

For multi-cluster deployments:
- Use Istio multi-cluster mesh
- Configure cross-cluster service discovery
- Implement disaster recovery across regions

## Additional Resources

- [Istio Documentation](https://istio.io/latest/docs/)
- [Argo Rollouts Documentation](https://argoproj.github.io/argo-rollouts/)
- [ArgoCD Documentation](https://argo-cd.readthedocs.io/)
- [OPA Gatekeeper Documentation](https://open-policy-agent.github.io/gatekeeper/)
- [Kubernetes Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/)

## Support

For issues related to civic-attest deployment, please open an issue at:
https://github.com/IAmSoThirsty/civic-attest/issues

---

**Last Updated:** 2026-02-23
**Version:** 1.0.0
