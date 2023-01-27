package security

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
)

func SecureScoreControls() *schema.Table {
	return &schema.Table{
		Name:        "azure_security_secure_score_controls",
		Resolver:    fetchSecureScoreControls,
		Description: "https://learn.microsoft.com/en-us/rest/api/defenderforcloud/secure-score-controls/list?tabs=HTTP#securescorecontroldetails",
		Multiplex:   client.SubscriptionMultiplexRegisteredNamespace("azure_security_secure_score_controls", client.Namespacemicrosoft_security),
		Transform:   transformers.TransformWithStruct(&armsecurity.SecureScoreControlDetails{}),
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

func fetchSecureScoreControls(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	svc, err := armsecurity.NewSecureScoreControlsClient(cl.SubscriptionId, cl.Creds, cl.Options)
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
