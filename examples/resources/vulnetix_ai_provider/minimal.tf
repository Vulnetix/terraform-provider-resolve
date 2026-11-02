# Only the arguments the kind requires. From config-schemas examples/AiProvider/minimal.yaml.

resource "vulnetix_ai_provider" "acme_vllm" {
  name     = "acme-vllm"
  auth     = "bearer"
  base_url = "https://llm.acme.example"
  protocol = "openai-chat"
}
