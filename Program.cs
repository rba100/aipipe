using System;
using System.IO;
using System.Text;
using System.Text.RegularExpressions;
using System.CommandLine;
using System.CommandLine.Invocation;

using aipipe.Llms;
using aipipe.Speech;

namespace aipipe;

class Program
{
    static async Task<int> Main(string[] args)
    {
        var options = new CommandLineOptions();
        var rootCommand = options.RootCommand;

        rootCommand.SetHandler(async (InvocationContext context) =>
        {
            bool isStream = context.ParseResult.GetValueForOption(options.StreamOption);
            bool isCodeBlock = context.ParseResult.GetValueForOption(options.CodeBlockOption);
            bool isReasoning = context.ParseResult.GetValueForOption(options.ReasoningOption);
            bool fast = context.ParseResult.GetValueForOption(options.FastOption);
            bool mic = context.ParseResult.GetValueForOption(options.MicOption);
            bool useOpenRouter = context.ParseResult.GetValueForOption(options.OpenRouterOption);
            ModelType modelType = isReasoning ? ModelType.Reasoning : fast ? ModelType.Fast : ModelType.Default;

            string? prompt = context.ParseResult.GetValueForArgument(options.PromptArgument);

            await RunAIQuery(new Config
            {
                IsStream = isStream,
                IsCodeBlock = isCodeBlock,
                IsMic = mic,
                UseOpenRouter = useOpenRouter,
                ModelType = modelType
            }.WithUserProfile(), prompt);
        });

        return await rootCommand.InvokeAsync(args);
    }

    static async Task RunAIQuery(Config config, string? argPrompt)
    {
        if ((string.IsNullOrEmpty(config.GroqEndpoint) || string.IsNullOrEmpty(config.GroqToken))
            && (string.IsNullOrEmpty(config.OpenRouterApiKey) || !config.UseOpenRouter))
        {
            Console.Error.WriteLine("Invalid configuration: missing API keys.");
            Environment.Exit(1);
        }

        ILLMClient llmClient;
        try
        {
            llmClient = LLMClientFactory.CreateClient(config);
        }
        catch (Exception ex)
        {
            Console.Error.WriteLine(ex.Message);
            Environment.Exit(1);
            return;
        }

        StringBuilder promptBuilder = new();

        if (Console.IsInputRedirected)
        {
            var input = await Console.In.ReadToEndAsync();
            promptBuilder.AppendLine(input);
        }

        if (argPrompt != null)
        {
            if (promptBuilder.Length > 0)
                promptBuilder.AppendLine("-----");
            promptBuilder.AppendLine(argPrompt);
        }

        if (config.IsMic)
        {
            var micObj = new SpeechToText(config);
            var micInput = await micObj.GetMicInput(useKeyboardInput: !Console.IsInputRedirected);
            if (micInput is null)
                Environment.Exit(0);
            if (promptBuilder.Length > 0)
                promptBuilder.AppendLine("-----");
            promptBuilder.AppendLine(micInput);
        }

        if (promptBuilder.Length == 0)
        {
            Console.Error.WriteLine("No input provided.");
            Environment.Exit(1);
        }

        if (config.IsStream)
        {
            var stream = llmClient.CreateCompletionStreamAsync(promptBuilder.ToString());
            if (config.IsCodeBlock)
            {
                stream = new CodeBlockStreamHandler(stream).Stream();
            }

            using var writer = new StreamWriter(Console.OpenStandardOutput(), new UTF8Encoding(false));
            bool endsWithNewline = false;
            await foreach (var part in stream)
            {
                writer.Write(part);
                endsWithNewline = part.EndsWith("\n");
            }
            if (!endsWithNewline)
                writer.Write(Environment.NewLine);
        }
        else
        {
            string response = await llmClient.CreateCompletionAsync(promptBuilder.ToString());
            if (config.IsCodeBlock)
                response = ExtractCodeBlock(response);

            using var writer = new StreamWriter(Console.OpenStandardOutput(), new UTF8Encoding(false));
            if (!response.EndsWith("\n"))
                response += Environment.NewLine;
            writer.Write(response);

        }
    }

    static string ExtractCodeBlock(string input)
    {
        var m = Regex.Match(input, @"```[a-zA-Z0-9.]*\n([\s\S]+?)\n```", RegexOptions.Compiled);
        return m.Success ? m.Groups[1].Value : input;
    }
}
