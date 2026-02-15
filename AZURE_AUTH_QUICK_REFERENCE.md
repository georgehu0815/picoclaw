# Azure OpenAI Authentication - Quick Reference Card

## 🎯 TL;DR - Start Testing Now

```bash
# Install SDKs
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity@latest
go get github.com/Azure/azure-sdk-for-go/sdk/azcore@latest

# Uncomment code in pkg/providers/codex_provider.go (lines 404-450)

# Login to Azure
az login

# Empty this in .env for local testing
# AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID=

# Test
go build -o picoclaw ./cmd/picoclaw/
./picoclaw agent -m "test"
```

---

## 📋 Authentication Methods

### Method 1: Local Testing (Azure CLI)

**When**: Development on your local machine
**Requires**: Azure CLI + `az login`
**Configuration**:
```bash
# .env file
AZURE_OPENAI_ENDPOINT=https://your-resource.cognitiveservices.azure.com/
AZURE_OPENAI_DEPLOYMENT=gpt-5.2-chat
AZURE_OPENAI_API_VERSION=2025-01-01-preview
AZURE_OPENAI_SCOPE=https://cognitiveservices.azure.com/.default
# AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID=  ← EMPTY or commented out
AZURE_OPENAI_VERBOSE=true
```

**How it works**: `DefaultAzureCredential()` → Azure CLI credentials

---

### Method 2: Azure with System-Assigned Managed Identity

**When**: Production deployment on Azure (App Service, VM, Container)
**Requires**: Managed identity enabled on Azure resource
**Configuration**:
```bash
# .env file (same as Method 1)
AZURE_OPENAI_ENDPOINT=https://your-resource.cognitiveservices.azure.com/
AZURE_OPENAI_DEPLOYMENT=gpt-5.2-chat
AZURE_OPENAI_API_VERSION=2025-01-01-preview
AZURE_OPENAI_SCOPE=https://cognitiveservices.azure.com/.default
# AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID=  ← EMPTY or commented out
```

**Azure Setup**:
```bash
# Enable system-assigned identity
az webapp identity assign --resource-group <rg> --name <app-name>

# Grant permissions
az role assignment create \
  --assignee <principal-id> \
  --role "Cognitive Services OpenAI User" \
  --scope <azure-openai-resource-id>
```

**How it works**: `DefaultAzureCredential()` → Managed Identity (system-assigned)

---

### Method 3: Azure with User-Assigned Managed Identity

**When**: Production with specific managed identity
**Requires**: User-assigned managed identity created and assigned
**Configuration**:
```bash
# .env file
AZURE_OPENAI_ENDPOINT=https://your-resource.cognitiveservices.azure.com/
AZURE_OPENAI_DEPLOYMENT=gpt-5.2-chat
AZURE_OPENAI_API_VERSION=2025-01-01-preview
AZURE_OPENAI_SCOPE=https://cognitiveservices.azure.com/.default
AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID=c9427d44-98e2-406a-9527-f7fa7059f984  ← SET
```

**Azure Setup**:
```bash
# Assign user-assigned identity to resource
az webapp identity assign \
  --resource-group <rg> \
  --name <app-name> \
  --identities <identity-resource-id>

# Grant permissions to the identity
az role assignment create \
  --assignee <identity-principal-id> \
  --role "Cognitive Services OpenAI User" \
  --scope <azure-openai-resource-id>
```

**How it works**: `ManagedIdentityCredential(clientID)` → User-assigned MI

---

## 🔄 Decision Tree

```
Need Azure OpenAI authentication?
│
├─ Testing locally?
│  └─ YES → Use Method 1 (Azure CLI)
│     • az login
│     • Leave MANAGED_IDENTITY_CLIENT_ID empty
│     • DefaultAzureCredential → Azure CLI
│
└─ Deploying to Azure?
   │
   ├─ Have specific managed identity requirement?
   │  │
   │  ├─ YES → Use Method 3 (User-Assigned MI)
   │  │  • Set MANAGED_IDENTITY_CLIENT_ID
   │  │  • Assign identity to resource
   │  │
   │  └─ NO → Use Method 2 (System-Assigned MI)
   │     • Leave MANAGED_IDENTITY_CLIENT_ID empty
   │     • Enable system identity on resource
```

---

## ✅ Checklist

### For Local Testing

