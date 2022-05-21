#!/bin/bash


CFN_TEMPLATE=dynamo_ddl.yaml
CFN_STACK_NAME=sc-ddl-deploy

aws cloudformation deploy --stack-name ${CFN_STACK_NAME} --template-file ${CFN_TEMPLATE} \
--role-arn ${AWS_CFN_ARN}