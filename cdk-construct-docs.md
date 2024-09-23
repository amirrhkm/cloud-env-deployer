# Amazon Web Services (AWS) Cloud Development Kit (CDK)

#### Local Development Environment 
Code → `cdk deploy` → CloudFormation Template → CloudFormation Deployment

#### Automated Deployments Environment
Code → Pull Request → Trigger Github Action → Trigger CloudFormation Deployment

---

    CDK construct: create AWS service resources with code which then compiled into AWS CloudFormation (IaC) templates.

---

### Core Concepts

#### L1 Constructs
- CloudFormation (CFN) (Define a single AWS resource)
- Define every settings associated with the resource
- Use Case: need to control every single setting of the resource, already extremely familiar with CloudFormation
- **Not Recommended**

#### L2 Constructs |
- Curated CloudFormation (Define a single AWS resource)
- Provide sensible defaults, fill empty fields with reasonable values
- Provide security best practices
- Provide helper methods which can change the state of resources much later after being initialised
- **Most Commonly Used**

#### L3 Constructs
- Pattern: Architectural Solution (Define a group of AWS resources)
- Define a very large infrustructure with minimal lines of code
- Lost control over the elements of the infrustructure, in exchange for speed of developmentk
- L3 architecture can be recreate using L1 or L2 constructs
- Use Case: already have a pre-established pattern, willing to sacrifice configuration for speed/performance
- **Only for Pre-made Architectures**

---

### Key Concepts

1. Construct Hub
- It is an open source collection of CDK constructs that are maintained by AWS which consists of pre-made architecture patterns.

2. Stacks
- A container of pieces of infrastructure or constructs
- Used to organize resources in a more intuitive way
- Usually seperated to group resources that related to an infrustructure or application

3. App
- A collection of stacks containing the constructs of the infrastructure with the resources required to deploy the application

# AWS CloudFormation (CFN)

AWS CloudFormation is a service that helps model and set up your AWS resources using infrastructure-as-code (IaC). It automates the creation, configuration, and management of AWS resources, making deployment and infrastructure management more efficient and reproducible.

---

## 1. CloudFormation Templates

A **CloudFormation template** is a declarative JSON or YAML file that describes the resources and configuration in the AWS environment. It’s the blueprint for an infrastructure.

- **Structure**:
  - **Resources**: The only mandatory section. Defines the AWS resources (like EC2 instances, VPCs, RDS databases) you want to create.
  - **Parameters**: Allows input of custom values at stack creation time (e.g., instance types, AMI IDs).
  - **Outputs**: Specifies return values (e.g., the URL of a deployed web application or instance IDs).
  - **Mappings**: Defines static mappings, such as region-specific settings.
  - **Conditions**: Used to create resources conditionally, based on parameter values or other logic.
  - **Metadata**: Provides additional information about the template.

```yaml
Resources:
  MyEC2Instance:
    Type: 'AWS::EC2::Instance'
    Properties:
      InstanceType: t2.micro
      ImageId: ami-0c55b159cbfafe1f0
```

## 2. CFN Stacks

A **CloudFormation stac** is a collection of AWS resources that manage a single unit. It create, update, or delete a collection of resources by managing a stack.
- **Creating a stack**: A stack is created from a template. When the stack is created, AWS CloudFormation provisions and configures the resources defined in the template.
- **Updating a stack**: You can modify a stack by updating its template or parameters. CloudFormation manages resource dependencies and applies only the necessary changes.
- **Deleting a stack**: When a stack is deleted, CloudFormation deletes all associated resources unless they are protected.

## 3. Resources

**Resources** are the core component of CloudFormation. A resource represents an AWS service (e.g., EC2 instance, S3 bucket, VPC, Lambda function). Each resource has a unique logical ID and defined properties.
- Type: Specifies the resource type (e.g., `AWS::EC2::Instance`, `AWS::S3::Bucket`).
- Properties: Defines the configuration of the resource (e.g., instance size, security group, VPC ID).
```yaml
Resources:
  MyBucket:
    Type: "AWS::S3::Bucket"
    Properties:
      BucketName: "my-cloudformation-bucket"
```

