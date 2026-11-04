# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/AiProvider/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_ai_provider" "acme_vllm" {
  name        = "acme-vllm"
  auth        = "bearer"
  base_url    = "https://llm.acme.example"
  protocol    = "openai-chat"
  server_kind = "openai-compatible"
}
