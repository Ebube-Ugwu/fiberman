package com.fiberman.fiberman_java_backend.controller;

import com.fiberman.fiberman_java_backend.dto.PlaygroundSettingsResponse;
import com.fiberman.fiberman_java_backend.dto.UpdatePlaygroundSettingsRequest;
import com.fiberman.fiberman_java_backend.service.RuntimeFiberSettingsService;
import jakarta.validation.Valid;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/settings")
public class PlaygroundSettingsController {
    private final RuntimeFiberSettingsService runtimeFiberSettingsService;

    public PlaygroundSettingsController(RuntimeFiberSettingsService runtimeFiberSettingsService) {
        this.runtimeFiberSettingsService = runtimeFiberSettingsService;
    }

    @GetMapping
    public PlaygroundSettingsResponse getSettings() {
        return runtimeFiberSettingsService.getSettings();
    }

    @PutMapping
    public PlaygroundSettingsResponse updateSettings(@Valid @RequestBody UpdatePlaygroundSettingsRequest request) {
        return runtimeFiberSettingsService.updateSettings(request);
    }
}
