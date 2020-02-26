# Prisma Cloud Remediation
## Dispatcher
Dispatcher integrate with AWS SQS to receive message from Prisma.
Dispatcher invoke the corresponding receiver (lambda function) with
the SQS Message to perform the remediation
```bash
.
├── Makefile                    <-- Make to automate build
├── README.md                   <-- This instructions file
├── dispatcher                  <-- Source code for a lambda function
│   └── main.go                 <-- Lambda function code
│   └── README.md               <-- dispatcher instruction file
└── template.yaml
```
