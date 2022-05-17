provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_edge_redirector_match_rule" "test" {

  match_rules {
    redirect_url     = "/ddd"
    status_code      = 301
    use_relative_url = "copy_scheme_hostname"
    matches {
      match_type     = "clientip"
      match_operator = "equals"
      check_ips      = "invalid"
      object_match_value {
        type  = "simple"
        value = ["fghi"]
      }
    }
  }
}