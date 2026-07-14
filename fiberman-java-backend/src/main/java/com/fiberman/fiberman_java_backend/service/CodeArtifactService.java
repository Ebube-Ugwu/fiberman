package com.fiberman.fiberman_java_backend.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fiberman.fiberman_java_backend.dto.CodeArtifacts;
import java.util.Map;
import org.springframework.stereotype.Service;

@Service
public class CodeArtifactService {
    private final ObjectMapper objectMapper;
    private final RuntimeFiberSettingsService runtimeFiberSettingsService;

    public CodeArtifactService(
            ObjectMapper objectMapper,
            RuntimeFiberSettingsService runtimeFiberSettingsService
    ) {
        this.objectMapper = objectMapper;
        this.runtimeFiberSettingsService = runtimeFiberSettingsService;
    }

    public CodeArtifacts generate(String backendPath, String sdkMethodCall, Object requestBody) {
        return new CodeArtifacts(
                buildCurl(backendPath, requestBody),
                buildJavaSnippet(sdkMethodCall),
                buildGoSnippet(sdkMethodCall));
    }

    private String buildCurl(String backendPath, Object requestBody) {
        StringBuilder builder = new StringBuilder()
                .append("curl -X ")
                .append(requestBody == null ? "GET" : "POST")
                .append(" \"")
                .append(runtimeFiberSettingsService.playgroundBaseUrl())
                .append(backendPath)
                .append("\"");

        if (requestBody != null) {
            builder.append(" \\\n  -H \"Content-Type: application/json\"")
                    .append(" \\\n  -d '")
                    .append(escapeSingleQuotes(toJson(requestBody)))
                    .append("'");
        }

        return builder.toString();
    }

    private String buildJavaSnippet(String sdkMethodCall) {
        String authToken = runtimeFiberSettingsService.authToken();
        String authLine = authToken == null || authToken.isBlank()
                ? ""
                : "    .authToken(\"" + escapeJava(authToken) + "\")\n";

        return """
                FiberClient client = FiberClient.builder()
                    .baseUrl("%s")
                %s    .build();

                var response = %s;
                """.formatted(
                escapeJava(runtimeFiberSettingsService.nodeUrl()),
                authLine,
                sdkMethodCall);
    }

    private String buildGoSnippet(String sdkMethodCall) {
        String authToken = runtimeFiberSettingsService.authToken();
        String authLine = authToken == null || authToken.isBlank()
                ? ""
                : "    AuthToken: \"" + escapeJava(authToken) + "\",\n";

        return """
                sdk, err := client.New(client.Config{
                    BaseURL: "%s",
                %s})
                if err != nil {
                    log.Fatal(err)
                }

                response, err := %s
                if err != nil {
                    log.Fatal(err)
                }
                """.formatted(
                escapeJava(runtimeFiberSettingsService.nodeUrl()),
                authLine,
                sdkMethodCall.replace("client.", "sdk."));
    }

    private String toJson(Object requestBody) {
        try {
            return objectMapper.writeValueAsString(requestBody);
        } catch (JsonProcessingException exception) {
            throw new IllegalStateException("Failed to serialize request body for code generation", exception);
        }
    }

    private String escapeSingleQuotes(String value) {
        return value.replace("'", "'\"'\"'");
    }

    private String escapeJava(String value) {
        return value.replace("\\", "\\\\").replace("\"", "\\\"");
    }
}
