Build and run as a docker image
```bash
docker build --tag rental-rewards .
docker run -ti -p 9090:9090 \
  -e PORT=9090 \
  -e ENVIRONMENT=development \
  rental-rewards

```