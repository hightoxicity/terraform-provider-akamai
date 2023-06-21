provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_dns_record" "txt_record" {
  zone       = "exampleterraform.io"
  name       = "exampleterraform.io"
  recordtype = "TXT"
  active     = true
  ttl        = 300
  target     = ["v=DKIM1\\; k=rsa\\; p=MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmMZAR79x/6UHyyz6INnpuDC0dAMXUqcF6xE4a0nRN8R9FXfGRYhUHIOLCYTtj0PBG39A82lQAb/IB8epeEHkiJBye7/X8Khf4NsuQd2mkJuBgmSGsDXRI9evWE7+LcyxJaiZK/qKBAzVx37iZtbw7KhKimXhq+UztjmkVJ4qTIEkqa1z467Fw3Yyrr70JDv0aorve7Fs94v4Lr4/NTWHi7wVLUHl6TpBhqfJir7xVupeMLCcm2pbKkMd8eyeDDhYcrKTnubiuNGO/hqw7Sjt6WoVo8srz3+cvkEPzQbw0NRN4MVUTkcr4XGQjl3C2XSD7Gmtvjrm7sPuvdYtCADGJQIDAQAB\\010"]
}
