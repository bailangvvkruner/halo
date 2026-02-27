FROM eclipse-temurin:21-jdk AS builder

WORKDIR /build

RUN apt-get update && apt-get install -y curl ca-certificates && \
    curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g pnpm@10 && \
    rm -rf /var/lib/apt/lists/*

COPY . .

RUN ./gradlew clean downloadPluginPresets build -x check --no-daemon

WORKDIR /build/application
RUN java -Djarmode=layertools -jar build/libs/*.jar extract

################################

FROM ibm-semeru-runtimes:open-21-jre

LABEL maintainer="bailangvvking"

WORKDIR /application

COPY --from=builder /build/application/dependencies/ ./
COPY --from=builder /build/application/spring-boot-loader/ ./
COPY --from=builder /build/application/snapshot-dependencies/ ./
COPY --from=builder /build/application/application/ ./

ENV JVM_OPTS="" \
    HALO_WORK_DIR="/root/.halo2" \
    SPRING_CONFIG_LOCATION="optional:classpath:/;optional:file:/root/.halo2/" \
    TZ=Asia/Shanghai

RUN ln -sf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

EXPOSE 8090

ENTRYPOINT ["sh", "-c", "java ${JVM_OPTS} org.springframework.boot.loader.launch.JarLauncher ${0} ${@}"]
