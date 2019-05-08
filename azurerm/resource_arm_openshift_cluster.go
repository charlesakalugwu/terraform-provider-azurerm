package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-06-01/containerservice"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmOpenshiftCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmOpenshiftClusterCreateUpdate,
		Read:   resourceArmOpenshiftClusterRead,
		Update: resourceArmOpenshiftClusterCreateUpdate,
		Delete: resourceArmOpenshiftClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			return nil
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"purchase_plan": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"product": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"promotion_code": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"publisher": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"openshift_version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"cluster_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_profile": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vnet_cidr": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.CIDR,
						},

						"vnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"peer_vnet_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"router_profiles": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"public_subdomain": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"master_pool_profile": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "master",
							ValidateFunc: validate.KubernetesAgentPoolName,
						},

						"count": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},

						"vm_size": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validate.NoEmptyStrings, // add proper validation here to limit to the list of vm sizes
						},

						"subnet_cidr": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.CIDR,
						},

						"os_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  string(containerservice.Linux),
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.Linux),
								string(containerservice.Windows),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},

			"infra_pool_profile": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "infra",
							ValidateFunc: validate.KubernetesAgentPoolName,
						},

						"count": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},

						"vm_size": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validate.NoEmptyStrings, // add proper validation here to limit to the list of vm sizes
						},

						"subnet_cidr": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.CIDR,
						},

						"os_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  string(containerservice.Linux),
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.Linux),
								string(containerservice.Windows),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"role": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "infra",
							ValidateFunc: validate.KubernetesAgentPoolName,
						},
					},
				},
			},

			"compute_pool_profile": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "compute",
							ValidateFunc: validate.KubernetesAgentPoolName,
						},

						"count": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  4,
						},

						"vm_size": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validate.NoEmptyStrings, // add proper validation here to limit to the list of vm sizes
						},

						"subnet_cidr": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.CIDR,
						},

						"os_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  string(containerservice.Linux),
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.Linux),
								string(containerservice.Windows),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"role": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "compute",
							ValidateFunc: validate.KubernetesAgentPoolName,
						},
					},
				},
			},

			"azure_active_directory": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.UUID,
						},

						"client_secret": {
							Type:         schema.TypeString,
							ForceNew:     true,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"tenant_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validate.UUIDOrEmpty,
						},

						"customer_admin_group_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmOpenshiftClusterCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containers.OpenshiftClustersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure Red Hat OpenShift Cluster create/update.")

	rg := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, rg, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("error checking for presence of existing Azure Red Hat OpenShift Cluster %q (Resource Group %q): %s", name, rg, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_openshift_cluster", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	openshiftVersion := d.Get("openshift_version").(string)
	fqdn := d.Get("fqdn").(string)
	publicHostname := d.Get("public_hostname").(string)

	plan := expandOpenShiftClusterPurchasePlan(d)
	masterPoolProfile := expandOpenShiftClusterMasterPoolProfile(d)
	agentPoolProfiles := expandOpenShiftClusterAgentPoolProfiles(d)
	routerProfiles := expandOpenShiftClusterRouterProfiles(d)
	authProfile := expandOpenShiftClusterAuthrofile(d)
	networkProfile := expandOpenShiftClusterNetworkProfile(d)

	tags := d.Get("tags").(map[string]interface{})

	parameters := containerservice.OpenShiftManagedCluster{
		Name:     &name,
		Location: &location,
		Plan:     plan,
		OpenShiftManagedClusterProperties: &containerservice.OpenShiftManagedClusterProperties{
			OpenShiftVersion:  &openshiftVersion,
			Fqdn:              &fqdn,
			PublicHostname:    &publicHostname,
			MasterPoolProfile: masterPoolProfile,
			AgentPoolProfiles: &agentPoolProfiles,
			RouterProfiles:    &routerProfiles,
			AuthProfile:       authProfile,
			NetworkProfile:    networkProfile,
		},
		Tags: expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, rg, name, parameters)
	if err != nil {
		return fmt.Errorf("error creating/updating Azure Red Hat OpenShift Cluster %q (Resource Group %q): %+v", name, rg, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("error waiting for completion of Azure Red Hat OpenShift Cluster %q (Resource Group %q): %+v", name, rg, err)
	}

	read, err := client.Get(ctx, rg, name)
	if err != nil {
		return fmt.Errorf("error retrieving Azure Red Hat OpenShift Cluster %q (Resource Group %q): %+v", name, rg, err)
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read ID for Azure Red Hat OpenShift Cluster %q (Resource Group %q)", name, rg)
	}

	d.SetId(*read.ID)

	return resourceArmOpenshiftClusterRead(d, meta)
}

func resourceArmOpenshiftClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containers.OpenshiftClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["openShiftManagedClusters"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Azure Red Hat OpenShift Cluster %q was not found in Resource Group %q - removing from state!", name, resGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("error retrieving Azure Red Hat OpenShift Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err := d.Set("name", resp.Name); err != nil {
		return fmt.Errorf("error setting `name`: %+v", err)
	}

	if err := d.Set("resource_group_name", resGroup); err != nil {
		return fmt.Errorf("error setting `resource_group_name`: %+v", err)
	}

	if location := resp.Location; location != nil {
		if err := d.Set("location", azure.NormalizeLocation(*location)); err != nil {
			return fmt.Errorf("error setting `location`: %+v", err)
		}
	}

	if props := resp.OpenShiftManagedClusterProperties; props != nil {
		if err := d.Set("fqdn", props.Fqdn); err != nil {
			return fmt.Errorf("error setting `fqdn`: %+v", err)
		}
		if err := d.Set("openshift_version", props.OpenShiftVersion); err != nil {
			return fmt.Errorf("error setting `openshift_version`: %+v", err)
		}
		if err := d.Set("public_hostname", props.PublicHostname); err != nil {
			return fmt.Errorf("error setting `public_hostname`: %+v", err)
		}

		masterPoolProfile := flattenOpenShiftClusterMasterPoolProfile(props.MasterPoolProfile)
		if err := d.Set("master_pool_profile", masterPoolProfile); err != nil {
			return fmt.Errorf("error setting `master_pool_profile`: %+v", err)
		}

		agentPoolProfiles := flattenOpenShiftClusterAgentPoolProfiles(props.AgentPoolProfiles)
		if err := d.Set("agent_pool_profile", agentPoolProfiles); err != nil {
			return fmt.Errorf("error setting `agent_pool_profile`: %+v", err)
		}

		routerProfiles := flattenOpenShiftClusterRouterProfiles(props.RouterProfiles)
		if err := d.Set("router_profiles", routerProfiles); err != nil {
			return fmt.Errorf("error setting `router_profiles`: %+v", err)
		}

		authProfile := flattenOpenShiftClusterAuthProfile(props.AuthProfile)
		if err := d.Set("auth_profile", authProfile); err != nil {
			return fmt.Errorf("error setting `auth_profile`: %+v", err)
		}

		networkProfile := flattenOpenShiftNetworkProfile(props.NetworkProfile)
		if err := d.Set("network_profile", networkProfile); err != nil {
			return fmt.Errorf("error setting `network_profile`: %+v", err)
		}
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmOpenshiftClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containers.OpenshiftClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["openShiftManagedClusters"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("error deleting Azure Red Hat OpenShift Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("error waiting for the deletion of Azure Red Hat OpenShift Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func flattenOpenShiftNetworkProfile(profile *containerservice.NetworkProfile) interface{} {
	return nil
}

func flattenOpenShiftClusterAuthProfile(profile *containerservice.OpenShiftManagedClusterAuthProfile) []interface{} {
	return nil
}

func flattenOpenShiftClusterRouterProfiles(profile *[]containerservice.OpenShiftRouterProfile) []interface{} {
	return nil
}

func flattenOpenShiftClusterAgentPoolProfiles(profile *[]containerservice.OpenShiftManagedClusterAgentPoolProfile) []interface{} {
	return nil
}

func flattenOpenShiftClusterMasterPoolProfile(profile *containerservice.OpenShiftManagedClusterMasterPoolProfile) interface{} {
	return nil
}

func expandOpenShiftClusterNetworkProfile(data *schema.ResourceData) *containerservice.NetworkProfile {
	return nil
}

func expandOpenShiftClusterAuthrofile(data *schema.ResourceData) *containerservice.OpenShiftManagedClusterAuthProfile {
	return nil
}

func expandOpenShiftClusterRouterProfiles(data *schema.ResourceData) []containerservice.OpenShiftRouterProfile {
	return nil
}

func expandOpenShiftClusterAgentPoolProfiles(data *schema.ResourceData) []containerservice.OpenShiftManagedClusterAgentPoolProfile {
	return nil
}

func expandOpenShiftClusterMasterPoolProfile(data *schema.ResourceData) *containerservice.OpenShiftManagedClusterMasterPoolProfile {
	return nil
}

func expandOpenShiftClusterPurchasePlan(data *schema.ResourceData) *containerservice.PurchasePlan {
	return nil
}
