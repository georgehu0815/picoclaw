# Azure OpenAI - Local Testing Guide

**Test Azure OpenAI with Managed Identity on your local machine using Azure CLI**

## 🎯 Overview

With `DefaultAzureCredential()`, you can test Azure OpenAI authentication locally without deploying to Azure infrastructure. This works by using your Azure CLI credentials.

## ✅ Prerequisites

1. **Azure CLI installed**
   ```bash
   # macOS
   brew install azure-cli

   # Linux
   curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash

   # Windows
   # Download from: https://aka.ms/installazurecliwindows
   ```

2. **Azure subscription with Azure OpenAI access**

3. **RBAC permissions** - Your Azure account needs "Cognitive Services OpenAI User" role

## 🚀 Quick Start - Local Testing

### Step 1: Install Azure SDK Packages

```bash
cd /Users/ghu/aiworker/picoclaw

# Install required Azure packages
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity@latest
go get github.com/Azure/azure-sdk-for-go/sdk/azcore@latest

# Update dependencies
go mod tidy
```

### Step 2: Uncomment the Azure Authentication Code

Edit [pkg/providers/codex_provider.go](pkg/providers/codex_provider.go):

1. Find the `createAzureManagedIdentityTokenSource()` function (around line 388)
2. Remove the `/*` and `*/` comment markers around the implementation (lines 404-450)
3. The code is ready to use - no other changes needed!

### Step 3: Login with Azure CLI

```bash
# Login to Azure
az login

# Verify you're logged in
az account show

# List your subscriptions
az account list --output table

# Set active subscription (if you have multiple)
az account set --subscription "Your Subscription Name"
```

### Step 4: Grant RBAC Permissions (if not already set)

```bash
# Get your user principal ID
USER_PRINCIPAL=$(az ad signed-in-user show --query id -o tsv)

# Get your Azure OpenAI resource ID
RESOURCE_ID="/subscriptions/<subscription-id>/resourceGroups/<rg>/providers/Microsoft.CognitiveServices/accounts/datacopilothub8882317788"

# Grant "Cognitive Services OpenAI User" role
az role assignment create \
  --assignee $USER_PRINCIPAL \
  --role "Cognitive Services OpenAI User" \
  --scope $RESOURCE_ID

# Verify role assignment
az role assignment list \
  --assignee $USER_PRINCIPAL \
  --scope $RESOURCE_ID \
  --output table
```

### Step 5: Configure Environment Variables

Your `.env` file already has the configuration:

```bash
AZURE_OPENAI_ENDPOINT=https://datacopilothub8882317788.cognitiveservices.azure.com/
AZURE_OPENAI_DEPLOYMENT=gpt-5.2-chat
AZURE_OPENAI_API_VERSION=2025-01-01-preview
AZURE_OPENAI_SCOPE=https://cognitiveservices.azure.com/.default
# AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID=  # Leave empty for local testing
AZURE_OPENAI_VERBOSE=true
```

**Important**: For local testing, **do NOT set** `AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID`. Leave it empty or comment it out. This tells the code to use `DefaultAzureCredential()` which will use your Azure CLI credentials.

### Step 6: Build and Test

```bash
# Rebuild with Azure SDK
go build -o picoclaw ./cmd/picoclaw/

# Test with verbose logging
export AZURE_OPENAI_VERBOSE=true
./picoclaw agent -d -m "What is 2+2?"
```

**Expected Output:**
```
[CodexProvider] Using Azure OpenAI configuration
[AzureAuth] Using DefaultAzureCredential (supports local Azure CLI auth)
[AzureAuth] Retrieved token for scope: https://cognitiveservices.azure.com/.default
🦞 4
```

## 🔍 How DefaultAzureCredential Works

`DefaultAzureCredential()` tries multiple authentication methods in order:

1. **Environment Variables** - Service principal credentials
2. **Managed Identity** - When running in Azure (App Service, VM, Container Instance)
3. **Azure CLI** - Your local `az login` credentials ✅ (Used for local testing)
4. **Azure PowerShell** - PowerShell module credentials
5. **Interactive Browser** - Opens browser for authentication (if enabled)

