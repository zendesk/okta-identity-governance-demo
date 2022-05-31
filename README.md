# Okta Identity Governance Demo

## Background
> :yellow_circle: Update with link to Medium article for the full background

This repository allows you to run a simplified version of Zendesk's internal Identity Governance tooling. It is designed to demonstrate how powerful [Attribute Based Access Control](https://developer.okta.com/books/api-security/authz/attribute-based/) combined with GitOps can be for speeding up onboarding of employees, allowing self service requests for permission changes, and providing a platform for adding authz for internal applications.

### Structure

#### Terraform
Terraform is used to manage all the infrastructure within Okta. This includes Application, Groups, Rules, and Authorization Servers. This demo includes a single SAML application. Terraform also fully suppports Okta [OAuth applications](https://registry.terraform.io/providers/okta/okta/latest/docs/resources/app_oauth). You can even provide custom claims for OAuth by purchasing [API Access Management](https://developer.okta.com/docs/concepts/api-access-management/) from Okta.

A Group Rule is used to assign the application to users automatically once the syncer assigns a value for the corresponding attribute on their user profile.

The Demo Application receives each value in the user's profile as separate attribute statements that can be used for authorization within the application. This can include RBAC type roles or other more discrete attributes.

For splitting multiple values into separate statements, you may need to enable the `SAML_SUPPORT_ARRAY_ATTRIBUTES` [Feature Flag](SAML_SUPPORT_ARRAY_ATTRIBUTES).

#### YAML
For this demo, we have included structure for [Teams](teams/) and [Attributes](attributes/).

Teams are simple collections of emails. This could include actual teams within your organization that are assigned identical permissions for different applications or more role based groups.

Attributes are used to assign values to each team. The values are combined across teams. Each Okta application could have a separate attribute to allow for delegated approvals and discrete 

#### Syncer
In order to generate each user's attributes, a simple Golang program is included within the [syncer](syncer/) directory. It reads the data stored in the `attributes` and `teams` folder to determine the desired values for all Custom Attributes in Okta. It then pulls all the current values and updates any user profiles that are different than the desired state. This allows for idempotent and concurrent actions, which is ideal for GitOps workflows.

This demo does not include any linting, continuous delivery, approval procecess, or flexibility of the schema. All of those can be added to suite your systems' needs.

## Running the demo
### Set Up

1. Local dependencies

    - Clone this repo locally
    - Install [Taskfile](https://taskfile.dev/#/installation)
    - Configure [Go](https://go.dev/doc/install)
    - Install [tfenv](https://github.com/tfutils/tfenv#installation)
        - You can also just install [Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)
1. Request a free trial from Okta
    - Sign up for a trial: https://www.okta.com/free-trial/
    - Log into the Okta Admin Portal
        - Browse to the **Directory** -> **Profile Editor** -> Click on the **User (default)** profile
            - `https://trial-{your account number}-admin.okta.com/admin/universaldirectory`
        - Click the **Add Attribute** button
        - Add an attribute for the Demo Application `DemoApplicationAttribute`
            - Insert Screenshot
        - *Optional* Edit the new attribute and set the **Source Priority** to **Inherit from Okta**
            - This is necessary for Okta accounts where other systems source profile attributes
            - Insert Screenshot
    - Create an API Token
        - Browse to **Security** -> **API** -> **Tokens**
            - `https://trial-{your account number}-admin.okta.com/admin/access/api/tokens`
        - Create a new token and save the value for the next step
    - Set the Token locally
        - Copy `.env.sample` to `.env`
        - Update the `OKTA_ORG_NAME` to match your Okta account
        - Update the `OKTA_API_TOKEN` to match the token from the previous step

### Applying to Okta
1. Applying Terraform
    1. Run `task terraform-apply`
    1. Type `yes` to apply
    1. Confirm that a new application has been created via the Okta Admin Console

1. Applying Attributes
    1. Add your email address to the [demo-admins]. You can create new users and teams. New users with `example.com` work within the Okta test accounts.
    1. Run `task build` to create the syncer binary for your platform
    1. Run `task sync` to update 
    1. Log in to your Okta account and see that you are assigned the `Demo Application` app
    1. Remove your email address from [demo-admins]
    1. Run `task sync` to remove the attribute from your account
    1. Log in to your Okta account and see that you no longer have access to the `Demo Application` app 

1. Further features
    1. If you you're interested in seeing how this works with more applications, add more teams and attributes
        1. In Okta's profile editor, add an additional attribute
        2. Create a new attribute file in this repo and map it to existing teams
        3. Run `task sync` and check each user's profile in Okta