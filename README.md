# pulumi-ec2-go

# Deploying EC2 with Pulumi and Go

## 🛠 Requirements
- [Pulumi](https://www.pulumi.com/docs/get-started/install/) installed
- [AWS CLI](https://aws.amazon.com/cli/) configured
- AWS Credentials with EC2 Permissions
- Go 1.16+ (for Go projects)

## 🚀 Rapid deployment

1. **Clone the repository:**
   ```bash
   git clone https://github.com/tu-usuario/ec2-pulumi-go.git
   cd ec2-pulumi-go

2. **Setting AWS Credentials:**  
    ```bash
    aws configure
    ```

3. **Deploy:**  
    ```bash
    pulumi up
    ```

4. **Get the public IP:**
    ```bash
    pulumi stack output instancePublicIP

5. **Connection to the Instance:**
    ```bash
    ssh -i ~/.ssh/my-key.pem ubuntu@publicIP

## ⚙️ Handling Dependency Issues  

If you encounter an error due to outdated dependencies in `go.mod`, run the following command to update them:  

```bash
go mod tidy
