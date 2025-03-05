# go-auth0

go-auth0 is a Go client to Auth0 APIs.

## Using

Use go modules. Import just the module(s) you need as follows:

```go
import (
	"github.com/zenoss/go-auth0/auth0"
	"github.com/zenoss/go-auth0/auth0/mgmt"
)
```

Then run `go mod tidy` to add the requirement to go.mod and go.sum. You will
also have to run `go mod vendor` if your project vendors its dependencies.

## Contributing

1. Create a branch based on master named ZING-???? after a JIRA issue.
2. Make your changes in the branch.
3. Add unit tests for your changes.
4. Run all the unit tests: `make test`
5. Commit your changes.
6. Push your changes: `git push -u origin ZING-????`
7. Open a pull request against master: `gh pr create -f`
8. Merge your pull request after getting it reviewed and approved.

## Testing

The tests for this project are true integration tests that use the real Auth0
and authorization extension APIs. Before running tests, you must ensure that
the `.env` files contains the following variables.

```sh
# Get these values from the auth0-common ConfigMap.
AUTH0_DOMAIN=
AUTH0_MANAGEMENT_API_URL="<"
AUTH0_MANAGEMENT_API_AUDIENCE=
AUTH0_AUTHORIZATION_API_URL=
AUTH0_AUTHORIZATION_API_AUDIENCE=

# Get these from the auth0-authz machine-to-machine application.
AUTH0_MANAGEMENT_CLIENT_ID=
AUTH0_MANAGEMENT_CLIENT_SECRET=

# These are the same as above.
AUTH0_AUTHORIZATION_CLIENT_ID=
AUTH0_AUTHORIZATION_CLIENT_SECRET=
```

When I last tried this, it didn't work because no machine-to-machine
application had all of the permissions needed to run the tests. Specifically,
the `auth0-authz` application didn't have the following permissions. The only
application that had these permissions was `auth0-deploy-cli-extension`, and it
didn't have a lot of different permissions the tests required. For this reason,
I temporarily added the following permissions to the `auth0-authz` application
for the `Auth0 Management API` API in the dev stack, ran the tests, then
removed the extra permissions.

- `create:connections`
- `update:connections`
- `delete:connections`

You can run the integration tests as follows.

```sh
make test
```

## Releasing

Releases of this library should be made the same way as releases for any go
library.

1. Create a tag for the version on the master branch: `git tag vX.Y.Z`
2. Push the tag: `git push --tags`

_P.S._ Use the vX.Y.Z format for your tag. This is the convention preferred by
go mod.

_P.P.S._ Avoid making backwards-incompatible changes if at all possible. The
major version (e.g. 4 in 4.3.2) should only be incremented when a
backwards-incompatible change has been made. This major version change also
requires the go package have its prefix (e.g. v4/) incremented. All users of
the library will then have to update their import paths to opt-in to the new
major version.
