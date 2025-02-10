import logging
from typing import Any, List
import boto3
from botocore.exceptions import ClientError

logger = logging.getLogger(__name__)

class ElasticIpWrapper:
    
    def __init__(self, ec2_client: Any) -> None:
        self.ec2_client = ec2_client
        self.elastic_ips: List = []

    def allocate(self):
        try:
            allocation = self.ec2.allocate_address(Domain='vpc')
            elastic_ip = {
                'allocation_id' : allocation['AllocationId'],
                'public_ip' : allocation['Public_Ip'],
                'instance_id' : None
            }
            self.elastic_ips.append(elastic_ip)
            return elastic_ip

        except ClientError as e:
            print(f"Failed to allocate Elastic IP: {e}")
            raise
    
    def associate(self, allocation_id, instance_id):
        try:
            response = self.ec2.associate_address(AllocationId=allocation['AllocationId'], InstanceId='INSTANCE_ID')
            for ip in self.elastic_ips:
                if ip['allocation_id'] == allocation_id:
                    ip['instance_id'] = instance_id
            return response
        except ClientError as e:
            print(f"Failed to associate Elastic IP: {e}")
            raise

    def disassociate(self, allocation_id):
        try:
            response = self.ec2.describe_addresses(AllocationIds = ['allocation_id'])
            association_id = response['Addresses'][0].get('AssociationId')
            if association_id:
                self.ec2_client.disassociate_address(AssociationId=association_id)
                for ip in self.elastic_ips:
                    if ip['allocation_id'] == allocation_id:
                        ip['instance_id'] = None
        except ClientError as e:
            print(f"Failed to disassociate Elastic IP: {e}")
            raise
    
    def release(self, allocation_id):
        try:
            self.ec2_client.release_address(AllocationId=allocation_id)
            self.elastic_ips = [ip for ip in self.elastic_ips if ip['allocation_id'] != allocation_id]
        except ClientError as e:
            print(f"Failed to release Elastic IP: {e}")
            raise


