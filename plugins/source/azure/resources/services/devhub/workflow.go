package devhub

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devhub/armdevhub"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func Workflow() *schema.Table {
	return &schema.Table{
		Name:        "azure_devhub_workflow",
		Resolver:    fetchWorkflow,
		Description: "https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devhub/armdevhub@v0.2.0#Workflow",
		Multiplex:   client.SubscriptionMultiplexRegisteredNamespace("azure_devhub_workflow", client.Namespacemicrosoft_devhub),
		Transform:   transformers.TransformWithStruct(&armdevhub.Workflow{}),
		Columns: []schema.Column{
			client.SubscriptionID,
			{
				Name:     "id",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("ID"),
				CreationOptions: schema.ColumnCreationOptions{
					PrimaryKey: true,
				},
			},
		},
	}
}

func fetchWorkflow(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	svc, err := armdevhub.NewWorkflowClient(cl.SubscriptionId, cl.Creds, cl.Options)
	if err != nil {
		return err
	}
	pager := svc.NewListPager(nil)
	for pager.More() {
		p, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		res <- p.Value
	}
	return nil
}
