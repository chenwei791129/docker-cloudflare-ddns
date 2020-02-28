#!/bin/bash

DATETIME=$(date +"%Y-%m-%d %H:%M:%S %z")

LASEST_IP=$(curl --silent -s ${CHECK_URL})
echo "[${DATETIME}] Checking current Public IP from ${CHECK_URL}..."

CURRENT_IP=$(cat /tmp/current_ip 2>/dev/null)

if [ "${LASEST_IP}" == "${CURRENT_IP}" ]; then
    echo "[${DATETIME}] Current public IP matches cached IP recorded. No update required! IP: ${LASEST_IP}"
    exit 0;
fi

# Get RecordID
RECORD_LIST=$(curl --silent -X GET "https://api.cloudflare.com/client/v4/zones/${CLOUDFLARE_ZONE_ID}/dns_records" \
     -H "Authorization: Bearer ${CLOUDFLARE_TOKEN}" \
     -H "Content-Type:application/json")

if [ $(echo ${RECORD_LIST} | jq '.success') == "true" ]; then
    RecordID=$(echo ${RECORD_LIST} | jq -r '.result[] | select(.name == "'${CLOUDFLARE_DOMAIN_NAME}'" and .type == "A").id')
    if [ ${RecordID} != "" ]; then
        echo "[${DATETIME}] Updating new ip: ${LASEST_IP}"
        # Update A Record
        RESULT=$(curl --silent -X PUT "https://api.cloudflare.com/client/v4/zones/${CLOUDFLARE_ZONE_ID}/dns_records/${RecordID}" \
                    -H "Authorization: Bearer ${CLOUDFLARE_TOKEN}" \
                    -H "Content-Type:application/json" \
                    --data '{"type":"A","name":"'${CLOUDFLARE_DOMAIN_NAME}'","content":"'${LASEST_IP}'","ttl":'${TTL}',"proxied":'${CLOUDFLARE_PROXIED}'}')
        if [ $(echo ${RESULT} | jq '.success') == "true" ]; then
            # Write to cache ip
            echo ${LASEST_IP} > /tmp/current_ip
            echo "[${DATETIME}] Update record successfully, new IP: ${LASEST_IP}"
            exit 0;
        else
            echo "[${DATETIME}] Update record fail!, please check your setting"
            echo "errors: $(echo ${RESULT} | jq '.errors[]')"
            exit 1;
        fi
    else
        echo "[${DATETIME}] Query record id Fail!, please check your record name and Type=A"
        exit 1;
    fi
else
    echo "[${DATETIME}] Query record list fail!, please check your zone ID or token"
    echo "errors: $(echo ${RECORD_LIST} | jq '.errors[]')"
    exit 1;
fi
