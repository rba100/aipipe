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

        var systemMessage = _config.IsCodeBlock
            ? "You are a helpful assistant. If the user has asked for something written, put it in a code block (```), otherwise just provide the answer."
              + " If you do use a codeblock, all other text is ignored."
            : "You are a helpful assistant.";

        var options = new ChatCompletionOptions { };

        var messages = new ChatMessage[]
        {
            new SystemChatMessage(systemMessage),
            new UserChatMessage(prompt),
        };

        var response = await client.CompleteChatAsync(messages, options, CancellationToken.None);
        return response.Value.Content.Single().Text;
    }
}
