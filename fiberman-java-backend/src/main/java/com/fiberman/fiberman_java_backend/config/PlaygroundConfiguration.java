package com.fiberman.fiberman_java_backend.config;

import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Configuration;

@Configuration
@EnableConfigurationProperties(FiberPlaygroundProperties.class)
public class PlaygroundConfiguration {
}
