package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	edgeos "github.com/montanadev/edgeos-config-api"
)

func resourceFirewallDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*edgeos.Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("failed to coerce provider client"))
	}
	return diag.FromErr(client.DeleteFirewallRule(&edgeos.DeleteFirewallRuleRequest{
		FirewallName: d.Get("firewall_name").(string),
		RuleId:       d.Id(),
	}))
}

func resourceFirewallUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceFirewallCreate(ctx, d, meta)
}

func resourceFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO - implement read
	return nil
}

func resourceFirewallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*edgeos.Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("failed to coerce provider client"))
	}

	var states []string
	stateInterfaces := d.Get("states").([]interface{})
	for _, s := range stateInterfaces {
		states = append(states, s.(string))
	}

	req := &edgeos.CreateOrUpdateFirewallRuleRequest{
		FirewallName:  d.Get("firewall_name").(string),
		RuleName:      d.Get("name").(string),
		RuleId:        d.Get("rule_id").(string),
		Action:        d.Get("action").(string),
		Protocol:      d.Get("protocol").(string),
		States:        states,
		EnableLogging: d.Get("enable_logging").(bool),
	}

	if d.Get("destination_address").(string) != "" {
		req.DestinationAddress = edgeos.String(d.Get("destination_address").(string))
	}
	if d.Get("destination_port").(int) != 0 {
		req.DestinationPort = edgeos.Int(d.Get("destination_port").(int))
	}
	if d.Get("destination_mac").(string) != "" {
		req.DestinationMACAddress = edgeos.String(d.Get("destination_mac").(string))
	}
	if d.Get("source_address").(string) != "" {
		req.SourceAddress = edgeos.String(d.Get("source_address").(string))
	}
	if d.Get("source_port").(int) != 0 {
		req.SourcePort = edgeos.Int(d.Get("source_port").(int))
	}
	if d.Get("source_mac").(string) != "" {
		req.SourceMACAddress = edgeos.String(d.Get("source_mac").(string))
	}

	if err := client.CreateOrUpdateFirewallRule(req); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(req.RuleId)
	d.Set("firewall_name", req.FirewallName)
	d.Set("name", req.RuleName)
	d.Set("rule_id", req.RuleId)
	d.Set("action", req.Action)
	d.Set("protocol", req.Protocol)
	d.Set("states", stringSliceToSet(states))
	d.Set("enable_logging", req.EnableLogging)
	d.Set("destination_address", req.DestinationAddress)
	d.Set("destination_port", req.DestinationPort)
	d.Set("destination_mac", req.DestinationMACAddress)
	d.Set("source_address", req.SourceAddress)
	d.Set("source_port", req.SourcePort)
	d.Set("source_mac", req.SourceMACAddress)

	return nil
}

func dataSourceFirewalls() *schema.Resource {
	return &schema.Resource{
		Description:   "asdf",
		CreateContext: resourceFirewallCreate,
		ReadContext:   resourceFirewallRead,
		DeleteContext: resourceFirewallDelete,
		UpdateContext: resourceFirewallUpdate,
		Schema: map[string]*schema.Schema{
			"firewall_name": {
				Description: "firewall_name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"id": {
				Description: "id",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"rule_id": {
				Description: "rule_id",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"action": {
				Description: "action",
				Type:        schema.TypeString,
				Required:    true,
			},
			"protocol": {
				Description: "protocol",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"states": {
				Description: "states",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"source_address": {
				Description: "source_address",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"source_port": {
				Description: "source_port",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"source_mac": {
				Description: "source_mac",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"destination_address": {
				Description: "destination_address",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"destination_port": {
				Description: "destination_port",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"destination_mac": {
				Description: "destination_mac",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enable_logging": {
				Description: "enable_logging",
				Type:        schema.TypeBool,
				Optional:    true,
			},
		},
	}
}
