package appstream

import (
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appstream"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAppstreamFleet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppstreamFleetCreate,
		Read:   resourceAppstreamFleetRead,
		Update: resourceAppstreamFleetUpdate,
		Delete: resourceAppstreamFleetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"compute_capacity": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"desired_instances": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"disconnect_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"domain_info": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"directory_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"organizational_unit_distinguished_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"enable_default_internet_access": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"fleet_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"iam_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"idle_disconnect_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			// TO-BE-REPAIRED: there is some inconsistency between sdk and api that makes it messy right now
			//"image_arn": {
			//	Type:          schema.TypeString,
			//	Optional:      true,
			//	ConflictsWith: []string{"image_name"},
			//},

			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
				//ConflictsWith: []string{"image_arn"},
			},

			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"max_user_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"stream_view": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "APP",
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"vpc_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"subnet_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceAppstreamFleetCreate(d *schema.ResourceData, meta interface{}) error {
	svc := meta.(*AWSClient).appstreamconn
	CreateFleetInputOpts := &appstream.CreateFleetInput{}

	ComputeConfig := &appstream.ComputeCapacity{}
	if a, ok := d.GetOk("compute_capacity"); ok {
		ComputeAttributes := a.([]interface{})
		attr := ComputeAttributes[0].(map[string]interface{})

		if v, ok := attr["desired_instances"]; ok {
			ComputeConfig.DesiredInstances = aws.Int64(int64(v.(int)))
		}

		CreateFleetInputOpts.ComputeCapacity = ComputeConfig
	}

	if v, ok := d.GetOk("description"); ok {
		CreateFleetInputOpts.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("disconnect_timeout"); ok {
		CreateFleetInputOpts.DisconnectTimeoutInSeconds = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("display_name"); ok {
		CreateFleetInputOpts.DisplayName = aws.String(v.(string))
	}

	DomainJoinInfoConfig := &appstream.DomainJoinInfo{}
	if dom, ok := d.GetOk("domain_info"); ok {
		DomainAttributes := dom.([]interface{})
		attr := DomainAttributes[0].(map[string]interface{})

		if v, ok := attr["directory_name"]; ok {
			DomainJoinInfoConfig.DirectoryName = aws.String(v.(string))
		}

		if v, ok := attr["organizational_unit_distinguished_name"]; ok {
			DomainJoinInfoConfig.OrganizationalUnitDistinguishedName = aws.String(v.(string))
		}

		CreateFleetInputOpts.DomainJoinInfo = DomainJoinInfoConfig
	}

	if v, ok := d.GetOk("enable_default_internet_access"); ok {
		CreateFleetInputOpts.EnableDefaultInternetAccess = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("fleet_type"); ok {
		CreateFleetInputOpts.FleetType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("iam_role_arn"); ok {
		CreateFleetInputOpts.IamRoleArn = aws.String(v.(string))
	}

	if v, ok := d.GetOk("idle_disconnect_timeout"); ok {
		CreateFleetInputOpts.IdleDisconnectTimeoutInSeconds = aws.Int64(int64(v.(int)))
	}

	//if v, ok := d.GetOk("image_arn"); ok {
	//	CreateFleetInputOpts.ImageArn = aws.String(v.(string))
	//}

	if v, ok := d.GetOk("image_name"); ok {
		CreateFleetInputOpts.ImageName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("instance_type"); ok {
		CreateFleetInputOpts.InstanceType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("max_user_duration"); ok {
		CreateFleetInputOpts.MaxUserDurationInSeconds = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("name"); ok {
		CreateFleetInputOpts.Name = aws.String(v.(string))
	}

	if v, ok := d.GetOk("stream_view"); ok {
		CreateFleetInputOpts.StreamView = aws.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_config"); ok {
		CreateFleetInputOpts.VpcConfig = expandVpcConfigs(v.([]interface{}))
	}

	log.Printf("[DEBUG] Run configuration: %s", CreateFleetInputOpts)

	resp, err := svc.CreateFleet(CreateFleetInputOpts)
	if err != nil {
		log.Printf("[ERROR] Error creating Appstream Fleet: %s", err)
		return err
	}

	log.Printf("[DEBUG] Appstream Fleet created %s ", resp)

	if v, ok := d.GetOk("tags"); ok {
		time.Sleep(2 * time.Second)

		fleet_name := aws.StringValue(CreateFleetInputOpts.Name)
		get, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{
			Names: aws.StringSlice([]string{fleet_name}),
		})

		if err != nil {
			log.Printf("[ERROR] Error describing Appstream Fleet: %s", err)
			return err
		}

		if get.Fleets == nil {
			log.Printf("[DEBUG] Appstream Fleet (%s) not found", d.Id())
		}

		tag, err := svc.TagResource(&appstream.TagResourceInput{
			ResourceArn: get.Fleets[0].Arn,
			Tags:        aws.StringMap(expandTags(v.(map[string]interface{}))),
		})

		if err != nil {
			log.Printf("[ERROR] Error tagging Appstream Fleet: %s", err)
			return err
		}

		log.Printf("[DEBUG] %s", tag)
	}

	if v, ok := d.GetOk("state"); ok {
		if v == "RUNNING" {
			resp, err := svc.StartFleet(&appstream.StartFleetInput{
				Name: CreateFleetInputOpts.Name,
			})

			if err != nil {
				log.Printf("[ERROR] Error satrting Appstream Fleet: %s", err)
				return err
			}

			log.Printf("[DEBUG] %s", resp)

			for {
				resp, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{
					Names: aws.StringSlice([]string{*CreateFleetInputOpts.Name}),
				})

				if err != nil {
					log.Printf("[ERROR] Error describing Appstream Fleet: %s", err)
					return err
				}

				curr_state := resp.Fleets[0].State
				if aws.StringValue(curr_state) == v {
					break
				}

				if aws.StringValue(curr_state) != v {
					time.Sleep(20 * time.Second)
					continue
				}
			}
		}
	}

	d.SetId(*CreateFleetInputOpts.Name)

	return resourceAppstreamFleetRead(d, meta)
}

func resourceAppstreamFleetRead(d *schema.ResourceData, meta interface{}) error {
	svc := meta.(*AWSClient).appstreamconn

	resp, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{})
	if err != nil {
		log.Printf("[ERROR] Error reading Appstream Fleet: %s", err)
		return err
	}

	for _, v := range resp.Fleets {
		if aws.StringValue(v.Name) == d.Get("name") {
			if v.ComputeCapacityStatus != nil {
				comp_attr := map[string]interface{}{}
				comp_attr["desired_instances"] = aws.Int64Value(v.ComputeCapacityStatus.Desired)
				d.Set("compute_capacity", comp_attr)
			}

			d.Set("description", v.Description)
			d.Set("disconnect_timeout", v.DisconnectTimeoutInSeconds)
			d.Set("display_name", v.DisplayName)

			if v.DomainJoinInfo != nil {
				dom_attr := map[string]interface{}{}
				dom_attr["directory_name"] = v.DomainJoinInfo.DirectoryName
				dom_attr["organizational_unit_distinguished_name"] = v.DomainJoinInfo.OrganizationalUnitDistinguishedName
				d.Set("domain_info", dom_attr)
			}

			d.Set("enable_default_internet_access", v.EnableDefaultInternetAccess)
			d.Set("fleet_type", v.FleetType)
			d.Set("iam_role_arn", v.IamRoleArn)
			d.Set("idle_disconnect_timeout", v.IdleDisconnectTimeoutInSeconds)
			//d.Set("image_arn", v.ImageArn)
			d.Set("image_name", v.ImageName)
			d.Set("instance_type", v.InstanceType)
			d.Set("max_user_duration", v.MaxUserDurationInSeconds)
			d.Set("name", v.Name)
			d.Set("stream_view", v.StreamView)

			tg, err := svc.ListTagsForResource(&appstream.ListTagsForResourceInput{
				ResourceArn: v.Arn,
			})

			if err != nil {
				log.Printf("[ERROR] Error listing stack tags: %s", err)
				return err
			}

			if tg.Tags == nil {
				log.Printf("[DEBUG] Apsstream Stack tags (%s) not found", d.Id())
				return nil
			}

			if len(tg.Tags) > 0 {
				tags_attr := make(map[string]string)
				tags := tg.Tags
				for k, v := range tags {
					tags_attr[k] = aws.StringValue(v)
				}

				d.Set("tags", tags_attr)
			}

			if v.VpcConfig != nil {
				vpc_attr := map[string]interface{}{}
				vpc_config_sg := aws.StringValueSlice(v.VpcConfig.SecurityGroupIds)
				vpc_config_sub := aws.StringValueSlice(v.VpcConfig.SubnetIds)
				vpc_attr["security_group_ids"] = aws.String(strings.Join(vpc_config_sg, ","))
				vpc_attr["subnet_ids"] = aws.String(strings.Join(vpc_config_sub, ","))
				d.Set("vpc_config", vpc_attr)
			}

			d.Set("state", v.State)

			return nil
		}
	}

	d.SetId("")

	return nil
}

func resourceAppstreamFleetUpdate(d *schema.ResourceData, meta interface{}) error {
	svc := meta.(*AWSClient).appstreamconn
	UpdateFleetInputOpts := &appstream.UpdateFleetInput{}

	d.Partial(true)

	if d.HasChange("description") {
		d.SetPartial("description")
		log.Printf("[DEBUG] Modify Fleet")
		description := d.Get("description").(string)
		UpdateFleetInputOpts.Description = aws.String(description)
	}

	if d.HasChange("disconnect_timeout") {
		d.SetPartial("disconnect_timeout")
		log.Printf("[DEBUG] Modify Fleet")
		disconnect_timeout := d.Get("disconnect_timeout").(int)
		UpdateFleetInputOpts.DisconnectTimeoutInSeconds = aws.Int64(int64(disconnect_timeout))
	}

	if d.HasChange("display_name") {
		d.SetPartial("display_name")
		log.Printf("[DEBUG] Modify Fleet")
		display_name := d.Get("display_name").(string)
		UpdateFleetInputOpts.DisplayName = aws.String(display_name)
	}

	if d.HasChange("enable_default_internet_access") {
		d.SetPartial("enable_default_internet_access")
		log.Printf("[DEBUG] Modify Fleet")
		enable_default_internet_access := d.Get("enable_default_internet_access").(bool)
		UpdateFleetInputOpts.EnableDefaultInternetAccess = aws.Bool(enable_default_internet_access)
	}

	if d.HasChange("iam_role_arn") {
		d.SetPartial("iam_role_arn")
		log.Printf("[DEBUG] Modify Fleet")
		iam_role_arn := d.Get("iam_role_arn").(string)
		UpdateFleetInputOpts.IamRoleArn = aws.String(iam_role_arn)
	}

	if d.HasChange("idle_disconnect_timeout") {
		d.SetPartial("idle_disconnect_timeout")
		log.Printf("[DEBUG] Modify Fleet")
		idle_disconnect_timeout_in_seconds := d.Get("idle_disconnect_timeout").(int)
		UpdateFleetInputOpts.IdleDisconnectTimeoutInSeconds = aws.Int64(int64(idle_disconnect_timeout_in_seconds))
	}

	// TO-BE-REPAIRED: there is some inconsistency between sdk and api that makes it hard right now
	//if d.HasChange("image_arn") {
	//	d.SetPartial("image_arn")
	//	log.Printf("[DEBUG] Modify Fleet")
	//	image_arn := d.Get("image_arn").(string)
	//	UpdateFleetInputOpts.ImageArn = aws.String(image_arn)
	//}

	if d.HasChange("image_name") {
		d.SetPartial("image_name")
		log.Printf("[DEBUG] Modify Fleet")
		image_name := d.Get("image_name").(string)
		UpdateFleetInputOpts.ImageName = aws.String(image_name)
	}

	if d.HasChange("instance_type") {
		d.SetPartial("instance_type")
		log.Printf("[DEBUG] Modify Fleet")
		instance_type := d.Get("instance_type").(string)
		UpdateFleetInputOpts.InstanceType = aws.String(instance_type)
	}

	if d.HasChange("max_user_duration") {
		d.SetPartial("max_user_duration")
		log.Printf("[DEBUG] Modify Fleet")
		max_user_duration := d.Get("max_user_duration").(int)
		UpdateFleetInputOpts.MaxUserDurationInSeconds = aws.Int64(int64(max_user_duration))
	}

	if v, ok := d.GetOk("name"); ok {
		UpdateFleetInputOpts.Name = aws.String(v.(string))
	}

	if d.HasChange("stream_view") {
		d.SetPartial("stream_view")
		log.Printf("[DEBUG] Modify Fleet")
		stream_view := d.Get("stream_view").(string)
		UpdateFleetInputOpts.StreamView = aws.String(stream_view)
	}

	resp, err := svc.UpdateFleet(UpdateFleetInputOpts)
	if err != nil {
		log.Printf("[ERROR] Error updating Appstream Fleet: %s", err)
		return err
	}

	log.Printf("[DEBUG] Appstream Fleet updated %s ", resp)

	if v, ok := d.GetOk("tags"); ok && d.HasChange("tags") {
		time.Sleep(2 * time.Second)

		fleet_name := aws.StringValue(UpdateFleetInputOpts.Name)
		get, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{
			Names: aws.StringSlice([]string{fleet_name}),
		})

		if err != nil {
			log.Printf("[ERROR] Error describing Appstream Fleet: %s", err)
			return err
		}

		if get.Fleets == nil {
			log.Printf("[DEBUG] Appstream Fleet (%s) not found", d.Id())
		}

		tag, err := svc.TagResource(&appstream.TagResourceInput{
			ResourceArn: get.Fleets[0].Arn,
			Tags:        aws.StringMap(expandTags(v.(map[string]interface{}))),
		})

		if err != nil {
			log.Printf("[ERROR] Error tagging Appstream Fleet: %s", err)
			return err
		}

		log.Printf("[DEBUG] %s", tag)
	}

	desired_state := d.Get("state")
	if d.HasChange("state") {
		d.SetPartial("state")

		if desired_state == "STOPPED" {
			svc.StopFleet(&appstream.StopFleetInput{
				Name: aws.String(d.Id()),
			})

			for {
				resp, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{
					Names: aws.StringSlice([]string{*UpdateFleetInputOpts.Name}),
				})

				if err != nil {
					log.Printf("[ERROR] Error describing Appstream Fleet: %s", err)
					return err
				}

				curr_state := resp.Fleets[0].State
				if aws.StringValue(curr_state) == desired_state {
					break
				}

				if aws.StringValue(curr_state) != desired_state {
					time.Sleep(20 * time.Second)
					continue
				}
			}
		} else if desired_state == "RUNNING" {
			svc.StartFleet(&appstream.StartFleetInput{
				Name: aws.String(d.Id()),
			})

			for {
				resp, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{
					Names: aws.StringSlice([]string{*UpdateFleetInputOpts.Name}),
				})

				if err != nil {
					log.Printf("[ERROR] Error describing Appstream Fleet: %s", err)
					return err
				}

				curr_state := resp.Fleets[0].State
				if aws.StringValue(curr_state) == desired_state {
					break
				}

				if aws.StringValue(curr_state) != desired_state {
					time.Sleep(20 * time.Second)
					continue
				}
			}
		}
	}

	d.Partial(false)

	return resourceAppstreamFleetRead(d, meta)
}

func resourceAppstreamFleetDelete(d *schema.ResourceData, meta interface{}) error {
	svc := meta.(*AWSClient).appstreamconn

	resp, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{
		Names: aws.StringSlice([]string{*aws.String(d.Id())}),
	})

	if err != nil {
		log.Printf("[ERROR] Error reading Appstream Fleet: %s", err)
		return err
	}

	curr_state := aws.StringValue(resp.Fleets[0].State)

	if curr_state == "RUNNING" {
		desired_state := "STOPPED"
		svc.StopFleet(&appstream.StopFleetInput{
			Name: aws.String(d.Id()),
		})
		for {

			resp, err := svc.DescribeFleets(&appstream.DescribeFleetsInput{
				Names: aws.StringSlice([]string{*aws.String(d.Id())}),
			})
			if err != nil {
				log.Printf("[ERROR] Error describing Appstream Fleet: %s", err)
				return err
			}

			curr_state := resp.Fleets[0].State
			if aws.StringValue(curr_state) == desired_state {
				break
			}
			if aws.StringValue(curr_state) != desired_state {
				time.Sleep(20 * time.Second)
				continue
			}
		}
	}

	del, err := svc.DeleteFleet(&appstream.DeleteFleetInput{
		Name: aws.String(d.Id()),
	})
	if err != nil {
		log.Printf("[ERROR] Error deleting Appstream Fleet: %s", err)
		return err
	}
	log.Printf("[DEBUG] %s", del)
	return nil
}
