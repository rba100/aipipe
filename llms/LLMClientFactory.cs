namespace aipipe.llms;

public static class LLMClientFactory
{
    public static ILLMClient CreateClient(Config config)
    {
        if (config.UseOpenRouter && !string.IsNullOrEmpty(config.OpenRouterApiKey))
        {
            return new OpenRouterClient(config);
        }
        else if (!string.IsNullOrEmpty(config.GroqEndpoint) && !string.IsNullOrEmpty(config.GroqToken))
        {
            return new GroqClient(config);
        }
        else
        {
            throw new InvalidOperationException("Invalid configuration. Must set either GROQ_ENDPOINT/GROQ_API_KEY or OPENROUTER_API_KEY environment variables and specify --or for OpenRouter.");
        }
    }
}
