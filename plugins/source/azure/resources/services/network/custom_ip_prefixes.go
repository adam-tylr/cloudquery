package network

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func CustomIpPrefixes() *schema.Table {
	return &schema.Table{
		Name:        "azure_network_custom_ip_prefixes",
		Resolver:    fetchCustomIpPrefixes,
		Description: "https://learn.microsoft.com/en-us/rest/api/virtualnetwork/custom-ip-prefixes/list?tabs=HTTP#customipprefix",
		Multiplex:   client.SubscriptionMultiplexRegisteredNamespace("azure_network_custom_ip_prefixes", client.Namespacemicrosoft_network),
		Transform:   transformers.TransformWithStruct(&armnetwork.CustomIPPrefix{}),
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

func fetchCustomIpPrefixes(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	svc, err := armnetwork.NewCustomIPPrefixesClient(cl.SubscriptionId, cl.Creds, cl.Options)
	if err != nil {
		return err
	}
	pager := svc.NewListAllPager(nil)
	for pager.More() {
		p, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		res <- p.Value
	}
	return nil
}
