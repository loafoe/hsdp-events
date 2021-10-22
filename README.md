# hsdp-metrics-teams-webhook

Microservice helper to translate HSDP Metrics webhooks to Microsoft Teams webhooks

## Configuration

| Environment | Description |
|-------------|-------------|
| EVENTS_TOKEN | Random token to protect the endpoint |
| EVENTS_WEBHOOK_URL | Microsoft Teams Webhook URL to use |

## Deployment

```shell
cf create-app hsdp-events
cf set-env hsdp-events EVENTS_TOKEN secret
cf set-env hsdp-events EVENTS_WEBHOOK_URL https://company.webhook.office.com/webhookb2/...
cf push hsdp-events -o loafoe/hsdp-events:latest
```

## License

License is MIT