## 4. Parameters

**Parameters** allow you to pass in dynamic values to your templates, making them reusable and customizable. You can prompt for instance types, AMI IDs, or other inputs at stack creation time.
- Parameters can have default values, allowed values, and constraints.
- Users specify parameters when creating or updating a stack.
```yaml
Parameters:
  InstanceTypeParameter:
    Type: String
    Default: t2.micro
    AllowedValues:
      - t2.micro
      - t2.small
    Description: "EC2 instance type"

Resources:
  MyInstance:
    Type: "AWS::EC2::Instance"
    Properties:
      InstanceType: !Ref InstanceTypeParameter
```

## 5. Outputs

**Outputs** provide information about resources created by the stack. These are useful for exporting values (like resource ARNs, URLs, or IDs) that can be used by other stacks or displayed after stack creation.
- Useful for cross-stack references or sharing information between stacks.
- Outputs are visible in the AWS CloudFormation console.

```yaml
Outputs:
  InstanceID:
    Description: "The ID of the EC2 instance"
    Value: !Ref MyInstance
```

## 6. Mappings
**Mappings** allow you to define static mappings, such as region-specific settings. They provide a way to organize data based on keys and values.
- You cannot modify mappings during stack creation or updates.
- Useful for region-specific AMI IDs or environment configurations.
```yaml
Mappings:
  RegionMap:
    us-east-1:
      AMI: "ami-0ff8a91507f77f867"
    us-west-1:
      AMI: "ami-0bdb828fd58c52235"

Resources:
  MyInstance:
    Type: "AWS::EC2::Instance"
    Properties:
      InstanceType: t2.micro
      ImageId: !FindInMap [RegionMap, !Ref "AWS::Region", AMI]
```

# AWS Virtual Private Cloud (VPC)
A VPC is a logically isolated network where AWS resources can be launch.
- **Subnets**: These are segments of a VPC that allow resources to be organised.
- **Route Tables**: Define the routes for traffic within and outside of the VPC.
- **Internet Gateway (IGW)**: Enables your VPC to communicate with the internet.
- **NAT Gateway**: Allows instances in a private subnet to access the internet without being directly exposed.
- **Security Groups**: Virtual firewalls controlling inbound and outbound traffic to your instances.
- **Network ACLs**: Optional stateless filters controlling traffic at the subnet level.

## Cdk Stack Components

### VPC
- Creates a VPC with 1 Availability Zone
- Includes 1 public and 1 private subnets
- CIDR mask of 24 for each subnet
- Includes a NAT Gateway for private subnet internet access

### EC2 Instances
1. Public Instance
   - Deployed in the public subnet
   - Uses t2.micro instance type
   - Latest Amazon Linux 2 AMI
   - Includes a security group allowing inbound HTTP traffic
   - User data script installs and starts Nginx

2. Private Instance
   - Deployed in the private subnet
   - Uses t2.micro instance type
   - Latest Amazon Linux 2 AMI
   - Includes a security group allowing inbound HTTP traffic from within the VPC
   - User data script installs and starts Nginx

### Security Groups
1. Public Security Group
   - Allows inbound HTTP traffic (port 80) from any IPv4 address
   - Allows all outbound traffic

2. Private Security Group
   - Allows inbound HTTP traffic (port 80) only from within the VPC
   - Allows all outbound traffic

### User Data
Both instances use a user data script that:
- Updates the system
- Installs Nginx using Amazon Linux Extras
- Starts the Nginx service
- Enables Nginx to start on boot

## Deployment
- Configure AWS CLI credentials and enter the required fields
    ```bash
    aws configure
    ```
- Deploy the stack using the CDK CLI
    ```bash
    cdk deploy
    ```