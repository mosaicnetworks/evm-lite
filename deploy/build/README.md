# BUILD

This folder contains scripts to generate different types of machine images for
with evm-lite installed.

## Docker

Produce a versioned Docker image with `evml` using a classic Dockerfile. 

Requires Docker Engine. 

```bash
$ make docker
```

## Amazon Machine Image (AMI)

Produce an AMI required to run instances (virtual servers) in the AWS cloud with 
`evml` pre-installed.

Requires Hashicorp Packer, an AWS account, and AWS access keys.

From Packer documentation:

> This builder builds an AMI by launching an EC2 instance from a source AMI, 
> provisioning that running machine, and then creating an AMI from that machine. 
> This is all done in your own AWS account. The builder will create temporary 
> keypairs, security group rules, etc. that provide it temporary access to the 
> instance while the image is being created. This simplifies configuration quite 
> a bit.
> 
> The builder does not manage AMIs. Once it creates an AMI and stores it in your 
> account, it is up to you to use, delete, etc. the AMI.

Provide the AWS access key in the `ami/secret.json` file: 

```json
{
    "aws_access_key" : "...",
    "aws_secret_key" : "..."
}
```

```bash
$ make ami
```


