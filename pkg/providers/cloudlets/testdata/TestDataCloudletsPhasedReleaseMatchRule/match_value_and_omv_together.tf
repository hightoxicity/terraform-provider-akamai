provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_phased_release_match_rule" "test" {

  match_rules {
    name = "rule 2"
    matches {
      match_type     = "hostname"
      match_operator = "equals"
      match_value    = "example.ex"
      object_match_value {
        type  = "simple"
        value = ["abc"]
      }
    }
    forward_settings {
      origin_id = "1234"
      percent   = 10
    }
  }
}