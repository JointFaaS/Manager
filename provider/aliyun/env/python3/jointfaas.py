"""
Aliyun python function handler wrapper

"""
from . import index

def handler(event, context):
    return index.handler(event)