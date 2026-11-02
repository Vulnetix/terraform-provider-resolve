# Every argument the kind supports. From config-schemas examples/AiProvider/advanced.yaml.

variable "acme_vllm_api_key" {
  type      = string
  sensitive = true
  ephemeral = true
}
variable "acme_ollama_api_key" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "vulnetix_ai_provider" "acme_vllm" {
  name               = "acme-vllm"
  api_key_wo         = var.acme_vllm_api_key
  api_key_wo_version = 1
  auth               = "bearer"
  base_url           = "https://llm.acme.example"
  display_name       = "ACME vLLM cluster"
  models = [
    {
      context_window = 262144
      id             = "qwen3-coder-30b"
      images         = false
      max_tokens     = 65536
      name           = "Qwen3 Coder 30B"
    },
    {
      context_window = 131072
      id             = "qwen2.5-vl-72b"
      images         = true
      max_tokens     = 8192
      name           = "Qwen2.5 VL 72B"
    },
  ]
  protocol    = "openai-responses"
  server_kind = "openai-compatible"
  labels = {
    team = "ml-platform"
  }
  annotations = {
    runbook = "https://wiki.example.com/ml-platform/vllm"
  }
}

resource "vulnetix_ai_provider" "acme_ollama" {
  name               = "acme-ollama"
  api_key_wo         = var.acme_ollama_api_key
  api_key_wo_version = 1
  auth               = "x-api-key"
  base_url           = "https://ollama.acme.example:11434"
  display_name       = "ACME Ollama"
  models = [
    {
      context_window = 131072
      id             = "llama3.1:70b"
      max_tokens     = 0
      name           = "Llama 3.1 70B"
    },
  ]
  protocol    = "openai-chat"
  server_kind = "ollama"
}
