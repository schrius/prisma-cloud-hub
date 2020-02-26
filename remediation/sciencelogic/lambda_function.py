import requests
import json
import os

def lambda_handler(event, context):
    print("--- Prisma High Alerts ---")
    send_alert(event)
    
def send_alert(msg):
    if 'hasFinding' in msg: del msg['hasFinding']
    if 'alertRemediationCli' in msg: del msg['alertRemediationCli']
    if 'source' in msg: del msg['source']
    if 'complianceMetadata' in msg: del msg['complianceMetadata']
    if 'policyLabels' in msg: del msg['policyLabels']
    if 'resource' in msg: del msg['resource']
    if 'resourceName' in msg: del msg['resourceName']
    if 'alertAttribution' in msg: del msg['alertAttribution']
    if 'riskRating' in msg: del msg['riskRating']
    if 'resourceRegion' in msg: del msg['resourceRegion']
    if 'policyDescription' in msg: del msg['policyDescription']
    if 'policyRecommendation' in msg: del msg['policyRecommendation']
    if 'accountId' in msg: del msg['accountId']
    if 'resourceConfig' in msg: del msg['resourceConfig']
    if 'policyId' in msg: del msg['policyId']
    if 'resourceCloudService' in msg: del msg['resourceCloudService']
    if 'alertTs' in msg: del msg['alertTs']
    if 'findingSummary' in msg: del msg['findingSummary']
    if 'hasFinresourceTypeding' in msg: del msg['hasFinresourceTypeding']
    
    payload = {}
    payload = {
        "force_ytype":"0",
        "force_yid":"0",
        "force_yname":"",
        "message":json.dumps(msg),
        "value":"0",
        "threshold":"0",
        "message_time":"0",
        "aligned_resource": "/device/{}".format(os.environ['SL_DEVICE_ID'])
    }
    
    headers = {
        'content-type': "application/json",
        'authorization': "Basic {}".format(os.environ['SL_AUTH'])
    }

    print(json.dumps(payload))
    print("--------")
    print(json.dumps(msg))
    
    url = "https://{}/api/alert".format(os.environ['SL_SERVER'])
    response = requests.request("POST", url, data=json.dumps(payload), headers=headers, verify=False)
    print(response)
    

    