For local development, it automatically uses your Azure CLI credentials (#3).

## 🧪 Testing Scenarios

### Scenario 1: Local Development (Your Case)

```bash
# Configuration
# AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID not set (empty)

# Auth flow
DefaultAzureCredential → Azure CLI → Your user account → ✅ Token

# Test
./picoclaw agent -d -m "Hello from local dev"
```

### Scenario 2: Azure Deployment with System-Assigned Managed Identity

```bash
# Configuration
# AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID not set (empty)
# Running on Azure App Service with managed identity enabled

# Auth flow
DefaultAzureCredential → Managed Identity → System identity → ✅ Token

# Same binary works without code changes!
```

### Scenario 3: Azure Deployment with User-Assigned Managed Identity

```bash
# Configuration
AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID=c9427d44-98e2-406a-9527-f7fa7059f984

# Auth flow
ManagedIdentityCredential → User-assigned identity → ✅ Token
```

## 📋 Troubleshooting

### Error: "failed to create Azure credential"

**Solution**: Ensure Azure CLI is logged in
```bash
az login
az account show
```

### Error: "failed to get Azure access token: authorization denied"

**Solution**: Check RBAC permissions
```bash
# Verify your role assignment
az role assignment list \
  --assignee $(az ad signed-in-user show --query id -o tsv) \
  --all \
  --output table | grep "Cognitive"
```

### Error: "Azure authentication requires Azure SDK packages"

**Solution**: Install and rebuild
```bash
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity@latest
go get github.com/Azure/azure-sdk-for-go/sdk/azcore@latest
go mod tidy
go build -o picoclaw ./cmd/picoclaw/
```

### Error: "package github.com/Azure/azure-sdk-for-go/sdk/azidentity: unrecognized import path"

**Solution**: Uncomment the import statements in the code
```go
// In createAzureManagedIdentityTokenSource(), uncomment:
import (
    "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)
```

### Verbose Logging Shows Nothing

**Solution**: Enable verbose mode
```bash
export AZURE_OPENAI_VERBOSE=true
./picoclaw agent -d -m "Test"
```

## 🎯 Configuration Matrix

| Environment | MANAGED_IDENTITY_CLIENT_ID | Auth Method Used | Works? |
|-------------|---------------------------|------------------|--------|
| Local (with az login) | Empty/Not set | DefaultAzureCredential → Azure CLI | ✅ Yes |
| Local (no az login) | Empty/Not set | DefaultAzureCredential → Error | ❌ No |
| Azure App Service | Empty/Not set | DefaultAzureCredential → System MI | ✅ Yes |
| Azure VM | Empty/Not set | DefaultAzureCredential → System MI | ✅ Yes |
| Azure Container | Set to client ID | ManagedIdentityCredential → User MI | ✅ Yes |
| Azure App Service | Set to client ID | ManagedIdentityCredential → User MI | ✅ Yes |

## 💡 Best Practices

### Local Development

```bash
# 1. Use DefaultAzureCredential (leave MANAGED_IDENTITY_CLIENT_ID empty)
# 2. Stay logged in to Azure CLI
az login

# 3. Enable verbose logging
export AZURE_OPENAI_VERBOSE=true

# 4. Test regularly
./picoclaw agent -d -m "Test"
```

### CI/CD Pipeline

```bash
# Use service principal authentication
export AZURE_CLIENT_ID=<client-id>
export AZURE_CLIENT_SECRET=<client-secret>
export AZURE_TENANT_ID=<tenant-id>

# DefaultAzureCredential will use these automatically
```

### Production Deployment

```bash
# Enable managed identity on Azure resource
az webapp identity assign --resource-group <rg> --name <app-name>

# Don't set AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID
# Let DefaultAzureCredential auto-detect the system identity

# Or use user-assigned identity:
AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID=c9427d44-98e2-406a-9527-f7fa7059f984
```

## 🔐 Security Notes

1. **Local Development**: Your personal Azure account is used - ensure it has minimal necessary permissions
2. **Production**: Use managed identity - never embed credentials in code or config files
3. **CI/CD**: Use service principals with limited scope and time-bound credentials
4. **Secrets**: Never commit Azure credentials to version control

## 📊 Complete Testing Workflow

```bash
# 1. Install dependencies
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity@latest
go get github.com/Azure/azure-sdk-for-go/sdk/azcore@latest
go mod tidy

# 2. Uncomment Azure implementation in codex_provider.go
# (Remove /* and */ around lines 404-450)

# 3. Configure environment
cat > .env <<EOF
AZURE_OPENAI_ENDPOINT=https://datacopilothub8882317788.cognitiveservices.azure.com/
AZURE_OPENAI_DEPLOYMENT=gpt-5.2-chat
AZURE_OPENAI_API_VERSION=2025-01-01-preview
AZURE_OPENAI_SCOPE=https://cognitiveservices.azure.com/.default
AZURE_OPENAI_VERBOSE=true
# AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID=  # Empty for local testing
EOF

# 4. Login to Azure
az login

# 5. Grant permissions (if needed)
az role assignment create \
  --assignee $(az ad signed-in-user show --query id -o tsv) \
  --role "Cognitive Services OpenAI User" \
  --scope /subscriptions/<sub-id>/resourceGroups/<rg>/providers/Microsoft.CognitiveServices/accounts/datacopilothub8882317788

# 6. Build
go build -o picoclaw ./cmd/picoclaw/

# 7. Test
export AZURE_OPENAI_VERBOSE=true
./picoclaw agent -d -m "What is 2+2?"

# Expected output:
# [CodexProvider] Using Azure OpenAI configuration
# [AzureAuth] Using DefaultAzureCredential (supports local Azure CLI auth)
# [AzureAuth] Retrieved token for scope: https://cognitiveservices.azure.com/.default
# 🦞 4
```

## 🎉 Success!

Once you see the response, you're successfully using Azure OpenAI with your local Azure CLI credentials!

The same code will work when deployed to Azure, automatically switching from Azure CLI auth to Managed Identity auth - no code changes needed.

---

**Key Takeaway**: By using `DefaultAzureCredential()`, you get seamless authentication that works locally (Azure CLI) and in production (Managed Identity) with the same code.

🔗 **Next Steps**:
- [Azure Managed Identity Test Results](AZURE_MANAGED_IDENTITY_TEST_RESULTS.md)
- [Codex Azure Usage Guide](pkg/providers/CODEX_AZURE_USAGE.md)
