import logging
import time
import urllib.request
import uuid
import os

import boto3
from alive_progress import alive_bar
from rich.console import Console
from typing import Any

from elastic_ip import ElasticIpWrapper
from instance import EC2InstanceWrapper
from key_pair import KeyPairWrapper
from security_groups import SecurityGroupWrapper

logger = logging.getLogger(__name__)
console = Console()


class EC2_Deployment:
    """
    Creates a key pair, security group, launches an instance, associates
    an Elastic IP, and cleans up resources.
    """

    def __init__(
        self,
        inst_wrapper: EC2InstanceWrapper,
        key_wrapper: KeyPairWrapper,
        sg_wrapper: SecurityGroupWrapper,
        eip_wrapper: ElasticIpWrapper,
        ssm_client: Any,
        remote_exec: bool = False,
    ):
        """
        Initializes the EC2InstanceScenario with the necessary AWS service wrappers.

        :param inst_wrapper: Wrapper for EC2 instance operations.
        :param key_wrapper: Wrapper for key pair operations.
        :param sg_wrapper: Wrapper for security group operations.
        :param eip_wrapper: Wrapper for Elastic IP operations.
        :param ssm_client: Boto3 client for accessing SSM to retrieve AMIs.
        :param remote_exec: Flag to indicate if the scenario is running in a remote execution
                            environment. Defaults to False. If True, the script won't prompt
                            for user interaction.
        """
        self.inst_wrapper = inst_wrapper
        self.key_wrapper = key_wrapper
        self.sg_wrapper = sg_wrapper
        self.eip_wrapper = eip_wrapper
        self.ssm_client = ssm_client
        self.remote_exec = remote_exec

    
    def create_key_pair(self):
        #You can change the name to whatever you want
        key_name = f"MyUniqueKeyPair-{uuid.uuid4().hex[:8]}"
        self.key_wrapper.create(key_name)
        return key_name

        #You can assign group name
    def create_security_group(self, vpc_id):
        group_name = f"MySecurityGroup-{uuid.uuid4().hex[:8]}"
        description = 'Description of your group!!!'
        sg_id = self.sg_wrapper.create(vpc_id, group_name, description)
        self.sg_wrapper.add_inbound_rule('tcp', (22, 22), '0.0.0.0/0')  # Sample you can change this
        return sg_id
    
    def launch_instance(self, ami_id, instance_type, key_name, security_group_id):
        return self.inst_wrapper.launch_instance(ami_id, instance_type, key_name, security_group_id)

    def save_resource_ids(self, instance_id, key_name, sg_id, allocation_id):
        resource_ids = {
            "instance_id": instance_id,
            "key_name": key_name,
            "sg_id": sg_id,
            "allocation_id": allocation_id
        }
        with open("resource_ids.txt", "w") as file:
            json.dump(resource_ids, file)
        console.print("Resource IDs saved to resource_ids.txt")

    def associate_elastic_ip(self, instance_id):
        elastic_ip = self.eip_wrapper.allocate()
        self.eip_wrapper.associate(elastic_ip['allocation_id'], instance_id)
        return elastic_ip

    def manage_instance(self, instance_id, action):
        if action == 'start':
            self.inst_wrapper.start_instance(instance_id)
        elif action == 'stop':
            self.inst_wrapper.stop_instance(instance_id)
        elif action == 'reboot':
            self.inst_wrapper.reboot_instance(instance_id)

    
    def cleanup_resources(self, instance_id, key_name, sg_id, allocation_id):
        self.inst_wrapper.terminate_instance(instance_id)
        self.key_wrapper.delete(key_name)
        self.sg_wrapper.delete(sg_id)
        self.eip_wrapper.release(allocation_id)
        if self.allocation_id:
            self.eip_wrapper.release(self.allocation_id)

def setup_deployment():
    os.environ["AWS_PROFILE"] = "AdministratorAccess-746669232870"
    region_name = 'us-east-1'
    vpc_id = 'vpc-08293ca8a1ffcef5c'

    ec2_client = boto3.client('ec2', region_name=region_name)
    ssm_client = boto3.client('ssm', region_name=region_name)

    inst_wrapper = EC2InstanceWrapper(ec2_client)
    key_wrapper = KeyPairWrapper(ec2_client)
    sg_wrapper = SecurityGroupWrapper(ec2_client)
    eip_wrapper = ElasticIpWrapper(ec2_client)

    deployment = EC2_Deployment(ec2_client, inst_wrapper, key_wrapper, sg_wrapper, eip_wrapper, ssm_client)

    key_name = deployment.create_key_pair()
    console.print(f"Key pair created: {key_name}")

    sg_id = deployment.create_security_group(vpc_id)
    console.print(f"Security group created: {sg_id}")

    ami_id = 'ami-04681163a08179f28'  # Use an appropriate AMI ID from your AWS account
    instance_type = 't2.micro'
    instance_id = deployment.launch_instance(ami_id, instance_type, key_name, sg_id)
    console.print(f"Instance launched: {instance_id}")

    allocation_id = None  # Replace with your actual allocation ID if you have one

    deployment.save_resource_ids(instance_id, key_name, sg_id, allocation_id)

    return deployment
    
if __name__ == "__main__":
    deployment = setup_deployment
    deployment.cleanup_resources()

    

    

    '''
    elastic_ip = deployment.associate_elastic_ip(instance_id)
    console.print(f"Elastic IP associated: {elastic_ip['public_ip']}")

    
    '''