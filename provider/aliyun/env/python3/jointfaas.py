"""
Aliyun python function handler wrapper

"""
import index

def handler(event, context):
    return index.handler(event)