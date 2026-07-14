package fiberman

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CodeArtifactService struct {
	settings *RuntimeSettings
}

func NewCodeArtifactService(settings *RuntimeSettings) *CodeArtifactService {
	return &CodeArtifactService{settings: settings}
}

func (s *CodeArtifactService) Generate(backendPath string, goMethodCall string, requestBody any) *CodeArtifacts {
	return &CodeArtifacts{
		Curl:        s.buildCurl(backendPath, requestBody),
		JavaSnippet: s.buildJavaSnippet(goMethodCall),
		GoSnippet:   s.buildGoSnippet(goMethodCall),
	}
}

func (s *CodeArtifactService) buildCurl(backendPath string, requestBody any) string {
	method := "GET"
	builder := strings.Builder{}
	if requestBody != nil {
		method = "POST"
	}

	builder.WriteString("curl -X ")
	builder.WriteString(method)
	builder.WriteString(" \"")
	builder.WriteString(s.settings.PlaygroundBaseURL())
	builder.WriteString(backendPath)
	builder.WriteString("\"")

	if requestBody != nil {
		payload, _ := json.Marshal(requestBody)
		builder.WriteString(" \\\n  -H \"Content-Type: application/json\"")
		builder.WriteString(" \\\n  -d '")
		builder.WriteString(strings.ReplaceAll(string(payload), "'", "'\"'\"'"))
		builder.WriteString("'")
	}

	return builder.String()
}

func (s *CodeArtifactService) buildGoSnippet(goMethodCall string) string {
	authLine := ""
	if authToken := s.settings.AuthToken(); authToken != "" {
		authLine = fmt.Sprintf("    AuthToken: %q,\n", authToken)
	}

	return fmt.Sprintf(`sdk, err := client.New(client.Config{
    BaseURL: %q,
%s})
if err != nil {
    log.Fatal(err)
}

response, err := %s
if err != nil {
    log.Fatal(err)
}
	`, s.settings.NodeURL(), authLine, goMethodCall)
}

func (s *CodeArtifactService) buildJavaSnippet(goMethodCall string) string {
	authLine := ""
	if authToken := s.settings.AuthToken(); authToken != "" {
		authLine = fmt.Sprintf("    .authToken(%q)\n", authToken)
	}

	javaCall := strings.ReplaceAll(goMethodCall, "sdk.", "client.")
	javaCall = strings.ReplaceAll(javaCall, "model.", "")
	javaCall = strings.ReplaceAll(javaCall, "stringPtr(", "")
	javaCall = strings.ReplaceAll(javaCall, "int64Ptr(", "")
	javaCall = strings.ReplaceAll(javaCall, "int64PtrOrNil(", "")
	javaCall = strings.ReplaceAll(javaCall, ")", ")")

	return fmt.Sprintf(`FiberClient client = FiberClient.builder()
    .baseUrl(%q)
%s    .build();

var response = %s;
`, s.settings.NodeURL(), authLine, javaCall)
}
