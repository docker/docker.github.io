Your configuration options can contain environment variables. Compose uses the
variable values from the shell environment in which `docker-compose` is run. For
example, suppose the shell contains `POSTGRES_VERSION=9.3` and you supply this
configuration:

```yaml
db:
  image: "postgres:${POSTGRES_VERSION}"
```

When you run `docker-compose up` with this configuration, Compose looks for the
`POSTGRES_VERSION` environment variable in the shell and substitutes its value
in. For this example, Compose resolves the `image` to `postgres:9.3` before
running the configuration.

If an environment variable is not set, Compose substitutes with an empty
string. In the example above, if `POSTGRES_VERSION` is not set, the value for
the `image` option is `postgres:`.

You can set default values for environment variables using a
[`.env` file](/compose/env-file/), which Compose automatically looks for in
project directory (parent folder of your Compose file). 
Values set in the shell environment override those set in the `.env` file.

> Note when using docker stack deploy
>
> The `.env file` feature only works when you use the `docker-compose up` command
> and does not work with `docker stack deploy`.
{: .important }

Both `$VARIABLE` and `${VARIABLE}` syntax are supported. Additionally when using
the [2.1 file format](/compose/compose-file/compose-versioning/#version-21), it is possible to
provide inline default values using typical shell syntax:

- `${VARIABLE:-default}` evaluates to `default` if `VARIABLE` is unset or
  empty in the environment.
- `${VARIABLE-default}` evaluates to `default` only if `VARIABLE` is unset
  in the environment.

Similarly, the following syntax allows you to specify mandatory variables when
using **Compose 1.19.0+** or `docker stack deploy` **18.04.0+**:

- `${VARIABLE:?err}` exits with an error message containing `err` if
  `VARIABLE` is unset or empty in the environment.
- `${VARIABLE?err}` exits with an error message containing `err` if
  `VARIABLE` is unset in the environment.

Other extended shell-style features, such as `${VARIABLE/foo/bar}`, are not
supported.

You can use a `$$` (double-dollar sign) when your configuration needs a literal
dollar sign. This also prevents Compose from interpolating a value, so a `$$`
allows you to refer to environment variables that you don't want processed by
Compose.

```yaml
web:
  build: .
  command: "$$VAR_NOT_INTERPOLATED_BY_COMPOSE"
```

If you forget and use a single dollar sign (`$`), Compose interprets the value
as an environment variable and warns you:

```
The VAR_NOT_INTERPOLATED_BY_COMPOSE is not set. Substituting an empty string.
```
