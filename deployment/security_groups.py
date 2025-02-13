import boto3
from botocore.exceptions import ClientError

class SecurityGroupWrapper:
    def __init__(self, ec2_client):
        self.ec2_client = ec2_client
        self.security_group_id = None

    # Creates Security Group
    def create(self, vpc_id, group_name, description):
        try:
            response = self.ec2_client.create_security_group(
                GroupName=group_name,
                Description=description,
                VpcId=vpc_id
            )
            self.security_group_id = response['GroupId']
            print(f"Created Security Group {group_name} with ID {self.security_group_id}")
            return self.security_group_id
        except ClientError as e:
            print(f"Failed to create security group: {e}")
            raise

    def add_inbound_rule(self, protocol, port_range, cidr_block):
        try:
            self.ec2_client.authorize_security_group_ingress(
                GroupId=self.security_group_id,
                IpPermissions=[
                    {
                        'IpProtocol': protocol,
                        'FromPort': port_range[0],
                        'ToPort': port_range[1],
                        'IpRanges': [{'CidrIp': cidr_block}]
                    }
                ]
            )
            print(f"Added inbound rule to Security Group {self.security_group_id}")
        except ClientError as e:
            print(f"Failed to add inbound rule: {e}")
            raise



    def delete(self, group_id):
        try:
            self.ec2_client.delete_security_group(GroupId=group_id)
            print(f"Deleted security group: {group_id}")
        except ClientError as e:
            print(f"Failed to delete security group: {e}")
            