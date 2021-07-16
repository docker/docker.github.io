---
title: Managing access tokens
description: Learn how to create and manage your personal Docker Hub access tokens to securely push and pull images programmatically.
keywords: docker hub, hub, security, PAT, personal access token
---

Docker Hub lets you create personal access tokens as alternatives to your password. You can use tokens to access Hub images from the Docker CLI.

Using personal access tokens provides some advantages over a password:
* You can investigate the last usage of the access token and disable or delete it if you find any suspicious activity.
* When logged in with an access token, you can't perform any admin activity on the account, including changing the password.

>**Important**
>
> Treat access tokens like your password and keep them secret. Store your
> tokens securely (for example, in a credential manager).
{: .important}

Access tokens are valuable for building integrations, as you can issue
multiple tokens &ndash; one for each integration &ndash; and revoke them at
any time.

   >**Note**
   >
   > If you have [two-factor authentication (2FA)](/docker-hub/2fa) enabled on
   > your account, you must create at least one personal access token. Otherwise,
   > you will be unable to log in to your account from the Docker CLI.


## Create an access token

You can create as many tokens as you need.

1. Log in to [hub.docker.com](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}.

2. Click on your username in the top right corner and select **[Account Settings](https://hub.docker.com/settings/general)**.

3. Select **[Security](https://hub.docker.com/settings/security) > New Access Token**.

4. Add a description for your token. Use something that indicates where
the token will be used, or set a purpose for the token.

      ![PAT Menu](images/hub-create-token.png){:width="700px"}

5. You can view the following permissions from the drop-down:

    | Features   | Free | Pro | Team |
    | --------------------- | ---- | ----- | ----- |
    | Personal Access Tokens (read, write, delete scope only)*    |  x   |   x   |   x   |
    | Scoped Personal Access Tokens  |      | x   |   x   |
    | Unlimited Personal Access Tokens   |      |       |   x   |

      >The following are the available access tokens for each plan: Free (1), Pro (5) and Team (Unlimited). For more information about a Pro or a Team plan, see [Docker Pricing](https://www.docker.com/pricing?utm_source=docker&utm_medium=webreferral&utm_campaign=docs_driven_upgrade){: target="_blank" rel="noopener" class="_"}.

6. Copy the token that appears on the screen and save it. You will not be able to retrieve the token once you close this prompt.

      ![Copy Token](images/hub-copy-token.png){:width="700px"}

## Use an access token

You can use an access token anywhere that requires your Docker Hub
password.

When logging in from your Docker CLI client (`docker login --username <username>`),
omit the password in the login command. Instead, enter your token when asked for a password.

   > **Note**
   >
   >If you have 2FA enabled, you must use a personal access token when logging in from the Docker CLI. 2FA is an optional, but more secure method of authentication.

## Modify existing tokens

You can rename, activate, deactivate, or delete a token as needed.

1. Access your tokens under **[Account Settings](https://hub-stage.docker.com/settings/general) > [Security](https://hub-stage.docker.com/settings/security)**.

2. Select a token and click **Delete** or **Edit**, or use the menu on
the far right of a token row to bring up the edit screen. You can also
select multiple tokens to delete at once.

      ![Delete or Edit](images/hub-delete-edit-token.png){: {:width="700px"}

   > **Note**
   >
   > You can also view the number of tokens that are activated and deactivated in the toolbar.

      ![Modify Token](images/hub-edit-token.png){:width="700px"}