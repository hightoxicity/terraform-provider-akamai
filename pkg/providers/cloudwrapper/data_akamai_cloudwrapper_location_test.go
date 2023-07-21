package cloudwrapper

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataLocation(t *testing.T) {
	location := testDataForCWLocation{
		locations: []cloudwrapper.Location{
			{
				LocationID:   1,
				LocationName: "US East",
				TrafficTypes: []cloudwrapper.TrafficTypeItem{
					{
						TrafficTypeID: 1,
						TrafficType:   "LIVE",
						MapName:       "cw-s-use-live",
					},
					{
						TrafficTypeID: 2,
						TrafficType:   "LIVE_VOD",
						MapName:       "cw-s-use",
					},
				},
				MultiCDNLocationID: "018",
			},
			{
				LocationID:   2,
				LocationName: "US West",
				TrafficTypes: []cloudwrapper.TrafficTypeItem{
					{
						TrafficTypeID: 3,
						TrafficType:   "LIVE_VOD",
						MapName:       "cw-s-usw",
					},
					{
						TrafficTypeID: 4,
						TrafficType:   "LIVE",
						MapName:       "cw-s-usw-live",
					},
				},
				MultiCDNLocationID: "020",
			},
		},
	}
	tests := map[string]struct {
		configPath string
		init       func(*testing.T, *cloudwrapper.Mock, testDataForCWLocation)
		mockData   testDataForCWLocation
		error      *regexp.Regexp
	}{
		"happy path": {
			configPath: "testdata/TestDataLocation/location.tf",
			init: func(t *testing.T, m *cloudwrapper.Mock, testData testDataForCWLocation) {
				expectListLocations(m, testData, 5)
			},
			mockData: location,
		},
		"no location": {
			configPath: "testdata/TestDataLocation/no_location.tf",
			init: func(t *testing.T, m *cloudwrapper.Mock, testData testDataForCWLocation) {
				expectListLocations(m, testData, 1)
			},
			mockData: location,
			error:    regexp.MustCompile("no location with given location name and traffic type"),
		},
		"invalid type": {
			configPath: "testdata/TestDataLocation/invalid_type.tf",
			error:      regexp.MustCompile(`Attribute traffic_type value must be one of: \["\\"LIVE\\"" "\\"LIVE_VOD\\""\s+"\\"WEB_STANDARD_TLS\\"" "\\"WEB_ENHANCED_TLS\\""], got: "TEST"`),
		},
		"error listing locations": {
			configPath: "testdata/TestDataLocation/location.tf",
			init: func(t *testing.T, m *cloudwrapper.Mock, testData testDataForCWLocation) {
				expectListLocationsWithError(m, 1)
			},
			error: regexp.MustCompile("list locations failed"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cloudwrapper.Mock{}
			if test.init != nil {
				test.init(t, client, test.mockData)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: newProviderFactory(withMockClient(client)),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, test.configPath),
						Check:       checkCloudWrapperLocationAttrs(),
						ExpectError: test.error,
					},
				},
			})

			client.AssertExpectations(t)
		})
	}
}

type testDataForCWLocation struct {
	locations []cloudwrapper.Location
}

func checkCloudWrapperLocationAttrs() resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_location.test", "traffic_type_id", "3"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudwrapper_location.test", "location_id", "cw-s-usw"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttrSet("data.akamai_cloudwrapper_location.test", "id"))

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func expectListLocations(client *cloudwrapper.Mock, data testDataForCWLocation, timesToRun int) {
	listLocationsRes := cloudwrapper.ListLocationResponse{
		Locations: data.locations,
	}
	client.On("ListLocations", mock.Anything).Return(&listLocationsRes, nil).Times(timesToRun)
}

func expectListLocationsWithError(client *cloudwrapper.Mock, timesToRun int) {
	client.On("ListLocations", mock.Anything).Return(nil, fmt.Errorf("list locations failed")).Times(timesToRun)
}
