package com.fiberman.fiberman_java_backend.config;

import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.ViewControllerRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

@Configuration
public class SpaWebConfiguration implements WebMvcConfigurer {
    @Override
    public void addViewControllers(ViewControllerRegistry registry) {
        registry.addViewController("/").setViewName("forward:/index.html");
        registry.addViewController("/dashboard").setViewName("forward:/index.html");
        registry.addViewController("/explorer").setViewName("forward:/index.html");
        registry.addViewController("/invoice").setViewName("forward:/index.html");
        registry.addViewController("/topology").setViewName("forward:/index.html");
        registry.addViewController("/payments").setViewName("forward:/index.html");
        registry.addViewController("/logs").setViewName("forward:/index.html");
        registry.addViewController("/settings").setViewName("forward:/index.html");
    }
}
