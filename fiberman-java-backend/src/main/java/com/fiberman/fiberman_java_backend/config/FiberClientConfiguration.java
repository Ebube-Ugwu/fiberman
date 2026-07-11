package com.fiberman.fiberman_java_backend.config;

import com.fiberman.sdk.client.FiberClient;
import java.time.Duration;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
@EnableConfigurationProperties(FiberNodeProperties.class)
public class FiberClientConfiguration {

    @Bean
    FiberClient fiberClient(FiberNodeProperties properties) {
        FiberClient.Builder builder = FiberClient.builder()
                .baseUrl(properties.url())
                .connectTimeout(Duration.ofSeconds(properties.timeoutSeconds()))
                .requestTimeout(Duration.ofSeconds(properties.timeoutSeconds()));

        if (properties.authToken() != null && !properties.authToken().isBlank()) {
            builder.authToken(properties.authToken());
        }

        return builder.build();
    }
}
