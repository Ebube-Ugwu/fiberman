package com.fiberman.fiberman_java_backend.service;

import com.fiberman.fiberman_java_backend.config.FiberNodeProperties;
import com.fiberman.fiberman_java_backend.config.FiberPlaygroundProperties;
import com.fiberman.fiberman_java_backend.dto.PlaygroundSettingsResponse;
import com.fiberman.fiberman_java_backend.dto.UpdatePlaygroundSettingsRequest;
import com.fiberman.sdk.client.FiberClient;
import java.time.Duration;
import org.springframework.stereotype.Service;

@Service
public class RuntimeFiberSettingsService {
    private final FiberPlaygroundProperties playgroundProperties;

    private volatile String nodeUrl;
    private volatile String authToken;
    private volatile long timeoutSeconds;
    private volatile String defaultInvoiceCurrency;

    public RuntimeFiberSettingsService(
            FiberNodeProperties nodeProperties,
            FiberPlaygroundProperties playgroundProperties
    ) {
        this.playgroundProperties = playgroundProperties;
        this.nodeUrl = nodeProperties.url();
        this.authToken = blankToNull(nodeProperties.authToken());
        this.timeoutSeconds = nodeProperties.timeoutSeconds();
        this.defaultInvoiceCurrency = "";
    }

    public FiberClient createClient() {
        FiberClient.Builder builder = FiberClient.builder()
                .baseUrl(nodeUrl)
                .connectTimeout(Duration.ofSeconds(timeoutSeconds))
                .requestTimeout(Duration.ofSeconds(timeoutSeconds));

        if (authToken != null && !authToken.isBlank()) {
            builder.authToken(authToken);
        }

        return builder.build();
    }

    public PlaygroundSettingsResponse getSettings() {
        return new PlaygroundSettingsResponse(
                nodeUrl,
                authToken == null ? "" : authToken,
                timeoutSeconds,
                defaultInvoiceCurrency,
                playgroundProperties.baseUrl());
    }

    public PlaygroundSettingsResponse updateSettings(UpdatePlaygroundSettingsRequest request) {
        this.nodeUrl = request.nodeUrl().trim();
        this.authToken = blankToNull(request.authToken());
        this.timeoutSeconds = request.timeoutSeconds();
        this.defaultInvoiceCurrency = request.defaultInvoiceCurrency() == null ? "" : request.defaultInvoiceCurrency().trim();
        return getSettings();
    }

    public String nodeUrl() {
        return nodeUrl;
    }

    public String authToken() {
        return authToken;
    }

    public String defaultInvoiceCurrency() {
        return defaultInvoiceCurrency;
    }

    public String playgroundBaseUrl() {
        return playgroundProperties.baseUrl();
    }

    private String blankToNull(String value) {
        return value == null || value.isBlank() ? null : value.trim();
    }
}
