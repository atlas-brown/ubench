## Build the client image
Run the following command from the `slowpoke` root directory:
```bash
docker build -f script/ClientDockerfile . -t your_tag
```

## Build the application with poker runtime
Run the following command from the `slowpoke` root directory with the `APP_NAME=benchmark_name` where `benchmark` name is one of `boutique`, `social`, `hotel`, and `movie`:
```
docker build --build-arg BENCHMARK=boutique -f scripts/PrebuiltDockerfile . -t yizhengx/ubench:boutique
```