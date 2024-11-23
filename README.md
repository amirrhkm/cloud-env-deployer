**Infrastructure-as-Code (IaC) Repository**
===========================================

**Overview**
------------

This repository contains Infrastructure-as-Code (IaC) scripts and configurations for provisioning and managing cloud resources across various cloud platforms. The repository includes **AWS CDK** to support cloud deployments and provide flexibility for different use cases.

### **Key Technologies**:

*   **AWS CDK**: Imperative code-based AWS infrastructure provisioning.
    

**Structure**
-------------

The repository is organized into separate folders for different IaC tools and environments:
 
```bash
├── cdk/                # AWS CDK stacks and constructs
│   ├── lib/            # CDK VPC, EC2, S3 configurations
│   ├── bin/            # CDK application entry points
│   └── ...             # Other CDK resources
├── scripts/            # Utility scripts (e.g., deployment, testing)
└── README.md           # Documentation and instructions
```
**Getting Started**
-------------------

### **Prerequisites**

Before using the IaC scripts in this repository, ensure that you have the following installed:

*   **AWS CDK**: [Install AWS CDK](https://docs.aws.amazon.com/cdk/latest/guide/getting_started.html)
    
*   **AWS CLI**: [Install AWS CLI](https://aws.amazon.com/cli/)
    
*   **Node.js**: Required for CDK (Install from [Node.js website](https://nodejs.org/))

*   **GoLang**: Required for CDK (Install from [Go website](https://go.dev/dl/))
    

### **AWS Authentication**

Before deploying any infrastructure, you need to configure your AWS credentials:

```bash
aws configure
```

Provide your AWS access key, secret access key, and region.

**Deploying with AWS CDK**
--------------------------

To deploy infrastructure using **AWS CDK**, follow these steps:

1.  Install dependencies:
    ```bash
    npm install
    ```
    
2.  Bootstrap the CDK environment (only required for first-time use):
    ```bash
    cdk bootstrap
    ```
    
3.  Deploy the CDK stack:
    ```bash
    cdk deploy
    ```
    

**Repository Structure**
------------------------

### **AWS CDK Stacks**

The **AWS CDK** stack files are located in the cdk/ directory. These are structured with reusable constructs in lib/ and entry-point scripts in bin/.

**Useful Commands**
-------------------

### **AWS CDK Commands**

*   `cdk bootstrap`: Set up the environment for CDK deployment.
    
*   `cdk synth`: Synthesize the CloudFormation template.
    
*   `cdk deploy`: Deploy the infrastructure defined in the CDK app.
    
*   `cdk destroy`: Tear down the infrastructure.
    

**Cleanup**
-----------

To remove resources created by **AWS CDK**:

```bash
cdk destroy
```

**SSB RONPOS Staging Environment**
----------------------------------

Reconfigure AWS SSO to use the latest session token
1. set AWS environment variable with the latest session token in the project
```bash
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_SESSION_TOKEN=
```

2. reconfigure AWS credentials to include the new profile
```bash
rm -rf ~/.aws/credentials
echo "[Profile Name]
aws_access_key_id = 
aws_secret_access_key = 
aws_session_token = 
" >> ~/.aws/credentials
```
