using System;
using aipipe.llms;
using YamlDotNet.Serialization;
using YamlDotNet.Serialization.NamingConventions;
using System.IO;

namespace aipipe;

public class Config
{
    public string? GroqEndpoint { get; set; }
    public string? GroqToken { get; set; }
    public string? OpenRouterApiKey { get; set; }
    public bool UseOpenRouter { get; set; }
    public string GroqDefaultModel { get; set; } = "llama-3.3-70b-versatile";
    public string OpenRouterDefaultModel { get; set; } = "google/gemini-2.0-flash-001";
    public string OpenRouterFastModel { get; set; } = "meta-llama/llama-3-8b-instruct";
    public string OpenRouterReasoningModel { get; set; } = "deepseek/deepseek-r1-distill-llama-70b:free";
    public string GroqFastModel { get; set; } = "llama-3.1-8b-instant";
    public string GroqReasoningModel { get; set; } = "deepseek-r1-distill-llama-70b";
    public ModelType ModelType { get; set; } = ModelType.Default;
    public bool IsCodeBlock { get; set; }
    public bool IsMic { get; set; }

    public Config()
    {
        GroqEndpoint = Environment.GetEnvironmentVariable("GROQ_ENDPOINT") ?? "https://api.groq.com/openai/v1";
        GroqToken = Environment.GetEnvironmentVariable("GROQ_API_KEY");
        OpenRouterApiKey = Environment.GetEnvironmentVariable("OPENROUTER_API_KEY");

        // Load from YAML file
        var homeDir = Environment.GetFolderPath(Environment.SpecialFolder.UserProfile);
        var configDir = Path.Combine(homeDir, ".aipipe");
        var configFile = Path.Combine(configDir, "config.yaml");

        if (File.Exists(configFile))
        {
            try
            {
                var deserializer = new DeserializerBuilder()
                    .WithNamingConvention(CamelCaseNamingConvention.Instance)
                    .Build();

                using (var reader = new StreamReader(configFile))
                {
                    var yamlConfig = deserializer.Deserialize<Config>(reader);

                    UseOpenRouter = UseOpenRouter || yamlConfig.UseOpenRouter;
                    OpenRouterApiKey ??= yamlConfig.OpenRouterApiKey;
                    OpenRouterDefaultModel = yamlConfig.OpenRouterDefaultModel ?? OpenRouterDefaultModel;
                    OpenRouterFastModel = yamlConfig.OpenRouterFastModel ?? OpenRouterFastModel;
                    OpenRouterReasoningModel = yamlConfig.OpenRouterReasoningModel ?? OpenRouterReasoningModel;

                    GroqEndpoint ??= yamlConfig.GroqEndpoint;
                    GroqToken ??= yamlConfig.GroqToken;
                    GroqDefaultModel = yamlConfig.GroqDefaultModel ?? GroqDefaultModel;
                    GroqFastModel = yamlConfig.GroqFastModel ?? GroqFastModel;
                    GroqReasoningModel = yamlConfig.GroqReasoningModel ?? GroqReasoningModel;
                }
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"Error loading config file: {ex.Message}");
            }
        }
    }
}
