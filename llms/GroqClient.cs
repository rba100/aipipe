using OpenAI;
using OpenAI.Chat;
using System;
using System.Threading.Tasks;

namespace aipipe.llms;

public class GroqClient : ILLMClient
{
    private readonly string _groqToken;
    private readonly string _groqEndpoint;
    private readonly Config _config;

    public GroqClient(Config config)
    {
        _groqEndpoint = config.GroqEndpoint!;
        _groqToken = config.GroqToken!;
        _config = config;
    }

    public async Task<string> CreateCompletionAsync(string prompt)
    {
        string groqModel = _config.GroqDefaultModel;
        if (_config.ModelType == ModelType.Fast) groqModel =  _config.GroqFastModel;
        if (_config.ModelType == ModelType.Reasoning) groqModel = _config.GroqReasoningModel;

        ChatClient client = new(model: groqModel, credential: _groqToken!, new OpenAIClientOptions
        {
            Endpoint = new Uri(_groqEndpoint!),
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
        string groqModel = _config.GroqDefaultModel;
        if (_config.ModelType == ModelType.Fast) groqModel =  _config.GroqFastModel;
        if (_config.ModelType == ModelType.Reasoning) groqModel = _config.GroqReasoningModel;

        ChatClient client = new(model: groqModel, credential: _groqToken!, new OpenAIClientOptions
        {
            Endpoint = new Uri(_groqEndpoint!),
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
            yield return string.Join("", thing.ContentUpdate.Select(x => x.Text));
        }
    }
}
