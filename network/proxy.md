---
title: Configure Docker to use a proxy server
description: How to configure the Docker client to use a proxy server
keywords: network, networking, proxy, client
---

If your container needs to use an HTTP, HTTPS, or FTP proxy server, you can
configure it in different ways:

    - In Docker 17.07 and higher, you can
  [configure the Docker client](#configure-the-docker-client) to pass
  proxy information to containers automatically.

    - In Docker 17.06 and lower, you must
  [set appropriate environment variables](#use-environment-variables)
  within the container. You can do this when you build the image (which makes
  the image less portable) or when you create or run the container.
  
After configuring the Docker client, [configure proxy settings](#configure-proxy-settings). 
 

## Configure the Docker client

1.  On the Docker client, create or edit the file `~/.docker/config.json` in the
    home directory of the user which starts containers. Add JSON such as the
    following, substituting the type of proxy with `httpsProxy` or `ftpProxy` if
    necessary, and substituting the address and port of the proxy server. You
    can configure multiple proxy servers at the same time.

    You can optionally exclude hosts or ranges from going through the proxy
    server by setting a `noProxy` key to one or more comma-separated IP
    addresses or hosts. Using the `*` character as a wildcard for hosts and using CIDR notation for IP addresses is supported as
    shown in this example.

    ```json
    {
     "proxies":
     {
       "default":
       {
         "httpProxy": "http://127.0.0.1:3001",
         "httpsProxy": "http://127.0.0.1:3001",
         "noProxy": "*.test.example.com,.example2.com,127.0.0.0/8"
       }
     }
    }
    ```

    Save the file.

 2. When you create or start new containers, the environment variables are
    set automatically within the container.


## Use environment variables

### Set the environment variables manually

When you build the image, or using the `--env` flag when you create or run the
container, you can set one or more of the following variables to the appropriate
value. This method makes the image less portable, so if you have Docker 17.07
or higher, you should [configure the Docker client](#configure-the-docker-client)
instead.

| Variable      | Dockerfile example                                | `docker run` Example                                |
|:--------------|:--------------------------------------------------|:----------------------------------------------------|
| `HTTP_PROXY`  | `ENV HTTP_PROXY="http://127.0.0.1:3001"`          | `--env HTTP_PROXY="http://127.0.0.1:3001"`          |
| `HTTPS_PROXY` | `ENV HTTPS_PROXY="https://127.0.0.1:3001"`        | `--env HTTPS_PROXY="https://127.0.0.1:3001"`        |
| `FTP_PROXY`   | `ENV FTP_PROXY="ftp://127.0.0.1:3001"`            | `--env FTP_PROXY="ftp://127.0.0.1:3001"`            |
| `NO_PROXY`    | `ENV NO_PROXY="*.test.example.com,.example2.com"` | `--env NO_PROXY="*.test.example.com,.example2.com"` |


## Configure proxy settings 

Proxy servers can block connections to your web application once it's up and running. If you are behind a proxy server, you can specify proxy server host and port information in one of the following ways:

    -  Use  the `build-args` command in your Dockerfile for build-time argument passing.
    
        ```
        docker build --build-arg HTTP_PROXY=http://myyproxy.example.com
        ```

        The proxy information is only available during the docker build itself, and doesn’t persist in the resulting image (including     leaking through docker image history).
    
        This is similar to how `docker run -e` works. Refer to the docker run documentation f additional details.
    
        You can use the --build-arg flag without a value, in which case the value from the local environment is propagated into the Docker container being built:

        ```shell
        export HTTPS_PROXY=https://proxy.corp.example.com:8080

        docker build --build-arg HTTPS_PROXY  ......
        ```
        Using this flag does not alter the output you see when the ARG lines from the Dockerfile are echoed during the build process. The ARG instruction lets Dockerfile authors define values that users can set at build-time using the --build-arg flag. See the Dockerfile reference for more information about ARG.

    - Set `proxy` configuration defaults in the CLI configuration file. 
    
    - Use the `--env` option at runtime

> Warning: Using ENV instructions in a Dockerfile to set proxy-server values in the image is not recommended, because those variable values persist in the intermediate and final images. This is not desirable for the following reasons:
The resulting image can only be run in the environment in which it is built.
Information about the proxy-server can leak security-sensitive information (name of the proxy-server, or even the name and password if the proxy server requires authentication).
