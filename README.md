# hsdp-events

Microservice helper to translate HSDP Metrics webhooks to Microsoft Teams webhooks

## Configuration

| Environment | Description |
|-------------|-------------|
| EVENTS_TOKEN | Random token to protect the endpoint |
| EVENTS_WEBHOOK_URL | Microsoft Teams Webhook URL to use |

## Deployment

```shell
cf push hsdp-events -o loafoe/hsdp-events:latest
cf set-env hsdp-events EVENTS_TOKEN secret
cf set-env hsdp-events EVENTS_WEBHOOK_URL https://company.webhook.office.com/webhookb2/...
cf restart hsdp-events
```

## License

License is MIT
