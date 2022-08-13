package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	edgeos "github.com/montanadev/edgeos-config-api"
)

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	host := d.Get("host").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (username != "") && (password != "") && (host != "") {
		c, err := edgeos.NewClient(host, username, password)
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("failed to create new edgerouter client: %s", err))
		}

		if err := c.Login(); err != nil {
			return nil, diag.FromErr(fmt.Errorf("failed to login to edgerouter: %s", err))
		}

		return c, diags
	}

	return nil, diags
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Description: "username",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("EDGEROUTER_USERNAME", nil),
			},
			"password": {
				Description: "password",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("EDGEROUTER_PASSWORD", nil),
			},
			"host": {
				Description: "host",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("EDGEROUTER_HOST", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"edgeos_firewall_rule": dataSourceFirewalls(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}
