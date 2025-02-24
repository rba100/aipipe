namespace aipipe.Llms;

public static class LLMClientFactory
{
    private static string GetModelForConfig(Config config, string defaultModel, string fastModel, string reasoningModel)
    {
        return config.ModelType switch
        {
            ModelType.Fast => fastModel,
            ModelType.Reasoning => reasoningModel,
            _ => defaultModel
        };
    }

    public static ILLMClient CreateClient(Config config)
    {
        if (config.UseOpenRouter)        
        {
            if(string.IsNullOrEmpty(config.OpenRouterApiKey))
            {
                throw new InvalidOperationException("Invalid configuration. Must set OPENROUTER_API_KEY environment variable.");
            }

            string model = GetModelForConfig(
                config,
                config.OpenRouterDefaultModel,
                config.OpenRouterFastModel,
                config.OpenRouterReasoningModel
            );

            return new GenericClient(
                config.OpenRouterApiKey,
                "https://openrouter.ai/api/v1/",
                config,
                model
            );
        }
        else if (!string.IsNullOrEmpty(config.GroqEndpoint) && !string.IsNullOrEmpty(config.GroqToken))
        {
            string model = GetModelForConfig(
                config,
                config.GroqDefaultModel,
                config.GroqFastModel,
                config.GroqReasoningModel
            );

            return new GenericClient(
                config.GroqToken,
                config.GroqEndpoint,
                config,
                model
            );
        }
        else
        {
            throw new InvalidOperationException("Invalid configuration. Must set either GROQ_ENDPOINT/GROQ_API_KEY or OPENROUTER_API_KEY environment variables and specify --or for OpenRouter.");
        }
    }
}
