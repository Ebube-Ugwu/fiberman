FROM docker.io/nervos/fiber:0.9.0-rc7 AS fiber-runtime

FROM docker.io/library/maven:3.9.11-eclipse-temurin-21 AS sdk-build
WORKDIR /workspace
COPY fiber-java-sdk ./fiber-java-sdk
RUN mvn -f fiber-java-sdk/pom.xml -DskipTests install

FROM docker.io/library/node:22-bookworm-slim AS frontend-build
WORKDIR /workspace/fiberman-frontend
COPY fiberman-frontend/package*.json ./
RUN npm ci
COPY fiberman-frontend ./
RUN npm run build

FROM docker.io/library/eclipse-temurin:21-jdk AS backend-build
WORKDIR /workspace
COPY --from=sdk-build /root/.m2 /root/.m2
COPY fiberman-java-backend/gradlew ./fiberman-java-backend/gradlew
COPY fiberman-java-backend/gradle ./fiberman-java-backend/gradle
COPY fiberman-java-backend/build.gradle ./fiberman-java-backend/build.gradle
COPY fiberman-java-backend/settings.gradle ./fiberman-java-backend/settings.gradle
COPY fiberman-java-backend/src ./fiberman-java-backend/src
COPY --from=frontend-build /workspace/fiberman-frontend/dist/fiberman-frontend/browser/ ./fiberman-java-backend/src/main/resources/static/
WORKDIR /workspace/fiberman-java-backend
RUN chmod +x ./gradlew && GRADLE_USER_HOME=/tmp/gradle-home ./gradlew bootJar --no-daemon

FROM docker.io/library/eclipse-temurin:21-jre
WORKDIR /app
COPY --from=fiber-runtime /usr/bin/tini /usr/bin/tini
COPY --from=fiber-runtime /usr/local/bin/fnn /usr/local/bin/fnn
COPY --from=fiber-runtime /usr/local/bin/fnn-cli /usr/local/bin/fnn-cli
COPY --from=fiber-runtime /usr/local/share/fiber /usr/local/share/fiber
COPY --from=backend-build /workspace/fiberman-java-backend/build/libs/*.jar /app/fiberman.jar
COPY deploy/fiberman-container-entrypoint.sh /usr/local/bin/fiberman-container-entrypoint.sh
RUN chmod +x /usr/local/bin/fiberman-container-entrypoint.sh
VOLUME ["/fiber"]
EXPOSE 9010 8228
ENV FIBER_HOME=/fiber
ENV FIBER_CONFIG_TEMPLATE=/usr/local/share/fiber/config/testnet/config.yml
ENV SERVER_PORT=9010
ENV FIBER_NODE_URL=http://127.0.0.1:8227
ENV FIBER_NODE_AUTH_TOKEN=
ENV FIBER_NODE_TIMEOUT_SECONDS=30
ENV FIBER_SECRET_KEY_PASSWORD=fiberman-demo-password
ENV FIBER_RUST_LOG=info
ENV FIBER_PLAYGROUND_BASE_URL=http://localhost:9010
ENTRYPOINT ["/usr/bin/tini", "--", "/usr/local/bin/fiberman-container-entrypoint.sh"]
