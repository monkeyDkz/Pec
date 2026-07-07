# Scripts

| Script | Description |
|---|---|
| `load-test.js` | k6 load test (TICK-074). Runs three scenarios: browse, 100 simultaneous listeners, and a login burst that proves the rate limiter triggers. |

## Run the load test

```bash
# Local
k6 run scripts/load-test.js

# Against staging
API_URL=https://staging.streampulse.example k6 run scripts/load-test.js
```

The test exits non-zero if any SLO is violated:

- `GET /streams` p95 < 500 ms
- Global error rate < 1 %
- 100 simultaneous listeners can subscribe without timeouts
