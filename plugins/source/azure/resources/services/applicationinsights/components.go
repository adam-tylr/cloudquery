package applicationinsights

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func Components() *schema.Table {
	return &schema.Table{
		Name:        "azure_applicationinsights_components",
		Resolver:    fetchComponents,
		Description: "https://learn.microsoft.com/en-us/rest/api/application-insights/components/list?tabs=HTTP#applicationinsightscomponent",
		Multiplex:   client.SubscriptionMultiplexRegisteredNamespace("azure_applicationinsights_components", client.Namespacemicrosoft_insights),
		Transform:   transformers.TransformWithStruct(&armapplicationinsights.Component{}),
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

func fetchComponents(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	svc, err := armapplicationinsights.NewComponentsClient(cl.SubscriptionId, cl.Creds, cl.Options)
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
