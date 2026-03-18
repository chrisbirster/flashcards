# Slide Text Extraction

## Slide 14

Getting started with AWS

## Slide 15

AWS Cloud History2002: Internally launched
2003: Amazon infrastructure is one of their core strength. Idea to market
2004: Launched publicly with SQS
2006: Re-launched publicly with SQS, S3 & EC2
2007: Launched in Europe

## Slide 16

AWS Cloud Number Facts
- In 2023, AWS had $90 billion in annual revenue
- AWS accounts for 31% of the market in Q1 2024 (Microsoft is 2nd with 25%)
- Pioneer and Leader of the AWS Cloud Market for the 13th consecutive year
- Over 1,000,000 active usersGartner Magic Quadrant

## Slide 17

AWS Cloud Use Cases
- AWS enables you to build sophisticated, scalable applications
- Applicable to a diverse set of industries
- Use cases include
- Enterprise IT, Backup & Storage, Big Data analytics
- Website hosting, Mobile & Social Apps
- Gaming

## Slide 18

AWS Global Infrastructure
- AWS Regions
- AWS Availability Zones
- AWS Data Centers
- AWS Edge Locations / Points of Presence
- https://infrastructure.aws/

## Slide 19

AWS Regions
- AWS has Regions all around the world
- Names can be us-east-1, eu-west-3…
- A region is a cluster of data centers
- Most AWS services are region-scoped
https://aws.amazon.com/about-aws/global-infrastructure/

## Slide 20

How to choose an AWS Region?
- Compliance with data governance and legal requirements: data never leaves a region without your explicit permission
- Proximity to customers: reduced latency
- Available services within a Region: new services and new features aren’t available in every Region
- Pricing: pricing varies region to region and is transparent in the service pricing page
If you need to launch a new application, where should you do it?
?
?
?
?

## Slide 21

AWS Availability Zones
- Each region has many availability zones (usually 3, min is 3, max is 6). Example:
- ap-southeast-2a
- ap-southeast-2b
- ap-southeast-2c
- Each availability zone (AZ) is one or more discrete data centers with redundant power, networking, and connectivity
- They’re separate from each other, so that they’re isolated from disasters
- They’re connected with high bandwidth, ultra-low latency networking
AWS RegionSydney: ap-southeast-2ap-southeast-2a
ap-southeast-2bap-southeast-2c

## Slide 22

AWS Points of Presence (Edge Locations)
- Amazon has 400+ Points of Presence (400+ Edge Locations & 10+ Regional Caches) in 90+ cities across 40+ countries
- Content is delivered to end users with lower latency
https://aws.amazon.com/cloudfront/features/

## Slide 23

To u r o f t h e AW S C o n s o l e
- AWS has Global Services:
- Identity and Access Management (IAM)
- Route 53 (DNS service)
- CloudFront (Content Delivery Network)
- WAF (Web Application Firewall)
- Most AWS services are Region-scoped:
- Amazon EC2 (Infrastructure as a Service)
- Elastic Beanstalk (Platform as a Service)
- Lambda (Function as a Service)
- Rekognition (Software as a Service)
- Region Table: https://aws.amazon.com/about-aws/global-infrastructure/regional-product-services

## Slide 24

AWS Identity and Access Management (AWS IAM)

## Slide 25

IAM: Users & Groups
- IAM = Identity and Access Management, Global service
- Root account created by default, shouldn’t be used or shared
- Users are people within your organization, and can be grouped
- Groups only contain users, not other groups
- Users don’t have to belong to a group, and user can belong to multiple groups
Alice
Bob
Charles
David
Edward
Group: DevelopersGroup: OperationsGroup Audit Team
Fred

## Slide 26

IAM: Permissions
- Users or Groups can be assigned JSON documents called policies
- These policies define the permissions of the users
- In AWS you apply the least privilege principle: don’t give more permissions than a user needs
{ "Version": "2012-10-17", "Statement": [ { "Effect": "Allow", "Action": "ec2:Describe*", "Resource": "*" }, { "Effect": "Allow", "Action": "elasticloadbalancing:Describe*", "Resource": "*" }, { "Effect": "Allow", "Action": [ "cloudwatch:ListMetrics", "cloudwatch:GetMetricStatistics", "cloudwatch:Describe*" ], "Resource": "*" } ]}

## Slide 27

IAM Policies inheritance
Alice
Bob
Charles
David
Edward
Developers OperationsAudit Team
Fred
inline

## Slide 28

IAM Policies Structure
- Consists of
- Version: policy language version, always include “2012-10-17”
- Id: an identifier for the policy (optional)
- Statement: one or more individual statements (required)
- Statements consists of
- Sid: an identifier for the statement (optional)
- Effect: whether the statement allows or denies access (Allow, Deny)
- Principal: account/user/role to which this policy applied to
- Action: list of actions this policy allows or denies
- Resource: list of resources to which the actions applied to
- Condition: conditions for when this policy is in effect (optional)

## Slide 29

IAM – Password Policy
- Strong passwords = higher security for your account
- In AWS, you can setup a password policy:
- Set a minimum password length
- Require specific character types:
- including uppercase letters
- lowercase letters
- numbers
- non-alphanumeric characters
- Allow all IAM users to change their own passwords
- Require users to change their password after some time (password expiration)
- Prevent password re-use

## Slide 30

Multi Factor Authentication - MFA
- Users have access to your account and can possibly change configurations or delete resources in your AWS account
- Yo u w a n t t o p r o t e c t yo u r R o o t A c c o u n t s a n d I A M u s e r s
- MFA = password you know + security device you own
- Main benefit of MFA: if a password is stolen or hacked, the account is not compromised
Alice+Password
=> Successful login

## Slide 31

MFA devices options in AWSVirtual MFA device
Google Authenticator(phone only)
Authy(phone only)
Universal 2nd Factor (U2F) Security Key
YubiKey by Yubico (3rd party)
Support for multiple tokens on a single device.Support for multiple root and IAM users using a single security key

## Slide 32

MFA devices options in AWS
Hardware Key Fob MFA Device
Provided by Gemalto (3rd party)
Hardware Key Fob MFA Device forAWS GovCloud (US)
Provided by SurePassID (3rd party)

## Slide 33

How can users access AWS ?
- To a c c e s s AW S , yo u h ave t h r e e o p t i o n s :
- AWS Management Console (protected by password + MFA)
- AWS Command Line Interface (CLI): protected by access keys
- AWS Software Developer Kit (SDK) - for code: protected by access keys
- Access Keys are generated through the AWS Console
- Users manage their own access keys
- Access Keys are secret, just like a password. Don’t share them
- Access Key ID ~= username
- Secret Access Key ~= password

## Slide 34

Example (Fake) Access Keys
- Access key ID: AKIASK4E37PV4983d6C
- Secret Access Key: AZPN3zojWozWCndIjhB0Unh8239a1bzbzO5fqqkZq
- Remember: don’t share your access keys

## Slide 35

What’s the AWS CLI?
- A tool that enables you to interact with AWS services using commands in your command-line shell
- Direct access to the public APIs of AWS services
- Yo u c a n d e ve l o p s c r i p t s t o m a n a g e yo u r r e s o u r c e s
- It’s open-source https://github.com/aws/aws-cli
- Alternative to using AWS Management Console

## Slide 36

What’s the AWS SDK?
- AWS Software Development Kit (AWS SDK)
- Language-specific APIs (set of libraries)
- Enables you to access and manage AWS services programmatically
- Embedded within your application
- Supports
- SDKs (JavaScript, Python, PHP , .NET, Ruby, Java, Go, Node.js, C++)
- Mobile SDKs (Android, iOS, …)
- IoT Device SDKs (Embedded C, Arduino, …)
- Example: AWS CLI is built on AWS SDK for Python
AWS SDK
Your Application

## Slide 37

IAM Roles for Services
- Some AWS service will need to perform actions on your behalf
- To d o s o, we w i l l a s s i g n permissions to AWS services with IAM Roles
- Common roles:
- EC2 Instance Roles
- Lambda Function Roles
- Roles for CloudFormation
EC2 Instance(virtual server)
IAM Role
Access AWS

## Slide 38

IAM Security T ools
- IAM Credentials Report (account-level)
- a report that lists all your account's users and the status of their various credentials
- IAM Access Advisor (user-level)
- Access advisor shows the service permissions granted to a user and when those services were last accessed.
- Yo u c a n u s e t h i s i n fo r m a t i o n t o r e v i s e yo u r p o l i c i e s .

## Slide 39

IAM Guidelines & Best Practices
- Don’t use the root account except for AWS account setup
- One physical user = One AWS user
- Assign users to groups and assign permissions to groups
- Create a strong password policy
- Use and enforce the use of Multi Factor Authentication (MFA)
- Create and use Roles for giving permissions to AWS services
- Use Access Keys for Programmatic Access (CLI / SDK)
- Audit permissions of your account using IAM Credentials Report & IAM Access Advisor
- Never share IAM users & Access Keys

## Slide 40

IAM Section – Summary
- Users: mapped to a physical user, has a password for AWS Console
- Groups: contains users only
- Policies: JSON document that outlines permissions for users or groups
- Roles: for EC2 instances or AWS services
- Security: MFA + Password Policy
- AWS CLI: manage your AWS services using the command-line
- AWS SDK: manage your AWS services using a programming language
- Access Keys: access AWS using the CLI or SDK
- Audit: IAM Credential Reports & IAM Access Advisor

## Slide 41

Amazon EC2 – Basics

## Slide 42

Amazon EC2
- EC2 is one of the most popular of AWS’ offering
- EC2 = Elastic Compute Cloud = Infrastructure as a Service
- It mainly consists in the capability of :
- Renting virtual machines (EC2)
- Storing data on virtual drives (EBS)
- Distributing load across machines (ELB)
- Scaling the services using an auto-scaling group (ASG)
- Knowing EC2 is fundamental to understand how the Cloud works

## Slide 43

EC2 sizing & configuration options
- Operating System (OS): Linux, Windows or Mac OS
- How much compute power & cores (CPU)
- How much random-access memory (RAM)
- How much storage space:
- Network-attached (EBS & EFS)
- hardware (EC2 Instance Store)
- Network card: speed of the card, Public IP address
- Firewall rules: security group
- Bootstrap script (configure at first launch): EC2 User Data

## Slide 44

EC2 User Data
- It is possible to bootstrap our instances using an EC2 User data script.
- bootstrapping means launching commands when a machine starts
- That script is only run once at the instance first start
- EC2 user data is used to automate boot tasks such as:
- Installing updates
- Installing software
- Downloading common files from the internet
- Anything you can think of
- The EC2 User Data Script runs with the root user

## Slide 45

Hands-On:Launching an EC2 Instance running Linux
- We’ll be launching our first vir tual ser ver using the AWS Console
- We’ll get a first high-level approach to the various parameters
- We’ll see that our web server is launched using EC2 user data
- We’ll learn how to star t / stop / terminate our instance.

## Slide 46

EC2 Instance Types - Overview
- Y ou can use different types of EC2 instances that are optimised for different use cases (https://aws.amazon.com/ec2/instance-types/)
- AWS has the following naming convention:m5.2xlarge
- m: instance class
- 5: generation (AWS improves them over time)
- 2xlarge: size within the instance class

## Slide 47

EC2 Instance Types – General Purpose
- Great for a diversity of workloads such as web servers or code repositories
- Balance between:
- Compute
- Memory
- Networking
- In the course, we will be using the t2.micro which is a General Purpose EC2 instance
* this list will evolve over time, please check the AWS website for the latest information

## Slide 48

EC2 Instance Types – Compute Optimized
- Great for compute-intensive tasks that require high performance processors:
- Batch processing workloads
- Media transcoding
- High performance web servers
- High performance computing (HPC)
- Scientific modeling & machine learning
- Dedicated gaming servers
* this list will evolve over time, please check the AWS website for the latest information

## Slide 49

EC2 Instance Types – Memory Optimized
- Fast performance for workloads that process large data sets in memory
- Use cases:
- High performance, relational/non-relational databases
- Distributed web scale cache stores
- In-memory databases optimized for BI (business intelligence)
- Applications performing real-time processing of big unstructured data
* this list will evolve over time, please check the AWS website for the latest information

## Slide 50

EC2 Instance Types – Storage Optimized
- Great for storage-intensive tasks that require high, sequential read and write access to large data sets on local storage
- Use cases:
- High frequency online transaction processing (OLTP) systems
- Relational & NoSQL databases
- Cache for in-memory databases (for example, Redis)
- Data warehousing applications
- Distributed file systems
* this list will evolve over time, please check the AWS website for the latest information

## Slide 51

EC2 Instance Types: exampleInstancevCPUMem (GiB)StorageNetwork PerformanceEBS Bandwidth (Mbps)t2.micro11EBS-OnlyLow to Moderatet2.xlarge416EBS-OnlyModeratec5d.4xlarge16321 x 400 NVMe SSDUp to 10 Gbps4,750r5.16xlarge64512EBS Only20 Gbps13,600m5.8xlarge32128EBS Only10 Gbps6,800
Great website: https://instances.vantage.sh

## Slide 52

Introduction to Security Groups
- Security Groups are the fundamental of network security in AWS
- They control how traffic is allowed into or out of our EC2 Instances.
- Security groups only contain rules
- Security groups rules can reference by IP or by security group
Inbound trafficOutbound traffic
Security GroupWWW
EC2 Instance

## Slide 53

Security GroupsDeeper Dive
- Security groups are acting as a “firewall” on EC2 instances
- They regulate:
- Access to Ports
- Authorised IP ranges – IPv4 and IPv6
- Control of inbound network (from other to the instance)
- Control of outbound network (from the instance to other)

## Slide 54

Security GroupsDiagram
EC2 InstanceIP XX.XX.XX.XX
Port 22Port 22Security Group 1InboundFilter IP / Port with Rules
Any PortSecurity Group 1OutboundFilter IP / Port with Rules
Your Computer - IP XX.XX.XX.XX(authorised port 22)Other computer(not authorised port 22)
WWWAny IP – Any Port

## Slide 55

Security GroupsGood to know
- Can be attached to multiple instances
- Locked down to a region / VPC combination
- Does live “outside” the EC2 – if traffic is blocked the EC2 instance won’t see it
- It’s good to maintain one separate security group for SSH access
- If your application is not accessible (time out), then it’s a security group issue
- If your application gives a “connection refused“ error, then it’s an application error or it’s not launched
- All inbound traffic is blocked by default
- All outbound traffic is authorised by default

## Slide 56

Referencing other security groupsDiagram
EC2 InstanceIP XX.XX.XX.XX
Port 123
Port 123
Port 123
Security Group 1InboundAuthorising Security Group 1Authorising Security Group 2
EC2 InstanceIP XX.XX.XX.XXSecurity Group 2 (attached)
EC2 InstanceIP XX.XX.XX.XXSecurity Group 1 (attached)
EC2 InstanceIP XX.XX.XX.XXSecurity Group 3 (attached)

## Slide 57

Classic Ports to know
- 22 = SSH (Secure Shell) - log into a Linux instance
- 21 = FTP (File Transfer Protocol) – upload files into a file share
- 22 = SFTP (Secure File Transfer Protocol) – upload files using SSH
- 80 = HTTP – access unsecured websites
- 443 = HTTPS – access secured websites
- 3389 = RDP (Remote Desktop Protocol) – log into a Windows instance

## Slide 58

SSH Summary TableMac
Linux
Windows < 10
Windows >= 10
SSHPuttyEC2 Instance Connect

## Slide 59

Which Lectures to watch
- Mac / Linux:
- SSH on Mac/Linux lecture
- Windows:
- Putty Lecture
- If Windows 10: SSH on Windows 10 lecture
- All:
- EC2 Instance Connect lecture

## Slide 60

SSH troubleshooting
- Students have the most problems with SSH
- If things don’t work…1.Re-watch the lecture. Y ou may have missed something2.Read the troubleshooting guide3.Try EC2 Instance Connect
- If one method works (SSH, Putty or EC2 Instance Connect) you’re good
- If no method works, that’s okay, the course won’t use SSH much

## Slide 61

How to SSH into your EC2 InstanceLinux / Mac OS X
- We’ll learn how to SSH into your EC2 instance using Linux / Mac
- SSH is one of the most important function. It allows you to control a remote machine, all using the command line.
- We will see how we can configure OpenSSH ~/.ssh/config to facilitate the SSH into our EC2 instances
EC2 InstanceLinuxPublic IPSSH – Port 22
WWW

## Slide 62

How to SSH into your EC2 InstanceWindows
- We’ll learn how to SSH into your EC2 instance using Windows
- SSH is one of the most important function. It allows you to control a remote machine, all using the command line.
- We will configure all the required parameters necessar y for doing SSH on Windows using the free tool Putty.
EC2 InstanceLinuxPublic IPSSH – Port 22
WWW

## Slide 63

EC2 Instance Connect
- Connect to your EC2 instance within your browser
- No need to use your key file that was downloaded
- The “magic” is that a temporary key is uploaded onto EC2 by AWS
- Works only out-of-the-box with Amazon Linux 2
- Need to make sure the port 22 is still opened!

## Slide 64

EC2 Instances Purchasing Options
- On-Demand Instances – short workload, predictable pricing, pay by second
- Reserved (1 & 3 years)
- Reserved Instances – long workloads
- Convertible Reserved Instances – long workloads with flexible instances
- Savings Plans (1 & 3 years) –commitment to an amount of usage, long workload
- Spot Instances – short workloads, cheap, can lose instances (less reliable)
- Dedicated Hosts – book an entire physical server, control instance placement
- Dedicated Instances – no other customers will share your hardware
- Capacity Reservations – reserve capacity in a specific AZ for any duration

## Slide 65

EC2 On Demand
- Pay for what you use:
- Linux or Windows - billing per second, after the first minute
- All other operating systems - billing per hour
- Has the highest cost but no upfront payment
- No long-term commitment
- Recommended for short-term and un-interrupted workloads, where you can't predict how the application will behave

## Slide 66

EC2 Reserved Instances
- Up to 72% discount compared to On-demand
- Yo u r e s e r ve a s p e c i fi c i n s t a n c e a t t r i bu t e s (Instance Type, Region, Tenancy, OS)
- Reservation Period – 1 year (+discount) or 3 years (+++discount)
- Payment Options – No Upfront (+), Partial Upfront (++), All Upfront (+++)
- Reserved Instance’s Scope – Regional or Zonal (reserve capacity in an AZ)
- Recommended for steady-state usage applications (think database)
- Yo u c a n bu y a n d s e l l i n t h e R e s e r ve d I n s t a n c e M a r k e t p l a c e
- Convertible Reserved Instance
- Can change the EC2 instance type, instance family, OS, scope and tenancy
- Up to 66% discount
Note: the % discounts are different from the video as AWS change them over time – the exact numbers are not needed for the exam. This is just for illustrative purposes J

## Slide 67

EC2 Savings Plans
- Get a discount based on long-term usage (up to 72% - same as RIs)
- Commit to a certain type of usage ($10/hour for 1 or 3 years)
- Usage beyond EC2 Savings Plans is billed at the On-Demand price
- Locked to a specific instance family & AWS region (e.g., M5 in us-east-1)
- Flexible across:
- Instance Size (e.g., m5.xlarge, m5.2xlarge)
- OS (e.g., Linux, Windows)
- Tenancy (Host, Dedicated, Default)

## Slide 68

EC2 Spot Instances
- Can get a discount of up to 90% compared to On-demand
- Instances that you can “lose” at any point of time if your max price is less than the current spot price
- The MOST cost-efficient instances in AWS
- Useful for workloads that are resilient to failure
- Batch jobs
- Data analysis
- Image processing
- Any distributed workloads
- Workloads with a flexible star t and end time
- Not suitable for critical jobs or databases

## Slide 69

EC2 Dedicated Hosts
- A physical server with EC2 instance capacity fully dedicated to your use
- Allows you address compliance requirements and use your existing server-bound software licenses (per-socket, per-core, pe—VM software licenses)
- Purchasing Options:
- On-demand – pay per second for active Dedicated Host
- Reserved - 1 or 3 years (No Upfront, Partial Upfront, All Upfront)
- The most expensive option
- Useful for software that have complicated licensing model (BYOL – Bring Y our Own License)
- Or for companies that have strong regulatory or compliance needs

## Slide 70

EC2 Dedicated Instances
- Instances run on hardware that’s dedicated to you
- May share hardware with other instances in same account
- No control over instance placement (can move hardware after Stop / Start)

## Slide 71

EC2 Capacity Reservations
- Reserve On-Demand instances capacity in a specific AZ for any duration
- Y ou always have access to EC2 capacity when you need it
- No time commitment (create/cancel anytime), no billing discounts
- Combine with Regional Reserved Instances and Savings Plans to benefit from billing discounts
- Yo u ’r e c h a r g e d a t O n-Demand rate whether you run instances or not
- Suitable for short-term, uninterrupted workloads that needs to be in a specific AZ

## Slide 72

Which purchasing option is right for me?
- On demand: coming and staying in resort whenever we like, we pay the full price
- Reserved: like planning ahead and if we plan to stay for a long time, we may get a good discount.
- Savings Plans: pay a certain amount per hour for certain period and stay in any room type (e.g., King, Suite, Sea View, …)
- Spot instances: the hotel allows people to bid for the empty rooms and the highest bidder keeps the rooms. Y ou can get kicked out at any time
- Dedicated Hosts: We book an entire building of the resort
- Capacity Reservations: you book a room for a period with full price even you don’t stay in it

## Slide 73

Price ComparisonExample – m4.large – us-east-1Price TypePrice (per hour)On-Demand $0.10 Spot Instance (Spot Price)$0.038 - $0.039 (up to 61% off)Reserved Instance (1 year)$0.062 (No Upfront) - $0.058 (All Upfront)Reserved Instance (3 years)$0.043 (No Upfront) - $0.037 (All Upfront)EC2 Savings Plan (1 year)$0.062 (No Upfront) - $0.058 (All Upfront)Reserved Convertible Instance (1 year)$0.071 (No Upfront) - $0.066 (All Upfront)Dedicated HostOn-Demand PriceDedicated Host ReservationUp to 70% offCapacity ReservationsOn-Demand Price

## Slide 74

EC2 Spot Instance Requests
- Can get a discount of up to 90% compared to On-demand
- Define max spot price and get the instance while current spot price < max
- The hourly spot price varies based on offer and capacity
- If the current spot price > your max price you can choose to stop or terminate your instance with a 2 minutes grace period.
- Other strategy: Spot Block
- “block” spot instance during a specified time frame (1 to 6 hours) without interruptions
- In rare situations, the instance may be reclaimed
- Used for batch jobs, data analysis, or workloads that are resilient to failures.
- Not great for critical jobs or databases

## Slide 75

EC2 Spot Instances Pricing
https://console.aws.amazon.com/ec2sp/v1/spot/home?region=us-east-1#
User-defined max price

## Slide 76

How to terminate Spot Instances?
https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/spot-requests.html
You can only cancel Spot Instance requests that areopen,active, ordisabled. Cancelling a Spot Request does not terminate instancesYou must first cancel a Spot Request, and then terminate the associated Spot Instances

## Slide 77

Spot Fleets
- Spot Fleets = set of Spot Instances + (optional) On-Demand Instances
- The Spot Fleet will try to meet the target capacity with price constraints
- Define possible launch pools: instance type (m5.large), OS, Availability Zone
- Can have multiple launch pools, so that the fleet can choose
- Spot Fleet stops launching instances when reaching capacity or max cost
- Strategies to allocate Spot Instances:
- lowestPrice: from the pool with the lowest price (cost optimization, short workload)
- diversified: distributed across all pools (great for availability, long workloads)
- capacityOptimized: pool with the optimal capacity for the number of instances
- priceCapacityOptimized (recommended): pools with highest capacity available, then select the pool with the lowest price (best choice for most workloads)
- Spot Fleets allow us to automatically request Spot Instances with the lowest price

## Slide 78

Amazon EC2 – Associate

## Slide 79

Private vs Public IP (IPv4)
- Networking has two sorts of IPs. IPv4 and IPv6:
- IPv4: 1.160.10.240
- IPv6: 3ffe:1900:4545:3:200:f8ff:fe21:67cf
- In this course, we will only be using IPv4.
- IPv4 is still the most common format used online.
- IPv6 is newer and solves problems for the Internet of Things (IoT).
- IPv4 allows for 3.7 billion different addresses in the public space
- IPv4: [0-255].[0-255].[0-255].[0-255].

## Slide 80

Company BPrivate Network192.168.0.1/22Company APrivate Network192.168.0.1/22
Private vs Public IP (IPv4)Example
Internet Gateway (public):149.140.72.10
Web Server (public):79.216.59.75Server (public):211.139.37.43
Internet Gateway (public):253.144.139.205WWW

## Slide 81

Private vs Public IP (IPv4)Fundamental Differences
- Public IP:
- Public IP means the machine can be identified on the internet (WWW)
- Must be unique across the whole web (not two machines can have the same public IP).
- Can be geo-located easily
- Private IP:
- Private IP means the machine can only be identified on a private network only
- The IP must be unique across the private network
- BUT two different private networks (two companies) can have the same IPs.
- Machines connect to WWW using a NAT + internet gateway (a proxy)
- Only a specified range of IPs can be used as private IP

## Slide 82

Elastic IPs
- When you stop and then start an EC2 instance, it can change its public IP .
- If you need to have a fixed public IP for your instance, you need an Elastic IP
- An Elastic IP is a public IPv4 IP you own as long as you don’t delete it
- Yo u c a n a t t a c h i t t o o n e i n s t a n c e a t a t i m e

## Slide 83

Elastic IP
- With an Elastic IPaddress, you can mask the failure of an instance or software by rapidly remapping the address to another instance in your account.
- Y ou can only have 5 Elastic IP in your account (you can ask AWS to increase that).
- Overall, try to avoid using Elastic IP:
- They often reflect poor architectural decisions
- Instead, use a random public IP and register a DNS name to it
- Or, as we’ll see later, use a Load Balancer and don’t use a public IP

## Slide 84

Private vs Public IP (IPv4)In AWS EC2 – Hands On
- By default, your EC2 machine comes with:
- A private IP for the internal AWS Network
- A public IP , for the WWW.
- When we are doing SSH into our EC2 machines:
- We can’t use a private IP, because we are not in the same network
- We can only use the public IP.
- If your machine is stopped and then started, the public IP can change

## Slide 85

Placement Groups
- Sometimes you want control over the EC2 Instance placement strategy
- That strategy can be defined using placement groups
- When you create a placement group, you specify one of the following strategies for the group:
- Cluster—clusters instances into a low-latency group in a single Availability Zone
- Spread—spreads instances across underlying hardware (max 7 instances per group per AZ)
- Partition—spreads instances across many different partitions (which rely on different sets of racks) within an AZ. Scales to 100s of EC2 instances per group (Hadoop, Cassandra, Kafka)

## Slide 86

Same AZ
Placement GroupsCluster
- Pros: Great network (10 Gbps bandwidth between instances with Enhanced Networking enabled - recommended)
- Cons: If the AZ fails, all instances fails at the same time
- Use case:
- Big Data job that needs to complete fast
- Application that needs extremely low latency and high network throughput
EC2EC2
EC2EC2EC2
EC2 Placement group ClusterLow latency10 Gbps network

## Slide 87

Us-east-1a
Hardware 1
Placement GroupsSpread
- Pros:
- Can span across Availability Zones (AZ)
- Reduced risk is simultaneous failure
- EC2 Instances are on different physical hardware
- Cons:
- Limited to 7 instances per AZ per placement group
- Use case:
- Application that needs to maximize high availability
- Critical Applications where each instance must be isolated from failure from each other
EC2
Hardware 2EC2
Us-east-1b
Hardware 3EC2
Hardware 4EC2
Us-east-1c
Hardware 5EC2
Hardware 6EC2

## Slide 88

Placements GroupsPartition
- Up to 7 partitions per AZ
- Can span across multiple AZs in the same region
- Up to 100s of EC2 instances
- The instances in a partition do not share racks with the instances in the other partitions
- A partition failure can affect many EC2 but won’t affect other partitions
- EC2 instances get access to the partition information as metadata
- Use cases: HDFS, HBase, Cassandra, Kafka
us-east-1a
Partition 1
EC2
EC2
EC2
EC2Partition 2
EC2
EC2
EC2
EC2Partition 3
EC2
EC2
EC2
EC2
us-east-1b

## Slide 89

Elastic Network Interfaces (ENI)
- Logical component in a VPC that represents a virtual network card
- The ENI can have the following attributes:
- Primary private IPv4, one or more secondary IPv4
- One Elastic IP (IPv4) per private IPv4
- One Public IPv4
- One or more security groups
- A MAC address
- Yo u c a n c r e a t e E N I i n d e p e n d e n t ly a n d a t t a c h them on the fly (move them) on EC2 instances for failover
- Bound to a specific availability zone (AZ)
EC2
Eth0 – primary ENI192.168.0.31
Eth1 – secondary ENI192.168.0.42
EC2
Eth0 – primary ENI
Can be moved
Availability Zone

## Slide 90

EC2 Hibernate
- We know we can stop, terminate instances
- Stop – the data on disk (EBS) is kept intact in the next start
- Ter minate – any EBS volumes (root) also set-up to be destroyed is lost
- On start, the following happens:
- First start: the OS boots & the EC2 User Data script is run
- Following starts: the OS boots up
- Then your application starts, caches get warmed up, and that can take time!

## Slide 91

EC2 Hibernate
- Introducing EC2 Hibernate:
- The in-memory (RAM) state is preserved
- The instance boot is much faster! (the OS is not stopped / restarted)
- Under the hood: the RAM state is written to a file in the root EBS volume
- The root EBS volume must be encrypted
- Use cases:
- Long-running processing
- Saving the RAM state
- Services that take time to initialize
RAM
EC2 InstanceRoot EBS Volume(Encrypted)
RAM
RAM
RAM
RAM
Running
Stopping
Stopped
Running
Hibernate
Shutdown
Start
Hibernation

## Slide 92

EC2 Hibernate – Good to know
- Supported Instance Families – C3, C4, C5, I3, M3, M4, R3, R4, T2, T3, …
- Instance RAM Size – must beless than 150 GB.
- Instance Size – not supported for bare metal instances.
- AMI – Amazon Linux 2, Linux AMI, Ubuntu, RHEL, CentOS & Windows…
- Root Volume – must be EBS, encrypted, not instance store, and large
- Available for On-Demand, Reserved and Spot Instances
- An instance can NOT be hibernated more than 60 days

## Slide 93

Amazon EC2 – Instance Storage

## Slide 94

What’s an EBS Volume?
- An EBS (Elastic Block Store) Volume is a network drive you can attach to your instances while they run
- It allows your instances to persist data, even after their termination
- They can only be mounted to one instance at a time (at the CCP level)
- They are bound to a specific availability zone
- Analogy: Think of them as a “network USB stick”

## Slide 95

EBS Volume
- It’s a network drive (i.e. not a physical drive)
- It uses the network to communicate the instance, which means there might be a bit of latency
- It can be detached from an EC2 instance and attached to another one quickly
- It’s locked to an Availability Zone (AZ)
- An EBS Volume in us-east-1a cannot be attached to us-east-1b
- To m ove a vo l u m e a c r o s s , yo u fi r s t n e e d t o s n a p s h o t i t
- Have a provisioned capacity (size in GBs, and IOPS)
- Yo u g e t b i l l e d fo r a l l t h e p r ov i s i o n e d c a p a c i t y
- Yo u c a n i n c r e a s e t h e c a p a c i t y o f t h e d r i ve ove r t i m e

## Slide 96

US-EAST-1BUS-EAST-1AEBS Volume - Example
EBS(10 GB)
EBS(100 GB)
EBS(50 GB)
EBS(50 GB)
EBS(10 GB)unattached

## Slide 97

EBS – Delete on Termination attribute
- Controls the EBS behaviour when an EC2 instance terminates
- By default, the root EBS volume is deleted (attribute enabled)
- By default, any other attached EBS volume is not deleted (attribute disabled)
- This can be controlled by the AWS console / AWS CLI
- Use case: preserve root volume when instance is terminated

## Slide 98

EBS Snapshots
- Make a backup (snapshot) of your EBS volume at a point in time
- Not necessary to detach volume to do snapshot, but recommended
- Can copy snapshots across AZ or RegionUS-EAST-1A
EBS(50 GB)
US-EAST-1B
EBS(50 GB)
EBS Snapshotsnapshotrestore

## Slide 99

EBS Snapshots Features
- EBS Snapshot Archive
- Move a Snapshot to an ”archive tier” that is 75% cheaper
- Takes within 24 to 72 hours for restoring the archive
- Recycle Bin for EBS Snapshots
- Setup rules to retain deleted snapshots so you can recover them after an accidental deletion
- Specify retention (from 1 day to 1 year)
- Fast Snapshot Restore (FSR)
- Force full initialization of snapshot to have no latency on the first use ($$$)
EBS Snapshot
EBS SnapshotArchive
archive
EBS SnapshotRecycle Bindelete

## Slide 100

AMI Overview
- AMI = Amazon Machine Image
- AMI are a customization of an EC2 instance
- Yo u a d d yo u r ow n s o f t w a r e , c o n fi g u r a t i o n , o p e r a t i n g s y s t e m , m o n i t o r i n g …
- Faster boot / configuration time because all your software is pre-packaged
- AMI are built for a specific region (and can be copied across regions)
- Y ou can launch EC2 instances from:
- A Public AMI: AWS provided
- Yo u r ow n A M I : you make and maintain them yourself
- An AWS Marketplace AMI: an AMI someone else made (and potentially sells)

## Slide 101

AMI Process (from an EC2 instance)
- Start an EC2 instance and customize it
- Stop the instance (for data integrity)
- Build an AMI – this will also create EBS snapshots
- Launch instances from other AMIsUS-EAST-1A
US-EAST-1B
Custom AMI
Create AMILaunch from AMI

## Slide 102

EC2 Instance Store
- EBS volumes are network drives with good but “limited” performance
- If you need a high-performance hardware disk, use EC2 Instance Store
- Better I/O performance
- EC2 Instance Store lose their storage if they’re stopped (ephemeral)
- Good for buffer / cache / scratch data / temporary content
- Risk of data loss if hardware fails
- Backups and Replication are your responsibility

## Slide 103

Local EC2 Instance Store
Very high IOPS

## Slide 104

EBS Volume Types
- EBS Volumes come in 6 types
- gp2 / gp3 (SSD): General purpose SSD volume that balances price and performance for a wide variety of workloads
- io1 / io2 Block Express (SSD): Highest-performance SSD volume for mission-critical low-latency or high-throughput workloads
- st1 (HDD): Low cost HDD volume designed for frequently accessed, throughput-intensive workloads
- sc1 (HDD): Lowest cost HDD volume designed for less frequently accessed workloads
- EBS Volumes are characterized in Size | Throughput | IOPS (I/O Ops Per Sec)
- When in doubt always consult the AWS documentation – it’s good!
- Only gp2/gp3 and io1/io2 Block Express can be used as boot volumes

## Slide 105

EBS Volume Types Use casesGeneral Purpose SSD
- Cost effective storage, low-latency
- System boot volumes, Virtual desktops, Development and test environments
- 1 GiB - 16 TiB
- gp3:
- Baseline of 3,000 IOPS and throughput of 125 MiB/s
- Can increase IOPS up to 16,000 and throughput up to 1000 MiB/s independently
- gp2:
- Small gp2 volumes can burst IOPS to 3,000
- Size of the volume and IOPS are linked, max IOPS is 16,000
- 3 IOPS per GiB, means at 5,334 GiB we are at the max IOPS

## Slide 106

EBS Volume Types Use casesProvisioned IOPS (PIOPS) SSD
- Critical business applications with sustained IOPS performance
- Or applications that need more than 16,000 IOPS
- Great for databases workloads (sensitive to storage perf and consistency)
- io1 (4 GiB - 16 TiB):
- Max PIOPS: 64,000 for Nitro EC2 instances & 32,000 for other
- Can increase PIOPS independently from storage size
- io2 Block Express (4 GiB – 64 TiB):
- Sub-millisecond latency
- Max PIOPS: 256,000 with an IOPS:GiB ratio of 1,000:1
- Supports EBS Multi-attach

## Slide 107

EBS Volume Types Use casesHard Disk Drives (HDD)
- Cannot be a boot volume
- 125 GiB to 16 TiB
- Throughput Optimized HDD (st1)
- Big Data, Data Warehouses, Log Processing
- Max throughput 500 MiB/s – max IOPS 500
- Cold HDD (sc1):
- For data that is infrequently accessed
- Scenarios where lowest cost is important
- Max throughput 250 MiB/s – max IOPS 250

## Slide 108

EBS – Volume Types Summar y
https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html#solid-state-drives

## Slide 109

EBS Multi-Attach – io1/io2 family
- Attach the same EBS volume to multiple EC2 instances in the same AZ
- Each instance has full read & write permissions to the high-performance volume
- Use case:
- Achieve higher application availability in clustered Linux applications (ex: Teradata)
- Applications must manage concurrent write operations
- Up to 16 EC2 Instances at a time
- Must use a file system that’s cluster-aware (not XFS, EXT4, etc…)
io2 volume with Multi-Attach
Availability Zone 1

## Slide 110

EBS Encryption
- When you create an encrypted EBS volume, you get the following:
- Data at rest is encrypted inside the volume
- All the data in flight moving between the instance and the volume is encrypted
- All snapshots are encrypted
- All volumes created from the snapshot
- Encryption and decryption are handled transparently (you have nothing to do)
- Encryption has a minimal impact on latency
- EBS Encryption leverages keys from KMS (AES-256)
- Copying an unencrypted snapshot allows encryption
- Snapshots of encrypted volumes are encrypted

## Slide 111

Encryption: encrypt an unencrypted EBS volume
- Create an EBS snapshot of the volume
- Encrypt the EBS snapshot ( using copy )
- Create new ebs volume from the snapshot ( the volume will also be encrypted )
- Now you can attach the encrypted volume to the original instance

## Slide 112

Amazon EFS – Elastic File System
- Managed NFS (network file system) that can be mounted on many EC2
- EFS works with EC2 instances in multi-AZ
- Highly available, scalable, expensive (3x gp2), pay per use us-east-1a
EC2 Instancesus-east-1b
EC2 Instancesus-east-1c
EC2 Instances
Security Group
EFS FileSystem

## Slide 113

Amazon EFS – Elastic File System
- Use cases: content management, web serving, data sharing, Wordpress
- Uses NFSv4.1 protocol
- Uses security group to control access to EFS
- Compatible with Linux based AMI (not Windows)
- Encryption at rest using KMS
- POSIX file system (~Linux) that has a standard file API
- File system scales automatically, pay-per-use, no capacity planning!

## Slide 114

EFS – Performance & Storage Classes
- EFS Scale
- 1000s of concurrent NFS clients, 10 GB+ /s throughput
- Grow to Petabyte-scale network file system, automatically
- Performance Mode (set at EFS creation time)
- General Purpose (default) – latency-sensitive use cases (web server, CMS, etc…)
- Max I/O – higher latency, throughput, highly parallel (big data, media processing)
- Throughput Mode
- Bursting – 1 TB = 50MiB/s + burst of up to 100MiB/s
- Provisioned – set your throughput regardless of storage size, ex: 1 GiB/s for 1 TB storage
- Elastic – automatically scales throughput up or down based on your workloads
- Up to 3GiB/s for reads and 1GiB/s for writes
- Used for unpredictable workloads

## Slide 115

EFS – Storage Classes
- Storage Tiers (lifecycle management feature – move file after N days)
- Standard: for frequently accessed files
- Infrequent access (EFS-IA): cost to retrieve files, lower price to store.
- Archive: rarely accessed data (few times each year), 50% cheaper
- Implement lifecycle policies to move files between storage tiers
- Availability and durability
- Standard: Multi-AZ, great for prod
- One Zone: One AZ, great for dev, backup enabled by default, compatible with IA (EFS One Zone-IA)
- Over 90% in cost savings
Amazon EFS File System
EFS Standard
EFS IA
Lifecycle Policy
no access for 60 days
move

## Slide 116

EBS vs EFS – Elastic Block Storage
- EBS volumes…
- one instance (except multi-attach io1/io2)
- are locked at the Availability Zone (AZ) level
- gp2: IO increases if the disk size increases
- gp3 & io1: can increase IO independently
- To m i g r a t e a n E B S vo l u m e a c r o s s A Z
- Take a snapshot
- Restore the snapshot to another AZ
- EBS backups use IO and you shouldn’t run them while your application is handling a lot of traffic
- Root EBS Volumes of instances get terminated by default if the EC2 instance gets terminated. (you can disable that)
snapshot
restore
Availability Zone 1
EBS Snapshot
Availability Zone 2
EBSEBS

## Slide 117

EBS vs EFS – Elastic File System
- Mounting 100s of instances across AZ
- EFS share website files (WordPress)
- Only for Linux Instances (POSIX)
- EFS has a higher price point than EBS
- Can leverage Storage Tiers for cost savings
- Remember: EFS vs EBS vs Instance Store
Availability Zone 1
EFS
Availability Zone 2
EFSMountTarget
LinuxLinux
EFSMountTarget

## Slide 118

High Availability & Scalability

## Slide 119

Scalability & High Availability
- Scalability means that an application / system can handle greater loads by adapting.
- There are two kinds of scalability:
- Ver tical Scalability
- Horizontal Scalability (= elasticity)
- Scalability is linked but different to High Availability
- Let’s deep dive into the distinction, using a call center as an example

## Slide 120

Ver tical Scalability
- Ver tically scalability means increasing the size of the instance
- For example, your application runs on a t2.micro
- Scaling that application vertically means running it on a t2.large
- Ver tical scalability is ver y common for non distributed systems, such as a database.
- RDS, ElastiCache are services that can scale vertically.
- There’s usually a limit to how much you can vertically scale (hardware limit)
junior operatorsenior operator

## Slide 121

Horizontal Scalability
- Horizontal Scalability means increasing the number of instances / systems for your application
- Horizontal scaling implies distributed systems.
- This is very common for web applications / modern applications
- It’s easy to horizontally scale thanks the cloud offerings such as Amazon EC2
operatoroperatoroperator
operatoroperatoroperator

## Slide 122

second building in San Francisco
first building in New YorkHigh Availability
- High Availability usually goes hand in hand with horizontal scaling
- High availability means running your application / system in at least 2 data centers (== Availability Zones)
- The goal of high availability is to survive a data center loss
- The high availability can be passive (for RDS Multi AZ for example)
- The high availability can be active (for horizontal scaling)

## Slide 123

High Availability & Scalability For EC2
- Ver tical Scaling: Increase instance size (= scale up / down)
- From: t2.nano - 0.5G of RAM, 1 vCPU
- To : u-12tb1.metal – 12.3 TB of RAM, 448 vCPUs
- Horizontal Scaling: Increase number of instances (= scale out / in)
- Auto Scaling Group
- Load Balancer
- High Availability: Run instances for the same application across multi AZ
- Auto Scaling Group multi AZ
- Load Balancer multi AZ

## Slide 124

What is load balancing?
Elastic Load Balancer
EC2 Instance
EC2 Instance
EC2 Instance
- Load Balances are servers that forward traffic to multiple servers (e.g., EC2 instances) downstream

## Slide 125

Why use a load balancer?
- Spread load across multiple downstream instances
- Expose a single point of access (DNS) to your application
- Seamlessly handle failures of downstream instances
- Do regular health checks to your instances
- Provide SSL termination (HTTPS) for your websites
- Enforce stickiness with cookies
- High availability across zones
- Separate public traffic from private traffic

## Slide 126

Why use an Elastic Load Balancer?
- An Elastic Load Balancer is a managed load balancer
- AWS guarantees that it will be working
- AWS takes care of upgrades, maintenance, high availability
- AWS provides only a few configuration knobs
- It costs less to setup your own load balancer but it will be a lot more effort on your end
- It is integrated with many AWS offerings / services
- EC2, EC2 Auto Scaling Groups, Amazon ECS
- AWS Cer tificate Manager (ACM), CloudWatch
- Route 53, AWS WAF, AWS Global Accelerator

## Slide 127

Health Checks
- Health Checks are crucial for Load Balancers
- They enable the load balancer to know if instances it forwards traffic to are available to reply to requests
- The health check is done on a port and a route (/health is common)
- If the response is not 200 (OK), then the instance is unhealthy
Elastic Load BalancerEC2 Instance
Health ChecksProtocol: HTTPPort: 4567Endpoint: /health

## Slide 128

Types of load balancer on AWS
- AWS has 4 kinds of managed Load Balancers
- Classic Load Balancer (v1 - old generation) – 2009 – CLB
- HTTP , HTTPS, TCP , SSL (secure TCP)
- Application Load Balancer (v2 - new generation) – 2016 – ALB
- HTTP , HTTPS, WebSocket
- Network Load Balancer (v2 - new generation) – 2017 – NLB
- TCP , TLS (secure TCP), UDP
- Gateway Load Balancer – 2020 – GWLB
- Operates at layer 3 (Network layer) – IP Protocol
- Overall, it is recommended to use the newer generation load balancers as they provide more features
- Some load balancers can be setup as internal (private) or external (public) ELBs

## Slide 129

Load Balancer Security GroupsUsersHTTPS / HTTPFrom anywhereHTTP Restricted to Load balancerLOAD BALANCER
Load Balancer Security Group:
Application Security Group: Allow traffic only from Load Balancer
EC2

## Slide 130

Classic Load Balancers (v1)
ClientCLBEC2
listenerinternal
- Supports TCP (Layer 4), HTTP & HTTPS (Layer 7)
- Health checks are TCP or HTTP based
- Fixed hostname XXX.region.elb.amazonaws.com

## Slide 131

Application Load Balancer (v2)
- Application load balancers is Layer 7 (HTTP)
- Load balancing to multiple HTTP applications across machines (target groups)
- Load balancing to multiple applications on the same machine (ex: containers)
- Support for HTTP/2 and WebSocket
- Support redirects (from HTTP to HTTPS for example)

## Slide 132

Application Load Balancer (v2)
- Routing tables to different target groups:
- Routing based on path in URL (example.com/users & example.com/posts)
- Routing based on hostname in URL (one.example.com & other.example.com)
- Routing based on Query String, Headers (example.com/users?id=123&order=false)
- ALB are a great fit for micro services & container-based application (example: Docker & Amazon ECS)
- Has a port mapping feature to redirect to a dynamic port in ECS
- In comparison, we’d need multiple Classic Load Balancer per application

## Slide 133

Target Group for Users application
Application Load Balancer (v2)HTTP Based Traffic
External Application Load Balancer (v2)
HTTPWWWRoute /user
Target Group for Search applicationHTTPWWWRoute /search
Health CheckHealth Check

## Slide 134

Application Load Balancer (v2)Tar get Groups
- EC2 instances (can be managed by an Auto Scaling Group) – HTTP
- ECS tasks (managed by ECS itself) – HTTP
- Lambda functions – HTTP request is translated into a JSON event
- IP Addresses – must be private IPs
- ALB can route to multiple target groups
- Health checks are at the target group level

## Slide 135

Target Group 1AWS – EC2 based
Application Load Balancer (v2)Query Strings/Parameters Routing
External Application Load Balancer (v2)
?Platform=Mobile
WWWRequests
Target Group 2On-premises – Private IP routing?Platform=Desktop

## Slide 136

Application Load Balancer (v2)Good to Know
- Fixed hostname (XXX.region.elb.amazonaws.com)
- The application servers don’t see the IP of the client directly
- The true IP of the client is inserted in the header X-Forwarded-For
- We can also get Por t (X-Forwarded-Port) and proto (X-Forwarded-Proto)
EC2 InstanceClient IP 12.34.56.78Load Balancer IP (Private IP)
Connection termination

## Slide 137

Network Load Balancer (v2)
- Network load balancers (Layer 4) allow to:
- Forward TCP & UDP traffic to your instances
- Handle millions of request per seconds
- Ultra-low latency
- NLB has one static IP per AZ, and suppor ts assigning Elastic IP(helpful for whitelisting specific IP)
- NLB are used for extreme performance, TCP or UDP traffic

## Slide 138

Network Load Balancer (v2)TCP (Layer 4) Based Traffic
Target Group for Users applicationExternal Network Load Balancer (v2)
TCPWWWTCP + Rules
Target Group for Search applicationHTTPWWWTCP + Rules
Health CheckHealth Check

## Slide 139

Network Load Balancer – Tar get Groups
- EC2 instances
- IP Addresses – must be private IPs
- Application Load Balancer
- Health Checks support the TCP , HTTP and HTTPS Protocols
Target Group (EC2 Instances)
i-1234567890abcdef0
i-1234567890abcdef0
NetworkLoad Balancer
Target Group (IP Addresses)
192.168.1.11810.0.4.21
NetworkLoad Balancer
Target Group (Application Load Balancer)
NetworkLoad Balancer

## Slide 140

Gateway Load Balancer
- Deploy, scale, and manage a fleet of 3rd party network virtual appliances in AWS
- Example: Firewalls, Intrusion Detection and Prevention Systems, Deep Packet Inspection Systems, payload manipulation, …
- Operates at Layer 3 (Network Layer) – IP Packets
- Combines the following functions:
- Transparent Networ k Gateway – single entry/exit for all traffic
- Load Balancer – distributes traffic to your virtual appliances
- Uses the GENEVE protocol on port 6081
Users(source)
GatewayLoad Balancer
Application(destination)
Target Group
3rd Party SecurityVirtual Appliances
traffictraffic
RouteTable

## Slide 141

Gateway Load Balancer – Tar get Groups
- EC2 instances
- IP Addresses – must be private IPs
Target Group (EC2 Instances)
i-1234567890abcdef0
i-1234567890abcdef0
GatewayLoad Balancer
Target Group (IP Addresses)
192.168.1.11810.0.4.21
GatewayLoad Balancer

## Slide 142

Sticky Sessions (Session Affinity)
- It is possible to implement stickiness so that the same client is always redirected to the same instance behind a load balancer
- This works for Classic Load Balancer, Application Load Balancer, and Network Load Balancer
- For both CLB & ALB, the “cookie” used for stickiness has an expiration date you control
- Use case: make sure the user doesn’t lose his session data
- Enabling stickiness may bring imbalance to the load over the backend EC2 instances
EC2 Instance
EC2 Instance
Client 1Client 2Client 3

## Slide 143

Sticky Sessions – Cookie Names
- Application-based Cookies
- Custom cookie
- Generated by the target
- Can include any custom attributes required by the application
- Cookie name must be specified individually for each target group
- Don’t use AWSALB, AWSALBAPP, or AWSALBTG (reserved for use by the ELB)
- Application cookie
- Generated by the load balancer
- Cookie name is AWSALBAPP
- Duration-based Cookies
- Cookie generated by the load balancer
- Cookie name is AWSALB for ALB, AWSELB for CLB

## Slide 144

Cross-Zone Load Balancing
Availability Zone 1
Availability Zone 2
101010101010
10101010
Availability Zone 1
Availability Zone 2
25256.256.256.256.25
6.256.256.256.25
With Cross Zone Load Balancing:each load balancer instance distributes evenly across all registered instances in all AZWithout Cross Zone Load Balancing:Requests are distributed in the instances of the node of the Elastic Load Balancer
50 5050 50

## Slide 145

Cross-Zone Load Balancing
- Application Load Balancer
- Enabled by default (can be disabled at the Target Group level)
- No charges for inter AZ data
- Network Load Balancer & Gateway Load Balancer
- Disabled by default
- Yo u p ay c h a r g e s ( $ ) fo r i n t e r A Z d a t a i f e n a bl e d
- Classic Load Balancer
- Disabled by default
- No charges for inter AZ data if enabled

## Slide 146

SSL/TLS - Basics
- An SSL Certificate allows traffic between your clients and your load balancer to be encrypted in transit (in-flight encryption)
- SSLrefers to Secure Sockets Layer, used to encrypt connections
- TLSrefers to Transport Layer Security, which is a newer version
- Nowadays, TLS certificates are mainly used, but people still refer as SSL
- Public SSL certificates are issued by Certificate Authorities (CA)
- Comodo, Symantec, GoDaddy, GlobalSign, Digicert, Letsencrypt, etc…
- SSL certificates have an expiration date (you set) and must be renewed

## Slide 147

Load Balancer - SSL Certificates
- The load balancer uses an X.509 certificate (SSL/TLS server certificate)
- Yo u c a n m a n a g e c e r t i fi c a t e s u s i n g AC M ( AW S C e r t i fi c a t e M a n a g e r )
- Yo u c a n c r e a t e u p l o a d yo u r ow n c e r t i fi c a t e s a l t e r n a t i ve ly
- HTTPS listener:
- Yo u mu s t s p e c i f y a d e f a u l t c e r t i fi c a t e
- Yo u c a n a d d a n o p t i o n a l l i s t o f c e r t s t o s u p p o r t mu l t i p l e d o m a i n s
- Clients can use SNI (Server Name Indication) to specify the hostname they reach
- Ability to specify a security policy to support older versions of SSL / TLS (legacy clients)
EC2 InstanceUsersHTTPS (encrypted)Over wwwHTTPOver private VPCLOAD BALANCER

## Slide 148

SSL – Server Name Indication (SNI)
- SNI solves the problem of loading multiple SSL certificates onto one web server (to serve multiple websites)
- It’s a “newer” protocol, and requires the client to indicate the hostname of the target server in the initial SSL handshake
- The server will then find the correct certificate, or return the default oneNote:
- Only works for ALB & NLB (newer generation), CloudFront
- Does not work for CLB (older gen)
Client
ALB
SSL Cert:Domain1.example.com
SSL Cert:www.mycorp.com
I would likewww.mycorp.com
Target group forwww.mycorp.com
Target group forDomain1.example.com
….
Use the correctSSL cert

## Slide 149

Elastic Load Balancers – SSL Certificates
- Classic Load Balancer (v1)
- Support only one SSL certificate
- Must use multiple CLB for multiple hostname with multiple SSL certificates
- Application Load Balancer (v2)
- Supports multiple listeners with multiple SSL certificates
- Uses Server Name Indication (SNI) to make it work
- Network Load Balancer (v2)
- Supports multiple listeners with multiple SSL certificates
- Uses Server Name Indication (SNI) to make it work

## Slide 150

Connection Draining
- Feature naming
- Connection Draining – for CLB
- Deregistration Delay – for ALB & NLB
- Time to complete “in-flight requests” while the instance is de-registering or unhealthy
- Stops sending new requests to the EC2 instance which is de-registering
- Between 1 to 3600 seconds (default: 300 seconds)
- Can be disabled (set value to 0)
- Set to a low value if your requests are short
ELB
waiting for existingconnections to complete
new connections established to all other instances
Users
EC2 InstanceDRAINING
EC2 Instance
EC2 Instance

## Slide 151

What’s an Auto Scaling Group?
- In real-life, the load on your websites and application can change
- In the cloud, you can create and get rid of servers very quickly
- The goal of an Auto Scaling Group (ASG) is to:
- Scale out (add EC2 instances) to match an increased load
- Scale in (remove EC2 instances) to match a decreased load
- Ensure we have a minimum and a maximum number of EC2 instances running
- Automatically register new instances to a load balancer
- Re-create an EC2 instance in case a previous one is terminated (ex: if unhealthy)
- ASG are free (you only pay for the underlying EC2 instances)

## Slide 152

Auto Scaling Group in AWS
EC2Instance
EC2Instance
EC2Instance
EC2InstanceEC2InstanceEC2InstanceEC2Instance
Auto Scaling Group
Minimum Capacity
Desired Capacity
Maximum Capacity
Scale Out as Needed

## Slide 153

Auto Scaling Group in AWS With Load Balancer
Elastic Load Balancer
Users
ELB can check the health of your EC2 instances!
EC2Instance
EC2Instance
EC2Instance
EC2InstanceEC2InstanceEC2InstanceEC2Instance
Auto Scaling Group

## Slide 154

Auto Scaling Group Attributes
- A Launch Template (older “Launch Configurations” are deprecated)
- AMI + Instance Type
- EC2 User Data
- EBS Volumes
- Security Groups
- SSH Key Pair
- IAM Roles for your EC2 Instances
- Network + Subnets Information
- Load Balancer Information
- Min Size / Max Size / Initial Capacity
- Scaling Policies
ASG Launch Template
AMI
InstanceType
EBS Volumes
Security Groups
SSH Key Pair
IAM Role
VPC + Subnets…
LoadBalancer

## Slide 155

Auto Scaling - CloudWatch Alarms & Scaling
- It is possible to scale an ASG based on CloudWatch alarms
- An alarm monitors a metric (such as Average CPU, or a custom metric)
- Metrics such as Average CPU are computed for the overall ASG instances
- Based on the alarm:
- We can create scale-out policies (increase the number of instances)
- We can create scale-in policies (decrease the number of instances)
EC2Instance
EC2Instance
EC2InstanceEC2InstanceEC2Instance
Auto Scaling Group
CloudWatchAlarm
trigger Scaling

## Slide 156

Auto Scaling Groups – Scaling Policies
- Dynamic Scaling
- Tar get Tr acking Scaling
- Simple to set-up
- Example: I want the average ASG CPU to stay at around 40%
- Simple / Step Scaling
- When a CloudWatch alarm is triggered (example CPU > 70%), then add 2 units
- When a CloudWatch alarm is triggered (example CPU < 30%), then remove 1
- Scheduled Scaling
- Anticipate a scaling based on known usage patterns
- Example: increase the min capacity to 10 at 5 pm on Fridays

## Slide 157

Auto Scaling Groups – Scaling Policies
- Predictive scaling: continuously forecast load and schedule scaling ahead

## Slide 158

Good metrics to scale on
- CPUUtilization: Average CPU utilization across your instances
- RequestCountPerTarget: to make sure the number of requests per EC2 instances is stable
- Average Network In / Out (if you’re application is network bound)
- Any custom metric (that you push using CloudWatch)
Auto Scaling group
Users
ApplicationLoad Balancer
RequestCountPerTargetTarget Value: 3

## Slide 159

Auto Scaling Groups - Scaling Cooldowns
- After a scaling activity happens, you are in the cooldown period (default 300 seconds)
- During the cooldown period, the ASG will not launch or terminate additional instances (to allow for metrics to stabilize)
- Advice: Use a ready-to-use AMI to reduce configuration time in order to be serving request fasters and reduce the cooldown period
Scaling ActionOccurs
Default Cooldown in effect?
Ignore Action
Launch or Teminate InstanceNo
Yes

## Slide 160

RDS, Aurora & ElastiCache

## Slide 161

Amazon RDS Overview
- RDS stands for Relational Database Service
- It’s a managed DB service for DB use SQL as a query language.
- It allows you to create databases in the cloud that are managed by AWS
- Postgres
- MySQL
- MariaDB
- Oracle
- Microsoft SQL Server
- IBM DB2
- Aurora (AWS Proprietary database)

## Slide 162

Advantage over using RDS versus deploying DB on EC2
- RDS is a managed service:
- Automated provisioning, OS patching
- Continuous backups and restore to specific timestamp (Point in Time Restore)!
- Monitoring dashboards
- Read replicas for improved read performance
- Multi AZ setup for DR (Disaster Recovery)
- Maintenance windows for upgrades
- Scaling capability (vertical and horizontal)
- Storage backed by EBS
- BUT you can’t SSH into your instances

## Slide 163

RDS – Storage Auto Scaling
- Helps you increase storage on your RDS DB instance dynamically
- When RDS detects you are running out of free database storage, it scales automatically
- Avoid manually scaling your database storage
- Yo u h ave t o s e t Maximum Storage Threshold (maximum limit for DB storage)
- Automatically modify storage if:
- Free storage is less than 10% of allocated storage
- Low-storage lasts at least 5 minutes
- 6 hours have passed since last modification
- Useful for applications with unpredictable workloads
- Supports all RDS database engines
Read/Write
Application
Storage

## Slide 164

RDS Read Replicas for read scalability
RDS DB instanceRDS DB instance read replica
RDS DB instance read replica
Application
writesreads readsreads
ASYNC replicationASYNC replication
- Up to 15 Read Replicas
- Within AZ, Cross AZ or Cross Region
- Replication is ASYNC, so reads are eventually consistent
- Replicas can be promoted to their own DB
- Applications must update the connection string to leverage read replicas

## Slide 165

RDS Read Replicas – Use Cases
- Yo u h ave a p r o d u c t i o n d a t a b a s e that is taking on normal load
- Yo u w a n t t o r u n a r e p o r t i n g application to run some analytics
- Yo u c r e a t e a R e a d R e p l i c a t o r u n the new workload there
- The production application is unaffected
- Read replicas are used for SELECT (=read) only kind of statements (not INSERT, UPDATE, DELETE)
RDS DB instance
ProductionApplication
reads
RDS DB instance read replica
readsreads
ASYNC replication
Reporting Application

## Slide 166

RDS Read Replicas – Network Cost
RDS DB instanceRDS DB instance read replica
ASYNC ReplicationCross-Region$$$
Region/AZus-east-1a Region/AZeu-west-1b
- In AWS there’s a network cost when data goes from one AZ to another
- For RDS Read Replicas within the same region, you don’t pay that fee
VS
RDS DB instanceRDS DB instance read replica
ASYNC ReplicationSame Region Free
Same Region / Different AZus-east-1a us-east-1b

## Slide 167

RDS Multi AZ (Disaster Recovery)
RDS Master DB instance (AZ A)
Application
writesreads
SYNCreplication
- SYNC replication
- One DNS name – automatic app failover to standby
- Increase availability
- Failover in case of loss of AZ, loss of network, instance or storage failure
- No manual intervention in apps
- Not used for scaling
- Note: The Read Replicas be setup as Multi AZ for Disaster Recovery (DR)RDS DB instance standby (AZ B)
One DNS name – automatic failover

## Slide 168

RDS – From Single-AZ to Multi-AZ
- Zero downtime operation (no need to stop the DB)
- Just click on “modify” for the database
- The following happens internally:
- A snapshot is taken
- A new DB is restored from the snapshot in a new AZ
- Synchronization is established between the two databases
RDS DB instance
SYNC Replication
Standby DB
snapshotrestore
DB snapshot

## Slide 169

RDS Custom
- Managed Oracle and Microsoft SQL Server Database with OS and database customization
- RDS: Automates setup, operation, and scaling of database in AWS
- Custom: access to the underlying database and OS so you can
- Configure settings
- Install patches
- Enable native features
- Access the underlying EC2 Instance using SSH or SSM Session Manager
- De-activate Automation Mode to perform your customization, better to take a DB snapshot before
- RDS vs. RDS Custom
- RDS: entire database and the OS to be managed by AWS
- RDS Custom: full admin access to the underlying OS and the database
EC2 Instance
User
SSHapplycstomizations
AutomationMode disabled

## Slide 170

Amazon Aurora
- Aurora is a proprietary technology from AWS (not open sourced)
- Postgres and MySQL are both supported as Aurora DB (that means your drivers will work as if Aurora was a Postgres or MySQL database)
- Aurora is “AWS cloud optimized” and claims 5x performance improvement over MySQL on RDS, over 3x the performance of Postgres on RDS
- Aurora storage automatically grows in increments of 10GB, up to 256 TB.
- Aurora can have up to 15 replicas and the replication process is faster than MySQL (sub 10 ms replica lag)
- Failover in Aurora is instantaneous. It’s HA (High Availability) native.
- Aurora costs more than RDS (20% more) – but is more efficient

## Slide 171

Aurora High Availability and Read Scaling
- 6 copies of your data across 3 AZ:
- 4 copies out of 6 needed for writes
- 3 copies out of 6 need for reads
- Self healing with peer-to-peer replication
- Storage is striped across 100s of volumes
- One Aurora Instance takes writes (master)
- Automated failover for master in less than 30 seconds
- Master + up to 15 Aurora Read Replicas serve reads
- Support for Cross Region Replication
Shared storage VolumeReplication + Self Healing + Auto Expanding
AZ 1 AZ 2 AZ 3
WRRRRR

## Slide 172

Auto Scaling
Aurora DB Cluster
Shared storage VolumeAuto Expanding from 10G to 256 TB
WRRRRR
Reader EndpointConnection Load BalancingWriter EndpointPointing to the master
client

## Slide 173

Features of Aurora
- Automatic fail-over
- Backup and Recovery
- Isolation and security
- Industry compliance
- Push-button scaling
- Automated Patching with Zero Downtime
- Advanced Monitoring
- Routine Maintenance
- Backtrack: restore data at any point of time without using backups

## Slide 174

Replicas Auto Scaling
Aurora Replicas - Auto Scaling
Writer EndpointReader Endpoint
Shared Storage VolumeW RRRR
Client
CPUUsageCPUUsage
Many Requests
Endpoint Extended

## Slide 175

Aurora – Custom Endpoints
- Define a subset of Aurora Instances as a Custom Endpoint
- Example: Run analytical queries on specific replicas
- The Reader Endpoint is generally not used after defining Custom Endpoints
db.r5.2xlarge
Writer EndpointReader Endpoint
Shared Storage VolumeWRRRR
ClientQueriesCustom Endpoint
db.r3.largedb.r3.large
db.r5.2xlarge
Analytical Queries

## Slide 176

Aurora Serverless
Shared storage Volume
Proxy Fleet(managed by Aurora)
Client
- Automated database instantiation and auto-scaling based on actual usage
- Good for infrequent, intermittent or unpredictable workloads
- No capacity planning needed
- Pay per second, can be more cost-effective

## Slide 177

- Aurora Cross Region Read Replicas:
- Useful for disaster recovery
- Simple to put in place
- Aurora Global Database (recommended):
- 1 Primary Region (read / write)
- Up to 10 secondary (read-only) regions, replication lag is less than 1 second
- Up to 16 Read Replicas per secondary region
- Helps for decreasing latency
- Promoting another region (for disaster recovery) has an RTO of < 1 minute
- Typical cross-region replication takes less than 1 second
us-east-1 - PRIMARY region
ApplicationsRead / Write
eu-west-1 - SECONDARY region
ApplicationsRead Only
replication
Global Aurora

## Slide 178

Aurora Machine Learning
- Enables you to add ML-based predictions to your applications via SQL
- Simple, optimized, and secure integration between Aurora and AWS ML services
- Supported services
- Amazon SageMaker (use with any ML model)
- Amazon Comprehend (for sentiment analysis)
- Yo u d o n ’ t n e e d t o h ave M L e x p e r i e n c e
- Use cases: fraud detection, ads targeting, sentiment analysis, product recommendations
Amazon SageMakerAmazon Comprehend
Amazon Aurora
Application
SQL query(Recommended products?)
data(user’s profile, shopping history, …)predictions(red shirt, blue pants, …)
query results(red shirt, blue …)

## Slide 179

Babelfish for Aurora PostgreSQL
- Allows Aurora PostgreSQL to understand commands targeted for MS SQL Server(e.g., T-SQL)
- Therefore Microsoft SQL Server based applications can work on Aurora PostgreSQL
- Requires no to little code changes (using the same MS SQL Server client driver)
- The same applications can be used after a migration of your database (using AWS SCT and DMS)
Aurora PostgreSQL
Babelfish
ApplicationSQL Server Client Driver
T-SQL PostgreSQL
ApplicationPostgreSQL Driver
PL/pgSQL
migrate
T-SQL

## Slide 180

RDS Backups
- Automated backups:
- Daily full backup of the database (during the backup window)
- Transaction logs are backed-up by RDS every 5 minutes
- => ability to restore to any point in time (from oldest backup to 5 minutes ago)
- 1 to 35 days of retention, set 0 to disable automated backups
- Manual DB Snapshots
- Manually triggered by the user
- Retention of backup for as long as you want
- Trick: in a stopped RDS database , you will still pay for storage . If you plan on stopping it for a long time, you should snapshot & restore instead

## Slide 181

Aurora Backups
- Automated backups
- 1 to 35 days (cannot be disabled)
- point-in-time recovery in that timeframe
- Manual DB Snapshots
- Manually triggered by the user
- Retention of backup for as long as you want

## Slide 182

RDS & Aurora Restore options
- Restoring a RDS / Aurora backup or a snapshot creates a new database
- Restoring MySQL RDS database from S3
- Create a backup of your on-premises database
- Store it on Amazon S3 (object storage)
- Restore the backup file onto a new RDS instance running MySQL
- Restoring MySQL Aurora cluster from S3
- Create a backup of your on-premises database using Percona XtraBackup
- Store the backup file on Amazon S3
- Restore the backup file onto a new Aurora cluster running MySQL

## Slide 183

Aurora Database Cloning
- Create a new Aurora DB Cluster from an existing one
- Faster than snapshot & restore
- Uses copy-on-write protocol
- Initially, the new DB cluster uses the same data volume as the original DB cluster (fast and efficient – no copying is needed)
- When updates are made to the new DB cluster data, then additional storage is allocated and data is copied to be separated
- Ver y fast & cost-effective
- Useful to create a “staging” database from a “production” database without impacting the production database
Production AuroraStaging Aurora
clone

## Slide 184

RDS & Aurora Security
- At-rest encryption:
- Database master & replicas encryption using AWS KMS – must be defined as launch time
- If the master is not encrypted, the read replicas cannot be encrypted
- To e n c r y p t a n u n-encrypted database, go through a DB snapshot & restore as encrypted
- In-flight encryption: TLS-ready by default, use the AWS TLS root certificates client-side
- IAM Authentication: IAM roles to connect to your database (instead of username/pw)
- Security Groups: Control Network access to your RDS / Aurora DB
- No SSH available except on RDS Custom
- Audit Logs can be enabled and sent to CloudWatch Logs for longer retention

## Slide 185

Amazon RDS Proxy
- Fully managed database proxy for RDS
- Allows apps to pool and share DB connections established with the database
- Improving database efficiency by reducing the stress on database resources (e.g., CPU, RAM) and minimize open connections (and timeouts)
- Serverless, autoscaling, highly available (multi-AZ)
- Reduced RDS & Aurora failover time by up 66%
- Supports RDS (MySQL, PostgreSQL, MariaDB, MS SQL Server) and Aurora (MySQL, PostgreSQL)
- No code changes required for most apps
- Enforce IAM Authentication for DB, and securely store credentials in AWS Secrets Manager
- RDS Proxy is never publicly accessible (must be accessed from VPC)
IAMAuthentication
VPC
Private subnet
…Lambda functions
RDS Proxy
RDS DBInstance

## Slide 186

Amazon ElastiCache Overview
- The same way RDS is to get managed Relational Databases…
- ElastiCache is to get managed Redis or Memcached
- Caches are in-memory databases with really high performance, low latency
- Helps reduce load off of databases for read intensive workloads
- Helps make your application stateless
- AWS takes care of OS maintenance / patching, optimizations, setup, configuration, monitoring, failure recovery and backups
- Using ElastiCache involves heavy application code changes

## Slide 187

ElastiCache Solution Architecture - DB Cache
- Applications queries ElastiCache, if not available, get from RDS and store in ElastiCache.
- Helps relieve load in RDS
- Cache must have an invalidation strategy to make sure only the most current data is used in there.
application
Cache hit
AmazonRDS
Amazon ElastiCache
Cache missRead from DB
Write to cache

## Slide 188

ElastiCache Solution Architecture – User Session Store
- User logs into any of the application
- The application writes the session data into ElastiCache
- The user hits another instance of our application
- The instance retrieves the data and the user is already logged in
application
Retrieve sessionAmazon ElastiCacheapplication
application
Write session
User

## Slide 189

ElastiCache – Redis vs MemcachedREDIS
- Multi AZ with Auto-Failover
- Read Replicas to scale reads and have high availability
- Data Durability using AOF persistence
- Backup and restore features
- Supports Sets and Sorted Sets
MEMCACHED
- Multi-node for partitioning of data (sharding)
- No high availability (replication)
- Non persistent
- Backup and restore (Serverless)
- Multi-threaded architecture
+Replication sharding

## Slide 190

ElastiCache – Cache Security
- ElastiCache supports IAM Authentication for Redis
- IAM policies on ElastiCache are only used for AWS API-level security
- Redis AUTH
- Yo u c a n s e t a “ p a s s wo r d / t o ke n ” w h e n yo u c r e a t e a Redis cluster
- This is an extra level of security for your cache (on top of security groups)
- Support SSL in flight encryption
- Memcached
- Supports SASL-based authentication (advanced)
Client
EC2
Redis Security group
SSL encryptionRedis AUTH
EC2 Security group

## Slide 191

Patterns for ElastiCache
- Lazy Loading: all the read data is cached, data can become stale in cache
- Write Through: Adds or update data in the cache when written to a DB (no stale data)
- Session Store: store temporary session data in a cache (using TTL features)
- Quote: There are only two hard things in Computer Science: cache invalidation and naming things
application
Cache hit
AmazonRDS
Amazon ElastiCache
Cache missRead from DB
Write to cache
Lazy Loading illustrated

## Slide 192

ElastiCache – Redis Use Case
- Gaming Leaderboards are computationally complex
- Redis Sorted sets guarantee both uniqueness and element ordering
- Each time a new element added, it’s ranked in real time, then added in correct order
Clients
ElastiCache for Redis
ElastiCache for Redis
ElastiCache for Redis
1
2
3
Real-time Leaderboard

## Slide 193

Amazon Route 53

## Slide 194

What is DNS?
- Domain Name System which translates the human friendly hostnames into the machine IP addresses
- www.google.com => 172.217.18.36
- DNS is the backbone of the Internet
- DNS uses hierarchical naming structure
example.com.com
api.example.comwww.example.com

## Slide 195

DNS Terminologies
http://api.www.example.com.
- Domain Registrar: Amazon Route 53, GoDaddy, …
- DNS Records: A, AAAA, CNAME, NS, …
- Zone File: contains DNS records
- Name Server: resolves DNS queries (Authoritative or Non-Authoritative)
- To p L e ve l D o m a i n ( T L D ) : .com, .us, .in, .gov, .org, …
- Second Level Domain (SLD): amazon.com, google.com, …
TLDSLDSub DomainFQDN (Fully Qualified Domain Name)
Protocol Root
URL

## Slide 196

How DNS Works
Web BrowserLocal DNS ServerYou want to access example.comAssigned and Managed by your company or assigned by your ISP dynamically
example.com?
Root DNS Serverexample.com?
TLD DNS Server(.com)
SLD DNS Server (example.com)
.com NS 1.2.3.4example.com?example.com NS 5.6.7.8
example.com IP 9.10.11.12
Managed by ICANN
Managed by Domain Registrar(e.g., Amazon Registrar, Inc.)
Managed by IANA (Branch of ICANN)9.10.11.12
Web Server (example.com) (IP: 9.10.11.12)
TTL
TTL
example.com?

## Slide 197

Amazon Route 53
- A highly available, scalable, fully managed and Authoritative DNS
- Authoritative = the customer (you) can update the DNS records
- Route 53 is also a Domain Registrar
- Ability to check the health of your resources
- The only AWS service which provides 100% availability SLA
- Why Route 53? 53 is a reference to the traditional DNS port
example.com?
54.22.33.44Client
AmazonRoute 53
EC2 Instance
AWS Cloud
Public IP54.22.33.44

## Slide 198

Route 53 – Records
- How you want to route traffic for a domain
- Each record contains:
- Domain/subdomain Name – e.g., example.com
- Record Type – e.g., A or AAAA
- Value – e.g., 12.34.56.78
- Routing Policy – how Route 53 responds to queries
- TTL – amount of time the record cached at DNS Resolvers
- Route 53 supports the following DNS record types:
- (must know) A / AAAA / CNAME / NS
- (advanced) CAA / DS / MX / NAPTR / PTR / SOA / TXT / SPF / SRV

## Slide 199

Route 53 – Record Types
- A – maps a hostname to IPv4
- AAAA – maps a hostname to IPv6
- CNAME – maps a hostname to another hostname
- The target is a domain name which must have an A or AAAA record
- Can’t create a CNAME record for the top node of a DNS namespace (Zone Apex)
- Example: you can’t create for example.com, but you can create for www.example.com
- NS – Name Servers for the Hosted Zone
- Control how traffic is routed for a domain

## Slide 200

Route 53 – Hosted Zones
- A container for records that define how to route traffic to a domain and its subdomains
- Public Hosted Zones – contains records that specify how to route traffic on the Internet (public domain names)application1.mypublicdomain.com
- Private Hosted Zones – contain records that specify how you route traffic within one or more VPCs (private domain names)application1.company.internal
- Y ou pay $0.50 per month per hosted zone

## Slide 201

Route 53 – Public vs. Private Hosted Zones
VPC
EC2 Instance(Public IP)ApplicationLoad Balancer
Public Hosted ZoneClient
S3 BucketAmazonCloudFront
example.com?54.22.33.44 VPC
EC2 Instance(webapp.example.internal)(Private IP)
Private Hosted Zone
EC2 Instance(api.example.internal)(Private IP)
DB Instance(db.example.internal)(Private IP)api.example.internal?10.0.0.10
db.example.internal?10.0.0.35
Public Hosted ZonePrivate Hosted Zone

## Slide 202

Route 53 – Records TTL (Time T o Live)
Client
Web Server
DNS Requestmyapp.example.com?
TTL
AmazonRoute 53A 12.34.56.78(with TTL)HTTP RequestHTTP ResponseWill cache the result for The TTL of the record
- High TTL – e.g., 24 hr
- Less traffic on Route 53
- Possibly outdated records
- Low TTL – e.g., 60 sec.
- More traffic on Route 53 ($$)
- Records are outdated for less time
- Easy to change records
- Except for Alias records, TTL is mandatory for each DNS record

## Slide 203

CNAME vs Alias
- AWS Resources (Load Balancer, CloudFront...) expose an AWS hostname:
- lb1-1234.us-east-2.elb.amazonaws.com and you want myapp.mydomain.com
- CNAME:
- Points a hostname to any other hostname. (app.mydomain.com => blabla.anything.com)
- ONLY FOR NON ROOT DOMAIN (aka. something.mydomain.com)
- Alias:
- Points a hostname to an AWS Resource (app.mydomain.com => blabla.amazonaws.com)
- Works for ROOT DOMAIN and NON ROOT DOMAIN (aka mydomain.com)
- Free of charge
- Native health check

## Slide 204

Route 53 – Alias Records
- Maps a hostname to an AWS resource
- An extension to DNS functionality
- Automatically recognizes changes in the resource’s IP addresses
- Unlike CNAME, it can be used for the top node of a DNS namespace (Zone Apex), e.g.: example.com
- Alias Record is always of type A/AAAA for AWS resources (IPv4 / IPv6)
- Yo u c a n ’ t s e t t h e T T L
MyALB-123456789.us-east-1.elb.amazonaws.com
AmazonRoute 53
ApplicationLoad Balancer
Alias Record (Enabled)Record NameTypeValueexample.comAMyALB-123456789.us-east-1.elb.amazonaws.com
AWS-Managed(IP Addresses might change)

## Slide 205

Route 53 – Alias Records Targets
- Elastic Load Balancers
- CloudFront Distributions
- API Gateway
- Elastic Beanstalk environments
- S3 Websites
- VPC Interface Endpoints
- Global Accelerator accelerator
- Route 53 record in the same hosted zone
- Y ou cannot set an ALIAS record for an EC2 DNS name
ElasticLoad BalancerAmazonCloudFrontAmazonAPI Gateway
Elastic BeanstalkS3 WebsitesVPC InterfaceEndpoints
Global AcceleratorRoute 53 Record(same Hosted Zone)

## Slide 206

Route 53 – Routing Policies
- Define how Route 53 responds to DNS queries
- Don’t get confused by the word “Routing”
- It’s not the same as Load balancer routing which routes the traffic
- DNS does not route any traffic, it only responds to the DNS queries
- Route 53 Supports the following Routing Policies
- Simple
- Weighted
- Failover
- Latency based
- Geolocation
- Multi-Value Answer
- Geoproximity (using Route 53 Traffic Flow feature)

## Slide 207

Routing Policies – Simple
- Typically, route tr affic to a single resource
- Can specify multiple values in the same record
- If multiple values are returned, a random one is chosen by the client
- When Alias enabled, specify only one AWS resource
- Can’t be associated with Health Checks
foo.example.comA 11.22.33.44Client
AmazonRoute 53
Single Value
foo.example.com
A 11.22.33.44A 55.66.77.88A 99.11.22.33Client
AmazonRoute 53
Multiple Value
chooses a random value

## Slide 208

Routing Policies – Weighted
- Control the % of the requests that go to each specific resource
- Assign each record a relative weight:
- 𝑡𝑟𝑎𝑓𝑓𝑖𝑐 (%)=!"#$%& ()* + ,-".#(#. *".)*/012 )( +33 &%" 4"#$%&, ()* +33 *".)*/,
- Weights don’t need to sum up to 100
- DNS records must have the same name and type
- Can be associated with Health Checks
- Use cases: load balancing between regions, testing new application versions…
- Assign a weight of 0 to a record to stop sending traffic to a resource
- If all records have weight of 0, then all records will be returned equally
AmazonRoute 53
Weight: 70
Weight: 10
Weight: 20
70%
20%
10%

## Slide 209

Routing Policies – Latency-based
- Redirect to the resource that has the least latency close to us
- Super helpful when latency for users is a priority
- Latency is based on traffic between users and AWS Regions
- Germany users may be directed to the US (if that’s the lowest latency)
- Can be associated with Health Checks (has a failover capability)
ALB(us-east-1)
ALB(ap-southeast-1)

## Slide 210

Route 53 – Health Checks
- HTTP Health Checks are only for public resources
- Health Check => Automated DNS Failover:1.Health checks that monitor an endpoint (application, server, other AWS resource)2.Health checks that monitor other health checks (Calculated Health Checks)3.Health checks that monitor CloudWatch Alarms (full control !!) – e.g., throttles of DynamoDB, alarms on RDS, custom metrics, … (helpful for private resources)
- Health Checks are integrated with CW metrics
Amazon Route 53DNS Record(latency, geoproximity, …)
ALB
Auto Scaling group
EC2 Instance
Health Check
Health Check
us-east-1
ALB
Auto Scaling group
EC2 Instance
eu-west-1

## Slide 211

eu-west-1
Health Checker(us-east-1)
Health Checks – Monitor an Endpoint
- About 15 global health checkers will check the endpoint health
- Healthy/Unhealthy Threshold – 3 (default)
- Interval – 30 sec (can set to 10 sec – higher cost)
- Supported protocol: HTTP , HTTPS and TCP
- If > 18% of health checkers report the endpoint is healthy, Route 53 considers it Healthy. Otherwise, it’s Unhealthy
- Ability to choose which locations you want Route 53 to use
- Health Checks pass only when the endpoint responds with the 2xx and 3xx status codes
- Health Checks can be setup to pass / fail based on the text in the first 5120 bytes of the response
- Configure you router/firewall to allow incoming requests from Route 53 Health Checkers
HTTP requestto /health
ALB
Auto Scaling group
EC2 Instance
Health Checker(us-west-1)
Health Checker(sa-east-1)
200 code
Must allow incomingrequests from Route 53Health Checkers IP address range
https://ip-ranges.amazonaws.com/ip-ranges.json

## Slide 212

Amazon Route 53Route 53 – Calculated Health Checks
- Combine the results of multiple Health Checks into a single Health Check
- Yo u c a n u s e OR, AND, or NOT
- Can monitor up to 256 Child Health Checks
- Specify how many of the health checks need to pass to make the parent pass
- Usage: perform maintenance to your website without causing all health checks to fail
Health Check(Parent)
Health Check(Child)
Health Check(Child)
Health Check(Child)
EC2 Instance
EC2 Instance
EC2 Instance
monitormonitormonitor

## Slide 213

Private subnet
Health Checks – Private Hosted Zones
- Route 53 health checkers are outside the VPC
- They can’t access private endpoints (private VPC or on-premises resource)
- Yo u c a n c r e a t e a CloudWatch Metric and associate a CloudWatch Alarm, then create a Health Check that checks the alarm itself
CloudWatch Alarm
VPC
Health Checker(us-east-1)
monitor
monitor

## Slide 214

Routing Policies – Failover (Active-Passive)EC2 Instance(Primary)Health Check(mandatory)Failover
Client
AmazonRoute 53
EC2 Instance(Secondary – Disaster Recovery)
DNS Requests

## Slide 215

Routing Policies – Geolocation
- Different from Latency-based!
- This routing is based on user location
- Specify location by Continent, Country or by US State (if there’s overlapping, most precise location selected)
- Should create a “Default” record (in case there’s no match on location)
- Use cases: website localization, restrict content distribution, load balancing, …
- Can be associated with Health Checks
A 11.22.33.44
A 55.66.77.88
DefaultA 99.11.22.33

## Slide 216

Routing Policies – Geoproximity
- Route traffic to your resources based on the geographic location of users and resources
- Ability to shift more traffic to resources based on the defined bias
- To c h a n g e t h e s i ze o f t h e g e o g r a p h i c r e g i o n , s p e c i f y bias values:
- T o expand (1 to 99) – more traffic to the resource
- To s h r i n k (-1 to -99) – less traffic to the resource
- Resources can be:
- AWS resources (specify AWS region)
- Non-AWS resources (specify Latitude and Longitude)
- Y ou must use Route 53 Traffic Flow to use this feature

## Slide 217

Routing Policies – Geoproximity
us-east-1
us-west-1
Bias: 0 Bias: 0

## Slide 218

Routing Policies – Geoproximity
us-east-1
us-west-1
Bias: 0 Bias: 50Higher bias in us-east-1

## Slide 219

Routing Policies – IP-based Routing
- Routing is based on clients’ IP addresses
- You provide a list of CIDRs for your clients and the corresponding endpoints/locations (user-IP-to-endpoint mappings)
- Use cases: Optimize performance, reduce network costs…
- Example: route end users from a particular ISP to a specific endpoint
LocationsCIDR blockslocation-1203.0.113.0/24location-2200.5.4.0/24
CIDR Collection
Route 53
Record NameValueIP-basedexample.com1.2.3.4location-1example.com5.6.7.8location-2
Records
EC2 Instance(1.2.3.4)
EC2 Instance(5.6.7.8)
User A(203.0.113.56)
User B(200.5.4.100)

## Slide 220

Routing Policies – Multi-Value
- Use when routing traffic to multiple resources
- Route 53 return multiple values/resources
- Can be associated with Health Checks (return only values for healthy resources)
- Up to 8 healthy records are returned for each Multi-Value query
- Multi-Value is not a substitute for having an ELB

## Slide 221

Domain Registar vs. DNS Service
- Yo u bu y o r r e g i s t e r yo u r d o m a i n n a m e w i t h a D o m a i n R e g i s t r a r t y p i c a l ly by paying annual charges (e.g., GoDaddy, Amazon Registrar Inc., …)
- The Domain Registrar usually provides you with a DNS service to manage your DNS records
- But you can use another DNS service to manage your DNS records
- Example: purchase the domain from GoDaddy and use Route 53 to manage your DNS records
AmazonRoute 53
purchaseexample.comUsermanage DNS records

## Slide 222

GoDaddy as Registrar & Route 53 as DNS Service
AmazonRoute 53
Public Hosted Zonestephanetheteacher.com

## Slide 223

3rd Party Registrar with Amazon Route 53
- If you buy your domain on a 3rd party registrar, you can still use Route 53 as the DNS Service provider1.Create a Hosted Zone in Route 532.Update NS Records on 3rd party website to use Route 53 Name Servers
- Domain Registrar != DNS Service
- But every Domain Registrar usually comes with some DNS features

## Slide 224

Route 53 – Hybrid DNS
- By default, Route 53 Resolver automatically answers DNS queries for:
- Local domain names for EC2 instances
- Records in Private Hosted Zones
- Records in public Name Servers
- Hybrid DNS – resolving DNS queries between VPC (Route 53 Resolver) and your networks (other DNS Resolvers)
- Networks can be:
- VPC itself / Peered VPC
- On-premises Network (connected through Direct Connect or AWS VPN)
VPC
Region
EC2 Instance(ec2-192-0-2-44.compute-1.amazonaws.com)
Private Hosted Zone
Route 53Resolver
Public Name Server

## Slide 225

Route 53 – Resolver Endpoints
- Inbound Endpoint – allows your DNS Resolvers to resolve domain names for AWS resources (e.g., EC2 instances) and records in Private Hosted Zones
Private Subnet
VPCus-east-1 On-Premises Data Center
DNS Resolvers(onpremise.private)
Server(web.onpremise.private)
ResolverInbound Endpoint
Route 53Resolver
Private Hosted Zone(aws.private)
EC2 Instance(app.aws.private) DNS Queryapp.aws.private?
DNS Queryapp.aws.private?
VPN or DX connection
lookup

## Slide 226

Route 53 – Resolver Endpoints
- Outbound Endpoint
- Route 53 Resolver forwards DNS queries to your DNS Resolvers
Private SubnetVPCus-east-1 On-Premises Data Center
DNS Resolvers(onpremise.private)
Server(web.onpremise.private)
ResolverOutbound Endpoint
Route 53Resolver
Private Hosted Zone(aws.private)
EC2 Instance(app.aws.private)DNS Queryweb.onpremise.private?
VPN or DX connection
DNS Queryweb.onpremise.private?

## Slide 227

Classic Solutions Architecture

## Slide 228

Section Introduction
- These solutions architectures are the best part of this course
- Let’s understand how all the technologies we’ve seen work together
- This is a section you need to be 100% comfortable with
- We’ll see the progression of a Solution’s architect mindset through many sample case studies:
- WhatIsTheTime.Com
- MyClothes.Com
- MyWordPress.Com
- Instantiating applications quickly
- Beanstalk

## Slide 229

Stateless Web App: WhatIsTheTime.com
- WhatIsTheTime.com allows people to know what time it is
- We don’t need a database
- We want to star t small and can accept downtime
- We want to fully scale ver tically and horizontally, no downtime
- Let’s go through the Solutions Architect journey for this app
- Let’s see how we can proceed!

## Slide 230

Stateless web app: What time is it?Starting simple
Public EC2
Elastic IP Address
User
What time is it?
5:30 pm!

## Slide 231

Stateless web app: What time is it?Scaling vertically
Public EC2
Elastic IP Address
User
What time is it?
5:30 pm!
What time is it?
6:30 pm!
What time is it?
7:30 pm!
Downtime while upgrading to M5

## Slide 232

Stateless web app: What time is it?Scaling horizontally
User
What time is it?
5:30 pm!What time is it?
6:30 pm!
What time is it?
7:30 pm!

## Slide 233

Stateless web app: What time is it?Scaling horizontally
What time is it?
5:30 pm!What time is it?
6:30 pm!
What time is it?
7:30 pm!
Public EC2 instance,No Elastic IPDNS QueryFor api.whatisthetime.comA RecordTTL 1 hour

## Slide 234

Stateless web app: What time is it?Scaling horizontally, adding and removing instances
What time is it?
5:30 pm!What time is it?
6:30 pm!
What time is it?
7:30 pm!
Public EC2 instance,No Elastic IP
DNS QueryFor api.whatisthetime.comA RecordTTL 1 hour INSTANCE IS GONE!

## Slide 235

Stateless web app: What time is it?Scaling horizontally, with a load balancerWhat time is it?
Private EC2 instances
DNS QueryFor api.whatisthetime.comAlias Record
Availability zone 1Availability zone 1
ELB +Health Checks
RestrictedSecurity groups rules

## Slide 236

Stateless web app: What time is it?Scaling horizontally, with an auto-scaling groupWhat time is it?
Private EC2 instances
DNS QueryFor api.whatisthetime.comAlias Record
Availability zone 1Availability zone 1
ELB +Health Checks
Auto Scaling group

## Slide 237

Stateless web app: What time is it?Making our app multi-AZWhat time is it?
DNS QueryFor api.whatisthetime.comAlias Record
Availability zone 1 to 3Availability zone 1
ELB +Health Checks+ Multi AZ
Auto Scaling group
Availability zone 2
Availability zone 3

## Slide 238

Minimum 2 AZ => Let’s reserve capacity
DNS QueryFor api.whatisthetime.comAlias Record
Availability zone 1 to 3Availability zone 1
ELB +Health Checks+ Multi AZ
Auto Scaling group
Availability zone 2
Minimum capacity= reserved instances = cost savings

## Slide 239

In this lecture we’ve discussed…
- Public vs Private IP and EC2 instances
- Elastic IP vs Route 53 vs Load Balancers
- Route 53 TTL, A records and Alias Records
- Maintaining EC2 instances manually vs Auto Scaling Groups
- Multi AZ to survive disasters
- ELB Health Checks
- Security Group Rules
- Reservation of capacity for costing savings when possible

## Slide 240

Stateful Web App: MyClothes.com
- MyClothes.com allows people to buy clothes online.
- There’s a shopping cart
- Our website is having hundreds of users at the same time
- We need to scale, maintain horizontal scalability and keep our web application as stateless as possible
- Users should not lose their shopping cart
- Users should have their details (address, etc) in a database
- Let’s see how we can proceed!

## Slide 241

Stateful Web App: MyClothes.com
Multi AZAvailability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3

## Slide 242

Stateful Web App: MyClothes.comIntroduce Stickiness (Session Affinity)
Multi AZAvailability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3
ELB Stickiness

## Slide 243

Stateful Web App: MyClothes.comIntroduce User Cookies
Multi AZAvailability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3
Send shopping cart content in Web CookiesStatelessHTTP requests are heavierSecurity risk (cookies can be altered)Cookies must be validatedCookies must be less than 4KB

## Slide 244

Stateful Web App: MyClothes.comIntroduce Server Session
Multi AZAvailability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3
Send session_id inWeb Cookies
ElastiCache
Store / retrieve session data
Amazon DynamoDB(alternative)

## Slide 245

Stateful Web App: MyClothes.comStoring User Data in a database
Multi AZAvailability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3
ElastiCache
Amazon RDS
Store / retrieve user data(address, name, etc)

## Slide 246

Stateful Web App: MyClothes.comScaling Reads
Multi AZAvailability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3
ElastiCache
RDSRead Replicas
RDSMaster(writes)replication

## Slide 247

Stateful Web App: MyClothes.comScaling Reads (Alternative) – Lazy Loading
Multi AZAvailability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3
ElastiCache
RDSRead/write
cache
Read from cachehit

## Slide 248

Stateful Web App: MyClothes.comMulti AZ – Survive disasters
Multi AZAvailability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3
ElastiCacheMulti AZ
RDSMulti AZ

## Slide 249

Stateful Web App: MyClothes.comSecurity Groups
Multi AZAvailability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3
ElastiCache
RDS
Restrict traffic to EC2Security group from the LBRestrict traffic to RDSSecurity group from theEC2 security group
Restrict traffic to ElastiCacheSecurity group from theEC2 security group
Open HTTP / HTTPS to 0.0.0.0/0

## Slide 250

In this lecture we’ve discussed…3-tier architectures for web applications
- ELB sticky sessions
- Web clients for storing cookies and making our web app stateless
- ElastiCache
- For storing sessions (alternative: DynamoDB)
- For caching data from RDS
- Multi AZ
- RDS
- For storing user data
- Read replicas for scaling reads
- Multi AZ for disaster recovery
- Tight Security with security groups referencing each other

## Slide 251

Stateful Web App: MyWordPress.com
- We are tr ying to create a fully scalable WordPress website
- We want that website to access and correctly display picture uploads
- Our user data, and the blog content should be stored in a MySQL database.
- Let’s see how we can achieve this!

## Slide 252

Stateful Web App: MyWordPress.comRDS layer
Multi AZAvailability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3
RDSMulti AZ

## Slide 253

Stateful Web App: MyWordPress.comScaling with Aurora: Multi AZ & Read Replicas
Multi AZAvailability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3
Aurora MySQLMulti AZRead Replicas

## Slide 254

Stateful Web App: MyWordPress.comStoring images with EBS
Multi AZ
Amazon EBS Volume
Availability zone 1
Send image

## Slide 255

Stateful Web App: MyWordPress.comStoring images with EBS
Multi AZ
Amazon EBS Volume
Amazon EBS Volume
Availability zone 1
Availability zone 2Send image

## Slide 256

Stateful Web App: MyWordPress.comStoring images with EFS
Multi AZ
Availability zone 1
Availability zone 2
EFS
ENI
ENI
Send image

## Slide 257

In this lecture we’ve discussed…
- Aurora Database to have easy Multi-AZ and Read-Replicas
- Storing data in EBS (single instance application)
- Vs Storing data in EFS (distributed application)

## Slide 258

Instantiating Applications quickly
- When launching a full stack (EC2, EBS, RDS), it can take time to:
- Install applications
- Insert initial (or recovery) data
- Configure everything
- Launch the application
- We can take advantage of the cloud to speed that up!

## Slide 259

Instantiating Applications quickly
- EC2 Instances:
- Use a Golden AMI: Install your applications, OS dependencies etc.. beforehand and launch your EC2 instance from the Golden AMI
- Bootstrap using User Data: For dynamic configuration, use User Data scripts
- Hybrid: mix Golden AMI and User Data (Elastic Beanstalk)
- RDS Databases:
- Restore from a snapshot: the database will have schemas and data ready!
- EBS Volumes:
- Restore from a snapshot: the disk will already be formatted and have data!

## Slide 260

Typical architecture: Web App 3-tier
Multi AZAvailability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3
ElastiCache
Store / retrieve session data+ Cached data
PUBLIC SUBNETPRIVATE SUBNETDATA SUBNET
Amazon RDSRead / write data
ELB
Route 53

## Slide 261

Developer problems on AWS
- Managing infrastructure
- Deploying Code
- Configuring all the databases, load balancers, etc
- Scaling concerns
- Most web apps have the same architecture (ALB + ASG)
- All the developers want is for their code to run!
- Possibly, consistently across different applications and environments

## Slide 262

Elastic Beanstalk – Overview
- Elastic Beanstalk is a developer centric view of deploying an application on AWS
- It uses all the component’s we’ve seen before: EC2, ASG, ELB, RDS, …
- Managed service
- Automatically handles capacity provisioning, load balancing, scaling, application health monitoring, instance configuration, …
- Just the application code is the responsibility of the developer
- We still have full control over the configuration
- Beanstalk is free but you pay for the underlying instances

## Slide 263

Elastic Beanstalk – Components
- Application: collection of Elastic Beanstalk components (environments, versions, configurations, …)
- Application Version: an iteration of your application code
- Environment
- Collection of AWS resources running an application version (only one application version at a time)
- Tiers: Web Ser ver Environment Tier & Worker Environment Tier
- Y ou can create multiple environments (dev, test, prod, …)
CreateApplicationUploadVersionLaunchEnvironmentManageEnvironment
update version
deploy new version

## Slide 264

Elastic Beanstalk – Supported Platforms
- Go
- Java SE
- Java with T omcat
- .NET Core on Linux
- .NET on Windows Server
- Node.js
- PHP
- Python
- Ruby
- Packer Builder
- Single Container Docker
- Multi-container Docker
- Preconfigured Docker

## Slide 265

Web Ser ver Tier vs. Worker Tier
Auto Scaling group
Web Environment(myapp.us-east-1.elasticbeanstalk.com)
Security Group
Availability Zone 1
Availability Zone 2
EC2 Instance(Web Server)
ELB
Security Group
EC2 Instance(Web Server)
Auto Scaling group
Worker Environment
Availability Zone 1Availability Zone 2
SQS Queue
SQS messageSQS messagepullmessages
EC2 Instance(Worker)
EC2 Instance(Worker)
- Scale based on the number of SQS messages
- Can push messages to SQS queue from another Web Server Tier

## Slide 266

Elastic Beanstalk Deployment ModesSingle InstanceGreat for devHigh Availability with Load BalancerGreat for prod
Availability Zone 1Auto Scaling GroupAvailability Zone 2
RDS Master
RDS Standby
EC2 Instance
EC2 Instance
ALBAvailability Zone 1
RDS Master
EC2 Instance
Elastic IP

## Slide 267

Amazon S3

## Slide 268

Section introduction
- Amazon S3 is one of the main building blocks of AWS
- It’s advertised as ”infinitely scaling” storage
- Many websites use Amazon S3 as a backbone
- Many AWS services use Amazon S3 as an integration as well
- We’ll have a step-by-step approach to S3

## Slide 269

Amazon S3 Use cases
- Backup and storage
- Disaster Recovery
- Archive
- Hybrid Cloud storage
- Application hosting
- Media hosting
- Data lakes & big data analytics
- Software delivery
- Static website
Nasdaq stores 7 years of data into S3 Glacier
Sysco runs analytics on its data and gain business insights

## Slide 270

Amazon S3 - Buckets
- Amazon S3 allows people to store objects (files) in “buckets” (directories)
- Buckets must have a globally unique name (across all regions all accounts)
- Buckets are defined at the region level
- S3 looks like a global service but buckets are created in a region
- Naming convention
- No uppercase, No underscore
- 3-63 characters long
- Not an IP
- Must start with lowercase letter or number
- Must NOT start with the prefix xn--
- Must NOT end with the suffix -s3alias
S3 Bucket

## Slide 271

Amazon S3 - Objects
- Objects (files) have a Key
- The key is the FULL path:
- s3://my-bucket/my_file.txt
- s3://my-bucket/my_folder1/another_folder/my_file.txt
- The key is composed of prefix + object name
- s3://my-bucket/my_folder1/another_folder/my_file.txt
- There’s no concept of “directories” within buckets(although the UI will trick you to think otherwise)
- Just keys with very long names that contain slashes (“/”)
S3 Bucketwith Objects
Object

## Slide 272

Amazon S3 – Objects (cont.)
- Object values are the content of the body:
- Max. Object Size is 5TB (5000GB)
- If uploading more than 5GB, must use “multi-part upload”
- Metadata (list of text key / value pairs – system or user metadata)
- Ta g s ( Un icode key / v a lu e pa ir – up to 10) – useful for security / lifecycle
- Version ID (if versioning is enabled)

## Slide 273

Amazon S3 – Security
- User-Based
- IAM Policies – which API calls should be allowed for a specific user from IAM
- Resource-Based
- Bucket Policies – bucket wide rules from the S3 console - allows cross account
- Object Access Control List (ACL) – finer grain (can be disabled)
- Bucket Access Control List (ACL) – less common (can be disabled)
- Note: an IAM principal can access an S3 object if
- The user IAM permissions ALLOW it OR the resource policy ALLOWS it
- AND there’s no explicit DENY
- Encryption: encrypt objects in Amazon S3 using encryption keys

## Slide 274

S3 Bucket Policies
- JSON based policies
- Resources: buckets and objects
- Effect: Allow / Deny
- Actions: Set of API to Allow or Deny
- Principal: The account or user to apply the policy to
- Use S3 bucket for policy to:
- Grant public access to the bucket
- Force objects to be encrypted at upload
- Grant access to another account (Cross Account)

## Slide 275

Example: Public Access - Use Bucket Policy
S3 Bucket
Anonymous www website visitor
S3 Bucket PolicyAllows Public Access

## Slide 276

Example: User Access to S3 – IAM permissions
S3 Bucket
IAM Policy
IAM User

## Slide 277

Example: EC2 instance access - Use IAM Roles
S3 Bucket
EC2 Instance Role
IAM permissions
EC2 Instance

## Slide 278

Advanced: Cross-Account Access – Use Bucket Policy
S3 Bucket
IAM UserOther AWS account
S3 Bucket PolicyAllows Cross-Account

## Slide 279

Bucket settings for Block Public Access
- These settings were created to prevent company data leaks
- If you know your bucket should never be public, leave these on
- Can be set at the account level

## Slide 280

Amazon S3 – Static Website Hosting
- S3 can host static websites and have them accessible on the Internet
- The website URL will be (depending on the region)
- http://bucket-name.s3-website-aws-region.amazonaws.comOR
- http://bucket-name.s3-website.aws-region.amazonaws.com
- If you get a 403 Forbidden error, make sure the bucket policy allows public reads!
S3 Bucket(demo-bucket)
us-west-2
User
http://demo-bucket.s3-website-us-west-2.amazonaws.comhttp://demo-bucket.s3-website.us-west-2.amazonaws.com

## Slide 281

Amazon S3 - Versioning
- Yo u c a n ve r s i o n yo u r fi l e s i n A m a zo n S 3
- It is enabled at the bucket level
- Same key overwrite will change the “version”: 1, 2, 3….
- It is best practice to version your buckets
- Protect against unintended deletes (ability to restore a version)
- Easy roll back to previous version
- Notes:
- Any file that is not versioned prior to enabling versioning will have version “null”
- Suspending versioning does not delete the previous versions
S3 Bucket(my-bucket)
s3://my-bucket/my-file.docx
User
upload
Version 1Version 2Version 3

## Slide 282

Amazon S3 – Replication (CRR & SRR)
- Must enable Versioning in source and destination buckets
- Cross-Region Replication (CRR)
- Same-Region Replication (SRR)
- Buckets can be in different AWS accounts
- Copying is asynchronous
- Must give proper IAM permissions to S3
- Use cases:
- CRR – compliance, lower latency access, replication across accounts
- SRR – log aggregation, live replication between production and test accounts
S3 Bucket(eu-west-1)
S3 Bucket(us-east-2)
asynchronousreplication

## Slide 283

Amazon S3 – Replication (Notes)
- After you enable Replication, only new objects are replicated
- Optionally, you can replicate existing objects using S3 Batch Replication
- Replicates existing objects and objects that failed replication
- For DELETE operations
- Can replicate delete markers from source to target (optional setting)
- Deletions with a version ID are not replicated (to avoid malicious deletes)
- There is no “chaining” of replication
- If bucket 1 has replication into bucket 2, which has replication into bucket 3
- Then objects created in bucket 1 are not replicated to bucket 3

## Slide 284

S3 Storage Classes
- Amazon S3 Standard - General Purpose
- Amazon S3 Standard-Infrequent Access (IA)
- Amazon S3 One Zone-Infrequent Access
- Amazon S3 Glacier Instant Retrieval
- Amazon S3 Glacier Flexible Retrieval
- Amazon S3 Glacier Deep Archive
- Amazon S3 Intelligent Tiering
- Can move between classes manually or using S3 Lifecycle configurations

## Slide 285

S3 Durability and Availability
- Durability:
- High durability (99.999999999%, 11 9’s) of objects across multiple AZ
- If you store 10,000,000 objects with Amazon S3, you can on average expect to incur a loss of a single object once every 10,000 years
- Same for all storage classes
- Availability:
- Measures how readily available a service is
- Varies depending on storage class
- Example: S3 standard has 99.99% availability = not available 53 minutes a year

## Slide 286

S3 Standard – General Purpose
- 99.99% Availability
- Used for frequently accessed data
- Low latency and high throughput
- Sustain 2 concurrent facility failures
- Use Cases: Big Data analytics, mobile & gaming applications, content distribution…

## Slide 287

S3 Storage Classes – Infrequent Access
- For data that is less frequently accessed, but requires rapid access when needed
- Lower cost than S3 Standard
- Amazon S3 Standard-Infrequent Access (S3 Standard-IA)
- 99.9% Availability
- Use cases: Disaster Recovery, backups
- Amazon S3 One Zone-Infrequent Access (S3 One Zone-IA)
- High durability (99.999999999%) in a single AZ; data lost when AZ is destroyed
- 99.5% Availability
- Use Cases: Storing secondary backup copies of on-premises data, or data you can recreate

## Slide 288

Amazon S3 Glacier Storage Classes
- Low-cost object storage meant for archiving / backup
- Pricing: price for storage + object retrieval cost
- Amazon S3 Glacier Instant Retrieval
- Millisecond retrieval, great for data accessed once a quarter
- Minimum storage duration of 90 days
- Amazon S3 Glacier Flexible Retrieval (formerly Amazon S3 Glacier):
- Expedited (1 to 5 minutes), Standard (3 to 5 hours), Bulk (5 to 12 hours) – free
- Minimum storage duration of 90 days
- Amazon S3 Glacier Deep Archive – for long term storage:
- Standard (12 hours), Bulk (48 hours)
- Minimum storage duration of 180 days

## Slide 289

S3 Intelligent-Tiering
- Small monthly monitoring and auto-tiering fee
- Moves objects automatically between Access Tiers based on usage
- There are no retrieval charges in S3 Intelligent-Tiering
- Frequent Access tier (automatic): default tier
- Infrequent Access tier (automatic): objects not accessed for 30 days
- Archive Instant Access tier (automatic): objects not accessed for 90 days
- Archive Access tier (optional): configurable from 90 days to 700+ days
- Deep Archive Access tier (optional): config. from 180 days to 700+ days

## Slide 290

S3 Storage Classes Comparison
https://aws.amazon.com/s3/storage-classes/
StandardIntelligent-TieringStandard-IAOne Zone-IAGlacier Instant RetrievalGlacier Flexible RetrievalGlacier Deep ArchiveDurability 99.999999999% == (11 9’s)Availability99.99%99.9%99.9%99.5%99.9%99.99%99.99%Availability SLA99.9%99%99%99%99%99.9%99.9%Availability Zones>= 3>= 3>= 31>= 3>= 3>= 3Min. Storage Duration ChargeNoneNone30 Days30 Days90 Days90 Days180 DaysMin. Billable Object SizeNoneNone128 KB128 KB128 KB40 KB40 KBRetrieval FeeNoneNonePer GB retrievedPer GB retrievedPer GB retrievedPer GB retrievedPer GB retrieved

## Slide 291

S3 Storage Classes – Price ComparisonExample: us-east-1
https://aws.amazon.com/s3/pricing/
StandardIntelligent-TieringStandard-IAOne Zone-IAGlacier Instant RetrievalGlacier Flexible RetrievalGlacier Deep ArchiveStorage Cost(per GB per month)$0.023$0.0025 - $0.023$0.0125$0.01$0.004$0.0036$0.00099
Retrieval Cost(per 1000 request)GET: $0.0004POST: $0.005GET: $0.0004POST: $0.005GET: $0.001POST: $0.01GET: $0.001POST: $0.01GET: $0.01POST: $0.02
GET: $0.0004POST: $0.03Expedited: $10Standard: $0.05Bulk: free
GET: $0.0004POST: $0.05Standard: $0.10Bulk: $0.025
Retrieval TimeInstantaneousExpedited (1 – 5 mins)Standard (3 – 5 hours)Bulk (5 – 12 hours)Standard (12 hours)Bulk (48 hours)Monitoring Cost (per 1000 objects)$0.0025

## Slide 292

S3 Express One Zone
- High performance, single Availability Zone storage class
- Objects stored in a Directory Bucket (bucket in a single AZ)
- Handle 100,000s requests per second with single-digit millisecond latency
- Up to 10x better performance than S3 Standard (50% lower costs)
- High Durability (99.999999999%) and Availability (99.95%)
- Co-locate your storage and compute resources in the same AZ (reduces latency)
- Use cases: latency-sensitive apps, data-intensive apps, AI & ML training, financial modeling, media processing, HPC…
- Best integrated with SageMaker Model Training, Athena, EMR, Glue…
stephane--use1-az4--x-s3
Region (us-east-1)
Availability Zone (AZ 4)

## Slide 293

Amazon S3 – Advanced

## Slide 294

Amazon S3 – Moving between Storage Classes
- Yo u c a n t r a n s i t i o n o b j e c t s b e t we e n storage classes
- For infrequently accessed object, move them to Standard IA
- For archive objects that you don’t need fast access to, move them to Glacier or Glacier Deep Archive
- Moving objects can be automated using a Lifecycle Rules
StandardStandard IA
Intelligent TieringOne-Zone IA
Glacier Instant RetrievalGlacier Flexible Retrieval
Glacier Deep Archive

## Slide 295

Amazon S3 – Lifecycle Rules
- Transition Actions – configure objects to transition to another storage class
- Move objects to Standard IA class 60 days after creation
- Move to Glacier for archiving after 6 months
- Expiration actions – configure objects to expire (delete) after some time
- Access log files can be set to delete after a 365 days
- Can be used to delete old versions of files (if versioning is enabled)
- Can be used to delete incomplete Multi-Part uploads
- Rules can be created for a certain prefix (example: s3://mybucket/mp3/*)
- Rules can be created for certain objects Tags (example: Department: Finance)

## Slide 296

Amazon S3 – Lifecycle Rules (Scenario 1)
- Y our application on EC2 creates images thumbnails after profile photos are uploaded to Amazon S3. These thumbnails can be easily recreated, and only need to be kept for 60 days. The source images should be able to be immediately retrieved for these 60 days, and afterwards, the user can wait up to 6 hours. How would you design this?
- S3 source images can be on Standard, with a lifecycle configuration to transition them to Glacier after 60 days
- S3 thumbnails can be on One-Zone IA, with a lifecycle configuration to expire them (delete them) after 60 days

## Slide 297

Amazon S3 – Lifecycle Rules (Scenario 2)
- A rule in your company states that you should be able to recover your deleted S3 objects immediately for 30 days, although this may happen rarely. After this time, and for up to 365 days, deleted objects should be recoverable within 48 hours.
- Enable S3 Versioning in order to have object versions, so that “deleted objects” are in fact hidden by a “delete marker” and can be recovered
- Transition the “noncurrent ver sions” of the object to Standard IA
- Transition afterwards the “noncurrent ver sions” to Glacier Deep Archive

## Slide 298

Amazon S3 Analytics – Storage Class Analysis
- Help you decide when to transition objects to the right storage class
- Recommendations for Standard and Standard IA
- Does NOT work for One-Zone IA or Glacier
- Report is updated daily
- 24 to 48 hours to start seeing data analysis
- Good first step to put together Lifecycle Rules (or improve them)!
S3 Bucket
S3 Analytics
DateStorageClassObjectAge8/22/2022STANDARD000-0148/25/2022STANDARD030-0449/6/2022STANDARD120-149
.csv report

## Slide 299

S3 – Requester Pays
- In general, bucket owners pay for all Amazon S3 storage and data transfer costs associated with their bucket
- With Requester Pays buckets, the requester instead of the bucket owner pays the cost of the request and the data download from the bucket
- Helpful when you want to share large datasets with other accounts
- The requester must be authenticated in AWS (cannot be anonymous)
Owner$$ Storage Cost
Owner$$ Networking Cost
download
Requester
Owner$$ Storage CostRequester$$ Networking Cost
download
Requester Pays Bucket
Standard Bucket

## Slide 300

S3 Event Notifications
Amazon S3
events
Lambda Function
SQS
SNS
- S3:ObjectCreated, S3:ObjectRemoved, S3:ObjectRestore, S3:Replication…
- Object name filtering possible (*.jpg)
- Use case: generate thumbnails of images uploaded to S3
- Can create as many “S3 events” as desired
- S3 event notifications typically deliver events in seconds but can sometimes take a minute or longer

## Slide 301

S3 Event Notifications – IAM Permissions
Amazon S3
events
Lambda Function
SQS
SNS
Lambda Resource Policy
SNS Resource (Access) Policy
SQS Resource (Access) Policy

## Slide 302

S3 Event Notifications with Amazon EventBridge
Amazon S3bucket
eventsAll events
AmazonEventBridge
rulesOver 18 AWS servicesas destinations
- Advanced filtering options with JSON rules (metadata, object size, name...)
- Multiple Destinations – ex Step Functions, Kinesis Streams / Firehose…
- EventBridge Capabilities – Archive, Replay Events, Reliable delivery

## Slide 303

S3 – Baseline Performance
- Amazon S3 automatically scales to high request rates, latency 100-200 ms
- Yo u r a p p l i c a t i o n c a n a c h i e ve a t l e a s t 3,500 PUT/COPY/POST/DELETE or 5,500 GET/HEAD requests per second per prefix in a bucket.
- There are no limits to the number of prefixes in a bucket.
- Example (object path => prefix):
- bucket/folder1/sub1/file => /folder1/sub1/
- bucket/folder1/sub2/file => /folder1/sub2/
- bucket/1/file => /1/
- bucket/2/file => /2/
- If you spread reads across all four prefixes evenly, you can achieve 22,000 requests per second for GET and HEAD

## Slide 304

S3 Performance
- Multi-Part upload:
- recommended for files > 100MB, must use for files > 5GB
- Can help parallelize uploads (speed up transfers)
- S3 Transfer Acceleration
- Increase transfer speed by transferring file to an AWS edge location which will forward the data to the S3 bucket in the target region
- Compatible with multi-part upload
Amazon S3
Parallel uploads
DivideIn parts
BIG file
S3 BucketAustralia
Edge LocationUSA
Fast (public www)Fast (private AWS)File in USA

## Slide 305

S3 Performance – S3 Byte-Range Fetches
- Parallelize GETs by requesting specific byte ranges
- Better resilience in case of failuresCan be used to speed up downloadsFile in S3Byte-range request for header(first XX bytes)header
File in S3
Part 1Part 2Part N…
Can be used to retrieve only partial data (for example the head of a file)
Requests in parallel

## Slide 306

S3 Batch Operations
- Perform bulk operations on existing S3 objects with a single request, example:
- Modify object metadata & properties
- Copy objects between S3 buckets
- Encrypt un-encrypted objects
- Modify ACLs, tags
- Restore objects from S3 Glacier
- Invoke Lambda function to perform custom action on each object
- A job consists of a list of objects, the action to perform, and optional parameters
- S3 Batch Operations manages retries, tracks progress, sends completion notifications, generate reports …
- Y ou can use S3 Inventory to get object list and use Athena to query and filter your objects
S3 Inventory
Objects List Report
filtered list
filter
S3 Batch Operations
operation+parametersUser
Processed Objects
…
Athena

## Slide 307

S3 – Storage Lens
- Understand, analyze, and optimize storage across entire AWS Organization
- Discover anomalies, identify cost efficiencies, and apply data protection best practices across entire AWS Organization (30 days usage & activity metrics)
- Aggregate data for Organization, specific accounts, regions, buckets, or prefixes
- Default dashboard or create your own dashboards
- Can be configured to export metrics daily to an S3 bucket (CSV, Parquet)
S3 Storage Lens
OrganizationAccountsRegionsBucketsAggregateAnalyze(Dashboard)
Summary InsightsData ProtectionCost Efficiency
Configure Optimize

## Slide 308

Storage Lens – Default Dashboard
- Visualize summarized insights and trends for both free and advanced metrics
- Default dashboard shows Multi-Region and Multi-Account data
- Preconfigured by Amazon S3
- Can’t be deleted, but can be disabled
https://aws.amazon.com/blogs/aws/s3-storage-lens/https://aws.amazon.com/blogs/aws/s3-storage-lens/

## Slide 309

Storage Lens – Metrics
- Summary Metrics
- General insights about your S3 storage
- StorageBytes, ObjectCount…
- Use cases: identify the fastest-growing (or not used) buckets and prefixes
- Cost-Optimization Metrics
- Provide insights to manage and optimize your storage costs
- NonCurrentVersionStorageBytes, IncompleteMultipartUploadStorageBytes…
- Use cases: identify buckets with incomplete multipart uploaded older than 7 days, Identify which objects could be transitioned to lower-cost storage class

## Slide 310

Storage Lens – Metrics
- Data-Protection Metrics
- Provide insights for data protection features
- VersioningEnabledBucketCount, MFADeleteEnabledBucketCount, SSEKMSEnabledBucketCount, CrossRegionReplicationRuleCount…
- Use cases: identify buckets that aren’t following data-protection best practices
- Access-management Metrics
- Provide insights for S3 Object Ownership
- ObjectOwnershipBucketOwnerEnforcedBucketCount…
- Use cases: identify which Object Ownership settings your buckets use
- Event Metrics
- Provide insights for S3 Event Notifications
- EventNotificationEnabledBucketCount (identify which buckets have S3 Event Notifications configured)

## Slide 311

Storage Lens – Metrics
- Performance Metrics
- Provide insights for S3 Transfer Acceleration
- TransferAccelerationEnabledBucketCount (identify which buckets have S3 Transfer Acceleration enabled)
- Activity Metrics
- Provide insights about how your storage is requested
- AllRequests, GetRequests, PutRequests, ListRequests, BytesDownloaded…
- Detailed Status Code Metrics
- Provide insights for HTTP status codes
- 200OKStatusCount, 403ForbiddenErrorCount, 404NotFoundErrorCount…

## Slide 312

Storage Lens – Free vs. Paid
- Free Metrics
- Automatically available for all customers
- Contains around 28 usage metrics
- Data is available for queries for 14 days
- Advanced Metrics and Recommendations
- Additional paid metrics and features
- Advanced Metrics – Activity, Advanced Cost Optimization, Advanced Data Protection, Status Code
- CloudWatch Publishing – Access metrics in CloudWatch without additional charges
- Prefix Aggregation – Collect metrics at the prefix level
- Data is available for queries for 15 months

## Slide 313

Amazon S3 – Security

## Slide 314

Amazon S3 – Object Encryption
- Y ou can encrypt objects in S3 buckets using one of 4 methods
- Server-Side Encryption (SSE)
- Server-Side Encryption with Amazon S3-Managed Keys (SSE-S3) – Enabled by Default
- Encrypts S3 objects using keys handled, managed, and owned by AWS
- Server-Side Encryption with KMS Keys stored in AWS KMS (SSE-KMS)
- Leverage AWS Key Management Service (AWS KMS) to manage encryption keys
- Server-Side Encryption with Customer-Provided Keys (SSE-C)
- When you want to manage your own encryption keys
- Client-Side Encryption
- It’s important to understand which ones are for which situation for the exam

## Slide 315

Amazon S3 Encryption – SSE-S3
- Encryption using keys handled, managed, and owned by AWS
- Object is encrypted server-side
- Encryption type is AES-256
- Must set header "x-amz-server-side-encryption": "AES256"
- Enabled by default for new buckets & new objectsAmazon S3
User
uploadHTTP(S) + Header
S3 Bucket
+Object
S3 Owned Key
Encryption

## Slide 316

Amazon S3 Encryption – SSE-KMS
- Encryption using keys handled and managed by AWS KMS (Key Management Service)
- KMS advantages: user control + audit key usage using CloudTrail
- Object is encrypted server side
- Must set header "x-amz-server-side-encryption": "aws:kms"Amazon S3
User
uploadHTTP(S) + Header
S3 Bucket
+Object
KMS Key
AWS KMS
Encryption

## Slide 317

SSE-KMS Limitation
- If you use SSE-KMS, you may be impacted by the KMS limits
- When you upload, it calls the GenerateDataKey KMS API
- When you download, it calls the Decrypt KMS API
- Count towards the KMS quota per second (5500, 10000, 30000 req/s based on region)
- Yo u c a n r e q u e s t a q u o t a i n c r e a s e u s i n g t h e Service Quotas Console
KMS KeyS3 Bucket
Users
Upload / downloadSSE-KMS
API call

## Slide 318

Amazon S3 Encryption – SSE-C
- Server-Side Encryption using keys fully managed by the customer outside of AWS
- Amazon S3 does NOT store the encryption key you provide
- HTTPS must be used
- Encryption key must provided in HTTP headers, for every HTTP request madeAmazon S3
User
uploadHTTPS ONLY+ Key in Header
S3 Bucket
+Object
Client-Provided Key
+ Encryption

## Slide 319

Amazon S3 Encryption – Client-Side Encryption
Amazon S3
S3 Bucket
- Use client libraries such as Amazon S3 Client-Side Encryption Library
- Clients must encrypt data themselves before sending to Amazon S3
- Clients must decrypt data themselves when retrieving from Amazon S3
- Customer fully manages the keys and encryption cycle
+File
Client Key
File(encrypted)
uploadHTTP(S)Encryption

## Slide 320

Amazon S3 – Encryption in transit (SSL/TLS)
- Encryption in flight is also called SSL/TLS
- Amazon S3 exposes two endpoints:
- HTTP Endpoint – non encrypted
- HTTPS Endpoint – encryption in flight
- HTTPS is recommended
- HTTPS is mandatory for SSE-C
- Most clients would use the HTTPS endpoint by default

## Slide 321

Amazon S3 – Force Encryption in Transitaws:SecureTransport
Account B
S3 Bucket(my-bucket)
Bucket Policy
User
User
http
https

## Slide 322

Amazon S3 – Default Encryption vs. Bucket Policies
- SSE-S3 encryption is automatically applied to new objects stored in S3 bucket
- Optionally, you can “force encryption” using a bucket policy and refuse any API call to PUT an S3 object without encryption headers (SSE-KMS or SSE-C)
- Note: Bucket Policies are evaluated before “Default Encryption”

## Slide 323

What is CORS?
- Cross-Origin Resource Sharing (CORS)
- Origin = scheme (protocol) + host (domain) + port
- example: https://www.example.com (implied port is 443 for HTTPS, 80 for HTTP)
- Web Browser based mechanism to allow requests to other origins while visiting the main origin
- Same origin: http://example.com/app1 & http://example.com/app2
- Different origins: http://www.example.com & http://other.example.com
- The requests won’t be fulfilled unless the other origin allows for the requests, using CORS Headers (example: Access-Control-Allow-Origin)

## Slide 324

What is CORS?
Web Server(Origin)https://www.example.com
Web Server(Cross-Origin)https://www.other.com
Web BrowserHTTPS Request
OPTIONS /Host: www.other.comOrigin: https://www.example.comPreflight Request
Access-Control-Allow-Origin: https://www.example.comAccess-Control-Allow-Methods: GET, PUT, DELETEPreflight Response
GET /Host: www.other.comOrigin: https://www.example.comCORS Headers received already by the OriginThe Web Browser can make requests

## Slide 325

Amazon S3 – CORS
- If a client makes a cross-origin request on our S3 bucket, we need to enable the correct CORS headers
- It’s a popular exam question
- Yo u c a n a l l ow fo r a s p e c i fi c o r i g i n o r fo r * ( a l l o r i g i n s )
Web Browser
S3 Bucket(my-bucket-html)(Static Website Enabled)
S3 Bucket(my-bucket-assets)(Static Website Enabled)
GET /index.htmlHost: http://my-bucket-html.s3-website.us-west-2.amazonaws.com
index.htmlGET /images/coffee.jpgHost: http://my-bucket-assets.s3-website.us-west-2.amazonaws.comOrigin: http://my-bucket-html.s3-website.us-west-2.amazonaws.com
Access-Control-Allow-Origin: http://my-bucket-html.s3-website.us-west-2.amazonaws.com

## Slide 326

Amazon S3 – MFA Delete
- MFA (Multi-Factor Authentication) – force users to generate a code on a device (usually a mobile phone or hardware) before doing important operations on S3
- MFA will be required to:
- Permanently delete an object version
- Suspend Versioning on the bucket
- MFA won’t be required to:
- Enable Versioning
- List deleted versions
- To u s e M FA D e l e t e , Versioning must be enabled on the bucket
- Only the bucket owner (root account) can enable/disable MFA Delete
Google Authenticator
MFA Hardware Device

## Slide 327

S3 Access Logs
- For audit purpose, you may want to log all access to S3 buckets
- Any request made to S3, from any account, authorized or denied, will be logged into another S3 bucket
- That data can be analyzed using data analysis tools…
- The target logging bucket must be in the same AWS region
- The log format is at: https://docs.aws.amazon.com/AmazonS3/latest/dev/LogFormat.html
My-bucket
Logging Bucket
requests
Log all requests

## Slide 328

S3 Access Logs: Warning
- Do not set your logging bucket to be the monitored bucket
- It will create a logging loop, and your bucket will grow exponentially
App Bucket &Logging Bucket
Logging loop
PutObject
Do not try this at home J

## Slide 329

Amazon S3 – Pre-Signed URLs
- Generate pre-signed URLs using the S3 Console, AWS CLI or SDK
- URL Expiration
- S3 Console – 1 min up to 720 mins (12 hours)
- AWS CLI – configure expiration with --expires-in parameter in seconds (default 3600 secs, max. 604800 secs ~ 168 hours)
- Users given a pre-signed URL inherit the permissions of the user that generated the URL for GET / PUT
- Examples:
- Allow only logged-in users to download a premium video from your S3 bucket
- Allow an ever-changing list of users to download files by generating URLs dynamically
- Allow temporarily a user to upload a file to a precise location in your S3 bucket
S3 Bucket(Private)
Owner
User generatepre-signed URL
URL
URL
URL

## Slide 330

S3 Glacier Vault Lock
- Adopt a WORM (Write Once Read Many) model
- Create a Vault Lock Policy
- Lock the policy for future edits (can no longer be changed or deleted)
- Helpful for compliance and data retention
Object
Vault Lock PolicyObject can’t be deleted

## Slide 331

S3 Object Lock (versioning must be enabled)
- Adopt a WORM (Write Once Read Many) model
- Block an object version deletion for a specified amount of time
- Retention mode - Compliance:
- Object versions can't be overwritten or deleted by any user, including the root user
- Objects retention modes can't be changed, and retention periods can't be shortened
- Retention mode - Governance:
- Most users can't overwrite or delete an object version or alter its lock settings
- Some users have special permissions to change the retention or delete the object
- Retention Period: protect the object for a fixed period, it can be extended
- Legal Hold:
- protect the object indefinitely, independent from retention period
- can be freely placed and removed using the s3:PutObjectLegalHold IAM permission

## Slide 332

S3 – Access Points
S3 Bucket
/finance/…
/sales/…
Simple BucketPolicy
FinanceAccess Point
SalesAccess Point
AnalyticsAccess Point
PolicyGrant R/W to/finance prefix
PolicyGrant R/W to/sales prefix
PolicyGrant R toentire bucket
Users(Finance)Users(Sales)Users(Analytics)
- Access Points simplify security management for S3 Buckets
- Each Access Point has:
- its own DNS name (Internet Origin or VPC Origin)
- an access point policy (similar to bucket policy) – manage security at scale

## Slide 333

S3 – Access Points – VPC Origin
- We can define the access point to be accessible only from within the VPC
- Yo u mu s t c r e a t e a V P C Endpoint to access the Access Point (Gateway or Interface Endpoint)
- The VPC Endpoint Policy must allow access to the target bucket and Access Point
S3 Bucket
EC2 InstanceAccess PointVPC Origin
VPC Endpoint
VPC
Endpoint PolicyAccess Point PolicyBucketPolicy

## Slide 334

S3 Object Lambda
- Use AWS Lambda Functions to change the object before it is retrieved by the caller application
- Only one S3 bucket is needed, on top of which we create S3 Access Point and S3 Object Lambda Access Points.
- Use Cases:
- Redacting personally identifiable information for analytics or non-production environments.
- Converting across data formats, such as converting XML to JSON.
- Resizing and watermarking images on the fly using caller-specific details, such as the user who requested the object.
AWS Cloud
E-CommerceApplication
AnalyticsApplication
MarketingApplication
S3 Bucket
SupportingS3 Access Point
Customer LoyaltyDatabase
RedactingLambda Function
EnrichingLambda Function
S3 Object LambdaAccess Point
S3 Object LambdaAccess Point
OriginalObject
RedactedObject
EnrichedObject

## Slide 335

CloudFront & Global Accelerator

## Slide 336

Amazon CloudFront
- Content Delivery Network (CDN)
- Improves read performance, content is cached at the edge
- Improves users experience
- Hundreds of Points of Presence globally (edge locations, caches)
- DDoS protection (because worldwide), integration with Shield, AWS Web Application Firewall
Source: https://aws.amazon.com/cloudfront/features/?nc=sn&loc=2

## Slide 337

CloudFront – Origins
- S3 bucket
- For distributing files and caching them at the edge
- For uploading files to S3 through CloudFront
- Secured using Origin Access Control (OAC)
- VPC Origin
- For applications hosted in VPC private subnets
- Private Application Load Balancer / Network Load Balancer / EC2 Instances
- Custom Origin (HTTP)
- S3 website (must first enable the bucket as a static S3 website)
- Any public HTTP backend you want (example: Public ALB)

## Slide 338

CloudFront at a high level
CloudFront Edge Location
Forward Request to your Origin
Local Cache
S3
HTTPor
Origin
Client
GET /beach.jpg?size=300x300 HTTP/1.1User-Agent: Mozilla/4.0 (compatible; MSIE5.01; Windows NT)Host: www.example.comAccept-Encoding: gzip, deflate

## Slide 339

CloudFront – S3 as an Origin
EdgeLos Angeles
EdgeMumbai
EdgeMelbourne
Origin (S3 bucket)
Public www
Public wwwEdgeSão Paulo
Origin Access Control+ S3 bucket policy
OAC
Private AWS
Private AWS
Private AWS
Private AWS
AWS Cloud

## Slide 340

CloudFront vs S3 Cross Region Replication
- CloudFront:
- Global Edge network
- Files are cached for a TTL (maybe a day)
- Great for static content that must be available everywhere
- S3 Cross Region Replication:
- Must be setup for each region you want replication to happen
- Files are updated in near real-time
- Read only
- Great for dynamic content that needs to be available at low-latency in few regions

## Slide 341

CloudFront – ALB or EC2 as an originUsing VPC Origins
- Allows you to deliver content from your applications hosted in your VPC private subnets (no need to expose them on the Internet)
- Deliver traffic to private:
- Application Load Balancer
- Network Load Balancer
- EC2 Instances
UsersCloudFront
Edge Location
Private Subnet
VPC
VPCOrigin
Application Load BalancerNetwork Load BalancerEC2 Instance

## Slide 342

CloudFront – ALB or EC2 as an originUsing Public Networkhttp://d7uri8nf7uskq.cloudfront.net/tools/list-cloudfront-ips
Edge LocationPublic IPs
Application Load BalancerMust be Public
EC2 InstancesCan be Private
Allow Security Group of Load BalancerAllow Public IP of Edge LocationsSecurity groupSecurity group
Edge LocationEC2 InstancesMust be Public
Allow Public IP of Edge LocationsSecurity group

## Slide 343

CloudFront Geo Restriction
- Yo u c a n r e s t r i c t w h o c a n a c c e s s yo u r d i s t r i bu t i o n
- Allowlist: Allow your users to access your content only if they're in one of the countries on a list of approved countries.
- Blocklist: Prevent your users from accessing your content if they're in one of the countries on a list of banned countries.
- The “country” is determined using a 3rd party Geo-IP database
- Use case: Copyright Laws to control access to content

## Slide 344

CloudFront – Cache Invalidations
- In case you update the back-end origin, CloudFront doesn’t know about it and will only get the refreshed content after the TTL has expired
- However, you can force an entire or partial cache refresh (thus bypassing the TTL) by performing a CloudFront Invalidation
- Yo u c a n i nv a l i d a t e a l l fi l e s ( * ) o r a special path (/images/*)
CloudFront
Edge Location
index.html/images/
Cache
Edge Location
Invalidate- /index.html- /images/*
invalidate
S3 Bucket(origin)
index.html/images/
Cache
GET /index.html
update files

## Slide 345

Global users for our application
- Yo u h ave d e p l oye d a n application and have global users who want to access it directly.
- They go over the public internet, which can add a lot of latency due to many hops
- We wish to go as fast as possible through AWS network to minimize latency
America
Australia
Europe
India
Public ALB
hops

## Slide 346

Unicast IP vs Anycast IP
- Unicast IP: one server holds one IP address
- Anycast IP: all servers hold the same IP address and the client is routed to the nearest one
12.34.56.7898.76.54.32
Client
12.34.56.7812.34.56.78
Client

## Slide 347

AWS Global Accelerator
- Leverage the AWS internal network to route to your application
- 2 Anycast IP are created for your application
- The Anycast IP send traffic directly to Edge Locations
- The Edge locations send the traffic to your application
America
Australia
Europe
IndiaPublic ALB
Edge location
Private AWS

## Slide 348

AWS Global Accelerator
- Works with Elastic IP , EC2 instances, ALB, NLB, public or private
- Consistent Performance
- Intelligent routing to lowest latency and fast regional failover
- No issue with client cache (because the IP doesn’t change)
- Internal AWS network
- Health Checks
- Global Accelerator performs a health check of your applications
- Helps make your application global (failover less than 1 minute for unhealthy)
- Great for disaster recovery (thanks to the health checks)
- Security
- only 2 external IP need to be whitelisted
- DDoS protection thanks to AWS Shield

## Slide 349

AWS Global Accelerator vs CloudFront
- They both use the AWS global network and its edge locations around the world
- Both services integrate with AWS Shield for DDoS protection.
- CloudFront
- Improves performance for both cacheable content (such as images and videos)
- Dynamic content (such as API acceleration and dynamic site delivery)
- Content is served at the edge
- Global Accelerator
- Improves performance for a wide range of applications over TCP or UDP
- Proxying packets at the edge to applications running in one or more AWS Regions.
- Good fit for non-HTTP use cases, such as gaming (UDP), IoT (MQTT), or Voice over IP
- Good for HTTP use cases that require static IP addresses
- Good for HTTP use cases that required deterministic, fast regional failover

## Slide 350

AWS Storage Extras

## Slide 351

AWS Snowball
- Highly-secure, portable devices to collect and process data at the edge, and migrate data into and out of AWS
- Helps migrate up to Petabytes of data
Snowball Edge
DeviceComputeMemoryStorage (SSD)Snowball Edge Storage Optimized104 vCPUs416 GB210 TBSnowball Edge Compute Optimized104 vCPUs416 GB28 TB

## Slide 352

Data Migrations with SnowballTime to Transfer100 Mbps1Gbps10Gbps10 TB12 days30 hours3 hours100 TB124 days12 days30 hours1 PB3 years124 days12 days
Challenges:
- Limited connectivity
- Limited bandwidth
- High network cost
- Shared bandwidth (can’t maximize the line)
- Connection stabilityAWS Snowball: offline devices to perform data migrationsIf it takes more than a week to transfer over the network, use Snowball devices!

## Slide 353

Diagrams
- Direct upload to S3:
- With Snowball:clientAmazon S3 bucket
import/exportAmazon S3 bucketAWS SnowballclientAWS Snowball
www: 10Gbit/s
ship

## Slide 354

What is Edge Computing?
- Process data while it’s being created on an edge location
- A truck on the road, a ship on the sea, a mining station underground...
- These locations may have limited internet and no access to computing power
- We setup a Snowball Edge device to do edge computing
- Snowball Edge Compute Optimized (dedicated for that use case) & Storage Optimized
- Run EC2 Instances or Lambda functions at the edge
- Use cases: preprocess data, machine learning, transcoding media

## Slide 355

Solution Architecture: Snowball into Glacier
- Snowball cannot import to Glacier directly
- Y ou must use Amazon S3 first, in combination with an S3 lifecycle policy
Snowball
Amazon S3
Amazon Glacier
S3 lifecycle policyimport

## Slide 356

Amazon FSx – Overview
- Launch 3rd party high-performance file systems on AWS
- Fully managed service
FSx for Lustre
FSx for Windows File Server
FSx for NetApp ONTAP
FSx for OpenZFS

## Slide 357

Amazon FSx for Windows (File Server)
- FSx for Windows is a fully managed Windows file system share drive
- Supports SMB protocol & Windows NTFS
- Microsoft Active Directory integration, ACLs, user quotas
- Can be mounted on Linux EC2 instances
- Supports Microsoft's Distributed File System (DFS) Namespaces (group files across multiple FS)
- Scale up to 10s of GB/s, millions of IOPS, 100s PB of data
- Storage Options:
- SSD – latency sensitive workloads (databases, media processing, data analytics, …)
- HDD – broad spectrum of workloads (home directory, CMS, …)
- Can be accessed from your on-premises infrastructure (VPN or Direct Connect)
- Can be configured to be Multi-AZ (high availability)
- Data is backed-up daily to S3

## Slide 358

Amazon FSx for Lustre
- Lustre is a type of parallel distributed file system, for large-scale computing
- The name Lustre is derived from “Linux” and “cluster
- Machine Learning, High Performance Computing (HPC)
- Video Processing, Financial Modeling, Electronic Design Automation
- Scales up to 100s GB/s, millions of IOPS, sub-ms latencies
- Storage Options:
- SSD – low-latency, IOPS intensive workloads, small & random file operations
- HDD – throughput-intensive workloads, large & sequential file operations
- Seamless integration with S3
- Can “read S3” as a file system (through FSx)
- Can write the output of the computations back to S3 (through FSx)
- Can be used from on-premises servers (VPN or Direct Connect)

## Slide 359

FSx Lustre - File System Deployment Options
- Scratch File System
- Tempor ar y stor age
- Data is not replicated (doesn’t persist if file server fails)
- High burst (6x faster, 200MBps per TiB)
- Usage: short-term processing, optimize costs
- Persistent File System
- Long-term storage
- Data is replicated within same AZ
- Replace failed files within minutes
- Usage: long-term processing, sensitive data
RegionAvailability Zone 1
Availability Zone 2
Compute instances Compute instances
S3 bucket(optional data repository)FSx For Lustre(Scratch file system)RegionAvailability Zone 1
Availability Zone 2
Compute instances Compute instances
S3 bucket(optional data repository)FSx For Lustre(Persistent file system)
ENI
ENI

## Slide 360

Amazon FSx for NetApp ONTAP
- Managed NetApp ONTAP on AWS
- File System compatible with NFS, SMB, iSCSI protocol
- Move workloads running on ONTAP or NAS to AWS
- Works with:
- Linux
- Windows
- MacOS
- VMware Cloud on AWS
- Amazon Workspaces & AppStream 2.0
- Amazon EC2, ECS and EKS
- Storage shrinks or grows automatically
- Snapshots, replication, low-cost, compression and data de-duplication
- Point-in-time instantaneous cloning (helpful for testing new workloads)
Amazon FSx forNetApp ONTAP FS
VMware Cloudon AWS
AmazonAppStream 2.0AmazonWorkSpaces
EC2
ECSEKS
On-premisesServer
NFS, SMB, iSCSI

## Slide 361

Amazon FSx for OpenZFS
- Managed OpenZFS file system on AWS
- File System compatible with NFS (v3, v4, v4.1, v4.2)
- Move workloads running on ZFS to AWS
- Works with:
- Linux
- Windows
- MacOS
- VMware Cloud on AWS
- Amazon Workspaces & AppStream 2.0
- Amazon EC2, ECS and EKS
- Up to 1,000,000 IOPS with < 0.5ms latency
- Snapshots, compression and low-cost
- Point-in-time instantaneous cloning (helpful for testing new workloads)
Amazon FSxfor OpenZFS
VMware Cloudon AWS
AmazonAppStream 2.0AmazonWorkSpaces
EC2
ECSEKS
On-premisesServer
NFS (v3, v4, v4.1, v4.2)

## Slide 362

Hybrid Cloud for Storage
- AWS is pushing for ”hybrid cloud”
- Part of your infrastructure is on the cloud
- Part of your infrastructure is on-premises
- This can be due to
- Long cloud migrations
- Security requirements
- Compliance requirements
- IT strategy
- S3 is a proprietary storage technology (unlike EFS / NFS), so how do you expose the S3 data on-premises?
- AWS Storage Gateway!

## Slide 363

AWS Storage Cloud Native OptionsBlock
FileObject
Amazon EBSEC2 InstanceStoreAmazon EFSAmazon FSxAmazon S3Amazon Glacier

## Slide 364

AWS Storage Gateway
- Bridge between on-premises data and cloud data
- Use cases:
- disaster recovery
- backup & restore
- tiered storage
- on-premises cache & low-latency files access
- Types of Stor age Gateway:
- S3 File Gateway
- Volume Gateway
- Tape Gateway
Storage Gateway

## Slide 365

Amazon S3 File Gateway
- Configured S3 buckets are accessible using the NFS and SMB protocol
- Most recently used data is cached in the file gateway
- Supports S3 Standard, S3 Standard IA, S3 One Zone A, S3 Intelligent Tiering
- Transition to S3 Glacier using a Lifecycle Policy
- Bucket access using IAM roles for each File Gateway
- SMB Protocol has integration with Active Directory (AD) for user authenticationCorporate Data Center
AWS Cloud .
HTTPSApplicationServerS3 FileGateway
S3 StandardS3 Standard-IAS3 One Zone-IAS3 Intelligent-TieringS3 GlacierNFS or SMB
Lifecycle policy

## Slide 366

Volume Gateway
- Block storage using iSCSI protocol backed by S3
- Backed by EBS snapshots which can help restore on-premises volumes!
- Cached volumes: low latency access to most recent data
- Stored volumes: entire dataset is on premise, scheduled backups to S3Corporate Data Center
AWS Cloud
Region
S3 Bucket
Amazon EBSSnapshots
ApplicationServer
iSCSI
Volume GatewayHTTPS

## Slide 367

Tape Gateway
- Some companies have backup processes using physical tapes (!)
- With Tape Gateway, companies use the same processes but, in the cloud
- Virtual Tape Library (VTL) backed by Amazon S3 and Glacier
- Back up data using existing tape-based processes (and iSCSI interface)
- Works with leading backup software vendors Corporate Data Center
AWS Cloud
Region
Virtual Tapesstored inAmazon S3
Archived Tapesstored inAmazon Glacier
BackupServer
iSCSI HTTPS
TapeGateway
Media ChangerTape Drive

## Slide 368

AWS Storage Gateway
On-Premises
File Gatewaylocal cache
Volume Gatewaylocal cache
Tape Gatewaylocal cache
Application Server
Backup Application
User/group file shares
AWS Cloud
NFS/SMB
iSCSI
iSCSI VTL
Gateway Deployment OptionsVM(VMware, Hyper-V, KVM)
Amazon S3excluding Glacier &Glacier Deep Archive
Amazon S3
Amazon S3Tape Library
Any S3 Storage ClassIncluding Glacier
AWS EBS
Tape ArchiveGlacier &Glacier Deep Archive
Eject from backup applicationStorage GatewayEncryption in TransitInternet or Direct Connect

## Slide 369

AWS Transfer Family
- A fully-managed service for file transfers into and out of Amazon S3 or Amazon EFS using the FTP protocol
- Supported Protocols
- AWS Transfer for FTP (File Transfer Protocol (FTP))
- AWS Transfer for FTPS (File Transfer Protocol over SSL (FTPS))
- AWS Transfer for SFTP (Secure File Transfer Protocol (SFTP))
- Managed infrastructure, Scalable, Reliable, Highly Available (multi-AZ)
- Pay per provisioned endpoint per hour + data transfers in GB
- Store and manage users’ credentials within the service
- Integrate with existing authentication systems (Microsoft Active Directory, LDAP , Okta, Amazon Cognito, custom)
- Usage: sharing files, public datasets, CRM, ERP , …

## Slide 370

AWS Transfer Family
AWS Transfer Family
AWS Transfer for SFTP
AWS Transfer for FTPS
AWS Transfer for FTP(only within VPC)
Amazon S3
Amazon EFS
authenticateMS Active DirectoryLDAP …
Route 53(optional)Users(FTP client)
IAM Role

## Slide 371

AWS DataSync
- Move large amount of data to and from
- On-premises / other cloud to AWS (NFS, SMB, HDFS, S3 API…) – needs agent
- AWS to AWS (different storage services) – no agent needed
- Can synchronize to:
- Amazon S3 (any storage classes – including Glacier)
- Amazon EFS
- Amazon FSx (Windows, Lustre, NetApp, OpenZFS...)
- Replication tasks can be scheduled hourly, daily, weekly
- File permissions and metadata are preserved (NFS POSIX, SMB…)
- One agent task can use 10 Gbps, can setup a bandwidth limit

## Slide 372

AWS DataSyncNFS / SMB to AWS (S3, EFS, FSx…)
On-Premises
NFS or SMBServerAWS DataSyncAgent
NFS or SMB
Region
AWSDataSync
AWS Storage Resources
S3 Standard
S3 Intelligent-TieringS3 Standard-IA
S3 OneZone-IAS3 GlacierS3 GlacierDeep Archive
AWS EFS
Amazon FSx
TLS
AWS Snowcone(agent pre-installed)

## Slide 373

AWS DataSyncTransfer between AWS storage ser vices
AWS DataSynccopy data and metadatabetween AWS Storage Services
Amazon FSx
Amazon EFS
Amazon S3
Amazon FSx
Amazon EFS
Amazon S3

## Slide 374

Storage Comparison
- S3: Object Storage
- S3 Glacier: Object Archival
- EBS volumes: Network storage for one EC2 instance at a time
- Instance Storage: Physical storage for your EC2 instance (high IOPS)
- EFS: Network File System for Linux instances, POSIX filesystem
- FSx for Windows: Network File System for Windows servers
- FSx for Lustre: High Performance Computing Linux file system
- FSx for NetApp ONTAP: High OS Compatibility
- FSx for OpenZFS: Managed ZFS file system
- Storage Gateway: S3 & FSx File Gateway, Volume Gateway (cache & stored), Tape Gateway
- Transfer Family: FTP , FTPS, SFTP interface on top of Amazon S3 or Amazon EFS
- DataSync: Schedule data sync from on-premises to AWS, or AWS to AWS
- Snowcone / Snowball / Snowmobile: to move large amount of data to the cloud, physically
- Database: for specific workloads, usually with indexing and querying

## Slide 375

AWS Integration & MessagingSQS, SNS & Kinesis

## Slide 376

Section Introduction
- When we start deploying multiple applications, they will inevitably need to communicate with one another
- There are two patterns of application communication1) Synchronous communications(application to application)2) Asynchronous / Event based (application to queue to application)
Buying ServiceShipping ServiceBuying ServiceShipping ServiceQueue

## Slide 377

Section Introduction
- Synchronous between applications can be problematic if there are sudden spikes of traffic
- What if you need to suddenly encode 1000 videos but usually it’s 10?
- In that case, it’s better to decouple your applications,
- using SQS: queue model
- using SNS: pub/sub model
- using Kinesis: real-time streaming model
- These services can scale independently from our application!

## Slide 378

Amazon SQSWhat’s a queue?Producer
Producer
ProducerSQS Queue
Consumer
Consumer
Consumer
Consumer
Send messagesPoll messages

## Slide 379

Amazon SQS – Standard Queue
- Oldest offering (over 10 years old)
- Fully managed service, used to decouple applications
- Attributes:
- Unlimited throughput, unlimited number of messages in queue
- Default retention of messages: 4 days, maximum of 14 days
- Low latency (<10 ms on publish and receive)
- Limitation of 1,024 KB per message sent
- Can have duplicate messages (at least once delivery, occasionally)
- Can have out of order messages (best effort ordering)

## Slide 380

SQS – Producing Messages
- Produced to SQS using the SDK (SendMessage API)
- The message is persisted in SQS until a consumer deletes it
- Message retention: default 4 days, up to 14 days
- Example: send an order to be processed
- Order id
- Customer id
- Any attributes you want
- SQS standard: unlimited throughput
Sent to SQS
MessageUp to 1024 KB

## Slide 381

SQS – Consuming Messages
- Consumers (running on EC2 instances, servers, or AWS Lambda)…
- Poll SQS for messages (receive up to 10 messages at a time)
- Process the messages (example: insert the message into an RDS database)
- Delete the messages using the DeleteMessage API
ConsumerPoll / Receive messages
DeleteMessage
insert

## Slide 382

SQS – Multiple EC2 Instances Consumers
SQS Queue
poll
- Consumers receive and process messages in parallel
- At least once delivery
- Best-effort message ordering
- Consumers delete messages after processing them
- We can scale consumers horizontally to improve throughput of processing

## Slide 383

SQS with Auto Scaling Group (ASG)
SQS Queue
CloudWatch Metric – Queue LengthApproximateNumberOfMessages
Auto Scaling Group
EC2 InstancesPoll for messages
CloudWatch Alarm
scaleAlarm for breach

## Slide 384

SQS to decouple between application tiers
SQS Queue(infinitely scalable)
Auto-Scaling
Back-end processingapplication
Front-end web apprequestsSendMessage
Auto-Scaling
ReceiveMessages

## Slide 385

Amazon SQS - Security
- Encryption:
- In-flight encryption using HTTPS API
- At-rest encryption using KMS keys
- Client-side encryption if the client wants to perform encryption/decryption itself
- Access Controls: IAM policies to regulate access to the SQS API
- SQS Access Policies (similar to S3 bucket policies)
- Useful for cross-account access to SQS queues
- Useful for allowing other services (SNS, S3…) to write to an SQS queue

## Slide 386

SQS – Message Visibility Timeout
- After a message is polled by a consumer, it becomes invisible to other consumers
- By default, the “message visibility timeout” is 30 seconds
- That means the message has 30 seconds to be processed
- After the message visibility timeout is over, the message is “visible” in SQS
Time
ReceiveMessageRequest
Visibility timeout
Message returned
ReceiveMessageRequest
Not returned
ReceiveMessageRequest
Not returned
ReceiveMessageRequest
Message returned (again)

## Slide 387

SQS – Message Visibility Timeout
Time
ReceiveMessageRequest
Visibility timeout
Message returned
ReceiveMessageRequest
Not returned
ReceiveMessageRequest
Not returned
ReceiveMessageRequest
Message returned (again)
- If a message is not processed within the visibility timeout, it will be processed twice
- A consumer could call the ChangeMessageVisibility API to get more time
- If visibility timeout is high (hours), and consumer crashes, re-processing will take time
- If visibility timeout is too low (seconds), we may get duplicates

## Slide 388

Amazon SQS - Long Polling
- When a consumer requests messages from the queue, it can optionally “wait” for messages to arrive if there are none in the queue
- This is called Long Polling
- LongPolling decreases the number of API calls made to SQS while increasing the efficiency and reducing latency of your application
- The wait time can be between 1 sec to 20 sec (20 sec preferable)
- Long Polling is preferable to Short Polling
- Long polling can be enabled at the queue level or at the API level using WaitTimeSecondsConsumer
SQS Queue
poll
message

## Slide 389

Amazon SQS – FIFO Queue
- FIFO = First In First Out (ordering of messages in the queue)Producer ConsumerSend messagesPoll messages
1234 1234
- Limited throughput: 300 msg/s without batching, 3000 msg/s with
- Exactly-once send capability (by removing duplicates using Deduplication ID)
- Messages are processed in order by the consumer
- Ordering by Message Group ID (all messages in the same group are ordered) – mandatory parameter

## Slide 390

SQS with Auto Scaling Group (ASG)
SQS Queue
CloudWatch Metric – Queue LengthApproximateNumberOfMessages
Auto Scaling Group
EC2 InstancesPoll for messages
CloudWatch Alarm
scaleAlarm for breach

## Slide 391

If the load is too big, some transactions may be lost
Applicationrequests
Auto-Scaling
Insert transactions
Amazon RDS
Amazon Aurora
Amazon DynamoDB

## Slide 392

SQS as a buffer to database writes
SQS Queue(infinitely scalable)
Auto-Scaling
Dequeue message
Enqueue messagerequestsSendMessage
Auto-Scaling
ReceiveMessages
insert

## Slide 393

SQS to decouple between application tiers
SQS Queue(infinitely scalable)
Auto-Scaling
Back-end processingapplication
Front-end web apprequestsSendMessage
Auto-Scaling
ReceiveMessages

## Slide 394

Amazon SNS
- What if you want to send one message to many receivers?
Buying Service
Email notificationFraud ServiceShipping Service
SQS Queue
Email notificationFraud ServiceShipping Service
SQS Queue
SNS Topic
Buying Service
Direct integrationPub / Sub

## Slide 395

Amazon SNS
- The “event producer” only sends message to one SNS topic
- As many “event receivers” (subscriptions) as we want to listen to the SNS topic notifications
- Each subscriber to the topic will get all the messages (note: new feature to filter messages)
- Up to 12,500,000 subscriptions per topic
- 100,000 topics limit
SNS
Subscribers
publish
SQSLambdaKinesis DataFirehose
HTTP(S)EndpointsSMS &Mobile Notifications
Emails

## Slide 396

SNS integrates with a lot of AWS services
- Many AWS services can send data directly to SNS for notifications
SNS
CloudWatch Alarms
S3 Bucket(Events)Auto Scaling Group(Notifications)
CloudFormation(State Changes)
AWS Budgets
Lambda
AWS DMS(New Replic)
DynamoDB
RDS Events
publish
……
…

## Slide 397

Amazon SNS – How to publish
- To p i c P u bl i s h ( u s i n g t h e S D K )
- Create a topic
- Create a subscription (or many)
- Publish to the topic
- Direct Publish (for mobile apps SDK)
- Create a platform application
- Create a platform endpoint
- Publish to the platform endpoint
- Works with Google GCM, Apple APNS, Amazon ADM…

## Slide 398

Amazon SNS – Security
- Encryption:
- In-flight encryption using HTTPS API
- At-rest encryption using KMS keys
- Client-side encryption if the client wants to perform encryption/decryption itself
- Access Controls: IAM policies to regulate access to the SNS API
- SNS Access Policies (similar to S3 bucket policies)
- Useful for cross-account access to SNS topics
- Useful for allowing other services ( S3…) to write to an SNS topic

## Slide 399

SNS + SQS: Fan Out
- Push once in SNS, receive in all SQS queues that are subscribers
- Fully decoupled, no data loss
- SQS allows for: data persistence, delayed processing and retries of work
- Ability to add more SQS subscribers over time
- Make sure your SQS queue access policy allows for SNS to write
- Cross-Region Delivery: works with SQS Queues in other regions
Fraud Service
Shipping ServiceSNS TopicBuying ServiceSQS Queue
SQS Queue

## Slide 400

Application: S3 Events to multiple queues
- For the same combination of: event type (e.g. object create) and prefix (e.g. images/) you can only have one S3 Event rule
- If you want to send the same S3 event to many SQS queues, use fan-out
SNS Topic
SQS Queues
Amazon S3
eventsS3 Object created…
Lambda Function
Fan-out

## Slide 401

Application: SNS to Amazon S3 through Kinesis Data Firehose
SNS TopicBuying Service
Kinesis DataFirehose
Amazon S3
- SNS can send to Kinesis and therefore we can have the following solutions architecture:
Any supported KDFDestination

## Slide 402

Amazon SNS – FIFO T opic
- FIFO = First In First Out (ordering of messages in the topic)Producer SubscribersSQS FIFOSend messagesReceive messages1234 1234
- Similar features as SQS FIFO:
- Ordering by Message Group ID (all messages in the same group are ordered)
- Deduplication using a Deduplication ID or Content Based Deduplication
- Can have SQS Standard and FIFO queues as subscribers
- Limited throughput (same throughput as SQS FIFO)

## Slide 403

SNS FIFO + SQS FIFO: Fan Out
- In case you need fan out + ordering + deduplicationFraud Service
Shipping ServiceSNS FIFO TopicBuying ServiceSQS FIFO Queue
SQS FIFO Queue

## Slide 404

SNS – Message Filtering
- JSON policy used to filter messages sent to SNS topic’s subscriptions
- If a subscription doesn’t have a filter policy, it receives every message
SNS TopicBuying Service
New transactionOrder: 1036Product: PencilQty: 4State: Placed
SQS Queue(Placed orders)
Email Subscription(Cancelled orders)SQS Queue(Declined orders)
State: PlacedFilter Policy
State: DeclinedFilter PolicyState: CancelledFilter Policy
SQS Queue(All)
SQS Queue(Cancelled orders)

## Slide 405

Amazon Kinesis Data Streams
- Collect and store streaming data in real-time
Click Streams
IoT devices
Metrics & Logs
Amazon KinesisData Streams
Real-time data Consumers
ProducersApplicationsKinesis Agent
Lambda
Application
AmazonData Firehose
Managed Service for Apache Flink

## Slide 406

Kinesis Data Streams
- Retention between up to 365 days
- Ability to reprocess (replay) data by consumers
- Data can’t be deleted from Kinesis (until it expires)
- Data up to 10MiB (typical use case is lot of “small” real-time data)
- Data ordering guarantee for data with the same “Partition ID”
- At-rest KMS encryption, in-flight HTTPS encryption
- Kinesis Producer Library (KPL) to write an optimized producer application
- Kinesis Client Library (KCL) to write an optimized consumer application

## Slide 407

Kinesis Data Streams – Capacity Modes
- Provisioned mode:
- Choose number of shards
- Each shard gets 1MB/s in (or 1000 records per second)
- Each shard gets 2MB/s out
- Scale manually to increase or decrease the number of shards
- Yo u p ay p e r s h a r d p r ov i s i o n e d p e r h o u r
- On-demand mode:
- No need to provision or manage the capacity
- Default capacity provisioned (4 MB/s in or 4000 records per second)
- Scales automatically based on observed throughput peak during the last 30 days
- Pay per stream per hour & data in/out per GB

## Slide 408

Amazon Data Firehose
AWS Destinations
Amazon S3Amazon RedshiftAmazon OpenSearch
HTTP Endpoint
Producers
Applications
Client
Kinesis Agent
SDK
KinesisData Streams
Amazon CloudWatch(Logs & Events)
AWS IoT
AmazonData Firehose
Data transformation
Lambdafunction
Up to 1 MBRecordBatch writes
S3 backup bucketAll or Failed dataCustom Destinations
Datadog
3rd-party Partner Destinations

## Slide 409

Amazon Data Firehose
- Note: used to be called “Kinesis Data Firehose”
- Fully Managed Service
- Amazon Redshift / Amazon S3 / Amazon OpenSearch Service
- 3rd party: Splunk / MongoDB / Datadog / NewRelic / …
- Custom HTTP Endpoint
- Automatic scaling, serverless, pay for what you use
- Near Real-Time with buffering capability based on size / time
- Supports CSV, JSON, Parquet, Avro, Raw Text, Binary data
- Conversions to Parquet / ORC, compressions with gzip / snappy
- Custom data transformations using AWS Lambda (ex: CSV to JSON)

## Slide 410

Kinesis Data Streams vs Amazon Data Firehose
Kinesis Data StreamsAmazon Data Firehose
- Streaming data collection
- Producer & Consumer code
- Real-time
- Provisioned / On-Demand mode
- Data storage up to 365 days
- Replay Capability
- Load streaming data into S3 / Redshift / OpenSearch / 3rd party / custom HTTP
- Fully managed
- Near real-time
- Automatic scaling
- No data storage
- Doesn’t support replay capability

## Slide 411

SQS vs SNS vs KinesisSQS:
- Consumer “pull data”
- Data is deleted after being consumed
- Can have as many workers (consumers) as we want
- No need to provision throughput
- Ordering guarantees only on FIFO queues
- Individual message delay capability
SNS:
- Push data to many subscribers
- Up to 12,500,000 subscribers
- Data is not persisted (lost if not delivered)
- Pub/Sub
- Up to 100,000 topics
- No need to provision throughput
- Integrates with SQS for fan-out architecture pattern
- FIFO capability for SQS FIFO
Kinesis:
- Standard: pull data
- 2 MB per shard
- Enhanced-fan out: push data
- 2 MB per shard per consumer
- Possibility to replay data
- Meant for real-time big data, analytics and ETL
- Ordering at the shard level
- Data expires after X days
- Provisioned mode or on-demand capacity mode

## Slide 412

Amazon MQ
- SQS, SNS are “cloud-native” services: proprietary protocols from AWS
- Traditional applications r unning from on-premises may use open protocols such as: MQTT, AMQP , STOMP , Openwire, WSS
- When migrating to the cloud, instead of re-engineering the application to use SQS and SNS, we can use Amazon MQ
- Amazon MQ is a managed message broker service for
- Amazon MQ doesn’t “scale” as much as SQS / SNS
- Amazon MQ runs on servers, can run in Multi-AZ with failover
- Amazon MQ has both queue feature (~SQS) and topic features (~SNS)

## Slide 413

Amazon MQ – High Availability
Region(us-east-1)Availability Zone(us-east-1a)
Availability Zone(us-east-1b)
ACTIVE
STANDBY
Amazon EFS(storage)
Amazon MQ Broker
Amazon MQ Broker
Clientfailover

## Slide 414

Containers on AWS

## Slide 415

What is Docker?
- Docker is a software development platform to deploy apps
- Apps are packaged in containers that can be run on any OS
- Apps run the same, regardless of where they’re run
- Any machine
- No compatibility issues
- Predictable behavior
- Less work
- Easier to maintain and deploy
- Works with any language, any OS, any technology
- Use cases: microservices architecture, lift-and-shift apps from on-premises to the AWS cloud, …

## Slide 416

Docker on an OSServer (e.g., EC2 instance)

## Slide 417

Where are Docker images stored?
- Docker images are stored in Docker Repositories
- Docker Hub (https://hub.docker.com)
- Public repository
- Find base images for many technologies or OS (e.g., Ubuntu, MySQL, …)
- Amazon ECR (Amazon Elastic Container Registry)
- Private repository
- Public repository (Amazon ECR Public Gallery https://gallery.ecr.aws)

## Slide 418

Docker vs. Virtual Machines
- Docker is ”sort of” a virtualization technology, but not exactly
- Resources are shared with the host => many containers on one server
InfrastructureHost OSHypervisor
AppsGuest OS (VM)
AppsGuest OS (VM)
AppsGuest OS (VM)
InfrastructureHost OS (EC2 Instance)Docker Daemon

## Slide 419

Getting Started with Docker
Dockerfile
imageBuild Runcontainer
Docker Repository
AmazonECR
PushPull

## Slide 420

Docker Containers Management on AWS
- Amazon Elastic Container Service (Amazon ECS)
- Amazon’s own container platform
- Amazon Elastic Kubernetes Service (Amazon EKS)
- Amazon’s managed Kubernetes (open source)
- AWS Fargate
- Amazon’s own Serverless container platform
- Works with ECS and with EKS
- Amazon ECR:
- Store container images
Amazon ECS
AWS Fargate
Amazon EKS
Amazon ECR

## Slide 421

Amazon ECS - EC2 Launch Type
- ECS = Elastic Container Service
- Launch Docker containers on AWS = Launch ECS Tasks on ECS Clusters
- EC2 Launch Type: you must provision & maintain the infrastructure (the EC2 instances)
- Each EC2 Instance must run the ECS Agent to register in the ECS Cluster
- AWS takes care of starting / stopping containers
Amazon ECS / ECS Cluster
New DockerContainer
EC2 Instance
EC2 Instance
EC2 Instance
ECS Agent
ECS Agent
ECS Agent

## Slide 422

Amazon ECS – Fargate Launch Type
- Launch Docker containers on AWS
- Yo u d o n o t p r ov i s i o n t h e i n f r a s t r u c t u r e (no EC2 instances to manage)
- It’s all Serverless!
- Yo u j u s t c r e a t e t a s k d e fi n i t i o n s
- AWS just runs ECS Tasks for you based on the CPU / RAM you need
- To s c a l e , j u s t i n c r e a s e t h e n u m b e r o f tasks. Simple - no more EC2 instances
AWS Fargate / ECS Cluster
New DockerContainer

## Slide 423

Amazon ECS – IAM Roles for ECS
- EC2 Instance Profile (EC2 Launch Type only):
- Used by the ECS agent
- Makes API calls to ECS service
- Send container logs to CloudWatch Logs
- Pull Docker image from ECR
- Reference sensitive data in Secrets Manager or SSM Parameter Store
- ECS Task Role:
- Allows each task to have a specific role
- Use different roles for the different ECS Services you run
- Task Role is defined in the task definition
EC2 Instance
ECS Agent
Task A
Task B
ECSECRCloudWatchLogs
EC2 Instance Profile
S3
DynamoDB
ECS Task A Role
ECS Task B Role

## Slide 424

Amazon ECS – Load Balancer Integrations
- Application Load Balancer supported and works for most use cases
- Network Load Balancer recommended only for high throughput / high performance use cases, or to pair it with AWS Private Link
- Classic Load Balancer supported but not recommended (no advanced features – no Fargate) ECS Cluster
EC2 Instance
EC2 Instance
ApplicationLoad Balancer
Users
ECS Task
ECS Task
ECS Task
ECS Task
80/443

## Slide 425

Amazon ECS – Data Volumes (EFS)
- Mount EFS file systems onto ECS tasks
- Works for both EC2 and Fargate launch types
- Tasks r unning in any AZ will share the same data in the EFS file system
- Fargate + EFS = Serverless
- Use cases: persistent multi-AZ shared storage for your containers
- Note:
- Amazon S3 cannot be mounted as a file system
ECS Cluster
EC2 Instance
Fargate
Amazon EFS
File System
mountmount

## Slide 426

ECS Service Auto Scaling
- Automatically increase/decrease the desired number of ECS tasks
- Amazon ECS Auto Scaling uses AWS Application Auto Scaling
- ECS Service Average CPU Utilization
- ECS Service Average Memory Utilization - Scale on RAM
- ALB Request Count Per Target – metric coming from the ALB
- Tar get Tr acking – scale based on target value for a specific CloudWatch metric
- Step Scaling – scale based on a specified CloudWatch Alarm
- Scheduled Scaling – scale based on a specified date/time (predictable changes)
- ECS Service Auto Scaling (task level) ≠ EC2 Auto Scaling (EC2 instance level)
- Fargate Auto Scaling is much easier to setup (because Serverless)

## Slide 427

EC2 Launch Type – Auto Scaling EC2 Instances
- Accommodate ECS Service Scaling by adding underlying EC2 Instances
- Auto Scaling Group Scaling
- Scale your ASG based on CPU Utilization
- Add EC2 instances over time
- ECS Cluster Capacity Provider
- Used to automatically provision and scale the infrastructure for your ECS Tasks
- Capacity Provider paired with an Auto Scaling Group
- Add EC2 Instances when you’re missing capacity (CPU, RAM…)

## Slide 428

ECS Scaling – Service CPU Usage Example
Service A
Auto Scaling Group
Scale ECS Capacity Providers(optional)
Auto Scaling
Task 1
Task 2
CPU Usage
CloudWatch Metric(ECS Service CPU Usage)
CloudWatch AlarmTriggerScale
Task 3(new)

## Slide 429

ECS tasks invoked by Event Bridge
Client
Upload object
AmazonEventBridge
VPC
Region
Amazon ECS Cluster
AWS Fargate
Task (new)S3 BucketEventRule: Run ECS Task
AmazonDynamoDB
Save result
ECS Task Role(Access S3 & DynamoDB)
Get object

## Slide 430

ECS tasks invoked by Event Bridge Schedule
AmazonEventBridgeAmazon ECS Cluster
AWS Fargate
Task (new)Rule: Run ECS Task
Amazon S3Batch Processing
Every 1 hour
ECS Task RoleAccess S3

## Slide 431

ECS – SQS Queue Example
Service A
ECS Service Auto Scaling
Task 1
Task 2
Task 3
SQS QueueMessagesPoll for messages

## Slide 432

ECS – Intercept Stopped Tasks using EventBridge
ECS Task
Containers
exited
EventBridge
event
Event Pattern
SNStrigger
Administratoremail

## Slide 433

Amazon ECR
- ECR = Elastic Container Registry
- Store and manage Docker images on AWS
- Private and Public repository (Amazon ECR Public Gallery https://gallery.ecr.aws)
- Fully integrated with ECS, backed by Amazon S3
- Access is controlled through IAM (permission errors => policy)
- Supports image vulnerability scanning, versioning, image tags, image lifecycle, …
ECS Cluster
EC2 Instance
ECR Repository
DockerImage ADockerImage B
IAM Rolepullpull

## Slide 434

Amazon EKS Overview
- Amazon EKS = Amazon Elastic Kubernetes Service
- It is a way to launch managed Kubernetes clusters on AWS
- Kubernetes is an open-source system for automatic deployment, scaling and management of containerized (usually Docker) application
- It’s an alternative to ECS, similar goal but different API
- EKS supports EC2 if you want to deploy worker nodes or Fargate to deploy serverless containers
- Use case: if your company is already using Kubernetes on-premises or in another cloud, and wants to migrate to AWS using Kubernetes
- Kubernetes is cloud-agnostic (can be used in any cloud – Azure, GCP…)
- For multiple regions, deploy one EKS cluster per region
- Collect logs and metrics using CloudWatch Container Insights

## Slide 435

Private subnet 3Private subnet 1Private subnet 2
Amazon EKS - Diagram
AWS CloudVPC
Availability Zone 1Availability Zone 2Availability Zone 3Public subnet 1
Public subnet 2
Public subnet 3
NGW
NGW
NGW
ELB
ELB
EKSPublicService LB
ELB
Auto Scaling GroupEKS node
EKS Pods
EKS node
EKS Pods
EKS node
EKS PodsEKS Worker Nodes
EKSPrivateService LB
ELB

## Slide 436

Amazon EKS – Node Types
- Managed Node Groups
- Creates and manages Nodes (EC2 instances) for you
- Nodes are part of an ASG managed by EKS
- Supports On-Demand or Spot Instances
- Self-Managed Nodes
- Nodes created by you and registered to the EKS cluster and managed by an ASG
- Yo u c a n u s e p r e bu i l t A M I - Amazon EKS Optimized AMI
- Supports On-Demand or Spot Instances
- AWS Fargate
- No maintenance required; no nodes managed

## Slide 437

Amazon EKS – Data Volumes
- Need to specify StorageClass manifest on your EKS cluster
- Leverages a Container Storage Interface (CSI) compliant driver
- Support for…
- Amazon EBS
- Amazon EFS (works with Fargate)
- Amazon FSx for Lustre
- Amazon FSx for NetApp ONTAP

## Slide 438

AWS App Runner
- Fully managed service that makes it easy to deploy web applications and APIs at scale
- No infrastructure experience required
- Start with your source code or container image
- Automatically builds and deploy the web app
- Automatic scaling, highly available, load balancer, encr yption
- VPC access support
- Connect to database, cache, and message queue services
- Use cases: web apps, APIs, microservices, rapid production deployments
SourceCodeContainerImage (Docker)
Configure SettingsvCPU, RAM,Auto Scaling, Health Check
Create & Deploy
Access using URL

## Slide 439

AWS App2Container (A2C)
- CLI tool for migrating and modernizing Java and .NET web apps into Docker Containers
- Lift-and-shift your apps running in on-premises bare metal, virtual machines, or in any Cloud to AWS
- Accelerate modernization, no code changes, migrate legacy apps…
- Generates CloudFormation templates (compute, network…)
- Register generated Docker containers to ECR
- Deploy to ECS, EKS, or App Runner
- Supports pre-built CI/CD pipelines

## Slide 440

AWS App2Container (A2C)
Discover & Analyzecreate app inventoryand analyze runtimedependencies
Extract & Containerizeextract an app withdependencies and createa Docker image
Create DeploymentArtifactsgenerate ECS Task andEKS Pod definitions, andcreate CI/CD pipelines, andother infrastructure
Deploy to AWSstore Docker image in ECR,and deploy to ECS, EKS, orApp Runner
CloudFormationTemplate
Amazon ECR(store image)
Amazon ECS(deploy)
Amazon EKS(deploy)
App Runner(deploy)

## Slide 441

Serverless Overview

## Slide 442

What’s serverless?
- Serverless is a new paradigm in which the developers don’t have to manage servers anymore…
- They just deploy code
- They just deploy… functions !
- Initially... Serverless == FaaS (Function as a Service)
- Serverless was pioneered by AWS Lambda but now also includes anything that’s managed: “databases, messaging, storage, etc.”
- Serverless does not mean there are no servers…it means you just don’t manage / provision / see them

## Slide 443

Serverless in AWS
- AWS Lambda
- DynamoDB
- AWS Cognito
- AWS API Gateway
- Amazon S3
- AWS SNS & SQS
- AWS Kinesis Data Firehose
- Aurora Serverless
- Step Functions
- Fargate
Static content
REST API
Log in
S3 bucketAPI Gateway
Cognito
Users
Lambda
DynamoDB

## Slide 444

Why AWS Lambda
- Virtual Servers in the Cloud
- Limited by RAM and CPU
- Continuously running
- Scaling means intervention to add / remove servers
Amazon Lambda
Amazon EC2
- Virtual functions – no servers to manage!
- Limited by time - short executions
- Run on-demand
- Scaling is automated!

## Slide 445

Benefits of AWS Lambda
- Easy Pricing:
- Pay per request and compute time
- Free tier of 1,000,000 AWS Lambda requests and 400,000 GBs of compute time
- Integrated with the whole AWS suite of services
- Integrated with many programming languages
- Easy monitoring through AWS CloudWatch
- Easy to get more resources per functions (up to 10GB of RAM!)
- Increasing RAM will also improve CPU and network!

## Slide 446

AWS Lambda language support
- Node.js (JavaScript)
- Python
- Java
- C# (.NET Core) / Powershell
- Ruby
- Custom Runtime API (community supported, example Rust or Golang)
- Lambda Container Image
- The container image must implement the Lambda Runtime API
- ECS / Fargate is preferred for running arbitrary Docker images

## Slide 447

AWS Lambda IntegrationsMain ones
CloudWatch LogsSNS CognitoSQS
S3KinesisAPI GatewayDynamoDB
CloudFront
CloudWatch EventsEventBridge

## Slide 448

Example: Serverless Thumbnail creation
New image in S3AWS Lambda FunctionCreates a Thumbnail
triggerpushNew thumbnail in S3
Metadata in DynamoDB
push Image nameImage sizeCreation dateetc…

## Slide 449

Example: Serverless CRON Job
CloudWatch EventsEventBridgeAWS Lambda FunctionPerform a task
TriggerEvery 1 hour

## Slide 450

AWS Lambda Pricing: example
- Yo u c a n fi n d ove r a l l p r i c i n g i n fo r m a t i o n h e r e : https://aws.amazon.com/lambda/pricing/
- Pay per calls:
- First 1,000,000 requests are free
- $0.20 per 1 million requests thereafter ($0.0000002 per request)
- Pay per duration: (in increment of 1 ms)
- 400,000 GB-seconds of compute time per month for FREE
- == 400,000 seconds if function is 1GB RAM
- == 3,200,000 seconds if function is 128 MB RAM
- After that $1.00 for 600,000 GB-seconds
- It is usually very cheap to run AWS Lambda so it’s very popular

## Slide 451

AWS Lambda Limits to Know - per region
- Execution:
- Memory allocation: 128 MB – 10GB (1 MB increments)
- Maximum execution time: 900 seconds (15 minutes)
- Environment variables (4 KB)
- Disk capacity in the “function container” (in /tmp): 512 MB to 10GB
- Concurrency executions: 1000 (can be increased)
- Deployment:
- Lambda function deployment size (compressed .zip): 50 MB
- Size of uncompressed deployment (code + dependencies): 250 MB
- Can use the /tmp directory to load other files at startup
- Size of environment variables: 4 KB

## Slide 452

Lambda Concurrency and Throttling
- Concurrency limit: up to 1000 concurrent executions
- Can set a “reserved concurrency” at the function level (=limit)
- Each invocation over the concurrency limit will trigger a “Throttle”
- Throttle behavior:
- If synchronous invocation => return ThrottleError - 429
- If asynchronous invocation => retry automatically and then go to DLQ
- If you need a higher limit, open a support ticket

## Slide 453

Lambda Concurrency Issue
- If you don’t reserve (=limit) concurrency, the following can happen:
Application Load Balancer
API Gateway
SDK / CLI
1000 concurrent executions
THROTTLE!
THROTTLE!
Many users
Few users

## Slide 454

Concurrency and Asynchronous Invocations
S3 bucket
New file event
New file event
New file event
- If the function doesn't have enough concurrency available to process all events, additional requests are throttled.
- For throttling errors (429) and system errors (500-series), Lambda returns the event to the queue and attempts to run the function again for up to 6 hours.
- The retry interval increases exponentially from 1 second after the first attempt to a maximum of 5 minutes.

## Slide 455

Cold Starts & Provisioned Concurrency
- Cold Start:
- New instance => code is loaded and code outside the handler run (init)
- If the init is large (code, dependencies, SDK…) this process can take some time.
- First request served by new instances has higher latency than the rest
- Provisioned Concurrency:
- Concurrency is allocated before the function is invoked (in advance)
- So the cold start never happens and all invocations have low latency
- Application Auto Scaling can manage concurrency (schedule or target utilization)
- Note:
- Note: cold starts in VPC have been dramatically reduced in Oct & Nov 2019
- https://aws.amazon.com/blogs/compute/announcing-improved-vpc-networking-for-aws-lambda-functions/

## Slide 456

Reserved and Provisioned Concurrency
https://docs.aws.amazon.com/lambda/latest/dg/configuration-concurrency.html

## Slide 457

Lambda SnapStart
- Improves your Lambda functions performance up to 10x at no extra cost for Java, Python & .NET
- When enabled, function is invoked from a pre-initialized state (no function initialization from scratch)
- When you publish a new version:
- Lambda initializes your function
- Takes a snapshot of memor y and disk state of the initialized function
- Snapshot is cached for low-latency access
Lambda InvocationLifecycle Phases
SnapStartdisabled
invoke
InitInvokeShutdown
SnapStartenabled
invoke
InvokeShutdown
function ispre-initialized
LambdaLambda

## Slide 458

Customization At The Edge
- Many modern applications execute some form of the logic at the edge
- Edge Function:
- A code that you write and attach to CloudFront distributions
- Runs close to your users to minimize latency
- CloudFront provides two types: CloudFront Functions & Lambda@Edge
- Yo u d o n ’ t h ave t o m a n a g e a ny s e r ve r s , d e p l oye d g l o b a l ly
- Use case: customize the CDN content
- Pay only for what you use
- Fully serverless

## Slide 459

CloudFront Functions & Lambda@Edge Use Cases
- Website Security and Privacy
- Dynamic Web Application at the Edge
- Search Engine Optimization (SEO)
- Intelligently Route Across Origins and Data Centers
- Bot Mitigation at the Edge
- Real-time Image Transformation
- A/B Testing
- User Authentication and Authorization
- User Prioritization
- User Tracking and Analytics

## Slide 460

CloudFront Functions
- Lightweight functions written in JavaScript
- For high-scale, latency-sensitive CDN customizations
- Sub-ms startup times, millions of requests/second
- Used to change Viewer requests and responses:
- Viewer Request: after CloudFront receives a request from a viewer
- Viewer Response: before CloudFront forwards the response to the viewer
- Native feature of CloudFront (manage code entirely within CloudFront)
Origin
CloudFront
Client
ViewerRequest ViewerResponse
OriginRequest OriginResponse

## Slide 461

Lambda@Edge
- Lambda functions written in NodeJS or Python
- Scales to 1000s of requests/second
- Used to change CloudFront requests and responses:
- Viewer Request – after CloudFront receives a request from a viewer
- Origin Request – before CloudFront forwards the request to the origin
- Origin Response – after CloudFront receives the response from the origin
- Viewer Response – before CloudFront forwards the response to the viewer
- Author your functions in one AWS Region (us-east-1), then CloudFront replicates to its locations
Origin
CloudFront
Client
ViewerRequest ViewerResponse
OriginRequest OriginResponse

## Slide 462

CloudFront Functions vs. Lambda@EdgeCloudFront FunctionsLambda@EdgeRuntime SupportJavaScriptNode.js, Python# of RequestsMillions of requests per secondThousands of requests per secondCloudFront Triggers-Viewer Request/Response-Viewer Request/Response-Origin Request/ResponseMax. Execution Time< 1 ms 5 – 10 secondsMax. Memory2 MB 128 MB up to 10 GBTotal Package Size10 KB 1 MB – 50 MBNetwork Access, File System AccessNo YesAccess to the Request BodyNo YesPricingFree tier available, 1/6th price of @EdgeNo free tier, charged per request & duration

## Slide 463

CloudFront Functions vs. Lambda@Edge - Use CasesCloudFront Functions
- Cache key normalization
- Transform request attributes (header s, cookies, query strings, URL) to create an optimal Cache Key
- Header manipulation
- Insert/modify/delete HTTP headers in the request or response
- URL rewrites or redirects
- Request authentication & authorization
- Create and validate user-generated tokens (e.g., JWT) to allow/deny requests
Lambda@Edge
- Longer execution time (several ms)
- Adjustable CPU or memory
- Y our code depends on a 3rd libraries (e.g., AWS SDK to access other AWS services)
- Network access to use external services for processing
- File system access or access to the body of HTTP requests

## Slide 464

Lambda by default
VPC & Private Subnet
AWS Cloud
Private RDS
DynamoDB
Not working
Default Lambda Deployment
Public www works
- By default, your Lambda function is launched outside your own VPC (in an AWS-owned VPC)
- Therefore, it cannot access resources in your VPC (RDS, ElastiCache, internal ELB…)

## Slide 465

Private subnet
Lambda in VPC
- Yo u mu s t d e fi n e t h e V P C I D, t h e Subnets and the Security Groups
- Lambda will create an ENI (Elastic Network Interface) in your subnets
Lambda Function
Amazon RDS In VPC
Lambda Security group
RDS Security group
Elastic Network Interface (ENI)

## Slide 466

Lambda with RDS Proxy
- If Lambda functions directly access your database, they may open too many connections under high load
- RDS Proxy
- Improve scalability by pooling and sharing DB connections
- Improve availability by reducing by 66% the failover time and preserving connections
- Improve security by enforcing IAM authentication and storing credentials in Secrets Manager
- The Lambda function must be deployed in your VPC, because RDS Proxy is never publicly accessible
VPC
Private subnet
…Lambda functions
RDS Proxy
RDS DBInstance

## Slide 467

Invoking Lambda from RDS & Aurora
- Invoke Lambda functions from within your DB instance
- Allows you to process data events from within a database
- Supported for RDS for PostgreSQL and Aurora MySQL
- Must allow outbound traffic to your Lambda function from within your DB instance (Public, NAT GW, VPC Endpoints)
- DB instance must have the required permissions to invoke the Lambda function (Lambda Resource-based Policy & IAM Policy)
RDS DBInstance
Amazon SES
Lambdafunction
Userregister(INSERT)
invoke
send Email
Permissions

## Slide 468

RDS Event Notifications
- Notifications that tells information about the DB instance itself (created, stopped, start, …)
- Yo u d o n ’ t h ave a ny i n fo r m a t i o n a b o u t t h e d a t a i t s e l f
- Subscribe to the following event categories: DB instance, DB snapshot, DB Parameter Group, DB Security Group, RDS Proxy, Custom Engine Version
- Near real-time events (up to 5 minutes)
- Send notifications to SNS or subscribe to events using EventBridge
RDS DBInstance
EventBridge
SNS
Lambdafunction
Lambdafunction
SQSQueue…

## Slide 469

Amazon DynamoDB
- Fully managed, highly available with replication across multiple AZs
- NoSQL database - not a relational database - with transaction support
- Scales to massive workloads, distributed database
- Millions of requests per seconds, trillions of row, 100s of TB of storage
- Fast and consistent in performance (single-digit millisecond)
- Integrated with IAM for security, authorization and administration
- Low cost and auto-scaling capabilities
- No maintenance or patching, always available
- Standard & Infrequent Access (IA) Table Class

## Slide 470

DynamoDB - Basics
- DynamoDB is made of Tables
- Each table has a Primary Key (must be decided at creation time)
- Each table can have an infinite number of items (= rows)
- Each item has attributes (can be added over time – can be null)
- Maximum size of an item is 400KB
- Data types supported are:
- Scalar Types – String, Number, Binary, Boolean, Null
- Document Types – List, Map
- Set Types – String Set, Number Set, Binary Set
- Therefore, in DynamoDB you can rapidly evolve schemas

## Slide 471

DynamoDB – Table example
User_IDGame_IDScoreResultPartition Key
7791a3d6-…873e0634-…873e0634-…
Primary Key
442192Win
452177Win1894 Lose
AttributesSort Key
14

## Slide 472

DynamoDB – Read/Write Capacity Modes
- Control how you manage your table’s capacity (read/write throughput)
- Provisioned Mode (default)
- Yo u s p e c i f y t h e nu m b e r o f r e a d s / w r i t e s p e r s e c o n d
- Yo u n e e d t o p l a n c a p a c i t y b e fo r e h a n d
- Pay for provisioned Read Capacity Units (RCU) & Write Capacity Units (WCU)
- Possibility to add auto-scaling mode for RCU & WCU
- On-Demand Mode
- Read/writes automatically scale up/down with your workloads
- No capacity planning needed
- Pay for what you use, more expensive ($$$)
- Great for unpredictable workloads, steep sudden spikes

## Slide 473

DynamoDB Accelerator (DAX)
- Fully-managed, highly available, seamless in-memory cache for DynamoDB
- Help solve read congestion by caching
- Microseconds latency for cached data
- Doesn’t require application logic modification (compatible with existing DynamoDB APIs)
- 5 minutes TTL for cache (default)
Application
Amazon DynamoDB
Tables
…
DAX Cluster
…Nodes

## Slide 474

DynamoDB Accelerator (DAX) vs. ElastiCache
DynamoDB Accelerator (DAX)AmazonDynamoDB
AmazonElastiCache
- Individual objects cache- Query & Scan cache
Store Aggregation Result
Application

## Slide 475

DynamoDB – Stream Processing
- Ordered stream of item-level modifications (create/update/delete) in a table
- Use cases:
- React to changes in real-time (welcome email to users)
- Real-time usage analytics
- Insert into derivative tables
- Implement cross-region replication
- Invoke AWS Lambda on changes to your DynamoDB tableDynamoDB Streams
- 24 hours retention
- Limited # of consumers
- Process using AWS Lambda Triggers, or DynamoDB Stream Kinesis adapter
Kinesis Data Streams (newer)
- 1 year retention
- High # of consumers
- Process using AWS Lambda, Kinesis Data Analytics, Kineis Data Firehose, AWS Glue Streaming ETL…

## Slide 476

DynamoDB Streams
Application
TableDynamoDBStreams
Kinesis DataStreams
Processing Layer
DynamoDBKCL AdapterLambda
Amazon SNSmessaging, notifications
DDB Tablefiltering, transforming, …
Kinesis DataFirehose
AmazonRedshiftanalytics
Amazon S3archiving
AmazonOpenSearchindexing
create/update/delete

## Slide 477

DynamoDB Global Tables
- Make a DynamoDB table accessible with low latency in multiple-regions
- Active-Active replication
- Applications can READ and WRITE to the table in any region
- Must enable DynamoDB Streams as a pre-requisite
TableUS-EAST-1 TableAP-SOUTHEAST-2two-wayreplication
GLOBAL TABLE

## Slide 478

DynamoDB – Time T o Live (TTL)
- Automatically delete items after an expiry timestamp
- Use cases: reduce stored data by keeping only current items, adhere to regulatory obligations, web session handling…User_IDSession_IDExpTime (TTL)7791a3d6-…746865726521631188571873e0634-…6e6f74686961631274971a80f73a1-…746f20736561631102171
SessionData (Table)
Expiration Process
Friday, September 10, 2021, 11:56:11 AM(Epoch timestamp: 1631274971)
Current Time
scan &expire items
Deletion Process
scan &delete items

## Slide 479

DynamoDB – Backups for disaster recovery
- Continuous backups using point-in-time recovery (PITR)
- Optionally enabled for the last 35 days
- Point-in-time recovery to any time within the backup window
- The recovery process creates a new table
- On-demand backups
- Full backups for long-term retention, until explicitely deleted
- Doesn’t affect performance or latency
- Can be configured and managed in AWS Backup (enables cross-region copy)
- The recovery process creates a new table

## Slide 480

DynamoDB – Integration with Amazon S3
- Export to S3 (must enable PITR)
- Works for any point of time in the last 35 days
- Doesn’t affect the read capacity of your table
- Perform data analysis on top of DynamoDB
- Retain snapshots for auditing
- ETL on top of S3 data before importing back into DynamoDB
- Export in DynamoDB JSON or ION format
- Import from S3
- Import CSV, DynamoDB JSON or ION format
- Doesn’t consume any write capacity
- Creates a new table
- Import errors are logged in CloudWatch Logs
DynamoDBS3Athena
S3(.csv, .json, .ion)DynamoDBimport
exportquery

## Slide 481

Example: Building a Serverless API
REST APIPROXY REQUESTSCRUD
API GatewayClient Lambda DynamoDB

## Slide 482

AWS API Gateway
- AWS Lambda + API Gateway: No infrastructure to manage
- Support for the WebSocket Protocol
- Handle API versioning (v1, v2…)
- Handle different environments (dev, test, prod…)
- Handle security (Authentication and Authorization)
- Create API keys, handle request throttling
- Swagger / Open API import to quickly define APIs
- Transform and validate requests and responses
- Generate SDK and API specifications
- Cache API responses

## Slide 483

API Gateway – Integrations High Level
- Lambda Function
- Invoke Lambda function
- Easy way to expose REST API backed by AWS Lambda
- HTTP
- Expose HTTP endpoints in the backend
- Example: internal HTTP API on premise, Application Load Balancer…
- Why? Add rate limiting, caching, user authentications, API keys, etc…
- AWS Service
- Expose any AWS API through the API Gateway
- Example: start an AWS Step Function workflow, post a message to SQS
- Why? Add authentication, deploy publicly, rate control…

## Slide 484

API Gateway – AWS Service IntegrationKinesis Data Streams example
API GatewayKinesis DataStreamsKinesis DataFirehoseAmazon S3sendrecordsstore .jsonfiles
Client
requests

## Slide 485

API Gateway - Endpoint Types
- Edge-Optimized (default): For global clients
- Requests are routed through the CloudFront Edge locations (improves latency)
- The API Gateway still lives in only one region
- Regional:
- For clients within the same region
- Could manually combine with CloudFront (more control over the caching strategies and the distribution)
- Private:
- Can only be accessed from your VPC using an interface VPC endpoint (ENI)
- Use a resource policy to define access

## Slide 486

API Gateway – Security
- User Authentication through
- IAM Roles (useful for internal applications)
- Cognito (identity for external users – example mobile users)
- Custom Authorizer (your own logic)
- Custom Domain Name HTTPS security through integration with AWS Certificate Manager (ACM)
- If using Edge-Optimized endpoint, then the certificate must be in us-east-1
- If using Regional endpoint, the certificate must be in the API Gateway region
- Must setup CNAME or A-alias record in Route 53

## Slide 487

AWS Step Functions
- Build serverless visual workflow to orchestrate your Lambda functions
- Features: sequence, parallel, conditions, timeouts, error handling, …
- Can integrate with EC2, ECS, On-premises servers, API Gateway, SQS queues, etc…
- Possibility of implementing human approval feature
- Use cases: order fulfillment, data processing, web applications, any workflow

## Slide 488

Amazon Cognito
- Give users an identity to interact with our web or mobile application
- Cognito User Pools:
- Sign in functionality for app users
- Integrate with API Gateway & Application Load Balancer
- Cognito Identity Pools (Federated Identity):
- Provide AWS credentials to users so they can access AWS resources directly
- Integrate with Cognito User Pools as an identity provider
- Cognito vs IAM: “hundreds of users”, ”mobile users”, “authenticate with SAML”

## Slide 489

Cognito User Pools (CUP) – User Features
- Create a serverless database of user for your web & mobile apps
- Simple login: Username (or email) / password combination
- Password reset
- Email & Phone Number Verification
- Multi-factor authentication (MFA)
- Federated Identities: users from Facebook, Google, SAML…

## Slide 490

Cognito User Pools (CUP) - Integrations
- CUP integrates with API Gateway and Application Load Balancer
REST API + Pass TokenEvaluate Cognito Token
AuthenticateRetrieve token
Cognito User Pools
Cognito User Pools
Application Load Balancer+ Listeners & Rules
Authenticate
Target GroupBackend
API Gatewaybackend

## Slide 491

Cognito Identity Pools (Federated Identities)
- Get identities for “users” so they obtain temporary AWS credentials
- Users source can be Cognito User Pools, 3rd party logins, etc…
- Users can then access AWS services directly or through API Gateway
- The IAM policies applied to the credentials are defined in Cognito
- They can be customized based on the user_id for fine grained control
- Default IAM roles for authenticated and guest users

## Slide 492

Cognito Identity Pools – Diagram
Cognito Identity Pools
Web & Mobile Applications
Social Identity Provider
Login and Get Token
Exchange token for temporaryAWS credentialsvalidate
Private S3 Bucket DynamoDB Table
Direct access to AWS
CognitoUser Pools

## Slide 493

Cognito Identity PoolsRow Level Security in DynamoDB

## Slide 494

Serverless Architectures

## Slide 495

Mobile application: MyT odoList
- We want to create a mobile application with the following requirements
- Expose as REST API with HTTPS
- Serverless architecture
- Users should be able to directly interact with their own folder in S3
- Users should authenticate through a managed serverless service
- The users can write and read to-dos, but they mostly read them
- The database should scale, and have some high read throughput

## Slide 496

Mobile app: REST API layer
Amazon API Gateway
Amazon Cognito
AWS LambdaAmazon DynamoDB
Mobile client
REST HTTPSinvokequery
authenticateVerify authentication

## Slide 497

Mobile app: giving users access to S3
Amazon API Gateway
Amazon Cognito
AWS LambdaAmazon DynamoDB
Amazon S3
Mobile client
Store/retrieve filesPermissions
authenticate

## Slide 498

Mobile app: high read throughput, static data
Amazon API Gateway
Amazon Cognito
AWS LambdaDynamoDB
Amazon S3
Mobile client
Store/retrieve filesPermissions
REST HTTPSinvokeQuery / read
authenticateVerify authentication
DAXCaching layer

## Slide 499

Mobile app: caching at the API Gateway
Amazon API Gateway
Amazon Cognito
AWS LambdaDynamoDB
Amazon S3
Mobile client
Store/retrieve filesPermissions
REST HTTPSinvokeQuery / read
authenticateVerify authentication
DAXCaching layer
CACHING OF RESPONSES

## Slide 500

In this lecture
- Serverless REST API: HTTPS, API Gateway, Lambda, DynamoDB
- Using Cognito to generate temporary credentials to access S3 bucket with restricted policy. App users can directly access AWS resources this way. Pattern can be applied to DynamoDB, Lambda…
- Caching the reads on DynamoDB using DAX
- Caching the REST requests at the API Gateway level
- Security for authentication and authorization with Cognito

## Slide 501

Serverless hosted website: MyBlog.com
- This website should scale globally
- Blogs are rarely written, but often read
- Some of the website is purely static files, the rest is a dynamic REST API
- Caching must be implement where possible
- Any new users that subscribes should receive a welcome email
- Any photo uploaded to the blog should have a thumbnail generated

## Slide 502

Serving static content, globallyAmazon CloudFrontGlobal distribution
Amazon S3
Client
Interaction with edge locations

## Slide 503

Serving static content, globally, securelyAmazon CloudFrontGlobal distribution
Amazon S3
Client
Interaction with edge locationsOAC: Origin Access Control
Bucket policyOnly authorize fromCloudFront Distribution

## Slide 504

Adding a public serverless REST APIAmazon CloudFrontGlobal distribution
Amazon S3
Client
Interaction with edge locationsOAC: Origin Access Control
Bucket policyOnly authorize fromCloudFront Distribution
Amazon API Gateway
AWS LambdaDynamoDB
REST HTTPSinvokeQuery / read
DAXCaching layer

## Slide 505

Leveraging DynamoDB Global TablesAmazon CloudFrontGlobal distribution
Amazon S3
Client
Interaction with edge locationsOAC: Origin Access Control
Bucket policyOnly authorize fromCloudFront Distribution
Amazon API Gateway
AWS LambdaDynamoDBGlobal Tables
REST HTTPSinvokeQuery / read
DAXCaching layer

## Slide 506

User Welcome email flowAmazon CloudFrontGlobal distribution
Amazon S3
Client
Interaction with edge locationsOAC: Origin Access Control
Bucket policyOnly authorize fromCloudFront Distribution
Amazon API Gateway
AWS LambdaDynamoDB
REST HTTPSinvokeQuery / read
DAXCaching layer
DynamoDBStream
Stream changes
AWS LambdaInvoke lambda
Amazon Simple Email Service (SES)
SDK to send emailIAM Role

## Slide 507

Thumbnail Generation flowAmazon CloudFrontGlobal distribution
Amazon S3
Client
Interaction with edge locationsOAC: Origin Access Control
Bucket policyOnly authorize fromCloudFront Distribution
Amazon API Gateway
AWS LambdaDynamoDB
REST HTTPSinvokeQuery / read
DAXCaching layer
Amazon CloudFrontGlobal distribution
Amazon S3
OACUpload photosTransfer accelerationtrigger
Amazon S3
thumbnailAWS Lambda
SQS
SNS
optional

## Slide 508

AWS Hosted Website Summary
- We’ve seen static content being distributed using CloudFront with S3
- The REST API was serverless, didn’t need Cognito because public
- We leveraged a Global DynamoDB table to ser ve the data globally
- (we could have used Aurora Global Database)
- We enabled DynamoDB streams to trigger a Lambda function
- The lambda function had an IAM role which could use SES
- SES (Simple Email Service) was used to send emails in a serverless way
- S3 can trigger SQS / SNS / Lambda to notify of events

## Slide 509

Micro Services architecture
- We want to switch to a micro ser vice architecture
- Many services interact with each other directly using a REST API
- Each architecture for each micro service may vary in form and shape
- We want a micro-service architecture so we can have a leaner development lifecycle for each service

## Slide 510

Micro Services Environment
Amazon API Gateway
AWS LambdaElastiCache
Elastic Load Balancing
Amazon EC2 Auto ScalingAmazon RDS
Elastic Load Balancing
ECS
DynamoDB
Amazon Route 53
service1.example.com
service2.example.com
service3.example.com
Users
DNS Query
HTTPS

## Slide 511

Discussions on Micro Services
- Yo u a r e f r e e t o d e s i g n e a c h m i c r o-service the way you want
- Synchronous patterns: API Gateway, Load Balancers
- Asynchronous patterns: SQS, Kinesis, SNS, Lambda triggers (S3)
- Challenges with micro-services:
- repeated overhead for creating each new microservice,
- issues with optimizing server density/utilization
- complexity of running multiple versions of multiple microservices simultaneously
- proliferation of client-side code requirements to integrate with many separate services.
- Some of the challenges are solved by Serverless patterns:
- API Gateway, Lambda scale automatically and you pay per usage
- Yo u c a n e a s i ly c l o n e A P I , r e p r o d u c e e nv i r o n m e n t s
- Generated client SDK through Swagger integration for the API Gateway

## Slide 512

Software updates offloading
- We have an application running on EC2, that distributes software updates once in a while
- When a new software update is out, we get a lot of request and the content is distributed in mass over the network. It’s very costly
- We don’t want to change our application, but want to optimize our cost and CPU, how can we do it?

## Slide 513

Our application current state
Availability zone 1 to 3
Availability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3
Amazon Elastic File System

## Slide 514

Easy way to fix things!
Availability zone 1 to 3
Availability zone 1
Auto Scaling group
Availability zone 2
Availability zone 3
Amazon Elastic File System
Amazon CloudFront

## Slide 515

Why CloudFront?
- No changes to architecture
- Will cache software update files at the edge
- Software update files are not dynamic, they’re static (never changing)
- Our EC2 instances aren’t serverless
- But CloudFront is, and will scale for us
- Our ASG will not scale as much, and we’ll save tremendously in EC2
- We’ll also save in availability, network bandwidth cost, etc
- Easy way to make an existing application more scalable and cheaper!

## Slide 516

Databases in AWS

## Slide 517

Choosing the Right Database
- We have a lot of managed databases on AWS to choose from
- Questions to choose the right database based on your architecture:
- Read-heavy, write-heavy, or balanced workload? Throughput needs? Will it change, does it need to scale or fluctuate during the day?
- How much data to store and for how long? Will it grow? Average object size? How are they accessed?
- Data durability? Source of truth for the data ?
- Latency requirements? Concurrent users?
- Data model? How will you query the data? Joins? Structured? Semi-Structured?
- Strong schema? More flexibility? Reporting? Search? RDBMS / NoSQL?
- License costs? Switch to Cloud Native DB such as Aurora?

## Slide 518

Database Types
- RDBMS (= SQL / OLTP): RDS, Aurora – great for joins
- NoSQL database – no joins, no SQL : DynamoDB (~JSON), ElastiCache (key / value pairs), Neptune (graphs), DocumentDB (for MongoDB), Keyspaces (for Apache Cassandra)
- Object Store: S3 (for big objects) / Glacier (for backups / archives)
- Data Warehouse (= SQL Analytics / BI): Redshift (OLAP), Athena, EMR
- Search: OpenSearch (JSON) – free text, unstructured searches
- Graphs: Amazon Neptune – displays relationships between data
- Ledger: Amazon Quantum Ledger Database
- Time series: Amazon Timestream
- Note: some databases are being discussed in the Data & Analytics section

## Slide 519

Amazon RDS – Summary
- Managed PostgreSQL / MySQL / Oracle / SQL Server / DB2 / MariaDB / Custom
- Provisioned RDS Instance Size and EBS Volume Type & Size
- Auto-scaling capability for Storage
- Support for Read Replicas and Multi AZ
- Security through IAM, Security Groups, KMS , SSL in transit
- Automated Backup with Point in time restore feature (up to 35 days)
- Manual DB Snapshot for longer-term recovery
- Managed and Scheduled maintenance (with downtime)
- Support for IAM Authentication, integration with Secrets Manager
- RDS Custom for access to and customize the underlying instance (Oracle & SQL Server)
- Use case: Store relational datasets (RDBMS / OLTP), perform SQL queries, transactions

## Slide 520

- Compatible API for PostgreSQL / MySQL, separation of storage and compute
- Storage: data is stored in 6 replicas, across 3 AZ – highly available, self-healing, auto-scaling
- Compute: Cluster of DB Instance across multiple AZ, auto-scaling of Read Replicas
- Cluster: Custom endpoints for writer and reader DB instances
- Same security / monitoring / maintenance features as RDS
- Know the backup & restore options for Aurora
- Aurora Serverless – for unpredictable / intermittent workloads, no capacity planning
- Aurora Global: up to 16 DB Read Instances in each region, < 1 second storage replication
- Aurora Machine Learning: perform ML using SageMaker & Comprehend on Aurora
- Aurora Database Cloning: new cluster from existing one, faster than restoring a snapshot
- Use case: same as RDS, but with less maintenance / more flexibility / more performance / more features
Amazon Aurora – Summary

## Slide 521

Amazon ElastiCache – Summary
- Managed Redis / Memcached (similar offering as RDS, but for caches)
- In-memory data store, sub-millisecond latency
- Select an ElastiCache instance type (e.g., cache.m6g.large)
- Support for Clustering (Redis) and Multi AZ, Read Replicas (sharding)
- Security through IAM, Security Groups, KMS, Redis Auth
- Backup / Snapshot / Point in time restore feature
- Managed and Scheduled maintenance
- Requires some application code changes to be leveraged
- Use Case: Key/Value store, Frequent reads, less writes, cache results for DB queries, store session data for websites, cannot use SQL.

## Slide 522

Amazon DynamoDB – Summary
- AWS proprietary technology, managed serverless NoSQL database, millisecond latency
- Capacity modes: provisioned capacity with optional auto-scaling or on-demand capacity
- Can replace ElastiCache as a key/value store (storing session data for example, using TTL feature)
- Highly Available, Multi AZ by default, Read and Writes are decoupled, transaction capability
- DAX cluster for read cache, microsecond read latency
- Security, authentication and authorization is done through IAM
- Event Processing: DynamoDB Streams to integrate with AWS Lambda, or Kinesis Data Streams
- Global Table feature: active-active setup
- Automated backups up to 35 days with PITR (restore to new table), or on-demand backups
- Export to S3 without using RCU within the PITR window, import from S3 without using WCU
- Great to rapidly evolve schemas
- Use Case: Serverless applications development (small documents 100s KB), distributed serverless cache

## Slide 523

Amazon S3 – Summary
- S3 is a… key / value store for objects
- Great for bigger objects, not so great for many small objects
- Serverless, scales infinitely, max object size is 5 TB, versioning capability
- Tiers: S3 Standard, S3 Infrequent Access, S3 Intelligent, S3 Glacier + lifecycle policy
- Features: Versioning, Encr yption, Replication, MFA-Delete, Access Logs…
- Security: IAM, Bucket Policies, ACL, Access Points, Object Lambda, CORS, Object/Vault Lock
- Encryption: SSE-S3, SSE-KMS, SSE-C, client-side, TLS in transit, default encryption
- Batch operations on objects using S3 Batch, listing files using S3 Inventory
- Performance: Multi-part upload, S3 Transfer Acceleration, S3 Select
- Automation: S3 Event Notifications (SNS, SQS, Lambda, EventBridge)
- Use Cases: static files, key value store for big files, website hosting

## Slide 524

DocumentDB
- Aurora is an “AWS-implementation” of PostgreSQL / MySQL …
- DocumentDB is the same for MongoDB (which is a NoSQL database)
- MongoDB is used to store, quer y, and index JSON data
- Similar “deployment concepts” as Aurora
- Fully Managed, highly available with replication across 3 AZ
- DocumentDB storage automatically grows in increments of 10GB
- Automatically scales to workloads with millions of requests per seconds

## Slide 525

Amazon Neptune
- Fully managed graph database
- A popular graph dataset would be a social network
- Users have friends
- Posts have comments
- Comments have likes from users
- Users share and like posts…
- Highly available across 3 AZ, with up to 15 read replicas
- Build and run applications working with highly connected datasets – optimized for these complex and hard queries
- Can store up to billions of relations and query the graph with milliseconds latency
- Highly available with replications across multiple AZs
- Great for knowledge graphs (Wikipedia), fraud detection, recommendation engines, social networking

## Slide 526

Amazon Neptune – Streams
- Real-time ordered sequence of every change to your graph data
- Changes are available immediately after writing
- No duplicates, strict order
- Streams data is accessible in an HTTP REST API
- Use cases:
- Send notifications when certain changes are made
- Maintain your graph data synchronized in another data store (e.g., S3, OpenSearch, ElastiCache)
- Replicate data across regions in Neptune
Neptune Cluster
Neptune StreamsStreams API HTTP Get Request
…
OpenSearch
S3
ElastiCache…
Streams reader application
writes

## Slide 527

Amazon Keyspaces (for Apache Cassandra)
- Apache Cassandra is an open-source NoSQL distributed database
- A managed Apache Cassandra-compatible database service
- Serverless, Scalable, highly available, fully managed by AWS
- Automatically scale tables up/down based on the application’s traffic
- Tables are replicated 3 times across multiple AZ
- Using the Cassandra Query Language (CQL)
- Single-digit millisecond latency at any scale, 1000s of requests per second
- Capacity: On-demand mode or provisioned mode with auto-scaling
- Encryption, backup, Point-In-Time Recovery (PITR) up to 35 days
- Use cases: store IoT devices info, time-series data, …

## Slide 528

Amazon Timestream
- Fully managed, fast, scalable, serverless time series database
- Automatically scales up/down to adjust capacity
- Store and analyze trillions of events per day
- 1000s times faster & 1/10th the cost of relational databases
- Scheduled queries, multi-measure records, SQL compatibility
- Data storage tiering: recent data kept in memory and historical data kept in a cost-optimized storage
- Built-in time series analytics functions (helps you identify patterns in your data in near real-time)
- Encryption in transit and at rest
- Use cases: IoT apps, operational applications, real-time analytics, …

## Slide 529

Amazon Timestream – Architecture
AmazonTimestream
AWS IoT
Kinesis DataStreams
Lambda
Prometheus
Amazon MSK
Kinesis Data AnalyticsFor Apache Flink
Kinesis DataStreams
AmazonQuickSight
AmazonSageMaker
Any JDBC connection

## Slide 530

Data & Analytics

## Slide 531

Amazon Athena
- Serverless query service to analyze data stored in Amazon S3
- Uses standard SQL language to query the files (built on Presto)
- Supports CSV, JSON, ORC, Avro, and Parquet
- Pricing: $5.00 per TB of data scanned
- Commonly used with Amazon Quicksight for reporting/dashboards
- Use cases: Business intelligence / analytics / reporting, analyze & query VPC Flow Logs, ELB Logs, CloudTrail trails, etc...
- Exam Tip: analyze data in S3 using serverless SQL, use Athena
S3 Bucket
load data
AmazonAthena
Query & Analyze
AmazonQuickSight
Reporting & Dashboards

## Slide 532

Amazon Athena – Performance Improvement
- Use columnar data for cost-savings (less scan)
- Apache Parquet or ORC is recommended
- Huge performance improvement
- Use Glue to convert your data to Parquet or ORC
- Compress data for smaller retrievals (bzip2, gzip, lz4, snappy, zlip, zstd…)
- Partition datasets in S3 for easy querying on virtual columns
- s3://yourBucket/pathT oTable /<PARTITION_COLUMN_NAME>=<VALUE> /<PARTITION_COLUMN_NAME>=<VALUE> /<PARTITION_COLUMN_NAME>=<VALUE> /etc…
- Example: s3://athena-examples/flight/parquet/year=1991/month=1/day=1/
- Use larger files (> 128 MB) to minimize overhead

## Slide 533

Amazon Athena – Federated Query
- Allows you to run SQL queries across data stored in relational, non-relational, object, and custom data sources (AWS or on-premises)
- Uses Data Source Connectors that run onAWS Lambdato run Federated Queries (e.g., CloudWatch Logs, DynamoDB, RDS, …)
- Store the results back in Amazon S3
AmazonAthena
Database(On-Premises)
S3 BucketLambda(Data SourceConnector)ElastiCache
DocumentDB
DynamoDB
Redshift
HBase in EMR
MySQLAuroraSQL Server

## Slide 534

Redshift Overview
- Redshift is based on PostgreSQL, but it’s not used for OLTP
- It’s OLAP – online analytical processing (analytics and data warehousing)
- 10x better performance than other data warehouses, scale to PBs of data
- Columnar storage of data (instead of row based) & parallel query engine
- Two modes: Provisioned cluster or Ser ver less cluster
- Has a SQL interface for performing the queries
- BI tools such as Amazon Quicksight or Tableau integrate with it
- vs Athena: faster queries / joins / aggregations thanks to indexes

## Slide 535

Redshift Cluster
- Leader node: for query planning, results aggregation
- Compute node: for performing the queries, send results to leader
- Provisioned mode:
- Choose instance types in advance
- Can reserve instances for cost savings
QuerySELECT COUNT (*), …FROM MY_TABLEGROUP BY …
Amazon Redshift Cluster
JDBC/ODBC
Leader Node
Compute Nodes

## Slide 536

Redshift – Snapshots & DR
- Redshift has “Multi-AZ” mode for some clusters
- Snapshots are point-in-time backups of a cluster, stored internally in S3
- Snapshots are incremental (only what has changed is saved)
- Yo u c a n r e s t o r e a s n a p s h o t i n t o a new cluster
- Automated: every 8 hours, every 5 GB, or on a schedule. Set retention between 1 to 35 days
- Manual: snapshot is retained until you delete it
- Yo u c a n c o n fi g u r e A m a zo n R e d s h i f t t o automatically copy snapshots (automated or manual) of a cluster to another AWS Region
Region(us-east-1)
Redshift Cluster(Original)
Cluster SnapshotTake Snapshot
Region(eu-west-1)
Redshift Cluster(New)
Copied SnapshotRestore
Automated/ ManualCopy

## Slide 537

Loading data into Redshift: Large inserts are MUCH betterAmazon Kinesis Data FirehoseS3 using COPY command
Amazon KinesisData Firehose
Amazon RedshiftCluster(through S3 copy)
Amazon RedshiftClusterS3 Bucket(mybucket)
copy customerfrom 's3://mybucket/mydata’ iam_role 'arn:aws:iam::0123456789012:role/MyRedshiftRole';
InternetWithout Enhanced VPC Routing
With Enhanced VPC RoutingThrough VPC
EC2 InstanceJDBC driver
Amazon RedshiftClusterEC2 InstanceBetter to write Data in batches

## Slide 538

Redshift Spectrum
- Query data that is already in S3 without loading it
- Must have a Redshift cluster available to start the query
- The query is then submitted to thousands of Redshift Spectrum nodes
QuerySELECT COUNT (*), …FROM S3.EXT_TABLEGROUP BY …
Amazon Redshift Cluster
JDBC/ODBC
Leader Node
Compute Nodes
12….N
Redshift Spectrum
Amazon S3

## Slide 539

Amazon OpenSearch Service
- Amazon OpenSearch is successor to Amazon ElasticSearch
- In DynamoDB, queries only exist by primary key or indexes…
- With OpenSearch, you can search any field, even partially matches
- It’s common to use OpenSearch as a complement to another database
- Two modes: managed cluster or ser ver less cluster
- Does not natively support SQL (can be enabled via a plugin)
- Ingestion from Kinesis Data Firehose, AWS IoT, and CloudWatch Logs
- Security through Cognito & IAM, KMS encryption, TLS
- Comes with OpenSearch Dashboards (visualization)

## Slide 540

OpenSearch patternsDynamoDB
DynamoDB TableDynamoDB Stream
Lambda Function
Amazon OpenSearch
API to retrieve itemsAPI to search items
CRUD

## Slide 541

OpenSearch patternsCloudWatch LogsCloudWatch LogsSubscription FilterLambda Function(managed by AWS)
Amazon OpenSearch
CloudWatch LogsSubscription FilterKinesis Data FirehoseAmazon OpenSearch
Real time
Near Real Time

## Slide 542

OpenSearch patternsKinesis Data Streams & Kinesis Data Firehose
Kinesis DataStreams
Kinesis DataFirehose(near real time)
AmazonOpenSearch
datatransformation
LambdaFunction
AmazonOpenSearch
LambdaFunction(real time)
Kinesis DataStreams

## Slide 543

Amazon EMR
- EMR stands for “Elastic MapReduce”
- EMR helps creating Hadoop clusters (Big Data) to analyze and process vast amount of data
- The clusters can be made of hundreds of EC2 instances
- EMR comes bundled with Apache Spark, HBase, Presto, Flink…
- EMR takes care of all the provisioning and configuration
- Auto-scaling and integrated with Spot instances
- Use cases: data processing, machine learning, web indexing, big data…

## Slide 544

Amazon EMR – Node types & purchasing
- Master Node: Manage the cluster, coordinate, manage health – long running
- Core Node: Run tasks and store data – long running
- Task Node (optional): Just to run tasks – usually Spot
- Purchasing options:
- On-demand: reliable, predictable, won’t be terminated
- Reserved (min 1 year): cost savings (EMR will automatically use if available)
- Spot Instances: cheaper, can be terminated, less reliable
- Can have long-running cluster, or transient (temporary) cluster

## Slide 545

Amazon QuickSight
- Serverless machine learning-powered business intelligence service to create interactive dashboards
- Fast, automatically scalable, embeddable, with per-session pricing
- Use cases:
- Business analytics
- Building visualizations
- Perform ad-hoc analysis
- Get business insights using data
- Integrated with RDS, Aurora,Athena, Redshift, S3…
- In-memory computation using SPICEengine if data is imported into QuickSight
- Enterprise edition: Possibility to setup Column-Level security (CLS)
https://aws.amazon.com/quicksight/

## Slide 546

QuickSight Integrations
QuickSight
RDS
Redshift
Athena
S3
ELF & CLF(Log Format)
OpenSearch
Aurora
Timestream
Data Sources (AWS Services)
Data Sources (SaaS)
Data Sources (Imports)
On-PremisesDatabases (JDBC)

## Slide 547

QuickSight – Dashboard & Analysis
- Define Users (standard versions) and Groups (enterprise version)
- These users & groups only exist within QuickSight, not IAM !!
- Adashboard…
- is a read-only snapshot of an analysis that you can share
- preserves the configuration of the analysis (filtering, parameters, controls, sort)
- Yo u c a n s h a r e t h e a n a ly s i s o r t h e d a s h b o a r d w i t h U s e r s o r G r o u p s
- To s h a r e a d a s h b o a r d , yo u mu s t fi r s t p u bl i s h i t
- Users who see the dashboard can also see the underlying data

## Slide 548

AWS Glue
- Managed extract, transform, and load (ETL) service
- Useful to prepare and transform data for analytics
- Fully serverless service
S3 Bucket
Amazon RDS
ExtractTransform
Glue ETL
RedshiftData WarehouseLoad

## Slide 549

AWS Glue – Convert data into Parquet format
InputS3 Bucket
Import CSVParquet
OutputS3 Bucket
Glue ETL
AmazonAthenaAnalyze
Event notificationsOn S3 PUTLambda Function(EventBridge works as an alternative)
Trigger Glue ETL Job
S3 Put

## Slide 550

Glue Data Catalog: catalog of datasets
JDBC
Amazon S3
Amazon RDS
Amazon DynamoDB
AWS GlueData CatalogWrites Metadata
Database
Tables(Metadata)
Database
Tables(Metadata)
Amazon Athena
AmazonRedshiftSpectrum
Amazon EMR
AWS GlueData Crawler
Glue Jobs (ETL)
Data discovery

## Slide 551

Glue – things to know at a high-level
- Glue Job Bookmarks: prevent re-processing old data
- Glue DataBrew: clean and normalize data using pre-built transformation
- Glue Studio: new GUI to create, run and monitor ETL jobs in Glue
- Glue Streaming ETL (built on Apache Spark Structured Streaming): compatible with Kinesis Data Streaming, Kafka, MSK (managed Kafka)

## Slide 552

AWS Lake Formation
- Data lake = central place to have all your data for analytics purposes
- Fully managed service that makes it easy to setup a data lake in days
- Discover, cleanse, transform, and ingest data into your Data Lake
- It automates many complex manual steps (collecting, cleansing, moving, cataloging data, …) and de-duplicate (using ML Transforms)
- Combine structured and unstructured data in the data lake
- Out-of-the-box source blueprints: S3, RDS, Relational & NoSQL DB…
- Fine-grained Access Control for your applications (row and column-level)
- Built on top of AWS Glue

## Slide 553

AWS Lake Formation
Data Sources
Amazon S3
RDS
On-PremisesDatabase (SQL & NoSQL)AWS Lake Formation
Source CrawlersETL and Data Prep.Data CatalogSecurity SettingsAccess Control
Data Lake(stored in S3)
Athena
Redshift
EMR
Users
ingest
Aurora

## Slide 554

AWS Lake Formation Centralized Permissions Example
Data Sources
Amazon S3
RDS AWS Lake Formation
Access ControlColumn-level security
Data Lake(stored in S3)
Athena
Usersingest
Aurora
Quicksight

## Slide 555

Amazon Managed Service for Apache Flink
- Previously named: Kinesis Data Analytics for Apache Flink
- Flink (Java, Scala or SQL) is a framework for processing data streams
- Run any Apache Flink application on a managed cluster on AWS
- Provisioned compute resources, parallel computation, automatic scaling
- Application backups (implemented as checkpoints and snapshots)
- Use any Apache Flink programming features to transform data
- Important: Flink does not read from Amazon Data Firehose
Amazon MSK(Apache Kafka)
Amazon Managed Servicefor Apache Flink
Kinesis DataStreams

## Slide 556

Amazon Managed Streaming for Apache Kafka (Amazon MSK)
- Alternative to Amazon Kinesis
- Fully managed Apache Kafka on AWS
- Allow you to create, update, delete clusters
- MSK creates & manages Kafka brokers nodes & Zookeeper nodes for you
- Deploy the MSK cluster in your VPC, multi-AZ (up to 3 for HA)
- Automatic recovery from common Apache Kafka failures
- Data is stored on EBS volumes for as long as you want
- MSK Serverless
- Run Apache Kafka on MSK without managing the capacity
- MSK automatically provisions resources and scales compute & storage

## Slide 557

Apache Kafka at a high level
Kinesis
IoT
RDS
Etc…
MSK Cluster
Etc…
Broker 2
Broker 3
Broker 1Producers(your code)Consumers(your code)
EMR
S3
SageMaker
Kinesis
RDS
Write to topicPoll from topicreplication
replication

## Slide 558

Kinesis Data Streams vs. Amazon MSK
Kinesis Data Streams
- 1 MB message size limit
- Data Streams with Shards
- Shard Splitting & Merging
- TLS In-flight encryption
- KMS at-rest encryption
- 1MB default, configure for higher (ex: 10MB)
- Kafka T opics with Partitions
- Can only add partitions to a topic
- PLAINTEXT or TLS In-flight Encryption
- KMS at-rest encryption
Amazon MSK

## Slide 559

Amazon MSK Consumers
Amazon MSK
Kinesis Data Analytics for Apache Flink
AWS GlueStreaming ETL JobsPowered by Apache Spark StreamingLambda
Amazon EC2
Applications Running on
ECSEKS

## Slide 560

Big Data Ingestion Pipeline
- We want the ingestion pipeline to be fully ser verless
- We want to collect data in real time
- We want to transform the data
- We want to quer y the transformed data using SQL
- The reports created using the queries should be in S3
- We want to load that data into a warehouse and create dashboards

## Slide 561

Big Data Ingestion Pipeline
Amazon Kinesis Data Streams
Amazon Kinesis Data Firehose
AWS Lambda
Amazon Simple Storage Service (S3)Amazon Simple Queue Service
AWS Lambda
Amazon Athena
Amazon Simple Storage Service (S3)
Reporting Bucket
IoT Devices
Ingestion Bucket
(optional)
Every 1 minuteReal-time trigger
Pull data
Amazon QuickSight
Amazon RedshiftServerless

## Slide 562

Big Data Ingestion Pipeline discussion
- IoT Core allows you to harvest data from IoT devices
- Kinesis is great for real-time data collection
- Firehose helps with data delivery to S3 in near real-time (1 minute)
- Lambda can help Firehose with data transformations
- Amazon S3 can trigger notifications to SQS
- Lambda can subscribe to SQS (we could have connecter S3 to Lambda)
- Athena is a serverless SQL service and results are stored in S3
- The reporting bucket contains analyzed data and can be used by reporting tool such as AWS QuickSight, Redshift, etc…

## Slide 563

Machine Learning

## Slide 564

Amazon Rekognition
- Find objects, people, text, scenes in images and videos using ML
- Facial analysis and facial search to do user verification, people counting
- Create a database of “familiar faces” or compare against celebrities
- Use cases:
- Labeling
- Content Moderation
- Text Detection
- Face Detection and Analysis (gender, age range, emotions…)
- Face Search and Verification
- Celebrity Recognition
- Pathing (ex: for sports game analysis)

## Slide 565

Amazon Rekognition – Content Moderation
- Detect content that is inappropriate, unwanted, or offensive (image and videos)
- Used in social media, broadcast media, advertising, and e-commerce situations to create a safer user experience
- Set a Minimum Confidence Threshold for items that will be flagged
- Flag sensitive content for manual review in Amazon Augmented AI (A2I)
- Help comply with regulations
Image
AmazonRekognition
Confidence Leveland Threshold
Optional Manual review in A2I

## Slide 566

Amazon Transcribe
- Automatically convert speech to text
- Uses a deep learning process called automatic speech recognition (ASR) to convert speech to text quickly and accurately
- Automatically remove Personally Identifiable Information (PII) using Redaction
- Supports Automatic Language Identification for multi-lingual audio
- Use cases:
- transcribe customer service calls
- automate closed captioning and subtitling
- generate metadata for media assets to create a fully searchable archive
”Hello my name is Stéphane.I hope you’re enjoying the course!

## Slide 567

Amazon Polly
- Tur n text into lifelike speech using deep lear ning
- Allowing you to create applications that talk
Hi! My name is Stéphane and this is a demo of Amazon Polly

## Slide 568

Amazon Polly – Lexicon & SSML
- Customize the pronunciation of words with Pronunciation lexicons
- Stylized words: St3ph4ne => “Stephane”
- Acronyms: AWS => “Amazon Web Services”
- Upload the lexicons and use them in the SynthesizeSpeech operation
- Generate speech from plain text or from documents marked up with Speech Synthesis Markup Language (SSML) – enables more customization
- emphasizing specific words or phrases
- using phonetic pronunciation
- including breathing sounds, whispering
- using the Newscaster speaking style

## Slide 569

Amazon Translate
- Natural and accurate language translation
- Amazon Translate allows you to localize content - such as websites and applications - for international users, and to easily translate large volumes of text efficiently.

## Slide 570

Amazon Lex & Connect
- Amazon Lex: (same technology that powers Alexa)
- Automatic Speech Recognition (ASR) to convert speech to text
- Natural Language Understanding to recognize the intent of text, callers
- Helps build chatbots, call center bots
- Amazon Connect:
- Receive calls, create contact flows, cloud-based virtual contact center
- Can integrate with other CRM systems or AWS
- No upfront payments, 80% cheaper than traditional contact center solutions
Phone CallSchedule an Appointment
Connect
LexIntent recognized
Lambda
CRM
callstreaminvokeschedule

## Slide 571

Amazon Comprehend
- For Natural Language Processing – NLP
- Fully managed and serverless service
- Uses machine learning to find insights and relationships in text
- Language of the text
- Extracts key phrases, places, people, brands, or events
- Understands how positive or negative the text is
- Analyzes text using tokenization and parts of speech
- Automatically organizes a collection of text files by topic
- Sample use cases:
- analyze customer interactions (emails) to find what leads to a positive or negative experience
- Create and groups articles by topics that Comprehend will uncover

## Slide 572

Amazon Comprehend Medical
- Amazon Comprehend Medical detects and returns useful information in unstructured clinical text:
- Physician’s notes
- Discharge summaries
- Test results
- Case notes
- Uses NLP to detect Protected Health Information (PHI) – DetectPHI API
- Store your documents in Amazon S3, analyze real-time data with Kinesis Data Firehose, or use Amazon Transcribe to transcribe patient narratives into text that can be analyzed by Amazon Comprehend Medical.

## Slide 573

Amazon SageMaker AI
- Fully managed service for developers / data scientists to build ML models
- Typically, difficult to do all the processes in one place + provision ser ver s
- Machine learning process (simplified): predicting your exam score
Historical Data:# years of experience in IT# years of experience with AWSTime spent on the course…
label
670890934
build
ML model
Train and Tune
New data
Apply model
PredictionPASS WITH 906
score

## Slide 574

Amazon Kendra
- Fully managed document search service powered by Machine Learning
- Extract answers from within a document (text, pdf, HTML, PowerPoint, MS Word, FAQs…)
- Natural language search capabilities
- Learn from user interactions/feedback to promote preferred results (Incremental Learning)
- Ability to manually fine-tune search results (importance of data, freshness, custom, …)
Data Sources
Amazon S3Amazon RDSGoogle DriveMS SharePoint
MS OneDrive3rd party, APNs, Custom
Knowledge Index(powered by ML)Amazon Kendra
indexing
User
Where is the IT support desk?
1st floor

## Slide 575

Amazon Personalize
- Fully managed ML-service to build apps with real-time personalized recommendations
- Example: personalized product recommendations/re-ranking, customized direct marketing
- Example: User bought gardening tools, provide recommendations on the next one to buy
- Same technology used by Amazon.com
- Integrates into existing websites, applications, SMS, email marketing systems, …
- Implement in days, not months (you don’t need to build, train, and deploy ML solutions)
- Use cases: retail stores, media and entertainment…
Amazon S3Amazon Personalize API
read data from S3
real-time data integrationAmazon Personalize
Customized personalized API
Websites & AppsMobile Apps
SMSEmails

## Slide 576

Amazon Textract
- Automatically extracts text, handwriting, and data from any scanned documents using AI and ML
- Extract data from forms and tables
- Read and process any type of document (PDFs, images, …)
- Use cases:
- Financial Services (e.g., invoices, financial reports)
- Healthcare (e.g., medical records, insurance claims)
- Public Sector (e.g., tax forms, ID documents, passports)
Amazon Textract
{ “Document ID”: “123456789-005”, “Name”: “”, “SEX”: “F”, “DOB”: “23.05.1997”, …}
analyzeresult

## Slide 577

AWS Machine Learning - Summary
- Rekognition: face detection, labeling, celebrity recognition
- Transcribe: audio to text (ex: subtitles)
- Polly: text to audio
- Translate: translations
- Lex: build conversational bots – chatbots
- Connect: cloud contact center
- Comprehend: natural language processing
- SageMaker: machine learning for ever y developer and data scientist
- Kendra: ML-powered search engine
- Personalize: real-time personalized recommendations
- Textr act: detect text and data in documents

## Slide 578

AWS Monitoring, Audit and PerformanceCloudWatch, CloudTrail & AWS Config

## Slide 579

Amazon CloudWatch Metrics
- CloudWatch provides metrics for every services in AWS
- Metric is a variable to monitor (CPUUtilization, NetworkIn…)
- Metrics belong to namespaces
- Dimension is an attribute of a metric (instance id, environment, etc…).
- Up to 30 dimensions per metric
- Metrics have timestamps
- Can create CloudWatch dashboards of metrics
- Can create CloudWatch Custom Metrics (for the RAM for example)

## Slide 580

CloudWatch Metric Streams
- Continually stream CloudWatch metrics to a destination of your choice, with near-real-time delivery and low latency.
- Amazon Kinesis Data Firehose (and then its destinations)
- 3rd party service provider: Datadog, Dynatrace, New Relic, Splunk, Sumo Logic…
- Option to filter metrics to only stream a subset of them
CloudWatch Metrics
Amazon S3AmazonRedshiftAmazon OpenSearch
Kinesis Data Firehose
Athena
Stream near-real-time

## Slide 581

CloudWatch Logs
- Log groups: arbitrary name, usually representing an application
- Log stream: instances within application / log files / containers
- Can define log expiration policies (never expire, 1 day to 10 years…)
- CloudWatch Logs can send logs to:
- Amazon S3 (exports)
- Kinesis Data Streams
- Kinesis Data Firehose
- AWS Lambda
- OpenSearch
- Logs are encrypted by default
- Can setup KMS-based encryption with your own keys

## Slide 582

CloudWatch Logs - Sources
- SDK, CloudWatch Logs Agent, CloudWatch Unified Agent
- Elastic Beanstalk: collection of logs from application
- ECS: collection from containers
- AWS Lambda: collection from function logs
- VPC Flow Logs: VPC specific logs
- API Gateway
- CloudTrail based on filter
- Route53: Log DNS queries

## Slide 583

CloudWatch Logs Insights
https://mng.workshop.aws/operations-2022/detect/cwlogs.html

## Slide 584

CloudWatch Logs Insights
- Search and analyze log data stored in CloudWatch Logs
- Example: find a specific IP inside a log, count occurrences of “ERROR” in your logs…
- Provides a purpose-built query language
- Automatically discovers fields from AWS services and JSON log events
- Fetch desired event fields, filter based on conditions, calculate aggregate statistics, sort events, limit number of events…
- Can save queries and add them to CloudWatch Dashboards
- Can query multiple Log Groups in different AWS accounts
- It’s a query engine, not a real-time engine

## Slide 585

CloudWatch Logs – S3 Export
- Log data can take up to 12 hours to become available for export
- The API call is CreateExportTask
- Not near-real time or real-time… use Logs Subscriptions insteadCloudWatch Logs
Amazon S3

## Slide 586

CloudWatch Logs Subscriptions
- Get a real-time log events from CloudWatch Logs for processing and analysis
- Send to Kinesis Data Streams, Kinesis Data Firehose, or Lambda
- Subscription Filter – filter which logs are events delivered to your destination
CloudWatch LogsSubscription Filter
Lambda
Kinesis Data Firehose
Kinesis Data Streams
OpenSearchService
S3logs
real-time
nearreal-time
KDFKDAEC2Lambda…

## Slide 587

CloudWatch Logs AggregationMulti-Account & Multi Region
CloudWatch LogsSubscription Filter
CloudWatch LogsSubscription Filter
CloudWatch LogsSubscription Filter
Kinesis Data StreamsKinesis Data Firehose
Near Real Time
Amazon S3
ACCOUNT AREGION 1
ACCOUNT BREGION 2
ACCOUNT BREGION 3

## Slide 588

CloudWatch Logs Subscriptions
- Cross-Account Subscription – send log events to resources in a different AWS account (KDS, KDF)Account – Sender(111111111111)
CloudWatchLogsSubscriptionFilter
logs
Account – Recipient(999999999999)
Kinesis Data Streams(RecipientStream)
SubscriptionDestination
DestinationAccess PolicyIAM Role
IAM Role(Cross-Account)
DestinationAccess Policy
logs
Can be assumedallow PutRecord

## Slide 589

CloudWatch Logs for EC2
- By default, no logs from your EC2 machine will go to CloudWatch
- Yo u n e e d t o r u n a C l o u d Wa t c h agent on EC2 to push the log files you want
- Make sure IAM permissions are correct
- The CloudWatch log agent can be setup on-premises tooEC2 Instance
CloudWatch Logs AgentOn Premise Server
CloudWatch Logs Agent
CloudWatch Logs

## Slide 590

CloudWatch Logs Agent & Unified Agent
- For virtual servers (EC2 instances, on-premises servers…)
- CloudWatch Logs Agent
- Old version of the agent
- Can only send to CloudWatch Logs
- CloudWatch Unified Agent
- Collect additional system-level metrics such as RAM, processes, etc…
- Collect logs to send to CloudWatch Logs
- Centralized configuration using SSM Parameter Store

## Slide 591

CloudWatch Unified Agent – Metrics
- Collected directly on your Linux server / EC2 instance
- CPU (active, guest, idle, system, user, steal)
- Disk metrics (free, used, total), Disk IO (writes, reads, bytes, iops)
- RAM (free, inactive, used, total, cached)
- Netstat (number of TCP and UDP connections, net packets, bytes)
- Processes (total, dead, bloqued, idle, running, sleep)
- Swap Space (free, used, used %)
- Reminder: out-of-the box metrics for EC2 – disk, CPU, network (high level)

## Slide 592

CloudWatch Alarms
- Alarms are used to trigger notifications for any metric
- Various options (sampling, %, max, min, etc…)
- Alarm States:
- OK
- INSUFFICIENT_DATA
- ALARM
- Period:
- Length of time in seconds to evaluate the metric
- High resolution custom metrics: 10 sec, 30 sec or multiples of 60 sec

## Slide 593

CloudWatch Alarm Targets
- Stop, Terminate, Reboot, or Recover an EC2 Instance
- Trigger Auto Scaling Action
- Send notification to SNS (from which you can do pretty much anything)
Amazon EC2EC2 Auto ScalingAmazon SNS

## Slide 594

CloudWatch Alarms – Composite Alarms
- CloudWatch Alarms are on a single metric
- Composite Alarms are monitoring the states of multiple other alarms
- AND and OR conditions
- Helpful to reduce “alarm noise” by creating complex composite alarms
CW Alarm - A
EC2 Instance
CW Alarm - B
monitor CPU
monitor IOPS
Composite Alarm
Amazon SNS
ALARM
ALARM
trigger

## Slide 595

EC2 Instance Recovery
CloudWatch AlarmStatusCheckFailed_System
- Status Check:
- Instance status = check the EC2 VM
- System status = check the underlying hardware
- Attached EBS status = check attached EBS volumes
EC2 Instance
monitor
EC2 Instance Recovery
SNS Topic
alert
- Recovery: Same Private, Public, Elastic IP , metadata, placement group

## Slide 596

CloudWatch Alarm: good to know
- Alarms can be created based on CloudWatch Logs Metrics Filters
- To t e s t a l a r m s a n d n o t i fi c a t i o n s , s e t t h e a l a r m s t a t e t o A l a r m u s i n g C L Iaws cloudwatch set-alarm-state --alarm-name "myalarm" --state-value ALARM --state-reason "testing purposes"
CW Logs
CloudWatch
CW AlarmAlert
Amazon SNS
Metric Filter

## Slide 597

CloudWatch Network Synthetic Monitor
- Monitor and detects network issues between your apps hosted on AWS and your on-premises data center
- Identify any network performance degradation (e.g., packet loss, latency, jitter…)
- No agents required to be installed
- Tests ICMP or TCP traffic to IPv4/IPv4 on-premises destinations through Direct Connect or S2S VPN connections
- Publishes data to CloudWatch Metrics
AWS Cloud
Corporate Data Center
Private Subnet
Server
VPN ConnectionorDX Connection
EC2 instance
CloudWatchMetrics

## Slide 598

Amazon EventBridge (formerly CloudWatch Events)
- Schedule: Cron jobs (scheduled scripts)
- Event Pattern: Event rules to react to a service doing something
- Trigger Lambda functions, send SQS/SNS messages…
IAM Root User Sign in EventSNS Topic with Email Notification
Schedule Every hourTrigger script on Lambda function

## Slide 599

Amazon EventBridge Rules
Example Source
EC2 Instance(ex: Start Instance)CodeBuild(ex: failed build)
S3 Event(ex: upload object)Trusted Advisor(ex: new Finding)
CloudTrail(any API call)Schedule or Cron(ex: every 4 hours)
AmazonEventBridge
Filter events(optional)
Example Destinations
{ "version": "0", "id": "6a7e8feb-b491", "detail-type": "EC2 Instance State-change Notification", ….}
JSON
LambdaAWS BatchECS Task
SQSSNSKinesis DataStreams
StepFunctionsCodePipelineCodeBuild
SSMEC2 Actions
ComputeIntegrationOrchestrationMaintenance

## Slide 600

Amazon EventBridge
- Event buses can be accessed by other AWS accounts using Resource-based Policies
- Yo u c a n archive events (all/filter) sent to an event bus (indefinitely or set period)
- Ability to replay archived events
AWS Services
Default Event BusPartner Event BusCustom Event BusAWS SaaS Partners
Custom Apps

## Slide 601

Amazon EventBridge – Schema Registry
- EventBridge can analyze the events in your bus and infer the schema
- The Schema Registry allows you to generate code for your application, that will know in advance how data is structured in the event bus
- Schema can be versioned

## Slide 602

Amazon EventBridge – Resource-based Policy
- Manage permissions for a specific Event Bus
- Example: allow/deny events from another AWS account or AWS region
- Use case: aggregate all events from your AWS Organization in a single AWS account or AWS region
Allow events from another AWS account
AWS Account(123456789012)
AWS Account(111122223333)
EventBridge Bus(central-event-bus)
Lambda functionPutEvents

## Slide 603

CloudWatch Container Insights
- Collect, aggregate, summarize metrics and logs from containers
- Available for containers on…
- Amazon Elastic Container Service (Amazon ECS)
- Amazon Elastic Kubernetes Services (Amazon EKS)
- Kubernetes platforms on EC2
- Fargate (both for ECS and EKS)
- In Amazon EKS and Kubernetes, CloudWatch Insights is using a containerized version of the CloudWatch Agent to discover containers
ECS Container
EKS Container
Metrics and logs
CloudWatch Container Insights

## Slide 604

CloudWatch Lambda Insights
- Monitoring and troubleshooting solution for serverless applications running on AWS Lambda
- Collects, aggregates, and summarizes system-level metrics including CPU time, memory, disk, and network
- Collects, aggregates, and summarizes diagnostic information such as cold starts and Lambda worker shutdowns
- Lambda Insights is provided as a Lambda Layer

## Slide 605

CloudWatch Contributor Insights
- Analyze log data and create time series that display contributor data.
- See metrics about the top-N contributors
- The total number of unique contributors, and their usage.
- This helps you find top talkers and understand who or what is impacting system performance.
- Works for any AWS-generated logs (VPC, DNS, etc..)
- For example, you can find bad hosts, identify the heaviest network users, or find the URLs that generate the most errors.
- Yo u c a n bu i l d yo u r r u l e s f r o m s c r a t c h , o r yo u c a n a l s o use sample rules that AWS has created – leverages your CloudWatch Logs
- CloudWatch also provides built-in rules that you can use to analyze metrics from other AWS services.
VPC Flow Logs
CloudWatch Logs
CloudWatch Contributor Insights
Top-10 IP addresses

## Slide 606

CloudWatch Application Insights
- Provides automated dashboards that show potential problems with monitored applications, to help isolate ongoing issues
- Y our applications run on Amazon EC2 Instances with select technologies only (Java, .NET, Microsoft IIS Web Server, databases…)
- And you can use other AWS resources such as Amazon EBS, RDS, ELB, ASG, Lambda, SQS, DynamoDB, S3 bucket, ECS, EKS, SNS, API Gateway…
- Powered by SageMaker
- Enhanced visibility into your application health to reduce the time it will take you to troubleshoot and repair your applications
- Findings and alerts are sent to Amazon EventBridge and SSM OpsCenter

## Slide 607

CloudWatch Insights and Operational Visibility
- CloudWatch Container Insights
- ECS, EKS, Kubernetes on EC2, Fargate, needs agent for Kubernetes
- Metrics and logs
- CloudWatch Lambda Insights
- Detailed metrics to troubleshoot serverless applications
- CloudWatch Contributors Insights
- Find “T op-N” Contributors through CloudWatch Logs
- CloudWatch Application Insights
- Automatic dashboard to troubleshoot your application and related AWS services

## Slide 608

AWS CloudTrail
- Provides governance, compliance and audit for your AWS Account
- CloudTrail is enabled by default!
- Get an history of events / API calls made within your AWS Account by:
- Console
- SDK
- CLI
- AWS Services
- Can put logs from CloudTrail into CloudWatch Logs or S3
- A trail can be applied to All Regions (default) or a single Region.
- If a resource is deleted in AWS, investigate CloudTrail first!

## Slide 609

CloudTrail Diagram
SDK
CLI
Console
CloudTrail Console
Inspect & Audit
CloudWatch Logs
S3 Bucket
IAM Users &IAM Roles

## Slide 610

CloudTrail Events
- Management Events:
- Operations that are performed on resources in your AWS account
- Examples:
- Configuring security (IAMAttachRolePolicy)
- Configuring rules for routing data (Amazon EC2CreateSubnet)
- Setting up logging (AWS CloudTrailCreateTrail)
- By default, trails are configured to log management events.
- Can separate Read Events (that don’t modify resources) from Write Events (that may modify resources)
- Data Events:
- By default, data events are not logged (because high volume operations)
- Amazon S3 object-level activity (ex: GetObject, DeleteObject, PutObject): can separate Read and Write Events
- AWS Lambda function execution activity (the Invoke API)
- CloudTrail Insights Events:
- See next slide J

## Slide 611

CloudTrail Insights
- Enable CloudTrail Insights to detect unusual activity in your account:
- inaccurate resource provisioning
- hitting service limits
- Bursts of AWS IAM actions
- Gaps in periodic maintenance activity
- CloudTrail Insights analyzes normal management events to create a baseline
- And then continuously analyzes write events to detect unusual patterns
- Anomalies appear in the CloudTrail console
- Event is sent to Amazon S3
- An EventBridge event is generated (for automation needs)
Management EventsCloudTrail InsightsContinous analysisInsights Events
generate
CloudTrail Console
S3 Bucket
EventBridge event

## Slide 612

CloudTrail Events Retention
- Events are stored for 90 days in CloudTrail
- T o keep events beyond this period, log them to S3 and use Athena
CloudTrail
S3 BucketLong-term retention
Data Events
Management Events
Insights Events
log
90 days retention
analyzeAthena

## Slide 613

Amazon EventBridge – Intercept API Calls
CloudTrail(any API call)AmazonEventBridge
SNS
event
DynamoDB
Log API call
User
alert
DeleteTable API Call
💥

## Slide 614

Amazon EventBridge + CloudTrail
CloudTrail
EventBridge
IAMIAM Role
UserAssumeRole
API Call logs
event
SNS
CloudTrailEC2
UserAuthorizeSecurityGroupIngressAPI Call logs
event
Security Groupedit SGInbound Rules
EventBridge
SNS

## Slide 615

AWS Config
- Helps with auditing and recording compliance of your AWS resources
- Helps record configurations and changes over time
- Questions that can be solved by AWS Config:
- Is there unrestricted SSH access to my security groups?
- Do my buckets have any public access?
- How has my ALB configuration changed over time?
- Yo u c a n r e c e i ve a l e r t s ( S N S n o t i fi c a t i o n s ) fo r a ny c h a n g e s
- AWS Config is a per-region service
- Can be aggregated across regions and accounts
- Possibility of storing the configuration data into S3 (analyzed by Athena)

## Slide 616

Config Rules
- Can use AWS managed config rules (over 75)
- Can make custom config rules (must be defined in AWS Lambda)
- Ex: evaluate if each EBS disk is of type gp2
- Ex: evaluate if each EC2 instance is t2.micro
- Rules can be evaluated / triggered:
- For each config change
- And / or: at regular time intervals
- AWS Config Rules does not prevent actions from happening (no deny)
- Pricing: no free tier, $0.003 per configuration item recorded per region, $0.001 per config rule evaluation per region

## Slide 617

AWS Config Resource
- View compliance of a resource over time
- View configuration of a resource over time
- View CloudTrail API calls of a resource over time

## Slide 618

Config Rules – Remediations
- Automate remediation of non-compliant resources using SSM Automation Documents
- Use AWS-Managed Automation Documents or create custom Automation Documents
- Tip: you can create custom Automation Documents that invokes Lambda function
- Yo u c a n s e t Remediation Retries if the resource is still non-compliant after auto-remediation
IAM Access Key(NON_COMPLIANT)
expired
AWS Config
monitor Auto-Remediation Action(SSM Document: AWSConfigRemediation-RevokeUnusedIAMUserCredentials)trigger
Retries: 5
deactivate

## Slide 619

Config Rules – Notifications
- Use EventBridge to trigger notifications when AWS resources are non-compliant
- Ability to send configuration changes and compliance state notifications to SNS (all events – use SNS Filtering or filter at client-side)
AWS Resources
Security group…monitor
AWS ConfigEventBridge
NON_COMPLIANT
trigger
LambdaSNSSQS…
monitor
AWS ConfigSNS
trigger
AdminAll events (configuration changes,compliance state…)
AWS Resources
Security group… notification

## Slide 620

CloudWatch vs CloudTrail vs Config
- CloudWatch
- Performance monitoring (metrics, CPU, network, etc…) & dashboards
- Events & Alerting
- Log Aggregation & Analysis
- CloudTrail
- Record API calls made within your Account by everyone
- Can define trails for specific resources
- Global Service
- Config
- Record configuration changes
- Evaluate resources against compliance rules
- Get timeline of changes and compliance

## Slide 621

For an Elastic Load Balancer
- CloudWatch:
- Monitoring Incoming connections metric
- Visualize error codes as % over time
- Make a dashboard to get an idea of your load balancer performance
- Config:
- Track security group r ules for the Load Balancer
- Track configuration changes for the Load Balancer
- Ensure an SSL certificate is always assigned to the Load Balancer (compliance)
- CloudTrail:
- Track who made any changes to the Load Balancer with API calls

## Slide 622

Advanced Identity in AWS

## Slide 623

AWS Organizations
- Global service
- Allows to manage multiple AWS accounts
- The main account is the management account
- Other accounts are member accounts
- Member accounts can only be part of one organization
- Consolidated Billing across all accounts - single payment method
- Pricing benefits from aggregated usage (volume discount for EC2, S3…)
- Shared reserved instances and Savings Plans discounts across accounts
- API is available to automate AWS account creation

## Slide 624

AWS OrganizationsRoot Organizational Unit (OU)
Management AccountOU (Dev)
Member Accounts
OU (Prod)
OU (HR)
OU (Finance)

## Slide 625

Organizational Units (OU) - ExamplesBusiness Unit
Management Account
Sales OU
Retail OU
Finance OU
SalesAccount 1SalesAccount 2
RetailAccount 1RetailAccount 2
FinanceAccount 1FinanceAccount 2
Environmental Lifecycle
Management Account
Prod OU
Dev OU
Test OU
ProdAccount 1ProdAccount 2
DevAccount 1DevAccount 2
TestAccount 1TestAccount 2
Project-Based
Management Account
Project 1 OU
Project 2 OU
Project 3 OU
Project 1Account 1Project 1Account 2
Project 2Account 1Project 2Account 2
Project 3Account 1Project 3Account 2

## Slide 626

AWS Organizations
- Advantages
- Multi Account vs One Account Multi VPC
- Use tagging standards for billing purposes
- Enable CloudTrail on all accounts, send logs to central S3 account
- Send CloudWatch Logs to central logging account
- Establish Cross Account Roles for Admin purposes
- Security: Service Control Policies (SCP)
- IAM policies applied to OU or Accounts to restrict Users and Roles
- They do not apply to the management account (full admin power)
- Must have an explicit allow from the root through each OU in the direct path to the target account (does not allow anything by default – like IAM)

## Slide 627

SCP HierarchyOU (Root)
Management Account
OU (Sandbox)
OU (Test)
Account AAccount D
FullAWSAccessDeny AthenaFullAWSAccess + Deny S3
FullAWSAccess + Deny EC2
- Management Account
- Can do anything (no SCP apply)
- Account A
- Can do anything
- EXCEPT S3 (explicit Deny from Sandbox OU)
- EXCEPT EC2 (explicit Deny)
- Account B & C
- Can do anything
- EXCEPT S3 (explicit Deny from Sandbox OU)
- Account D
- Can access EC2
- Prod OU & Account E & F
- Can do anything
Account B
Account C
OU (Workloads)
OU (Prod)
Account E
Account F
FullAWSAccessAllow EC2
FullAWSAccess

## Slide 628

SCP ExamplesBlocklist and Allowlist strategies
More examples: https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_example-scps.html

## Slide 629

AWS Organizations – Tag Policies
- Helps you standardize tags across resources in an AWS Organization
- Ensure consistent tags, audit tagged resources, maintain proper resources categorization, …
- Yo u d e fi n e t a g ke y s a n d t h e i r a l l owe d v a l u e s
- Helps with AWS Cost Allocation Tags and Attribute-based Access Control
- Prevent any non-compliant tagging operations on specified services and resources (has no effect on resources without tags)
- Generate a report that lists all tagged/non-compliant resources
- Use EventBridge to monitor non-compliant tags

## Slide 630

IAM Conditionsaws:SourceIprestrict the client IP from which the API calls are being made
aws:RequestedRegionrestrict the region theAPI calls are made to

## Slide 631

IAM Conditionsec2:ResourceTagrestrict based on tags aws:MultiFactorAuthPresentto force MFA

## Slide 632

IAM for S3
- s3:ListBucket permission applies toarn:aws:s3:::test
- => bucket level permission
- s3:GetObject, s3:PutObject, s3:DeleteObject applies toarn:awn:s3:::test/*
- => object level permission

## Slide 633

Resource Policies & aws:PrincipalOrgID
- aws:PrincipalOrgID can be used in any resource policies to restrict access to accounts that are member of an AWS Organization
S3 Bucket(2022-financial-data)
AWS Organization(o-yyyyyyyyyy)
…Member Accounts
User outside Organization

## Slide 634

IAM Roles vs Resource Based Policies
- Cross account:
- attaching a resource-based policy to a resource (example: S3 bucket policy)
- OR using a role as a proxy
UserAccount A
Amazon S3
RoleAccount B
UserAccount A
Amazon S3
S3 Bucket Policy

## Slide 635

IAM Roles vs Resource-Based Policies
- When you assume a role (user, application or service), you give up your original permissions and take the permissions assigned to the role
- When using a resource-based policy, the principal doesn’t have to give up his permissions
- Example: User in account A needs to scan a DynamoDB table in Account A and dump it in an S3 bucket in Account B.
- Supported by: Amazon S3 buckets, SNS topics, SQS queues, etc…

## Slide 636

Amazon EventBridge – Security
- When a rule runs, it needs permissions on the target
- Resource-based policy: Lambda, SNS, SQS, S3 buckets, API Gateway…
- IAM role: EC2 Auto Scaling, Systems Manager Run Command, ECS task…
EventBridgeRule
Lambda with Resource based Policye.g. Allow EventBridge
IAM Role
EC2 Auto ScalingEventBridgeRule

## Slide 637

IAM Permission Boundaries
- IAM Permission Boundaries are supported for users and roles (not groups)
- Advanced feature to use a managed policy to set the maximum permissions an IAM entity can get.
IAM Permission Boundary
+
IAM Permissions Through IAM Policy
=No PermissionsExample:

## Slide 638

IAM Permission Boundaries
- Can be used in combinations of AWS Organizations SCP
https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_boundaries.html
Use cases
- Delegate responsibilities to non administrators within their permission boundaries, for example create new IAM users
- Allow developers to self-assign policies and manage their own permissions, while making sure they can’t “escalate” their privileges (= make themselves admin)
- Useful to restrict one specific user (instead of a whole account using Organizations & SCP)

## Slide 639

IAM Policy Evaluation Logic
https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_evaluation-logic.html

## Slide 640

Example IAM Policy
- Can you perform sqs:CreateQueue?
- Can you perform sqs:DeleteQueue?
- Can you perform ec2:DescribeInstances?

## Slide 641

AWS IAM Identity Center(successor to AWS Single Sign-On)
- One login (single sign-on) for all your
- AWS accounts in AWS Organizations
- Business cloud applications (e.g., Salesforce, Box, Microsoft 365, …)
- SAML2.0-enabled applications
- EC2 Windows Instances
- Identity providers
- Built-in identity store in IAM Identity Center
- 3rd party: Active Directory (AD), OneLogin, Okta…

## Slide 642

AWS IAM Identity Center – Login Flow
AWS IAM Identity Center

## Slide 643

AWS IAM Identity Center
AWS Organization
Business Cloud Apps
Custom SAML2.0-enabled Apps
AWS IAM Identity Center
Permission Sets
Browser Interface
Windows EC2login
Store / retrieveUser identities
SSO
Active DirectoryUsers & groups(On-premises, cloud)
IAM Identity CenterBuilt-in Identity Store
AWS Cloud

## Slide 644

IAM Identity CenterAWS Organization
Management Account
OU (Development)
Dev Account A
Dev Account B
OU (Production)
Prod Account A
Prod Account B
IAM Identity Center(in Management account)Group (Developers)
Bob
Alice
Permission SetReadOnlyAccess
Permission SetFullAccess
assignassign

## Slide 645

AWS IAM Identity CenterFine-grained Permissions and Assignments
- Multi-Account Permissions
- Manage access across AWS accounts in your AWS Organization
- Permission Sets – a collection of one or more IAM Policies assigned to users and groups to define AWS access
- Application Assignments
- SSO access to many SAML 2.0 business applications (Salesforce, Box, Microsoft 365, …)
- Provide required URLs, certificates, and metadata
- Attribute-Based Access Control (ABAC)
- Fine-grained permissions based on users’ attributes stored in IAM Identity Center Identity Store
- Example: cost center, title, locale, …
- Use case: Define permissions once, then modify AWS access by changing the attributes
IAM Identity Center
Permission Sets(DB Admins)
AWS Organization
DevAccount
ProdAccount
RDSRDS
DatabaseAdmins
Aurora
Aurora
IAM RoleIAM Role
assume
Permission Sets(DB Admins)

## Slide 646

What is Microsoft Active Directory (AD)?
- Found on any Windows Server with AD Domain Services
- Database of objects: User Accounts, Computers, Printers, File Shares, Security Groups
- Centralized security management, create account, assign permissions
- Objects are organized in trees
- A group of trees is a forest
Domain Controller
JohnPassword

## Slide 647

AWS Directory Services
- AWS Managed Microsoft AD
- Create your own AD in AWS, manage users locally, supports MFA
- Establish “trust” connections with your on-premises AD
- AD Connector
- Directory Gateway (proxy) to redirect to on-premises AD, supports MFA
- Users are managed on the on-premises AD
- Simple AD
- AD-compatible managed directory on AWS
- Cannot be joined with on-premises AD
On-prem ADAWS Managed AD
trust
On-prem ADAD Connector
proxy
Simple AD
authauth
auth

## Slide 648

IAM Identity Center – Active Directory Setup
- Connect to an AWS Managed Microsoft AD (Directory Service)
- Integration is out of the box
- Connect to a Self-Managed Directory
- Create Two-way Trust Relationship using AWS Managed Microsoft AD
- Create an AD Connector
IAM IdentityCenter
AWS ManagedMicrosoft ADconnect
IAM IdentityCenter
AWS ManagedMicrosoft ADtwo-way trust relationship
AD Connectorproxy
connectconnect

## Slide 649

AWS Control Tower
- Easy way to set up and govern a secure and compliant multi-account AWS environment based on best practices
- AWS Control Tower uses AWS Organizations to create accounts
- Benefits:
- Automate the set up of your environment in a few clicks
- Automate ongoing policy management using guardrails
- Detect policy violations and remediate them
- Monitor compliance through an interactive dashboard

## Slide 650

AWS Control Tower – Guardrails
- Provides ongoing governance for your Control T ower environment (AWS Accounts)
- Preventive Guardrail – using SCPs (e.g., Restrict Regions across all your accounts)
- Detective Guardrail – using AWS Config (e.g., identify untagged resources)
AWS Control Tower
MemberAccounts
Guardrail(Detective) SNS
Lambda
monitor un-taggedresources
trigger(NON_COMPLIANT)Adminnotify
invokeremediate(add tags)
AWS Config

## Slide 651

AWS Security & EncryptionKMS, Encryption SDK, SSM Parameter Store

## Slide 652

Why encryption?Encryption in flight (TLS / SSL)
- Data is encrypted before sending and decrypted after receiving
- TLS certificates help with encryption (HTTPS)
- Encryption in flight ensures no MITM (man in the middle attack) can happen
HTTPS Website
TLS EncryptionTLS Decryption
Username: adminPassword: supersecretUsername: adminPassword: supersecretaGVsbG8gd29ybGQgZWh…
Client ServerClient

## Slide 653

Why encryption?Server-side encryption at rest
- Data is encrypted after being received by the server
- Data is decrypted before being sent
- It is stored in an encrypted form thanks to a key (usually a data key)
- The encryption / decryption keys must be managed somewhere, and the server must have access to it
AWS Service (e.g., S3)
+HTTP(S) HTTP(S)ObjectData key Object
+Data key
Encryption
Decryption

## Slide 654

Why encryption?Client-side encryption
- Data is encrypted by the client and never decrypted by the server
- Data will be decrypted by a receiving client
- The server should not be able to decrypt the data
- Could leverage Envelope Encryption
Client Any storage service (FTP , S3, …)
+Data key(client-side)
Encryption
Object
store
Encrypted object
+
Data key(client-side)
Object
Decryptionretrieve

## Slide 655

AWS KMS (Key Management Service)
- Anytime you hear “encryption” for an AWS service, it’s most likely KMS
- AWS manages encryption keys for us
- Fully integrated with IAM for authorization
- Easy way to control access to your data
- Able to audit KMS Key usage using CloudTrail
- Seamlessly integrated into most AWS services (EBS, S3, RDS, SSM…)
- Never ever store your secrets in plaintext, especially in your code!
- KMS Key Encryption also available through API calls (SDK, CLI)
- Encrypted secrets can be stored in the code / environment variables

## Slide 656

KMS Keys Types
- KMS Keys is the new name of KMS Customer Master Key
- Symmetric (AES-256 keys)
- Single encryption key that is used to Encrypt and Decrypt
- AWS services that are integrated with KMS use Symmetric CMKs
- Yo u n e ve r g e t a c c e s s t o t h e K M S K e y u n e n c r y p t e d ( mu s t c a l l K M S A P I t o u s e )
- Asymmetric (RSA & ECC key pairs)
- Public (Encrypt) and Private Key (Decrypt) pair
- Used for Encrypt/Decrypt, or Sign/Verify operations
- The public key is downloadable, but you can’t access the Private Key unencrypted
- Use case: encryption outside of AWS by users who can’t call the KMS API

## Slide 657

AWS KMS (Key Management Service)
- Types of KMS Keys:
- AWS Owned Keys (free): SSE-S3, SSE-SQS, SSE-DDB (default key)
- AWS Managed Key: free (aws/service-name, example: aws/rds or aws/ebs)
- Customer managed keys created in KMS: $1 / month
- Customer managed keys imported: $1 / month
- + pay for API call to KMS ($0.03 / 10000 calls)
- Automatic Key rotation:
- AWS-managed KMS Key: automatic every 1 year
- Customer-managed KMS Key: (must be enabled) automatic & on-demand
- Imported KMS Key: only manual rotation possible using alias

## Slide 658

Copying Snapshots across regions
EBS VolumeEncryptedWith KMS
EBS SnapshotEncryptedWith KMS
KMS Key A
KMS Key A
Region eu-west-2
EBS VolumeEncryptedWith KMS
EBS SnapshotEncryptedWith KMS
KMS Key B
KMS Key B
Region ap-southeast-2
KMS ReEncrypt with KMS Key B

## Slide 659

KMS Key Policies
- Control access to KMS keys, “similar” to S3 bucket policies
- Difference: you cannot control access without them
- Default KMS Key Policy:
- Created if you don’t provide a specific KMS Key Policy
- Complete access to the key to the root user = entire AWS account
- Custom KMS Key Policy:
- Define users, roles that can access the KMS key
- Define who can administer the key
- Useful for cross-account access of your KMS key

## Slide 660

Copying Snapshots across accounts1.Create a Snapshot, encrypted with your own KMS Key (Customer Managed Key)2.Attach a KMS Key Policy to authorize cross-account access3.Share the encrypted snapshot4.(in target) Create a copy of the Snapshot, encrypt it with a CMK in your account5.Create a volume from the snapshot
KMS Key Policy

## Slide 661

KMS Multi-Region KeysAWS KMS
us-east-1
multi-Region Primary keyarn:aws:kms:us-east-1:111122223333:key/mrk-1234abcd12ab34cd56ef1234567890ab
us-west-2
multi-Region Replica keyarn:aws:kms:us-west-2:111122223333:key/mrk-1234abcd12ab34cd56ef1234567890abeu-west-1
multi-Region Replica keyarn:aws:kms:eu-west-1:111122223333:key/mrk-1234abcd12ab34cd56ef1234567890abap-southeast-2
multi-Region Replica keyarn:aws:kms:ap-southeast-2:111122223333:key/mrk-1234abcd12ab34cd56ef1234567890ab
sync

## Slide 662

KMS Multi-Region Keys
- Identical KMS keys in different AWS Regions that can be used interchangeably
- Multi-Region keys have the same key ID, key material, automatic rotation…
- Encrypt in one Region and decrypt in other Regions
- No need to re-encrypt or making cross-Region API calls
- KMS Multi-Region are NOT global (Primary + Replicas)
- Each Multi-Region key is managed independently
- Use cases: global client-side encryption, encryption on Global DynamoDB, Global Aurora

## Slide 663

DynamoDB Global Tables and KMS Multi-Region Keys Client-Side encryption
- We can encr ypt specific attributes client-side in our DynamoDB table using the Amazon DynamoDB Encryption Client
- Combined with Global Tables, the client-side encrypted data is replicated to other regions
- If we use a multi-region key, replicated in the same region as the DynamoDB Global table, then clients in these regions can use low-latency API calls to KMS in their region to decrypt the data client-side
- Using client-side encryption we can protect specific fields and guarantee only decryption if the client has access to an API key
us-east-1
DDB Table
KMS
MRK
Client App
1. Encrypt attribute with primary MRK2. Put encryptedattribute
ap-southeast-2
DDB Table
KMS
MRK
Client App5. Decrypt attribute with replica MRK
4. Get encryptedattribute
3. Global TableReplication
Attr(SSN)
Attr(SSN)
replication

## Slide 664

Global Aurora and KMS Multi-Region Keys Client-Side encryption
- We can encr ypt specific attributes client-side in our Aurora table using the AWS Encryption SDK
- Combined with Aurora Global Tables, the client-side encrypted data is replicated to other regions
- If we use a multi-region key, replicated in the same region as the Global Aurora DB, then clients in these regions can use low-latency API calls to KMS in their region to decrypt the data client-side
- Using client-side encryption we can protect specific fields and guarantee only decryption if the client has access to an API key, we can protect specific fields even from database admins
us-east-1
Table
KMS
MRK
Client App
1. Encrypt attribute with primary MRK2. Put encryptedcolumn
ap-southeast-2
Table
KMS
MRK
Client App5. Decrypt attribute with replica MRK
4. Get encryptedcolumn
3. Global DBReplication
Col(SSN)
Col(SSN)
replication

## Slide 665

S3 Replication Encryption Considerations
- Unencrypted objects and objects encrypted with SSE-S3 are replicated by default
- Objects encrypted with SSE-C (customer provided key) can be replicated
- For objects encrypted with SSE-KMS, you need to enable the option
- Specify which KMS Key to encrypt the objects within the target bucket
- Adapt the KMS Key Policy for the target key
- An IAM Role with kms:Decrypt for the source KMS Key and kms:Encrypt for the target KMS Key
- Yo u m i g h t g e t K M S t h r o t t l i n g e r r o r s , i n w h i c h c a s e yo u c a n a s k fo r a S e r v i c e Q u o t a s i n c r e a s e
- Yo u c a n u s e mu l t i-region AWS KMS Keys, but they are currently treated as independent keys by Amazon S3 (the object will still be decrypted and then encrypted)

## Slide 666

AMI Sharing Process Encrypted via KMS1.AMI in Source Account is encrypted with KMS Key from Source Account2.Must modify the image attribute to add a Launch Permission which corresponds to the specified target AWS account3.Must share the KMS Keys used to encrypted the snapshot the AMI references with the target account / IAM Role4.The IAM Role/User in the target account must have the permissions to DescribeKey, ReEncrypt*, CreateGrant, Decrypt5.When launching an EC2 instance from the AMI, optionally the target account can specify a new KMS key in its own account to re-encrypt the volumes
KMS
AMI
Account - A
Key
AMI
Account - B
share
EC2 Instancelaunch
share

## Slide 667

SSM Parameter Store
- Secure storage for configuration and secrets
- Optional Seamless Encryption using KMS
- Serverless, scalable, durable, easy SDK
- Version tracking of configurations / secrets
- Security through IAM
- Notifications with Amazon EventBridge
- Integration with CloudFormationSSM Parameter Store
Applications
Plaintext configurationEncrypted configuration
Check IAMpermissions
AWS KMS
DecryptionService

## Slide 668

SSM Parameter Store Hierarchy
- /my-department/
- my-app/
- dev/
- db-url
- db-password
- prod/
- db-url
- db-password
- other-app/
- /other-department/
- /aws/reference/secretsmanager/secret_ID_in_Secrets_Manager
- /aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2 (public)
Dev Lambda Function
GetParameters orGetParametersByPath API
Prod Lambda Function

## Slide 669

Standard and advanced parameter tiersStandardAdvancedTotal number of parameters allowed(per AWS account and Region)
10,000100,000
Maximum size of a parameter value4 KB 8 KBParameter policies availableNo YesCostNo additional chargeCharges applyStorage PricingFree $0.05 per advanced parameter per month

## Slide 670

Parameters Policies (for advanced parameters)
- Allow to assign a TTL to a parameter (expiration date) to force updating or deleting sensitive data such as passwords
- Can assign multiple policies at a timeExpiration (to delete a parameter)
ExpirationNotification (EventBridge)
NoChangeNotification (EventBridge)

## Slide 671

AWS Secrets Manager
- Newer service, meant for storing secrets
- Capability to force rotation of secrets every X days
- Automate generation of secrets on rotation (uses Lambda)
- Integration with Amazon RDS (MySQL, PostgreSQL, Aurora)
- Secrets are encrypted using KMS
- Mostly meant for RDS integration

## Slide 672

AWS Secrets Manager – Multi-Region Secrets
- Replicate Secrets across multiple AWS Regions
- Secrets Manager keeps read replicas in sync with the primary Secret
- Ability to promote a read replica Secret to a standalone Secret
- Use cases: multi-region apps, disaster recovery strategies, multi-region DB…us-east-1 (Primary)
SecretsManager
MySecret-A(primary)
us-west-2 (Secondary)
SecretsManagerMySecret-A(replica)
replicate

## Slide 673

AWS Certificate Manager (ACM)
- Easily provision, manage, and deploy TLS Certificates
- Provide in-flight encryption for websites (HTTPS)
- Supports both public and private TLS certificates
- Free of charge for public TLS certificates
- Automatic TLS certificate renewal
- Integrations with (load TLS certificates on)
- Elastic Load Balancers (CLB, ALB, NLB)
- CloudFront Distributions
- APIs on API Gateway
- Cannot use ACM with EC2 (can’t be extracted)
Auto Scaling group
EC2 InstanceEC2 Instance
AWS Certificate Manager
ApplicationLoadBalancer
HTTPS
HTTP
provision andmaintain TLS certs

## Slide 674

ACM – Requesting Public Certificates1.List domain names to be included in the certificate
- Fully Qualified Domain Name (FQDN): corp.example.com
- Wildcard Domain: *.example.com 2.Select Validation Method: DNS Validation or Email validation
- DNS Validation is preferred for automation purposes
- Email validation will send emails to contact addresses in the WHOIS database
- DNS Validation will leverage a CNAME record to DNS config (ex: Route 53)3.It will take a few hours to get verified4.The Public Certificate will be enrolled for automatic renewal
- ACM automatically renews ACM-generated certificates 60 days before expiry

## Slide 675

ACM – Importing Public Certificates
- Option to generate the certificate outside of ACM and then import it
- No automatic renewal, must impor t a new certificate before expiry
- ACM sends daily expiration events starting 45 days prior to expiration
- The # of days can be configured
- Events are appearing in EventBridge
- AWS Config has a managed rule named acm-certificate-expiration-check to check for expiring certificates (configurable number of days)
AWS Config
EventBridge
Rule events:Non-compliance
Lambda
SNS
SQS
ACMACM Events:Daily Certificate Expiry
Rule check

## Slide 676

ACM – Integration with ALB
Auto Scaling group
EC2 InstanceEC2 Instance
AWS Certificate Manager
Application Load BalancerWith HTTP -> HTTPS redirect ruleHTTP
provision andmaintain TLS certs
Redirect to HTTPSHTTPS

## Slide 677

API Gateway - Endpoint Types
- Edge-Optimized (default): For global clients
- Requests are routed through the CloudFront Edge locations (improves latency)
- The API Gateway still lives in only one region
- Regional:
- For clients within the same region
- Could manually combine with CloudFront (more control over the caching strategies and the distribution)
- Private:
- Can only be accessed from your VPC using an interface VPC endpoint (ENI)
- Use a resource policy to define access

## Slide 678

ACM – Integration with API Gateway
- Create a Custom Domain Name in API Gateway
- Edge-Optimized (default): For global clients
- Requests are routed through the CloudFront Edge locations (improves latency)
- The API Gateway still lives in only one region
- The TLS Certificate must be in the same region as CloudFront, in us-east-1
- Then setup CNAME or (better) A-Alias record in Route 53
- Regional:
- For clients within the same region
- The TLS Certificate must be imported on API Gateway, in the same region as the API Stage
- Then setup CNAME or (better) A-Alias record in Route 53
us-east-1
CloudFrontACMlinkedcertificate
API GatewayEdge-Optimized
ap-southeast-2
ACMlinkedcertificateAPI GatewayRegional

## Slide 679

CloudHSM
- KMS => AWS manages the software for encryption
- CloudHSM => AWS provisions encryption hardware
- Dedicated Hardware (HSM = Hardware Security Module)
- Yo u m a n a g e yo u r ow n e n c r y p t i o n k e y s e n t i r e ly ( n o t AW S )
- HSM device is tamper resistant, FIPS 140-2 Level 3 compliance
- Supports both symmetric and asymmetric encryption (SSL/TLS keys)
- No free tier available
- Must use the CloudHSM Client Software
- Redshift supports CloudHSM for database encryption and key management
- Good option to use with SSE-C encryption

## Slide 680

CloudHSM Diagram
AWS CloudHSMCloudHSM Client
AWS manages the Hardware
SSL ConnectionUser manages the Keys
IAM permissions:
- CRUD an HSM ClusterCloudHSM Software:
- Manage the Keys
- Manage the Users

## Slide 681

CloudHSM – High Availability
- CloudHSM clusters are spread across Multi AZ (HA)
- Great for availability and durabilityCloudHSM 1
CloudHSM Client
Availability Zone 1
CloudHSM 2
Availability Zone 2

## Slide 682

CloudHSM – Integration with AWS Services
- Through integration with AWS KMS
- Configure KMS Custom Key Store with CloudHSM
- Example: EBS, S3, RDS …
RDS DBInstanceEBS Volume
AWS KMS(Custom Key Store)
CloudHSM
KMS Encryption
Connector
CloudTrailkeys usage logs

## Slide 683

CloudHSM vs. KMSFeatureAWS KMSAWS CloudHSMTenancyMulti-TenantSingle-TenantStandardFIPS 140-2 Level 3 FIPS 140-2 Level 3Master Keys
- AWS Owned CMK
- AWS Managed CMK
- Customer Managed CMKCustomer Managed CMK
Key Types
- Symmetric
- Asymmetric
- Digital Signing
- Symmetric
- Asymmetric
- Digital Signing & HashingKey AccessibilityAccessible in multiple AWS regions (can’t access keys outside the region it’s created in)
- Deployed and managed in a VPC
- Can be shared across VPCs (VPC Peering)Cryptographic AccelerationNone
- SSL/TLS Acceleration
- Oracle TDE AccelerationAccess & AuthenticationAWS IAMYou create users and manage their permissions

## Slide 684

CloudHSM vs. KMSFeatureAWS KMSAWS CloudHSMHigh AvailabilityAWS Managed ServiceAdd multiple HSMs over different AZsAudit Capability
- CloudTrail
- CloudWatch
- CloudTrail
- CloudWatch
- MFA supportFree TierYes No

## Slide 685

AWS WAF – Web Application Firewall
- Protects your web applications from common web exploits (Layer 7)
- Layer 7 is HTTP (vs Layer 4 is TCP/UDP)
- Deploy on
- Application Load Balancer
- API Gateway
- CloudFront
- AppSync GraphQL API
- Cognito User Pool

## Slide 686

AWS WAF – Web Application Firewall
- Define Web ACL (Web Access Control List) Rules:
- IP Set: up to 10,000 IP addresses – use multiple Rules for more IPs
- HTTP headers, HTTP body, or URI strings Protects from common attack - SQL injection and Cross-Site Scripting (XSS)
- Size constraints, geo-match (block countries)
- Rate-based rules (to count occurrences of events) – for DDoS protection
- Web ACL are Regional except for CloudFront
- A rule group isa reusable set of rules that you can add to a web ACL

## Slide 687

WAF – Fixed IP while using WAF with a Load Balancer
- WAF does not support the Network Load Balancer (Layer 4)
- We can use Global Accelerator for fixed IP and WAF on the ALB
Global AcceleratorFixed IPv4: 1.2.3.4
us-east-1
Application LoadBalancer
AWS WAFWebACL
attached
Users
EC2 Instances
WebACL must be in the sameAWS Region as ALB

## Slide 688

AWS Shield: protect from DDoS attack
- DDoS: Distributed Denial of Service – many requests at the same time
- AWS Shield Standard:
- Free service that is activated for every AWS customer
- Provides protection from attacks such as SYN/UDP Floods, Reflection attacks and other layer 3/layer 4 attacks
- AWS Shield Advanced:
- Optional DDoS mitigation service ($3,000 per month per organization)
- Protect against more sophisticated attack on Amazon EC2, Elastic Load Balancing (ELB), Amazon CloudFront, AWS Global Accelerator, and Route 53
- 24/7 access to AWS DDoS response team (DRP)
- Protect against higher fees during usage spikes due to DDoS
- Shield Advanced automatic application layer DDoS mitigation automatically creates, evaluates and deploys AWS WAF rules to mitigate layer 7 attacks

## Slide 689

AWS Firewall Manager
- Manage rules in all accounts of an AWS Organization
- Security policy: common set of security rules
- WAF rules (Application Load Balancer, API Gateways, CloudFront)
- AWS Shield Advanced (ALB, CLB, NLB, Elastic IP, CloudFront)
- Security Groups for EC2, Application Load BAlancer and ENI resources in VPC
- AWS Network Firewall (VPC Level)
- Amazon Route 53 Resolver DNS Firewall
- Policies are created at the region level
- Rules are applied to new resources as they are created (good for compliance) across all and future accounts in your Organization

## Slide 690

WAF vs. Firewall Manager vs. Shield
- WAF, Shield and Firewall Manager are used together for comprehensive protection
- Define your Web ACL rules in WAF
- For granular protection of your resources, WAF alone is the correct choice
- If you want to use AWS WAF across accounts, accelerate WAF configuration, automate the protection of new resources, use Firewall Manager with AWS WAF
- Shield Advanced adds additional features on top of AWS WAF, such as dedicated support from the Shield Response Team (SRT) and advanced reporting.
- If you’re prone to frequent DDoS attacks, consider purchasing Shield Advanced
AWS Firewall ManagerAWS WAF
AWS Shield

## Slide 691

AWS Best Practices for DDoS ResiliencyEdge Location Mitigation (BP1, BP3)
- BP1 – CloudFront
- Web Application deliver y at the edge
- Protect from DDoS Common Attacks (SYN floods, UDP reflection…)
- BP1 – Global Accelerator
- Access your application from the edge
- Integration with Shield for DDoS protection
- Helpful if your backend is not compatible with CloudFront
- BP3 – Route 53
- Domain Name Resolution at the edge
- DDoS Protection mechanism

## Slide 692

AWS Best Practices for DDoS ResiliencyBest pratices for DDoS mitigation
- Infrastructure layer defense (BP1, BP3, BP6)
- Protect Amazon EC2 against high traffic
- That includes using Global Accelerator, Route 53, CloudFront, Elastic Load Balancing
- Amazon EC2 with Auto Scaling (BP7)
- Helps scale in case of sudden traffic surges including a flash crowd or a DDoS attack
- Elastic Load Balancing (BP6)
- Elastic Load Balancing scales with the traffic increases and will distribute the traffic to many EC2 instances

## Slide 693

AWS Best Practices for DDoS ResiliencyApplication Layer Defense
- Detect and filter malicious web requests (BP1, BP2)
- CloudFront cache static content and serve it from edge locations, protecting your backend
- AWS WAF is used on top of CloudFront and Application Load Balancer to filter and block requests based on request signatures
- WAF rate-based rules can automatically block the IPs of bad actors
- Use managed rules on WAF to block attacks based on IP reputation, or block anonymous Ips
- CloudFront can block specific geographies
- Shield Advanced (BP1, BP2, BP6)
- Shield Advanced automatic application layer DDoS mitigation automatically creates, evaluates and deploys AWS WAF rules to mitigate layer 7 attacks

## Slide 694

AWS Best Practices for DDoS ResiliencyAttack surface reduction
- Obfuscating AWS resources (BP1, BP4, BP6)
- Using CloudFront, API Gateway, Elastic Load Balancing to hide your backend resources (Lambda functions, EC2 instances)
- Security groups and Network ACLs (BP5)
- Use security groups and NACLs to filter traffic based on specific IP at the subnet or ENI-level
- Elastic IP are protected by AWS Shield Advanced
- Protecting API endpoints (BP4)
- Hide EC2, Lambda, elsewhere
- Edge-optimized mode, or CloudFront + regional mode (more control for DDoS)
- WAF + API Gateway: burst limits, headers filtering, use API keys

## Slide 695

Amazon GuardDuty
- Intelligent Threat discovery to protect your AWS Account
- Uses Machine Learning algorithms, anomaly detection, 3rd party data
- One click to enable (30 days trial), no need to install software
- Input data includes:
- CloudTrail Events Logs – unusual API calls, unauthorized deployments
- CloudTrail Management Events – create VPC subnet, create trail, …
- CloudTrail S3 Data Events – get object, list objects, delete object, …
- VPC Flow Logs – unusual internal traffic, unusual IP address
- DNS Logs – compromised EC2 instances sending encoded data within DNS queries
- Optional Features – EKS Audit Logs, RDS & Aurora, EBS, Lambda, S3 Data Events…
- Can setup EventBridge rules to be notified in case of findings
- EventBridge rules can target AWS Lambda or SNS
- Can protect against CryptoCurrency attacks (has a dedicated “finding” for it)

## Slide 696

Amazon GuardDuty
VPC Flow Logs
CloudTrail Logs
DNS Logs (AWS DNS)GuardDuty
EventBridge
SNS
Lambda
Optional Features
EKS Audit Logs &Runtime Monitoring
RDS & AuroraLogin Activity
S3 Logs
EBS Volumes
Lambda NetworkActivity

## Slide 697

Amazon Inspector
- Automated Security Assessments
- For EC2 instances
- Leveraging the AWS System Manager (SSM) agent
- Analyze against unintended network accessibility
- Analyze the running OS against known vulnerabilities
- For Container Images push to Amazon ECR
- Assessment of Container Images as they are pushed
- For Lambda Functions
- Identifies software vulnerabilities in function code and package dependencies
- Assessment of functions as they are deployed
- Reporting & integration with AWS Security Hub
- Send findings to Amazon Event Bridge
SSM Agent
AmazonInspector
EventBridge
assessment run state& findings
Security Hub
LambdaFunction
Amazon ECRContainer Image

## Slide 698

What does Amazon Inspector evaluate?
- Remember: only for EC2 instances, Container Images & Lambda functions
- Continuous scanning of the infrastructure, only when needed
- Package vulnerabilities (EC2, ECR & Lambda) – database of CVE
- Network reachability (EC2)
- A risk score is associated with all vulnerabilities for prioritization

## Slide 699

AWS Macie
- Amazon Macie is a fully managed data security and data privacy service that uses machine learning and pattern matching to discover and protect your sensitive data in AWS.
- Macie helps identify and alert you to sensitive data, such as personally identifiable information (PII)
S3 Buckets
MacieDiscover Sensitive Data (PII)
AmazonEventBridge
analyzenotifyintegrations

## Slide 700

Amazon VPC

## Slide 701

VPC Components Diagram
Region
Availability Zone
VPC
VPC PeeringConnections
Private SubnetPublic Subnet
InternetGatewayRouter
RouteTable
NAT Gateway
Security Group
Public EC2 Instance
NACL
NACL
Security Group
Private EC2 InstanceSecurity Group
Private EC2 Instance
RouteTable
VPC Flow Logs
AmazonDynamoDB
Internetwww Corporate Data Center
Server
VPNGateway
CustomerGateway
S2S VPNConnection
CloudWatch
S3VPCEndpoint
DX Location
Direct ConnectConnection
TransitGatewayVPN
DX

## Slide 702

Understanding CIDR – IPv4
- Classless Inter-Domain Routing – a method for allocating IP addresses
- Used in Security Groups rules and AWS networking in general
- They help to define an IP address range:
- We’ve seen WW.XX.YY.ZZ/32 => one IP
- We’ve seen 0.0.0.0/0 => all IPs
- But we can define:192.168.0.0/26 =>192.168.0.0 – 192.168.0.63 (64 IP addresses)

## Slide 703

Understanding CIDR – IPv4
- A CIDR consists of two components
- Base IP
- Represents an IP contained in the range (XX.XX.XX.XX)
- Example: 10.0.0.0, 192.168.0.0, …
- Subnet Mask
- Defines how many bits can change in the IP
- Example: /0, /24, /32
- Can take two forms:
- /8 ó 255.0.0.0
- /16 ó 255.255.0.0
- /24 ó 255.255.255.0
- /32 ó 255.255.255.255

## Slide 704

Understanding CIDR – Subnet Mask
- The Subnet Mask basically allows part of the underlying IP to get additional next values from the base IP192.168.0.0/32 => allows for 1 IP (2!)192.168.0.0192.168.0.0/31 => allows for 2 IP (2")192.168.0.0 -> 192.168.0.1192.168.0.0/30 => allows for 4 IP (2#)192.168.0.0 -> 192.168.0.3192.168.0.0/29 => allows for 8 IP (2$)192.168.0.0 -> 192.168.0.7192.168.0.0/28 => allows for 16 IP (2%)192.168.0.0 -> 192.168.0.15192.168.0.0/27 => allows for 32 IP (2&)192.168.0.0 -> 192.168.0.31192.168.0.0/26 => allows for 64 IP (2')192.168.0.0 -> 192.168.0.63192.168.0.0/25 => allows for 128 IP (2()192.168.0.0 -> 192.168.0.127192.168.0.0/24 => allows for 256 IP (2))192.168.0.0 -> 192.168.0.255…192.168.0.0/16 => allows for 65,536 IP (2"')192.168.0.0 -> 192.168.255.255…0.0.0.0/0 => allows for All IPs0.0.0.0 -> 255.255.255.255
Quick Memo
- /32 – no octet can change
- /24 – last octet can change
- /16 – last 2 octets can change
- /8 – last 3 octets can change
- /0 – all octets can change
...1*+2,-3.-4+/Octets

## Slide 705

Understanding CIDR – Little Exercise
- 192.168.0.0/24 = … ?
- 192.168.0.0 – 192.168.0.255 (256 IPs)
- 192.168.0.0/16 = … ?
- 192.168.0.0 – 192.168.255.255 (65,536 IPs)
- 134.56.78.123/32 = … ?
- Just 134.56.78.123
- 0.0.0.0/0
- All IPs!
- When in doubt, use this website https://www.ipaddressguide.com/cidr

## Slide 706

Public vs. Private IP (IPv4)
- The Internet Assigned Numbers Authority (IANA) established certain blocks of IPv4 addresses for the use of private (LAN) and public (Internet) addresses
- Private IP can only allow certain values:
- 10.0.0.0 – 10.255.255.255 (10.0.0.0/8) ç in big networks
- 172.16.0.0 – 172.31.255.255 (172.16.0.0/12) ç AWS default VPC in that range
- 192.168.0.0 – 192.168.255.255 (192.168.0.0/16) ç e.g., home networks
- All the rest of the IP addresses on the Internet are Public

## Slide 707

Default VPC Walkthrough
- All new AWS accounts have a default VPC
- New EC2 instances are launched into the default VPC if no subnet is specified
- Default VPC has Internet connectivity and all EC2 instances inside it have public IPv4 addresses
- We also get a public and a private IPv4 DNS names

## Slide 708

VPC in AWS – IPv4
- VPC = Virtual Private Cloud
- Y ou can have multiple VPCs in an AWS region (max. 5 per region – soft limit)
- Max. CIDR per VPC is 5, for each CIDR:
- Min. size is /28 (16 IP addresses)
- Max. size is /16 (65536 IP addresses)
- Because VPC is private, only the Private IPv4 ranges are allowed:
- 10.0.0.0 – 10.255.255.255 (10.0.0.0/8)
- 172.16.0.0 – 172.31.255.255 (172.16.0.0/12)
- 192.168.0.0 – 192.168.255.255 (192.168.0.0/16)
- Yo u r V P C C I D R s h o u l d N OT ove r l a p w i t h yo u r o t h e r n e t wo r k s ( e . g . , c o r p o r a t e )

## Slide 709

State of Hands-on
Region
VPC

## Slide 710

Adding Subnets
Region
VPC
Availability Zone
Private SubnetPublic Subnet

## Slide 711

VPC – Subnet (IPv4)
- AWS reserves 5 IP addresses (first 4 & last 1) in each subnet
- These 5 IP addresses are not available for use and can’t be assigned to an EC2 instance
- Example: if CIDR block 10.0.0.0/24, then reserved IP addresses are:
- 10.0.0.0 – Network Address
- 10.0.0.1 – reserved by AWS for the VPC router
- 10.0.0.2 – reserved by AWS for mapping to Amazon-provided DNS
- 10.0.0.3 – reserved by AWS for future use
- 10.0.0.255 – Network Broadcast Address. AWS does not support broadcast in a VPC, therefore the address is reserved
- Exam Tip, if you need 29 IP addresses for EC2 instances:
- Y ou can’t choose a subnet of size /27 (32 IP addresses, 32 – 5 = 27 < 29)
- Y ou need to choose a subnet of size /26 (64 IP addresses, 64 – 5 = 59 > 29)

## Slide 712

Internet Gateway (IGW)
- Allows resources (e.g., EC2 instances) in a VPC connect to the Internet
- It scales horizontally and is highly available and redundant
- Must be created separately from a VPC
- One VPC can only be attached to one IGW and vice versa
- Internet Gateways on their own do not allow Internet access…
- Route tables must also be edited!

## Slide 713

Adding Internet Gateway
Region
VPC
Availability Zone
Private SubnetPublic Subnet
InternetGateway

## Slide 714

Editing Route Tables
Region
VPC
Availability Zone
Private SubnetPublic Subnet
InternetGateway
Router
RouteTableSecurity Group
Public EC2 Instance
Internetwww

## Slide 715

Bastion Hosts
- We can use a Bastion Host to SSH into our private EC2 instances
- The bastion is in the public subnet which is then connected to all other private subnets
- Bastion Host security group must allow inbound from the internet on port 22 from restricted CIDR, for example the public CIDR of your corporation
- Security Group of the EC2 Instances must allow the Security Group of the Bastion Host, or the private IP of the Bastion host
VPC
Private Subnet
Public Subnet
EC2 Instance(Bastion Host)
Security Group (BastionHost-SG)
Security Group (LinuxInstance-SG)
Users
SSH
SSH

## Slide 716

NAT Instance (outdated, but still at the exam)
- NAT = Network Address Translation
- Allows EC2 instances in private subnets to connect to the Internet
- Must be launched in a public subnet
- Must disable EC2 setting: Source / destination Check
- Must have Elastic IP attached to it
- Route Tables must be configured to route traffic from private subnets to the NAT Instance
VPC
Private Subnet
Public Subnet
NAT Instance
Security Group (NATInstance-SG)
EIP(IP: 12.34.56.78)
IP: 10.0.0.10IP: 10.0.0.20
Server(IP: 50.60.4.10)
Src.: 10.0.0.20Dest.: 50.60.4.10
Src.: 12.34.56.78Dest.: 50.60.4.10
Src.: 50.60.4.10Dest.: 12.34.56.78
Src.: 50.60.4.10Dest.: 10.0.0.20

## Slide 717

NAT Instance
Region
VPC
Availability Zone
Private SubnetPublic Subnet
InternetGateway
Router
RouteTableSecurity Group
Public EC2 Instance
Internetwww Security Group
NAT Instance
Security Group
Private EC2 Instance
RouteTable
EIP

## Slide 718

NAT Instance – Comments
- Pre-configured Amazon Linux AMI is available
- Reached the end of standard support on December 31, 2020
- Not highly available / resilient setup out of the box
- Yo u n e e d t o c r e a t e a n A S G i n mu l t i-AZ + resilient user-data script
- Internet traffic bandwidth depends on EC2 instance type
- Yo u mu s t m a n a g e S e c u r i t y G r o u p s & r u l e s :
- Inbound:
- Allow HTTP / HTTPS traffic coming from Private Subnets
- Allow SSH from your home network (access is provided through Internet Gateway)
- Outbound:
- Allow HTTP / HTTPS traffic to the Internet

## Slide 719

NAT Gateway
- AWS-managed NAT, higher bandwidth, high availability, no administration
- Pay per hour for usage and bandwidth
- NATGW is created in a specific Availability Zone, uses an Elastic IP
- Can’t be used by EC2 instance in the same subnet (only from other subnets)
- Requires an IGW (Private Subnet => NATGW => IGW)
- 5 Gbps of bandwidth with automatic scaling up to 100 Gbps
- No Security Groups to manage / required

## Slide 720

NAT Gateway
Region
VPC
Availability Zone
Private SubnetPublic Subnet
InternetGateway
Router
RouteTableSecurity Group
Public EC2 Instance
Internetwww
NAT GatewaySecurity Group
Private EC2 Instance
RouteTable

## Slide 721

NAT Gateway with High Availability
- NAT Gateway is resilient within a single Availability Zone
- Must create multiple NAT Gateways in multiple AZs for fault-tolerance
- There is no cross-AZ failover needed because if an AZ goes down it doesn't need NAT
Region
VPC
AZ - BAZ - A
Public Subnet
Private Subnet
NAT Gateway
EC2 Instance
Router
InternetGateway
Internetwww
Public Subnet
Private Subnet
NAT Gateway
EC2 Instance

## Slide 722

NAT Gateway vs. NAT InstanceNAT GatewayNAT InstanceAvailabilityHighly available within AZ (create in another AZ)Use a script to manage failover between instancesBandwidthUp to 100 Gbps Depends on EC2 instance typeMaintenanceManaged by AWSManaged by you (e.g., software, OS patches, …)CostPer hour & amount of data transferred Per hour, EC2 instance type and size, + network $Public IPv4Private IPv4Security GroupsUse as Bastion Host?
More at: https://docs.aws.amazon.com/vpc/latest/userguide/vpc-nat-comparison.html

## Slide 723

Subnet
Security Groups & NACLs
Security Group
EC2 Instance
NACL
NACL InboundRules
SG InboundRules
Outbound Allowed(Stateful)
NACL OutboundRules (Stateless)
Incoming RequestSubnetSecurity Group
EC2 Instance
NACL
NACL OutboundRules
SG OutboundRules
Inbound Allowed(Stateful)
NACL InboundRules (Stateless)
Outgoing Request
1 2
3
12
3

## Slide 724

Network Access Control List (NACL)
- NACL are like a firewall which control traffic from and to subnets
- One NACL per subnet, new subnets are assigned the Default NACL
- Yo u d e fi n e NACL Rules:
- Rules have a number (1-32766), higher precedence with a lower number
- First rule match will drive the decision
- Example: if you define #100 ALLOW 10.0.0.10/32 and #200 DENY 10.0.0.10/32, the IP address will be allowed because 100 has a higher precedence over 200
- The last rule is an asterisk (*) and denies a request in case of no rule match
- AWS recommends adding rules by increment of 100
- Newly created NACLs will deny everything
- NACL are a great way of blocking a specific IP address at the subnet level

## Slide 725

NACLs
Region
VPC
Availability Zone
Private SubnetPublic Subnet
InternetGateway
Router
RouteTableSecurity Group
Public EC2 Instance
Internetwww
NAT GatewaySecurity Group
Private EC2 Instance
RouteTable
NACL
NACL

## Slide 726

Default NACL
- Accepts everything inbound/outbound with the subnets it’s associated with
- Do NOT modify the Default NACL, instead create custom NACLs
Rule #TypeProtocolPort RangeSourceAllow/Deny100All IPv4 TrafficAllAll0.0.0.0/0ALLOW*All IPv4 TrafficAllAll0.0.0.0/0DENY
Inbound Rules
Rule #TypeProtocolPort RangeDestinationAllow/Deny100All IPv4 TrafficAllAll0.0.0.0/0ALLOW*All IPv4 TrafficAllAll0.0.0.0/0DENY
Outbound Rules
Default NACL for a VPC that supports IPv4

## Slide 727

Response
Ephemeral Ports
- For any two endpoints to establish a connection, they must use ports
- Clients connect to a defined port, and expect a response on an ephemeral port
- Different Operating Systems use different port ranges, examples:
- IANA & MS Windows 10 è 49152 – 65535
- Many Linux Kernels è 32768 – 60999
Web ServerIP: 55.66.77.88Fixed Port: 443
Dest. IP 55.66.77.88Dest. Port 443Src. IP11.22.33.44
ClientIP: 11.22.33.44Ephemeral Port: 50105
Src. Port50105 Payload …
Src. IP55.66.77.88Src. Port443Dest. IP 11.22.33.44Dest. Port 50105Payload …
Request

## Slide 728

NACL with Ephemeral PortsVPC
DB Subnet (Private)Web Subnet (Public)
Web Tier Database Tier
DB InstancePort 3306
Web-NACL
DB-NACL
Allow Outbound TCPOn port 3306To DB Subnet CIDR
Allow Inbound TCPOn port 3306From Web Subnet CIDR
Allow Inbound TCPOn port 1024-65535From DB Subnet CIDRAllow Outbound TCPOn port 1024-65535To Web Subnet CIDR
https://docs.aws.amazon.com/vpc/latest/userguide/vpc-network-acls.html#nacl-ephemeral-ports
ClientEphemeral Port

## Slide 729

Create NACL rules for each target subnets CIDRVPC
DB Subnet – A (Private)Web Subnet - A (Public)
Web Subnet - B (Public)
DB Subnet – B (Private)
Web Tier Database Tier
DB Instance
DB Instance
Web-NACL
DB-NACL

## Slide 730

Security Group vs. NACLsSecurity GroupNACLOperates at the instance levelOperates at the subnet levelSupports allow rules onlySupports allow rules and deny rulesStateful: return traffic is automatically allowed, regardless of any rulesStateless: return traffic must be explicitly allowed by rules (think of ephemeral ports)All rules are evaluated before deciding whether to allow trafficRules are evaluated in order (lowest to highest) when deciding whether to allow traffic, first match winsApplies to an EC2 instance when specified by someone Automatically applies to all EC2 instances in the subnet that it’s associated withNACL Examples: https://docs.aws.amazon.com/vpc/latest/userguide/vpc-network-acls.html

## Slide 731

VPC Peering
- Privately connect two VPCs using AWS’ network
- Make them behave as if they were in the same network
- Must not have overlapping CIDRs
- VPC Peering connection is NOT transitive (must be established for each VPC that need to communicate with one another)
- Yo u mu s t u p d a t e r o u t e t a bl e s i neach VPC’s subnets to ensure EC2 instances can communicate with each other
VPC - A
VPC - B
VPC - C
VPC Peering(A – B)
VPC Peering(B – C)
VPC Peering(A – C)

## Slide 732

VPC Peering – Good to know
- Yo u c a n c r e a t e V P C Pe e r i n g c o n n e c t i o n b e t we e n V P C s i n different AWS accounts/regions
- Yo u c a n r e fe r e n c e a s e c u r i t y g r o u p i n a p e e r e d V P C ( wo r k s c r o s s accounts – same region)
Account ID

## Slide 733

VPC Peering
Region
VPC
Availability Zone
Private SubnetPublic Subnet
InternetGateway
Router
RouteTableSecurity Group
Public EC2 Instance
Internetwww
NAT GatewaySecurity Group
Private EC2 Instance
RouteTable
NACL
NACL
VPC PeeringConnections

## Slide 734

VPC Endpoints
Region
VPC
Availability Zone
Private SubnetPublic Subnet
InternetGateway
Router
RouteTableSecurity Group
Public EC2 Instance
Internetwww
NAT GatewaySecurity Group
Private EC2 Instance
RouteTable
NACL
NACL
VPC PeeringConnections
CloudWatch
S3VPCEndpoint
AmazonDynamoDB

## Slide 735

VPC Endpoints (AWS PrivateLink)
- Every AWS service is publicly exposed (public URL)
- VPC Endpoints (powered by AWS PrivateLink) allows you to connect to AWS services using a private network instead of using the public Internet
- They’re redundant and scale horizontally
- They remove the need of IGW, NATGW, … to access AWS Services
- In case of issues:
- Check DNS Setting Resolution in your VPC
- Check Route Tables
VPC
Private Subnet
Public Subnet
InternetGatewayRegion
VPC Endpoint
EC2 Instance
www
NATGateway
EC2 Instance
Amazon SNS
Option 1Option 2
Amazon SNS

## Slide 736

Private Subnet
Types of Endpoints
- Interface Endpoints (powered by PrivateLink)
- Provisions an ENI (private IP address) as an entry point (must attach a Security Group)
- Supports most AWS services
- $ per hour + $ per GB of data processed
- Gateway Endpoints
- Provisions a gateway and must be used as a target in a route table (does not use security groups)
- Supports both S3 and DynamoDB
- Free
Region
VPC
EC2 Instance
VPC Endpoint(Interface)ENI (PrivateLink)
Amazon SNS
Private Subnet
Region
VPC
EC2 Instance
Amazon S3
VPC Endpoint(Gateway)
AmazonDynamoDBOR

## Slide 737

Gateway or Interface Endpoint for S3?
- Gateway is most likely going to be preferred all the time at the exam
- Cost: free for Gateway, $ for interface endpoint
- Interface Endpoint is preferred access is required from on-premises (Site to Site VPN or Direct Connect), a different VPC or a different region
Region
AWS Cloud
VPC
Direct ConnectS2S VPN
Amazon S3
InterfaceEndpointIn-VPCApps
GatewayEndpoint
PrivateLink
Users

## Slide 738

Private subnet
Lambda in VPC accessing DynamoDBAWS Cloud
DynamoDB
Public subnet
NATIGW
VPC Gateway EndpointFor DynamoDB
- DynamoDB is a public service from AWS
- Option 1: Access from the public internet
- Because Lambda is in a VPC, it needs a NAT Gateway in a public subnet and an internet gateway
- Option 2 (better & free): Access from the private VPC network
- Deploy a VPC Gateway endpoint for DynamoDB
- Change the Route Tables

## Slide 739

VPC Flow Logs
- Capture information about IP traffic going into your interfaces:
- VPC Flow Logs
- Subnet Flow Logs
- Elastic Network Interface (ENI) Flow Logs
- Helps to monitor & troubleshoot connectivity issues
- Flow logs data can go to S3, CloudWatch Logs, and Kinesis Data Firehose
- Captures network information from AWS managed interfaces too: ELB, RDS, ElastiCache, Redshift, WorkSpaces, NATGW, Transit Gateway…

## Slide 740

VPC Flow Logs
Region
VPC
Availability Zone
Private SubnetPublic Subnet
InternetGateway
Router
RouteTableSecurity Group
Public EC2 Instance
Internetwww
NAT GatewaySecurity Group
Private EC2 Instance
RouteTable
NACL
NACL
VPC PeeringConnections
AmazonDynamoDB
CloudWatch
S3VPCEndpoint
VPC Flow Logs

## Slide 741

VPC Flow Logs Syntax
- srcaddr & dstaddr – help identify problematic IP
- srcport & dstport – help identity problematic ports
- Action – success or failure of the request due to Security Group / NACL
- Can be used for analytics on usage patterns, or malicious behavior
- Query VPC flow logs using Athena on S3 or CloudWatch Logs Insights
- Flow Logs examples: https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs-records-examples.html
version
account-id
interface-id
srcaddr
dstaddr
srcport
dstport
protocol
packets
bytes
start
end
action
log-status

## Slide 742

VPC Flow Logs – Troubleshoot SG & NACL issues
SubnetSecurity Group
EC2 Instance
NACL
NACL InboundRules
SG InboundRules
Outbound Allowed(Stateful)
NACL OutboundRules (Stateless)
SubnetSecurity Group
EC2 Instance
NACL
NACL OutboundRules
SG OutboundRules
Inbound Allowed(Stateful)
NACL InboundRules (Stateless)
Incoming Requests
- Inbound REJECT => NACL or SG
- Inbound ACCEPT, Outbound REJECT => NACL
Outgoing Requests
- Outbound REJECT => NACL or SG
- Outbound ACCEPT, Inbound REJECT => NACL
Look at the “ACTION” field

## Slide 743

VPC Flow Logs – Architectures
VPC Flow Logs
CloudWatch LogsCloudWatch Contributor Insights
Top-10 IP addresses
CW Alarm
Alert
Amazon SNS
Metric Filter
VPC Flow Logs
CloudWatch Logs
VPC Flow Logs
S3 Bucket
AmazonAthena
AmazonQuickSight
SSH, RDP…

## Slide 744

VPC Flow Logs – CloudWatch Permissions
- IAM Service Role associated with VPC Flow Logs must have the required permissions to publish logs to CloudWatch Logs
- logs:CreateLogGroup, logs:CreateLogStream, or logs:PutLogEvents
VPC Flow Logs
IAM Service Role
CloudWatch Logslogs

## Slide 745

AWS Site-to-Site VPN
Region
Availability Zone
VPC
VPC PeeringConnections
Private SubnetPublic Subnet
InternetGatewayRouter
RouteTable
NAT Gateway
Security Group
Public EC2 Instance
NACL
NACL
Security Group
Private EC2 InstanceSecurity Group
Private EC2 Instance
RouteTable
VPC Flow Logs
AmazonDynamoDB
Internetwww Corporate Data Center
Server
VPNGateway
CustomerGateway
S2S VPNConnection
CloudWatch
S3VPCEndpoints

## Slide 746

AWS Site-to-Site VPN
- Virtual Private Gateway (VGW)
- VPN concentrator on the AWS side of the VPN connection
- VGW is created and attached to the VPC from which you want to create the Site-to-Site VPN connection
- Possibility to customize the ASN (Autonomous System Number)
- Customer Gateway (CGW)
- Software application or physical device on customer side of the VPN connection
- https://docs.aws.amazon.com/vpn/latest/s2svpn/your-cgw.html#DevicesTested

## Slide 747

Private SubnetSite-to-Site VPN Connections
- Customer Gateway Device (On-premises)
- What IP address to use?
- Public Internet-routable IP address for your Customer Gateway device
- If it’s behind a NAT device that’s enabled for NAT traversal (NAT-T), use the public IP address of the NAT device
- Important step: enable Route Propagation for the Virtual Private Gateway in the route table that is associated with your subnets
- If you need to ping your EC2 instances from on-premises, make sure you add the ICMP protocol on the inbound of your security groups
VPC
Security Group
Corporate Data Center
Server
Virtual PrivateGateway
CustomerGateway(Public IP)
Route Table(Route Propagation enabled)
NAT Device(Public IP)
CustomerGateway(Private IP)
OR

## Slide 748

AWS VPN CloudHub
- Provide secure communication between multiple sites, if you have multiple VPN connections
- Low-cost hub-and-spoke model for primary or secondary network connectivity between different locations (VPN only)
- It’s a VPN connection so it goes over the public Internet
- To s e t i t u p, c o n n e c t mu l t i p l e V P N connections on the same VGW, setup dynamic routing and configure route tables
VPC
Private Subnet 1
Availability Zone
Private Subnet 2
Availability Zone
EC2 Instances
EC2 Instances
Customer Network
Customer Network
Customer Network
VirtualPrivateGateway(VGW)
CustomerGateway
CustomerGateway
CustomerGateway

## Slide 749

Direct Connect (DX)
- Provides a dedicated private connection from a remote network to your VPC
- Dedicated connection must be setup between your DC and AWS Direct Connect locations
- Yo u n e e d t o s e t u p a V i r t u a l P r i v a t e G a t e w ay o n yo u r V P C
- Access public resources (S3) and private (EC2) on same connection
- Use Cases:
- Increase bandwidth throughput - working with large data sets – lower cost
- More consistent network experience - applications using real-time data feeds
- Hybrid Environments (on prem + cloud)
- Supports both IPv4 and IPv6

## Slide 750

Direct Connect Diagram
Region(us-east-1)
VPC
Private Subnet
EC2 Instances
Virtual Private Gateway
Amazon S3Amazon Glacier
AWS Direct Connect LocationAWS CageCustomer orpartner cage
AWS DirectConnect EndpointCustomer orpartner router
Corporatedata center
Customer Network
Customerrouter/firewall
VLAN 1VLAN 2
Private virtual interfacePublic virtual interface

## Slide 751

Direct Connect Gateway
- If you want to setup a Direct Connect to one or more VPC in many different regions (same account), you must use a Direct Connect Gateway
Region(us-east-1)
VPC
Region(us-west-1)
VPC
10.0.0.0/16172.16.0.0/16
Direct Connect Gateway
Customer network
AWS DirectConnectconnection
Private virtualinterfacePrivate virtualinterface
Private virtualinterface

## Slide 752

Direct Connect – Connection Types
- Dedicated Connections: 1Gbps,10 Gbps and 100 Gbps capacity
- Physical ethernet port dedicated to a customer
- Request made to AWS first, then completed by AWS Direct Connect Partners
- Hosted Connections: 50Mbps, 500 Mbps, to 10 Gbps
- Connection requests are made via AWS Direct Connect Partners
- Capacity can be added or removed on demand
- 1, 2, 5, 10 Gbps available at select AWS Direct Connect Partners
- Lead times are often longer than 1 month to establish a new connection

## Slide 753

Direct Connect – Encryption
- Data in transit is not encrypted but is private
- AWS Direct Connect + VPN provides an IPsec-encrypted private connection
- Good for an extra level of security, but slightly more complex to put in place
Region(us-east-1)
VPCAvailability Zone(us-east-1a)
Availability Zone(us-east-1b)
Private Subnet 1
Private Subnet 2
EC2 Instances
EC2 Instances
AWS DirectConnect Location
AWS DirectConnect Endpoint
Customer Network
Customerrouter/firewall
Corporatedata center
VPN Connection
Client
Client

## Slide 754

Direct Connect - ResiliencyHigh Resiliency for Critical Workloads
One connection at multiple locations
Maximum Resiliency for Critical Workloads
Maximum resilience is achieved by separate connections terminating on separate devices in more than one location.
Region
AWS DirectConnect Location - 1
Corporatedata center
AWS DirectConnect Location - 2
Corporatedata center
Region
AWS DirectConnect Location - 1
Corporatedata center
AWS DirectConnect Location - 2
Corporatedata center

## Slide 755

Site-to-Site VPN connection as a backup
- In case Direct Connect fails, you can set up a backup Direct Connect connection (expensive), or a Site-to-Site VPN connection
Direct ConnectPrimary Connection
Site-to-Site VPNBackup Connection
AWS Cloud
Corporate DC
VPC

## Slide 756

Network topologies can become complicated
Customer Gateway
Amazon VPC
VPN Connection
Amazon VPC
VPC PeeringConnection
Amazon VPC
VPC PeeringConnection
Amazon VPC
VPC PeeringConnection
VPC PeeringConnection
VPC PeeringConnection
VPN Connection
VPN Connection
Direct ConnectGateway

## Slide 757

Transit Gateway
- For having transitive peering between thousands of VPC and on-premises, hub-and-spoke (star) connection
- Regional resource, can work cross-region
- Share cross-account using Resource Access Manager (RAM)
- Yo u c a n p e e r Tr a n s i t G a t e w ay s a c r o s s r e g i o n s
- Route Tables: limit which VPC can talk with other VPC
- Works with Direct Connect Gateway, VPN connections
- Supports IP Multicast (not supported by any other AWS service)
AWS DirectConnect Gateway
Amazon VPC
Amazon VPC
Amazon VPC
Amazon VPCVPN Connection
Customer Gateway
Transit Gateway

## Slide 758

Transit Gateway: Site-to-Site VPN ECMP
- ECMP = Equal-cost multi-path routing
- Routing strategy to allow to forward a packet over multiple best path
- Use case: create multiple Site-to-Site VPN connections to increase the bandwidth of your connection to AWS
VPC
VPC
VPC
VPC
AWS Transit Gateway
VPC attachment
VPC attachment
VPC attachment
VPC attachment
Corporate data center172.16.0.0/16
VPN attachment

## Slide 759

Transit Gateway: throughput with ECMPVPN to virtual private gateway
1x=1xVPC
1x=1.25 Gbps
VPN connection(2 tunnels)
VPN to transit gateway
1x=1xVPC
VPC
VPC
VPC
1x=2.5 Gbps (ECMP) – 2 tunnels used
2x=5.0 Gbps (ECMP)
3x=7.5 Gbps (ECMP)
per GB of TGWprocessed data

## Slide 760

Transit Gateway – Share Direct Connect between multiple accounts
AWS DirectConnect Location
Corporatedata center
Clients
Clients
Servers
AWS DirectConnect endpointCustomer router/firewall
Transit VIFVLAN
Account 1
DirectConnectGateway
AWS CloudRegion
TransitGatewayAccount 2
VPC
VPC
You can use AWS Resource Access Manager to share Transit Gateway with other accounts.

## Slide 761

VPC – Traffic Mirroring
- Allows you to capture and inspect network traffic in your VPC
- Route the traffic to security appliances that you manage
- Capture the traffic
- From (Source) – ENIs
- To ( Ta r g e t s ) – an ENI or a Network Load Balancer
- Capture all packets or capture the packets of your interest (optionally, truncate packets)
- Source and T arget can be in the same VPC or different VPCs (VPC Peering)
- Use cases: content inspection, threat monitoring, troubleshooting, …
Source B
Source A
Auto Scaling group
Network LoadBalancer
Inbound &Outbound trafficInbound &Outbound traffic
Traffic Mirroring(filter traffic, optional)
EC2 instances with Security Appliances

## Slide 762

What is IPv6?
- IPv4 designed to provide 4.3 Billion addresses (they’ll be exhausted soon)
- IPv6 is the successor of IPv4
- IPv6 is designed to provide 3.4 × 10,- unique IP addresses
- Every IPv6 address in AWS is public and Internet-routable (no private range)
- Format è x.x.x.x.x.x.x.x (x is hexadecimal, range can be from 0000 to ffff)
- Examples:
- 2001:db8:3333:4444:5555:6666:7777:8888
- 2001:db8:3333:4444:cccc:dddd:eeee:ffff
- :: è all 8 segments are zero
- 2001:db8:: è the last 6 segments are zero
- ::1234:5678 è the first 6 segments are zero
- 2001:db8::1234:5678 è the middle 4 segments are zero

## Slide 763

IPv6 in VPC
- IPv4 cannot be disabled for your VPC and subnets
- Y ou can enable IPv6 (they’re public IP addresses) to operate in dual-stack mode
- Y our EC2 instances will get at least a private internal IPv4 and a public IPv6
- They can communicate using either IPv4 or IPv6 to the internet through an Internet Gateway
EC2 Instance(Private IP: 10.0.0.5)(IPv6: 2001:db8::ff00:42:8329)
VPC
InternetGatewayIPv4 & IPv6
Internet

## Slide 764

IPv4 Troubleshooting
- IPv4 cannot be disabled for your VPC and subnets
- So, if you cannot launch an EC2 instance in your subnet
- It’s not because it cannot acquire an IPv6 (the space is very large)
- It’s because there are no available IPv4 in your subnet
- Solution: create a new IPv4 CIDR in your subnet
VPC(IPv4: 192.168.0.0/24)(IPv4: 10.0.0.0/24)(IPv6: 2001:db8:1234:5678::/56)
…
User
create
192.168.0.10192.168.0.15
10.0.0.35

## Slide 765

Egress-only Internet Gateway
- Used for IPv6 only
- (similar to a NAT Gateway but for IPv6)
- Allows instances in your VPC outbound connections over IPv6 while preventing the internet to initiate an IPv6 connection to your instances
- Yo u mu s t u p d a t e t h e R o u t e Ta bl e s
VPC
InternetGateway
Public Subnet
Internet
IPv6: 2001:db8::b1c2
IPv6: 2001:db8::e1c3
Egress-onlyInternet Gateway
initiate connectionsfrom both sides
can’t initiateconnections fromInternet
Private Subnet

## Slide 766

Private Subnet(IPv4: 10.0.1.0/24)(IPv6: 2001:db8:1234:1a02::/64)
Public Subnet(IPv4: 10.0.0.0/24)(IPv6: 2001:db8:1234:1a00::/64)
IPv6 RoutingRegion
VPC(IPv4: 10.0.0.0/16)(IPv6: 2001:db8:1234:1a00::/56)
Web serverPrivate IPv4: 10.0.0.5EIP: 198.51.100.1IPv6: 2001:db8:1234:1a00::123
ServerPrivate IPv4: 10.0.1.5IPv6: 2001:db8:1234:1a02::456
NAT Gateway(IPv4)EIP: 198.51.100.1
InternetGateway(IPv4 & IPv6)
Egress-onlyInternet Gateway(IPv6)
Route Table(Public Subnet)DestinationTarget10.0.0.0/16local2001:db8:1234:1a00::/56local0.0.0.0/0igw-id::/0igw-id
DestinationTarget10.0.0.0/16local2001:db8:1234:1a00::/56local0.0.0.0/0nat-gateway-id::/0eigw-id
Internet
Route Table(Private Subnet)

## Slide 767

VPC Section Summary (1/3)
- CIDR – IP Range
- VPC – Virtual Private Cloud => we define a list of IPv4 & IPv6 CIDR
- Subnets – tied to an AZ, we define a CIDR
- Internet Gateway – at the VPC level, provide IPv4 & IPv6 Internet Access
- Route Tables – must be edited to add routes from subnets to the IGW, VPC Peering Connections, VPC Endpoints, …
- Bastion Host – public EC2 instance to SSH into, that has SSH connectivity to EC2 instances in private subnets
- NAT Instances – gives Internet access to EC2 instances in private subnets. Old, must be setup in a public subnet, disable Source / Destination check flag
- NAT Gateway – managed by AWS, provides scalable Internet access to private EC2 instances, when the target is an IPv4 address

## Slide 768

VPC Section Summary (2/3)
- NACL – stateless, subnet rules for inbound and outbound, don’t forget Ephemeral Ports
- Security Groups – stateful, operate at the EC2 instance level
- VPC Peering – connect two VPCs with non overlapping CIDR, non-transitive
- VPC Endpoints – provide private access to AWS Services (S3, DynamoDB, CloudFormation, SSM) within a VPC
- VPC Flow Logs – can be setup at the VPC / Subnet / ENI Level, for ACCEPT and REJECT traffic, helps identifying attacks, analyze using Athena or CloudWatch Logs Insights
- Site-to-Site VPN – setup a Customer Gateway on DC, a Virtual Private Gateway on VPC, and site-to-site VPN over public Internet
- AWS VPN CloudHub – hub-and-spoke VPN model to connect your sites

## Slide 769

VPC Section Summary (3/3)
- Direct Connect – setup a Virtual Private Gateway on VPC, and establish a direct private connection to an AWS Direct Connect Location
- Direct Connect Gateway – setup a Direct Connect to many VPCs in different AWS regions
- AWS PrivateLink / VPC Endpoint Services:
- Connect services privately from your service VPC to customers VPC
- Doesn’t need VPC Peering, public Internet, NAT Gateway, Route Tables
- Must be used with Network Load Balancer & ENI
- ClassicLink – connect EC2-Classic EC2 instances privately to your VPC
- Transit Gateway – transitive peering connections for VPC, VPN & DX
- Traffic Mirroring – copy network traffic from ENIs for further analysis
- Egress-only Internet Gateway – like a NAT Gateway, but for IPv6 targets

## Slide 770

Networking Costs in AWS per GB - SimplifiedRegion
Availability ZoneAvailability ZoneRegion
Availability Zone
$0.02Inter-region
$0.02 if usingPublic IP / Elastic IP
Free if using private IP$0.01 if Using private IP
- Use Private IP instead of Public IP for good savings and better network performance
- Use same AZ for maximum savings (at the cost of high availability)
Free for traffic in

## Slide 771

Minimizing egress traffic network cost
- Egress traffic: outbound traffic (from AWS to outside)
- Ingress traffic: inbound traffic - from outside to AWS (typically free)
- Tr y to keep as much internet traffic within AWS to minimize costs
- Direct Connect location that are co-located in the same AWS Region result in lower cost for egress networkCorporate data center
AWS Cloud
DatabaseApplication
Egress cost is minimized
Corporate data center
AWS Cloud
DB Query100 MB
DatabaseApplication
Query Results50 KB
DB Query100 MBQuery Results50 KB
Egress cost is high

## Slide 772

S3 Data Transfer Pricing – Analysis for USA
- S3 ingress: free
- S3 to Internet: $0.09 per GB
- S3 Transfer Acceleration:
- Faster transfer times (50 to 500% better)
- Additional cost on top of Data Transfer Pricing: +$0.04 to $0.08 per GB
- S3 to CloudFront: $0.00 per GB
- CloudFront to Internet: $0.085 per GB (slightly cheaper than S3)
- Caching capability (lower latency)
- Reduce costs associated with S3 Requests Pricing (7x cheaper with CloudFront)
- S3 Cross Region Replication: $0.02 per GB
Replication $0.02
Edge location
Transfer acceleration +$0.04
internet
$0.09
$0.00
CloudFront
$0.085

## Slide 773

Private subnet 2(10.0.1.0/24)
Public subnet
Pricing: NAT Gateway vs Gateway VPC EndpointRegion(us-east-1)
VPC(10.0.0.0/16)
Private subnet 1(10.0.0.0/24)
EC2 Instance
NAT GatewayInternetGateway
VPC Endpoint
Internet
S3 Bucket
EC2 Instance
DestinationTarget10.0.0.0/16Local0.0.0.0/0igw-id
DestinationTarget10.0.0.0/16Localpl-id for Amazon S3vpce-id
Subnet 1 route table
Subnet 2 route table
$0.045 NAT Gateway / hour$0.045 NAT Gateway data processed / GB$0.09 Data transfer out to S3 (cross-region)$0.00 Data transfer out to S3 (same-region)
No cost for using Gateway Endpoint.$0.01 Data transfer in/out (same-region)

## Slide 774

Network Protection on AWS
- To p r o t e c t n e t wo r k o n AW S , we ’ ve s e e n
- Network Access Control Lists (NACLs)
- Amazon VPC security groups
- AWS WAF (protect against malicious requests)
- AWS Shield & AWS Shield Advanced
- AWS Firewall Manager (to manage them across accounts)
- But what if we want to protect in a sophisticated way our entire VPC?

## Slide 775

AWS Network Firewall
- Protect your entire Amazon VPC
- From Layer 3 to Layer 7 protection
- Any direction, you can inspect
- VPC to VPC traffic
- Outbound to internet
- Inbound from internet
- To / f r o m D i r e c t C o n n e c t & S i t e-to-Site VPN
- Internally, the AWS Network Firewall uses the AWS Gateway Load Balancer
- Rules can be centrally managed cross-account by AWS Firewall Manager to apply to many VPCs
VPC
Private subnet
Peered VPC
VPN connection
Direct Connect
Corporate DC
AWS Network Firewall
internet

## Slide 776

Network Firewall – Fine Grained Controls
- Supports 1000s of rules
- IP & port - example: 10,000s of IPs filtering
- Protocol – example: block the SMB protocol for outbound communications
- Stateful domain list rule groups: only allow outbound traffic to *.mycorp.com or third-party software repo
- General pattern matching using regex
- Traffic filtering: Allow, drop, or aler t for the traffic that matches the r ules
- Active flow inspection to protect against network threats with intrusion-prevention capabilities (like Gateway Load Balancer, but all managed by AWS)
- Send logs of rule matches to Amazon S3, CloudWatch Logs, Kinesis Data Firehose

## Slide 777

Disaster Recovery & Migrations

## Slide 778

Disaster Recovery Overview
- Any event that has a negative impact on a company’s business continuity or finances is a disaster
- Disaster recovery (DR) is about preparing for and recovering from a disaster
- What kind of disaster recovery?
- On-premise => On-premise: traditional DR, and very expensive
- On-premise => AWS Cloud: hybrid recovery
- AWS Cloud Region A => AWS Cloud Region B
- Need to define two terms:
- RPO: Recovery Point Objective
- RTO: Recovery Time Objective

## Slide 779

RPO and RTO
RPODisaster
Data loss
RTO
Downtime

## Slide 780

Disaster Recovery Strategies
- Backup and Restore
- Pilot Light
- Warm Standby
- Hot Site / Multi Site Approach
Faster RTO

## Slide 781

Backup and Restore (High RPO)Corporate data center
AWS Cloud
AWS Cloud
AWS Storage Gateway
Glacier
Amazon S3
AWS Snowball
EBS
Snapshot
Redshift
RDS
lifecycle
AWS Cloud
Amazon EC2
Amazon RDS
AMIScheduled regularsnapshots

## Slide 782

Disaster Recovery – Pilot Light
- A small version of the app is always running in the cloud
- Useful for the critical core (pilot light)
- Ver y similar to Backup and Restore
- Faster than Backup and Restore as critical systems are already upCorporate data center
AWS Cloud
RDS (running)
Data ReplicationEC2 (not running)
Route 53

## Slide 783

Warm Standby
- Full system is up and running, but at minimum size
- Upon disaster, we can scale to production loadCorporate data center
AWS Cloud
RDS Secondary (running)
Data Replication
Route 53
ELB
EC2 Auto Scaling(minimum)
Reverse proxy
App Server
PrimaryDB
failover

## Slide 784

Multi Site / Hot Site Approach
- Ver y low RTO (minutes or seconds) – very expensive
- Full Production Scale is running AWS and On PremiseCorporate data center
AWS Cloud
RDS secondary (running)
Data Replication
Route 53
ELB
EC2 Auto Scaling(production)
Reverse proxy
App Server
PrimaryDB
failover
activeactive

## Slide 785

All AWS Multi RegionAWS Cloud
Data Replication
Route 53
ELB
EC2 Auto Scaling(production)
failover
activeactiveAWS Cloud
EC2 Auto Scaling(production)
ELB
Aurora Global (primary)
Aurora Global (secondary)

## Slide 786

Disaster Recovery Tips
- Backup
- EBS Snapshots, RDS automated backups / Snapshots, etc…
- Regular pushes to S3 / S3 IA / Glacier, Lifecycle Policy, Cross Region Replication
- From On-Premise: Snowball or Storage Gateway
- High Availability
- Use Route53 to migrate DNS over from Region to Region
- RDS Multi-AZ, ElastiCache Multi-AZ, EFS, S3
- Site to Site VPN as a recovery from Direct Connect
- Replication
- RDS Replication (Cross Region), AWS Aurora + Global Databases
- Database replication from on-premises to RDS
- Storage Gateway
- Automation
- CloudFormation / Elastic Beanstalk to re-create a whole new environment
- Recover / Reboot EC2 instances with CloudWatch if alarms fail
- AWS Lambda functions for customized automations
- Chaos
- Netflix has a “simian-army” randomly terminating EC2

## Slide 787

AWS Elastic Disaster Recovery (DRS)
- Used to be named “CloudEndure Disaster Recovery”
- Quickly and easily recover your physical, virtual, and cloud-based servers into AWS
- Example: protect your most critical databases (including Oracle, MySQL, and SQL Server), enterprise apps (SAP), protect your data from ransomware attacks, …
- Continuous block-level replication for your serversCorporate Data Center / Any cloud
Disks
OSAppsDB
AWS Cloud
AWS ReplicationAgent
Staging
Production
failover(minutes)Low-cost EC2 instances& EBS volumesTarget EC2 instances& EBS volumes
continuous replication(seconds)
Elastic Disaster Recovery
failback

## Slide 788

DMS – Database Migration Service
- Quickly and securely migrate databases to AWS, resilient, self healing
- The source database remains available during the migration
- Supports:
- Homogeneous migrations: ex Oracle to Oracle
- Heterogeneous migrations: ex Microsoft SQL Server to Aurora
- Continuous Data Replication using CDC
- Y ou must create an EC2 instance to perform the replication tasks
EC2 instanceRunning DMS
Source DB
Target DB

## Slide 789

DMS Sources and TargetsSOURCES:
- On-Premises and EC2 instances databases: Oracle, MS SQL Server, MySQL, MariaDB, PostgreSQL, MongoDB, SAP , DB2
- Azure: Azure SQL Database
- Amazon RDS: all including Aurora
- Amazon S3
- DocumentDB
TARGETS:
- On-Premises and EC2 instances databases: Oracle, MS SQL Server, MySQL, MariaDB, PostgreSQL, SAP
- Amazon RDS
- Redshift, DynamoDB, S3
- OpenSearch Service
- Kinesis Data Streams
- Apache Kafka
- DocumentDB & Amazon Neptune
- Redis & Babelfish

## Slide 790

AWS Schema Conversion Tool (SCT)
- Convert your Database’s Schema from one engine to another
- Example OLTP: (SQL Server or Oracle) to MySQL, PostgreSQL, Aurora
- Example OLAP: (Teradata or Oracle) to Amazon Redshift
- Prefer compute-intensive instances to optimize data conversions
- Yo u d o n o t n e e d t o u s e S C T i f yo u a r e m i g r a t i n g t h e s a m e D B e n g i n e
- Ex: On-Premise PostgreSQL => RDS PostgreSQL
- The DB engine is still PostgreSQL (RDS is the platform)
DMS + SCTSource DB
Target DB (different engine)

## Slide 791

DMS - Continuous ReplicationCorporate data center
Oracle DB(source)
VPC
Region
Public Subnet
AWS DMS Replication Instance
Private Subnet
Amazon RDS for MySQL DB(target)
Full load +CDC
Schema conversion
Data migration
Server withAWS SCT installed

## Slide 792

AWS DMS – Multi-AZ DeploymentAWS RegionAvailability Zone - A
DMS ReplicationInstance
Availability Zone - B
DMS ReplicationInstance(Standby Replica)
synchronousreplication
- When Multi-AZ Enabled, DMS provisions and maintains a synchronously stand replica in a different AZ
- Advantages:
- Provides Data Redundancy
- Eliminates I/O freezes
- Minimizes latency spikes

## Slide 793

RDS & Aurora MySQL Migrations
- RDS MySQL to Aurora MySQL
- Option 1: DB Snapshots from RDS MySQL restored as MySQL Aurora DB
- Option 2: Create an Aurora Read Replica from your RDS MySQL, and when the replication lag is 0, promote it as its own DB cluster (can take time and cost $)
- External MySQL to Aurora MySQL
- Option 1:
- Use Percona XtraBackup to create a file backup in Amazon S3
- Create an Aurora MySQL DB from Amazon S3
- Option 2:
- Create an Aurora MySQL DB
- Use the mysqldump utility to migrate MySQL into Aurora (slower than S3 method)
- Use DMS if both databases are up and running
Read Replica
PerconaXtraBackup
import
mysqldump

## Slide 794

RDS & Aurora PostgreSQL Migrations
- RDS PostgreSQL to Aurora PostgreSQL
- Option 1: DB Snapshots from RDS PostgreSQL restored as PostgreSQL Aurora DB
- Option 2: Create an Aurora Read Replica from your RDS PostgreSQL, and when the replication lag is 0, promote it as its own DB cluster (can take time and cost $)
- External PostgreSQL to Aurora PostgreSQL
- Create a backup and put it in Amazon S3
- Import it using the aws_s3 Aurora extension
- Use DMS if both databases are up and running
Read Replica
backup
import

## Slide 795

On-Premise strategy with AWS
- Ability to download Amazon Linux 2 AMI as a VM (.iso format)
- VMWare, KVM, VirtualBox (Oracle VM), Microsoft Hyper-V
- VM Import / Export
- Migrate existing applications into EC2
- Create a DR repository strategy for your on-premises VMs
- Can export back the VMs from EC2 to on-premises
- AWS Application Discovery Service
- Gather information about your on-premises servers to plan a migration
- Server utilization and dependency mappings
- Track with AWS Migration Hub
- AWS Database Migration Service (DMS)
- replicate On-premise => AWS , AWS => AWS, AWS => On-premise
- Works with various database technologies (Oracle, MySQL, DynamoDB, etc..)
- AWS Application Migration Service (MGN)
- Incremental replication of on-premises live servers to AWS

## Slide 796

AWS Backup
- Fully managed service
- Centrally manage and automate backups across AWS services
- No need to create custom scripts and manual processes
- Supported services:
- Amazon EC2 / Amazon EBS
- Amazon S3
- Amazon RDS (all DBs engines) / Amazon Aurora / Amazon DynamoDB
- Amazon DocumentDB / Amazon Neptune
- Amazon EFS / Amazon FSx (Lustre & Windows File Server)
- AWS Storage Gateway (Volume Gateway)
- Supports cross-region backups
- Supports cross-account backups

## Slide 797

AWS Backup
- Supports PITR for supported services
- On-Demand and Scheduled backups
- Ta g-based backup policies
- Yo u c r e a t e b a c k u p p o l i c i e s k n ow n a s Backup Plans
- Backup frequency (every 12 hours, daily, weekly, monthly, cron expression)
- Backup window
- Transition to Cold Storage (Never, Days, Weeks, Months, Year s)
- Retention Period (Always, Days, Weeks, Months, Y ears)

## Slide 798

AWS Backup
AWS Backup
Create Backup Plan(frequency, retention policy)
Assign AWS Resources
EC2EBS
DynamoDBRDS
EFS
Aurora
FSxStorageGateway
Amazon S3Automaticallybacked up to
S3
DocumentDB
Neptune

## Slide 799

AWS Backup Vault Lock
- Enforce a WORM (Write Once Read Many) state for all the backups that you store in your AWS Backup Vault
- Additional layer of defense to protect your backups against:
- Inadvertent or malicious delete operations
- Updates that shorten or alter retention periods
- Even the root user cannot delete backups when enabled
backup
Backup Vault Lock PolicyBackups can’t be deleted

## Slide 800

AWS Application Discovery Service
- Plan migration projects by gathering information about on-premises data centers
- Server utilization data and dependency mapping are important for migrations
- Agentless Discovery (AWS Agentless Discovery Connector)
- VM inventory, configuration, and performance history such as CPU, memory, and disk usage
- Agent-based Discovery (AWS Application Discovery Agent)
- System configuration, system performance, running processes, and details of the network connections between systems
- Resulting data can be viewed within AWS Migration Hub

## Slide 801

AWS Application Migration Service (MGN)
- The “AWS evolution” of CloudEndure Migration, replacing AWS Server Migration Service (SMS)
- Lift-and-shift (rehost) solution which simplify migrating applications to AWS
- Converts your physical, virtual, and cloud-based servers to run natively on AWS
- Supports wide range of platforms, Operating Systems, and databases
- Minimal downtime, reduced costsCorporate Data Center / Any cloud
Disks
OSAppsDB
AWS Cloud
AWS ReplicationAgent
Staging
Production
cutover
Low-cost EC2 instances& EBS volumesTarget EC2 instances& EBS volumes
continuous replication
Application Migration Service

## Slide 802

VMware Cloud on AWS
- Some customers use VMware Cloud to manage their on-premises Data Center
- They want to extend the Data Center capacity to AWS, but keep using the VMware Cloud software
- …Enter VMware Cloud on AWS
- Use cases
- Migrate your VMware vSphere-based workloads to AWS
- Run your production workloads across VMware vSphere-based private, public, and hybrid cloud environments
- Have a disaster recover strategy
Customer Data Center
AWS Cloud
On-Premises vCentervSphere-basedenvironment
vSphereVMware Cloudon AWS
AWS Services
AmazonEC2AmazonS3DirectConnectAmazonFSxAmazonRDSAmazonRedshift

## Slide 803

Transferring large amount of data into AWS
- Example: transfer 200 TB of data in the cloud. We have a 100 Mbps internet connection.
- Over the internet / Site-to-Site VPN:
- Immediate to setup
- Will take 200(TB)*1000(GB)*1000(MB)*8(Mb)/100 Mbps = 16,000,000s = 185d
- Over direct connect 1Gbps:
- Long for the one-time setup (over a month)
- Will take 200(TB)*1000(GB)*8(Gb)/1 Gbps = 1,600,000s = 18.5d
- Over Snowball:
- Takes about 1 week for the end-to-end transfer
- Can be combined with DMS
- For on-going replication / transfers: Site-to-Site VPN or DX with DMS or DataSync

## Slide 804

More Solutions Architecture

## Slide 805

Lambda, SNS & SQS
SQS
SNS
retries
SQS
DLQ
Try, retry
SQS FIFO
DLQ
Try, retryblocking DLQ
asynchronous(poll)
SQS + Lambda
SQS FIFO + LambdaSNS + Lambda

## Slide 806

Fan Out Pattern: deliver to multiple SQS
SDK
PUT #1
PUT #2
PUT #3
Option 1
SDK
Option 2 – Fan Out
SNS
SQS SQS
PUTsubscribe

## Slide 807

S3 Event Notifications
Amazon S3
events
Lambda Function
SQS
SNS
- S3:ObjectCreated, S3:ObjectRemoved, S3:ObjectRestore, S3:Replication…
- Object name filtering possible (*.jpg)
- Use case: generate thumbnails of images uploaded to S3
- Can create as many “S3 events” as desired
- S3 event notifications typically deliver events in seconds but can sometimes take a minute or longer

## Slide 808

S3 Event Notifications with Amazon EventBridge
Amazon S3bucket
eventsAll events
AmazonEventBridge
rulesOver 18 AWS servicesas destinations
- Advanced filtering options with JSON rules (metadata, object size, name...)
- Multiple Destinations – ex Step Functions, Kinesis Streams / Firehose…
- EventBridge Capabilities – Archive, Replay Events, Reliable delivery

## Slide 809

Amazon EventBridge – Intercept API Calls
CloudTrail(any API call)AmazonEventBridge
SNS
event
DynamoDB
Log API call
User
alert
DeleteTable API Call
💥

## Slide 810

API Gateway – AWS Service IntegrationKinesis Data Streams example
API GatewayKinesis DataStreamsKinesis DataFirehoseAmazon S3sendrecordsstore .jsonfiles
Client
requests

## Slide 811

Caching Strategies
API Gateway
CloudFront (edge)
Redis Memcached DAX
S3
CloudFront
Caching, TTL, Network, Computation, Cost, Latency
Client
App logicEC2 / Lambda
Database

## Slide 812

Blocking an IP address
Client
VPC
Public Subnet
Security Group (allow rules)
EC2 Instancepublic IP + Firewall Software (optional)
NACLDeny + Allow rules

## Slide 813

Blocking an IP address – with an ALB
Client
VPC
Public Subnet
EC2 Security Group
EC2 InstancePrivate IP
Private Subnet
ALB Security Group
Application Load BalancerConnection TerminationNACL

## Slide 814

Blocking an IP address – with an NLB
Client
VPC
Public Subnet
EC2 Security Group
EC2 InstancePrivate IP
NACL
Private Subnet
NLB Security Group
Network Load Balancer

## Slide 815

Blocking an IP address – ALB + WAF
Client
VPC
Public Subnet
EC2 Security Group
EC2 InstancePrivate IP
NACL
Private Subnet
ALB Security Group
Application Load Balancer
AWS WAFIP Address Filtering

## Slide 816

Blocking an IP address – ALB, CloudFront & WAF
Client
VPC
Public Subnet
EC2 Security Group
EC2 InstancePrivate IP
NACL
Private Subnet
ALB Security Group
Application Load BalancerPublic
CloudFront(Geo Restriction)
AWS WAF(IP Address Filtering)
CloudFrontPublic IPs
NOT helpful

## Slide 817

High Performance Computing (HPC)
- The cloud is the perfect place to perform HPC
- Yo u c a n c r e a t e a ve r y h i g h nu m b e r o f r e s o u r c e s i n n o t i m e
- Yo u c a n s p e e d u p t i m e t o r e s u l t s by a d d i n g m o r e r e s o u r c e s
- Yo u c a n p ay o n ly fo r t h e s y s t e m s yo u h ave u s e d
- Perform genomics, computational chemistry, financial risk modeling, weather prediction, machine learning, deep learning, autonomous driving
- Which services help perform HPC?

## Slide 818

Data Management & Transfer
- AWS Direct Connect:
- Move GB/s of data to the cloud, over a private secure network
- Snowball & Snowmobile
- Move PB of data to the cloud
- AWS DataSync
- Move large amount of data between on-premises and S3, EFS, FSx for Windows

## Slide 819

Compute and Networking
- EC2 Instances:
- CPU optimized, GPU optimized
- Spot Instances / Spot Fleets for cost savings + Auto Scaling
- EC2 Placement Groups: Cluster for good network performance
Same RackSame AZEC2EC2
EC2EC2EC2
EC2 Placement group ClusterLow latency10Gbps network

## Slide 820

Compute and Networking
- EC2 Enhanced Networking (SR-IOV)
- Higher bandwidth, higher PPS (packet per second), lower latency
- Option 1: Elastic Network Adapter (ENA) up to 100 Gbps
- Option 2: Intel 82599 VF up to 10 Gbps – LEGACY
- Elastic Fabric Adapter (EFA)
- Improved ENA for HPC, only works for Linux
- Great for inter-node communications, tightly coupled workloads
- Leverages Message Passing Interface (MPI) standard
- Bypasses the underlying Linux OS to provide low-latency, reliable transport

## Slide 821

Storage
- Instance-attached storage:
- EBS: scale up to 256,000 IOPS with io2 Block Express
- Instance Store: scale to millions of IOPS, linked to EC2 instance, low latency
- Network storage:
- Amazon S3: large blob, not a file system
- Amazon EFS: scale IOPS based on total size, or use provisioned IOPS
- Amazon FSx for Lustre:
- HPC optimized distributed file system, millions of IOPS
- Backed by S3

## Slide 822

Automation and Orchestration
- AWS Batch
- AWS Batchsupports multi-node parallel jobs, which enables you to run single jobs that span multipleEC2instances.
- Easily schedule jobs and launch EC2 instances accordingly
- AWS ParallelCluster
- Open-source cluster management tool to deploy HPC on AWS
- Configure with text files
- Automate creation of VPC, Subnet, cluster type and instance types
- Ability to enable EFA on the cluster (improves network performance)

## Slide 823

Creating a highly available EC2 instance
Public EC2
Elastic IP Address
What time is it?
5:30 pm!
Standby EC2 instance
Attachment
CloudWatch Event(or Alarm based on metric)
Start the instanceAttach the Elastic IP
monitor

## Slide 824

Creating a highly available EC2 instanceWith an Auto Scaling Group
Public EC2
Elastic IP Address
What time is it?
5:30 pm!
Replacement EC2 instance
Auto Scaling group
EC2 User DataAttachmentBased on Tag
Availability Zone 1
Availability Zone 2
ASG Settings1 min1 max1 desired>= 2 AZEC2 user data to attachThe Elastic IPEC2 instance role to Allow API calls to attachThe Elastic IP

## Slide 825

Creating a highly available EC2 instanceWith ASG + EBS
Public EC2
Elastic IP Address
What time is it?
5:30 pm!
Replacement EC2 instance
Auto Scaling group
EC2 User DataAttachmentBased on Tag
Availability Zone 1
Availability Zone 2
EBS Volume
EBS
EBS Snapshot+ tags
EBS SnapshotOn ASG Terminate lifecycle hook
EBS Volume created + attachedOn ASG Launch lifecycle hook

## Slide 826

Other ServicesOverview of Services that might come up in a few questions

## Slide 827

What is CloudFormation
- CloudFormation is a declarative way of outlining your AWS Infrastructure, for any resources (most of them are supported).
- For example, within a CloudFormation template, you say:
- I want a security group
- I want two EC2 instances using this security group
- I want an S3 bucket
- I want a load balancer (ELB) in front of these machines
- Then CloudFormation creates those for you, in the right order, with the exact configuration that you specify

## Slide 828

Benefits of AWS CloudFormation (1/2)
- Infrastructure as code
- No resources are manually created, which is excellent for control
- Changes to the infrastructure are reviewed through code
- Cost
- Each resources within the stack is tagged with an identifier so you can easily see how much a stack costs you
- Yo u c a n e s t i m a t e t h e c o s t s o f yo u r r e s o u r c e s u s i n g t h e C l o u d F o r m a t i o n t e m p l a t e
- Savings strategy: In Dev, you could automation deletion of templates at 5 PM and recreated at 8 AM, safely

## Slide 829

Benefits of AWS CloudFormation (2/2)
- Productivity
- Ability to destroy and re-create an infrastructure on the cloud on the fly
- Automated generation of Diagram for your templates!
- Declarative programming (no need to figure out ordering and orchestration)
- Don’t re-invent the wheel
- Leverage existing templates on the web!
- Leverage the documentation
- Supports (almost) all AWS resources:
- Everything we’ll see in this course is supported
- Yo u c a n u s e “ c u s t o m r e s o u r c e s ” fo r r e s o u r c e s t h a t a r e n o t s u p p o r t e d

## Slide 830

CloudFormation + Infrastructure Composer
- Example: WordPress CloudFormation Stack
- We can see all the resources
- We can see the relations between the components
CloudFormationInfrastructureComposer
+

## Slide 831

CloudFormation – Service Role
- IAM role that allows CloudFormation to create/update/delete stack resources on your behalf
- Give ability to users to create/update/delete the stack resources even if they don’t have permissions to work with the resources in the stack
- Use cases:
- Yo u w a n t t o a c h i e ve t h e l e a s t p r i v i l e g e p r i n c i p l e
- But you don’t want to give the user all the required permissions to create the stack resources
- User must have iam:PassRole permissions
-cloudformation:*-iam:PassRole
Permissions
- s3:*BucketService Role
User
Template
CloudFormation
Stack
S3 bucket

## Slide 832

Amazon Simple Email Service (Amazon SES)
- Fully managed service to send emails securely, globally and at scale
- Allows inbound/outbound emails
- Reputation dashboard, performance insights, anti-spam feedback
- Provides statistics such as email deliveries, bounces, feedback loop results, email open
- Supports DomainKeys Identified Mail (DKIM) and Sender Policy Framework (SPF)
- Flexible IP deployment: shared, dedicated, and customer-owned IPs
- Send emails using your application using AWS Console, APIs, or SMTP
- Use cases: transactional, marketing and bulk email communications
Amazon SES
Application
Users
APIsor SMTP
bulk emails

## Slide 833

Amazon Pinpoint
- Scalable 2-way (outbound/inbound) marketing communications service
- Supports email, SMS, push, voice, and in-app messaging
- Ability to segment and personalize messages with the right content to customers
- Possibility to receive replies
- Scales to billions of messages per day
- Use cases: run campaigns by sending marketing, bulk, transactional SMS messages
- Versus Amazon SNS or Amazon SES
- In SNS & SES you managed each message's audience, content, and delivery schedule
- In Amazon Pinpoint, you create message templates, delivery schedules, highly-targeted segments, and full campaigns
SNSKinesis DataFirehoseCloudWatchLogs
AmazonPinpoint
stream events(e.g., TEXT_SUCCESS,TEXT_DELIVERED, …)
Customers
SMS

## Slide 834

Systems Manager – SSM Session Manager
- Allows you to start a secure shell on your EC2 and on-premises servers
- No SSH access, bastion hosts, or SSH keys needed
- No port 22 needed (better security)
- Supports Linux, macOS, and Windows
- Send session log data to S3 or CloudWatch Logs
EC2 Instance(SSM Agent)
User
IAMPermissions
Execute commandsSession Manager

## Slide 835

Systems Manager – Run Command
- Execute a document (= script) or just run a command
- Run command across multiple instances (using resource groups)
- No need for SSH
- Command Output can be shown in the AWS Console, sent to S3 bucket or CloudWatch Logs
- Send notifications to SNS about command status (In progress, Success, Failed, …)
- Integrated with IAM & CloudTrail
- Can be invoked using EventBridge
EventBridge
Run CommandAmazon S3
CloudWatchLogs
trigger
output
Amazon SNSnotification
EC2 Instances(with SSM Agent)
EC2 Instances(with SSM Agent)

## Slide 836

Systems Manager – Patch Manager
- Automates the process of patching managed instances
- OS updates, applications updates, security updates
- Supports EC2 instances and on-premises servers
- Supports Linux, macOS, and Windows
- Patch on-demand or on a schedule using Maintenance Windows
- Scan instances and generate patch compliance report (missing patches)
AWS ConsoleAWS SDKMaintenanceWindows
Run Command
runAWS-RunBatchBaseline
EC2 Instances(with SSM Agent)
EC2 Instances(with SSM Agent)

## Slide 837

Systems Manager – Maintenance Windows
- Defines a schedule for when to perform actions on your instances
- Example: OS patching, updating drivers, installing software, …
- Maintenance Window contains
- Schedule
- Duration
- Set of registered instances
- Set of registered tasks
EC2 Instances(with SSM Agent)
EC2 Instances(with SSM Agent)
Maintenance WindowsRun Command
trigger every 24 hourupdate

## Slide 838

Systems Manager - Automation
- Simplifies common maintenance and deployment tasks of EC2 instances and other AWS resources
- Examples: restart instances, create an AMI, EBS snapshot
- Automation Runbook – SSM Documents to define actions preformed on your EC2 instances or AWS resources (pre-defined or custom)
- Can be triggered using:
- Manually using AWS Console, AWS CLI or SDK
- Amazon EventBridge
- On a schedule using Maintenance Windows
- By AWS Config for rules remediations
AWS ConsoleAWS SDKMaintenanceWindows
AmazonEventBridge
SSM AutomationRunbooks(automation documents)
execute automation(AWS-RestartEC2Instance)
EC2 Instances
AWS Resources
EBS
AMI
RDS…
execute
AWS ConfigRemediation

## Slide 839

Cost Explorer
- Visualize, understand, and manage your AWS costs and usage over time
- Create custom reports that analyze cost and usage data.
- Analyze your data at a high level: total costs and usage across all accounts
- Or Monthly, hourly, resource level granularity
- Choose an optimal Savings Plan (to lower prices on your bill)
- Forecast usage up to 18 months based on previous usage

## Slide 840

Cost Explorer – Monthly Cost by AWS Service

## Slide 841

Cost Explorer– Hourly & Resource Level

## Slide 842

Cost Explorer – Savings PlanAlternative to Reserved Instances

## Slide 843

Cost Explorer – Forecast Usage

## Slide 844

AWS Cost Anomaly Detection
- Continuously monitor your cost and usage using ML to detect unusual spends
- It learns your unique, historic spend patterns to detect one-time cost spike and/or continuous cost increases (you don’t need to define thresholds)
- Monitor AWS services, member accounts, cost allocation tags, or cost categories
- Sends you the anomaly detection report with root-cause analysis
- Get notified with individual alerts or daily/weekly summary (using SNS)
AWS Cost AnomalyDetectionreduce cost surpriseswith Machine Learning
Create Cost MonitorIdentify unusual spend atthe granularity levelthat you specify
Get AlertedReceive alerts whenunusual spend is detectedAnalyze Root CauseAnalyze the root causebehind the anomaly andthe impact on your costs

## Slide 845

AWS Outposts
- Hybrid Cloud: businesses that keep an on-premises infrastructure alongside a cloud infrastructure
- Therefore, two ways of dealing with IT systems:
- One for the AWS cloud (using the AWS console, CLI, and AWS APIs)
- One for their on-premises infrastructure
- AWS Outposts are “server racks” that offers the same AWS infrastructure, services, APIs & tools to build your own applications on-premises just as in the cloud
- AWS will setup and manage “Outposts Racks” within your on-premises infrastructure and you can start leveraging AWS services on-premises
- Yo u a r e r e s p o n s i bl e fo r t h e O u t p o s t s R a c k physical security
Corporate data center
AWS Cloud
On-prem serversOutpostsRacks
Extension of AWS services

## Slide 846

AWS Outposts
- Benefits:
- Low-latency access to on-premises systems
- Local data processing
- Data residency
- Easier migration from on-premises to the cloud
- Fully managed service
- Some services that work on Outposts:
Amazon EC2Amazon EBSAmazon S3Amazon EKSAmazon ECSAmazon RDSAmazon EMR

## Slide 847

AWS Batch
- Fully managed batch processing at any scale
- Efficiently run 100,000s of computing batch jobs on AWS
- A “batch” job is a job with a start and an end (opposed to continuous)
- Batch will dynamically launch EC2 instances or Spot Instances
- AWS Batch provisions the right amount of compute / memory
- Yo u s u b m i t o r s c h e d u l e b a t c h j o b s a n d AW S B a t c h d o e s t h e r e s t !
- Batch jobs are defined as Docker images and run on ECS
- Helpful for cost optimizations and focusing less on the infrastructure

## Slide 848

AWS Batch – Simplified Example
Trigger
AWS Batch
EC2 Instance
Spot InstanceAmazon S3 Insert processed objectAmazon S3
ECS

## Slide 849

Batch vs Lambda
- Lambda:
- Time limit
- Limited runtimes
- Limited temporary disk space
- Serverless
- Batch:
- No time limit
- Any runtime as long as it’s packaged as a Docker image
- Rely on EBS / instance store for disk space
- Relies on EC2 (can be managed by AWS)

## Slide 850

Amazon AppFlow
- Fully managed integration service that enables you to securely transfer data between Software-as-a-Service (SaaS) applications and AWS
- Sources: Salesforce, SAP, Zendesk, Slack, and Ser viceNow
- Destinations: AWS services like Amazon S3, Amazon Redshift or non-AWS such as SnowFlake and Salesforce
- Frequency: on a schedule, in response to events, or on demand
- Data transformation capabilities like filtering and validation
- Encrypted over the public internet or privately over AWS PrivateLink
- Don’t spend time writing the integrations and leverage APIs immediately

## Slide 851

Amazon AppFlow
AmazonAppFlow
Sources Destinations

## Slide 852

AWS Amplify - web and mobile applications
- A set of tools and services that helps you develop and deploy scalable full stack web and mobile applications
- Authentication, Storage, API (REST, GraphQL), CI/CD, PubSub, Analytics, AI/ML Predictions, Monitoring, …
- Connect your source code from GitHub, AWS CodeCommit, Bitbucket, GitLab, or upload directly
Amplify backendFrontend
configure backendusing Amplify CLI
Amazon S3Amazon CognitoAWS AppSyncAPIGateway
DynamoDBLambdaAmazonSageMakerAmazon Lex
…
connect frontend to backendusingAmplify Frontend Libraries
……
AmplifyConsoleAmazon CloudFront
build using Amplify Console & deploy

## Slide 853

Instance Scheduler on AWS
- AWS solution deployed through CloudFormation (not a service)
- Automatically start/stop your AWS services to reduce costs (up to 70%)
- Example: stop company’s EC2 instances outside business hours
- Supports EC2 instances, EC2 Auto Scaling Groups, and RDS instances
- Schedules are managed in a DynamoDB table
- Uses resources’ tags and Lambda to stop/start instances
- Supports cross-account and cross-region resources
- https://aws.amazon.com/solutions/implementations/instance-scheduler-on-aws/

## Slide 854

White Papers & ArchitecturesWell Architected Framework, Disaster Recover y, etc…

## Slide 855

Section Overview
- Well Architected Framework Whitepaper
- Well Architected Tool
- AWS Trusted Advisor
- Reference architectures resources (for real-world)
- Disaster Recovery on AWS Whitepaper

## Slide 856

Well Architected Framework General Guiding Principles
- https://aws.amazon.com/architecture/well-architected
- Stop guessing your capacity needs
- Test syst em s a t pr odu ct ion sca le
- Automate to make architectural experimentation easier
- Allow for evolutionary architectures
- Design based on changing requirements
- Drive architectures using data
- Improve through game days
- Simulate applications for flash sale days

## Slide 857

Well Architected Framework 6 Pillars
- 1) Operational Excellence
- 2) Security
- 3) Reliability
- 4) Performance Efficiency
- 5) Cost Optimization
- 6) Sustainability
- They are not something to balance, or trade-offs, they’re a synergy

## Slide 858

AWS Well-Architected T ool
- Free tool to review your architectures against the 6 pillars Well-Architected Framework and adopt architectural best practices
- How does it work?
- Select your workload and answer questions
- Review your answers against the 6 pillars
- Obtain advice: get videos and documentations, generate a report, see the results in a dashboard
- Let’s have a look: https://console.aws.amazon.com/wellarchitected
https://aws.amazon.com/blogs/aws/new-aws-well-architected-tool-review-workloads-against-best-practices/

## Slide 859

Tr usted Advisor
- No need to install anything – high level AWS account assessment
- Analyze your AWS accounts and provides recommendation on 6 categories:
- Cost optimization
- Performance
- Security
- Fault tolerance
- Service limits
- Operational Excellence
- Business & Enterprise Support plan
- Full Set of Checks
- Programmatic Access using AWS Suppor t API

## Slide 860

More Architecture Examples
- We’ve explored the most impor tant architectural patterns:
- Classic: EC2, ELB, RDS, ElastiCache, etc…
- Serverless: S3, Lambda, DynamoDB, CloudFront, API Gateway, etc…
- If you want to see more AWS architectures:
- https://aws.amazon.com/architecture/
- https://aws.amazon.com/solutions/

## Slide 861

Exam Review & Tips

## Slide 862

State of learning checkpoint
- Let’s look how far we’ve gone on our learning journey
- https://aws.amazon.com/certification/certified-solutions-architect-associate/

## Slide 863

Practice makes perfect
- If you’re new to AWS, take a bit of AWS practice thanks to this course before rushing to the exam
- The exam recommends you to have one or more years of hands-on experience on AWS
- Practice makes perfect!
- If you feel overwhelmed by the amount of knowledge you just learned, just go through it one more time

## Slide 864

Proceed by elimination
- Most questions are going to be scenario based
- For all the questions, rule out answers that you know for sure are wrong
- For the remaining answers, understand which one makes the most sense
- There are very few trick questions
- Don’t over-think it
- If a solution seems feasible but highly complicated, it’s probably wrong

## Slide 865

Skim the AWS Whitepapers
- Yo u c a n r e a d a b o u t s o m e AW S W h i t e Pa p e r s h e r e :
- Architecting for the Cloud: AWS Best Practices
- AWS Well-Architected Framework
- AWS Disaster Recovery (https://aws.amazon.com/disaster-recovery/)
- Overall we’ve explored all the most important concepts in the course
- It’s never bad to have a look at the whitepapers you think are interesting!

## Slide 866

Read each service’s FAQ
- FAQ = Frequently asked questions
- Example: https://aws.amazon.com/vpc/faqs/
- FAQ cover a lot of the questions asked at the exam
- They help confirm your understanding of a service