- [ ] Azure CLI installed (`brew install azure-cli`)
- [ ] Logged in to Azure (`az login`)
- [ ] Azure SDK packages installed
- [ ] Azure code uncommented in `codex_provider.go`
- [ ] `MANAGED_IDENTITY_CLIENT_ID` is empty/commented in `.env`
- [ ] RBAC permissions granted to your user account
- [ ] Build successful (`go build`)
- [ ] Test passes (`./picoclaw agent -m "test"`)

### For Azure Deployment

- [ ] Azure SDK packages installed
- [ ] Azure code uncommented in `codex_provider.go`
- [ ] Binary built and deployed
- [ ] Managed identity enabled on Azure resource
- [ ] `MANAGED_IDENTITY_CLIENT_ID` configured (or empty for system-assigned)
- [ ] RBAC permissions granted to managed identity
- [ ] Environment variables set
- [ ] Test in Azure environment

---

## 🐛 Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| "failed to create Azure credential" | Not logged in to Azure CLI | Run `az login` |
| "authorization denied" | Missing RBAC permissions | Grant "Cognitive Services OpenAI User" role |
| "Azure authentication requires Azure SDK" | Code still commented | Uncomment implementation in `codex_provider.go` |
| "unrecognized import path" | SDK not installed | Run `go get` commands |
| Works locally, fails in Azure | Using Azure CLI auth | Set managed identity on Azure resource |
| Works in Azure, fails locally | No Azure CLI login | Run `az login` |

---

## 📊 Configuration Comparison

| Setting | Local Dev | Azure (System MI) | Azure (User MI) |
|---------|-----------|-------------------|-----------------|
| **AZURE_OPENAI_ENDPOINT** | ✅ Set | ✅ Set | ✅ Set |
| **AZURE_OPENAI_DEPLOYMENT** | ✅ Set | ✅ Set | ✅ Set |
| **AZURE_OPENAI_API_VERSION** | ✅ Set | ✅ Set | ✅ Set |
| **AZURE_OPENAI_SCOPE** | ✅ Set | ✅ Set | ✅ Set |
| **MANAGED_IDENTITY_CLIENT_ID** | ❌ Empty | ❌ Empty | ✅ Set |
| **Azure CLI Login** | ✅ Required | ❌ Not used | ❌ Not used |
| **Managed Identity Enabled** | ❌ Not needed | ✅ System | ✅ User |
| **Auth Method** | CLI | System MI | User MI |

---

## 🎓 Key Concepts

### DefaultAzureCredential
- Tries multiple auth methods in order
- Works locally (Azure CLI) AND in Azure (Managed Identity)
- No code changes needed between environments
- Recommended for most use cases

### ManagedIdentityCredential
- Uses specific managed identity by client ID
- Only works in Azure environment
- Use when you need a specific identity
- Requires MANAGED_IDENTITY_CLIENT_ID to be set

### Azure CLI Authentication
- Uses your personal Azure account
- Perfect for local development
- Requires `az login`
- Same permissions as your user account

### Managed Identity
- Azure resource has its own identity
- No credentials needed in code/config
- Scoped permissions via RBAC
- More secure than service principals

---

## 🚀 Quick Commands

```bash
# Check Azure login status
az account show

# List subscriptions
az account list --output table

# Check RBAC permissions
az role assignment list --assignee <principal-id> --output table

# Test token acquisition (local)
az account get-access-token --resource https://cognitiveservices.azure.com

# Enable verbose Azure auth logs
export AZURE_OPENAI_VERBOSE=true

# Build and test
go build -o picoclaw ./cmd/picoclaw/
./picoclaw agent -d -m "test"
```

---

## 📚 Documentation Links

- **Local Testing Guide**: [AZURE_LOCAL_TESTING_GUIDE.md](AZURE_LOCAL_TESTING_GUIDE.md)
- **Test Results**: [AZURE_MANAGED_IDENTITY_TEST_RESULTS.md](AZURE_MANAGED_IDENTITY_TEST_RESULTS.md)
- **Usage Guide**: [pkg/providers/CODEX_AZURE_USAGE.md](pkg/providers/CODEX_AZURE_USAGE.md)
- **Implementation**: [pkg/providers/CODEX_IMPLEMENTATION_SUMMARY.md](pkg/providers/CODEX_IMPLEMENTATION_SUMMARY.md)

---

**💡 Pro Tip**: Start with Method 1 (Azure CLI) for local testing, then deploy with Method 2 (System-Assigned MI) for production. The same code works in both environments!
