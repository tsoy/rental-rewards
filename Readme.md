Build and run as a docker image
```bash
docker build --tag rental-rewards .
docker run -ti -p 9090:9090 \
  -e PORT=9090 \ 
  -e ENVIRONMENT=development \
  rental-rewards

```

Migrate installation 
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.0/migrate.linux-amd64.tar.gz | tar xvz
```


pubsub-init:
```bash
curl -X PUT "http://$(PUBSUB_EMULATOR_HOST)/v1/projects/test-project/topics/payment.completed"
curl -X PUT "http://$(PUBSUB_EMULATOR_HOST)/v1/projects/test-project/subscriptions/rewards-worker-sub" \
-H "Content-Type: application/json" \
-d '{"topic": "projects/test-project/topics/payment.completed"}'
```