using System.Text;
using System.Text.RegularExpressions;
using System.CommandLine;
using System.CommandLine.Invocation;

using aipipe.llms;

namespace aipipe;

class Program
{
    static async Task Main(string[] args)
    {
        var options = new CommandLineOptions();
        var rootCommand = options.RootCommand;

        rootCommand.SetHandler(async (InvocationContext context) =>
        {
            // Config
            var codeBlock = context.ParseResult.GetValueForOption(options.CodeBlockOption);
            var isReasoning = context.ParseResult.GetValueForOption(options.ReasoningOption);
            var fast = context.ParseResult.GetValueForOption(options.FastOption);
            var mic = context.ParseResult.GetValueForOption(options.MicOption);
            var useOpenRouter = context.ParseResult.GetValueForOption(options.OpenRouterOption);
            ModelType modelType = isReasoning ? ModelType.Reasoning : fast ? ModelType.Fast : ModelType.Default;

            // Prompt
            var prompt = context.ParseResult.GetValueForArgument(options.PromptArgument);


            await RunAIQuery(new Config
            {
                IsCodeBlock = codeBlock,
                IsMic = mic,
                UseOpenRouter = useOpenRouter,
                ModelType = modelType
            }, prompt);
        });

        await rootCommand.InvokeAsync(args);
    }

    static async Task RunAIQuery(Config config, string? argPrompt)
    {
        void PrintErrorAndExit(string message)
        {
            Console.Error.WriteLine(message);
            Environment.Exit(1);
        }

        if ((string.IsNullOrEmpty(config.GroqEndpoint) || string.IsNullOrEmpty(config.GroqToken)) && (string.IsNullOrEmpty(config.OpenRouterApiKey) || !config.UseOpenRouter))
        {
            PrintErrorAndExit("Must set either GROQ_ENDPOINT/GROQ_API_KEY or OPENROUTER_API_KEY environment variables and specify --or for OpenRouter.");
        }

        ILLMClient llmClient;
        try
        {
            llmClient = LLMClientFactory.CreateClient(config);
        }
        catch (Exception ex)
        {
            PrintErrorAndExit(ex.Message);
            return;
        }

        // Build prompt
        StringBuilder sb = new();

        // Is there a file being piped to stdin
        bool isFileStream = Console.IsInputRedirected;

        if (isFileStream)
        {
            var input = await Console.In.ReadToEndAsync();
            sb.AppendLine(input);
        }

        if (config.IsMic)
        {
            var mic = new Mic(config);
            var micInput = await mic.GetMicInput();
            if (micInput is null) // user aborted
            {
                Environment.Exit(0);
            }
            sb.AppendLine(micInput);
        }

        if (argPrompt is not null)
        {
            if (sb.Length > 0)
                sb.AppendLine("-----");
            sb.AppendLine(argPrompt);
        }

        if (sb.Length == 0)
        {
            PrintErrorAndExit("Error: Must provide prompt or pipe a file to stdin.");
        }

        string aiOutput = await llmClient.CreateCompletionAsync(sb.ToString());

        if (config.IsCodeBlock)
        {
            aiOutput = ExtractCodeBlock(aiOutput);
        }

        using (var writer = new StreamWriter(Console.OpenStandardOutput(), new UTF8Encoding(false)))
        {
            writer.Write(aiOutput);
        }
    }

    static string ExtractCodeBlock(string input)
    {
        var match = Regex.Match(input, @"```[a-zA-Z0-9.]*\n([\s\S]+?)\n```");
        if (match.Success)
        {
            return match.Groups[1].Value;
        }
        return input;
    }
}
