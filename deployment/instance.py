import boto3
import logging
from botocore.exceptions import ClientError

class EC2InstanceWrapper:
    def __init__(self,ec2_client):
        self.ec2_client = ec2_client
        self.instance_id = None
    
    # Launches EC2 Instance
    def launch_instance(self, ami_id, instance_type, key_name, security_group_id):
        try:
            response = self.ec2_client.run_instances(
                ImageId=ami_id,
                InstanceType=instance_type,
                KeyName=key_name,
                SecurityGroupIds=[security_group_id],
                MinCount=1,
                MaxCount=1
            )
            self.instance_id = response['Instances'][0]['InstanceId']
            print(f"Launched EC2 Instance with ID: {self.instance_id}")
            return self.instance_id
        except ClientError as e:
            print(f"Failed to launch EC2 instance: {e}")
            raise
    
    def terminate_instance(self, instance_id):
        #Terminates the specified EC2 instance
        try:
            self.ec2_client.terminate_instances(InstanceIds=[instance_id])
            print(f"Terminated EC2 Instance with ID: {instance_id}")
        except ClientError as e:
            print(f"Failed to terminate EC2 instance: {e}")
            raise

    def start_instance(self, instance_id):
        try:
            response = self.ec2_client.start_instances(InstanceIds=[instance_id])
            print(f"Started EC2 Instance with ID: {instance_id}")
            return response
        except ClientError as e:
            print(f"Failed to start EC2 instance: {e}")
            raise

    def stop_instance(self, instance_id):
        try:
            response = self.ec2_client.stop_instances(InstanceIds=[instance_id])
            print(f"Stopped EC2 Instance with ID: {instance_id}")
            return response
        except ClientError as e:
            print(f"Failed to stop EC2 instance: {e}")
            raise

    def reboot_instance(self, instance_id):
        try:
            response = self.ec2_client.reboot_instances(InstanceIds=[instance_id])
            print(f"Rebooted EC2 Instance with ID: {instance_id}")
            return response
        except ClientError as e:
            print(f"Failed to reboot EC2 instance: {e}")
            raise
            