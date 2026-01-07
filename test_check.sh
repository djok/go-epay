#!/bin/bash

# ePay CHECK Test Script
# Usage: ./test_check.sh <IDN> [BASE_URL]
# Example: ./test_check.sh 12345
# Example: ./test_check.sh 12345 https://test-epay.gw.fiber.bg

IDN="${1:-12345}"
SECRET="qk9bxHVnQH2gBzvfA9vyr5pnGLHyNBgC"
BASE_URL="${2:-http://localhost:8090}"

TYPE="CHECK"

# Create message for HMAC (sorted keys alphabetically)
MESSAGE="IDN${IDN}
TYPE${TYPE}
"

# Calculate HMAC-SHA1 checksum
CHECKSUM=$(echo -n "$MESSAGE" | openssl dgst -sha1 -hmac "$SECRET" | awk '{print $2}')

echo "=== ePay CHECK Test ==="
echo "IDN: $IDN"
echo "URL: $BASE_URL"
echo "Checksum: $CHECKSUM"
echo ""

RESPONSE=$(curl -s "${BASE_URL}/v1/pay/init?TYPE=${TYPE}&IDN=${IDN}&CHECKSUM=${CHECKSUM}")

echo "Response: $RESPONSE"
echo ""

# Parse status
STATUS=$(echo "$RESPONSE" | grep -o '"STATUS":"[^"]*"' | cut -d'"' -f4)

case "$STATUS" in
    "00") echo "✓ Сметка намерена" ;;
    "14") echo "✗ Непознат абонат (IDN не съществува)" ;;
    "62") echo "○ Няма текуща сметка" ;;
    "93") echo "✗ Грешна контролна сума (CHECKSUM)" ;;
    "96") echo "✗ Обща грешка" ;;
    *)    echo "? Неизвестен статус: $STATUS" ;;
esac
