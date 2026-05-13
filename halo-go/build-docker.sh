#!/bin/bash
set -euo pipefail

IMAGE_NAME="${IMAGE_NAME:-halohub/halo-go}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

echo "=========================================="
echo "  Halo Go Docker 镜像构建"
echo "  镜像: ${IMAGE_NAME}:${IMAGE_TAG}"
echo "=========================================="

docker build \
  --platform linux/amd64,linux/arm64 \
  --tag "${IMAGE_NAME}:${IMAGE_TAG}" \
  --file Dockerfile \
  .

echo ""
echo "✅ 镜像构建完成: ${IMAGE_NAME}:${IMAGE_TAG}"
echo ""
echo "运行方式（与 Java Halo 参数完全一致）:"
echo ""
echo "  # 基本运行"
echo "  docker run -d --name halo-go --restart always \\"
echo "    -p 8090:8090 \\"
echo "    -v ~/.halo2:/root/.halo2 \\"
echo "    ${IMAGE_NAME}:${IMAGE_TAG}"
echo ""
echo "  # 自定义端口"
echo "  docker run -d --name halo-go --restart always \\"
echo "    -p 8080:8080 \\"
echo "    -v ~/.halo2:/root/.halo2 \\"
echo "    ${IMAGE_NAME}:${IMAGE_TAG} --server.port=8080"
echo ""
echo "  # 自定义工作目录"
echo "  docker run -d --name halo-go --restart always \\"
echo "    -p 8090:8090 \\"
echo "    -v /data/halo:/data/halo \\"
echo "    ${IMAGE_NAME}:${IMAGE_TAG} --halo.work-dir=/data/halo"
echo ""
echo "=========================================="