resource "okta_app_saml" "demo-application" {
  label                    = "Demo Application"
  sso_url                  = "https://example.com"
  recipient                = "https://example.com"
  destination              = "https://example.com"
  audience                 = "https://example.com/audience"
  subject_name_id_template = "$${user.userName}"
  subject_name_id_format   = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
  response_signed          = true
  signature_algorithm      = "RSA_SHA256"
  digest_algorithm         = "SHA256"
  honor_force_authn        = false
  authn_context_class_ref  = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"

  attribute_statements {
      name   = "groups"
      type   = "EXPRESSION"
      values = ["Arrays.flatten(user.DemoApplicationAttribute)"]
  }
}

resource "okta_group" "demo-application-users" {
  name        = "Demo Application Users"
  description = "Users assigned to Demo Application"
}

resource "okta_group_rule" "demo-application-group-rule" {
  name              = "example"
  status            = "ACTIVE"
  group_assignments = [
    okta_group.demo-application-users.id
  ]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "! Arrays.isEmpty(user.DemoApplicationAttribute)"
}

resource "okta_app_group_assignment" "demo-application-assignment" {
  depends_on = [okta_group.demo-application-users]
  app_id     = okta_app_saml.demo-application.id
  group_id   = okta_group.demo-application-users.id

  lifecycle {
    ignore_changes = [
      priority,
      profile,
    ]
  }
}