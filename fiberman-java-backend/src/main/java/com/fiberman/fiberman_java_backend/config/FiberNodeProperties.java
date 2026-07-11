package com.fiberman.fiberman_java_backend.config;

import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.validation.annotation.Validated;

@Validated
@ConfigurationProperties(prefix = "fiber.node")
public record FiberNodeProperties(
        @NotBlank String url,
        String authToken,
        @Min(1) long timeoutSeconds
) {
}
