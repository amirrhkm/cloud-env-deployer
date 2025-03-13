**Infrastructure-as-Code (IaC) Repository**
===========================================

**Overview**
------------

This repository contains Infrastructure-as-Code (IaC) scripts and configurations for provisioning and managing cloud resources across various cloud platforms. The repository includes code for both **Terraform** and **AWS CDK** to support multi-cloud deployments and provide flexibility for different use cases.

### **Key Technologies**:

*   **Terraform**: Declarative multi-cloud infrastructure management.
    
*   **AWS CDK**: Imperative code-based AWS infrastructure provisioning.
    

**Structure**
-------------

The repository is organized into separate folders for different IaC tools and environments:
 
```bash
├── terraform/                      # Terraform scripts and modules
│   ├── vpc/                        # Terraform VPC configuration
│   ├── ec2/                        # Terraform EC2 configuration
│   └── ...                         # Other Terraform resources
├── aws/                            # AWS CDK stacks and constructs
│   ├── ec2/                        # CDK EC2 configuration
│   │   └── provisioning-cdk.go     # CDK provisioning scripts
│   └── s3/                         # CDK S3 configuration
│   │   └── provisioning-cdk.go     # CDK provisioning scripts
└── README.md                       # Documentation and instructions
```
**Getting Started**
-------------------

### **Prerequisites**

Before using the IaC scripts in this repository, ensure that you have the following installed:

*   **Terraform**: Install Terraform
    
*   **AWS CDK**: [Install AWS CDK](https://docs.aws.amazon.com/cdk/latest/guide/getting_started.html)
    
*   **AWS CLI**: [Install AWS CLI](https://aws.amazon.com/cli/)
    
*   **Node.js**: Required for CDK (Install from [Node.js website](https://nodejs.org/))
    

### **AWS Authentication**

Before deploying any infrastructure, you need to configure your AWS credentials:

```bash
aws configure
```

Provide your AWS access key, secret access key, and region.

**Deploying with Terraform**
----------------------------

To deploy infrastructure using **Terraform**, navigate to the desired module (e.g., terraform/vpc/):

```bash
codecd terraform/vpc/
terraform init          # Initialize the Terraform project
terraform plan          # Preview the changes that will be applied
terraform apply         # Apply the changes to provision resources
 ```

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

### **Terraform Modules**

Terraform modules are reusable building blocks for infrastructure. Each module defines a specific set of resources (e.g., VPC, EC2, RDS). The Terraform modules are located in the terraform/ directory.

### **AWS CDK Stacks**

The **AWS CDK** stack files are located in the cdk/ directory. These are structured with reusable constructs in lib/ and entry-point scripts in bin/.

**Useful Commands**
-------------------

### **Terraform Commands**

*   `terraform init`: Initialize the Terraform project.
    
*   `terraform plan`: Show the resources to be created/modified.
    
*   `terraform apply`: Apply the changes and deploy resources.
    
*   `terraform destroy`: Remove all resources managed by Terraform.
    

### **AWS CDK Commands**

*   `cdk bootstrap`: Set up the environment for CDK deployment.
    
*   `cdk synth`: Synthesize the CloudFormation template.
    
*   `cdk deploy`: Deploy the infrastructure defined in the CDK app.
    
*   `cdk destroy`: Tear down the infrastructure.
    

**Cleanup**
-----------

To remove resources created by **Terraform**:

```bash
terraform destroy
```

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
