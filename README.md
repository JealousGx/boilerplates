# Serverless AWS CDK Backend API Boilerplate

This repository provides a boilerplate for building a backend API using Go and AWS CDK in TypeScript. It is designed to help you quickly set up a serverless application with AWS Lambda, API Gateway, and other AWS services.

## Features

- **Go Backend**: The backend is written in Go, providing a fast and efficient environment for building APIs.
- **AWS CDK**: Infrastructure is defined using AWS CDK in TypeScript, allowing for infrastructure as code with modern programming constructs.
- **Serverless Architecture**: Utilizes AWS Lambda and API Gateway to create a serverless API.
- **Credential Management**: Uses AWS SDK for Go to manage credentials and access AWS services.
- **Environment Configuration**: Supports environment variables for configuration.
- **Testing and Development Tools**: Includes Jest for testing and SAM for local development.

API endpoints are defined in `template.yml` file for local development & for production, the cdk file is configured in `./lib/serverless-aws-cdk-stack.ts`. Right now, the project is configured for 2 endpoints: `/api/v1/test` & `/api/v1/users` with additional endpoints configured as proxy.

## Prerequisites

- **Go**: Ensure you have Go installed on your machine. This project is compatible with Go 1.22.5.
- **Node.js**: Node.js is required for AWS CDK. Ensure you have Node.js version 20.10.0 or later.
- **Docker**: Docker is required for running the DynamoDB local container.
- **AWS CLI**: AWS CLI should be configured with your AWS credentials.
- **AWS CDK**: Install AWS CDK globally using npm.
- **AWS SAM CLI**: AWS SAM CLI is required for running the API locally.

```bash
npm install -g aws-cdk
```

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/JealousGx/boilerplates/tree/serverless-aws-cdk.git
cd serverless-aws-cdk
```

### Install Dependencies

#### Go Dependencies

```bash
cd pkg && go mod tidy
```

#### Node Dependencies

Run the following command at the root of the project:

```bash
npm install
```

### Environment Configuration

Rename `.env.example` file in the root directory to `.env` and set the required variables.

### Build and Deploy

#### Build the Go Application

```bash
npm run build
```

```bash
npm run deploy
```

### Local Development

Run the following commands at the root of the project to create the necessary resources for local development such as the DynamoDB local container and the lambda network:

```bash
npm run create-dynamodb-image
npm run create-lambda-network
npm run create-table:local
```

Run the following command to start the API locally:

```bash
npm run start-sam
```

You should be able to access the API at `http://localhost:4000/api/v1/test` & `http://localhost:4000/api/v1/users`.

## Project Structure

- **pkg**: Contains Go module-packages for utilities, lambdas, and internal logic.
- **lambdas**: Contains the main entry point for AWS Lambda functions.
- **utils**: Utility functions for password hashing, response preparation, etc.
- **internal**: Internal controllers and logic for handling API requests.
- **cdk.json**: Configuration for AWS CDK.
- **package.json**: Node.js project configuration.
- **test**: Contains Jest test files.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [AWS SDK for Go](https://github.com/aws/aws-sdk-go)
- [AWS CDK](https://github.com/aws/aws-cdk)
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) for password hashing

## Security Considerations

- **Environment Variables**: Ensure that sensitive information such as AWS credentials is stored securely and not hardcoded in the codebase. Use environment variables or AWS Secrets Manager for managing sensitive data.
- **IAM Roles**: Follow the principle of least privilege when defining IAM roles and policies for your AWS Lambda functions.
- **HTTPS**: Always use HTTPS for API Gateway endpoints to ensure secure data transmission.

## Troubleshooting

- **Deployment Issues**: If you encounter issues during deployment, ensure that your AWS CLI is configured correctly and that you have the necessary permissions to deploy resources.
- **Local Development**: If `SAM` is not working as expected, check that all dependencies are installed and that your environment variables are correctly set.
- **Testing Failures**: Ensure that your test environment is correctly configured and that all necessary mock data and dependencies are available.

## Useful Commands

- **Build Go Functions**: `make build`
- **Clean Build Artifacts**: `make clean`
- **Zip Lambda Functions**: `make zip`
- **Run Local API**: `npm run start-sam`
- **Deploy with CDK**: `npm run deploy`
- **Destroy CDK Stack**: `cdk destroy`

## Additional Resources

- [AWS Lambda Documentation](https://docs.aws.amazon.com/lambda/latest/dg/welcome.html)
- [AWS CDK Documentation](https://docs.aws.amazon.com/cdk/v2/guide/home.html)
- [Go Programming Language](https://go.dev/doc/)
- [AWS SAM CLI](https://aws.amazon.com/serverless/sam/)

## Contact

For any questions or support, please use the discussion feature of the repository or join the [Discord Server](https://discord.gg/Pb3dJPdAQr).
