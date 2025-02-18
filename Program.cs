using System.ComponentModel.DataAnnotations;
using System.Text;
using System.Text.RegularExpressions;
using OpenAI;
using OpenAI.Chat;

namespace aipipe;

class Program
{
    static async Task Main(string[] args)
    {
        void PrintErrorAndExit(string message)
        {
            Console.Error.WriteLine(message);
            Environment.Exit(1);
        }

        if(args.Any(a=> new[] { "--help", "-h", "/?" }.Contains(a)))
        {
            Console.WriteLine("Usage: aipipe [--cb] [--r1] [--fast] [prompt]");
            Console.WriteLine("Options:");
            Console.WriteLine("  --cb     Extract code block from response");
            Console.WriteLine("  --r1     Use the DeepSeek model");
            Console.WriteLine("  --fast   Use the Fast model");
            Console.WriteLine("  --help   Display this help message");
            Environment.Exit(0);
        }

        // Handle flags
        string cbFlagStr = "--cb";
        string deepseek = "--r1";
        string fast = "--fast";
        bool isCodeBlock = false;
        bool isDeepSeek = false;
        bool isFast = false;
        var nonFlagArgs = args.Where(arg => !arg.StartsWith("--")).ToList();
        if(args.Contains(cbFlagStr)) isCodeBlock = true;
        if(args.Contains(deepseek)) isDeepSeek = true;
        if(args.Contains(fast)) isFast = true;

        // Build prompt
        StringBuilder sb = new();

        // Is there a file being piped to stdin
        bool isFileStream = Console.IsInputRedirected;

        if(isFileStream)
        {
            var input = await Console.In.ReadToEndAsync();
            sb.AppendLine(input);
        }
        
        var argPrompt = nonFlagArgs.FirstOrDefault();

        if(argPrompt is not null)
        {
            if(sb.Length > 0)
                sb.AppendLine("-----");
            sb.AppendLine(argPrompt);
        }

        if (sb.Length == 0)
        {
            PrintErrorAndExit("Error: Must provide prompt or pipe a file to stdin.");
        }

        // Run AI Query
        var groqEndpoint = Environment.GetEnvironmentVariable("GROQ_ENDPOINT");
        var groqToken = Environment.GetEnvironmentVariable("GROQ_API_KEY");
        var groqModel = Environment.GetEnvironmentVariable("GROQ_MODEL") ?? "llama-3.3-70b-versatile";

        if (string.IsNullOrEmpty(groqEndpoint))
        {
            PrintErrorAndExit("GROQ_ENDPOINT environment variable not set.");
        }
        if (string.IsNullOrEmpty(groqToken))
        {
            PrintErrorAndExit("GROQ_API_KEY environment variable not set.");
        }

        if(isDeepSeek) groqModel = "deepseek-r1-distill-llama-70b";
        if(isFast) groqModel = "llama-3.1-8b-instant";

        ChatClient client = new(model: groqModel, credential: groqToken!, new OpenAIClientOptions{
            Endpoint = new Uri(groqEndpoint!),
        });

        var systemMessage = isCodeBlock
                ? "You are a helpful assistant. If the user has asked for something written, put it in a code block (```), otherwise just provide the answer."
                 +" If you do use a codeblock, all other text is ignored."
                : "You are a helpful assistant.";

        var options = new ChatCompletionOptions
        {
            
        };
        var messages = new ChatMessage[]
        {
            new SystemChatMessage(systemMessage),
            new UserChatMessage(sb.ToString()),
        };
        var response = await client.CompleteChatAsync(messages, options, CancellationToken.None);
        var aiOutput = response.Value.Content.Single().Text;

        if(isCodeBlock)
        {
            aiOutput = ExtractCodeBlock(aiOutput);
        }

        using (var writer = new StreamWriter(Console.OpenStandardOutput(), Encoding.UTF8))
        {
            writer.WriteLine(aiOutput);
        }
    }

    static string ExtractCodeBlock(string input)
    {
        var match = Regex.Match(input, @"```[a-zA-Z0-9.]*\n([\s\S]+?)\n```");
        if(match.Success)
        {
            return match.Groups[1].Value;
        }
        return input;
    }
}
