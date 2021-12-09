package cloudlets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudletsForwardRewriteMatchRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudletsForwardRewriteMatchRuleRead,
		Schema: map[string]*schema.Schema{
			"match_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Defines a set of rules for policy",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the rule",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of Cloudlet the rule is for",
						},
						"start": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The start time for this match (in seconds since the epoch)",
						},
						"end": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The end time for this match (in seconds since the epoch)",
						},
						"matches": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Defines a set of match objects",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The type of match used",
									},
									"match_value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Depends on the matchType",
									},
									"match_operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Valid entries for this property: contains, exists, and equals",
									},
									"case_sensitive": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, the match is case sensitive",
									},
									"negate": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, negates the match",
									},
									"check_ips": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "For clientip, continent, countrycode, proxy, and regioncode match types, the part of the request that determines the IP address to use",
									},
									"object_match_value": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "An object used when a rule either includes more complex match criteria, like multiple value attributes",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Optional: true,
													Description: "If using a match type that supports name attributes, enter the value in the incoming request to match on. " +
														"The following match types support this property: cookie, header, parameter, and query",
												},
												"type": {
													Type:     schema.TypeString,
													Required: true,
													Description: "The array type, which can be one of the following: object or simple. " +
														"Use the simple option when adding only an array of string-based values",
												},
												"name_case_sensitive": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Set to true if the entry for the name property should be evaluated based on case sensitivity",
												},
												"name_has_wildcard": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Set to true if the entry for the name property includes wildcards",
												},
												"options": {
													Type:        schema.TypeSet,
													MaxItems:    1,
													Optional:    true,
													Description: "If using the object type, use this set to list the values to match on (use only with the object type)",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:        schema.TypeList,
																Elem:        &schema.Schema{Type: schema.TypeString},
																Optional:    true,
																Description: "The value attributes in the incoming request to match on",
															},
															"value_has_wildcard": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Set to true if the entries for the value property include wildcards",
															},
															"value_case_sensitive": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Set to true if the entries for the value property should be evaluated based on case sensitivity",
															},
															"value_escaped": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Set to true if provided value should be compared in escaped form",
															},
														},
													},
												},
												"value": {
													Type:        schema.TypeList,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Optional:    true,
													Description: "The value attributes in the incoming request to match on (use only with simple or range type)",
												},
											},
										},
									},
								},
							},
						},
						"match_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "If using a URL match, this property is the URL that the Cloudlet uses to match the incoming request",
						},
						"forward_settings": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 1,
							Description: "This property defines data used to construct a new request URL if all conditions are met. " +
								"If all of the conditions you set are true, then the Edge Server returns an HTTP response from the rewritten URL",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The ID of the Conditional Origin requests are forwarded to",
									},
									"use_incoming_query_string": {
										Type:     schema.TypeBool,
										Optional: true,
										Description: "If set to true, the Cloudlet includes the query string from the request " +
											"in the rewritten or forwarded URL.",
									},
									"path_and_qs": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "If a value is provided and match conditions are met, this property defines " +
											"the path/resource/query string to rewrite URL for the incoming request.",
									},
								},
							},
						},
					},
				},
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A match_rules JSON structure generated from the schema",
			},
		},
	}
}

func dataSourceCloudletsForwardRewriteMatchRuleRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	matchRulesList, err := tools.GetListValue("match_rules", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = setMatchRuleSchemaType(matchRulesList, cloudlets.MatchRuleTypeFR); err != nil {
		return diag.FromErr(err)
	}

	matchRules, err := getMatchRulesFR(matchRulesList)
	if err != nil {
		return diag.Errorf("'match_rules' - %s", err)
	}

	jsonBody, err := json.MarshalIndent(matchRules, "", "  ")
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	hashID, err := getMatchRulesHashID(matchRules)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hashID)

	return nil
}

func getMatchRulesFR(matchRules []interface{}) (*cloudlets.MatchRules, error) {
	result := make(cloudlets.MatchRules, 0, len(matchRules))
	for _, mr := range matchRules {
		matchRuleMap, ok := mr.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("match rule is of invalid type: %T", mr)
		}

		matches, err := getMatchCriteriaFR(matchRuleMap["matches"].([]interface{}))
		if err != nil {
			return nil, err
		}

		matchRule := cloudlets.MatchRuleFR{
			Name:     getStringValue(matchRuleMap, "name"),
			Type:     cloudlets.MatchRuleTypeFR,
			MatchURL: getStringValue(matchRuleMap, "match_url"),
			Start:    getIntValue(matchRuleMap, "start"),
			End:      getIntValue(matchRuleMap, "end"),
			Matches:  matches,
		}

		// Schema guarantees that "forward_settings" will be present and of type *schema.Set
		settings, ok := matchRuleMap["forward_settings"].(*schema.Set)
		if !ok {
			return nil, fmt.Errorf("%v: 'forward_settings' should be an *schema.Set", tools.ErrInvalidType)
		}
		for _, element := range settings.List() {
			entries := element.(map[string]interface{})
			matchRule.ForwardSettings = cloudlets.ForwardSettingsFR{
				OriginID:               entries["origin_id"].(string),
				PathAndQS:              entries["path_and_qs"].(string),
				UseIncomingQueryString: entries["use_incoming_query_string"].(bool),
			}
		}

		result = append(result, matchRule)
	}
	return &result, nil
}

func getMatchCriteriaFR(matches []interface{}) ([]cloudlets.MatchCriteriaFR, error) {
	result := make([]cloudlets.MatchCriteriaFR, 0, len(matches))
	for _, criteria := range matches {
		criteriaMap, ok := criteria.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("matches is of invalid type")
		}

		omv, err := parseObjectMatchValue(criteriaMap, getObjectMatchValueObjectOrSimple)
		if err != nil {
			return nil, err
		}

		matchCriterion := cloudlets.MatchCriteriaFR{
			MatchType:        getStringValue(criteriaMap, "match_type"),
			MatchValue:       getStringValue(criteriaMap, "match_value"),
			MatchOperator:    cloudlets.MatchOperator(getStringValue(criteriaMap, "match_operator")),
			CaseSensitive:    getBoolValue(criteriaMap, "case_sensitive"),
			Negate:           getBoolValue(criteriaMap, "negate"),
			CheckIPs:         cloudlets.CheckIPs(getStringValue(criteriaMap, "check_ips")),
			ObjectMatchValue: omv,
		}

		result = append(result, matchCriterion)
	}
	return result, nil
}
