---
description: Amazon Web Services driver for machine
keywords: machine, Amazon Web Services, driver
title: Amazon Web Services
---

Create machines on [Amazon Web Services](https://aws.amazon.com).

To create machines on [Amazon Web Services](https://aws.amazon.com), you must supply two parameters: the AWS Access Key ID and the AWS Secret Access Key.

## Configuring credentials

Before using the amazonec2 driver, ensure that you've configured credentials.

### AWS credential file

One way to configure credentials is to use the standard credential file for Amazon AWS `~/.aws/credentials` file, which might look like:

    [default]
    aws_access_key_id = AKID1234567890
    aws_secret_access_key = MY-SECRET-KEY

On Mac OS or various flavors of Linux you can install the [AWS Command Line Interface](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html#cli-quick-configuration) (`aws cli`) in the terminal and use the `aws configure` command which guides you through the creation of the credentials file.

This is the simplest method, you can then create a new machine with:

    $ docker-machine create --driver amazonec2 aws01

### Command line flags

Alternatively, you can use the flags `--amazonec2-access-key` and `--amazonec2-secret-key` on the command line:

    $ docker-machine create --driver amazonec2 --amazonec2-access-key AKI******* --amazonec2-secret-key 8T93C*******  aws01

### Environment variables

You can use environment variables:

    $ export AWS_ACCESS_KEY_ID=AKID1234567890
    $ export AWS_SECRET_ACCESS_KEY=MY-SECRET-KEY
    $ docker-machine create --driver amazonec2 aws01

## Options

-   `--amazonec2-access-key`: Your access key ID for the Amazon Web Services API.
-   `--amazonec2-ami`: The AMI ID of the instance to use.
-   `--amazonec2-block-duration-minutes`: AWS spot instance duration in minutes (60, 120, 180, 240, 300, or 360).
-   `--amazonec2-device-name`: The root device name of the instance.
-   `--amazonec2-endpoint`: Optional endpoint URL (hostname only or a fully qualified URI).
-   `--amazonec2-iam-instance-profile`: The AWS IAM role name to be used as the instance profile.
-   `--amazonec2-insecure-transport`: Disable SSL when sending requests.
-   `--amazonec2-instance-type`: The instance type to run.
-   `--amazonec2-keypair-name`: AWS keypair to use; requires `--amazonec2-ssh-keypath`.
-   `--amazonec2-monitoring`: Enable CloudWatch Monitoring.
-   `--amazonec2-open-port`: Make the specified port number accessible from the Internet.
-   `--amazonec2-private-address-only`: Use the private IP address only.
-   `--amazonec2-region`: The region to use when launching the instance.
-   `--amazonec2-request-spot-instance`: Use spot instances.
-   `--amazonec2-retries`: Set retry count for recoverable failures (use `-1` to disable).
-   `--amazonec2-root-size`: The root disk size of the instance (in GB).
-   `--amazonec2-secret-key`: Your secret access key for the Amazon Web Services API.
-   `--amazonec2-security-group`: AWS VPC security group name.
-   `--amazonec2-security-group-readonly`: Skip adding default rules to security groups.
-   `--amazonec2-session-token`: Your session token for the Amazon Web Services API.
-   `--amazonec2-spot-price`: Spot instance bid price in dollars. Requires the `--amazonec2-request-spot-instance` flag.
-   `--amazonec2-ssh-keypath`: Path to private key file to use for instance. Requires a matching public key with `.pub` extension to exist.
-   `--amazonec2-ssh-user`: The SSH login username, which must match the default SSH user set in the AMI being used.
-   `--amazonec2-subnet-id`: AWS VPC subnet ID.
-   `--amazonec2-tags`: A comma-separated list of AWS extra tag key-value pairs. For example, `key1,value1,key2,value2`.
-   `--amazonec2-use-ebs-optimized-instance`: Create an EBS Optimized Instance. Instance type must support it.
-   `--amazonec2-use-private-address`: Use the private IP address for docker-machine, but still create a public IP address.
-   `--amazonec2-userdata`: Path to file with cloud-init user data.
-   `--amazonec2-volume-type`: The Amazon EBS volume type to be attached to the instance.
-   `--amazonec2-vpc-id`: Your VPC ID to launch the instance in.
-   `--amazonec2-zone`: The AWS zone to launch the instance in (one of a,b,c,d, and e).



#### Environment variables and default values:

| CLI option                               | Environment variable          | Default          |
|:-----------------------------------------|:------------------------------|:-----------------|
| `--amazonec2-access-key`                 | `AWS_ACCESS_KEY_ID`           | -                |
| `--amazonec2-ami`                        | `AWS_AMI`                     | `ami-c60b90d1`   |
| `--amazonec2-block-duration-minutes`     | -                             | -                |
| `--amazonec2-device-name`                | `AWS_DEVICE_NAME`             | `/dev/sda1`      |
| `--amazonec2-endpoint`                   | `AWS_ENDPOINT`                | -                |
| `--amazonec2-iam-instance-profile`       | `AWS_INSTANCE_PROFILE`        | -                |
| `--amazonec2-insecure-transport`         | `AWS_INSECURE_TRANSPORT`      | -                |
| `--amazonec2-instance-type`              | `AWS_INSTANCE_TYPE`           | `t2.micro`       |
| `--amazonec2-keypair-name`               | `AWS_KEYPAIR_NAME`            | -                |
| `--amazonec2-monitoring`                 | -                             | `false`          |
| `--amazonec2-open-port`                  | -                             | -                |
| `--amazonec2-private-address-only`       | -                             | `false`          |
| `--amazonec2-region`                     | `AWS_DEFAULT_REGION`          | `us-east-1`      |
| `--amazonec2-request-spot-instance`      | -                             | `false`          |
| `--amazonec2-retries`                    | -                             | `5`              |
| `--amazonec2-root-size`                  | `AWS_ROOT_SIZE`               | `16`             |
| `--amazonec2-secret-key`                 | `AWS_SECRET_ACCESS_KEY`       | -                |
| `--amazonec2-security-group`             | `AWS_SECURITY_GROUP`          | `docker-machine` |
| `--amazonec2-security-group-readonly`    | `AWS_SECURITY_GROUP_READONLY` | `false`          |
| `--amazonec2-session-token`              | `AWS_SESSION_TOKEN`           | -                |
| `--amazonec2-spot-price`                 | -                             | `0.50`           |
| `--amazonec2-ssh-keypath`                | `AWS_SSH_KEYPATH`             | -                |
| `--amazonec2-ssh-user`                   | `AWS_SSH_USER`                | `ubuntu`         |
| `--amazonec2-subnet-id`                  | `AWS_SUBNET_ID`               | -                |
| `--amazonec2-tags`                       | `AWS_TAGS`                    | -                |
| `--amazonec2-use-ebs-optimized-instance` | -                             | `false`          |
| `--amazonec2-use-private-address`        | -                             | `false`          |
| `--amazonec2-userdata`                   | `AWS_USERDATA`                | -                |
| `--amazonec2-volume-type`                | `AWS_VOLUME_TYPE`             | `gp2`            |
| `--amazonec2-vpc-id`                     | `AWS_VPC_ID`                  | -                |
| `--amazonec2-zone`                       | `AWS_ZONE`                    | `a`              |

## Default AMIs

By default, the Amazon EC2 driver uses a daily image of `Ubuntu 16.04 LTS`.

| Region         | AMI ID       |
| -------------- | ------------ |
| ap-northeast-1 | ami-785c491f |
| ap-northeast-2 | ami-94d20dfa |
| ap-southeast-1 | ami-2378f540 |
| ap-southeast-2 | ami-e94e5e8a |
| ap-south-1     | ami-49e59a26 |
| ca-central-1   | ami-7ed56a1a |
| cn-north-1     | ami-a163b4cc |
| eu-central-1   | ami-1c45e273 |
| eu-west-1      | ami-6d48500b |
| eu-west-2      | ami-cc7066a8 |
| eu-west-3      | ami-c1cf79bc |
| sa-east-1      | ami-34afc458 |
| us-east-1      | ami-d15a75c7 |
| us-east-2      | ami-8b92b4ee |
| us-west-1      | ami-73f7da13 |
| us-west-2      | ami-835b4efa |
| us-gov-west-1  | ami-939412f2 |

## Security Group

A security group is created and associated to the host. This security group has the following ports opened inbound:

-   ssh (22/tcp)
-   docker (2376/tcp)
-   swarm (3376/tcp), only if the node is a swarm master

If you specify a security group yourself using the `--amazonec2-security-group` flag, the above ports are checked and opened and the security group is modified.
If you want more ports to be opened such as application-specific ports, use the AWS console and modify the configuration manually.

Port ranges will not be checked. Only SG rules having single ports are taken in account.

If the security needs to be modified, `0.0.0.0/0` will be used as source by default.

## VPC ID

Your default VPC ID is determined at the start of a command. In some cases, either because your account does not have a default VPC, or you do not want to use the default one, you can specify a VPC with the `--amazonec2-vpc-id` flag.

### To find the VPC ID:

1.  Login to the AWS console.
2.  Go to **Services -> VPC -> Your VPCs**.
3.  Locate the VPC ID you want from the *_VPC_* column.
4.  Go to **Services -> VPC -> Subnets**. Examine the _Availability Zone_ column to verify that zone `a` exists and matches your VPC ID. For example, `us-east1-a` is in the `a` availability zone. If the `a` zone is not present, you can create a new subnet in that zone or specify a different zone when you create the machine.

### To create a machine with a non-default VPC-ID:

    $ docker-machine create --driver amazonec2 --amazonec2-access-key AKI******* --amazonec2-secret-key 8T93C********* --amazonec2-vpc-id vpc-****** aws02

This example assumes the VPC ID was found in the `a` availability zone. Use the`--amazonec2-zone` flag to specify a zone other than the `a` zone. For example, `--amazonec2-zone c` signifies `us-east1-c`.

## VPC Connectivity
Docker Machine uses SSH to complete the set up of instances in EC2 and requires the ability to access the instance directly.

If you use the flag `--amazonec2-private-address-only`, ensure that you can access the new instance from within the internal network of the VPC, such as a corporate VPN to the VPC, a VPN instance inside the VPC, or using `docker-machine` from an instance within your VPC.

Configuration of VPCs is beyond the scope of this guide. However, the first step in troubleshooting is making sure that you are using private subnets that follow the design guidance in the [AWS VPC User Guide](http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/VPC_Scenario2.html) and have some form of NAT available so that the setup process can access the internet to complete the setup.

## Custom AMI and SSH username

The default SSH username for the default AMIs is `ubuntu`.

You need to change the SSH username only if the custom AMI you use has a different SSH username.

You can change the SSH username with the `--amazonec2-ssh-user` according to the AMI you selected with the `--amazonec2-ami` option.
