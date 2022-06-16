package main

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights"
	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/v1/operationalinsights"
)

var (
	subscriptionID    string
	resourceGroupName string
	workspaceName     string
)

func main() {
	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		log.Fatal("AZURE_SUBSCRIPTION_ID is not set.")
	}
	resourceGroupName = os.Getenv("AZURE_RESOURCE_GROUP")
	workspaceName = os.Getenv("AZURE_WORKSPACE_NAME")

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	workspace, err := getWorkspace(ctx, cred)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("operational insights get workspace:", *workspace.ID)

	kqlQuery := `
	AzureDiagnostics
		| where Category == "kube-audit-admin"
		| extend logs = parse_json(log_s)
		| where logs.verb == "delete"
		| project TimeGenerated, logs.kind, logs.level, logs.verb, logs
	`

	queryBody := operationalinsights.QueryBody{
		Query: &kqlQuery,
	}
	queryClient := operationalinsights.NewQueryClient()
	queryClient.Execute(ctx, *workspace.ID, queryBody)

	workspaces, err := listWorkspace(ctx, cred)
	if err != nil {
		log.Fatal(err)
	}
	for _, w := range workspaces {
		log.Println(*w.Name, *w.ID)
	}
}

func getWorkspace(ctx context.Context, cred azcore.TokenCredential) (*armoperationalinsights.Workspace, error) {
	workspacesClient, err := armoperationalinsights.NewWorkspacesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	resp, err := workspacesClient.Get(ctx, resourceGroupName, workspaceName, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Workspace, nil
}

func listWorkspace(ctx context.Context, cred azcore.TokenCredential) ([]*armoperationalinsights.Workspace, error) {
	workspacesClient, err := armoperationalinsights.NewWorkspacesClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, err
	}

	workspaceResp := workspacesClient.NewListByResourceGroupPager(resourceGroupName, nil)
	pager, err := workspaceResp.NextPage(ctx)
	if err != nil {
		return nil, err
	}
	workspaces := make([]*armoperationalinsights.Workspace, 0)
	workspaces = append(workspaces, pager.Value...)
	return workspaces, nil
}
