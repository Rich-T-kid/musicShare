import logging
import os
import boto3
from typing import Optional, Any
from botocore.exceptions import ClientError

logger = logging.getLogger(__name__)

class KeyPairWrapper:
    def __init__(self, ec2_client):
        self.ec2_client = ec2_client
        self.key_pair = None
        self.key_path = None
        self.key_dir = None
    
    def create(self,key_name:str) -> dict:
        try:
            response = self.ec2_client.create_key_pair(KeyName=key_name,KeyType='rsa')
            self.key_pair = response
            self.key_file_dir = os.path.dirname(__file__)
            self.key_file_path = os.path.join(self.key_file_dir, f"{self.key_pair['KeyName']}.pem")
            with open(self.key_file_path, "w") as key_file:
                key_file.write(self.key_pair["KeyMaterial"])
            os.chmod(self.key_file_path, 0o400)  
            return self.key_pair
        except ClientError as e:
            if e.response['Error']['Code'] == 'InvalidKeyPair.Duplicate':
                print(f"A key pair named {key_name} already exists. Choose a different name.")
            else:
                print(f"Failed to create key pair: {e}")
            raise
        


    def delete(self, key_name: str) -> None:
        try:
            self.ec2_client.delete_key_pair(KeyName=key_name)
            print(f"Deleted key pair: {key_name}")
        except ClientError as e:
            print(f"Failed to delete key pair: {e}")
            raise
    
    def list_keys(self, limit: Optional[int] = None):
        try:
            response = self.ec2_client.describe_key_pairs()
            key_pairs = response.get('KeyPairs', [])
            if limit:
                key_pairs = key_pairs[:limit]
            for key_pair in key_pairs:
                print(f"Key Pair Name: {key_pair['KeyName']}, Fingerprint: {key_pair['KeyFingerprint']}")
        except ClientError as e:
            print(f"Failed to list key pairs: {e}")
            raise
