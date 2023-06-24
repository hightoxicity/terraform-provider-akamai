package dns

import (
	"context"
	"net/http"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/dns"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/session"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResDnsRecord(t *testing.T) {
	dnsClient := dns.Client(session.Must(session.New()))

	var rec *dns.RecordBody

	notFound := &dns.Error{
		StatusCode: http.StatusNotFound,
	}

	// This test peforms a full life-cycle (CRUD) test
	t.Run("lifecycle test", func(t *testing.T) {
		client := &dns.Mock{}

		getCall := client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, notFound)

		parseCall := client.On("ParseRData",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]string"),
		).Return(nil)

		procCall := client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("[]string"),
			mock.AnythingOfType("string"),
		).Return(nil, nil)

		updateArguments := func(args mock.Arguments) {
			rec = args.Get(1).(*dns.RecordBody)
			getCall.ReturnArguments = mock.Arguments{rec, nil}
			parseCall.ReturnArguments = mock.Arguments{
				dnsClient.ParseRData(context.Background(), rec.RecordType, rec.Target),
			}
			procCall.ReturnArguments = mock.Arguments{rec.Target, nil}
		}

		client.On("CreateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(nil).Run(func(args mock.Arguments) {
			updateArguments(args)
		})

		client.On("UpdateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(nil).Run(func(args mock.Arguments) {
			updateArguments(args)
		})

		client.On("DeleteRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]bool"),
		).Return(nil).Run(func(mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{nil, notFound}
		})

		dataSourceName := "akamai_dns_record.a_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDnsRecord/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "A"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResDnsRecord/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "A"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("TXT record test", func(t *testing.T) {
		client := &dns.Mock{}
		escapedTarget := "\"Hel\\\\lo\\\"world\""

		client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"exampleterraform.io",
			"TXT",
		).Return(nil, notFound).Once()

		client.On("CreateRecord",
			mock.Anything, // ctx is irrelevant for this test
			&dns.RecordBody{
				Name:       "exampleterraform.io",
				RecordType: "TXT",
				TTL:        300,
				Active:     false,
				Target:     []string{escapedTarget},
			},
			"exampleterraform.io",
			[]bool{false},
		).Return(nil)

		client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"exampleterraform.io",
			"TXT",
		).Return(&dns.RecordBody{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{escapedTarget},
		}, nil).Once()

		client.On("ParseRData",
			mock.Anything,
			"TXT",
			[]string{escapedTarget},
		).Return(map[string]interface{}{
			"target": []string{escapedTarget},
		}).Once()

		client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			[]string{escapedTarget},
			"TXT",
		).Return([]string{escapedTarget}).Once()

		client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"exampleterraform.io",
			"TXT",
		).Return(&dns.RecordBody{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{escapedTarget},
		}, nil).Once()

		client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			[]string{escapedTarget},
			"TXT",
		).Return([]string{escapedTarget}).Once()

		client.On("ParseRData",
			mock.Anything,
			"TXT",
			[]string{escapedTarget},
		).Return(
			map[string]interface{}{
				"target": []string{escapedTarget},
			}).Once()

		client.On("DeleteRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]bool"),
		).Return(nil)
		dataSourceName := "akamai_dns_record.txt_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDnsRecord/create_basic_txt.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "TXT"),
							resource.TestCheckResourceAttr(dataSourceName, "target.#", "1"),
							resource.TestCheckResourceAttr(dataSourceName, "target.0", "\"Hel\\\\lo\\\"world\""),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("TXT record test with more than 255 chars target", func(t *testing.T) {
		client := &dns.Mock{}

		escapedTarget := "\"v=DKIM1; k=rsa; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmMZAR79x/6UHyyz6INnpuDC0dAMXUqcF6xE4a0nRN8R9FXfGRYhUHIOLCYTtj0PBG39A82lQAb/IB8epeEHkiJBye7/X8Khf4NsuQd2mkJuBgmSGsDXRI9evWE7+LcyxJaiZK/qKBAzVx37iZtbw7KhKimXhq+UztjmkVJ4qTIEkqa1z467Fw3Yyrr70JDv0a\" \"orve7Fs94v4Lr4/NTWHi7wVLUHl6TpBhqfJir7xVupeMLCcm2pbKkMd8eyeDDhYcrKTnubiuNGO/hqw7Sjt6WoVo8srz3+cvkEPzQbw0NRN4MVUTkcr4XGQjl3C2XSD7Gmtvjrm7sPuvdYtCADGJQIDAQAB\\010\""

		client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"exampleterraform.io",
			"TXT",
		).Return(nil, notFound).Once()

		client.On("CreateRecord",
			mock.Anything, // ctx is irrelevant for this test
			&dns.RecordBody{
				Name:       "exampleterraform.io",
				RecordType: "TXT",
				TTL:        300,
				Active:     false,
				Target:     []string{escapedTarget},
			},
			"exampleterraform.io",
			[]bool{false},
		).Return(nil)

		client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"exampleterraform.io",
			"TXT",
		).Return(&dns.RecordBody{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{escapedTarget},
		}, nil).Once()
		client.On("ParseRData",
			mock.Anything,
			"TXT",
			[]string{escapedTarget},
		).Return(map[string]interface{}{
			"target": []string{escapedTarget},
		}).Once()

		client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			[]string{escapedTarget},
			"TXT",
		).Return([]string{escapedTarget}).Once()

		client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"exampleterraform.io",
			"TXT",
		).Return(&dns.RecordBody{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{escapedTarget},
		}, nil).Once()

		client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			[]string{escapedTarget},
			"TXT",
		).Return([]string{escapedTarget}).Once()

		client.On("ParseRData",
			mock.Anything,
			"TXT",
			[]string{escapedTarget},
		).Return(
			map[string]interface{}{
				"target": []string{escapedTarget},
			}).Once()

		client.On("DeleteRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]bool"),
		).Return(nil)

		dataSourceName := "akamai_dns_record.txt_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDnsRecord/create_basic_txt_long.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "TXT"),
							resource.TestCheckResourceAttr(dataSourceName, "target.#", "1"),
							resource.TestCheckResourceAttr(dataSourceName, "target.0", "\"v=DKIM1; k=rsa; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmMZAR79x/6UHyyz6INnpuDC0dAMXUqcF6xE4a0nRN8R9FXfGRYhUHIOLCYTtj0PBG39A82lQAb/IB8epeEHkiJBye7/X8Khf4NsuQd2mkJuBgmSGsDXRI9evWE7+LcyxJaiZK/qKBAzVx37iZtbw7KhKimXhq+UztjmkVJ4qTIEkqa1z467Fw3Yyrr70JDv0a\" \"orve7Fs94v4Lr4/NTWHi7wVLUHl6TpBhqfJir7xVupeMLCcm2pbKkMd8eyeDDhYcrKTnubiuNGO/hqw7Sjt6WoVo8srz3+cvkEPzQbw0NRN4MVUTkcr4XGQjl3C2XSD7Gmtvjrm7sPuvdYtCADGJQIDAQAB\\010\""),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestTargetDiffSuppress(t *testing.T) {
	t.Run("target is computed and recordType is AAAA", func(t *testing.T) {
		config := schema.TestResourceDataRaw(t, getResourceDNSRecordSchema(), map[string]interface{}{"recordtype": "AAAA"})
		assert.False(t, dnsRecordTargetSuppress("target.#", "0", "", config))
	})
}
