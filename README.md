## go-epay
A generic ePay integration in Go

### Supported Billing Systems
 * **TelcoNG** - Integration with the online payment processor of TelcoNG
 * **UCRM/UISP** - Integration with Ubiquiti's UCRM/UISP billing system

### Deployment Options

#### Docker (Recommended)
```bash
# Copy and configure environment
cp .env.example .env
# Edit .env with your settings

# Start the service
docker-compose up -d
```

#### Google App Engine
Deploy to GAE with configuration stored in Datastore.

### Environment Variables

| Variable | Description |
|----------|-------------|
| `BILLING_SYSTEM` | `telcong` or `ucrm` - which billing system to use |
| `EPAY_SECRET` | ePay HMAC secret for request validation |
| `EPAY_MERCHANT_ID` | ePay merchant ID |
| `TELCONG_BILLING_URL` | TelcoNG API URL (if using TelcoNG) |
| `TELCONG_JWT_KEY` | TelcoNG JWT key JSON (if using TelcoNG) |
| `UCRM_BILLING_URL` | UCRM/UISP API URL (if using UCRM) |
| `UCRM_API_KEY` | UCRM/UISP API key (if using UCRM) |
| `UCRM_METHOD_ID` | UCRM payment method ID |
| `UCRM_PROVIDER_NAME` | Provider name for UCRM payments |
| `UCRM_ORGANIZATION_ID` | UCRM organization ID (optional) |

### UCRM/UISP Client Lookup

The system supports two types of subscriber identifiers (IDN):

| IDN Format | Search Method | API Parameter |
|------------|---------------|---------------|
| 7 digits with valid Luhn checksum | Contract ID | `?query=IDN` |
| 7 digits with invalid Luhn | Error: not found | - |
| Other lengths | User ID | `?userIdent=IDN` |

### Architecture
![Architecture](docs/architecture.png)

### Building

```bash
# Build
go build ./cmd/goepay

# Run tests
go test ./...
```

### Requirements
 * Go 1.19 or greater

### License
Copyright 2018 clouWay ood.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   https://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

