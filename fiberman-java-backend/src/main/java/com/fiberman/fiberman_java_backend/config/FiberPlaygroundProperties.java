package com.fiberman.fiberman_java_backend.config;

import jakarta.validation.constraints.NotBlank;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.validation.annotation.Validated;

@Validated
@ConfigurationProperties(prefix = "fiber.playground")
public record FiberPlaygroundProperties(@NotBlank String baseUrl) {
}
