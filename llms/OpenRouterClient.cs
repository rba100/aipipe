using System;
using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Text.Json.Serialization.Metadata;
using System.Threading.Tasks;

namespace aipipe.llms;

public class OpenRouterClient : ILLMClient
{
    private readonly string _openRouterApiKey;
    private readonly Config _config;
    private static readonly HttpClient _httpClient = new HttpClient();

    public OpenRouterClient(Config config)
    {
        _openRouterApiKey = config.OpenRouterApiKey!;
        _config = config;
        _httpClient.BaseAddress = new Uri("https://openrouter.ai/api/v1/");
        _httpClient.DefaultRequestHeaders.Add("Authorization", $"Bearer {_openRouterApiKey}");
    }

    public async Task<string> CreateCompletionAsync(string prompt)
    {
        string model = _config.OpenRouterDefaultModel;
        if (_config.ModelType == ModelType.Fast) model = _config.OpenRouterFastModel;
        if (_config.ModelType == ModelType.Reasoning) model = _config.OpenRouterReasoningModel;

        var systemMessage = Prompts.GetSystemPrompt(_config.IsCodeBlock);

        var requestBody = new ChatRequest(model, new[]
        {
            new Message("system", systemMessage),
            new Message("user", prompt)
        });

        var json = JsonSerializer.Serialize(requestBody, SourceGenerationContext.Default.ChatRequest);
        var content = new StringContent(json, Encoding.UTF8, "application/json");

        var response = await _httpClient.PostAsync("chat/completions", content);
        response.EnsureSuccessStatusCode();

        var responseJson = await response.Content.ReadAsStringAsync();
        using JsonDocument doc = JsonDocument.Parse(responseJson);
        JsonElement root = doc.RootElement;
        JsonElement choices = root.GetProperty("choices");
        JsonElement firstChoice = choices[0];
        JsonElement message = firstChoice.GetProperty("message");
        string aiOutput = message.GetProperty("content").GetString()!;

        return aiOutput;
    }

    public IAsyncEnumerable<string> CreateCompletionStreamAsync(string prompt)
    {
        throw new NotImplementedException();
    }

    public record ChatRequest(string? model, Message[]? messages);

    public record Message(string? role, string? content);
}

[System.Text.Json.Serialization.JsonSourceGenerationOptions(WriteIndented = true)]
[JsonSerializable(typeof(OpenRouterClient.ChatRequest))]
internal partial class SourceGenerationContext : JsonSerializerContext
{
}
