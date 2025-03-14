using OpenAI;
using OpenAI.Chat;
using System;
using System.ClientModel;
using System.Threading.Tasks;

namespace aipipe.Llms;

public class GenericClient : ILLMClient
{
    private readonly string _apiKey;
    private readonly string _baseUrl;
    private readonly Config _config;
    private readonly string _model;

    public GenericClient(string apiKey, string baseUrl, Config config, string model)
    {
        _apiKey = apiKey;
        _baseUrl = baseUrl;
        _config = config;
        _model = model;
    }

    public async Task<string> CreateCompletionAsync(string prompt)
    {
        ChatClient client = new(model: _model, credential: new ApiKeyCredential(_apiKey), new OpenAIClientOptions
        {
            Endpoint = new Uri(_baseUrl),
        });

        var systemMessage = Prompts.GetSystemPrompt(_config.IsCodeBlock);
        var options = new ChatCompletionOptions { };
        var messages = new ChatMessage[]
        {
            new SystemChatMessage(systemMessage),
            new UserChatMessage(prompt),
        };

        var response = await client.CompleteChatAsync(messages, options, CancellationToken.None);
        return response.Value.Content.Single().Text;
    }

    public async IAsyncEnumerable<string> CreateCompletionStreamAsync(string prompt)
    {
        ChatClient client = new(model: _model, credential: new ApiKeyCredential(_apiKey), new OpenAIClientOptions
        {
            Endpoint = new Uri(_baseUrl),
        });

        var systemMessage = Prompts.GetSystemPrompt(_config.IsCodeBlock);
        var options = new ChatCompletionOptions { };
        var messages = new ChatMessage[]
        {
            new SystemChatMessage(systemMessage),
            new UserChatMessage(prompt),
        };

        await foreach(var thing in client.CompleteChatStreamingAsync(messages, options, CancellationToken.None))
        {
            var textUpdates = thing.ContentUpdate.Where(u=>u.Kind == ChatMessageContentPartKind.Text);
            yield return string.Join("", textUpdates.Select(x => x.Text));
        }
    }
}